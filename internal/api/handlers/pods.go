package handlers

import (
	"net/http"
	"strconv"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

type PodsHandler struct {
	Connector kubernetes.Connector
}

func (h PodsHandler) Containers(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "failed to connect to cluster")
		return
	}
	containers, err := client.ListPodContainers(r.Context(), r.PathValue("namespace"), r.PathValue("name"))
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, containers)
}

func (h PodsHandler) Logs(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "failed to connect to cluster")
		return
	}
	tail, _ := strconv.ParseInt(r.URL.Query().Get("tail"), 10, 64)
	if tail == 0 {
		tail = 500
	}
	logs, err := client.GetPodLogs(r.Context(), r.PathValue("namespace"), r.PathValue("name"), kubernetes.LogOptions{
		Container: r.URL.Query().Get("container"),
		Previous:  r.URL.Query().Get("previous") == "true",
		TailLines: tail,
		Search:    r.URL.Query().Get("q"),
	})
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"logs": logs})
}
