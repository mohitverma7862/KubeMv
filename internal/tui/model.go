package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// Kind is one resource tab in the operator console.
type Kind struct {
	Label     string
	Mnemonic  string
	Path      kubernetes.ResourceKindPath
	Available bool
}

// Model is the terminal console state. It does not draw.
type Model struct {
	Cluster      string
	Context      string
	Namespace    string
	User         string
	Production   bool
	ReadOnly     bool
	Live         bool
	Note         string
	Kinds        []Kind
	Active       int
	Rows         map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow
	Cursor       int
	Filter       string
	Filtering    bool
	Help         bool
	Workspace    string
	Commanding   bool
	Command      string
	CommandError string
	Palette      bool
	Confirm      bool
	Forwards     int
	More         bool
	back         []string
}

// IsProduction reports whether the connected context should show the production badge.
// A nonprod or staging context stays unmarked.
func IsProduction(cluster, contextName string) bool {
	text := strings.ToLower(cluster + " " + contextName)
	text = strings.ReplaceAll(text, "non-prod", "")
	text = strings.ReplaceAll(text, "nonprod", "")
	return strings.Contains(text, "prod")
}

func (m *Model) kind() Kind {
	if len(m.Kinds) == 0 {
		return Kind{}
	}
	if m.Active < 0 || m.Active >= len(m.Kinds) {
		m.Active = 0
	}
	return m.Kinds[m.Active]
}

// Visible returns the filtered rows for the active resource.
func (m *Model) Visible() []kubernetes.ResourceRow {
	rows := m.Rows[m.kind().Path]
	if m.Filter == "" {
		return rows
	}
	status := ""
	text := m.Filter
	if strings.HasPrefix(strings.ToLower(m.Filter), "status=") {
		status = m.Filter[len("status="):]
		text = ""
	}
	needle := strings.ToLower(text)
	out := make([]kubernetes.ResourceRow, 0, len(rows))
	for _, row := range rows {
		if status != "" && !strings.EqualFold(row.Status, status) {
			continue
		}
		haystack := strings.ToLower(row.Name + " " + row.Namespace + " " + row.Status + " " + row.Kind)
		if needle != "" && !strings.Contains(haystack, needle) {
			continue
		}
		out = append(out, row)
	}
	return out
}

func (m *Model) clamp() {
	visible := m.Visible()
	if len(visible) == 0 {
		m.Cursor = 0
		return
	}
	if m.Cursor < 0 {
		m.Cursor = 0
	}
	if m.Cursor >= len(visible) {
		m.Cursor = len(visible) - 1
	}
}

func (m *Model) Selected() (kubernetes.ResourceRow, bool) {
	visible := m.Visible()
	m.clamp()
	if len(visible) == 0 {
		return kubernetes.ResourceRow{}, false
	}
	return visible[m.Cursor], true
}

// OnRune handles a typed key. quit is true when the operator leaves the console.
func (m *Model) OnRune(r rune) (quit bool) {
	if m.Commanding || m.Filtering || m.Confirm {
		if r >= 32 {
			if m.Commanding {
				m.Command += string(r)
			} else if m.Filtering {
				m.Filter += string(r)
				m.clamp()
			} else {
				m.Command += string(r)
			}
		}
		return false
	}
	switch r {
	case 'j':
		m.Cursor++
	case 'k':
		m.Cursor--
	case 'g':
		m.Cursor = 0
	case 'G':
		m.Cursor = len(m.Visible()) - 1
	case '/':
		m.Filtering = true
		m.Commanding = false
	case ':':
		m.Commanding = true
		m.Command = ""
		m.CommandError = ""
	case '?':
		m.Help = !m.Help
	case 'l':
		m.openWorkspace("logs")
	case 'd':
		m.openWorkspace("describe")
	case 'y':
		m.openWorkspace("yaml")
	case 's':
		m.openWorkspace("terminal")
	case 'D':
		if m.ReadOnly {
			m.CommandError = "readonly: delete is rejected"
			return false
		}
		m.Confirm = true
		m.Command = ""
	case '>':
		m.More = !m.More
	}
	m.clamp()
	return false
}

