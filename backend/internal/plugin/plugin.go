package plugin

import (
	"errors"
	"regexp"
)

// Category is a permission boundary for a plugin, not a guarantee that the
// category's product features are installed.
type Category string

const (
	CategoryObservability Category = "observability"
	CategorySecurity      Category = "security"
	CategoryCloud         Category = "cloud"
	CategoryGitOps        Category = "gitops"
	CategoryAI            Category = "ai"
	CategoryIncident      Category = "incident"
	CategoryCost          Category = "cost"
	CategoryBenchmark     Category = "benchmark"
	CategoryAutomation    Category = "automation"
)

// Permission is a host capability a plugin may request. Secret values and
// unrestricted command execution are not grantable.
type Permission string

const (
	PermClusterRead       Permission = "cluster.read"
	PermClusterRegister   Permission = "cluster.register"
	PermUIContribute      Permission = "ui.contribute"
	PermObservabilityRead Permission = "observability.read"
	PermGitRead           Permission = "git.read"
	PermAIAnalyze         Permission = "ai.analyze"
	PermAuditRead         Permission = "audit.read"
)

// Manifest describes a plugin. The host rejects manifests outside the catalog.
type Manifest struct {
	ID          string
	Name        string
	Version     string
	Category    Category
	Description string
	Permissions []Permission
}

// Plugin is the host contract. Phase 0 loads manifests only; plugins cannot
// execute code from the network.
type Plugin interface {
	Manifest() Manifest
}

var (
	ErrInvalidManifest = errors.New("invalid plugin manifest")
	ErrDuplicate       = errors.New("plugin already registered")
	ErrUnknown         = errors.New("plugin not found")

	idPattern      = regexp.MustCompile(`^[a-z][a-z0-9-]{1,63}$`)
	versionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
)

func KnownCategory(category Category) bool {
	switch category {
	case CategoryObservability, CategorySecurity, CategoryCloud, CategoryGitOps,
		CategoryAI, CategoryIncident, CategoryCost, CategoryBenchmark, CategoryAutomation:
		return true
	default:
		return false
	}
}

func KnownPermission(permission Permission) bool {
	switch permission {
	case PermClusterRead, PermClusterRegister, PermUIContribute, PermObservabilityRead,
		PermGitRead, PermAIAnalyze, PermAuditRead:
		return true
	default:
		return false
	}
}

func ValidateManifest(manifest Manifest) error {
	if !idPattern.MatchString(manifest.ID) || manifest.Name == "" || len(manifest.Name) > 80 {
		return ErrInvalidManifest
	}
	if !versionPattern.MatchString(manifest.Version) || !KnownCategory(manifest.Category) {
		return ErrInvalidManifest
	}
	if len(stringsTrim(manifest.Description)) == 0 || len(manifest.Description) > 400 {
		return ErrInvalidManifest
	}
	seen := map[Permission]struct{}{}
	if len(manifest.Permissions) == 0 {
		return ErrInvalidManifest
	}
	for _, permission := range manifest.Permissions {
		if !KnownPermission(permission) {
			return ErrInvalidManifest
		}
		if _, ok := seen[permission]; ok {
			return ErrInvalidManifest
		}
		seen[permission] = struct{}{}
	}
	return nil
}

func stringsTrim(value string) string {
	for len(value) > 0 && (value[0] == ' ' || value[0] == '\n' || value[0] == '\t') {
		value = value[1:]
	}
	for len(value) > 0 && (value[len(value)-1] == ' ' || value[len(value)-1] == '\n' || value[len(value)-1] == '\t') {
		value = value[:len(value)-1]
	}
	return value
}
