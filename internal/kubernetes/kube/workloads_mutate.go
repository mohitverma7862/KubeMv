package kube

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *liveClient) GetRolloutStatus(ctx context.Context, kind kubemvk8s.ResourceKindPath, namespace, name string) (kubemvk8s.RolloutStatus, error) {
	if kind != kubemvk8s.ResourceDeployments {
		return kubemvk8s.RolloutStatus{}, fmt.Errorf("rollout status only supported for deployments in Phase 3")
	}
	dep, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return kubemvk8s.RolloutStatus{}, err
	}
	revisions := []kubemvk8s.RolloutRevision{}
	progress := 0
	if dep.Status.Replicas > 0 {
		progress = int((float64(dep.Status.ReadyReplicas) / float64(dep.Status.Replicas)) * 100)
	}
	status := "Progressing"
	if len(dep.Status.Conditions) > 0 {
		status = string(dep.Status.Conditions[len(dep.Status.Conditions)-1].Type)
	}
	revisions = append(revisions, kubemvk8s.RolloutRevision{
		Revision: dep.Generation,
		Replicas: dep.Status.Replicas,
		Ready:    dep.Status.ReadyReplicas,
		Progress: progress,
		Status:   status,
	})
	return kubemvk8s.RolloutStatus{
		Kind: "Deployment", Name: name, Namespace: namespace,
		Replicas: dep.Status.Replicas, Ready: dep.Status.ReadyReplicas,
		Updated: dep.Status.UpdatedReplicas, Available: dep.Status.AvailableReplicas,
		Strategy: string(dep.Spec.Strategy.Type), Revisions: revisions,
	}, nil
}

func (c *liveClient) ScaleWorkload(ctx context.Context, kind kubemvk8s.ResourceKindPath, namespace, name string, replicas int32) (kubemvk8s.MutationResult, error) {
	switch kind {
	case kubemvk8s.ResourceDeployments:
		scale, err := c.clientset.AppsV1().Deployments(namespace).GetScale(ctx, name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.MutationResult{}, err
		}
		scale.Spec.Replicas = replicas
		_, err = c.clientset.AppsV1().Deployments(namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
		if err != nil {
			return kubemvk8s.MutationResult{}, err
		}
	case kubemvk8s.ResourceStatefulSets:
		scale, err := c.clientset.AppsV1().StatefulSets(namespace).GetScale(ctx, name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.MutationResult{}, err
		}
		scale.Spec.Replicas = replicas
		_, err = c.clientset.AppsV1().StatefulSets(namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
		if err != nil {
			return kubemvk8s.MutationResult{}, err
		}
	default:
		return kubemvk8s.MutationResult{}, fmt.Errorf("scale not supported for %s", kind)
	}
	return kubemvk8s.MutationResult{Action: "scale", Risk: "LOW", Status: "SUCCESS", Message: "scaled successfully"}, nil
}

func (c *liveClient) RestartWorkload(ctx context.Context, kind kubemvk8s.ResourceKindPath, namespace, name string) (kubemvk8s.MutationResult, error) {
	if kind != kubemvk8s.ResourceDeployments {
		return kubemvk8s.MutationResult{}, fmt.Errorf("restart only supported for deployments in Phase 3")
	}
	patch := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))
	_, err := c.clientset.AppsV1().Deployments(namespace).Patch(ctx, name, types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{})
	if err != nil {
		return kubemvk8s.MutationResult{}, err
	}
	return kubemvk8s.MutationResult{Action: "restart", Risk: "LOW", Status: "SUCCESS", Message: "rollout restart triggered"}, nil
}

func (c *liveClient) RollbackDeployment(ctx context.Context, namespace, name string, revision int64) (kubemvk8s.MutationResult, error) {
	dep, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return kubemvk8s.MutationResult{}, err
	}
	if dep.Annotations == nil {
		dep.Annotations = map[string]string{}
	}
	dep.Annotations["kubemv.io/rollback-revision"] = fmt.Sprintf("%d", revision)
	_, err = c.clientset.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
	if err != nil {
		return kubemvk8s.MutationResult{}, err
	}
	return kubemvk8s.MutationResult{Action: "rollback", Risk: "MEDIUM", Status: "SUCCESS", Message: fmt.Sprintf("rollback requested to revision %d", revision)}, nil
}
