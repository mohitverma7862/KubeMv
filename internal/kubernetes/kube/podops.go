package kube

import (
	"bytes"
	"context"
	"io"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (c *liveClient) ListPodContainers(ctx context.Context, namespace, podName string) ([]kubemvk8s.PodContainer, error) {
	pod, err := c.clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	out := make([]kubemvk8s.PodContainer, 0, len(pod.Spec.Containers))
	for i, container := range pod.Spec.Containers {
		ready := false
		restarts := int32(0)
		if i < len(pod.Status.ContainerStatuses) {
			st := pod.Status.ContainerStatuses[i]
			ready = st.Ready
			restarts = st.RestartCount
		}
		out = append(out, kubemvk8s.PodContainer{
			Name: container.Name, Image: container.Image, Ready: ready, Restarts: restarts,
		})
	}
	return out, nil
}

func (c *liveClient) GetPodLogs(ctx context.Context, namespace, podName string, opts kubemvk8s.LogOptions) (string, error) {
	tail := opts.TailLines
	if tail == 0 {
		tail = 500
	}
	req := c.clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Container: opts.Container,
		Previous:  opts.Previous,
		TailLines: &tail,
	})
	stream, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer stream.Close()
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, stream); err != nil {
		return "", err
	}
	text := buf.String()
	if opts.Search == "" {
		return text, nil
	}
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(strings.ToLower(line), strings.ToLower(opts.Search)) {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n"), nil
}

// RESTConfig exposes the underlying config for exec/port-forward handlers.
func (c *liveClient) RESTConfig() *rest.Config {
	return c.cfg
}
