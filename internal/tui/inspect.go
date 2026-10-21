package tui

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

func (m Model) inspectorBody() string {
	selected, ok := m.Selected()
	actions := inspectorActions()
	if !ok {
		return "Selected Resource\n\nNo resource selected.\n\n" + actions
	}
	selected = kubernetes.PublicRow(selected)
	body := "Selected Resource\n\n" + selected.Name + "  " + selected.Status
	if selected.Kind == "Secret" {
		body += "\nnamespace  " + empty(selected.Namespace)
		body += "\ntype       " + field(selected, "type")
		body += "\nkeys       " + field(selected, "keys")
		body += "\n\nmetadata only\nsecret values hidden"
	} else {
		body += inspectorLine("namespace", selected.Namespace)
		body += inspectorLine("node", fieldIfSet(selected, "node"))
		body += inspectorLine("IP", fieldIfSet(selected, "ip"))
		body += inspectorLine("ready", fieldIfSet(selected, "ready"))
		body += inspectorLine("restarts", fieldIfSet(selected, "restarts"))
		body += inspectorLine("age", selected.Age)
		body += inspectorLine("qos", fieldIfSet(selected, "qos"))
		body += inspectorLine("cpu", fieldIfSet(selected, "cpu"))
		body += inspectorLine("mem", fieldIfSet(selected, "mem"))
		body += inspectorLine("limits", fieldIfSet(selected, "limits"))
		body += inspectorLine("ports", fieldIfSet(selected, "ports"))
		body += inspectorLine("mounts", fieldIfSet(selected, "mounts"))
		if containers := containerLine(selected); containers != "" {
			body += "\ncontainers " + containers
		}
		if labels := labelLine(selected); labels != "" {
			body += "\nlabels     " + labels
		}
		if len(selected.Owners) > 0 {
			owner := selected.Owners[0]
			body += "\nowner      " + owner.Kind + "/" + owner.Name
		}
	}
	body += "\n\n" + actions
	if m.Confirm {
		body += "\n\nDELETE POD\n" + m.Namespace + " / " + selected.Name + "\nThis action cannot be automatically undone.\nType DELETE to confirm."
	}
	return body
}

func inspectorLine(label, value string) string {
	if value == "" {
		return ""
	}
	return "\n" + label + strings.Repeat(" ", max(1, 11-len(label))) + value
}

func fieldIfSet(row kubernetes.ResourceRow, key string) string {
	if row.Extra == nil || row.Extra[key] == "" {
		return ""
	}
	return row.Extra[key]
}

func containerLine(row kubernetes.ResourceRow) string {
	parts := make([]string, 0, len(row.Containers))
	for _, container := range row.Containers {
		if container.Image == "" {
			parts = append(parts, container.Name)
			continue
		}
		parts = append(parts, container.Name+"  "+container.Image)
	}
	return strings.Join(parts, ", ")
}

func labelLine(row kubernetes.ResourceRow) string {
	if len(row.Labels) == 0 {
		return ""
	}
	keys := make([]string, 0, len(row.Labels))
	for key, value := range row.Labels {
		if !safeValue(key) || !safeValue(value) {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+"="+row.Labels[key])
	}
	return strings.Join(parts, " ")
}

func safeValue(value string) bool {
	text := strings.ToLower(value)
	return !strings.Contains(text, "begin ") && !strings.Contains(text, "private key") && !strings.Contains(text, "bearer ")
}

func inspectorActions() string {
	shortcuts := ShortcutsFor("inspect")
	parts := make([]string, len(shortcuts))
	for i, shortcut := range shortcuts {
		parts[i] = shortcut.Keys + " " + shortcut.Label
	}
	return strings.Join(parts, "  ")
}

func (m Model) treeNodes() []kubernetes.TreeNode {
	selected, ok := m.Selected()
	if !ok {
		return nil
	}
	return kubernetes.BuildTree(selected, m.lookupRow)
}

