package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mohitverma7862/KubeMv/internal/api"
	"github.com/mohitverma7862/KubeMv/internal/auth"
	"github.com/mohitverma7862/KubeMv/internal/config"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// Server is the KubeMv HTTP API process.
type Server struct {
	cfg    config.Config
	http   *http.Server
	router http.Handler
}

type Options struct {
	Config        config.Config
	Authenticator auth.Authenticator
	Connector     kubernetes.Connector
}

func New(opts Options) *Server {
	router := api.NewRouter(api.Dependencies{
		Authenticator: opts.Authenticator,
		Connector:     opts.Connector,
		PrometheusURL: opts.Config.PrometheusURL,
		AIEnabled:     opts.Config.AIEnabled,
	})
	srv := &http.Server{
		Addr:         opts.Config.HTTPAddr,
		Handler:      router,
		ReadTimeout:  opts.Config.ReadTimeout,
		WriteTimeout: opts.Config.WriteTimeout,
	}
	return &Server{cfg: opts.Config, http: srv, router: router}
}

func (s *Server) ListenAndServe() error {
	fmt.Printf("kubemv-api listening on %s (env=%s)\n", s.cfg.HTTPAddr, s.cfg.Environment)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
