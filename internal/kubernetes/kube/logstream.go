package kube

import (
	"context"
	"io"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

func podLogOptions(opts kubemvk8s.LogOptions, follow bool) *corev1.PodLogOptions {
	tail := opts.TailLines
	if tail == 0 {
		tail = 500
	}
	return &corev1.PodLogOptions{
		Container:  opts.Container,
		Previous:   opts.Previous,
		Timestamps: opts.Timestamps,
		Follow:     follow || opts.Follow,
		TailLines:  &tail,
	}
}

// StreamPodLogs follows pod logs until the context ends.
func (c *liveClient) StreamPodLogs(ctx context.Context, namespace, podName string, opts kubemvk8s.LogOptions) (io.ReadCloser, error) {
	return c.clientset.CoreV1().Pods(namespace).GetLogs(podName, podLogOptions(opts, true)).Stream(ctx)
}
