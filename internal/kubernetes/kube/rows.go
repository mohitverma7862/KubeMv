package kube

import (
	"fmt"
	"strings"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func podRow(pod corev1.Pod) kubemvk8s.ResourceRow {
	restarts := 0
	ready := 0
	for _, cs := range pod.Status.ContainerStatuses {
		restarts += int(cs.RestartCount)
		if cs.Ready {
			ready++
		}
	}
	return kubemvk8s.ResourceRow{
		Kind: "Pod", Namespace: pod.Namespace, Name: pod.Name,
		Status: string(pod.Status.Phase), Age: age(pod.CreationTimestamp.Time),
		Labels: pod.Labels,
		Extra: map[string]string{
			"restarts": fmt.Sprintf("%d", restarts),
			"ready":    fmt.Sprintf("%d/%d", ready, len(pod.Spec.Containers)),
			"node":     pod.Spec.NodeName,
			"ip":       pod.Status.PodIP,
		},
	}
}

func deploymentRow(item appsv1.Deployment) kubemvk8s.ResourceRow {
	status := "Healthy"
	if item.Status.ReadyReplicas < item.Status.Replicas {
		status = "Degraded"
	}
	return kubemvk8s.ResourceRow{
		Kind: "Deployment", Namespace: item.Namespace, Name: item.Name,
		Status: status, Age: age(item.CreationTimestamp.Time), Labels: item.Labels,
		Extra: map[string]string{"ready": fmt.Sprintf("%d/%d", item.Status.ReadyReplicas, item.Status.Replicas)},
	}
}

func serviceRow(item corev1.Service) kubemvk8s.ResourceRow {
	ports := make([]string, 0, len(item.Spec.Ports))
	for _, port := range item.Spec.Ports {
		ports = append(ports, fmt.Sprintf("%d", port.Port))
	}
	return kubemvk8s.ResourceRow{
		Kind: "Service", Namespace: item.Namespace, Name: item.Name,
		Status: "Active", Age: age(item.CreationTimestamp.Time), Labels: item.Labels,
		Extra: map[string]string{"type": string(item.Spec.Type), "ports": strings.Join(ports, ",")},
	}
}

func nodeRow(item corev1.Node) kubemvk8s.ResourceRow {
	status := "NotReady"
	for _, condition := range item.Status.Conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
			status = "Ready"
		}
	}
	return kubemvk8s.ResourceRow{
		Kind: "Node", Name: item.Name, Status: status, Age: age(item.CreationTimestamp.Time), Labels: item.Labels,
	}
}

func eventRow(item corev1.Event) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "Event", Namespace: item.Namespace, Name: item.Name,
		Status: item.Type, Age: age(item.LastTimestamp.Time),
		Extra: map[string]string{"reason": item.Reason, "message": item.Message, "object": item.InvolvedObject.Name},
	}
}

func configMapRow(item corev1.ConfigMap) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "ConfigMap", Namespace: item.Namespace, Name: item.Name,
		Status: "Active", Age: age(item.CreationTimestamp.Time), Labels: item.Labels,
		Extra: map[string]string{"keys": fmt.Sprintf("%d", len(item.Data)+len(item.BinaryData))},
	}
}

func secretRow(item corev1.Secret) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "Secret", Namespace: item.Namespace, Name: item.Name,
		Status: "Active", Age: age(item.CreationTimestamp.Time), Labels: item.Labels,
		Extra: map[string]string{"type": string(item.Type), "keys": fmt.Sprintf("%d", len(item.Data))},
	}
}

func namespaceRow(item corev1.Namespace) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "Namespace", Name: item.Name, Status: string(item.Status.Phase), Age: age(item.CreationTimestamp.Time),
	}
}

func crdRow(item apiextensionsv1.CustomResourceDefinition) kubemvk8s.ResourceRow {
	version := ""
	if len(item.Spec.Versions) > 0 {
		version = item.Spec.Versions[0].Name
	}
	return kubemvk8s.ResourceRow{
		Kind: "CustomResourceDefinition", Name: item.Name, Status: string(item.Spec.Scope), Age: age(item.CreationTimestamp.Time),
		Extra: map[string]string{"group": item.Spec.Group, "version": version, "kind": item.Spec.Names.Kind},
	}
}

