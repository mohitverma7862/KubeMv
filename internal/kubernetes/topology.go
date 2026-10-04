package kubernetes

// TopologyMode selects graph layout semantics.
type TopologyMode string

const (
	TopologyWorkload   TopologyMode = "workload"
	TopologyNetwork    TopologyMode = "network"
	TopologyDependency TopologyMode = "dependency"
)

// TopologyQuery describes a graph request.
type TopologyQuery struct {
	Mode      TopologyMode
	Namespace string
	RootKind  string
	RootName  string
	Search    string
}

// GraphNode is a vertex in the resource graph.
type GraphNode struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Status      string `json:"status"`
	HealthScore int    `json:"healthScore"`
}

// GraphEdge connects two resources.
type GraphEdge struct {
	Source   string `json:"source"`
	Target   string `json:"target"`
	Relation string `json:"relation"`
}

// TopologyGraph is an XRay-style resource graph.
type TopologyGraph struct {
	Mode      string      `json:"mode"`
	Namespace string      `json:"namespace"`
	Root      string      `json:"root"`
	Nodes     []GraphNode `json:"nodes"`
	Edges     []GraphEdge `json:"edges"`
}
