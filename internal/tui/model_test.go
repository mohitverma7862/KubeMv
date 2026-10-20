package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestIsProduction(t *testing.T) {
	if !IsProduction("prod-cluster", "prod") {
		t.Fatal("prod context should be production")
	}
	if IsProduction("staging", "dev") {
		t.Fatal("staging context should not be production")
	}
	if IsProduction("nonprod", "nonprod-east") {
		t.Fatal("nonprod context should not be production")
	}
}

func TestMovementAndFilter(t *testing.T) {
	model := Model{
		Kinds: []Kind{{Label: "Pods", Path: kubernetes.ResourcePods}, {Label: "Deployments", Path: kubernetes.ResourceDeployments}},
		Rows: map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow{
			kubernetes.ResourcePods: {
				{Name: "api", Status: "Running"},
				{Name: "payments-worker", Status: "CrashLoopBackOff"},
				{Name: "web", Status: "Running"},
			},
		},
	}
	model.OnRune('j')
	selected, ok := model.Selected()
	if !ok || selected.Name != "payments-worker" {
		t.Fatalf("cursor = %s", selected.Name)
	}
	model.OnRune('G')
	selected, _ = model.Selected()
	if selected.Name != "web" {
		t.Fatalf("last = %s", selected.Name)
	}
	model.OnRune('g')
	model.OnRune('/')
	model.OnRune('w')
	model.OnRune('o')
	if got := model.Visible(); len(got) != 1 || got[0].Name != "payments-worker" {
		t.Fatalf("filter = %#v", got)
	}
	model.OnKey("esc")
	if model.Filtering || model.Filter != "" || len(model.Visible()) != 3 {
		t.Fatal("esc should clear the filter")
	}
	model.OnRune(':')
	model.OnRune('d')
	model.OnRune('p')
	if !model.OnKey("enter") && model.kind().Path != kubernetes.ResourceDeployments {
		t.Fatal("expected deployments")
	}
	model.OnRune(':')
	model.OnRune('q')
	if !model.OnKey("enter") {
		t.Fatal(":q should quit")
	}
}

func TestReferenceLayout(t *testing.T) {
	model := sampleModel("sample")
	if !strings.Contains(model.HeaderLine(time.Now()), "Fast Mode (K9s-Parity)") || !strings.Contains(model.HeaderLine(time.Now()), "PRODUCTION") {
		t.Fatal(model.HeaderLine(time.Now()))
	}
	if !strings.Contains(model.NavLine(), "[I] Ingresses") || !strings.Contains(model.NavLine(), "[>] More") {
		t.Fatal(model.NavLine())
	}
	if model.TabLine() != "[Logs] | Events | YAML | Describe | Terminal | Rollout | History" {
		t.Fatal(model.TabLine())
	}
	if !strings.Contains(model.StatusLine(), "R readonly") || !strings.Contains(model.StatusLine(), "0 forwards") {
		t.Fatal(model.StatusLine())
	}
	if !strings.Contains(QuickHelp(), "Ctrl-P    Palette") {
		t.Fatal(QuickHelp())
	}
	if !strings.Contains(strings.Join(columnsFor(kubernetes.ResourcePods), "|"), "CPU|MEM") {
		t.Fatal(columnsFor(kubernetes.ResourcePods))
	}
	selected, _ := model.Selected()
	if strings.Contains(model.WorkspaceText(), "BEGIN") || strings.Contains(model.InspectorText(), "BEGIN") {
		t.Fatal("secret material rendered")
	}
	_ = selected
}
