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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func podRow(pod corev1.Pod) kubemvk8s.ResourceRow {
	restarts := 0
	ready := 0
	status := string(pod.Status.Phase)
	for _, cs := range pod.Status.ContainerStatuses {
		restarts += int(cs.RestartCount)
		if cs.Ready {
			ready++
		}
		if cs.State.Waiting != nil && cs.State.Waiting.Reason != "" {
			status = cs.State.Waiting.Reason
		}
	}
	containers := make([]kubemvk8s.ContainerRef, 0, len(pod.Spec.Containers))
	var ports []string
	for _, container := range pod.Spec.Containers {
		containers = append(containers, kubemvk8s.ContainerRef{Name: container.Name, Image: container.Image})
		for _, port := range container.Ports {
			ports = append(ports, fmt.Sprintf("%d", port.ContainerPort))
		}
	}
	extra := map[string]string{
		"restarts": fmt.Sprintf("%d", restarts),
		"ready":    fmt.Sprintf("%d/%d", ready, len(pod.Spec.Containers)),
		"node":     pod.Spec.NodeName,
		"ip":       pod.Status.PodIP,
	}
	if qos := string(pod.Status.QOSClass); qos != "" {
		extra["qos"] = qos
	}
	if cpu := sumResource(pod, true, corev1.ResourceCPU); cpu != "" {
		extra["cpu"] = cpu
	}
	if mem := sumResource(pod, true, corev1.ResourceMemory); mem != "" {
		extra["mem"] = mem
	}
	limits := limitSummary(pod)
	if limits != "" {
		extra["limits"] = limits
	}
	if len(ports) > 0 {
		extra["ports"] = strings.Join(ports, ",")
	}
	if mounts := mountSummary(pod); mounts != "" {
		extra["mounts"] = mounts
	}
	return kubemvk8s.ResourceRow{
		Kind: "Pod", Namespace: pod.Namespace, Name: pod.Name,
		Status: status, Age: age(pod.CreationTimestamp.Time),
		Labels: pod.Labels, Extra: extra,
		Owners: ownerRefs(pod.OwnerReferences, pod.Namespace), Containers: containers,
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
		Extra:  map[string]string{"ready": fmt.Sprintf("%d/%d", item.Status.ReadyReplicas, item.Status.Replicas)},
		Owners: ownerRefs(item.OwnerReferences, item.Namespace),
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
		Extra:  map[string]string{"ready": fmt.Sprintf("%d/%d", item.Status.ReadyReplicas, item.Status.Replicas)},
		Owners: ownerRefs(item.OwnerReferences, item.Namespace),
	}
}

func daemonSetRow(item appsv1.DaemonSet) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "DaemonSet", Namespace: item.Namespace, Name: item.Name, Status: "Healthy", Age: age(item.CreationTimestamp.Time),
		Extra:  map[string]string{"ready": fmt.Sprintf("%d/%d", item.Status.NumberReady, item.Status.DesiredNumberScheduled)},
		Owners: ownerRefs(item.OwnerReferences, item.Namespace),
	}
}

func jobRow(item batchv1.Job) kubemvk8s.ResourceRow {
	status := "Running"
	if item.Status.Succeeded > 0 {
		status = "Complete"
	}
	return kubemvk8s.ResourceRow{
		Kind: "Job", Namespace: item.Namespace, Name: item.Name, Status: status, Age: age(item.CreationTimestamp.Time),
		Extra:  map[string]string{"succeeded": fmt.Sprintf("%d", item.Status.Succeeded)},
		Owners: ownerRefs(item.OwnerReferences, item.Namespace),
	}
}

func cronJobRow(item batchv1.CronJob) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "CronJob", Namespace: item.Namespace, Name: item.Name, Status: "Active", Age: age(item.CreationTimestamp.Time),
		Extra:  map[string]string{"schedule": item.Spec.Schedule},
		Owners: ownerRefs(item.OwnerReferences, item.Namespace),
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
	case *appsv1.ReplicaSet:
		return replicaSetRow(*item), true
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

func replicaSetRow(item appsv1.ReplicaSet) kubemvk8s.ResourceRow {
	return kubemvk8s.ResourceRow{
		Kind: "ReplicaSet", Namespace: item.Namespace, Name: item.Name, Status: "Active",
		Age: age(item.CreationTimestamp.Time), Labels: item.Labels,
		Owners: ownerRefs(item.OwnerReferences, item.Namespace),
	}
}

func ownerRefs(refs []metav1.OwnerReference, namespace string) []kubemvk8s.OwnerRef {
	controller := make([]kubemvk8s.OwnerRef, 0, 1)
	rest := make([]kubemvk8s.OwnerRef, 0, len(refs))
	for _, ref := range refs {
		item := kubemvk8s.OwnerRef{Kind: ref.Kind, Name: ref.Name, Namespace: namespace}
		if ref.Controller != nil && *ref.Controller {
			controller = append(controller, item)
			continue
		}
		rest = append(rest, item)
	}
	return append(controller, rest...)
}

func sumResource(pod corev1.Pod, requests bool, name corev1.ResourceName) string {
	total := resource.Quantity{}
	found := false
	for _, container := range pod.Spec.Containers {
		list := container.Resources.Limits
		if requests {
			list = container.Resources.Requests
		}
		if quantity, ok := list[name]; ok {
			total.Add(quantity)
			found = true
		}
	}
	if !found {
		return ""
	}
	return total.String()
}

func limitSummary(pod corev1.Pod) string {
	cpu := sumResource(pod, false, corev1.ResourceCPU)
	mem := sumResource(pod, false, corev1.ResourceMemory)
	switch {
	case cpu != "" && mem != "":
		return "cpu=" + cpu + " mem=" + mem
	case cpu != "":
		return "cpu=" + cpu
	case mem != "":
		return "mem=" + mem
	default:
		return ""
	}
}

func mountSummary(pod corev1.Pod) string {
	parts := make([]string, 0, len(pod.Spec.Volumes))
	for _, volume := range pod.Spec.Volumes {
		switch {
		case volume.Secret != nil:
			parts = append(parts, volume.Name+"=secret:"+volume.Secret.SecretName)
		case volume.ConfigMap != nil:
			parts = append(parts, volume.Name+"=configmap:"+volume.ConfigMap.Name)
		default:
			parts = append(parts, volume.Name)
		}
	}
	return strings.Join(parts, ",")
}
