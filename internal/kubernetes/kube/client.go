package kube

import (
	"context"
	"fmt"
	"strings"
	"time"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

type liveClient struct {
	ref       kubemvk8s.ClusterRef
	clientset k8sclientset.Interface
	ext       apiextensionsclientset.Interface
	cfg       *rest.Config
}

func (c *liveClient) Cluster() kubemvk8s.ClusterRef { return c.ref }

func (c *liveClient) ListNamespaces(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.Namespace, error) {
	list, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.Namespace, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, kubemvk8s.Namespace{
			Name:      item.Name,
			Status:    string(item.Status.Phase),
			CreatedAt: item.CreationTimestamp.Time,
		})
	}
	return out, nil
}

func (c *liveClient) Health(ctx context.Context) (kubemvk8s.ClusterHealth, error) {
	nodes, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return kubemvk8s.ClusterHealth{}, err
	}
	pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return kubemvk8s.ClusterHealth{}, err
	}
	healthy, warning, failed := 0, 0, 0
	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case corev1.PodRunning, corev1.PodSucceeded:
			healthy++
		case corev1.PodPending:
			warning++
		default:
			failed++
		}
	}
	score := 100
	if len(pods.Items) > 0 {
		score = int((float64(healthy) / float64(len(pods.Items))) * 100)
	}
	return kubemvk8s.ClusterHealth{
		Score:      score,
		NodeCount:  len(nodes.Items),
		PodCount:   len(pods.Items),
		Healthy:    healthy,
		Warning:    warning,
		Failed:     failed,
		CPUPercent: 0,
		MemPercent: 0,
	}, nil
}

func (c *liveClient) ListResources(ctx context.Context, kind kubemvk8s.ResourceKindPath, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	rows, _, err := c.listVersion(ctx, kind, opts)
	return rows, err
}

func (c *liveClient) listVersion(ctx context.Context, kind kubemvk8s.ResourceKindPath, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	switch kind {
	case kubemvk8s.ResourcePods:
		return c.listPods(ctx, opts)
	case kubemvk8s.ResourceDeployments:
		return c.listDeployments(ctx, opts)
	case kubemvk8s.ResourceServices:
		return c.listServices(ctx, opts)
	case kubemvk8s.ResourceNodes:
		return c.listNodes(ctx, opts)
	case kubemvk8s.ResourceEvents:
		return c.listEvents(ctx, opts)
	case kubemvk8s.ResourceConfigMaps:
		return c.listConfigMaps(ctx, opts)
	case kubemvk8s.ResourceSecrets:
		return c.listSecrets(ctx, opts)
	case kubemvk8s.ResourceNamespaces:
		return c.listNamespacesRows(ctx, opts)
	case kubemvk8s.ResourceCRDs:
		return c.listCRDRows(ctx)
	case kubemvk8s.ResourceStatefulSets:
		return c.listStatefulSets(ctx, opts)
	case kubemvk8s.ResourceDaemonSets:
		return c.listDaemonSets(ctx, opts)
	case kubemvk8s.ResourceJobs:
		return c.listJobs(ctx, opts)
	case kubemvk8s.ResourceCronJobs:
		return c.listCronJobs(ctx, opts)
	case kubemvk8s.ResourceHPA:
		return c.listHPA(ctx, opts)
	case kubemvk8s.ResourcePDB:
		return c.listPDB(ctx, opts)
	default:
		return nil, "", fmt.Errorf("unsupported kind %s", kind)
	}
}

func (c *liveClient) GetResource(ctx context.Context, kind kubemvk8s.ResourceKindPath, namespace, name string) (kubemvk8s.ResourceDetail, error) {
	rows, err := c.ListResources(ctx, kind, kubemvk8s.ListOptions{Namespace: namespace, Search: name})
	if err != nil {
		return kubemvk8s.ResourceDetail{}, err
	}
	var row kubemvk8s.ResourceRow
	for _, r := range rows {
		if r.Name == name {
			row = r
			break
		}
	}
	if row.Name == "" {
		return kubemvk8s.ResourceDetail{}, fmt.Errorf("resource not found")
	}
	yamlText, err := c.fetchYAML(ctx, kind, namespace, name)
	if err != nil {
		yamlText = "# yaml unavailable: " + err.Error()
	}
	events, _ := c.eventsForObject(ctx, namespace, name)
	return kubemvk8s.ResourceDetail{Row: row, YAML: yamlText, Events: events}, nil
}

