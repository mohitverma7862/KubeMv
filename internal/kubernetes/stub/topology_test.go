package stub

import (
	"context"
	"testing"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestStubTopologyWorkloadGraph(t *testing.T) {
	c := &client{}
	graph, err := c.GetTopology(context.Background(), kubernetes.TopologyQuery{
		Mode: kubernetes.TopologyWorkload, Namespace: "payments", RootName: "payment-api",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(graph.Nodes) < 5 {
		t.Fatalf("expected rich graph, got %d nodes", len(graph.Nodes))
	}
	if len(graph.Edges) < 4 {
		t.Fatalf("expected edges, got %d", len(graph.Edges))
	}
}
