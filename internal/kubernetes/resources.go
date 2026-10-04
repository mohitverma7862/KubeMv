package kubernetes

// ResourceKindPath is the URL segment for a listable resource type.
type ResourceKindPath string

const (
	ResourcePods         ResourceKindPath = "pods"
	ResourceDeployments  ResourceKindPath = "deployments"
	ResourceServices     ResourceKindPath = "services"
	ResourceNodes        ResourceKindPath = "nodes"
	ResourceEvents       ResourceKindPath = "events"
	ResourceConfigMaps   ResourceKindPath = "configmaps"
	ResourceSecrets      ResourceKindPath = "secrets"
	ResourceCRDs         ResourceKindPath = "crds"
	ResourceNamespaces   ResourceKindPath = "namespaces"
)

func ValidResourceKind(kind string) (ResourceKindPath, bool) {
	switch ResourceKindPath(kind) {
	case ResourcePods, ResourceDeployments, ResourceServices, ResourceNodes,
		ResourceEvents, ResourceConfigMaps, ResourceSecrets, ResourceCRDs, ResourceNamespaces:
		return ResourceKindPath(kind), true
	default:
		return "", false
	}
}

// ResourceRow is a normalized list row for tables and fast navigation.
type ResourceRow struct {
	Kind      string            `json:"kind"`
	Namespace string            `json:"namespace"`
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	Age       string            `json:"age"`
	Labels    map[string]string `json:"labels,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

// ResourceDetail is returned for the resource detail pane.
type ResourceDetail struct {
	Row        ResourceRow `json:"row"`
	YAML       string      `json:"yaml"`
	Events     []EventRow  `json:"events,omitempty"`
	Related    []ResourceRow `json:"related,omitempty"`
}

// EventRow is a compact Kubernetes event.
type EventRow struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Object    string `json:"object"`
	Age       string `json:"age"`
	Namespace string `json:"namespace"`
}

// CRDInfo describes a discovered custom resource definition.
type CRDInfo struct {
	Name       string `json:"name"`
	Group      string `json:"group"`
	Version    string `json:"version"`
	Kind       string `json:"kind"`
	Scope      string `json:"scope"`
	Namespaced bool   `json:"namespaced"`
}
