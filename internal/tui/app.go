package tui

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mohitverma7862/KubeMv/internal/kubernetes"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

// Run opens the Fast Mode console in this terminal, the same way k9s does.
func Run() error {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return fmt.Errorf("kubemv opens a terminal console; run it in a terminal")
	}
	model := Load(context.Background())
	return run(&model)
}

func run(model *Model) error {
	app := tview.NewApplication()
	header := tview.NewTextView()
	nav := tview.NewTextView()
	table := tview.NewTable().SetSelectable(false, false).SetFixed(1, 0)
	inspector := tview.NewTextView().SetWrap(true)
	status := tview.NewTextView()
	help := tview.NewTextView().SetText("j/k or arrows move\ng / G first and last\n1-9 switch resource\n/ filter   Esc clear\n? hide help   q quit")

	body := tview.NewFlex().
		AddItem(table, 0, 3, true).
		AddItem(inspector, 36, 0, false)
	pages := tview.NewPages().
		AddPage("main", body, true, true).
		AddPage("help", help, true, false)
	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(nav, 1, 0, false).
		AddItem(pages, 0, 1, true).
		AddItem(status, 1, 0, false)

	paint := func() {
		header.SetText(model.HeaderLine(time.Now()))
		nav.Clear()
		for i, kind := range model.Kinds {
			mark := " "
			if i == model.Active {
				mark = "*"
			}
			fmt.Fprintf(nav, "%s%d %s  ", mark, i+1, kind.Label)
		}
		table.Clear()
		columns := columnsFor(model.kind().Path)
		for col, title := range columns {
			table.SetCell(0, col, tview.NewTableCell(title).SetSelectable(false).SetTextColor(tcell.ColorCadetBlue))
		}
		visible := model.Visible()
		for rowIndex, item := range visible {
			for col, value := range cells(item, columns) {
				cell := tview.NewTableCell(value)
				if rowIndex == model.Cursor {
					cell.SetBackgroundColor(tcell.ColorDarkBlue)
				}
				table.SetCell(rowIndex+1, col, cell)
			}
		}
		if len(visible) == 0 {
			table.SetCell(1, 0, tview.NewTableCell("No resources").SetTextColor(tcell.ColorGray))
		}
		selected, ok := model.Selected()
		if !ok {
			inspector.SetText("Resource Inspector\n\nNo resource selected.")
		} else {
			inspector.SetText(fmt.Sprintf("Resource Inspector\n\nName       %s\nNamespace  %s\nStatus     %s\nAge        %s\nNode       %s\nIP         %s\nReady      %s\nRestarts   %s",
				selected.Name, empty(selected.Namespace), selected.Status, selected.Age,
				field(selected, "node"), field(selected, "ip"), field(selected, "ready"), field(selected, "restarts")))
		}
		status.SetText(model.StatusLine())
		if model.Help {
			pages.ShowPage("help")
		} else {
			pages.HidePage("help")
		}
	}
	paint()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		quit := false
		switch event.Key() {
		case tcell.KeyCtrlC:
			quit = model.OnKey("ctrlc")
		case tcell.KeyEscape:
			quit = model.OnKey("esc")
		case tcell.KeyEnter:
			quit = model.OnKey("enter")
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			quit = model.OnKey("backspace")
		case tcell.KeyDown:
			quit = model.OnKey("down")
		case tcell.KeyUp:
			quit = model.OnKey("up")
		case tcell.KeyRune:
			quit = model.OnRune(event.Rune())
		}
		if quit {
			app.Stop()
			return nil
		}
		paint()
		return nil
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			app.QueueUpdateDraw(paint)
		}
	}()

	return app.SetRoot(root, true).Run()
}

func columnsFor(kind kubernetes.ResourceKindPath) []string {
	if kind == kubernetes.ResourcePods {
		return []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP", "NODE"}
	}
	return []string{"NAMESPACE", "NAME", "STATUS", "AGE"}
}

func cells(row kubernetes.ResourceRow, columns []string) []string {
	out := make([]string, len(columns))
	for i, column := range columns {
		switch column {
		case "NAME":
			out[i] = row.Name
		case "NAMESPACE":
			out[i] = empty(row.Namespace)
		case "STATUS":
			out[i] = row.Status
		case "AGE":
			out[i] = row.Age
		case "READY":
			out[i] = field(row, "ready")
		case "RESTARTS":
			out[i] = field(row, "restarts")
		case "IP":
			out[i] = field(row, "ip")
		case "NODE":
			out[i] = field(row, "node")
		default:
			out[i] = ""
		}
	}
	return out
}

func field(row kubernetes.ResourceRow, key string) string {
	if row.Extra == nil || row.Extra[key] == "" {
		return "—"
	}
	return row.Extra[key]
}
