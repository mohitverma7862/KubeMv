package stub

import (
	"context"
	"fmt"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *client) GetRolloutStatus(_ context.Context, kind kubernetes.ResourceKindPath, namespace, name string) (kubernetes.RolloutStatus, error) {
	return kubernetes.RolloutStatus{
		Kind: string(kind), Name: name, Namespace: namespace,
		Replicas: 3, Ready: 2, Updated: 2, Available: 2, Strategy: "RollingUpdate",
		Revisions: []kubernetes.RolloutRevision{
			{Revision: 41, Replicas: 3, Ready: 3, Progress: 100, Status: "Complete"},
			{Revision: 42, Replicas: 3, Ready: 2, Progress: 67, Status: "Progressing"},
		},
	}, nil
}

func (c *client) ScaleWorkload(_ context.Context, kind kubernetes.ResourceKindPath, namespace, name string, replicas int32) (kubernetes.MutationResult, error) {
	return kubernetes.MutationResult{
		Action:  "scale",
		Risk:    "LOW",
		Status:  "SUCCESS",
		Message: fmt.Sprintf("scaled %s/%s/%s to %d replicas (stub)", kind, namespace, name, replicas),
		AuditID: "stub-audit",
	}, nil
}

func (c *client) RestartWorkload(_ context.Context, kind kubernetes.ResourceKindPath, namespace, name string) (kubernetes.MutationResult, error) {
	return kubernetes.MutationResult{
		Action:  "restart",
		Risk:    "LOW",
		Status:  "SUCCESS",
		Message: fmt.Sprintf("restarted %s/%s/%s (stub)", kind, namespace, name),
		AuditID: "stub-audit",
	}, nil
}

func (c *client) RollbackDeployment(_ context.Context, namespace, name string, revision int64) (kubernetes.MutationResult, error) {
	return kubernetes.MutationResult{
		Action:  "rollback",
		Risk:    "MEDIUM",
		Status:  "SUCCESS",
		Message: fmt.Sprintf("rolled back deployment %s/%s to revision %d (stub)", namespace, name, revision),
		AuditID: "stub-audit",
	}, nil
}
