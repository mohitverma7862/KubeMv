package cluster

import "context"

// Registry stores cluster metadata. It does not open a Kubernetes connection.
type Registry interface {
	List(ctx context.Context) ([]Cluster, error)
	Get(ctx context.Context, id string) (Cluster, error)
	Register(ctx context.Context, in Registration) (Cluster, error)
	Remove(ctx context.Context, id string) error
}
