package handlers

import (
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
)

type HealthHandler struct {
	Version string
}

func (h HealthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "kubemv-api",
		"version": h.Version,
	})
}
