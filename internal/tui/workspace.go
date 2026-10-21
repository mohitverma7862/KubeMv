package tui

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (m *Model) workspaceBody() string {
	selected, ok := m.Selected()
	switch m.workspace() {
	case "events":
		if m.EventsText != "" {
			return m.EventsText
		}
		return "No events loaded."
	case "yaml":
		if m.YAML != "" {
			return m.YAML
		}
		if ok && selected.Kind == "Secret" {
			return "apiVersion: v1\nkind: Secret\nmetadata:\n  name: " + selected.Name + "\n  # secret data omitted"
		}
		if !ok {
			return "No YAML loaded."
		}
		return "apiVersion: v1\nkind: " + emptyKind(selected.Kind) + "\nmetadata:\n  name: " + selected.Name + "\n  namespace: " + selected.Namespace
	case "describe":
		if !ok {
			return "No describe loaded."
		}
		if m.Describe != "" {
			return m.Describe
		}
		return describeRow(selected)
	case "terminal":
		return m.terminalText(selected, ok)
	case "rollout":
		if m.RolloutText != "" {
			return m.RolloutText
		}
		return "Rollout is not connected."
	case "history":
		if len(m.Visits) == 0 {
			return "No navigation history."
		}
		return strings.Join(m.Visits, "\n")
	default:
		m.seedLogs(selected, ok)
		return m.Logs.Render(m.LogWidth)
	}
}

func (m Model) WorkspaceTitle() string {
	if m.workspace() == "logs" {
		return " Logs " + m.Logs.Status() + " "
	}
	if m.workspace() == "terminal" {
		name := m.Logs.ContainerName()
		if name == "" {
			name = "container"
		}
		return " Terminal " + name + " "
	}
	return " Workspace "
}

func (m *Model) seedLogs(selected kubernetes.ResourceRow, ok bool) {
	if len(m.Logs.lines) > 0 || m.Logs.cleared || m.source != nil || !ok {
		return
	}
	m.Logs.SetText(selected.Name+"\nlevel=info msg=\"ready\" status="+selected.Status, time.Now())
	if len(selected.Containers) > 0 && len(m.Logs.Containers) == 0 {
		names := make([]string, len(selected.Containers))
		for i, container := range selected.Containers {
			names[i] = container.Name
		}
		m.Logs.SetContainers(names)
	}
}

func (m Model) terminalText(selected kubernetes.ResourceRow, ok bool) string {
	if m.ReadOnly {
		return "readonly: exec is rejected"
	}
	if m.exec == nil {
		return "Terminal is not connected."
	}
	if !ok || selected.Kind != "Pod" {
		return "Exec needs a pod."
	}
	container := m.Logs.ContainerName()
	if container == "" && len(selected.Containers) > 0 {
		container = selected.Containers[0].Name
	}
	return "Exec " + selected.Namespace + "/" + selected.Name + " container " + empty(container) + "\nresize follows the terminal"
}

func describeRow(row kubernetes.ResourceRow) string {
	row = kubernetes.PublicRow(row)
	lines := []string{
		"Name: " + row.Name,
		"Namespace: " + empty(row.Namespace),
		"Kind: " + emptyKind(row.Kind),
		"Status: " + row.Status,
		"Age: " + row.Age,
	}
	if value := fieldIfSet(row, "node"); value != "" {
		lines = append(lines, "Node: "+value)
	}
	if value := fieldIfSet(row, "ip"); value != "" {
		lines = append(lines, "IP: "+value)
	}
	if value := fieldIfSet(row, "qos"); value != "" {
		lines = append(lines, "QoS: "+value)
	}
	if len(row.Containers) > 0 {
		parts := make([]string, len(row.Containers))
		for i, container := range row.Containers {
			parts[i] = container.Name
			if container.Image != "" {
				parts[i] += " " + container.Image
			}
		}
		lines = append(lines, "Containers: "+strings.Join(parts, ", "))
	}
	if len(row.Owners) > 0 {
		lines = append(lines, "Owner: "+row.Owners[0].Kind+"/"+row.Owners[0].Name)
	}
	if labels := labelLine(row); labels != "" {
		lines = append(lines, "Labels: "+labels)
	}
	if row.Kind == "Secret" {
		lines = append(lines, "Type: "+field(row, "type"), "Keys: "+field(row, "keys"), "secret values hidden")
	}
	return strings.Join(lines, "\n")
}

