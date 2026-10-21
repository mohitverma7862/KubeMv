package kube

import (
	"context"
	"testing"
	"time"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestPodWatchDeliversKeyedAdd(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "api", Namespace: "default", ResourceVersion: "1", CreationTimestamp: metav1.NewTime(time.Now().Add(-time.Hour))},
		Status:     corev1.PodStatus{Phase: corev1.PodRunning},
	})
	client.PrependReactor("list", "pods", func(action k8stesting.Action) (bool, runtime.Object, error) {
		listAction := action.(k8stesting.ListAction)
		obj, err := client.Tracker().List(listAction.GetResource(), schema.GroupVersionKind{Version: "v1", Kind: "Pod"}, listAction.GetNamespace())
		if err != nil {
			return true, nil, err
		}
		list := obj.(*corev1.PodList)
		list.ResourceVersion = "9"
		return true, list, nil
	})
	live := &liveClient{clientset: client}
	source := live.WatchSource(kubemvk8s.ResourcePods, "default")
	snap, err := source.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(snap.Rows) != 1 || snap.Rows[0].Name != "api" || snap.Rows[0].Status != string(corev1.PodRunning) || snap.ResourceVersion != "9" {
		t.Fatalf("snapshot %#v", snap)
	}
	stream, err := source.Watch(context.Background(), snap.ResourceVersion)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	_, err = client.CoreV1().Pods("default").Create(context.Background(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "web", Namespace: "default", ResourceVersion: "10"},
		Status:     corev1.PodStatus{Phase: corev1.PodPending},
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.After(2 * time.Second)
	for {
		got := make(chan kubemvk8s.Delta, 1)
		failed := make(chan error, 1)
		go func() {
			delta, err := stream.Next()
			if err != nil {
				failed <- err
				return
			}
			got <- delta
		}()
		select {
		case delta := <-got:
			if delta.Type == kubemvk8s.WatchBookmark {
				continue
			}
			if delta.Type == kubemvk8s.WatchAdded && delta.Row.Name == "web" && delta.Row.Status == string(corev1.PodPending) && delta.ResourceVersion != "" {
				return
			}
			t.Fatalf("delta %#v", delta)
		case err := <-failed:
			t.Fatal(err)
		case <-deadline:
			t.Fatal("timed out waiting for the added pod")
		}
	}
}
