package kubernetes

// LogOptions controls pod log retrieval.
type LogOptions struct {
	Container string
	Previous  bool
	TailLines int64
	Search    string
}

// PodContainer describes a container in a pod.
type PodContainer struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Ready   bool   `json:"ready"`
	Restarts int32 `json:"restarts"`
}

// PortForwardRequest starts a local tunnel through the API server.
type PortForwardRequest struct {
	Namespace  string `json:"namespace"`
	Pod        string `json:"pod"`
	LocalPort  int    `json:"localPort"`
	RemotePort int    `json:"remotePort"`
}

// PortForwardSession is an active port-forward managed server-side.
type PortForwardSession struct {
	ID         string `json:"id"`
	Namespace  string `json:"namespace"`
	Pod        string `json:"pod"`
	LocalPort  int    `json:"localPort"`
	RemotePort int    `json:"remotePort"`
	Status     string `json:"status"`
	URL        string `json:"url"`
}
