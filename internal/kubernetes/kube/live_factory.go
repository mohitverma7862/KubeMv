package kube

import (
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func newLiveClient(id, clusterName string, cfg *rest.Config) (kubemvk8s.ClusterClient, error) {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	ext, err := apiextensionsclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &liveClient{
		ref:       kubemvk8s.ClusterRef{ID: id, Name: clusterName},
		clientset: clientset,
		ext:       ext,
		cfg:       cfg,
	}, nil
}
