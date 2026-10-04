package kube

import (
	"context"
	"fmt"
	"strings"
	"time"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *liveClient) GetGitOpsOverview(ctx context.Context, namespace string) (kubemvk8s.GitOpsOverview, error) {
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}
	apps := make([]kubemvk8s.GitOpsApplication, 0)
	drift := make([]kubemvk8s.GitOpsDrift, 0)

	deploys, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return kubemvk8s.GitOpsOverview{}, err
	}
	for _, d := range deploys.Items {
		if !isGitOpsManaged(d.ObjectMeta) {
			continue
		}
		app := appFromDeployment(d)
		apps = append(apps, app)
		if d.Status.ObservedGeneration < d.Generation {
			drift = append(drift, kubemvk8s.GitOpsDrift{
				ID:         "gen-" + d.Name,
				Severity:   "medium",
				Resource:   "Deployment/" + d.Name,
				Field:      "metadata.generation",
				GitValue:   fmt.Sprintf("observed %d", d.Status.ObservedGeneration),
				LiveValue:  fmt.Sprintf("desired %d", d.Generation),
				Suggestion: "Wait for rollout or investigate controller reconciliation.",
			})
		}
		if d.Spec.Replicas != nil && d.Status.Replicas != *d.Spec.Replicas {
			drift = append(drift, kubemvk8s.GitOpsDrift{
				ID:         "replicas-" + d.Name,
				Severity:   "high",
				Resource:   "Deployment/" + d.Name,
				Field:      "status.replicas",
				GitValue:   fmt.Sprintf("%d desired", *d.Spec.Replicas),
				LiveValue:  fmt.Sprintf("%d ready", d.Status.ReadyReplicas),
				Suggestion: "Replicas differ from spec — check HPA, manual scale, or sync from Git.",
			})
		}
	}

	stats := kubemvk8s.GitOpsStats{Applications: len(apps), DriftItems: len(drift)}
	for _, a := range apps {
		if strings.EqualFold(a.SyncStatus, "Synced") {
			stats.Synced++
		} else {
			stats.OutOfSync++
		}
	}

	return kubemvk8s.GitOpsOverview{
		Namespace:    namespace,
		Stats:        stats,
		Applications: apps,
		Drift:        drift,
		Pipelines:    []kubemvk8s.PipelineRun{},
	}, nil
}

func isGitOpsManaged(meta metav1.ObjectMeta) bool {
	labels := meta.Labels
	if labels == nil {
		return false
	}
	for k := range labels {
		if strings.Contains(k, "argocd") || strings.Contains(k, "flux") || k == "app.kubernetes.io/managed-by" {
			return true
		}
	}
	if meta.Annotations != nil {
		for k := range meta.Annotations {
			if strings.Contains(k, "argocd") || strings.Contains(k, "flux") {
				return true
			}
		}
	}
	return false
}

func appFromDeployment(d appsv1.Deployment) kubemvk8s.GitOpsApplication {
	provider := "GitOps"
	repo := ""
	rev := ""
	path := ""
	sync := "Unknown"
	if d.Annotations != nil {
		if v := d.Annotations["argocd.argoproj.io/sync-status"]; v != "" {
			provider = "Argo CD"
			sync = v
		}
		if v := d.Annotations["fluxcd.io/sync-gc-mark"]; v != "" {
			provider = "Flux"
			sync = "Synced"
		}
		repo = d.Annotations["kubemv.git/repository"]
		rev = d.Annotations["kubemv.git/revision"]
		path = d.Annotations["kubemv.git/path"]
	}
	if sync == "Unknown" {
		if d.Generation == d.Status.ObservedGeneration {
			sync = "Synced"
		} else {
			sync = "OutOfSync"
		}
	}
	health := "Healthy"
	if d.Status.ReadyReplicas < d.Status.Replicas || d.Status.UnavailableReplicas > 0 {
		health = "Degraded"
	}
	return kubemvk8s.GitOpsApplication{
		Name:         d.Name,
		Namespace:    d.Namespace,
		Provider:     provider,
		Repository:   repo,
		Revision:     rev,
		Path:         path,
		SyncStatus:   sync,
		Health:       health,
		LastSyncedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
