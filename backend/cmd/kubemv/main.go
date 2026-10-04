package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mohitverma7862/KubeMv/backend/internal/api"
	"github.com/mohitverma7862/KubeMv/backend/internal/config"
	"github.com/mohitverma7862/KubeMv/backend/internal/version"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		log.Error("configuration", "err", err.Error())
		os.Exit(1)
	}
	app, err := api.New(cfg, log)
	if err != nil {
		log.Error("startup", "err", err.Error())
		os.Exit(1)
	}
	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           app.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 14,
	}
	go func() {
		log.Info("api listening", "addr", cfg.Addr, "version", version.Version, "phase", version.Phase)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("listen", "err", err.Error())
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown", "err", err.Error())
		os.Exit(1)
	}
}
