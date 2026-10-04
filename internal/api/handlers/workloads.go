package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mohitverma7862/KubeMv/internal/api/middleware"
	"github.com/mohitverma7862/KubeMv/internal/audit"
	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

type WorkloadsHandler struct {
	Connector kubernetes.Connector
}

func (h WorkloadsHandler) Rollout(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	kind, ok := kubernetes.ValidResourceKind(r.PathValue("kind"))
	if !ok {
		httputil.WriteBadRequest(w, "invalid kind")
		return
	}
	status, err := client.GetRolloutStatus(r.Context(), kind, r.PathValue("namespace"), r.PathValue("name"))
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, status)
}

type scaleRequest struct {
	Replicas int32 `json:"replicas"`
}

func (h WorkloadsHandler) Scale(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		httputil.WriteUnauthorized(w, "missing principal")
		return
	}
	var req scaleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteBadRequest(w, "invalid JSON")
		return
	}
	kind, ok := kubernetes.ValidResourceKind(r.PathValue("kind"))
	if !ok {
		httputil.WriteBadRequest(w, "invalid kind")
		return
	}
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	result, err := client.ScaleWorkload(r.Context(), kind, r.PathValue("namespace"), r.PathValue("name"), req.Replicas)
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	result.AuditID = audit.LogMutation(principal, r.PathValue("clusterID"), "scale", r.PathValue("kind")+"/"+r.PathValue("name"), result.Status)
	httputil.WriteJSON(w, http.StatusOK, result)
}

func (h WorkloadsHandler) Restart(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		httputil.WriteUnauthorized(w, "missing principal")
		return
	}
	kind, ok := kubernetes.ValidResourceKind(r.PathValue("kind"))
	if !ok {
		httputil.WriteBadRequest(w, "invalid kind")
		return
	}
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	result, err := client.RestartWorkload(r.Context(), kind, r.PathValue("namespace"), r.PathValue("name"))
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	result.AuditID = audit.LogMutation(principal, r.PathValue("clusterID"), "restart", r.PathValue("kind")+"/"+r.PathValue("name"), result.Status)
	httputil.WriteJSON(w, http.StatusOK, result)
}

func (h WorkloadsHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		httputil.WriteUnauthorized(w, "missing principal")
		return
	}
	rev, _ := strconv.ParseInt(r.URL.Query().Get("revision"), 10, 64)
	if rev == 0 {
		httputil.WriteBadRequest(w, "revision query required")
		return
	}
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	result, err := client.RollbackDeployment(r.Context(), r.PathValue("namespace"), r.PathValue("name"), rev)
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	result.AuditID = audit.LogMutation(principal, r.PathValue("clusterID"), "rollback", "deployments/"+r.PathValue("name"), result.Status)
	httputil.WriteJSON(w, http.StatusOK, result)
}
