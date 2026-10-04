package handlers

import (
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

type SecurityHandler struct {
	Connector kubernetes.Connector
}

func (h SecurityHandler) Summary(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	summary, err := client.GetSecuritySummary(r.Context(), r.URL.Query().Get("namespace"))
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, summary)
}
