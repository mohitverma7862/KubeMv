package stub

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *client) ListPodContainers(_ context.Context, _ string, podName string) ([]kubernetes.PodContainer, error) {
	return []kubernetes.PodContainer{
		{Name: "app", Image: "payment-api:def456", Ready: false, Restarts: 17},
		{Name: "sidecar", Image: "envoy:1.29", Ready: true, Restarts: 0},
	}, nil
}

func (c *client) GetPodLogs(_ context.Context, namespace, podName string, opts kubernetes.LogOptions) (string, error) {
	container := opts.Container
	if container == "" {
		container = "app"
	}
	prefix := "current"
	if opts.Previous {
		prefix = "previous"
	}
	now := time.Now().UTC().Format("15:04:05")
	lines := []string{
		fmt.Sprintf("%s INFO  [%s] starting container %s", now, prefix, container),
		fmt.Sprintf("%s INFO  [%s] listening on :8080", now, prefix),
		fmt.Sprintf("%s ERROR [%s] OOMKilled - memory limit exceeded (512Mi)", now, prefix),
		fmt.Sprintf("%s WARN  [%s] Back-off restarting failed container", now, prefix),
		fmt.Sprintf("%s DEBUG [%s] namespace=%s pod=%s", now, prefix, namespace, podName),
	}
	if opts.Search != "" {
		var filtered []string
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), strings.ToLower(opts.Search)) {
				filtered = append(filtered, line)
			}
		}
		lines = filtered
	}
	return strings.Join(lines, "\n"), nil
}