// OnKey handles named keys from the terminal.
func (m *Model) OnKey(name string) (quit bool) {
	switch name {
	case "ctrlc":
		return true
	case "ctrlp":
		m.Palette = !m.Palette
		return false
	case "ctrlr":
		m.Note = "refreshed"
		return false
	case "esc":
		if m.Palette {
			m.Palette = false
			return false
		}
		if m.Confirm {
			m.Confirm = false
			m.Command = ""
			return false
		}
		if m.Help {
			m.Help = false
			return false
		}
		if m.Commanding {
			m.Commanding = false
			m.Command = ""
			m.CommandError = ""
			return false
		}
		if m.Filtering {
			m.Filtering = false
			m.Filter = ""
			m.clamp()
			return false
		}
		if len(m.back) > 0 {
			m.Workspace = m.back[len(m.back)-1]
			m.back = m.back[:len(m.back)-1]
		}
		return false
	case "enter":
		if m.Commanding {
			quit = m.applyCommand(m.Command)
			m.Commanding = false
			return quit
		}
		if m.Confirm {
			m.finishConfirm()
			return false
		}
		if m.Filtering {
			m.Filtering = false
			m.clamp()
			return false
		}
		m.openWorkspace("describe")
		return false
	case "backspace":
		if m.Commanding && m.Command != "" {
			m.Command = m.Command[:len(m.Command)-1]
		} else if m.Filtering && m.Filter != "" {
			m.Filter = m.Filter[:len(m.Filter)-1]
			m.clamp()
		} else if m.Confirm && m.Command != "" {
			m.Command = m.Command[:len(m.Command)-1]
		}
		return false
	case "down":
		if !m.Filtering && !m.Commanding {
			m.Cursor++
			m.clamp()
		}
	case "up":
		if !m.Filtering && !m.Commanding {
			m.Cursor--
			m.clamp()
		}
	}
	return false
}

func (m *Model) openWorkspace(tab string) {
	if m.workspace() == tab {
		return
	}
	m.back = append(m.back, m.workspace())
	m.Workspace = tab
}

func (m *Model) workspace() string {
	if m.Workspace == "" {
		return "logs"
	}
	return m.Workspace
}

func (m *Model) applyCommand(line string) bool {
	fields := strings.Fields(strings.TrimSpace(line))
	m.CommandError = ""
	if len(fields) == 0 {
		return false
	}
	switch strings.ToLower(fields[0]) {
	case "q", "quit":
		return true
	case "pods", "po":
		m.selectPath(kubernetes.ResourcePods)
	case "deploy", "deployments", "dp":
		m.selectPath(kubernetes.ResourceDeployments)
	case "svc", "services":
		m.selectPath(kubernetes.ResourceServices)
	case "nodes", "no":
		m.selectPath(kubernetes.ResourceNodes)
	case "events", "ev":
		m.selectPath(kubernetes.ResourceEvents)
	case "jobs", "job":
		m.selectPath(kubernetes.ResourceJobs)
	case "cronjobs", "cronjob":
		m.selectPath(kubernetes.ResourceCronJobs)
	case "cm", "configmap", "configmaps":
		m.selectPath(kubernetes.ResourceConfigMaps)
	case "secrets", "secret":
		m.selectPath(kubernetes.ResourceSecrets)
	case "ns":
		if len(fields) < 2 {
			m.CommandError = "usage: :ns <namespace>"
			return false
		}
		m.Namespace = fields[1]
	case "ctx":
		if len(fields) < 2 {
			m.CommandError = "usage: :ctx <context>"
			return false
		}
		m.Context = fields[1]
		m.Production = IsProduction(m.Cluster, m.Context)
	case "readonly":
		m.ReadOnly = true
	default:
		m.CommandError = "unknown command"
	}
	return false
}

func (m *Model) selectPath(path kubernetes.ResourceKindPath) {
	for i, kind := range m.Kinds {
		if kind.Path == path {
			m.Active = i
			m.Cursor = 0
			m.clamp()
			return
		}
	}
	m.CommandError = "unknown resource"
}

func (m *Model) finishConfirm() {
	if m.Command != "DELETE" {
		m.CommandError = "confirmation rejected"
		m.Confirm = false
		m.Command = ""
		return
	}
	m.Confirm = false
	m.Command = ""
	m.CommandError = "not sent: this console has no mutation boundary"
}

// HeaderLine matches the reference header.
func (m Model) HeaderLine(now time.Time) string {
	prod := ""
	if m.Production {
		prod = " | PRODUCTION"
	}
	conn := "✕ DISCONNECTED"
	if m.Live {
		conn = "● LIVE"
	} else if m.Note != "" {
		conn = "● SAMPLE"
	}
	return "KubeMv | Fast Mode (K9s-Parity) | " + empty(m.Cluster) + " | user: " + empty(m.User) + " | context: " + empty(m.Context) + " | namespace: " + empty(m.Namespace) + prod + " | " + conn + " | " + now.Format("15:04:05")
}

// StatusLine matches the reference command bar.
func (m Model) StatusLine() string {
	command := ": command"
	if m.Commanding {
		command = ":" + m.Command
	}
	if m.CommandError != "" {
		command += "  " + m.CommandError
	}
	filter := "/ filter"
	if m.Filter != "" || m.Filtering {
		filter = "/filter: " + m.Filter
	}
	mode := "R writable"
	if m.ReadOnly {
		mode = "R readonly"
	}
	return command + "   " + filter + "   ? help   ctx " + empty(m.Context) + "   ns " + empty(m.Namespace) + "   " + mode + "   " + fmt.Sprintf("%d forwards", m.Forwards)
}

func empty(value string) string {
	if value == "" {
		return "—"
	}
	return value
}
