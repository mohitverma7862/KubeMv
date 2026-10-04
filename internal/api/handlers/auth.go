package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/httputil"
	"github.com/mohitverma7862/KubeMv/internal/api/middleware"
	"github.com/mohitverma7862/KubeMv/internal/auth"
)

type AuthHandler struct {
	Authenticator auth.Authenticator
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteBadRequest(w, "invalid JSON body")
		return
	}
	session, err := h.Authenticator.Login(r.Context(), auth.Credentials{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		httputil.WriteUnauthorized(w, "login failed")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token":      session.Token,
		"expiresAt":  session.ExpiresAt,
		"principal":  session.Principal,
	})
}

func (h AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		httputil.WriteUnauthorized(w, "missing principal")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, principal)
}
