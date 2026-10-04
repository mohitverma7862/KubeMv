package kube

import (
	"context"
	"fmt"
	"strings"
	"time"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	corev1 "k8s.io/api/core/v1"
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
		return nil, fmt.Errorf("unsupported kind %s", kind)
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

func (c *liveClient) listPods(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().Pods(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, pod := range list.Items {
		if !matchSearch(pod.Name, opts.Search) {
			continue
		}
		restarts := 0
		for _, cs := range pod.Status.ContainerStatuses {
			restarts += int(cs.RestartCount)
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "Pod", Namespace: pod.Namespace, Name: pod.Name,
			Status: string(pod.Status.Phase), Age: age(pod.CreationTimestamp.Time),
			Labels: pod.Labels,
			Extra: map[string]string{"restarts": fmt.Sprintf("%d", restarts), "node": pod.Spec.NodeName},
		})
	}
	return out, nil
}

func (c *liveClient) listDeployments(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.AppsV1().Deployments(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, d := range list.Items {
		if !matchSearch(d.Name, opts.Search) {
			continue
		}
		ready := fmt.Sprintf("%d/%d", d.Status.ReadyReplicas, d.Status.Replicas)
		status := "Healthy"
		if d.Status.ReadyReplicas < d.Status.Replicas {
			status = "Degraded"
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "Deployment", Namespace: d.Namespace, Name: d.Name,
			Status: status, Age: age(d.CreationTimestamp.Time), Labels: d.Labels,
			Extra: map[string]string{"ready": ready},
		})
	}
	return out, nil
}

func (c *liveClient) listServices(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().Services(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, svc := range list.Items {
		if !matchSearch(svc.Name, opts.Search) {
			continue
		}
		ports := make([]string, 0, len(svc.Spec.Ports))
		for _, p := range svc.Spec.Ports {
			ports = append(ports, fmt.Sprintf("%d", p.Port))
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "Service", Namespace: svc.Namespace, Name: svc.Name,
			Status: "Active", Age: age(svc.CreationTimestamp.Time), Labels: svc.Labels,
			Extra: map[string]string{"type": string(svc.Spec.Type), "ports": strings.Join(ports, ",")},
		})
	}
	return out, nil
}

func (c *liveClient) listNodes(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	list, err := c.clientset.CoreV1().Nodes().List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, node := range list.Items {
		if !matchSearch(node.Name, opts.Search) {
			continue
		}
		status := "NotReady"
		for _, c := range node.Status.Conditions {
			if c.Type == corev1.NodeReady && c.Status == corev1.ConditionTrue {
				status = "Ready"
			}
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "Node", Name: node.Name, Status: status, Age: age(node.CreationTimestamp.Time), Labels: node.Labels,
		})
	}
	return out, nil
}

func (c *liveClient) listEvents(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().Events(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, ev := range list.Items {
		if !matchSearch(ev.Name, opts.Search) && !matchSearch(ev.InvolvedObject.Name, opts.Search) {
			continue
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "Event", Namespace: ev.Namespace, Name: ev.Name,
			Status: ev.Type, Age: age(ev.LastTimestamp.Time),
			Extra: map[string]string{"reason": ev.Reason, "message": ev.Message, "object": ev.InvolvedObject.Name},
		})
	}
	return out, nil
}

func (c *liveClient) listConfigMaps(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().ConfigMaps(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, cm := range list.Items {
		if !matchSearch(cm.Name, opts.Search) {
			continue
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "ConfigMap", Namespace: cm.Namespace, Name: cm.Name,
			Status: "Active", Age: age(cm.CreationTimestamp.Time), Labels: cm.Labels,
			Extra: map[string]string{"keys": fmt.Sprintf("%d", len(cm.Data)+len(cm.BinaryData))},
		})
	}
	return out, nil
}

func (c *liveClient) listSecrets(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	ns := namespaceOrAll(opts.Namespace)
	list, err := c.clientset.CoreV1().Secrets(ns).List(ctx, listOptions(opts))
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(list.Items))
	for _, secret := range list.Items {
		if !matchSearch(secret.Name, opts.Search) {
			continue
		}
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "Secret", Namespace: secret.Namespace, Name: secret.Name,
			Status: "Active", Age: age(secret.CreationTimestamp.Time), Labels: secret.Labels,
			Extra: map[string]string{"type": string(secret.Type), "keys": fmt.Sprintf("%d", len(secret.Data))},
		})
	}
	return out, nil
}

func (c *liveClient) listNamespacesRows(ctx context.Context, opts kubemvk8s.ListOptions) ([]kubemvk8s.ResourceRow, error) {
	namespaces, err := c.ListNamespaces(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(namespaces))
	for _, ns := range namespaces {
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "Namespace", Name: ns.Name, Status: ns.Status, Age: age(ns.CreatedAt),
		})
	}
	return out, nil
}

func (c *liveClient) listCRDRows(ctx context.Context) ([]kubemvk8s.ResourceRow, error) {
	crds, err := c.ListCRDs(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.ResourceRow, 0, len(crds))
	for _, crd := range crds {
		out = append(out, kubemvk8s.ResourceRow{
			Kind: "CustomResourceDefinition", Name: crd.Name, Status: crd.Scope, Age: "-",
			Extra: map[string]string{"group": crd.Group, "version": crd.Version, "kind": crd.Kind},
		})
	}
	return out, nil
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
