package tui

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes/kube"
	"github.com/rivo/tview"
	"golang.org/x/term"
	"k8s.io/client-go/rest"
)

func runPodExec(ctx context.Context, cfg *rest.Config, namespace, pod, container string) error {
	queue := kube.NewSizeQueue()
	defer queue.Close()
	push := func() {
		width, height, err := term.GetSize(int(os.Stdout.Fd()))
		if err == nil {
			queue.Push(width, height)
		}
	}
	push()
	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)
	defer signal.Stop(winch)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-winch:
				push()
			}
		}
	}()
	return kube.StreamExec(ctx, cfg, kube.ExecOptions{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		SizeQueue: queue,
	})
}

func runSideEffects(app *tview.Application, model *Model) {
	if model.LogSave {
		model.LogSave = false
		name := "logs"
		if row, ok := model.Selected(); ok && row.Name != "" {
			name = row.Name
		}
		path, err := SaveLogFile(os.TempDir(), name, model.Logs.Download())
		if err != nil {
			model.CommandError = err.Error()
		} else {
			model.CommandError = "saved " + path
		}
	}
	if !model.ExecRequested {
		return
	}
	model.ExecRequested = false
	row, ok := model.Selected()
	exec := model.exec
	if !ok || exec == nil {
		return
	}
	container := model.Logs.ContainerName()
	app.Suspend(func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := exec(ctx, row.Namespace, row.Name, container); err != nil {
			fmt.Fprintf(os.Stderr, "exec: %v\n", err)
		}
	})
}
