package tui

import (
	"strings"
	"unicode"
)

// lowerRuneBound reports runes already used by table/workspace actions (lowercase only).
func lowerRuneBound(r rune) bool {
	switch r {
	case 'j', 'k', 'g', 'l', 'd', 'y', 's', '/', ':', '?', '>':
		return true
	default:
		return false
	}
}

// trySelectMnemonic switches the active resource when the key matches a nav [Mnemonic].
func (m *Model) trySelectMnemonic(r rune) bool {
	if m.TreeFocus || m.Commanding || m.Filtering || m.Confirm || m.Palette {
		return false
	}
	if unicode.IsLower(r) && lowerRuneBound(r) {
		return false
	}
	letter := strings.ToUpper(string(r))
	for i, kind := range m.Kinds {
		if kind.Mnemonic == "" || !strings.EqualFold(kind.Mnemonic, letter) {
			continue
		}
		if !kind.Available {
			m.CommandError = kind.Label + " unavailable"
			return true
		}
		if kind.Path == "" {
			return false
		}
		m.Active = i
		m.Cursor = 0
		m.CommandError = ""
		return true
	}
	return false
}

// trySelectDigit switches to the nth available resource (1–9), k9s-style.
func (m *Model) trySelectDigit(r rune) bool {
	if r < '1' || r > '9' {
		return false
	}
	if m.TreeFocus || m.Commanding || m.Filtering || m.Confirm || m.Palette {
		return false
	}
	want := int(r - '1')
	seen := 0
	for i, kind := range m.Kinds {
		if !kind.Available || kind.Path == "" {
			continue
		}
		if seen == want {
			m.Active = i
			m.Cursor = 0
			m.CommandError = ""
			return true
		}
		seen++
	}
	return false
}
