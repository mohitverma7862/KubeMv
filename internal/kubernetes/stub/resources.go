package stub

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *client) ListResources(_ context.Context, kind kubernetes.ResourceKindPath, opts kubernetes.ListOptions) ([]kubernetes.ResourceRow, error) {
	rows := stubRows(kind, opts.Namespace)
	return filterRows(rows, opts), nil
}

func (c *client) GetResource(_ context.Context, kind kubernetes.ResourceKindPath, namespace, name string) (kubernetes.ResourceDetail, error) {
	rows, _ := c.ListResources(context.Background(), kind, kubernetes.ListOptions{Namespace: namespace, Search: name})
	var row kubernetes.ResourceRow
	for _, r := range rows {
		if r.Name == name && (namespace == "" || r.Namespace == namespace) {
			row = r
			break
		}
	}
	if row.Name == "" {
		row = kubernetes.ResourceRow{Kind: string(kind), Namespace: namespace, Name: name, Status: "NotFound"}
	}
	return kubernetes.ResourceDetail{
		Row: row,
		YAML: fmt.Sprintf("apiVersion: v1\nkind: %s\nmetadata:\n  name: %s\n  namespace: %s\n# stub data\n", titleKind(kind), name, namespace),
		Events: []kubernetes.EventRow{
			{Type: "Warning", Reason: "BackOff", Message: "Back-off restarting failed container", Object: name, Age: "2m", Namespace: namespace},
		},
	}, nil
}

func (c *client) ListCRDs(_ context.Context) ([]kubernetes.CRDInfo, error) {
	return []kubernetes.CRDInfo{
		{Name: "certificates.cert-manager.io", Group: "cert-manager.io", Version: "v1", Kind: "Certificate", Scope: "Namespaced", Namespaced: true},
	}, nil
}

func stubRows(kind kubernetes.ResourceKindPath, ns string) []kubernetes.ResourceRow {
	ns = defaultNS(ns)
	switch kind {
	case kubernetes.ResourcePods:
		return []kubernetes.ResourceRow{
			{Kind: "Pod", Namespace: ns, Name: "payment-api-7d9c8", Status: "CrashLoopBackOff", Age: "17m", Labels: map[string]string{"app": "payments"}, Extra: map[string]string{"restarts": "17", "node": "node-1"}},
			{Kind: "Pod", Namespace: ns, Name: "payment-api-7d9c9", Status: "Running", Age: "2h", Labels: map[string]string{"app": "payments"}, Extra: map[string]string{"restarts": "0", "node": "node-2"}},
			{Kind: "Pod", Namespace: "platform", Name: "api-gateway-0", Status: "Running", Age: "1d", Labels: map[string]string{"app": "gateway"}},
		}
	case kubernetes.ResourceDeployments:
		return []kubernetes.ResourceRow{
			{Kind: "Deployment", Namespace: ns, Name: "payment-api", Status: "Degraded", Age: "30d", Extra: map[string]string{"ready": "2/3"}},
			{Kind: "Deployment", Namespace: "platform", Name: "api-gateway", Status: "Healthy", Age: "60d", Extra: map[string]string{"ready": "3/3"}},
		}
	case kubernetes.ResourceServices:
		return []kubernetes.ResourceRow{
			{Kind: "Service", Namespace: ns, Name: "payment-api", Status: "Active", Age: "30d", Extra: map[string]string{"type": "ClusterIP", "ports": "8080"}},
		}
	case kubernetes.ResourceNodes:
		return []kubernetes.ResourceRow{
			{Kind: "Node", Name: "node-1", Status: "Ready", Age: "90d", Extra: map[string]string{"cpu": "64%", "memory": "71%"}},
			{Kind: "Node", Name: "node-2", Status: "Ready", Age: "90d"},
		}
	case kubernetes.ResourceEvents:
		return []kubernetes.ResourceRow{
			{Kind: "Event", Namespace: ns, Name: "payment-api-7d9c8.17f2", Status: "Warning", Age: "1m", Extra: map[string]string{"reason": "BackOff", "message": "Back-off restarting failed container"}},
		}
	case kubernetes.ResourceConfigMaps:
		return []kubernetes.ResourceRow{
			{Kind: "ConfigMap", Namespace: ns, Name: "payment-api-config", Status: "Active", Age: "20d", Extra: map[string]string{"keys": "4"}},
		}
	case kubernetes.ResourceSecrets:
		return []kubernetes.ResourceRow{
			{Kind: "Secret", Namespace: ns, Name: "payment-api-secrets", Status: "Active", Age: "20d", Extra: map[string]string{"type": "Opaque", "keys": "3"}},
		}
	case kubernetes.ResourceNamespaces:
		now := time.Now().UTC()
		_ = now
		return []kubernetes.ResourceRow{
			{Kind: "Namespace", Name: "default", Status: "Active", Age: "200d"},
			{Kind: "Namespace", Name: "payments", Status: "Active", Age: "120d"},
			{Kind: "Namespace", Name: "platform", Status: "Active", Age: "120d"},
		}
	default:
		return nil
	}
}

func defaultNS(ns string) string {
	if ns == "" || ns == "*" {
		return "payments"
	}
	return ns
}

func filterRows(rows []kubernetes.ResourceRow, opts kubernetes.ListOptions) []kubernetes.ResourceRow {
	var out []kubernetes.ResourceRow
	for _, row := range rows {
		if opts.Namespace != "" && opts.Namespace != "*" && row.Namespace != "" && row.Namespace != opts.Namespace {
			continue
		}
		if opts.Search != "" && !strings.Contains(strings.ToLower(row.Name), strings.ToLower(opts.Search)) {
			continue
		}
		if opts.Label != "" {
			parts := strings.SplitN(opts.Label, "=", 2)
			if len(parts) != 2 || row.Labels[parts[0]] != parts[1] {
				continue
			}
		}
		out = append(out, row)
	}
	return out
}

func titleKind(kind kubernetes.ResourceKindPath) string {
	switch kind {
	case kubernetes.ResourcePods:
		return "Pod"
	case kubernetes.ResourceDeployments:
		return "Deployment"
	default:
		return string(kind)
	}
}
