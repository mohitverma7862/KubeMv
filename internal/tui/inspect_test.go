package tui

import (
	"strings"
	"testing"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestInspectorAndOwnerTree(t *testing.T) {
	model := sampleModel("sample")
	body := model.InspectorText()
	for _, part := range []string{"Running", "namespace", "production", "ip-10-0-1-12", "10.8.1.14", "1/1", "4h", "Burstable", "42m", "128Mi", "api:1.2.3", "app=api", "ReplicaSet/api-7c9f8b6d4", "l Logs", "D Delete"} {
		if !strings.Contains(body, part) {
			t.Fatalf("inspector missing %q\n%s", part, body)
		}
	}
	tree := model.TreeText()
	want := "Deployment (api)\n└─ ReplicaSet (api-7c9f8b6d4)\n   └─ Pod (api-7c9f8b6d4-k2m9x)\n      └─ Containers\n         └─ api"
	if tree != want {
		t.Fatalf("tree:\n%s", tree)
	}
	model.OnKey("tab")
	if !model.TreeFocus || !strings.HasPrefix(model.TreeText(), "> Deployment") {
		t.Fatalf("focus %v\n%s", model.TreeFocus, model.TreeText())
	}
	model.OnKey("enter")
	if model.kind().Path != kubernetes.ResourceDeployments || model.Filter != "api" || model.TreeFocus {
		t.Fatalf("drill kind=%s filter=%s focus=%v", model.kind().Path, model.Filter, model.TreeFocus)
	}

	secret := kubernetes.ResourceRow{
		Kind: "Secret", Namespace: "production", Name: "db", Status: "Active",
		Extra: map[string]string{"type": "Opaque", "keys": "2", "data": "BEGIN PRIVATE KEY", "token": "supersecret"},
	}
	model.Rows[kubernetes.ResourceSecrets] = []kubernetes.ResourceRow{secret}
	model.selectPath(kubernetes.ResourceSecrets)
	model.Filter = ""
	inspected := model.InspectorText()
	if strings.Contains(inspected, "BEGIN") || strings.Contains(inspected, "supersecret") || strings.Contains(model.TreeText(), "supersecret") {
		t.Fatalf("secret material rendered\n%s\n%s", inspected, model.TreeText())
	}
	if !strings.Contains(inspected, "secret values hidden") || !strings.Contains(inspected, "Opaque") {
		t.Fatal(inspected)
	}
}
