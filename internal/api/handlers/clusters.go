package handlers

import (
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

type ClustersHandler struct {
	Connector kubernetes.Connector
}

func (h ClustersHandler) List(w http.ResponseWriter, r *http.Request) {
	clusters, err := h.Connector.ListClusters(r.Context())
	if err != nil {
		httputil.WriteInternal(w, "failed to list clusters")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, clusters)
}

func (h ClustersHandler) Overview(w http.ResponseWriter, r *http.Request) {
	clusterID := r.PathValue("clusterID")
	client, err := h.Connector.Connect(r.Context(), clusterID)
	if err != nil {
		httputil.WriteInternal(w, "failed to connect to cluster")
		return
	}
	health, err := client.Health(r.Context())
	if err != nil {
		httputil.WriteInternal(w, "failed to load cluster health")
		return
	}
	namespaces, err := client.ListNamespaces(r.Context(), kubernetes.ListOptions{})
	if err != nil {
		httputil.WriteInternal(w, "failed to list namespaces")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"cluster":    client.Cluster(),
		"health":     health,
		"namespaces": namespaces,
	})
}
