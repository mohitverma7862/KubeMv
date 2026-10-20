package tui

import (
	"strings"
	"time"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// Kind is one resource tab in the operator console.
type Kind struct {
	Label string
	Path  kubernetes.ResourceKindPath
}

// Model is the terminal console state. It does not draw.
type Model struct {
	Cluster    string
	Context    string
	Namespace  string
	User       string
	Production bool
	ReadOnly   bool
	Live       bool
	Note       string
	Kinds      []Kind
	Active     int
	Rows       map[kubernetes.ResourceKindPath][]kubernetes.ResourceRow
	Cursor     int
	Filter     string
	Filtering  bool
	Help       bool
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
	needle := strings.ToLower(m.Filter)
	out := make([]kubernetes.ResourceRow, 0, len(rows))
	for _, row := range rows {
		haystack := strings.ToLower(row.Name + " " + row.Namespace + " " + row.Status + " " + row.Kind)
		if strings.Contains(haystack, needle) {
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
	if m.Filtering {
		switch r {
		case '\x7f', '\b':
			if m.Filter != "" {
				m.Filter = m.Filter[:len(m.Filter)-1]
			}
		default:
			if r >= 32 {
				m.Filter += string(r)
			}
		}
		m.clamp()
		return false
	}
	switch r {
	case 'q', 'Q':
		return true
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
	case '?':
		m.Help = !m.Help
	default:
		if r >= '1' && r <= '9' {
			index := int(r - '1')
			if index < len(m.Kinds) {
				m.Active = index
				m.Cursor = 0
			}
		}
	}
	m.clamp()
	return false
}

// OnKey handles named keys from the terminal.
func (m *Model) OnKey(name string) (quit bool) {
	switch name {
	case "ctrlc":
		return true
	case "esc":
		if m.Help {
			m.Help = false
			return false
		}
		if m.Filtering {
			m.Filtering = false
			m.Filter = ""
			m.clamp()
		}
		return false
	case "enter":
		m.Filtering = false
		m.clamp()
		return false
	case "backspace":
		if m.Filtering && m.Filter != "" {
			m.Filter = m.Filter[:len(m.Filter)-1]
			m.clamp()
		}
		return false
	case "down":
		if !m.Filtering {
			m.Cursor++
			m.clamp()
		}
	case "up":
		if !m.Filtering {
			m.Cursor--
			m.clamp()
		}
	}
	return false
}

// HeaderLine is the top status line.
func (m Model) HeaderLine(now time.Time) string {
	prod := ""
	if m.Production {
		prod = "  PRODUCTION"
	}
	read := ""
	if m.ReadOnly {
		read = "  READ ONLY"
	}
	conn := "DISCONNECTED"
	if m.Live {
		conn = "LIVE"
	} else if m.Note != "" {
		conn = "SAMPLE"
	}
	return "KubeMv  Fast Mode  " + m.Cluster + "  user:" + empty(m.User) + "  context:" + empty(m.Context) + "  ns:" + empty(m.Namespace) + prod + read + "  " + conn + "  " + now.Format("15:04:05")
}

// StatusLine is the bottom command line.
func (m Model) StatusLine() string {
	filter := "/"
	if m.Filter != "" || m.Filtering {
		filter = "/" + m.Filter
	}
	return "j/k move  g/G ends  1-9 resource  " + filter + " filter  ? help  q quit  " + m.Note
}

func empty(value string) string {
	if value == "" {
		return "—"
	}
	return value
}
