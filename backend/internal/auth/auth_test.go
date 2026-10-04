package auth

import (
	"context"
	"testing"
	"time"
)

func TestRoleCatalog(t *testing.T) {
	want := []Role{
		RoleSuperAdmin,
		RolePlatformAdmin,
		RoleClusterAdmin,
		RoleSRE,
		RoleDevOps,
		RoleDeveloper,
		RoleSecurity,
		RoleViewer,
		RoleAuditor,
	}
	got := Catalog()
	if len(got) != len(want) {
		t.Fatalf("catalog length %d", len(got))
	}
	for i := range want {
		if got[i] != want[i] || !KnownRole(want[i]) {
			t.Fatalf("role %s missing", want[i])
		}
	}
	if KnownRole("OWNER") {
		t.Fatal("unknown role accepted")
	}
}

func TestLocalAuthenticator(t *testing.T) {
	authn, err := NewLocalAuthenticator("admin", "foundation-test-password", []Role{RolePlatformAdmin})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	principal, err := authn.Authenticate(ctx, "admin", "foundation-test-password")
	if err != nil {
		t.Fatal(err)
	}
	if principal.Username != "admin" || principal.Roles[0] != RolePlatformAdmin {
		t.Fatalf("principal %+v", principal)
	}
	if _, err := authn.Authenticate(ctx, "admin", "wrong-password-value"); err != ErrInvalidCredentials {
		t.Fatalf("wrong password err %v", err)
	}
	if _, err := authn.Authenticate(ctx, "other", "foundation-test-password"); err != ErrInvalidCredentials {
		t.Fatalf("wrong user err %v", err)
	}
}

func TestSessionExpiry(t *testing.T) {
	now := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	store := NewMemorySessionStore(time.Hour, func() time.Time { return now })
	session, err := store.Create(context.Background(), Principal{Username: "admin", Roles: []Role{RoleViewer}})
	if err != nil {
		t.Fatal(err)
	}
	if session.ID == "" || session.CSRFToken == "" || session.ID == session.CSRFToken {
		t.Fatal("session tokens must be present and distinct")
	}
	got, err := store.Get(context.Background(), session.ID)
	if err != nil || got.Principal.Username != "admin" {
		t.Fatalf("get %+v %v", got, err)
	}
	now = now.Add(2 * time.Hour)
	if _, err := store.Get(context.Background(), session.ID); err != ErrUnauthenticated {
		t.Fatalf("expired session err %v", err)
	}
}

func TestLoginLimiterCountsFailuresOnly(t *testing.T) {
	now := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	limiter := NewLoginLimiter(2, time.Minute, func() time.Time { return now })
	if ok, _ := limiter.Allowed("admin"); !ok {
		t.Fatal("expected budget")
	}
	limiter.RecordFailure("admin")
	limiter.RecordFailure("admin")
	if ok, wait := limiter.Allowed("admin"); ok || wait < time.Second {
		t.Fatalf("ok %v wait %s", ok, wait)
	}
	limiter.Reset("admin")
	if ok, _ := limiter.Allowed("admin"); !ok {
		t.Fatal("reset should restore budget")
	}
}
