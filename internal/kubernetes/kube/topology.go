package kube

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *liveClient) GetTopology(ctx context.Context, query kubemvk8s.TopologyQuery) (kubemvk8s.TopologyGraph, error) {
	ns := query.Namespace
	if ns == "" || ns == "all" {
		ns = ""
	}
	nodes := []kubemvk8s.GraphNode{}
	edges := []kubemvk8s.GraphEdge{}

	deployments, _ := c.clientset.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	for _, d := range deployments.Items {
		id := fmt.Sprintf("deploy/%s", d.Name)
		score := healthScore(d.Status.ReadyReplicas, d.Status.Replicas)
		nodes = append(nodes, kubemvk8s.GraphNode{
			ID: id, Kind: "Deployment", Name: d.Name, Namespace: d.Namespace,
			Status: workloadStatus(d.Status.ReadyReplicas, d.Status.Replicas), HealthScore: score,
		})
	}

	pods, _ := c.clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	for _, p := range pods.Items {
		id := fmt.Sprintf("pod/%s", p.Name)
		score := 100
		if p.Status.Phase != "Running" && p.Status.Phase != "Succeeded" {
			score = 40
		}
		nodes = append(nodes, kubemvk8s.GraphNode{
			ID: id, Kind: "Pod", Name: p.Name, Namespace: p.Namespace,
			Status: string(p.Status.Phase), HealthScore: score,
		})
		for _, owner := range p.OwnerReferences {
			edges = append(edges, kubemvk8s.GraphEdge{
				Source: ownerRefID(owner), Target: id, Relation: "owns",
			})
		}
	}

	services, _ := c.clientset.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
	for _, s := range services.Items {
		id := fmt.Sprintf("svc/%s", s.Name)
		nodes = append(nodes, kubemvk8s.GraphNode{
			ID: id, Kind: "Service", Name: s.Name, Namespace: s.Namespace,
			Status: "Active", HealthScore: 90,
		})
	}

	if query.Mode == kubemvk8s.TopologyNetwork {
		ingresses, _ := c.clientset.NetworkingV1().Ingresses(ns).List(ctx, metav1.ListOptions{})
		for _, ing := range ingresses.Items {
			id := fmt.Sprintf("ing/%s", ing.Name)
			nodes = append(nodes, kubemvk8s.GraphNode{
				ID: id, Kind: "Ingress", Name: ing.Name, Namespace: ing.Namespace,
				Status: "Active", HealthScore: 88,
			})
		}
	}

	return kubemvk8s.TopologyGraph{
		Mode: string(query.Mode), Namespace: ns, Root: query.RootName,
		Nodes: nodes, Edges: edges,
	}, nil
}

func ownerRefID(owner metav1.OwnerReference) string {
	kind := owner.Kind
	switch kind {
	case "ReplicaSet":
		return fmt.Sprintf("rs/%s", owner.Name)
	case "Deployment":
		return fmt.Sprintf("deploy/%s", owner.Name)
	default:
		return fmt.Sprintf("%s/%s", kind, owner.Name)
	}
}

func healthScore(ready, total int32) int {
	if total == 0 {
		return 100
	}
	return int((float64(ready) / float64(total)) * 100)
}

func workloadStatus(ready, total int32) string {
	if ready < total {
		return "Degraded"
	}
	return "Healthy"
}
