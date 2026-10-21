package kube

import (
	"context"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// ExecOptions is the shared SPDY exec request used by the terminal and the websocket handler.
type ExecOptions struct {
	Namespace string
	Pod       string
	Container string
	Command   []string
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	SizeQueue remotecommand.TerminalSizeQueue
}

// StreamExec runs a TTY shell through the API server. Credentials stay on the kube client.
func StreamExec(ctx context.Context, cfg *rest.Config, opts ExecOptions) error {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}
	command := opts.Command
	if len(command) == 0 {
		command = []string{"/bin/sh"}
	}
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(opts.Pod).
		Namespace(opts.Namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: opts.Container,
			Command:   command,
			Stdin:     opts.Stdin != nil,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	executor, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return err
	}
	return executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             opts.Stdin,
		Stdout:            opts.Stdout,
		Stderr:            opts.Stderr,
		Tty:               true,
		TerminalSizeQueue: opts.SizeQueue,
	})
}

// SizeQueue feeds terminal resize events to an exec stream.
type SizeQueue struct {
	ch chan remotecommand.TerminalSize
}

func NewSizeQueue() *SizeQueue {
	return &SizeQueue{ch: make(chan remotecommand.TerminalSize, 4)}
}

func (q *SizeQueue) Next() *remotecommand.TerminalSize {
	size, ok := <-q.ch
	if !ok {
		return nil
	}
	return &size
}

func (q *SizeQueue) Push(width, height int) {
	if width <= 0 || height <= 0 {
		return
	}
	select {
	case q.ch <- remotecommand.TerminalSize{Width: uint16(width), Height: uint16(height)}:
	default:
		select {
		case <-q.ch:
		default:
		}
		q.ch <- remotecommand.TerminalSize{Width: uint16(width), Height: uint16(height)}
	}
}

func (q *SizeQueue) Close() { close(q.ch) }
