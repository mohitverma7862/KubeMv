package tui

import (
	"strings"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// Command is one entry in the shared command registry.
type Command struct {
	Name    string
	Aliases []string
	Summary string
	Path    kubernetes.ResourceKindPath
	Action  string
}

// Shortcut is one entry in the shared shortcut registry.
type Shortcut struct {
	Keys  string
	Label string
	Scope string
}

// Commands is the only command list used by completion, the palette, and dispatch.
var Commands = []Command{
	{Name: "pods", Aliases: []string{"po"}, Summary: "Open pods", Path: kubernetes.ResourcePods, Action: "kind"},
	{Name: "deploy", Aliases: []string{"deployments", "dp"}, Summary: "Open deployments", Path: kubernetes.ResourceDeployments, Action: "kind"},
	{Name: "svc", Aliases: []string{"services"}, Summary: "Open services", Path: kubernetes.ResourceServices, Action: "kind"},
	{Name: "nodes", Aliases: []string{"no"}, Summary: "Open nodes", Path: kubernetes.ResourceNodes, Action: "kind"},
	{Name: "events", Aliases: []string{"ev"}, Summary: "Open events", Path: kubernetes.ResourceEvents, Action: "kind"},
	{Name: "jobs", Aliases: []string{"job"}, Summary: "Open jobs", Path: kubernetes.ResourceJobs, Action: "kind"},
	{Name: "cronjobs", Aliases: []string{"cronjob"}, Summary: "Open cronjobs", Path: kubernetes.ResourceCronJobs, Action: "kind"},
	{Name: "cm", Aliases: []string{"configmap", "configmaps"}, Summary: "Open configmaps", Path: kubernetes.ResourceConfigMaps, Action: "kind"},
	{Name: "secrets", Aliases: []string{"secret"}, Summary: "Open secrets", Path: kubernetes.ResourceSecrets, Action: "kind"},
	{Name: "ns", Summary: "Switch namespace", Action: "ns"},
	{Name: "ctx", Summary: "Switch context", Action: "ctx"},
	{Name: "readonly", Summary: "Reject mutations", Action: "readonly"},
	{Name: "quit", Aliases: []string{"q"}, Summary: "Leave the console", Action: "quit"},
}

// Shortcuts is the only shortcut list used by quick help and the help overlay.
var Shortcuts = []Shortcut{
	{Keys: "j/k", Label: "Select row", Scope: "global"},
	{Keys: "g/G", Label: "First/Last", Scope: "global"},
	{Keys: "Enter", Label: "Drill down", Scope: "global"},
	{Keys: "Esc", Label: "Back", Scope: "global"},
	{Keys: "Ctrl-R", Label: "Refresh", Scope: "global"},
	{Keys: "/", Label: "Filter", Scope: "global"},
	{Keys: ":", Label: "Command mode", Scope: "global"},
	{Keys: "?", Label: "Help", Scope: "global"},
	{Keys: "Ctrl-P", Label: "Palette", Scope: "global"},
	{Keys: "Tab", Label: "Complete command", Scope: "command"},
	{Keys: "l", Label: "Logs", Scope: "row"},
	{Keys: "s", Label: "Shell", Scope: "row"},
	{Keys: "d", Label: "Describe", Scope: "row"},
	{Keys: "y", Label: "YAML", Scope: "row"},
}

// Lookup finds a command by its name or an alias.
func Lookup(token string) (Command, bool) {
	token = strings.ToLower(strings.TrimSpace(token))
	for _, command := range Commands {
		if command.Name == token {
			return command, true
		}
		for _, alias := range command.Aliases {
			if alias == token {
				return command, true
			}
		}
	}
	return Command{}, false
}

// Complete returns commands whose name or alias starts with the token.
func Complete(token string) []Command {
	token = strings.ToLower(strings.TrimSpace(token))
	if token == "" {
		return append([]Command(nil), Commands...)
	}
	var out []Command
	for _, command := range Commands {
		if strings.HasPrefix(command.Name, token) {
			out = append(out, command)
			continue
		}
		for _, alias := range command.Aliases {
			if strings.HasPrefix(alias, token) {
				out = append(out, command)
				break
			}
		}
	}
	return out
}

// ShortcutsFor returns global shortcuts plus the requested scopes.
func ShortcutsFor(scopes ...string) []Shortcut {
	wanted := map[string]bool{"global": true}
	for _, scope := range scopes {
		wanted[scope] = true
	}
	var out []Shortcut
	for _, shortcut := range Shortcuts {
		if wanted[shortcut.Scope] {
			out = append(out, shortcut)
		}
	}
	return out
}
