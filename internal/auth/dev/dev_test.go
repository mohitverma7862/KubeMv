package dev

import (
	"context"
	"testing"

	"github.com/mohitverma7862/KubeMv/internal/auth"
)

func TestDevLoginAndValidate(t *testing.T) {
	a := NewAuthenticator()
	session, err := a.Login(context.Background(), auth.Credentials{Username: "sre"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	principal, err := a.ValidateSession(context.Background(), session.Token)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if principal.DisplayName != "sre" {
		t.Fatalf("unexpected principal: %+v", principal)
	}
}
