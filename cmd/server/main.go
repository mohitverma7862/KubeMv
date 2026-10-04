package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mohitverma7862/KubeMv/internal/auth/dev"
	"github.com/mohitverma7862/KubeMv/internal/config"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes/stub"
	"github.com/mohitverma7862/KubeMv/internal/server"
)

func main() {
	cfg := config.Load()
	if cfg.Environment == "production" && cfg.EnableDevAuth {
		_, _ = os.Stderr.WriteString("warning: KUBEMV_DEV_AUTH is enabled in production\n")
	}

	srv := server.New(server.Options{
		Config:        cfg,
		Authenticator: dev.NewAuthenticator(),
		Connector:     stub.NewConnector(),
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			_, _ = os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
