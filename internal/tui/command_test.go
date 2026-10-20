package tui

import (
	"strings"
	"testing"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestCommandRegistryAliasesAndTab(t *testing.T) {
	model := Model{Kinds: []Kind{{Label: "Pods", Path: kubernetes.ResourcePods}, {Label: "Deployments", Path: kubernetes.ResourceDeployments}}}
	model.OnRune(':')
	model.OnRune('p')
	if !strings.Contains(model.CommandError, "pods") {
		t.Fatalf("suggestions = %q", model.CommandError)
	}
	model.OnRune('o')
	model.OnKey("tab")
	if model.Command != "pods" {
		t.Fatalf("completion = %q", model.Command)
	}
	if model.OnKey("enter") {
		t.Fatal("pods should not quit")
	}
	if model.kind().Path != kubernetes.ResourcePods {
		t.Fatal(model.kind().Path)
	}
	model.OnRune(':')
	for _, char := range "deployments" {
		model.OnRune(char)
	}
	model.OnKey("enter")
	if model.kind().Path != kubernetes.ResourceDeployments {
		t.Fatal("alias deployments")
	}
	if _, ok := Lookup("po"); !ok {
		t.Fatal("po alias")
	}
	if !strings.Contains(PaletteText(), ":pods") || !strings.Contains(QuickHelp(), "Ctrl-P") {
		t.Fatal("palette and help must come from the registries")
	}
	if !strings.Contains(ContextualHelp(true), "Tab") {
		t.Fatal(ContextualHelp(true))
	}
}

func TestFilterGrammar(t *testing.T) {
	model := Model{
		Kinds: []Kind{{Label: "Pods", Path: kubernetes.ResourcePods}},
		Rows: map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow{
			kubernetes.ResourcePods: {
				{Name: "api", Namespace: "production", Status: "Running", Labels: map[string]string{"app": "api"}},
				{Name: "payments-worker", Namespace: "production", Status: "CrashLoopBackOff", Labels: map[string]string{"app": "payments"}},
				{Name: "web", Namespace: "staging", Status: "Running", Labels: map[string]string{"app": "web"}},
			},
		},
	}
	model.Filter = "status=Running ns=production"
	if got := model.Visible(); len(got) != 1 || got[0].Name != "api" {
		t.Fatalf("field filter %#v", got)
	}
	model.Filter = "app=payments"
	if got := model.Visible(); len(got) != 1 || got[0].Name != "payments-worker" {
		t.Fatalf("label filter %#v", got)
	}
	model.Filter = "/worker/"
	if got := model.Visible(); len(got) != 1 || got[0].Name != "payments-worker" {
		t.Fatalf("regex %#v", got)
	}
	model.Filter = "pwk"
	if got := model.Visible(); len(got) != 1 || got[0].Name != "payments-worker" {
		t.Fatalf("fuzzy %#v", got)
	}
}
