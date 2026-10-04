package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// LocalAuthenticator is the Phase 0 development identity provider. It holds one
// bootstrap operator whose password is stored only as a bcrypt hash in memory.
type LocalAuthenticator struct {
	username  string
	hash      []byte
	dummy     []byte
	principal Principal
}

func NewLocalAuthenticator(username, password string, roles []Role) (*LocalAuthenticator, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}
	if len(roles) == 0 {
		return nil, fmt.Errorf("at least one role is required")
	}
	for _, role := range roles {
		if !KnownRole(role) {
			return nil, fmt.Errorf("unknown role %q", role)
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	dummySecret := make([]byte, 32)
	if _, err := rand.Read(dummySecret); err != nil {
		return nil, fmt.Errorf("dummy secret: %w", err)
	}
	dummy, err := bcrypt.GenerateFromPassword(dummySecret, bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash dummy: %w", err)
	}
	copied := append([]Role(nil), roles...)
	return &LocalAuthenticator{
		username: username,
		hash:     hash,
		dummy:    dummy,
		principal: Principal{
			Username:    username,
			DisplayName: username,
			Roles:       copied,
		},
	}, nil
}

func (a *LocalAuthenticator) Authenticate(ctx context.Context, username, password string) (Principal, error) {
	if err := ctx.Err(); err != nil {
		return Principal{}, err
	}
	if len(password) == 0 || len(password) > 128 || len(username) > 32 {
		_ = bcrypt.CompareHashAndPassword(a.dummy, []byte("invalid"))
		return Principal{}, ErrInvalidCredentials
	}
	hash := a.dummy
	matched := subtle.ConstantTimeCompare([]byte(username), []byte(a.username)) == 1
	if matched {
		hash = a.hash
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil || !matched {
		return Principal{}, ErrInvalidCredentials
	}
	return a.principal, nil
}
