package handlers

import (
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
)

// MetaHandler exposes product metadata for the frontend shell.
type MetaHandler struct {
	Phase   string
	Version string
}

func (h MetaHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"name":        "KubeMv",
		"tagline":     "Advanced Kubernetes Operations Platform",
		"phase":       h.Phase,
		"version":     h.Version,
		"fastMode":    true,
		"visualMode":  true,
		"aiEnabled":   false,
		"description": "Phase 0 foundation — architecture, abstractions, and UI shell",
	})
}
