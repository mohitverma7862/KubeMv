package tui

import (
	"sort"
	"strings"

	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
)

// TablePlan says which visible rows the terminal table needs to redraw.
// A single watch event paints those rows and leaves the rest of the table in place.
type TablePlan struct {
	Full  bool
	Paint []int
	Drop  int
}

// TableSync remembers the last frame so one changed resource does not clear the table.
type TableSync struct {
	ready   bool
	columns string
	keys    []string
	text    []string
	cursor  int
}

// Plan compares the next visible page with the previous frame.
func (s *TableSync) Plan(columns []string, rows []kubernetes.ResourceRow, cursor int) TablePlan {
	keys := make([]string, len(rows))
	text := make([]string, len(rows))
	for i, row := range rows {
		keys[i] = row.Namespace + "/" + row.Name
		text[i] = strings.Join(cells(row, columns), "\t")
	}
	signature := strings.Join(columns, "|")
	plan := TablePlan{}
	if !s.ready || signature != s.columns {
		plan.Full = true
		plan.Paint = make([]int, len(rows))
		for i := range rows {
			plan.Paint[i] = i
		}
	} else {
		dirty := map[int]struct{}{}
		if cursor != s.cursor {
			if s.cursor >= 0 && s.cursor < len(rows) {
				dirty[s.cursor] = struct{}{}
			}
			if cursor >= 0 && cursor < len(rows) {
				dirty[cursor] = struct{}{}
			}
		}
		shared := len(rows)
		if len(s.keys) < shared {
			shared = len(s.keys)
		}
		for i := 0; i < shared; i++ {
			if s.keys[i] != keys[i] || s.text[i] != text[i] {
				dirty[i] = struct{}{}
			}
		}
		for i := len(s.keys); i < len(rows); i++ {
			dirty[i] = struct{}{}
		}
		if len(s.keys) > len(rows) {
			plan.Drop = len(s.keys) - len(rows)
		}
		for index := range dirty {
			plan.Paint = append(plan.Paint, index)
		}
		sort.Ints(plan.Paint)
	}
	s.ready = true
	s.columns = signature
	s.keys = keys
	s.text = text
	s.cursor = cursor
	return plan
}
