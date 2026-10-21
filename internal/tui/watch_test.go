package tui

import (
	"testing"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestWatchHeaderAndAge(t *testing.T) {
	model := Model{
		Connection:       string(kubernetes.ConnectionReconnecting),
		ReconnectAttempt: 2,
		RetryAfter:       4 * time.Second,
		Cluster:          "prod-cluster",
		Context:          "prod",
		Namespace:        "production",
	}
	if got := model.connectionLabel(); got != "◌ RECONNECTING 2 in 4s" {
		t.Fatal(got)
	}
	model.Connection = string(kubernetes.ConnectionLive)
	model.SyncedAt = time.Now().Add(-3 * time.Second)
	model.ResourceVersion = "42"
	title := model.TableTitle(time.Now())
	if !containsAll(title, "/filter", "◷", "rv 42") {
		t.Fatal(title)
	}
	if !containsAll(model.HeaderLine(time.Now()), "● LIVE") {
		t.Fatal(model.HeaderLine(time.Now()))
	}
	model.Connection = string(kubernetes.ConnectionDisconnected)
	if model.connectionLabel() != "✕ DISCONNECTED" {
		t.Fatal(model.connectionLabel())
	}
}

func TestSingleEventDoesNotRebuildTable(t *testing.T) {
	rows := []kubernetes.ResourceRow{
		{Namespace: "production", Name: "api", Status: "Running"},
		{Namespace: "production", Name: "worker", Status: "Running"},
	}
	columns := []string{"NAME", "STATUS"}
	var frame TableSync
	first := frame.Plan(columns, rows, 0)
	if !first.Full || len(first.Paint) != 2 {
		t.Fatalf("first %#v", first)
	}
	rows[1].Status = "CrashLoopBackOff"
	next := frame.Plan(columns, rows, 0)
	if next.Full || len(next.Paint) != 1 || next.Paint[0] != 1 {
		t.Fatalf("patch %#v", next)
	}
	cursor := frame.Plan(columns, rows, 1)
	if cursor.Full || len(cursor.Paint) != 2 {
		t.Fatalf("cursor %#v", cursor)
	}
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !stringsContains(value, part) {
			return false
		}
	}
	return true
}

func stringsContains(value, part string) bool {
	return len(part) == 0 || (len(value) >= len(part) && (value == part || len(value) > 0 && containsAt(value, part)))
}

func containsAt(value, part string) bool {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