func formatEvents(events []kubernetes.EventRow) string {
	if len(events) == 0 {
		return "No events loaded."
	}
	lines := make([]string, len(events))
	for i, event := range events {
		lines[i] = event.Type + "  " + event.Age + "  " + event.Reason + "  " + event.Message
	}
	return strings.Join(lines, "\n")
}

func formatRollout(status kubernetes.RolloutStatus) string {
	if status.Name == "" {
		return "Rollout is not connected."
	}
	lines := []string{
		status.Kind + " " + status.Namespace + "/" + status.Name,
		fmt.Sprintf("replicas %d  ready %d  updated %d  available %d", status.Replicas, status.Ready, status.Updated, status.Available),
		"strategy " + empty(status.Strategy),
	}
	for _, revision := range status.Revisions {
		lines = append(lines, fmt.Sprintf("revision %d  %s  %d%%", revision.Revision, revision.Status, revision.Progress))
	}
	return strings.Join(lines, "\n")
}

func (m *Model) noteVisit() {
	visible := m.Visible()
	if len(visible) == 0 || m.Cursor < 0 || m.Cursor >= len(visible) {
		return
	}
	row := visible[m.Cursor]
	label := emptyKind(row.Kind) + " " + empty(row.Namespace) + "/" + row.Name
	if n := len(m.Visits); n > 0 && m.Visits[n-1] == label {
		return
	}
	key := row.Namespace + "/" + row.Name
	if key != m.logPod {
		m.logPod = key
		m.Logs = LogBuffer{Containers: containerNames(row), Wrap: m.Logs.Wrap}
	}
	m.Visits = append(m.Visits, label)
	if len(m.Visits) > 40 {
		m.Visits = append([]string(nil), m.Visits[len(m.Visits)-40:]...)
	}
}

func containerNames(row kubernetes.ResourceRow) []string {
	names := make([]string, 0, len(row.Containers))
	for _, container := range row.Containers {
		if container.Name != "" {
			names = append(names, container.Name)
		}
	}
	return names
}

func (m *Model) onWorkspaceRune(r rune) bool {
	if m.TreeFocus || m.Commanding || m.Filtering || m.Confirm || m.Palette {
		return false
	}
	if m.Logs.Searching && m.workspace() == "logs" && r != 'n' && r != 'N' {
		if r >= 32 && r != 127 {
			m.Logs.Query += string(r)
			m.Logs.match = 0
			return true
		}
	}
	if m.workspace() == "logs" && r == '\'' {
		m.Logs.Searching = true
		m.Logs.Query = ""
		m.Logs.match = 0
		return true
	}
	if m.workspace() != "logs" && m.workspace() != "terminal" {
		return false
	}
	if r == 'c' {
		m.Logs.NextContainer()
		m.LogEpoch++
		return true
	}
	if m.workspace() != "logs" {
		return false
	}
	switch r {
	case 'p':
		m.Logs.TogglePause()
	case 'f':
		m.Logs.ToggleFollow()
	case 'P':
		m.Logs.TogglePrevious()
		m.LogEpoch++
	case 't':
		m.Logs.ToggleTimestamps()
	case 'w':
		m.Logs.ToggleWrap()
	case 'n':
		m.Logs.NextMatch()
	case 'N':
		m.Logs.PrevMatch()
	case 'C':
		m.Logs.Clear()
	case 'O':
		m.LogSave = true
	default:
		return false
	}
	return true
}

func (m *Model) logOptions() kubernetes.LogOptions {
	return kubernetes.LogOptions{
		Container:  m.Logs.ContainerName(),
		Previous:   m.Logs.Previous,
		Timestamps: m.Logs.Timestamps,
		TailLines:  500,
		Follow:     m.Logs.Follow && !m.Logs.Paused,
	}
}

