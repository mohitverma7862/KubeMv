package audit

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const maxEvents = 1000

// MemoryRecorder is a bounded in-process audit log.
type MemoryRecorder struct {
	now func() time.Time

	mu     sync.Mutex
	events []Event
}

func NewMemoryRecorder(now func() time.Time) *MemoryRecorder {
	if now == nil {
		now = time.Now
	}
	return &MemoryRecorder{now: now}
}

func (r *MemoryRecorder) Record(_ context.Context, event Event) {
	if event.Time.IsZero() {
		event.Time = r.now().UTC()
	}
	event.Detail = redact(event.Detail)
	if event.ID == "" {
		event.ID = newEventID()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
	if len(r.events) > maxEvents {
		r.events = r.events[len(r.events)-maxEvents:]
	}
}

func (r *MemoryRecorder) Events() []Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Event, len(r.events))
	copy(out, r.events)
	return out
}

func newEventID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "aud_unknown"
	}
	return "aud_" + hex.EncodeToString(buf)
}
