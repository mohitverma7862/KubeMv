package stub

import (
	"context"
	"testing"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestConnectorOverviewData(t *testing.T) {
	c := NewConnector()
	client, err := c.Connect(context.Background(), "local")
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	health, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("health: %v", err)
	}
	if health.Score != 87 {
		t.Fatalf("expected score 87, got %d", health.Score)
	}
	ns, err := client.ListNamespaces(context.Background(), kubernetes.ListOptions{})
	if err != nil {
		t.Fatalf("namespaces: %v", err)
	}
	if len(ns) < 3 {
		t.Fatalf("expected stub namespaces, got %d", len(ns))
	}
}
