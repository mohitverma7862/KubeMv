package handlers

import (
	"net/http"
	"strconv"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

type ResourcesHandler struct {
	Connector kubernetes.Connector
}

func (h ResourcesHandler) List(w http.ResponseWriter, r *http.Request) {
	kind, ok := kubernetes.ValidResourceKind(r.PathValue("kind"))
	if !ok {
		httputil.WriteBadRequest(w, "unknown resource kind")
		return
	}
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "failed to connect to cluster")
		return
	}
	opts := listOptionsFromQuery(r)
	rows, err := client.ListResources(r.Context(), kind, opts)
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, rows)
}

func (h ResourcesHandler) Get(w http.ResponseWriter, r *http.Request) {
	kind, ok := kubernetes.ValidResourceKind(r.PathValue("kind"))
	if !ok {
		httputil.WriteBadRequest(w, "unknown resource kind")
		return
	}
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "failed to connect to cluster")
		return
	}
	namespace := r.PathValue("namespace")
	name := r.PathValue("name")
	if namespace == "_" {
		namespace = ""
	}
	if kind == kubernetes.ResourceNodes || kind == kubernetes.ResourceNamespaces || kind == kubernetes.ResourceCRDs {
		namespace = ""
	}
	detail, err := client.GetResource(r.Context(), kind, namespace, name)
	if err != nil {
		httputil.WriteBadRequest(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, detail)
}

func (h ResourcesHandler) CRDs(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "failed to connect to cluster")
		return
	}
	crds, err := client.ListCRDs(r.Context())
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, crds)
}

func listOptionsFromQuery(r *http.Request) kubernetes.ListOptions {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	return kubernetes.ListOptions{
		Namespace: r.URL.Query().Get("namespace"),
		Search:    r.URL.Query().Get("q"),
		Label:     r.URL.Query().Get("labels"),
		Limit:     limit,
	}
}
