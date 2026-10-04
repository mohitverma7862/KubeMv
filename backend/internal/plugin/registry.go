package plugin

import (
	"sort"
	"sync"
)

// Registry is the in-process plugin host. Registration is a code path, not an
// HTTP upload, so the UI cannot install executable plugins.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Manifest
}

func NewRegistry() *Registry {
	return &Registry{plugins: map[string]Manifest{}}
}

func (r *Registry) Register(p Plugin) error {
	manifest := p.Manifest()
	if err := ValidateManifest(manifest); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.plugins[manifest.ID]; exists {
		return ErrDuplicate
	}
	manifest.Permissions = append([]Permission(nil), manifest.Permissions...)
	r.plugins[manifest.ID] = manifest
	return nil
}

func (r *Registry) List() []Manifest {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Manifest, 0, len(r.plugins))
	for _, manifest := range r.plugins {
		copied := manifest
		copied.Permissions = append([]Permission(nil), manifest.Permissions...)
		out = append(out, copied)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

func (r *Registry) Get(id string) (Manifest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	manifest, ok := r.plugins[id]
	if !ok {
		return Manifest{}, ErrUnknown
	}
	manifest.Permissions = append([]Permission(nil), manifest.Permissions...)
	return manifest, nil
}
