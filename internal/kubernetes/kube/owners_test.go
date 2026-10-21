package kube

import (
	"context"
	"testing"

	kubemvk8s "github.com/mohitverma7862/KubeMv/internal/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestLookupOwnerWalksToDeployment(t *testing.T) {
	controller := true
	client := fake.NewSimpleClientset(
		&appsv1.ReplicaSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: "api-7c9f8b6d4", Namespace: "production",
				OwnerReferences: []metav1.OwnerReference{{
					Kind: "Deployment", Name: "api", Controller: &controller,
				}},
			},
		},
		&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "api-7c9f8b6d4-k2m9x", Namespace: "production",
				OwnerReferences: []metav1.OwnerReference{{
					Kind: "ReplicaSet", Name: "api-7c9f8b6d4", Controller: &controller,
				}},
			},
			Spec:   corev1.PodSpec{Containers: []corev1.Container{{Name: "api", Image: "api:1.2.3"}}},
			Status: corev1.PodStatus{Phase: corev1.PodRunning, QOSClass: corev1.PodQOSBurstable},
		},
	)
	live := &liveClient{clientset: client}
	pod, err := live.LookupOwner(context.Background(), kubemvk8s.OwnerRef{Kind: "Pod", Namespace: "production", Name: "api-7c9f8b6d4-k2m9x"})
	if err != nil {
		t.Fatal(err)
	}
	if pod.Extra["qos"] != "Burstable" || len(pod.Containers) != 1 || pod.Containers[0].Image != "api:1.2.3" {
		t.Fatalf("pod %#v", pod)
	}
	rs, err := live.LookupOwner(context.Background(), pod.Owners[0])
	if err != nil || rs.Owners[0].Kind != "Deployment" || rs.Owners[0].Name != "api" {
		t.Fatalf("replicaset %#v %v", rs, err)
	}
	nodes := kubemvk8s.BuildTree(pod, func(ref kubemvk8s.OwnerRef) (kubemvk8s.ResourceRow, bool) {
		if ref.Kind == "ReplicaSet" {
			return rs, true
		}
		return kubemvk8s.ResourceRow{}, false
	})
	if nodes[0].Kind != "Deployment" || nodes[0].Name != "api" {
		t.Fatalf("nodes %#v", nodes)
	}
}
