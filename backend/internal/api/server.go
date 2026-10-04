package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mohitverma7862/KubeMv/backend/internal/audit"
	"github.com/mohitverma7862/KubeMv/backend/internal/auth"
	"github.com/mohitverma7862/KubeMv/backend/internal/cluster"
	"github.com/mohitverma7862/KubeMv/backend/internal/config"
	"github.com/mohitverma7862/KubeMv/backend/internal/plugin"
	"github.com/mohitverma7862/KubeMv/backend/internal/version"
)

const (
	sessionCookie = "kubemv_session"
	csrfHeader    = "X-KubeMv-CSRF"
	maxBodyBytes  = 1 << 16
)

// Server is the Phase 0 HTTP API. It serves identity, cluster metadata, and
// plugin manifests. It has no Kubernetes client.
type Server struct {
	cfg      config.Config
	log      *slog.Logger
	authn    auth.Authenticator
	sessions auth.SessionStore
	limiter  *auth.LoginLimiter
	clusters cluster.Registry
	plugins  *plugin.Registry
	audit    audit.Recorder
	handler  http.Handler
}

func New(cfg config.Config, log *slog.Logger) (*Server, error) {
	if log == nil {
		log = slog.Default()
	}
	authenticator, err := auth.NewLocalAuthenticator(cfg.BootstrapUsername, cfg.BootstrapPassword, []auth.Role{auth.RolePlatformAdmin})
	if err != nil {
		return nil, err
	}
	s := &Server{
		cfg:      cfg,
		log:      log,
		authn:    authenticator,
		sessions: auth.NewMemorySessionStore(cfg.SessionTTL, nil),
		limiter:  auth.NewLoginLimiter(5, time.Minute, nil),
		clusters: cluster.NewMemoryRegistry(nil),
		plugins:  plugin.NewRegistry(),
		audit:    audit.NewMemoryRecorder(nil),
	}
	s.handler = s.routes()
	return s, nil
}

// NewForTest builds a server around injected stores.
func NewForTest(cfg config.Config, authenticator auth.Authenticator, sessions auth.SessionStore, clusters cluster.Registry, plugins *plugin.Registry, recorder audit.Recorder, limiter *auth.LoginLimiter) *Server {
	if plugins == nil {
		plugins = plugin.NewRegistry()
	}
	if recorder == nil {
		recorder = audit.NewMemoryRecorder(nil)
	}
	if limiter == nil {
		limiter = auth.NewLoginLimiter(5, time.Minute, nil)
	}
	s := &Server{
		cfg:      cfg,
		log:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		authn:    authenticator,
		sessions: sessions,
		limiter:  limiter,
		clusters: clusters,
		plugins:  plugins,
		audit:    recorder,
	}
	s.handler = s.routes()
	return s
}

func (s *Server) Handler() http.Handler { return s.handler }

func (s *Server) Audit() audit.Recorder { return s.audit }