func statefulSetRow(item appsv1.StatefulSet) kubemvk8s.ResourceRow {
	status := "Healthy"
	if item.Status.ReadyReplicas < item.Status.Replicas {
		status = "Degraded"
	}
	return kubemvk8s.ResourceRow{
		Kind: "StatefulSet", Namespace: item.Namespace, Name: item.Name, Status: status, Age: age(item.CreationTimestamp.Time),
		Extra: map[string]string{"ready": fmt.Sprintf("%d/%d", item.Status.ReadyReplicas, item.Status.Replicas)},
	}
}

func daemonSetRow(item appsv1.DaemonSet) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "DaemonSet", Namespace: item.Namespace, Name: item.Name, Status: "Healthy", Age: age(item.CreationTimestamp.Time),
		Extra: map[string]string{"ready": fmt.Sprintf("%d/%d", item.Status.NumberReady, item.Status.DesiredNumberScheduled)},
	}
}

func jobRow(item batchv1.Job) kubemvk8s.ResourceRow {
	status := "Running"
	if item.Status.Succeeded > 0 {
		status = "Complete"
	}
	return kubemvk8s.ResourceRow{
		Kind: "Job", Namespace: item.Namespace, Name: item.Name, Status: status, Age: age(item.CreationTimestamp.Time),
		Extra: map[string]string{"succeeded": fmt.Sprintf("%d", item.Status.Succeeded)},
	}
}

func cronJobRow(item batchv1.CronJob) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "CronJob", Namespace: item.Namespace, Name: item.Name, Status: "Active", Age: age(item.CreationTimestamp.Time),
		Extra: map[string]string{"schedule": item.Spec.Schedule},
	}
}

func hpaRow(item autoscalingv2.HorizontalPodAutoscaler) kubemvk8s.ResourceRow {
	min := int32(0)
	if item.Spec.MinReplicas != nil {
		min = *item.Spec.MinReplicas
	}
	return kubemvk8s.ResourceRow{
		Kind: "HorizontalPodAutoscaler", Namespace: item.Namespace, Name: item.Name, Status: "Active", Age: age(item.CreationTimestamp.Time),
		Extra: map[string]string{
			"targets": fmt.Sprintf("%d-%d", min, item.Spec.MaxReplicas),
			"current": fmt.Sprintf("%d", item.Status.CurrentReplicas),
		},
	}
}

func pdbRow(item policyv1.PodDisruptionBudget) kubemvk8s.ResourceRow {
	min := ""
	if item.Spec.MinAvailable != nil {
		min = item.Spec.MinAvailable.String()
	}
	return kubemvk8s.ResourceRow{
		Kind: "PodDisruptionBudget", Namespace: item.Namespace, Name: item.Name, Status: "Healthy", Age: age(item.CreationTimestamp.Time),
		Extra: map[string]string{"minAvailable": min},
	}
}

func rowFromObject(obj runtime.Object) (kubemvk8s.ResourceRow, bool) {
	switch item := obj.(type) {
	case *corev1.Pod:
		return podRow(*item), true
	case *appsv1.Deployment:
		return deploymentRow(*item), true
	case *corev1.Service:
		return serviceRow(*item), true
	case *corev1.Node:
		return nodeRow(*item), true
	case *corev1.Event:
		return eventRow(*item), true
	case *corev1.ConfigMap:
		return configMapRow(*item), true
	case *corev1.Secret:
		return secretRow(*item), true
	case *corev1.Namespace:
		return namespaceRow(*item), true
	case *apiextensionsv1.CustomResourceDefinition:
		return crdRow(*item), true
	case *appsv1.StatefulSet:
		return statefulSetRow(*item), true
	case *appsv1.DaemonSet:
		return daemonSetRow(*item), true
	case *batchv1.Job:
		return jobRow(*item), true
	case *batchv1.CronJob:
		return cronJobRow(*item), true
	case *autoscalingv2.HorizontalPodAutoscaler:
		return hpaRow(*item), true
	case *policyv1.PodDisruptionBudget:
		return pdbRow(*item), true
	default:
		return kubemvk8s.ResourceRow{}, false
	}
}
