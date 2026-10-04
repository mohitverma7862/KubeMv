package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mohitverma7862/KubeMv/backend/internal/audit"
	"github.com/mohitverma7862/KubeMv/backend/internal/auth"
	"github.com/mohitverma7862/KubeMv/backend/internal/cluster"
	"github.com/mohitverma7862/KubeMv/backend/internal/config"
	"github.com/mohitverma7862/KubeMv/backend/internal/plugin"
)

type fakeAuthenticator struct{}

func (fakeAuthenticator) Authenticate(_ context.Context, username, password string) (auth.Principal, error) {
	if username == "admin" && password == "foundation-test-password" {
		return auth.Principal{
			Username:    "admin",
			DisplayName: "admin",
			Roles:       []auth.Role{auth.RolePlatformAdmin},
		}, nil
	}
	return auth.Principal{}, auth.ErrInvalidCredentials
}

func testConfig() config.Config {
	return config.Config{
		Addr:              "127.0.0.1:8787",
		AllowedOrigins:    []string{"http://127.0.0.1:1420"},
		BootstrapUsername: "admin",
		BootstrapPassword: "foundation-test-password",
		SessionTTL:        time.Hour,
	}
}

func newTestApp(t *testing.T) *Server {
	t.Helper()
	return NewForTest(
		testConfig(),
		fakeAuthenticator{},
		auth.NewMemorySessionStore(time.Hour, nil),
		cluster.NewMemoryRegistry(nil),
		plugin.NewRegistry(),
		audit.NewMemoryRecorder(nil),
		auth.NewLoginLimiter(3, time.Minute, nil),
	)
}

func TestHealthAndSessionFlow(t *testing.T) {
	app := newTestApp(t)
	health := httptest.NewRecorder()
	app.Handler().ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/api/v1/health", nil))
	if health.Code != http.StatusOK || !strings.Contains(health.Body.String(), `"phase":0`) {
		t.Fatalf("health %d %s", health.Code, health.Body.String())
	}

	denied := httptest.NewRecorder()
	app.Handler().ServeHTTP(denied, httptest.NewRequest(http.MethodGet, "/api/v1/clusters", nil))
	if denied.Code != http.StatusUnauthorized {
		t.Fatalf("clusters without session %d", denied.Code)
	}

	login := doLogin(t, app, "admin", "foundation-test-password", http.StatusOK)
	var session sessionResponse
	if err := json.Unmarshal(login.Body.Bytes(), &session); err != nil {
		t.Fatal(err)
	}
	if session.Principal.Username != "admin" || session.CSRFToken == "" {
		t.Fatalf("session %+v", session)
	}
	cookie := login.Result().Cookies()[0]
	if !cookie.HttpOnly || cookie.Name != sessionCookie {
		t.Fatalf("cookie %+v", cookie)
	}

	clusters := authed(t, app, http.MethodGet, "/api/v1/clusters", cookie, session.CSRFToken, nil)
	if clusters.Code != http.StatusOK || !strings.Contains(clusters.Body.String(), `"clusters":[]`) {
		t.Fatalf("clusters %d %s", clusters.Code, clusters.Body.String())
	}

	logout := authed(t, app, http.MethodPost, "/api/v1/auth/logout", cookie, session.CSRFToken, nil)
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout %d %s", logout.Code, logout.Body.String())
	}
	after := authed(t, app, http.MethodGet, "/api/v1/auth/session", cookie, session.CSRFToken, nil)
	if after.Code != http.StatusUnauthorized {
		t.Fatalf("session after logout %d", after.Code)
	}
}

func TestClusterRegistryDoesNotReturnKubeconfigRef(t *testing.T) {
	app := newTestApp(t)
	login := doLogin(t, app, "admin", "foundation-test-password", http.StatusOK)
	var session sessionResponse
	if err := json.Unmarshal(login.Body.Bytes(), &session); err != nil {
		t.Fatal(err)
	}
	cookie := login.Result().Cookies()[0]
	body := []byte(`{"name":"prod-eks","provider":"eks","context":"prod","kubeconfigRef":"secret://clusters/prod"}`)
	created := authed(t, app, http.MethodPost, "/api/v1/clusters", cookie, session.CSRFToken, body)
	if created.Code != http.StatusCreated {
		t.Fatalf("create %d %s", created.Code, created.Body.String())
	}
	if strings.Contains(created.Body.String(), "secret://") || strings.Contains(strings.ToLower(created.Body.String()), "kubeconfig") {
		t.Fatalf("response leaked ref: %s", created.Body.String())
	}
	var view cluster.Cluster
	if err := json.Unmarshal(created.Body.Bytes(), &view); err != nil {
		t.Fatal(err)
	}
	if view.ConnectionState != cluster.ConnectionNotConnected || view.KubeconfigRef != "" {
		t.Fatalf("view %+v", view)
	}

	rejected := authed(t, app, http.MethodPost, "/api/v1/clusters", cookie, session.CSRFToken, []byte(`{"name":"bad","provider":"eks","context":"prod","kubeconfigRef":"token: plain-text"}`))
	if rejected.Code != http.StatusBadRequest {
		t.Fatalf("reject inline token %d %s", rejected.Code, rejected.Body.String())
	}

	missing := authed(t, app, http.MethodPost, "/api/v1/clusters", cookie, "wrong-csrf-token-value", body)
	if missing.Code != http.StatusForbidden {
		t.Fatalf("csrf %d", missing.Code)
	}
}

