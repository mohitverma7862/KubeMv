package cluster

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateRegistrationRejectsSecrets(t *testing.T) {
	valid := Registration{Name: "prod-eks", Provider: ProviderEKS, Context: "prod", KubeconfigRef: "secret://clusters/prod"}
	if err := ValidateRegistration(valid); err != nil {
		t.Fatal(err)
	}
	cases := []Registration{
		{Name: "", Provider: ProviderGeneric, Context: "dev"},
		{Name: "bad name", Provider: ProviderGeneric, Context: "dev"},
		{Name: "dev", Provider: "openshift", Context: "dev"},
		{Name: "dev", Provider: ProviderGeneric, Context: "dev", KubeconfigRef: "apiVersion: v1\nkind: Config"},
		{Name: "dev", Provider: ProviderGeneric, Context: "dev", KubeconfigRef: "token: super-secret"},
		{Name: "dev", Provider: ProviderGeneric, Context: "dev", KubeconfigRef: "/tmp/../kubeconfig"},
		{Name: "dev", Provider: ProviderGeneric, Context: "dev", KubeconfigRef: "-----BEGIN PRIVATE KEY-----"},
	}
	for _, tc := range cases {
		if err := ValidateRegistration(tc); err != ErrInvalidRegistration {
			t.Fatalf("expected invalid for %+v, got %v", tc, err)
		}
	}
}

func TestMemoryRegistryHidesKubeconfigRef(t *testing.T) {
	registry := NewMemoryRegistry(nil)
	created, err := registry.Register(context.Background(), Registration{
		Name:          "staging",
		Provider:      ProviderGKE,
		Context:       "staging",
		KubeconfigRef: "secret://clusters/staging",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ConnectionState != ConnectionNotConnected {
		t.Fatalf("state %s", created.ConnectionState)
	}
	if created.KubeconfigRef != "secret://clusters/staging" {
		t.Fatal("server-side ref was not retained")
	}
	encoded, err := json.Marshal(created)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "secret://") || strings.Contains(string(encoded), "kubeconfig") {
		t.Fatalf("encoded cluster leaked reference: %s", encoded)
	}
	if _, err := registry.Register(context.Background(), Registration{Name: "Staging", Provider: ProviderGeneric, Context: "other"}); err != ErrConflict {
		t.Fatal("case-insensitive duplicate should conflict")
	}
	if err := registry.Remove(context.Background(), created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Get(context.Background(), created.ID); err != ErrNotFound {
		t.Fatal("expected not found")
	}
}
