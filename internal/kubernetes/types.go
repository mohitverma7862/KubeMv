package kubernetes

import "time"

// ClusterRef identifies a cluster connection managed server-side (never exposed to the browser).
type ClusterRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Namespace is a lightweight view model for API responses.
type Namespace struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// ResourceKind classifies Kubernetes objects for discovery and navigation.
type ResourceKind string

const (
	KindPod        ResourceKind = "Pod"
	KindDeployment ResourceKind = "Deployment"
	KindService    ResourceKind = "Service"
	KindNode       ResourceKind = "Node"
	KindNamespace  ResourceKind = "Namespace"
	KindEvent      ResourceKind = "Event"
	KindCRD        ResourceKind = "CustomResourceDefinition"
)

// ListOptions supports server-side filtering in later phases.
type ListOptions struct {
	Namespace string
	Label     string
	Search    string
	Limit     int
}

// ClusterHealth is a placeholder aggregate used by the overview shell.
type ClusterHealth struct {
	Score      int `json:"score"`
	NodeCount  int `json:"nodeCount"`
	PodCount   int `json:"podCount"`
	Healthy    int `json:"healthy"`
	Warning    int `json:"warning"`
	Failed     int `json:"failed"`
	CPUPercent int `json:"cpuPercent"`
	MemPercent int `json:"memPercent"`
}
