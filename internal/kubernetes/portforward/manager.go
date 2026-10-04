package portforward

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"sync"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"k8s.io/client-go/rest"
)

// Manager tracks server-side port-forward sessions (credentials stay on server).
type Manager struct {
	mu        sync.RWMutex
	sessions  map[string]*session
	connector kubemvk8s.Connector
}

type session struct {
	clusterID string
	info      kubemvk8s.PortForwardSession
	cancel    context.CancelFunc
}

func NewManager(connector kubemvk8s.Connector) *Manager {
	return &Manager{
		sessions:  map[string]*session{},
		connector: connector,
	}
}

func (m *Manager) List(_ context.Context, clusterID string) []kubemvk8s.PortForwardSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]kubemvk8s.PortForwardSession, 0)
	for _, s := range m.sessions {
		if s.clusterID == clusterID {
			out = append(out, s.info)
		}
	}
	return out
}

func (m *Manager) Create(ctx context.Context, clusterID string, req kubemvk8s.PortForwardRequest) (kubemvk8s.PortForwardSession, error) {
	if req.LocalPort == 0 {
		req.LocalPort = req.RemotePort
	}
	client, err := m.connector.Connect(ctx, clusterID)
	if err != nil {
		return kubemvk8s.PortForwardSession{}, err
	}
	id, err := randomID()
	if err != nil {
		return kubemvk8s.PortForwardSession{}, err
	}
	info := kubemvk8s.PortForwardSession{
		ID: id, Namespace: req.Namespace, Pod: req.Pod,
		LocalPort: req.LocalPort, RemotePort: req.RemotePort,
		Status: "active", URL: fmt.Sprintf("http://127.0.0.1:%d", req.LocalPort),
	}
	var cancel context.CancelFunc
	if provider, ok := client.(interface{ RESTConfig() *rest.Config }); ok && provider.RESTConfig() != nil {
		_, cancel, err = startListener(ctx, req)
		if err != nil {
			return kubemvk8s.PortForwardSession{}, err
		}
	}
	m.mu.Lock()
	m.sessions[id] = &session{clusterID: clusterID, info: info, cancel: cancel}
	m.mu.Unlock()
	return info, nil
}

func (m *Manager) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	s, ok := m.sessions[id]
	if ok {
		if s.cancel != nil {
			s.cancel()
		}
		delete(m.sessions, id)
	}
	m.mu.Unlock()
	if !ok {
		return fmt.Errorf("session not found")
	}
	return nil
}

func randomID() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func startListener(ctx context.Context, req kubemvk8s.PortForwardRequest) (int, context.CancelFunc, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", req.LocalPort))
	if err != nil {
		return 0, nil, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer listener.Close()
		for {
			if ctx.Err() != nil {
				return
			}
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_, _ = c.Write([]byte(fmt.Sprintf("KubeMv forward %s/%s:%d (Phase 2 tunnel placeholder)\n", req.Namespace, req.Pod, req.RemotePort)))
			}(conn)
		}
	}()
	return port, cancel, nil
}
