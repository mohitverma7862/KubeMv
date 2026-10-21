package kube

import (
	"context"
	"strings"
	"testing"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSecretYAMLOmitsValues(t *testing.T) {
	client := fake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "db", Namespace: "production"},
		Type:       corev1.SecretTypeOpaque,
		Data:       map[string][]byte{"password": []byte("BEGIN PRIVATE KEY supersecret")},
	})
	live := &liveClient{clientset: client}
	text, err := live.fetchYAML(context.Background(), kubemvk8s.ResourceSecrets, "production", "db")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(text, "supersecret") || strings.Contains(text, "BEGIN PRIVATE KEY") {
		t.Fatalf("secret yaml leaked:\n%s", text)
	}
	if !strings.Contains(text, "secret data omitted") || !strings.Contains(text, "name: db") {
		t.Fatalf("yaml:\n%s", text)
	}
}

func TestPodLogOptionsCarryFollowAndContainer(t *testing.T) {
	opts := podLogOptions(kubemvk8s.LogOptions{Container: "api", Previous: true, Timestamps: true, TailLines: 20}, true)
	if !opts.Follow || !opts.Previous || !opts.Timestamps || opts.Container != "api" || opts.TailLines == nil || *opts.TailLines != 20 {
		t.Fatalf("%#v", opts)
	}
}

func TestSizeQueueResize(t *testing.T) {
	queue := NewSizeQueue()
	queue.Push(80, 24)
	queue.Push(100, 40)
	first := queue.Next()
	if first == nil || first.Width != 80 || first.Height != 24 {
		t.Fatalf("first %#v", first)
	}
	queue.Close()
}
