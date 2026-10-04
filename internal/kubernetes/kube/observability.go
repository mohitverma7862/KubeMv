package kube

import (
	"context"
	"fmt"
	"strings"
	"time"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *liveClient) QueryMetrics(_ context.Context, _ kubemvk8s.MetricsQuery) (kubemvk8s.MetricsQueryResult, error) {
	return kubemvk8s.MetricsQueryResult{}, fmt.Errorf("use Prometheus proxy (set KUBEMV_PROMETHEUS_URL) for live metric queries")
}

func (c *liveClient) ListScrapeTargets(ctx context.Context, namespace string) ([]kubemvk8s.ScrapeTarget, error) {
	svcs, err := c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ScrapeTarget, 0)
	for _, svc := range svcs.Items {
		if !hasPrometheusAnnotations(svc.Annotations) {
			continue
		}
		job := svc.Annotations["prometheus.io/job"]
		if job == "" {
			job = svc.Name
		}
		out = append(out, kubemvk8s.ScrapeTarget{
			Namespace:  svc.Namespace,
			Kind:       "Service",
			Name:       svc.Name,
			Job:        job,
			Instance:   fmt.Sprintf("%s.%s.svc", svc.Name, svc.Namespace),
			Up:         true,
			LastScrape: time.Now().UTC().Format(time.RFC3339),
			Labels:     svc.Labels,
		})
	}
	return out, nil
}

func hasPrometheusAnnotations(ann map[string]string) bool {
	if ann == nil {
		return false
	}
	for k := range ann {
		if strings.HasPrefix(k, "prometheus.io/") {
			return true
		}
	}
	return false
}

func (c *liveClient) GetObservabilityDashboard(ctx context.Context, namespace, kind, name string) (kubemvk8s.ObservabilityDashboard, error) {
	targets, err := c.ListScrapeTargets(ctx, namespace)
	if err != nil {
		return kubemvk8s.ObservabilityDashboard{}, err
	}
	filtered := make([]kubemvk8s.ScrapeTarget, 0, len(targets))
	for _, t := range targets {
		if name == "" || strings.Contains(t.Name, name) || strings.Contains(t.Job, name) {
			filtered = append(filtered, t)
		}
	}
	return kubemvk8s.ObservabilityDashboard{
		Namespace:   namespace,
		Kind:        kind,
		Name:        name,
		Presets:     []kubemvk8s.ObservabilityPreset{},
		Targets:     filtered,
		LogDeepLink: "/fast?kind=pods&namespace=" + namespace,
	}, nil
}
