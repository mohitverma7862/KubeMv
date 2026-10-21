package tui

import (
	"strings"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

var workspaceTabs = []string{"logs", "events", "yaml", "describe", "terminal", "rollout", "history"}

// NavLine is the resource strip from the reference.
func (m Model) NavLine() string {
	var parts []string
	for i, kind := range m.Kinds {
		if !kind.primary() && !m.More && i != m.Active {
			continue
		}
		label := kind.Label
		if kind.Mnemonic != "" {
			label = "[" + kind.Mnemonic + "] " + kind.Label
		}
		if i == m.Active {
			label = "* " + label
		}
		if !kind.Available {
			label += " (unavailable)"
		}
		parts = append(parts, label)
	}
	if !m.More {
		parts = append(parts, "[>] More")
	}
	return strings.Join(parts, "  ")
}

func (kind Kind) primary() bool {
	switch kind.Path {
	case kubernetes.ResourcePods, kubernetes.ResourceDeployments, kubernetes.ResourceServices, kubernetes.ResourceNodes,
		kubernetes.ResourceEvents, kubernetes.ResourceJobs, kubernetes.ResourceCronJobs, kubernetes.ResourceConfigMaps, kubernetes.ResourceSecrets:
		return true
	default:
		return kind.Path == "" && kind.Label == "Ingresses"
	}
}

// TabLine is the lower workspace tab strip.
func (m Model) TabLine() string {
	active := m.workspace()
	parts := make([]string, len(workspaceTabs))
	for i, tab := range workspaceTabs {
		title := strings.ToUpper(tab[:1]) + tab[1:]
		if tab == "yaml" {
			title = "YAML"
		}
		if tab == active {
			title = "[" + title + "]"
		}
		parts[i] = title
	}
	return strings.Join(parts, " | ")
}

// QuickHelp is the persistent help block, drawn from the shortcut registry.
func QuickHelp() string {
	return formatShortcuts(ShortcutsFor("global"))
}

// ContextualHelp adds the shortcuts that apply to the current focus.
func ContextualHelp(commanding bool) string {
	scopes := []string{"row"}
	if commanding {
		scopes = append(scopes, "command")
	}
	return formatShortcuts(ShortcutsFor(scopes...))
}

func formatShortcuts(shortcuts []Shortcut) string {
	lines := make([]string, len(shortcuts))
	for i, shortcut := range shortcuts {
		lines[i] = shortcut.Keys + strings.Repeat(" ", max(1, 10-len(shortcut.Keys))) + shortcut.Label
	}
	return strings.Join(lines, "\n")
}

// FilterLine is the filter shown in the resource workspace.
func (m Model) FilterLine() string {
	if m.Filter == "" && !m.Filtering {
		return "/filter"
	}
	return "/filter: " + m.Filter
}

// InspectorText is the right-hand selected resource panel.
func (m Model) InspectorText() string {
	return m.inspectorBody()
}

// TreeText is the owner chain beside the lower workspace.
func (m Model) TreeText() string {
	return m.treeBody()
}

// WorkspaceText is the lower terminal workspace body.
func (m *Model) WorkspaceText() string {
	return m.workspaceBody()
}

func emptyKind(kind string) string {
	if kind == "" {
		return "Pod"
	}
	return kind
}

// PaletteText lists the shared command registry.
func PaletteText() string {
	lines := make([]string, len(Commands))
	for i, command := range Commands {
		line := ":" + command.Name
		if len(command.Aliases) > 0 {
			line += "  (" + strings.Join(command.Aliases, ", ") + ")"
		}
		line += "  " + command.Summary
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}
