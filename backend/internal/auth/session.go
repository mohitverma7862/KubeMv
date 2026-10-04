package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const (
	sessionIDBytes = 32
	csrfBytes      = 32
)

// Session is an opaque server-side login. The identifier travels in an
// HttpOnly cookie. The CSRF token is a separate value the browser must send
// on mutating requests.
type Session struct {
	ID        string
	Principal Principal
	CSRFToken string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionStore interface {
	Create(ctx context.Context, principal Principal) (Session, error)
	Get(ctx context.Context, id string) (Session, error)
	Delete(ctx context.Context, id string) error
}

type MemorySessionStore struct {
	ttl time.Duration
	now func() time.Time

	mu       sync.Mutex
	sessions map[string]Session
}

func NewMemorySessionStore(ttl time.Duration, now func() time.Time) *MemorySessionStore {
	if now == nil {
		now = time.Now
	}
	if ttl <= 0 {
		ttl = 8 * time.Hour
	}
	return &MemorySessionStore{
		ttl:      ttl,
		now:      now,
		sessions: make(map[string]Session),
	}
}

func (s *MemorySessionStore) Create(_ context.Context, principal Principal) (Session, error) {
	id, err := randomToken(sessionIDBytes)
	if err != nil {
		return Session{}, err
	}
	csrf, err := randomToken(csrfBytes)
	if err != nil {
		return Session{}, err
	}
	now := s.now()
	session := Session{
		ID:        id,
		Principal: principal,
		CSRFToken: csrf,
		CreatedAt: now,
		ExpiresAt: now.Add(s.ttl),
	}
	s.mu.Lock()
	s.sessions[id] = session
	s.mu.Unlock()
	return session, nil
}

func (s *MemorySessionStore) Get(_ context.Context, id string) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[id]
	if !ok {
		return Session{}, ErrUnauthenticated
	}
	if !s.now().Before(session.ExpiresAt) {
		delete(s.sessions, id)
		return Session{}, ErrUnauthenticated
	}
	return session, nil
}

func (s *MemorySessionStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
	return nil
}

func randomToken(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
