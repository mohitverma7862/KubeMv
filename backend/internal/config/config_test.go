package config

import (
	"testing"
	"time"
)

func TestValidateBootstrapPassword(t *testing.T) {
	cfg := Config{
		Addr:              "127.0.0.1:8787",
		AllowedOrigins:    []string{"http://127.0.0.1:1420"},
		BootstrapUsername: "admin",
		BootstrapPassword: "short",
		SessionTTL:        time.Hour,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected short password to fail")
	}
	cfg.BootstrapPassword = "foundation-test-password"
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	cfg.BootstrapUsername = "1admin"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected username to fail")
	}
}
