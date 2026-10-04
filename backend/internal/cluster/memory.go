package cluster

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sort"
	"strings"
	"sync"
	"time"
)

const maxClusters = 100

// MemoryRegistry is the Phase 0 cluster registry. Records live in the API
// process and disappear on restart.
type MemoryRegistry struct {
	now func() time.Time

	mu       sync.Mutex
	clusters map[string]Cluster
}

func NewMemoryRegistry(now func() time.Time) *MemoryRegistry {
	if now == nil {
		now = time.Now
	}
	return &MemoryRegistry{
		now:      now,
		clusters: make(map[string]Cluster),
	}
}

func (r *MemoryRegistry) List(_ context.Context) ([]Cluster, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Cluster, 0, len(r.clusters))
	for _, cluster := range r.clusters {
		out = append(out, cluster)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out, nil
}

func (r *MemoryRegistry) Get(_ context.Context, id string) (Cluster, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cluster, ok := r.clusters[id]
	if !ok {
		return Cluster{}, ErrNotFound
	}
	return cluster, nil
}

func (r *MemoryRegistry) Register(_ context.Context, in Registration) (Cluster, error) {
	if err := ValidateRegistration(in); err != nil {
		return Cluster{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.clusters) >= maxClusters {
		return Cluster{}, ErrLimit
	}
	for _, existing := range r.clusters {
		if strings.EqualFold(existing.Name, in.Name) {
			return Cluster{}, ErrConflict
		}
	}
	id, err := newID()
	if err != nil {
		return Cluster{}, err
	}
	cluster := Cluster{
		ID:               id,
		Name:             in.Name,
		Provider:         in.Provider,
		Context:          in.Context,
		ConnectionState:  ConnectionNotConnected,
		ConnectionDetail: ConnectionDetailPhase0,
		CreatedAt:        r.now().UTC(),
		KubeconfigRef:    in.KubeconfigRef,
	}
	r.clusters[id] = cluster
	return cluster, nil
}

func (r *MemoryRegistry) Remove(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.clusters[id]; !ok {
		return ErrNotFound
	}
	delete(r.clusters, id)
	return nil
}

func newID() (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return "cls_" + hex.EncodeToString(buf), nil
}
