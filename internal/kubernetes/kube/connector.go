package kube

import (
	"context"
	"fmt"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Connector loads kubeconfig contexts and returns live clients.
type Connector struct {
	contexts map[string]clientcmdapiContext
}

type clientcmdapiContext struct {
	name    string
	cluster string
}

// NewConnectorFromDefaultLoadingRules returns a connector or an error if kubeconfig is unavailable.
func NewConnectorFromDefaultLoadingRules() (*Connector, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	raw, err := loader.RawConfig()
	if err != nil {
		return nil, err
	}
	if len(raw.Contexts) == 0 {
		return nil, fmt.Errorf("kubeconfig has no contexts")
	}
	contexts := map[string]clientcmdapiContext{}
	for name, ctx := range raw.Contexts {
		contexts[name] = clientcmdapiContext{name: name, cluster: ctx.Cluster}
	}
	return &Connector{contexts: contexts}, nil
}

func (c *Connector) ListClusters(_ context.Context) ([]kubernetes.ClusterRef, error) {
	out := make([]kubernetes.ClusterRef, 0, len(c.contexts))
	for id, ctx := range c.contexts {
		out = append(out, kubernetes.ClusterRef{ID: id, Name: ctx.cluster})
	}
	return out, nil
}

func (c *Connector) Connect(ctx context.Context, clusterID string) (kubernetes.ClusterClient, error) {
	if _, ok := c.contexts[clusterID]; !ok {
		return nil, fmt.Errorf("unknown context %q", clusterID)
	}
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{CurrentContext: clusterID}
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
	cfg, err := loader.ClientConfig()
	if err != nil {
		return nil, err
	}
	return newClient(ctx, clusterID, c.contexts[clusterID].cluster, cfg)
}

func newClient(_ context.Context, id, clusterName string, cfg *rest.Config) (kubernetes.ClusterClient, error) {
	return newLiveClient(id, clusterName, cfg)
}
