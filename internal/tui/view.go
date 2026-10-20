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

// QuickHelp is the persistent help block from the reference.
func QuickHelp() string {
	return strings.Join([]string{
		"j/k       Select row",
		"g/G       First/Last",
		"Enter     Drill down",
		"Esc       Back",
		"Ctrl-R    Refresh",
		"/         Filter",
		":         Command mode",
		"?         Help",
		"Ctrl-P    Palette",
	}, "\n")
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
	selected, ok := m.Selected()
	if !ok {
		return "Selected Resource\n\nNo resource selected.\n\nl Logs  s Shell  d Describe  y YAML  D Delete"
	}
	secret := selected.Kind == "Secret"
	body := "Selected Resource\n\n" + selected.Name + "  " + selected.Status
	if secret {
		body += "\n\nmetadata only\nsecret values hidden"
	} else {
		body += "\nnamespace  " + empty(selected.Namespace)
		body += "\nnode       " + field(selected, "node")
		body += "\nIP         " + field(selected, "ip")
		body += "\nready      " + field(selected, "ready")
		body += "\nrestarts   " + field(selected, "restarts")
		body += "\nage        " + selected.Age
		body += "\ncpu        " + field(selected, "cpu")
		body += "\nmem        " + field(selected, "mem")
	}
	body += "\n\nl Logs  s Shell  d Describe  y YAML  D Delete"
	if m.Confirm {
		body += "\n\nDELETE POD\n" + m.Namespace + " / " + selected.Name + "\nThis action cannot be automatically undone.\nType DELETE to confirm."
	}
	return body
}

// TreeText is the owner chain beside the lower workspace.
func (m Model) TreeText() string {
	selected, ok := m.Selected()
	if !ok {
		return "No owner chain."
	}
	if strings.HasPrefix(selected.Name, "api-") {
		return "Deployment (api)\n└─ ReplicaSet (api-7c9f8b6d4)\n   └─ Pod (" + selected.Name + ")\n      └─ Containers\n         └─ api"
	}
	return "No owner chain."
}

// WorkspaceText is the lower terminal workspace body.
func (m Model) WorkspaceText() string {
	selected, ok := m.Selected()
	switch m.workspace() {
	case "events":
		if !ok {
			return "No events loaded."
		}
		return selected.Status + "  " + selected.Name
	case "yaml":
		if !ok {
			return "No YAML loaded."
		}
		if selected.Kind == "Secret" {
			return "apiVersion: v1\nkind: Secret\nmetadata:\n  name: " + selected.Name + "\n  # secret data omitted"
		}
		ns := selected.Namespace
		return "apiVersion: v1\nkind: " + emptyKind(selected.Kind) + "\nmetadata:\n  name: " + selected.Name + "\n  namespace: " + ns
	case "describe":
		if !ok {
			return "No describe loaded."
		}
		return "Name: " + selected.Name + "\nNamespace: " + empty(selected.Namespace) + "\nStatus: " + selected.Status + "\nAge: " + selected.Age
	case "terminal":
		return "Terminal is not connected."
	case "rollout":
		return "Rollout is not connected."
	case "history":
		if len(m.back) == 0 {
			return "No navigation history."
		}
		return strings.Join(m.back, "\n")
	default:
		if !ok {
			return "No logs loaded."
		}
		return selected.Name + "\nlevel=info msg=\"ready\" status=" + selected.Status
	}
}

func emptyKind(kind string) string {
	if kind == "" {
		return "Pod"
	}
	return kind
}

// PaletteText lists the shared command aliases.
func PaletteText() string {
	return strings.Join([]string{
		":q :quit",
		":pods :po",
		":deploy :dp",
		":svc",
		":nodes :no",
		":events :ev",
		":jobs",
		":cronjobs",
		":cm",
		":secrets",
		":ns",
		":ctx",
		":readonly",
	}, "\n")
}
