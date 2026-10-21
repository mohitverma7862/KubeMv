package kubernetes

import "testing"

func TestOwnerTreeAndSecretRedaction(t *testing.T) {
	pod := ResourceRow{
		Kind: "Pod", Namespace: "production", Name: "api-7c9f8b6d4-k2m9x", Status: "Running",
		Owners:     []OwnerRef{{Kind: "ReplicaSet", Namespace: "production", Name: "api-7c9f8b6d4"}},
		Containers: []ContainerRef{{Name: "api", Image: "api:1.2.3"}},
	}
	rs := ResourceRow{
		Kind: "ReplicaSet", Namespace: "production", Name: "api-7c9f8b6d4",
		Owners: []OwnerRef{{Kind: "Deployment", Namespace: "production", Name: "api"}},
	}
	lookup := func(ref OwnerRef) (ResourceRow, bool) {
		if ref.Kind == "ReplicaSet" {
			return rs, true
		}
		return ResourceRow{}, false
	}
	nodes := BuildTree(pod, lookup)
	text := FormatTree(nodes, -1)
	want := "Deployment (api)\n└─ ReplicaSet (api-7c9f8b6d4)\n   └─ Pod (api-7c9f8b6d4-k2m9x)\n      └─ Containers\n         └─ api"
	if text != want {
		t.Fatalf("tree:\n%s", text)
	}
	if len(nodes) != 5 || nodes[0].Kind != "Deployment" || nodes[4].Name != "api" {
		t.Fatalf("nodes %#v", nodes)
	}

	secret := PublicRow(ResourceRow{
		Kind: "Secret", Name: "db", Namespace: "production",
		Labels: map[string]string{"app": "db", "token": "Bearer secret"},
		Extra:  map[string]string{"type": "Opaque", "keys": "2", "data": "BEGIN PRIVATE KEY", "token": "abc"},
	})
	if secret.Extra["data"] != "" || secret.Extra["token"] != "" || secret.Extra["type"] != "Opaque" || secret.Extra["keys"] != "2" {
		t.Fatalf("secret extra %#v", secret.Extra)
	}
	if _, ok := secret.Labels["token"]; ok {
		t.Fatal("secret label leaked")
	}
}
