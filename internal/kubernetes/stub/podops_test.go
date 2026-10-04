package stub

import (
	"context"
	"strings"
	"testing"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestStubPodLogsSearch(t *testing.T) {
	c := &client{ref: kubernetes.ClusterRef{ID: "local", Name: "dev"}}
	logs, err := c.GetPodLogs(context.Background(), "payments", "payment-api", kubernetes.LogOptions{
		Container: "app",
		Search:    "oom",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(logs), "oom") {
		t.Fatalf("expected OOM line, got %q", logs)
	}
}
