package tui

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const logLimit = 1000

type logLine struct {
	At   time.Time
	Text string
}

// LogBuffer is the logs workspace: follow, pause, search, and a bounded line list.
type LogBuffer struct {
	lines      []logLine
	Follow     bool
	Paused     bool
	Previous   bool
	Timestamps bool
	Wrap       bool
	Containers []string
	Container  int
	Query      string
	match      int
	Searching  bool
	cleared    bool
}

func (b *LogBuffer) SetContainers(names []string) {
	b.Containers = append([]string(nil), names...)
	if b.Container >= len(b.Containers) {
		b.Container = 0
	}
}

func (b *LogBuffer) ContainerName() string {
	if len(b.Containers) == 0 || b.Container < 0 || b.Container >= len(b.Containers) {
		return ""
	}
	return b.Containers[b.Container]
}

func (b *LogBuffer) NextContainer() string {
	if len(b.Containers) == 0 {
		return ""
	}
	b.Container = (b.Container + 1) % len(b.Containers)
	return b.ContainerName()
}

func (b *LogBuffer) SetText(text string, now time.Time) {
	b.lines = nil
	b.cleared = false
	b.add(text, now)
}

func (b *LogBuffer) Append(text string, now time.Time) {
	if b.Paused {
		return
	}
	b.cleared = false
	b.add(text, now)
}

func (b *LogBuffer) add(text string, now time.Time) {
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			continue
		}
		b.lines = append(b.lines, logLine{At: now, Text: line})
	}
	if len(b.lines) > logLimit {
		b.lines = append([]logLine(nil), b.lines[len(b.lines)-logLimit:]...)
	}
}

func (b *LogBuffer) Stop() { b.Follow = false }

func (b *LogBuffer) Clear() {
	b.lines = nil
	b.match = 0
	b.cleared = true
}

func (b *LogBuffer) TogglePause() { b.Paused = !b.Paused }

func (b *LogBuffer) ToggleFollow() {
	b.Follow = !b.Follow
	if b.Follow {
		b.Paused = false
	}
}

func (b *LogBuffer) TogglePrevious()   { b.Previous = !b.Previous }
func (b *LogBuffer) ToggleTimestamps() { b.Timestamps = !b.Timestamps }
func (b *LogBuffer) ToggleWrap()       { b.Wrap = !b.Wrap }

func (b *LogBuffer) NextMatch() {
	hits := b.matchIndexes()
	if len(hits) == 0 {
		b.match = 0
		return
	}
	b.match = (b.match + 1) % len(hits)
}

func (b *LogBuffer) PrevMatch() {
	hits := b.matchIndexes()
	if len(hits) == 0 {
		return
	}
	b.match = (b.match - 1 + len(hits)) % len(hits)
}

func (b *LogBuffer) matchIndexes() []int {
	if strings.TrimSpace(b.Query) == "" {
		return nil
	}
	query := strings.ToLower(b.Query)
	var hits []int
	for i, line := range b.lines {
		if strings.Contains(strings.ToLower(line.Text), query) {
			hits = append(hits, i)
		}
	}
	return hits
}

func (b *LogBuffer) Status() string {
	state := "stopped"
	if b.Follow && b.Paused {
		state = "paused"
	} else if b.Follow {
		state = "follow"
	}
	if b.Previous {
		state += " previous"
	}
	if name := b.ContainerName(); name != "" {
		state += " " + name
	}
	if b.Timestamps {
		state += " timestamps"
	}
	if b.Wrap {
		state += " wrap"
	}
	if hits := b.matchIndexes(); len(hits) > 0 {
		state += " match " + itoa(b.match+1) + "/" + itoa(len(hits))
	}
	return state
}

func (b *LogBuffer) Render(width int) string {
	if len(b.lines) == 0 {
		return b.Status() + "\nNo logs loaded."
	}
	hits := b.matchIndexes()
	current := -1
	if len(hits) > 0 {
		if b.match >= len(hits) {
			b.match = 0
		}
		current = hits[b.match]
	}
	start := 0
	if b.Follow && b.Query == "" && len(b.lines) > 80 {
		start = len(b.lines) - 80
	}
	body := make([]string, 0, len(b.lines)-start+1)
	body = append(body, b.Status())
	for i := start; i < len(b.lines); i++ {
		line := b.display(b.lines[i])
		if i == current {
			line = "> " + line
		}
		body = append(body, fitLine(line, width, b.Wrap)...)
	}
	return strings.Join(body, "\n")
}

func (b *LogBuffer) Download() string {
	parts := make([]string, len(b.lines))
	for i, line := range b.lines {
		parts[i] = b.display(line)
	}
	return strings.Join(parts, "\n")
}

func (b *LogBuffer) display(line logLine) string {
	if !b.Timestamps || line.At.IsZero() {
		return line.Text
	}
	return line.At.Format("15:04:05") + " " + line.Text
}

func fitLine(line string, width int, wrap bool) []string {
	if width <= 0 || len(line) <= width {
		return []string{line}
	}
	if !wrap {
		return []string{line[:width]}
	}
	var out []string
	for len(line) > width {
		out = append(out, line[:width])
		line = line[width:]
	}
	if line != "" {
		out = append(out, line)
	}
	return out
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var digits []byte
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}

// SaveLogFile writes the current log text. The file mode keeps it private to the operator.
func SaveLogFile(dir, name, text string) (string, error) {
	safe := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, name)
	if safe == "" {
		safe = "logs"
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "kubemv-"+safe+".log")
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		return "", err
	}
	return path, nil
}
