package kubernetes

import (
	"context"
	"io"
	"sync"
	"time"
)

// ConnectionState is the watch connection shown in the console header.
type ConnectionState string

const (
	ConnectionLive         ConnectionState = "live"
	ConnectionReconnecting ConnectionState = "reconnecting"
	ConnectionDisconnected ConnectionState = "disconnected"
	ConnectionSample       ConnectionState = "sample"
)

const (
	WatchAdded    = "ADDED"
	WatchModified = "MODIFIED"
	WatchDeleted  = "DELETED"
	WatchBookmark = "BOOKMARK"
	WatchError    = "ERROR"
)

// Snapshot is a consistent list and the resourceVersion it was read at.
type Snapshot struct {
	Rows            []ResourceRow
	ResourceVersion string
}

// Delta is one watch event. Bookmark and error events carry no row.
type Delta struct {
	Type            string
	Row             ResourceRow
	ResourceVersion string
	Gone            bool
	Message         string
}

// WatchSource lists and watches one resource kind.
type WatchSource interface {
	List(ctx context.Context) (Snapshot, error)
	Watch(ctx context.Context, resourceVersion string) (WatchStream, error)
}

// WatchStream yields watch events until it closes.
type WatchStream interface {
	Next() (Delta, error)
	Close()
}

// WatchStatus is the latest connection state for the header and the watch-age indicator.
type WatchStatus struct {
	State           ConnectionState
	Attempt         int
	RetryIn         time.Duration
	ResourceVersion string
	SyncedAt        time.Time
	Err             string
	Generation      int
}

// Session keeps a keyed row set and resumes a watch from the last resourceVersion.
type Session struct {
	max     int
	backoff func(attempt int) time.Duration
	sleep   func(ctx context.Context, delay time.Duration) error

	mu       sync.Mutex
	rows     []ResourceRow
	index    map[string]int
	status   WatchStatus
	notify   chan struct{}
	closed   bool
	onStatus func(WatchStatus)
}

// NewSession builds a watch session. maxAttempts is the number of failed lists before DISCONNECTED.
func NewSession(maxAttempts int, backoff func(int) time.Duration, sleep func(context.Context, time.Duration) error) *Session {
	if maxAttempts < 1 {
		maxAttempts = 8
	}
	if backoff == nil {
		backoff = defaultBackoff
	}
	if sleep == nil {
		sleep = func(ctx context.Context, delay time.Duration) error {
			timer := time.NewTimer(delay)
			defer timer.Stop()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				return nil
			}
		}
	}
	return &Session{
		max:     maxAttempts,
		backoff: backoff,
		sleep:   sleep,
		index:   map[string]int{},
		notify:  make(chan struct{}, 1),
		status:  WatchStatus{State: ConnectionReconnecting, Attempt: 1},
	}
}

// Notify wakes when status or rows change. It closes when Run returns.
func (s *Session) Notify() <-chan struct{} { return s.notify }

// Read returns a copy of the rows and the current status.
func (s *Session) Read() (WatchStatus, []ResourceRow) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status, append([]ResourceRow(nil), s.rows...)
}

// Run lists, watches, and reconnects until the context ends or the attempt budget is spent.
func (s *Session) Run(ctx context.Context, source WatchSource) {
	defer s.close()
	attempt := 0
	rv := ""
	var lastErr string
	for {
		if ctx.Err() != nil {
			return
		}
		if attempt > 0 {
			delay := s.backoff(attempt)
			s.publish(WatchStatus{State: ConnectionReconnecting, Attempt: attempt, RetryIn: delay, ResourceVersion: rv, SyncedAt: s.synced(), Err: lastErr, Generation: s.generation()})
			if err := s.sleep(ctx, delay); err != nil {
				return
			}
		}
		snap, err := source.List(ctx)
		if err != nil {
			attempt++
			lastErr = err.Error()
			if attempt >= s.max {
				s.publish(WatchStatus{State: ConnectionDisconnected, Attempt: attempt, ResourceVersion: rv, SyncedAt: s.synced(), Err: lastErr, Generation: s.generation()})
				return
			}
			continue
		}
		attempt = 0
		lastErr = ""
		rv = snap.ResourceVersion
		s.replace(snap.Rows, rv)
		stream, err := source.Watch(ctx, rv)
		if err != nil {
			attempt = 1
			lastErr = err.Error()
			continue
		}
		gone := s.follow(ctx, stream, &rv)
		stream.Close()
		if ctx.Err() != nil {
			return
		}
		if gone {
			rv = ""
		}
		attempt = 1
		if lastErr == "" {
			lastErr = "watch closed"
		}
	}
}

