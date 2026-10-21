package kube

import (
	"context"
	"fmt"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OwnerLookup resolves one owner object. The returned row includes that object's own owners.
type OwnerLookup interface {
	LookupOwner(ctx context.Context, ref kubemvk8s.OwnerRef) (kubemvk8s.ResourceRow, error)
}

// LookupOwner reads a single owner so the relationship tree can walk Deployment → ReplicaSet → Pod.
func (c *liveClient) LookupOwner(ctx context.Context, ref kubemvk8s.OwnerRef) (kubemvk8s.ResourceRow, error) {
	switch ref.Kind {
	case "Pod":
		item, err := c.clientset.CoreV1().Pods(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.ResourceRow{}, err
		}
		return podRow(*item), nil
	case "ReplicaSet":
		item, err := c.clientset.AppsV1().ReplicaSets(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.ResourceRow{}, err
		}
		return replicaSetRow(*item), nil
	case "Deployment":
		item, err := c.clientset.AppsV1().Deployments(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.ResourceRow{}, err
		}
		return deploymentRow(*item), nil
	case "StatefulSet":
		item, err := c.clientset.AppsV1().StatefulSets(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.ResourceRow{}, err
		}
		return statefulSetRow(*item), nil
	case "DaemonSet":
		item, err := c.clientset.AppsV1().DaemonSets(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.ResourceRow{}, err
		}
		return daemonSetRow(*item), nil
	case "Job":
		item, err := c.clientset.BatchV1().Jobs(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.ResourceRow{}, err
		}
		return jobRow(*item), nil
	case "CronJob":
		item, err := c.clientset.BatchV1().CronJobs(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.ResourceRow{}, err
		}
		return cronJobRow(*item), nil
	default:
		return kubemvk8s.ResourceRow{}, fmt.Errorf("unsupported owner %s", ref.Kind)
	}
}
