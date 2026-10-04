package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes/portforward"
)

type PortForwardHandler struct {
	Manager *portforward.Manager
}

func (h PortForwardHandler) List(w http.ResponseWriter, r *http.Request) {
	sessions := h.Manager.List(r.Context(), r.PathValue("clusterID"))
	httputil.WriteJSON(w, http.StatusOK, sessions)
}

func (h PortForwardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req kubernetes.PortForwardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteBadRequest(w, "invalid JSON body")
		return
	}
	session, err := h.Manager.Create(r.Context(), r.PathValue("clusterID"), req)
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, session)
}

func (h PortForwardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.Manager.Delete(r.Context(), r.PathValue("id")); err != nil {
		httputil.WriteBadRequest(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "stopped"})
}
