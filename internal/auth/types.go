package auth

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")
)

// Principal represents an authenticated KubeMv user. Kubernetes RBAC is evaluated separately.
type Principal struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Roles       []AppRole `json:"roles"`
}

// AppRole is application-level RBAC (never bypasses Kubernetes RBAC).
type AppRole string

const (
	RoleSuperAdmin    AppRole = "SUPER_ADMIN"
	RolePlatformAdmin AppRole = "PLATFORM_ADMIN"
	RoleClusterAdmin  AppRole = "CLUSTER_ADMIN"
	RoleSRE           AppRole = "SRE"
	RoleDevOps        AppRole = "DEVOPS"
	RoleDeveloper     AppRole = "DEVELOPER"
	RoleViewer        AppRole = "VIEWER"
	RoleSecurity      AppRole = "SECURITY"
	RoleAuditor       AppRole = "AUDITOR"
)

// Session is an opaque authenticated session reference.
type Session struct {
	Token     string
	ExpiresAt time.Time
	Principal Principal
}

// Credentials are presented during login (local dev or future OIDC exchange).
type Credentials struct {
	Username string
	Password string
}

// Authenticator validates credentials and sessions. Implementations: local dev, OIDC, SAML (future).
type Authenticator interface {
	Login(ctx context.Context, creds Credentials) (Session, error)
	ValidateSession(ctx context.Context, token string) (Principal, error)
	Logout(ctx context.Context, token string) error
}

// SessionStore persists session metadata server-side (never store kube credentials here).
type SessionStore interface {
	Save(ctx context.Context, session Session) error
	Get(ctx context.Context, token string) (Session, error)
	Delete(ctx context.Context, token string) error
}
