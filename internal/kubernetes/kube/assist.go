package kube

import (
	"context"
	"fmt"
	"strings"

	"github.com/mohitverma7862/KubeMv/internal/ai"
	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *liveClient) GetAssistBundle(ctx context.Context, namespace, kind, name string) (kubemvk8s.AssistBundle, error) {
	if namespace == "" {
		namespace = metav1.NamespaceDefault
	}
	signals := ai.Signals{Namespace: namespace, Kind: kind, Name: name}

	switch strings.ToLower(kind) {
	case "pod", "pods":
		signals.Kind = "Pod"
		pod, err := c.clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.AssistBundle{}, err
		}
		signals.Status = string(pod.Status.Phase)
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.RestartCount > 0 {
				signals.Restarts = strings.TrimSpace(signals.Restarts + " " + cs.Name + "=" + fmt.Sprint(cs.RestartCount))
			}
			if cs.State.Waiting != nil {
				signals.Status = cs.State.Waiting.Reason
			}
		}
		logs, err := c.GetPodLogs(ctx, namespace, name, kubemvk8s.LogOptions{TailLines: 20})
		if err == nil {
			signals.LogSnippet = logs
		}
	case "deployment", "deployments":
		signals.Kind = "Deployment"
		dep, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return kubemvk8s.AssistBundle{}, err
		}
		signals.ReadyRatio = fmt.Sprintf("%d/%d", dep.Status.ReadyReplicas, dep.Status.Replicas)
		signals.Status = "Progressing"
		if dep.Status.ReadyReplicas < dep.Status.Replicas {
			signals.Status = "Degraded"
		}
	default:
		signals.Kind = kind
	}

	events, _ := c.clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: "involvedObject.name=" + name,
	})
	if events != nil && len(events.Items) > 0 {
		last := events.Items[len(events.Items)-1]
		signals.RecentEvent = last.Reason + ": " + last.Message
	}

	return ai.BuildBundle(signals), nil
}

func (c *liveClient) DryRunAssistHook(_ context.Context, hookID, namespace, kind, name string) (kubemvk8s.HookDryRunResult, error) {
	return ai.DryRunHook(hookID, ai.Signals{Namespace: namespace, Kind: kind, Name: name}), nil
}
