package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	minPasswordLength = 12
	maxPasswordLength = 128
	defaultAddr       = "127.0.0.1:8787"
	defaultTTL        = 8 * time.Hour
)

// Config is process configuration. Secrets come from the environment and are
// not written back to disk.
type Config struct {
	Addr              string
	AllowedOrigins    []string
	BootstrapUsername string
	BootstrapPassword string
	SessionTTL        time.Duration
	CookieSecure      bool
}

func Load() (Config, error) {
	cfg := Config{
		Addr:              getenv("KUBEMV_ADDR", defaultAddr),
		BootstrapUsername: strings.TrimSpace(os.Getenv("KUBEMV_BOOTSTRAP_USERNAME")),
		BootstrapPassword: os.Getenv("KUBEMV_BOOTSTRAP_PASSWORD"),
		SessionTTL:        defaultTTL,
		CookieSecure:      strings.EqualFold(os.Getenv("KUBEMV_COOKIE_SECURE"), "true"),
	}

	origins := getenv("KUBEMV_ALLOWED_ORIGINS", "http://127.0.0.1:1420,http://localhost:1420")
	for _, origin := range strings.Split(origins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, origin)
		}
	}

	if raw := strings.TrimSpace(os.Getenv("KUBEMV_SESSION_TTL")); raw != "" {
		ttl, err := time.ParseDuration(raw)
		if err != nil {
			return Config{}, fmt.Errorf("KUBEMV_SESSION_TTL: %w", err)
		}
		if ttl < time.Minute || ttl > 24*time.Hour {
			return Config{}, errors.New("KUBEMV_SESSION_TTL must be between 1m and 24h")
		}
		cfg.SessionTTL = ttl
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.Addr == "" {
		return errors.New("address is required")
	}
	if len(c.AllowedOrigins) == 0 {
		return errors.New("at least one allowed origin is required")
	}
	if err := validateUsername(c.BootstrapUsername); err != nil {
		return fmt.Errorf("KUBEMV_BOOTSTRAP_USERNAME: %w", err)
	}
	if len(c.BootstrapPassword) < minPasswordLength || len(c.BootstrapPassword) > maxPasswordLength {
		return fmt.Errorf("KUBEMV_BOOTSTRAP_PASSWORD must be %d to %d characters", minPasswordLength, maxPasswordLength)
	}
	return nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return errors.New("must be 3 to 32 characters")
	}
	for i, r := range username {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
		case i > 0 && (r >= '0' && r <= '9' || r == '.' || r == '_' || r == '-'):
		default:
			return errors.New("must start with a letter and use only letters, digits, '.', '_' or '-'")
		}
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

// ParseBool is used by tests and optional flags.
func ParseBool(raw string) bool {
	v, err := strconv.ParseBool(strings.TrimSpace(raw))
	return err == nil && v
}
