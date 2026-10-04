package audit

import (
	"context"
	"testing"
)

func TestRecorderRedactsSecrets(t *testing.T) {
	recorder := NewMemoryRecorder(nil)
	recorder.Record(context.Background(), Event{
		Actor:  "admin",
		Action: "cluster.register",
		Result: ResultSuccess,
		Detail: "kubeconfig password=hunter2 token: abc",
	})
	events := recorder.Events()
	if len(events) != 1 {
		t.Fatalf("events %d", len(events))
	}
	if events[0].Detail != "[redacted]" {
		t.Fatalf("detail %q", events[0].Detail)
	}
	if events[0].ID == "" {
		t.Fatal("missing id")
	}
}
