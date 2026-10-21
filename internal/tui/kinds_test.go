package tui

import (
	"testing"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestNavMnemonicsAndDigits(t *testing.T) {
	model := sampleModel("sample")
	if model.kind().Path != kubernetes.ResourcePods {
		t.Fatalf("start = %s", model.kind().Path)
	}
	if !model.trySelectMnemonic('E') {
		t.Fatal("E should select events")
	}
	if model.kind().Path != kubernetes.ResourceEvents {
		t.Fatalf("events = %s", model.kind().Path)
	}
	if len(model.Visible()) == 0 {
		t.Fatal("sample events should be visible")
	}
	if !model.trySelectMnemonic('e') {
		t.Fatal("e should select events")
	}
	if !model.trySelectMnemonic('D') {
		t.Fatal("D should select deployments")
	}
	if model.kind().Path != kubernetes.ResourceDeployments {
		t.Fatalf("deploy = %s", model.kind().Path)
	}
	if !model.trySelectDigit('1') {
		t.Fatal("1 should select first resource")
	}
	if model.kind().Path != kubernetes.ResourcePods {
		t.Fatalf("digit 1 = %s", model.kind().Path)
	}
	if !model.trySelectMnemonic('D') {
		t.Fatal("D should re-select deployments")
	}
	if !model.trySelectMnemonic('I') {
		t.Fatal("I should be handled")
	}
	if model.kind().Path != kubernetes.ResourceDeployments {
		t.Fatal("unavailable ingress should not change the active view")
	}
	if model.CommandError == "" {
		t.Fatal("unavailable ingress should set command error")
	}
}
