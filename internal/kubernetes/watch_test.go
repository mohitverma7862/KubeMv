package kubernetes

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"
)

func TestKeyedUpdateKeepsOtherRows(t *testing.T) {
	session := NewSession(3, func(int) time.Duration { return 0 }, func(context.Context, time.Duration) error { return nil })
	session.replace([]ResourceRow{
		{Namespace: "production", Name: "api", Status: "Running"},
		{Namespace: "production", Name: "worker", Status: "Running"},
	}, "1")
	session.apply(Delta{Type: WatchModified, Row: ResourceRow{Namespace: "production", Name: "worker", Status: "CrashLoopBackOff"}, ResourceVersion: "2"}, "2")
	status, rows := session.Read()
	if status.ResourceVersion != "2" || status.State != ConnectionLive {
		t.Fatalf("status %#v", status)
	}
	if len(rows) != 2 || rows[0].Status != "Running" || rows[1].Status != "CrashLoopBackOff" {
		t.Fatalf("rows %#v", rows)
	}
	before := len(rows)
	session.apply(Delta{Type: WatchDeleted, Row: ResourceRow{Namespace: "production", Name: "api"}, ResourceVersion: "3"}, "3")
	_, rows = session.Read()
	if len(rows) != before-1 || rows[0].Name != "worker" {
		t.Fatalf("after delete %#v", rows)
	}
}

func TestReconnectThenResume(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	source := &scriptSource{
		lists: []listResult{
			{err: errors.New("timeout")},
			{snap: Snapshot{ResourceVersion: "10", Rows: []ResourceRow{{Namespace: "ns", Name: "api", Status: "Running"}}}},
			{snap: Snapshot{ResourceVersion: "20", Rows: []ResourceRow{{Namespace: "ns", Name: "api", Status: "Failed"}}}},
		},
		streams: []WatchStream{
			&scriptStream{deltas: []Delta{{Type: WatchBookmark, ResourceVersion: "11"}}, err: io.EOF},
			&scriptStream{wait: ctx},
		},
	}
	session := NewSession(4, func(int) time.Duration { return time.Millisecond }, func(context.Context, time.Duration) error { return nil })
	var mu sync.Mutex
	var seen []WatchStatus
	session.onStatus = func(status WatchStatus) {
		mu.Lock()
		seen = append(seen, status)
		mu.Unlock()
	}
	done := make(chan struct{})
	go func() {
		session.Run(ctx, source)
		close(done)
	}()
	deadline := time.After(2 * time.Second)
	for {
		mu.Lock()
		got := append([]WatchStatus(nil), seen...)
		mu.Unlock()
		var reconnect, live bool
		for _, status := range got {
			if status.State == ConnectionReconnecting && status.Attempt >= 1 {
				reconnect = true
			}
			if status.State == ConnectionLive && status.ResourceVersion == "20" {
				live = true
			}
		}
		if reconnect && live {
			_, rows := session.Read()
			if len(rows) != 1 || rows[0].Status != "Failed" {
				t.Fatalf("rows %#v", rows)
			}
			break
		}
		select {
		case <-deadline:
			t.Fatalf("statuses %#v", got)
		case <-time.After(5 * time.Millisecond):
		}
	}
	cancel()
	<-done
}

func TestExpiredWatchRestartsWithoutResourceVersion(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	source := &scriptSource{
		lists: []listResult{
			{snap: Snapshot{ResourceVersion: "10", Rows: []ResourceRow{{Name: "api"}}}},
			{snap: Snapshot{ResourceVersion: "12", Rows: []ResourceRow{{Name: "api"}}}},
		},
		streams: []WatchStream{
			&scriptStream{deltas: []Delta{{Type: WatchBookmark, ResourceVersion: "11"}, {Type: WatchError, Gone: true, Message: "expired"}}},
			&scriptStream{wait: ctx},
		},
	}
	session := NewSession(3, func(int) time.Duration { return 0 }, func(context.Context, time.Duration) error { return nil })
	done := make(chan struct{})
	go func() {
		session.Run(ctx, source)
		close(done)
	}()
	deadline := time.After(2 * time.Second)
	var versions []string
	for {
		source.mu.Lock()
		versions = append([]string(nil), source.watchRVs...)
		source.mu.Unlock()
		if len(versions) >= 2 {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("watch resourceVersions %#v", versions)
		case <-time.After(5 * time.Millisecond):
		}
	}
	if versions[0] != "10" || versions[1] != "12" {
		t.Fatalf("watch resourceVersions %#v", versions)
	}
	cancel()
	<-done
}

func TestDisconnectsAfterAttemptBudget(t *testing.T) {
	source := &scriptSource{lists: []listResult{{err: errors.New("refused")}, {err: errors.New("refused")}}}
	session := NewSession(2, func(int) time.Duration { return 0 }, func(context.Context, time.Duration) error { return nil })
	session.Run(context.Background(), source)
	status, _ := session.Read()
	if status.State != ConnectionDisconnected || status.Attempt != 2 || status.Err == "" {
		t.Fatalf("status %#v", status)
	}
}

type listResult struct {
	snap Snapshot
	err  error
}

type scriptSource struct {
	mu       sync.Mutex
	lists    []listResult
	streams  []WatchStream
	listAt   int
	watchAt  int
	watchRVs []string
}

func (s *scriptSource) List(context.Context) (Snapshot, error) {
	if s.listAt >= len(s.lists) {
		return Snapshot{}, errors.New("no further list")
	}
	item := s.lists[s.listAt]
	s.listAt++
	return item.snap, item.err
}

func (s *scriptSource) Watch(_ context.Context, resourceVersion string) (WatchStream, error) {
	s.mu.Lock()
	s.watchRVs = append(s.watchRVs, resourceVersion)
	s.mu.Unlock()
	if s.watchAt >= len(s.streams) {
		return &scriptStream{err: io.EOF}, nil
	}
	stream := s.streams[s.watchAt]
	s.watchAt++
	return stream, nil
}

type scriptStream struct {
	deltas []Delta
	err    error
	at     int
	wait   context.Context
}

func (s *scriptStream) Next() (Delta, error) {
	if s.wait != nil {
		<-s.wait.Done()
		return Delta{}, s.wait.Err()
	}
	if s.at >= len(s.deltas) {
		if s.err == nil {
			return Delta{}, io.EOF
		}
		return Delta{}, s.err
	}
	delta := s.deltas[s.at]
	s.at++
	return delta, nil
}

func (s *scriptStream) Close() {}