func (s *Server) Plugins() *plugin.Registry { return s.plugins }

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", s.health)
	mux.HandleFunc("POST /api/v1/auth/login", s.login)
	mux.HandleFunc("POST /api/v1/auth/logout", s.logout)
	mux.HandleFunc("GET /api/v1/auth/session", s.session)
	mux.HandleFunc("GET /api/v1/clusters", s.listClusters)
	mux.HandleFunc("POST /api/v1/clusters", s.registerCluster)
	mux.HandleFunc("GET /api/v1/clusters/{id}", s.getCluster)
	mux.HandleFunc("DELETE /api/v1/clusters/{id}", s.removeCluster)
	mux.HandleFunc("GET /api/v1/plugins", s.listPlugins)
	return s.withSecurity(s.withCORS(http.MaxBytesHandler(mux, maxBodyBytes)))
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":       "ok",
		"service":      version.Service,
		"version":      version.Version,
		"phase":        version.Phase,
		"capabilities": version.Capabilities,
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	if !requireJSON(w, r) {
		return
	}
	var body loginRequest
	if err := decodeJSON(w, r, &body); err != nil {
		return
	}
	username := strings.TrimSpace(body.Username)
	limitKey := strings.ToLower(username)
	allowed, wait := s.limiter.Allowed(limitKey)
	if !allowed {
		s.audit.Record(r.Context(), audit.Event{
			Actor:  username,
			Action: "auth.login",
			Result: audit.ResultDenied,
		})
		seconds := int(wait.Round(time.Second) / time.Second)
		if seconds < 1 {
			seconds = 1
		}
		w.Header().Set("Retry-After", strconv.Itoa(seconds))
		writeError(w, http.StatusTooManyRequests, "rate_limited", "Too many sign-in attempts. Wait and try again.")
		return
	}
	principal, err := s.authn.Authenticate(r.Context(), username, body.Password)
	if err != nil {
		s.limiter.RecordFailure(limitKey)
		s.audit.Record(r.Context(), audit.Event{
			Actor:  username,
			Action: "auth.login",
			Result: audit.ResultFailure,
		})
		s.log.Info("sign-in failed", "username", username)
		writeError(w, http.StatusUnauthorized, "unauthorized", "Invalid username or password.")
		return
	}
	s.limiter.Reset(limitKey)
	session, err := s.sessions.Create(r.Context(), principal)
	if err != nil {
		s.log.Error("create session")
		writeError(w, http.StatusInternalServerError, "internal", "Could not start a session.")
		return
	}
	s.setSessionCookie(w, session)
	s.audit.Record(r.Context(), audit.Event{
		Actor:  principal.Username,
		Action: "auth.login",
		Result: audit.ResultSuccess,
	})
	writeJSON(w, http.StatusOK, sessionView(session))
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r, true)
	if !ok {
		return
	}
	_ = s.sessions.Delete(r.Context(), session.ID)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   s.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
	s.audit.Record(r.Context(), audit.Event{
		Actor:  session.Principal.Username,
		Action: "auth.logout",
		Result: audit.ResultSuccess,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) session(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r, false)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, sessionView(session))
}

type registerClusterRequest struct {
	Name          string `json:"name"`
	Provider      string `json:"provider"`
	Context       string `json:"context"`
	KubeconfigRef string `json:"kubeconfigRef"`
}

func (s *Server) listClusters(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r, false); !ok {
		return
	}
	clusters, err := s.clusters.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "Could not list clusters.")
		return
	}
	if clusters == nil {
		clusters = []cluster.Cluster{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"clusters": clusters})
}