func (c *liveClient) ListCRDs(ctx context.Context) ([]kubemvk8s.CRDInfo, error) {
	list, err := c.ext.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.CRDInfo, 0, len(list.Items))
	for _, item := range list.Items {
		version := ""
		if len(item.Spec.Versions) > 0 {
			version = item.Spec.Versions[0].Name
		}
		out = append(out, kubemvk8s.CRDInfo{
			Name:       item.Name,
			Group:      item.Spec.Group,
			Version:    version,
			Kind:       item.Spec.Names.Kind,
			Scope:      string(item.Spec.Scope),
			Namespaced: item.Spec.Scope == "Namespaced",
		})
	}
	return out, nil
}

func (c *liveClient) listPods(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().Pods(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, pod := range list.Items {
		if !matchSearch(pod.Name, opts.Search) {
			continue
		}
		out = append(out, podRow(pod))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listDeployments(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.AppsV1().Deployments(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, deploymentRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listServices(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().Services(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, serviceRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listNodes(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.clientset.CoreV1().Nodes().List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, nodeRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listEvents(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().Events(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) && !matchSearch(item.InvolvedObject.Name, opts.Search) {
			continue
		}
		out = append(out, eventRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listConfigMaps(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().ConfigMaps(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, configMapRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listSecrets(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().Secrets(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, secretRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listNamespacesRows(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		if !matchSearch(item.Name, opts.Search) {
			continue
		}
		out = append(out, namespaceRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) listCRDRows(ctx context.Context) ([]kubemvk8s.ResourceRow, string, error) {
	list, err := c.ext.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, "", err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, item := range list.Items {
		out = append(out, crdRow(item))
	}
	return out, list.ResourceVersion, nil
}

func (c *liveClient) fetchYAML(ctx context.Context, kind kubemvk8s.ResourceKindPath, namespace, name string) (string, error) {
	switch kind {
	case kubemvk8s.ResourcePods:
		obj, err := c.clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		return marshalYAML(obj)
	case kubemvk8s.ResourceDeployments:
		obj, err := c.clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		return marshalYAML(obj)
	case kubemvk8s.ResourceSecrets:
		obj, err := c.clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		obj.Data = nil
		obj.StringData = nil
		text, err := marshalYAML(obj)
		if err != nil {
			return "", err
		}
		return "# secret data omitted\n" + text, nil
	case kubemvk8s.ResourceServices:
		obj, err := c.clientset.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		return marshalYAML(obj)
	case kubemvk8s.ResourceConfigMaps:
		obj, err := c.clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		return marshalYAML(obj)
	default:
		return fmt.Sprintf("# detail YAML for %s/%s not implemented in Phase 1\n", kind, name), nil
	}
}

func (c *liveClient) eventsForObject(ctx context.Context, namespace, name string) ([]kubemvk8s.EventRow, error) {
	list, err := c.clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s", name),
	})
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.EventRow, 0, len(list.Items))
	for _, ev := range list.Items {
		out = append(out, kubemvk8s.EventRow{
			Type: ev.Type, Reason: ev.Reason, Message: ev.Message,
			Object: ev.InvolvedObject.Name, Age: age(ev.LastTimestamp.Time), Namespace: ev.Namespace,
		})
	}
	return out, nil
}

func marshalYAML(obj runtime.Object) (string, error) {
	b, err := yaml.Marshal(obj)
	return string(b), err
}

func listOptions(opts kubemvk8s.ListOptions) metav1.ListOptions {
	lo := metav1.ListOptions{}
	if opts.Label != "" {
		lo.LabelSelector = opts.Label
	}
	if opts.Limit > 0 {
		lo.Limit = int64(opts.Limit)
	}
	return lo
}

func namespaceOrAll(ns string) string {
	if ns == "" || ns == "*" {
		return ""
	}
	return ns
}

func matchSearch(name, search string) bool {
	if search == "" {
		return true
	}
	return strings.Contains(strings.ToLower(name), strings.ToLower(search))
}

func age(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	d := time.Since(t)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
