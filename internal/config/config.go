package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds runtime configuration for the KubeMv API server.
type Config struct {
	HTTPAddr        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	Environment     string
	EnableDevAuth   bool
	PrometheusURL   string
	GrafanaURL      string
	LokiURL         string
	AIEnabled       bool
}

// Load reads configuration from environment variables with safe defaults.
func Load() Config {
	return Config{
		HTTPAddr:        envOr("KUBEMV_HTTP_ADDR", ":8080"),
		ReadTimeout:     durationEnv("KUBEMV_HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    durationEnv("KUBEMV_HTTP_WRITE_TIMEOUT", 15*time.Second),
		ShutdownTimeout: durationEnv("KUBEMV_SHUTDOWN_TIMEOUT", 10*time.Second),
		Environment:     envOr("KUBEMV_ENV", "development"),
		EnableDevAuth:   envBool("KUBEMV_DEV_AUTH", true),
		PrometheusURL:   envOr("KUBEMV_PROMETHEUS_URL", ""),
		GrafanaURL:      envOr("KUBEMV_GRAFANA_URL", ""),
		LokiURL:         envOr("KUBEMV_LOKI_URL", ""),
		AIEnabled:       envBool("KUBEMV_AI_ENABLED", true),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
