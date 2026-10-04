package dev

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/auth"
)

// Authenticator is a development-only authenticator. Disable in production.
type Authenticator struct {
	store *memoryStore
	ttl   time.Duration
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{
		store: newMemoryStore(),
		ttl:   24 * time.Hour,
	}
}

func (a *Authenticator) Login(_ context.Context, creds auth.Credentials) (auth.Session, error) {
	if creds.Username == "" {
		return auth.Session{}, auth.ErrUnauthenticated
	}
	principal := auth.Principal{
		ID:          "dev:" + creds.Username,
		DisplayName: creds.Username,
		Email:       creds.Username + "@kubemv.local",
		Roles:       []auth.AppRole{auth.RolePlatformAdmin},
	}
	token, err := randomToken()
	if err != nil {
		return auth.Session{}, err
	}
	session := auth.Session{
		Token:     token,
		ExpiresAt: time.Now().Add(a.ttl),
		Principal: principal,
	}
	if err := a.store.Save(context.Background(), session); err != nil {
		return auth.Session{}, err
	}
	return session, nil
}

func (a *Authenticator) ValidateSession(_ context.Context, token string) (auth.Principal, error) {
	session, err := a.store.Get(context.Background(), token)
	if err != nil {
		return auth.Principal{}, auth.ErrUnauthenticated
	}
	if time.Now().After(session.ExpiresAt) {
		_ = a.store.Delete(context.Background(), token)
		return auth.Principal{}, auth.ErrUnauthenticated
	}
	return session.Principal, nil
}

func (a *Authenticator) Logout(_ context.Context, token string) error {
	return a.store.Delete(context.Background(), token)
}

type memoryStore struct {
	mu       sync.RWMutex
	sessions map[string]auth.Session
}

func newMemoryStore() *memoryStore {
	return &memoryStore{sessions: map[string]auth.Session{}}
}

func (m *memoryStore) Save(_ context.Context, session auth.Session) error {
	m.mu.Lock()
	m.sessions[session.Token] = session
	m.mu.Unlock()
	return nil
}

func (m *memoryStore) Get(_ context.Context, token string) (auth.Session, error) {
	m.mu.RLock()
	session, ok := m.sessions[token]
	m.mu.RUnlock()
	if !ok {
		return auth.Session{}, auth.ErrUnauthenticated
	}
	return session, nil
}

func (m *memoryStore) Delete(_ context.Context, token string) error {
	m.mu.Lock()
	delete(m.sessions, token)
	m.mu.Unlock()
	return nil
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
