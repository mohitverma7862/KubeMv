package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes/kube"
	"k8s.io/client-go/tools/clientcmd"
)

// Load reads the current kubeconfig context. Without a usable context it returns the sample console.
func Load(ctx context.Context) Model {
	model, err := loadKube(ctx)
	if err != nil {
		return sampleModel("sample — " + err.Error())
	}
	return model
}

func loadKube(ctx context.Context) (Model, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	raw, err := loader.RawConfig()
	if err != nil {
		return Model{}, err
	}
	if raw.CurrentContext == "" {
		return Model{}, fmt.Errorf("kubeconfig has no current context")
	}
	current, ok := raw.Contexts[raw.CurrentContext]
	if !ok || current == nil {
		return Model{}, fmt.Errorf("unknown current context %q", raw.CurrentContext)
	}
	connector, err := kube.NewConnectorFromDefaultLoadingRules()
	if err != nil {
		return Model{}, err
	}
	call, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	client, err := connector.Connect(call, raw.CurrentContext)
	if err != nil {
		return Model{}, err
	}
	model := Model{
		Cluster:    current.Cluster,
		Context:    raw.CurrentContext,
		Namespace:  current.Namespace,
		User:       current.AuthInfo,
		Production: IsProduction(current.Cluster, raw.CurrentContext),
		ReadOnly:   true,
		Live:       true,
		Note:       "context " + raw.CurrentContext,
		Kinds:      defaultKinds(),
		Rows:       map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow{},
	}
	if model.Namespace == "" {
		model.Namespace = "all"
	}
	var failed []string
	for _, kind := range model.Kinds {
		rows, err := client.ListResources(call, kind.Path, kubernetes.ListOptions{Namespace: current.Namespace})
		if err != nil {
			failed = append(failed, string(kind.Path))
			continue
		}
		model.Rows[kind.Path] = rows
	}
	if len(failed) > 0 {
		model.Note = "partial list: " + fmt.Sprint(failed)
	}
	return model, nil
}
