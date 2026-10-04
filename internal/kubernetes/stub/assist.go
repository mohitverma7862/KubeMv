package stub

import (
	"context"

	"github.com/mohitverma7862/KubeMv/internal/ai"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *client) GetAssistBundle(_ context.Context, namespace, kind, name string) (kubernetes.AssistBundle, error) {
	if namespace == "" {
		namespace = "payments"
	}
	if kind == "" {
		kind = "Deployment"
	}
	if name == "" {
		name = "payment-api"
	}
	signals := ai.Signals{
		Namespace:   namespace,
		Kind:        kind,
		Name:        name,
		Status:      "CrashLoopBackOff",
		Restarts:    "17",
		RecentEvent: "Back-off restarting failed container",
		ReadyRatio:  "2/3",
		LogSnippet:  "Error: connect ECONNREFUSED postgres.payments.svc:5432",
	}
	if kind == "Pod" {
		signals.Name = name
	}
	return ai.BuildBundle(signals), nil
}

func (c *client) DryRunAssistHook(_ context.Context, hookID, namespace, kind, name string) (kubernetes.HookDryRunResult, error) {
	return ai.DryRunHook(hookID, ai.Signals{Namespace: namespace, Kind: kind, Name: name}), nil
}
