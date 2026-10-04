package cluster

import (
	"errors"
	"strings"
	"time"
	"unicode"
)

// Provider identifies where a cluster runs. Phase 0 stores the provider as
// metadata and does not call a cloud or Kubernetes API.
type Provider string

const (
	ProviderGeneric Provider = "generic"
	ProviderEKS     Provider = "eks"
	ProviderGKE     Provider = "gke"
	ProviderAKS     Provider = "aks"
)

// ConnectionState is explicit so the UI cannot present a registry record as a
// live Kubernetes connection.
type ConnectionState string

const (
	ConnectionNotConnected ConnectionState = "not_connected"
)

// ConnectionDetailPhase0 is the only connection explanation Phase 0 produces.
const ConnectionDetailPhase0 = "The Kubernetes API is not contacted in Phase 0. This record is registry metadata only."

// Cluster is registry metadata. KubeconfigRef is server-side only and is
// omitted from every JSON encoding.
type Cluster struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Provider         Provider        `json:"provider"`
	Context          string          `json:"context"`
	ConnectionState  ConnectionState `json:"connectionState"`
	ConnectionDetail string          `json:"connectionDetail"`
	CreatedAt        time.Time       `json:"createdAt"`
	KubeconfigRef    string          `json:"-"`
}

// Registration is the operator input for a cluster record.
type Registration struct {
	Name          string
	Provider      Provider
	Context       string
	KubeconfigRef string
}

var (
	ErrInvalidRegistration = errors.New("invalid cluster registration")
	ErrNotFound            = errors.New("cluster not found")
	ErrConflict            = errors.New("cluster already registered")
	ErrLimit               = errors.New("cluster registry limit reached")
)

func KnownProvider(provider Provider) bool {
	switch provider {
	case ProviderGeneric, ProviderEKS, ProviderGKE, ProviderAKS:
		return true
	default:
		return false
	}
}

// ValidateRegistration rejects inline credentials and kubeconfig documents.
// A reference, when present, is an opaque server-side locator.
func ValidateRegistration(in Registration) error {
	if !validName(in.Name) || !validName(in.Context) || !KnownProvider(in.Provider) {
		return ErrInvalidRegistration
	}
	if err := validateKubeconfigRef(in.KubeconfigRef); err != nil {
		return err
	}
	return nil
}

func validName(value string) bool {
	if len(value) < 1 || len(value) > 63 {
		return false
	}
	for i, r := range value {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
		case i > 0 && (r == '.' || r == '_' || r == '-'):
		default:
			return false
		}
	}
	return true
}

func validateKubeconfigRef(ref string) error {
	if ref == "" {
		return nil
	}
	if len(ref) > 512 || strings.Contains(ref, "..") || strings.ContainsAny(ref, " \r\n\t") {
		return ErrInvalidRegistration
	}
	lowered := strings.ToLower(ref)
	banned := []string{
		"apiversion",
		"client-certificate-data",
		"client-key-data",
		"certificate-authority-data",
		"-----begin",
		"password=",
		"token:",
		"bearer ",
	}
	for _, marker := range banned {
		if strings.Contains(lowered, marker) {
			return ErrInvalidRegistration
		}
	}
	return nil
}
