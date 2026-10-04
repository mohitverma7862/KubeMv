package auth

import (
	"sync"
	"time"
)

// LoginLimiter slows password guessing for a username. It is in-memory and
// per process, which matches the Phase 0 single-operator API.
type LoginLimiter struct {
	limit  int
	window time.Duration
	now    func() time.Time

	mu   sync.Mutex
	hits map[string][]time.Time
}

func NewLoginLimiter(limit int, window time.Duration, now func() time.Time) *LoginLimiter {
	if limit <= 0 {
		limit = 5
	}
	if window <= 0 {
		window = time.Minute
	}
	if now == nil {
		now = time.Now
	}
	return &LoginLimiter{
		limit:  limit,
		window: window,
		now:    now,
		hits:   make(map[string][]time.Time),
	}
}

// Allowed reports whether a failed-attempt budget remains. It does not record
// the attempt. RecordFailure counts a miss, and Reset clears the budget after
// a successful sign-in.
func (l *LoginLimiter) Allowed(username string) (bool, time.Duration) {
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()
	recent := l.prune(username, now)
	if len(recent) >= l.limit {
		wait := l.window - now.Sub(recent[0])
		if wait < time.Second {
			wait = time.Second
		}
		return false, wait
	}
	return true, 0
}

func (l *LoginLimiter) RecordFailure(username string) {
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()
	recent := l.prune(username, now)
	l.hits[username] = append(recent, now)
}

func (l *LoginLimiter) Reset(username string) {
	l.mu.Lock()
	delete(l.hits, username)
	l.mu.Unlock()
}

func (l *LoginLimiter) prune(username string, now time.Time) []time.Time {
	cutoff := now.Add(-l.window)
	recent := l.hits[username][:0]
	for _, ts := range l.hits[username] {
		if ts.After(cutoff) {
			recent = append(recent, ts)
		}
	}
	l.hits[username] = recent
	return recent
}
