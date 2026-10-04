package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("unexpected addr: %s", cfg.HTTPAddr)
	}
	if cfg.Environment != "development" {
		t.Fatalf("unexpected env: %s", cfg.Environment)
	}
}