func (s *Session) follow(ctx context.Context, stream WatchStream, rv *string) bool {
	for {
		if ctx.Err() != nil {
			return false
		}
		delta, err := stream.Next()
		if err != nil {
			if err != io.EOF {
				s.setErr(err.Error())
			}
			return false
		}
		switch delta.Type {
		case WatchBookmark:
			if delta.ResourceVersion != "" {
				*rv = delta.ResourceVersion
			}
			s.touch(*rv)
		case WatchError:
			if delta.Message != "" {
				s.setErr(delta.Message)
			}
			return delta.Gone
		case WatchAdded, WatchModified, WatchDeleted:
			if delta.ResourceVersion != "" {
				*rv = delta.ResourceVersion
			}
			s.apply(delta, *rv)
		}
	}
}

func (s *Session) replace(rows []ResourceRow, rv string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows = append([]ResourceRow(nil), rows...)
	s.index = map[string]int{}
	for i, row := range s.rows {
		s.index[rowKey(row)] = i
	}
	s.status.State = ConnectionLive
	s.status.Attempt = 0
	s.status.RetryIn = 0
	s.status.ResourceVersion = rv
	s.status.SyncedAt = time.Now()
	s.status.Err = ""
	s.status.Generation++
	s.signalLocked()
}

func (s *Session) apply(delta Delta, rv string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := rowKey(delta.Row)
	changed := false
	switch delta.Type {
	case WatchAdded, WatchModified:
		if i, ok := s.index[key]; ok {
			if !sameRow(s.rows[i], delta.Row) {
				s.rows[i] = delta.Row
				changed = true
			}
		} else {
			s.index[key] = len(s.rows)
			s.rows = append(s.rows, delta.Row)
			changed = true
		}
	case WatchDeleted:
		i, ok := s.index[key]
		if !ok {
			break
		}
		last := len(s.rows) - 1
		if i != last {
			moved := s.rows[last]
			s.rows[i] = moved
			s.index[rowKey(moved)] = i
		}
		s.rows = s.rows[:last]
		delete(s.index, key)
		changed = true
	}
	if rv != "" {
		s.status.ResourceVersion = rv
	}
	s.status.State = ConnectionLive
	s.status.Attempt = 0
	s.status.RetryIn = 0
	s.status.SyncedAt = time.Now()
	s.status.Err = ""
	if changed {
		s.signalLocked()
	} else if rv != "" {
		s.signalLocked()
	}
}

func (s *Session) touch(rv string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rv != "" {
		s.status.ResourceVersion = rv
	}
	s.status.State = ConnectionLive
	s.status.SyncedAt = time.Now()
	s.signalLocked()
}

func (s *Session) setErr(message string) {
	s.mu.Lock()
	s.status.Err = message
	s.mu.Unlock()
}

func (s *Session) publish(status WatchStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if status.SyncedAt.IsZero() {
		status.SyncedAt = s.status.SyncedAt
	}
	if status.Generation == 0 {
		status.Generation = s.status.Generation
	}
	s.status = status
	s.signalLocked()
}

func (s *Session) synced() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status.SyncedAt
}

func (s *Session) generation() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status.Generation
}

func (s *Session) signalLocked() {
	if s.onStatus != nil {
		s.onStatus(s.status)
	}
	if s.closed {
		return
	}
	select {
	case s.notify <- struct{}{}:
	default:
	}
}

func (s *Session) close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	close(s.notify)
}

func rowKey(row ResourceRow) string {
	return row.Namespace + "/" + row.Name
}

func sameRow(left, right ResourceRow) bool {
	if left.Kind != right.Kind || left.Namespace != right.Namespace || left.Name != right.Name || left.Status != right.Status || left.Age != right.Age {
		return false
	}
	if len(left.Labels) != len(right.Labels) || len(left.Extra) != len(right.Extra) {
		return false
	}
	for key, value := range left.Labels {
		if right.Labels[key] != value {
			return false
		}
	}
	for key, value := range left.Extra {
		if right.Extra[key] != value {
			return false
		}
	}
	return true
}

func defaultBackoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	shift := attempt - 1
	if shift > 4 {
		shift = 4
	}
	return time.Second << shift
}