func (m Model) lookupRow(ref kubernetes.OwnerRef) (kubernetes.ResourceRow, bool) {
	if row, ok := m.Relations[ownerKey(ref)]; ok {
		return row, true
	}
	for _, rows := range m.Rows {
		for _, row := range rows {
			if row.Kind == ref.Kind && row.Name == ref.Name && (ref.Namespace == "" || row.Namespace == ref.Namespace) {
				return row, true
			}
		}
	}
	return kubernetes.ResourceRow{}, false
}

func ownerKey(ref kubernetes.OwnerRef) string {
	return ref.Kind + "/" + ref.Namespace + "/" + ref.Name
}

func (m Model) treeBody() string {
	nodes := m.treeNodes()
	cursor := -1
	if m.TreeFocus {
		cursor = m.TreeCursor
	}
	return kubernetes.FormatTree(nodes, cursor)
}

func (m *Model) clampTree() {
	nodes := m.treeNodes()
	if len(nodes) == 0 {
		m.TreeCursor = 0
		return
	}
	if m.TreeCursor < 0 {
		m.TreeCursor = 0
	}
	if m.TreeCursor >= len(nodes) {
		m.TreeCursor = len(nodes) - 1
	}
}

func (m *Model) drillTree() {
	nodes := m.treeNodes()
	m.clampTree()
	if len(nodes) == 0 {
		return
	}
	node := nodes[m.TreeCursor]
	path, ok := kindPath(node.Kind)
	if !ok || node.Name == "" {
		m.CommandError = "no list for " + node.Kind
		return
	}
	m.Filter = node.Name
	m.Filtering = false
	m.TreeFocus = false
	m.selectPath(path)
}

func kindPath(kind string) (kubernetes.ResourceKindPath, bool) {
	switch kind {
	case "Pod":
		return kubernetes.ResourcePods, true
	case "Deployment":
		return kubernetes.ResourceDeployments, true
	case "Service":
		return kubernetes.ResourceServices, true
	case "Node":
		return kubernetes.ResourceNodes, true
	case "Event":
		return kubernetes.ResourceEvents, true
	case "Job":
		return kubernetes.ResourceJobs, true
	case "CronJob":
		return kubernetes.ResourceCronJobs, true
	case "ConfigMap":
		return kubernetes.ResourceConfigMaps, true
	case "Secret":
		return kubernetes.ResourceSecrets, true
	case "StatefulSet":
		return kubernetes.ResourceStatefulSets, true
	case "DaemonSet":
		return kubernetes.ResourceDaemonSets, true
	default:
		return "", false
	}
}

func (m *Model) nextOwner() (kubernetes.OwnerRef, bool) {
	selected, ok := m.Selected()
	if !ok {
		return kubernetes.OwnerRef{}, false
	}
	current := selected
	seen := map[string]bool{}
	for i := 0; i < 8; i++ {
		if len(current.Owners) == 0 {
			return kubernetes.OwnerRef{}, false
		}
		owner := current.Owners[0]
		key := ownerKey(owner)
		if seen[key] || m.ownerMiss[key] {
			return kubernetes.OwnerRef{}, false
		}
		seen[key] = true
		parent, ok := m.lookupRow(owner)
		if !ok {
			return owner, true
		}
		current = parent
	}
	return kubernetes.OwnerRef{}, false
}

// ScheduleOwnerFetch loads the next missing owner. apply runs the mutation on the UI goroutine.
func (m *Model) ScheduleOwnerFetch(apply func(func())) {
	if m.fetchOwner == nil || apply == nil {
		return
	}
	owner, ok := m.nextOwner()
	if !ok {
		return
	}
	key := ownerKey(owner)
	if m.ownerBusy == nil {
		m.ownerBusy = map[string]bool{}
	}
	if m.ownerBusy[key] {
		return
	}
	m.ownerBusy[key] = true
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		row, err := m.fetchOwner(ctx, owner)
		apply(func() {
			delete(m.ownerBusy, key)
			if err != nil {
				if m.ownerMiss == nil {
					m.ownerMiss = map[string]bool{}
				}
				m.ownerMiss[key] = true
				return
			}
			if m.Relations == nil {
				m.Relations = map[string]kubernetes.ResourceRow{}
			}
			m.Relations[key] = row
		})
	}()
}
