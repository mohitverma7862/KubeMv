package connector

import (
	"context"
	"log"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes/kube"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes/stub"
)

// NewDefault prefers live kubeconfig contexts and falls back to stub data.
func NewDefault() kubernetes.Connector {
	kubeConnector, err := kube.NewConnectorFromDefaultLoadingRules()
	if err == nil {
		clusters, listErr := kubeConnector.ListClusters(context.Background())
		if listErr == nil && len(clusters) > 0 {
			log.Printf("kubemv: using live Kubernetes connector (%d contexts)", len(clusters))
			return kubeConnector
		}
	}
	log.Printf("kubemv: using stub Kubernetes connector (%v)", err)
	return stub.NewConnector()
}