// ScheduleWorkspace loads the open tab for the selected resource.
func (m *Model) ScheduleWorkspace(apply func(func())) {
	if m.source == nil || apply == nil {
		return
	}
	row, ok := m.Selected()
	if !ok {
		return
	}
	key := string(m.kind().Path) + "/" + row.Namespace + "/" + row.Name + "/" + m.workspace() + "/" + fmt.Sprint(m.LogEpoch)
	if m.workspace() == "logs" && m.Logs.Follow {
		return
	}
	if m.pageKey == key || m.pageBusy[key] {
		return
	}
	if m.pageBusy == nil {
		m.pageBusy = map[string]bool{}
	}
	m.pageBusy[key] = true
	kind := m.kind().Path
	tab := m.workspace()
	opts := m.logOptions()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		var (
			yamlText string
			events   string
			describe string
			rollout  string
			logs     string
			err      error
		)
		switch tab {
		case "yaml", "describe", "events":
			detail, detailErr := m.source.GetResource(ctx, kind, row.Namespace, row.Name)
			err = detailErr
			if err == nil {
				yamlText = detail.YAML
				events = formatEvents(detail.Events)
				describe = describeRow(detail.Row)
				if events != "No events loaded." {
					describe += "\n\nEvents:\n" + events
				}
			}
		case "logs":
			logs, err = m.source.GetPodLogs(ctx, row.Namespace, row.Name, opts)
		case "rollout":
			targetKind := kind
			targetName := row.Name
			if row.Kind != "Deployment" {
				err = fmt.Errorf("rollout follows a deployment")
				break
			}
			status, statusErr := m.source.GetRolloutStatus(ctx, targetKind, row.Namespace, targetName)
			err = statusErr
			if err == nil {
				rollout = formatRollout(status)
			}
		}
		apply(func() {
			delete(m.pageBusy, key)
			if err != nil {
				m.CommandError = err.Error()
				m.pageKey = key
				return
			}
			m.pageKey = key
			if yamlText != "" {
				m.YAML = yamlText
			}
			if events != "" {
				m.EventsText = events
			}
			if describe != "" {
				m.Describe = describe
			}
			if rollout != "" {
				m.RolloutText = rollout
			}
			if tab == "logs" {
				m.Logs.SetText(logs, time.Now())
			}
		})
	}()
}

// SyncLogFollow starts or stops the log stream. emit runs on the UI goroutine.
func (m *Model) SyncLogFollow(emit func(string)) {
	want := m.source != nil && m.workspace() == "logs" && m.Logs.Follow && !m.Logs.Paused
	row, ok := m.Selected()
	key := ""
	if ok {
		key = row.Namespace + "/" + row.Name + "/" + m.Logs.ContainerName() + "/" + fmt.Sprint(m.Logs.Previous) + "/" + fmt.Sprint(m.LogEpoch)
	}
	if !want || !ok || row.Kind != "Pod" {
		m.stopFollow()
		return
	}
	streamer, ok := m.source.(interface {
		StreamPodLogs(context.Context, string, string, kubernetes.LogOptions) (io.ReadCloser, error)
	})
	if !ok {
		return
	}
	if m.followKey == key && m.followCancel != nil {
		return
	}
	m.stopFollow()
	ctx, cancel := context.WithCancel(context.Background())
	m.followCancel = cancel
	m.followKey = key
	opts := m.logOptions()
	go func() {
		stream, err := streamer.StreamPodLogs(ctx, row.Namespace, row.Name, opts)
		if err != nil {
			emit("")
			return
		}
		defer stream.Close()
		buf := make([]byte, 4096)
		var pending string
		for {
			n, readErr := stream.Read(buf)
			if n > 0 {
				pending += string(buf[:n])
				for {
					index := strings.IndexByte(pending, '\n')
					if index < 0 {
						break
					}
					emit(pending[:index])
					pending = pending[index+1:]
				}
			}
			if readErr != nil {
				if pending != "" {
					emit(pending)
				}
				return
			}
		}
	}()
}

func (m *Model) stopFollow() {
	if m.followCancel != nil {
		m.followCancel()
		m.followCancel = nil
	}
	m.followKey = ""
}
