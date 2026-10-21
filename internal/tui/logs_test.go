package tui

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func TestLogWorkspaceControls(t *testing.T) {
	var logs LogBuffer
	logs.SetContainers([]string{"api", "proxy"})
	now := time.Date(2026, 10, 21, 15, 4, 5, 0, time.UTC)
	logs.SetText("level=info ready\nlevel=error boom\nlevel=info done", now)
	logs.ToggleFollow()
	logs.TogglePause()
	logs.Append("level=info later", now)
	if strings.Contains(logs.Download(), "later") {
		t.Fatal("paused logs accepted a line")
	}
	logs.TogglePause()
	logs.Append("level=info later", now)
	logs.ToggleTimestamps()
	if !strings.Contains(logs.Download(), "15:04:05 level=info later") {
		t.Fatal(logs.Download())
	}
	logs.Searching = true
	logs.Query = "error"
	if !strings.Contains(logs.Render(40), "> ") || !strings.Contains(logs.Status(), "1/1") {
		t.Fatal(logs.Render(40))
	}
	logs.Query = "info"
	logs.NextMatch()
	first := logs.match
	logs.PrevMatch()
	if logs.match == first {
		t.Fatal("match did not move")
	}
	if logs.NextContainer() != "proxy" {
		t.Fatal(logs.ContainerName())
	}
	logs.ToggleWrap()
	wrapped := logs.Render(10)
	if !strings.Contains(wrapped, "\n") {
		t.Fatal(wrapped)
	}
	logs.Clear()
	if !strings.Contains(logs.Render(80), "No logs loaded.") {
		t.Fatal(logs.Render(80))
	}
	path, err := SaveLogFile(t.TempDir(), "api pod", "line")
	if err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil || string(body) != "line" {
		t.Fatalf("saved %q %v", body, err)
	}
}

func TestWorkspaceTabsAndReadonlyExec(t *testing.T) {
	model := sampleModel("sample")
	model.clamp()
	text := model.WorkspaceText()
	if !strings.Contains(text, "level=info") || !strings.Contains(text, model.SelectedName()) {
		t.Fatal(text)
	}
	model.OnRune('\'')
	for _, char := range "ready" {
		model.OnRune(char)
	}
	if !strings.Contains(model.WorkspaceText(), "> ") {
		t.Fatal(model.WorkspaceText())
	}
	model.OnKey("esc")
	model.OnKey("ctrlc")
	model.Logs.ToggleFollow()
	if model.OnKey("ctrlc") {
		t.Fatal("ctrl-c should stop follow before quitting")
	}
	if model.Logs.Follow {
		t.Fatal("follow still on")
	}
	model.OnRune('d')
	if !strings.Contains(model.WorkspaceText(), "Name:") || !strings.Contains(model.WorkspaceText(), "api:1.2.3") {
		t.Fatal(model.WorkspaceText())
	}
	model.OnRune('y')
	if strings.Contains(model.WorkspaceText(), "BEGIN") {
		t.Fatal(model.WorkspaceText())
	}
	model.Rows[kubernetes.ResourceSecrets] = []kubernetes.ResourceRow{{
		Kind: "Secret", Namespace: "production", Name: "db", Status: "Active",
		Extra: map[string]string{"type": "Opaque", "keys": "1", "token": "supersecret"},
	}}
	model.selectPath(kubernetes.ResourceSecrets)
	model.Filter = ""
	model.YAML = ""
	model.Describe = ""
	model.OnRune('y')
	if strings.Contains(model.WorkspaceText(), "supersecret") || !strings.Contains(model.WorkspaceText(), "secret data omitted") {
		t.Fatal(model.WorkspaceText())
	}
	model.OnRune('d')
	if strings.Contains(model.WorkspaceText(), "supersecret") || !strings.Contains(model.WorkspaceText(), "secret values hidden") {
		t.Fatal(model.WorkspaceText())
	}
	model.OnKey("esc")
	model.selectPath(kubernetes.ResourcePods)
	model.Filter = ""
	model.openWorkspace("history")
	if !strings.Contains(model.WorkspaceText(), "api-7c9f8b6d4-k2m9x") {
		t.Fatal(model.WorkspaceText())
	}
	model.OnRune('s')
	if model.ExecRequested || !strings.Contains(model.WorkspaceText(), "readonly: exec is rejected") {
		t.Fatal(model.WorkspaceText())
	}
	model.ReadOnly = false
	model.exec = func(context.Context, string, string, string) error { return nil }
	model.OnRune('s')
	if !model.ExecRequested {
		t.Fatal("expected an exec request")
	}
	model.RolloutText = formatRollout(kubernetes.RolloutStatus{
		Kind: "Deployment", Namespace: "production", Name: "api", Strategy: "RollingUpdate",
		Replicas: 3, Ready: 3, Updated: 3, Available: 3,
		Revisions: []kubernetes.RolloutRevision{{Revision: 4, Status: "Available", Progress: 100}},
	})
	model.openWorkspace("rollout")
	if !strings.Contains(model.WorkspaceText(), "revision 4") {
		t.Fatal(model.WorkspaceText())
	}
}

func (m Model) SelectedName() string {
	row, ok := m.Selected()
	if !ok {
		return ""
	}
	return row.Name
}