func (s *Server) registerCluster(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r, true)
	if !ok {
		return
	}
	if !requireJSON(w, r) {
		return
	}
	var body registerClusterRequest
	if err := decodeJSON(w, r, &body); err != nil {
		return
	}
	created, err := s.clusters.Register(r.Context(), cluster.Registration{
		Name:          strings.TrimSpace(body.Name),
		Provider:      cluster.Provider(strings.TrimSpace(body.Provider)),
		Context:       strings.TrimSpace(body.Context),
		KubeconfigRef: strings.TrimSpace(body.KubeconfigRef),
	})
	if err != nil {
		result := audit.ResultFailure
		status, code, message := http.StatusBadRequest, "invalid_request", "Cluster registration is invalid. Use a name, context, and provider. kubeconfigRef must be an opaque server-side reference."
		switch {
		case errors.Is(err, cluster.ErrConflict):
			status, code, message = http.StatusConflict, "conflict", "A cluster with that name is already registered."
		case errors.Is(err, cluster.ErrLimit):
			status, code, message = http.StatusTooManyRequests, "rate_limited", "The cluster registry is full."
		default:
			result = audit.ResultDenied
		}
		s.audit.Record(r.Context(), audit.Event{
			Actor:      session.Principal.Username,
			Action:     "cluster.register",
			Result:     result,
			TargetType: "cluster",
		})
		writeError(w, status, code, message)
		return
	}
	s.audit.Record(r.Context(), audit.Event{
		Actor:      session.Principal.Username,
		Action:     "cluster.register",
		Result:     audit.ResultSuccess,
		TargetType: "cluster",
		TargetID:   created.ID,
		Detail:     created.Name,
	})
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) getCluster(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r, false); !ok {
		return
	}
	id := r.PathValue("id")
	if !validClusterID(id) {
		writeError(w, http.StatusNotFound, "not_found", "Cluster not found.")
		return
	}
	item, err := s.clusters.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Cluster not found.")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) removeCluster(w http.ResponseWriter, r *http.Request) {
	session, ok := s.requireSession(w, r, true)
	if !ok {
		return
	}
	id := r.PathValue("id")
	if !validClusterID(id) {
		writeError(w, http.StatusNotFound, "not_found", "Cluster not found.")
		return
	}
	if err := s.clusters.Remove(r.Context(), id); err != nil {
		s.audit.Record(r.Context(), audit.Event{
			Actor:      session.Principal.Username,
			Action:     "cluster.remove",
			Result:     audit.ResultFailure,
			TargetType: "cluster",
			TargetID:   id,
		})
		writeError(w, http.StatusNotFound, "not_found", "Cluster not found.")
		return
	}
	s.audit.Record(r.Context(), audit.Event{
		Actor:      session.Principal.Username,
		Action:     "cluster.remove",
		Result:     audit.ResultSuccess,
		TargetType: "cluster",
		TargetID:   id,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listPlugins(w http.ResponseWriter, r *http.Request) {
	if _, ok := s.requireSession(w, r, false); !ok {
		return
	}
	manifests := s.plugins.List()
	views := make([]pluginView, 0, len(manifests))
	for _, manifest := range manifests {
		permissions := make([]string, 0, len(manifest.Permissions))
		for _, permission := range manifest.Permissions {
			permissions = append(permissions, string(permission))
		}
		views = append(views, pluginView{
			ID:          manifest.ID,
			Name:        manifest.Name,
			Version:     manifest.Version,
			Category:    string(manifest.Category),
			Description: manifest.Description,
			Permissions: permissions,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"plugins": views})
}

type pluginView struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type sessionResponse struct {
	Principal principalView `json:"principal"`
	ExpiresAt time.Time     `json:"expiresAt"`
	CSRFToken string        `json:"csrfToken"`
}

type principalView struct {
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName"`
	Roles       []string `json:"roles"`
}

func sessionView(session auth.Session) sessionResponse {
	roles := make([]string, 0, len(session.Principal.Roles))
	for _, role := range session.Principal.Roles {
		roles = append(roles, string(role))
	}
	return sessionResponse{
		Principal: principalView{
			Username:    session.Principal.Username,
			DisplayName: session.Principal.DisplayName,
			Roles:       roles,
		},
		ExpiresAt: session.ExpiresAt,
		CSRFToken: session.CSRFToken,
	}
}

func (s *Server) requireSession(w http.ResponseWriter, r *http.Request, mutating bool) (auth.Session, bool) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil || cookie.Value == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Sign in required.")
		return auth.Session{}, false
	}
	session, err := s.sessions.Get(r.Context(), cookie.Value)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Sign in required.")
		return auth.Session{}, false
	}
	if mutating {
		if !requireJSON(w, r) && r.Method != http.MethodDelete {
			return auth.Session{}, false
		}
		header := r.Header.Get(csrfHeader)
		if subtle.ConstantTimeCompare([]byte(header), []byte(session.CSRFToken)) != 1 {
			s.audit.Record(r.Context(), audit.Event{
				Actor:  session.Principal.Username,
				Action: "auth.csrf",
				Result: audit.ResultDenied,
			})
			writeError(w, http.StatusForbidden, "forbidden", "The request failed the CSRF check.")
			return auth.Session{}, false
		}
	}
	return session, true
}

func (s *Server) setSessionCookie(w http.ResponseWriter, session auth.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.cfg.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt,
	})
}

func requireJSON(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		return true
	}
	contentType := r.Header.Get("Content-Type")
	media := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	if media != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "invalid_request", "Content-Type must be application/json.")
		return false
	}
	return true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Request body is not valid JSON.")
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_request", "Request body must contain one JSON object.")
		return errors.New("trailing json")
	}
	return nil
}

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorBody{Error: errorDetail{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(true)
	_ = encoder.Encode(value)
}

func validClusterID(id string) bool {
	if len(id) < 8 || len(id) > 80 {
		return false
	}
	for _, r := range id {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '_':
		default:
			return false
		}
	}
	return true
}

func (s *Server) withSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(s.cfg.AllowedOrigins))
	for _, origin := range s.cfg.AllowedOrigins {
		allowed[origin] = struct{}{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; !ok {
				if r.Method == http.MethodOptions {
					writeError(w, http.StatusForbidden, "forbidden", "Origin is not allowed.")
					return
				}
			} else {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+csrfHeader)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
				w.Header().Set("Vary", "Origin")
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Shutdown helper keeps the process entrypoint small.
func Shutdown(ctx context.Context, server *http.Server) error {
	return server.Shutdown(ctx)
}
