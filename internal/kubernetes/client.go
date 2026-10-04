package kubernetes

import "context"

// ClusterClient is the primary abstraction over Kubernetes API access.
// Implementations use client-go; credentials remain on the server.
type ClusterClient interface {
	Cluster() ClusterRef
	ListNamespaces(ctx context.Context, opts ListOptions) ([]Namespace, error)
	Health(ctx context.Context) (ClusterHealth, error)
	ListResources(ctx context.Context, kind ResourceKindPath, opts ListOptions) ([]ResourceRow, error)
	GetResource(ctx context.Context, kind ResourceKindPath, namespace, name string) (ResourceDetail, error)
	ListCRDs(ctx context.Context) ([]CRDInfo, error)
	ListPodContainers(ctx context.Context, namespace, podName string) ([]PodContainer, error)
	GetPodLogs(ctx context.Context, namespace, podName string, opts LogOptions) (string, error)
	GetRolloutStatus(ctx context.Context, kind ResourceKindPath, namespace, name string) (RolloutStatus, error)
	ScaleWorkload(ctx context.Context, kind ResourceKindPath, namespace, name string, replicas int32) (MutationResult, error)
	RestartWorkload(ctx context.Context, kind ResourceKindPath, namespace, name string) (MutationResult, error)
	RollbackDeployment(ctx context.Context, namespace, name string, revision int64) (MutationResult, error)
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
