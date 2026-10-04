package kubernetes

import "context"

// ClusterClient is the primary abstraction over Kubernetes API access.
// Implementations use client-go; credentials remain on the server.
type ClusterClient interface {
	Cluster() ClusterRef
	ListNamespaces(ctx context.Context, opts ListOptions) ([]Namespace, error)
	Health(ctx context.Context) (ClusterHealth, error)
}

// Connector establishes cluster clients from server-managed credentials.
type Connector interface {
	Connect(ctx context.Context, clusterID string) (ClusterClient, error)
	ListClusters(ctx context.Context) ([]ClusterRef, error)
}

// CredentialManager stores encrypted kubeconfig / cloud identity material (Phase 0: interface only).
type CredentialManager interface {
	RegisterCluster(ctx context.Context, ref ClusterRef, secretRef string) error
	ResolveSecretRef(ctx context.Context, clusterID string) (string, error)
}
