package tui

import (
	"context"
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
	Cluster          string
	Context          string
	Namespace        string
	User             string
	Production       bool
	ReadOnly         bool
	Live             bool
	Note             string
	Kinds            []Kind
	Active           int
	Rows             map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow
	Cursor           int
	Filter           string
	Filtering        bool
	Help             bool
	Workspace        string
	Commanding       bool
	Command          string
	CommandError     string
	Palette          bool
	Confirm          bool
	Forwards         int
	More             bool
	Connection       string
	ReconnectAttempt int
	RetryAfter       time.Duration
	SyncedAt         time.Time
	ResourceVersion  string
	WatchError       string
	RefreshCount     int
	Relations        map[string]kubernetes.ResourceRow
	TreeCursor       int
	TreeFocus        bool
	back             []string
	watch            func(kind kubernetes.ResourceKindPath, namespace string) (kubernetes.WatchSource, bool)
	fetchOwner       func(context.Context, kubernetes.OwnerRef) (kubernetes.ResourceRow, error)
	ownerBusy        map[string]bool
	ownerMiss        map[string]bool
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
	out := make([]kubernetes.ResourceRow, 0, len(rows))
	for _, row := range rows {
		if MatchFilter(row, m.Filter) {
			out = append(out, row)
		}
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
				m.suggest()
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
		if m.TreeFocus {
			m.TreeCursor++
			m.clampTree()
			return false
		}
		m.Cursor++
	case 'k':
		if m.TreeFocus {
			m.TreeCursor--
			m.clampTree()
			return false
		}
		m.Cursor--
	case 'g':
		if m.TreeFocus {
			m.TreeCursor = 0
			return false
		}
		m.Cursor = 0
	case 'G':
		if m.TreeFocus {
			m.TreeCursor = len(m.treeNodes()) - 1
			m.clampTree()
			return false
		}
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
		m.RefreshCount++
		m.ownerMiss = nil
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
		if m.TreeFocus {
			m.TreeFocus = false
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
		if m.TreeFocus {
			m.drillTree()
			return false
		}
		m.openWorkspace("describe")
		return false
	case "backspace":
		if m.Commanding && m.Command != "" {
			m.Command = m.Command[:len(m.Command)-1]
			m.suggest()
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
	case "tab":
		if m.Commanding {
			m.complete()
			return false
		}
		if !m.Filtering && !m.Confirm {
			m.TreeFocus = !m.TreeFocus
			m.clampTree()
		}
	}
	return false
}

func (m *Model) commandToken() string {
	fields := strings.Fields(m.Command)
	if len(fields) == 0 || strings.HasSuffix(m.Command, " ") {
		return ""
	}
	return fields[0]
}

func (m *Model) suggest() {
	token := m.commandToken()
	if token == "" {
		m.CommandError = ""
		return
	}
	m.CommandError = suggestionText(Complete(token), token)
}

func (m *Model) complete() {
	matches := Complete(m.commandToken())
	if len(matches) == 1 {
		m.Command = matches[0].Name
		if matches[0].Action == "ns" || matches[0].Action == "ctx" {
			m.Command += " "
		}
		m.CommandError = ""
		return
	}
	m.CommandError = suggestionText(matches, m.commandToken())
}

func suggestionText(matches []Command, token string) string {
	if len(matches) == 0 {
		return "no match"
	}
	if len(matches) == 1 && matches[0].Name == strings.ToLower(token) {
		return ""
	}
	names := make([]string, len(matches))
	for i, match := range matches {
		names[i] = match.Name
	}
	return strings.Join(names, " ")
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
	command, ok := Lookup(fields[0])
	if !ok {
		m.CommandError = "unknown command"
		return false
	}
	switch command.Action {
	case "quit":
		return true
	case "kind":
		m.selectPath(command.Path)
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
	return "KubeMv | Fast Mode (K9s-Parity) | " + empty(m.Cluster) + " | user: " + empty(m.User) + " | context: " + empty(m.Context) + " | namespace: " + empty(m.Namespace) + prod + " | " + m.connectionLabel() + " | " + now.Format("15:04:05")
}

func (m Model) connectionLabel() string {
	switch kubernetes.ConnectionState(m.Connection) {
	case kubernetes.ConnectionLive:
		return "● LIVE"
	case kubernetes.ConnectionReconnecting:
		label := "◌ RECONNECTING"
		if m.ReconnectAttempt > 0 {
			label += fmt.Sprintf(" %d", m.ReconnectAttempt)
		}
		if m.RetryAfter > 0 {
			label += " in " + m.RetryAfter.Truncate(time.Second).String()
		}
		return label
	case kubernetes.ConnectionDisconnected:
		return "✕ DISCONNECTED"
	default:
		if m.Live {
			return "● LIVE"
		}
		if m.Note != "" {
			return "● SAMPLE"
		}
		return "✕ DISCONNECTED"
	}
}

// WatchAge is the time since the last successful list or watch event.
func (m Model) WatchAge(now time.Time) string {
	if m.SyncedAt.IsZero() {
		return "◷ —"
	}
	age := now.Sub(m.SyncedAt)
	if age < 0 {
		age = 0
	}
	return "◷ " + age.Truncate(time.Second).String()
}

// TableTitle shows the filter, watch age, and resourceVersion.
func (m Model) TableTitle(now time.Time) string {
	title := m.FilterLine() + "  " + m.WatchAge(now)
	if m.ResourceVersion != "" {
		title += "  rv " + m.ResourceVersion
	}
	if kubernetes.ConnectionState(m.Connection) == kubernetes.ConnectionDisconnected && m.WatchError != "" {
		title += "  " + m.WatchError
	}
	return title
}

// OpenWatch returns the live watch for the active cluster, when one is connected.
func (m *Model) OpenWatch(kind kubernetes.ResourceKindPath, namespace string) (kubernetes.WatchSource, bool) {
	if m.watch == nil || kind == "" {
		return nil, false
	}
	return m.watch(kind, namespace)
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
