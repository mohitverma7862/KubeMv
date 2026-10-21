package kube

import (
	"context"
	"io"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

// ResourceWatch is the live-client capability used by the terminal console.
type ResourceWatch interface {
	WatchSource(kind kubemvk8s.ResourceKindPath, namespace string) kubemvk8s.WatchSource
}

// WatchSource lists and watches one kind from the current kubeconfig context.
func (c *liveClient) WatchSource(kind kubemvk8s.ResourceKindPath, namespace string) kubemvk8s.WatchSource {
	return kindSource{client: c, kind: kind, namespace: namespace}
}

type kindSource struct {
	client    *liveClient
	kind      kubemvk8s.ResourceKindPath
	namespace string
}

func (s kindSource) List(ctx context.Context) (kubemvk8s.Snapshot, error) {
	rows, rv, err := s.client.listVersion(ctx, s.kind, kubemvk8s.ListOptions{Namespace: s.namespace})
	if err != nil {
		return kubemvk8s.Snapshot{}, err
	}
	return kubemvk8s.Snapshot{Rows: rows, ResourceVersion: rv}, nil
}

func (s kindSource) Watch(ctx context.Context, resourceVersion string) (kubemvk8s.WatchStream, error) {
	watcher, err := s.client.openWatch(ctx, s.kind, s.namespace, resourceVersion)
	if err != nil {
		return nil, err
	}
	return &watchStream{watcher: watcher}, nil
}

func (c *liveClient) openWatch(ctx context.Context, kind kubemvk8s.ResourceKindPath, namespace, resourceVersion string) (watch.Interface, error) {
	opts := metav1.ListOptions{ResourceVersion: resourceVersion, AllowWatchBookmarks: true}
	ns := namespaceOrAll(namespace)
	switch kind {
	case kubemvk8s.ResourcePods:
		return c.clientset.CoreV1().Pods(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceDeployments:
		return c.clientset.AppsV1().Deployments(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceServices:
		return c.clientset.CoreV1().Services(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceNodes:
		return c.clientset.CoreV1().Nodes().Watch(ctx, opts)
	case kubemvk8s.ResourceEvents:
		return c.clientset.CoreV1().Events(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceConfigMaps:
		return c.clientset.CoreV1().ConfigMaps(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceSecrets:
		return c.clientset.CoreV1().Secrets(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceNamespaces:
		return c.clientset.CoreV1().Namespaces().Watch(ctx, opts)
	case kubemvk8s.ResourceCRDs:
		return c.ext.ApiextensionsV1().CustomResourceDefinitions().Watch(ctx, opts)
	case kubemvk8s.ResourceStatefulSets:
		return c.clientset.AppsV1().StatefulSets(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceDaemonSets:
		return c.clientset.AppsV1().DaemonSets(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceJobs:
		return c.clientset.BatchV1().Jobs(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceCronJobs:
		return c.clientset.BatchV1().CronJobs(ns).Watch(ctx, opts)
	case kubemvk8s.ResourceHPA:
		return c.clientset.AutoscalingV2().HorizontalPodAutoscalers(ns).Watch(ctx, opts)
	case kubemvk8s.ResourcePDB:
		return c.clientset.PolicyV1().PodDisruptionBudgets(ns).Watch(ctx, opts)
	default:
		return nil, errUnsupported(kind)
	}
}

func errUnsupported(kind kubemvk8s.ResourceKindPath) error {
	return &unsupportedKind{kind: kind}
}

type unsupportedKind struct{ kind kubemvk8s.ResourceKindPath }

func (e *unsupportedKind) Error() string { return "unsupported kind " + string(e.kind) }

type watchStream struct {
	watcher watch.Interface
}

func (s *watchStream) Close() { s.watcher.Stop() }

func (s *watchStream) Next() (kubemvk8s.Delta, error) {
	event, ok := <-s.watcher.ResultChan()
	if !ok {
		return kubemvk8s.Delta{}, io.EOF
	}
	switch event.Type {
	case watch.Bookmark:
		rv := ""
		if event.Object != nil {
			rv = objectVersion(event.Object)
		}
		return kubemvk8s.Delta{Type: kubemvk8s.WatchBookmark, ResourceVersion: rv}, nil
	case watch.Error:
		message := "watch error"
		gone := false
		if status, ok := event.Object.(*metav1.Status); ok {
			message = status.Message
			gone = status.Code == 410
		}
		return kubemvk8s.Delta{Type: kubemvk8s.WatchError, Gone: gone, Message: message}, nil
	default:
		row, ok := rowFromObject(event.Object)
		if !ok {
			return kubemvk8s.Delta{Type: string(event.Type), ResourceVersion: objectVersion(event.Object)}, nil
		}
		return kubemvk8s.Delta{Type: string(event.Type), Row: row, ResourceVersion: objectVersion(event.Object)}, nil
	}
}

func objectVersion(obj runtime.Object) string {
	accessor, err := meta.Accessor(obj)
	if err != nil || accessor == nil {
		return ""
	}
	return accessor.GetResourceVersion()
}
