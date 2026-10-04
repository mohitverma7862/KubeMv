package stub

import (
	"context"
	"strings"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *client) GetTopology(_ context.Context, query kubernetes.TopologyQuery) (kubernetes.TopologyGraph, error) {
	ns := query.Namespace
	if ns == "" || ns == "all" {
		ns = "payments"
	}
	root := query.RootName
	if root == "" {
		root = "payment-api"
	}

	nodes := []kubernetes.GraphNode{
		{ID: "ing/payment", Kind: "Ingress", Name: "payment", Namespace: ns, Status: "Healthy", HealthScore: 92},
		{ID: "svc/payment-api", Kind: "Service", Name: "payment-api", Namespace: ns, Status: "Active", HealthScore: 88},
		{ID: "deploy/payment-api", Kind: "Deployment", Name: "payment-api", Namespace: ns, Status: "Degraded", HealthScore: 61},
		{ID: "rs/payment-api-42", Kind: "ReplicaSet", Name: "payment-api-42", Namespace: ns, Status: "Progressing", HealthScore: 67},
		{ID: "pod/payment-api-7d9c8", Kind: "Pod", Name: "payment-api-7d9c8", Namespace: ns, Status: "CrashLoopBackOff", HealthScore: 31},
		{ID: "pod/payment-api-7d9c9", Kind: "Pod", Name: "payment-api-7d9c9", Namespace: ns, Status: "Running", HealthScore: 94},
		{ID: "cm/payment-api-config", Kind: "ConfigMap", Name: "payment-api-config", Namespace: ns, Status: "Active", HealthScore: 100},
		{ID: "secret/payment-api-secrets", Kind: "Secret", Name: "payment-api-secrets", Namespace: ns, Status: "Active", HealthScore: 100},
		{ID: "sa/payment-api", Kind: "ServiceAccount", Name: "payment-api", Namespace: ns, Status: "Active", HealthScore: 100},
	}

	if query.Mode == kubernetes.TopologyNetwork {
		nodes = append(nodes,
			kubernetes.GraphNode{ID: "net/internet", Kind: "Internet", Name: "internet", Namespace: "", Status: "Healthy", HealthScore: 100},
			kubernetes.GraphNode{ID: "svc/postgres", Kind: "Service", Name: "postgres", Namespace: ns, Status: "Active", HealthScore: 90},
		)
	}

	if query.Search != "" {
		q := strings.ToLower(query.Search)
		filtered := make([]kubernetes.GraphNode, 0)
		for _, n := range nodes {
			if strings.Contains(strings.ToLower(n.Name), q) || strings.Contains(strings.ToLower(n.Kind), q) {
				filtered = append(filtered, n)
			}
		}
		nodes = filtered
	}

	edges := []kubernetes.GraphEdge{
		{Source: "ing/payment", Target: "svc/payment-api", Relation: "routes"},
		{Source: "svc/payment-api", Target: "deploy/payment-api", Relation: "selects"},
		{Source: "deploy/payment-api", Target: "rs/payment-api-42", Relation: "owns"},
		{Source: "rs/payment-api-42", Target: "pod/payment-api-7d9c8", Relation: "owns"},
		{Source: "rs/payment-api-42", Target: "pod/payment-api-7d9c9", Relation: "owns"},
		{Source: "pod/payment-api-7d9c8", Target: "cm/payment-api-config", Relation: "mounts"},
		{Source: "pod/payment-api-7d9c8", Target: "secret/payment-api-secrets", Relation: "mounts"},
		{Source: "pod/payment-api-7d9c8", Target: "sa/payment-api", Relation: "uses"},
	}

	if query.Mode == kubernetes.TopologyNetwork {
		edges = append(edges,
			kubernetes.GraphEdge{Source: "net/internet", Target: "ing/payment", Relation: "ingress"},
			kubernetes.GraphEdge{Source: "pod/payment-api-7d9c9", Target: "svc/postgres", Relation: "egress"},
		)
	}

	return kubernetes.TopologyGraph{
		Mode: string(query.Mode), Namespace: ns, Root: root,
		Nodes: nodes, Edges: filterEdges(edges, nodes),
	}, nil
}

func filterEdges(edges []kubernetes.GraphEdge, nodes []kubernetes.GraphNode) []kubernetes.GraphEdge {
	ids := map[string]bool{}
	for _, n := range nodes {
		ids[n.ID] = true
	}
	out := make([]kubernetes.GraphEdge, 0, len(edges))
	for _, e := range edges {
		if ids[e.Source] && ids[e.Target] {
			out = append(out, e)
		}
	}
	return out
}