func TestLoginDoesNotEchoPassword(t *testing.T) {
	app := newTestApp(t)
	password := "foundation-test-password"
	response := doLogin(t, app, "admin", "not-the-real-password", http.StatusUnauthorized)
	if strings.Contains(response.Body.String(), "not-the-real-password") || strings.Contains(response.Body.String(), password) {
		t.Fatalf("password echoed: %s", response.Body.String())
	}
	events := app.Audit().Events()
	for _, event := range events {
		if strings.Contains(event.Detail, "not-the-real-password") {
			t.Fatal("audit stored password")
		}
	}
}

func TestLoginRateLimit(t *testing.T) {
	app := newTestApp(t)
	for i := 0; i < 3; i++ {
		doLogin(t, app, "admin", "wrong-password-value", http.StatusUnauthorized)
	}
	limited := doLogin(t, app, "admin", "wrong-password-value", http.StatusTooManyRequests)
	if limited.Header().Get("Retry-After") == "" {
		t.Fatal("missing retry-after")
	}
}

func TestPluginListAndPermissionBoundary(t *testing.T) {
	app := newTestApp(t)
	err := app.Plugins().Register(staticPlugin{plugin.Manifest{
		ID:          "foundation-notes",
		Name:        "Foundation Notes",
		Version:     "0.1.0",
		Category:    plugin.CategoryObservability,
		Description: "In-process manifest used to prove the registry API.",
		Permissions: []plugin.Permission{plugin.PermObservabilityRead},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if err := app.Plugins().Register(staticPlugin{plugin.Manifest{
		ID:          "leaky",
		Name:        "Leaky",
		Version:     "0.1.0",
		Category:    plugin.CategorySecurity,
		Description: "Requests a forbidden permission.",
		Permissions: []plugin.Permission{"secret.read"},
	}}); err != plugin.ErrInvalidManifest {
		t.Fatalf("boundary err %v", err)
	}
	login := doLogin(t, app, "admin", "foundation-test-password", http.StatusOK)
	var session sessionResponse
	_ = json.Unmarshal(login.Body.Bytes(), &session)
	response := authed(t, app, http.MethodGet, "/api/v1/plugins", login.Result().Cookies()[0], session.CSRFToken, nil)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "foundation-notes") {
		t.Fatalf("plugins %d %s", response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), "secret.read") {
		t.Fatal("forbidden plugin was listed")
	}
}

func TestOpenAPICoversRoutes(t *testing.T) {
	raw, err := os.ReadFile("../../../api/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	spec := string(raw)
	required := []string{
		"/api/v1/health:",
		"/api/v1/auth/login:",
		"/api/v1/auth/logout:",
		"/api/v1/auth/session:",
		"/api/v1/clusters:",
		"/api/v1/clusters/{id}:",
		"/api/v1/plugins:",
		"kubeconfigRef",
	}
	for _, item := range required {
		if !strings.Contains(spec, item) {
			t.Fatalf("openapi missing %s", item)
		}
	}
}

func TestSecurityHeadersAndUnknownFields(t *testing.T) {
	app := newTestApp(t)
	health := httptest.NewRecorder()
	app.Handler().ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/api/v1/health", nil))
	if health.Header().Get("X-Content-Type-Options") != "nosniff" || health.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatalf("security headers missing: %v", health.Header())
	}
	payload := []byte(`{"username":"admin","password":"foundation-test-password","kubeconfig":"secret"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("unknown field status %d", response.Code)
	}
	if strings.Contains(response.Body.String(), "foundation-test-password") || strings.Contains(response.Body.String(), "kubeconfig") {
		t.Fatalf("error leaked input: %s", response.Body.String())
	}
}

func TestCORSRejectsUnknownOrigin(t *testing.T) {
	app := newTestApp(t)
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	request.Header.Set("Origin", "http://evil.example")
	request.Header.Set("Access-Control-Request-Method", "POST")
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("cors %d", response.Code)
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("reflected unknown origin")
	}
}

type staticPlugin struct{ manifest plugin.Manifest }

func (p staticPlugin) Manifest() plugin.Manifest { return p.manifest }

func doLogin(t *testing.T, app *Server, username, password string, want int) *httptest.ResponseRecorder {
	t.Helper()
	payload, _ := json.Marshal(map[string]string{"username": username, "password": password})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "http://127.0.0.1:1420")
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)
	if response.Code != want {
		t.Fatalf("login status %d body %s", response.Code, response.Body.String())
	}
	return response
}

func authed(t *testing.T, app *Server, method, path string, cookie *http.Cookie, csrf string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	request := httptest.NewRequest(method, path, reader)
	request.AddCookie(cookie)
	request.Header.Set(csrfHeader, csrf)
	request.Header.Set("Origin", "http://127.0.0.1:1420")
	if method != http.MethodGet && method != http.MethodDelete {
		request.Header.Set("Content-Type", "application/json")
	}
	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, request)
	return response
}
