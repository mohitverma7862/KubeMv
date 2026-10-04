package stub

import (
	"context"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *client) GetGitOpsOverview(_ context.Context, namespace string) (kubernetes.GitOpsOverview, error) {
	if namespace == "" {
		namespace = "payments"
	}
	now := time.Now().UTC()
	apps := []kubernetes.GitOpsApplication{
		{
			Name:         "payment-api",
			Namespace:    namespace,
			Provider:     "Argo CD",
			Repository:   "github.com/acme/platform-gitops",
			Revision:     "main@a3f91c2",
			Path:         "clusters/dev/payments/payment-api",
			SyncStatus:   "OutOfSync",
			Health:       "Degraded",
			LastSyncedAt: now.Add(-18 * time.Minute).Format(time.RFC3339),
		},
		{
			Name:         "payment-worker",
			Namespace:    namespace,
			Provider:     "Flux",
			Repository:   "github.com/acme/platform-gitops",
			Revision:     "main@b8120de",
			Path:         "clusters/dev/payments/payment-worker",
			SyncStatus:   "Synced",
			Health:       "Healthy",
			LastSyncedAt: now.Add(-4 * time.Minute).Format(time.RFC3339),
		},
		{
			Name:         "ingress-payments",
			Namespace:    "platform",
			Provider:     "Argo CD",
			Repository:   "github.com/acme/ingress-gitops",
			Revision:     "main@cc104aa",
			Path:         "ingress/payments",
			SyncStatus:   "Synced",
			Health:       "Healthy",
			LastSyncedAt: now.Add(-11 * time.Minute).Format(time.RFC3339),
		},
	}
	filtered := make([]kubernetes.GitOpsApplication, 0)
	for _, a := range apps {
		if a.Namespace == namespace {
			filtered = append(filtered, a)
		}
	}
	if len(filtered) == 0 {
		filtered = apps[:2]
	}
	drift := []kubernetes.GitOpsDrift{
		{
			ID:         "deploy-replicas",
			Severity:   "high",
			Resource:   "Deployment/payment-api",
			Field:      "spec.replicas",
			GitValue:   "3",
			LiveValue:  "2",
			Suggestion: "Scale deployment to 3 or revert manual change and sync from Git.",
		},
		{
			ID:         "cm-log-level",
			Severity:   "medium",
			Resource:   "ConfigMap/payment-api-config",
			Field:      "data.LOG_LEVEL",
			GitValue:   "info",
			LiveValue:  "debug",
			Suggestion: "Restore ConfigMap from Git or commit intentional debug change.",
		},
	}
	pipelines := []kubernetes.PipelineRun{
		{
			ID:        "run-1042",
			Name:      "deploy-payments",
			Trigger:   "push main",
			Status:    "success",
			Commit:    "a3f91c2",
			StartedAt: now.Add(-22 * time.Minute).Format(time.RFC3339),
			URL:       "https://ci.example/runs/1042",
		},
		{
			ID:        "run-1043",
			Name:      "deploy-payments",
			Trigger:   "manual sync",
			Status:    "running",
			Commit:    "a3f91c2",
			StartedAt: now.Add(-2 * time.Minute).Format(time.RFC3339),
			URL:       "https://ci.example/runs/1043",
		},
	}
	stats := kubernetes.GitOpsStats{
		Applications: len(filtered),
		Synced:       1,
		OutOfSync:    1,
		DriftItems:   len(drift),
	}
	return kubernetes.GitOpsOverview{
		Namespace:    namespace,
		Stats:        stats,
		Applications: filtered,
		Drift:        drift,
		Pipelines:    pipelines,
	}, nil
}
