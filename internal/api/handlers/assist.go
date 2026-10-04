package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

type AssistHandler struct {
	Connector kubernetes.Connector
}

func (h AssistHandler) Bundle(w http.ResponseWriter, r *http.Request) {
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	bundle, err := client.GetAssistBundle(
		r.Context(),
		r.URL.Query().Get("namespace"),
		r.URL.Query().Get("kind"),
		r.URL.Query().Get("name"),
	)
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, bundle)
}

type hookDryRunBody struct {
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
}

func (h AssistHandler) HookDryRun(w http.ResponseWriter, r *http.Request) {
	var body hookDryRunBody
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body)
	}
	if body.Namespace == "" {
		body.Namespace = r.URL.Query().Get("namespace")
	}
	if body.Kind == "" {
		body.Kind = r.URL.Query().Get("kind")
	}
	if body.Name == "" {
		body.Name = r.URL.Query().Get("name")
	}
	client, err := h.Connector.Connect(r.Context(), r.PathValue("clusterID"))
	if err != nil {
		httputil.WriteInternal(w, "cluster connect failed")
		return
	}
	result, err := client.DryRunAssistHook(r.Context(), r.PathValue("hookID"), body.Namespace, body.Kind, body.Name)
	if err != nil {
		httputil.WriteInternal(w, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, result)
}
