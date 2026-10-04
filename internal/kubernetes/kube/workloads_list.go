package kube

import (
	"context"
	"fmt"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *liveClient) listStatefulSets(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	list, err := c.clientset.AppsV1().StatefulSets(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, sts := range list.Items {
		if !matchSearch(sts.Name, opts.Search) {
			continue
		}
		status := "Healthy"
		if sts.Status.ReadyReplicas < sts.Status.Replicas {
			status = "Degraded"
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "StatefulSet", Namespace: sts.Namespace, Name: sts.Name, Status: status, Age: age(sts.CreationTimestamp.Time),
			Extra: map[string]string{"ready": fmt.Sprintf("%d/%d", sts.Status.ReadyReplicas, sts.Status.Replicas)},
		})
	}
	return out, nil
}

func (c *liveClient) listDaemonSets(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	list, err := c.clientset.AppsV1().DaemonSets(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, ds := range list.Items {
		if !matchSearch(ds.Name, opts.Search) {
			continue
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "DaemonSet", Namespace: ds.Namespace, Name: ds.Name, Status: "Healthy", Age: age(ds.CreationTimestamp.Time),
			Extra: map[string]string{"ready": fmt.Sprintf("%d/%d", ds.Status.NumberReady, ds.Status.DesiredNumberScheduled)},
		})
	}
	return out, nil
}

func (c *liveClient) listJobs(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	list, err := c.clientset.BatchV1().Jobs(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, job := range list.Items {
		if !matchSearch(job.Name, opts.Search) {
			continue
		}
		status := "Running"
		if job.Status.Succeeded > 0 {
			status = "Complete"
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "Job", Namespace: job.Namespace, Name: job.Name, Status: status, Age: age(job.CreationTimestamp.Time),
			Extra: map[string]string{"succeeded": fmt.Sprintf("%d", job.Status.Succeeded)},
		})
	}
	return out, nil
}

func (c *liveClient) listCronJobs(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	list, err := c.clientset.BatchV1().CronJobs(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, cj := range list.Items {
		if !matchSearch(cj.Name, opts.Search) {
			continue
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "CronJob", Namespace: cj.Namespace, Name: cj.Name, Status: "Active", Age: age(cj.CreationTimestamp.Time),
			Extra: map[string]string{"schedule": cj.Spec.Schedule},
		})
	}
	return out, nil
}

func (c *liveClient) listHPA(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	list, err := c.clientset.AutoscalingV2().HorizontalPodAutoscalers(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, hpa := range list.Items {
		if !matchSearch(hpa.Name, opts.Search) {
			continue
		}
		min, max, current := int32(0), int32(0), int32(0)
		if hpa.Spec.MinReplicas != nil {
			min = *hpa.Spec.MinReplicas
		}
		max = hpa.Spec.MaxReplicas
		if hpa.Status.CurrentReplicas > 0 {
			current = hpa.Status.CurrentReplicas
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "HorizontalPodAutoscaler", Namespace: hpa.Namespace, Name: hpa.Name, Status: "Active", Age: age(hpa.CreationTimestamp.Time),
			Extra: map[string]string{"targets": fmt.Sprintf("%d-%d", min, max), "current": fmt.Sprintf("%d", current)},
		})
	}
	return out, nil
}

func (c *liveClient) listPDB(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	list, err := c.clientset.PolicyV1().PodDisruptionBudgets(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, pdb := range list.Items {
		if !matchSearch(pdb.Name, opts.Search) {
			continue
		}
		min := ""
		if pdb.Spec.MinAvailable != nil {
			min = pdb.Spec.MinAvailable.String()
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "PodDisruptionBudget", Namespace: pdb.Namespace, Name: pdb.Name, Status: "Healthy", Age: age(pdb.CreationTimestamp.Time),
			Extra: map[string]string{"minAvailable": min},
		})
	}
	return out, nil
}
