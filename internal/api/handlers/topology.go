package handlers

import (
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

type TopologyHandler struct {
	Connector kubernetes.Connector
}

func (h TopologyHandler) Get(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	mode := kubernetes.TopologyMode(r.URL.Query().Get("mode"))
	if mode == "" {
		mode = kubernetes.TopologyWorkload
	}
	graph, err := client.GetTopology(r.Context(), kubernetes.TopologyQuery{
		Mode:      mode,
		Namespace: r.URL.Query().Get("namespace"),
		RootKind:  r.URL.Query().Get("rootKind"),
		RootName:  r.URL.Query().Get("rootName"),
		Search:    r.URL.Query().Get("q"),
	})
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, graph)
}
