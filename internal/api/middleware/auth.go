package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mohitverma7862/KubeMv/internal/auth"
	"github.com/mohitverma7862/KubeMv/internal/httputil"
)

type contextKey string

const principalKey contextKey = "kubemv.principal"

// WithPrincipal stores the authenticated principal on the request context.
func WithPrincipal(ctx context.Context, principal auth.Principal) context.Context {
	return context.WithValue(ctx, principalKey, principal)
}

// PrincipalFromContext returns the principal if present.
func PrincipalFromContext(ctx context.Context) (auth.Principal, bool) {
	p, ok := ctx.Value(principalKey).(auth.Principal)
	return p, ok
}

// RequireAuth validates the session token and attaches the principal.
func RequireAuth(authenticator auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r.Header.Get("Authorization"))
			if token == "" {
				httputil.WriteUnauthorized(w, "missing bearer token")
				return
			}
			principal, err := authenticator.ValidateSession(r.Context(), token)
			if err != nil {
				httputil.WriteUnauthorized(w, "invalid session")
				return
			}
			ctx := WithPrincipal(r.Context(), principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(header string) string {
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
