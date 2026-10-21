package kube

import (
	"context"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *liveClient) listStatefulSets(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.clientset.AppsV1().StatefulSets(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, statefulSetRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listDaemonSets(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.clientset.AppsV1().DaemonSets(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, daemonSetRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listJobs(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.clientset.BatchV1().Jobs(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, jobRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listCronJobs(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.clientset.BatchV1().CronJobs(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, cronJobRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listHPA(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.clientset.AutoscalingV2().HorizontalPodAutoscalers(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, hpaRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listPDB(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.clientset.PolicyV1().PodDisruptionBudgets(namespaceOrAll(opts.Namespace)).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, pdbRow(item))
	}
	return out, list.ResourceVersion, nil
}
