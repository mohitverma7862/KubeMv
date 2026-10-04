package stub

import (
	"context"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// Connector returns deterministic stub data for UI and API development without a live cluster.
type Connector struct{}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) ListClusters(_ context.Context) ([]kubernetes.ClusterRef, error) {
	return []kubernetes.ClusterRef{
		{ID: "local", Name: "Development"},
		{ID: "staging", Name: "Staging"},
		{ID: "production", Name: "Production"},
	}, nil
}

func (c *Connector) Connect(_ context.Context, clusterID string) (kubernetes.ClusterClient, error) {
	name := clusterID
	switch clusterID {
	case "local":
		name = "Development"
	case "staging":
		name = "Staging"
	case "production":
		name = "Production"
	}
	return &client{ref: kubernetes.ClusterRef{ID: clusterID, Name: name}}, nil
}

type client struct {
	ref kubernetes.ClusterRef
}

func (c *client) Cluster() kubernetes.ClusterRef {
	return c.ref
}

func (c *client) ListNamespaces(_ context.Context, _ kubernetes.ListOptions) ([]kubernetes.Namespace, error) {
	now := time.Now().UTC()
	return []kubernetes.Namespace{
		{Name: "default", Status: "Active", CreatedAt: now},
		{Name: "payments", Status: "Active", CreatedAt: now},
		{Name: "platform", Status: "Active", CreatedAt: now},
		{Name: "monitoring", Status: "Active", CreatedAt: now},
	}, nil
}

func (c *client) Health(_ context.Context) (kubernetes.ClusterHealth, error) {
	return kubernetes.ClusterHealth{
		Score:      87,
		NodeCount:  3,
		PodCount:   42,
		Healthy:    39,
		Warning:    2,
		Failed:     1,
		CPUPercent: 64,
		MemPercent: 71,
	}, nil
}
