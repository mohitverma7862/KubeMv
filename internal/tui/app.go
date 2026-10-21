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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	model := Load(ctx)
	return run(ctx, &model)
}

func run(ctx context.Context, model *Model) error {
	app := tview.NewApplication()
	tview.Styles.PrimitiveBackgroundColor = tcell.NewRGBColor(2, 8, 23)
	tview.Styles.ContrastBackgroundColor = tcell.NewRGBColor(8, 67, 122)
	tview.Styles.BorderColor = tcell.NewRGBColor(18, 62, 115)
	tview.Styles.TitleColor = tcell.NewRGBColor(76, 201, 255)
	tview.Styles.PrimaryTextColor = tcell.NewRGBColor(185, 200, 232)

	header := tview.NewTextView()
	nav := tview.NewTextView()
	table := tview.NewTable().SetSelectable(false, false).SetFixed(1, 0)
	inspector := tview.NewTextView().SetWrap(true)
	inspector.SetBorder(true).SetTitle(" Selected Resource ")
	tabs := tview.NewTextView()
	workspace := tview.NewTextView().SetWrap(true)
	workspace.SetBorder(true).SetTitle(" Workspace ")
	tree := tview.NewTextView().SetWrap(true)
	tree.SetBorder(true).SetTitle(" Relationship ")
	quick := tview.NewTextView().SetWrap(true)
	quick.SetBorder(true).SetTitle(" Quick Help ")
	status := tview.NewTextView()
	palette := tview.NewTextView().SetText(PaletteText())
	palette.SetBorder(true).SetTitle(" Command palette ")

	upper := tview.NewFlex().
		AddItem(table, 0, 3, true).
		AddItem(inspector, 36, 0, false)
	lower := tview.NewFlex().
		AddItem(workspace, 0, 3, false).
		AddItem(tree, 0, 2, false).
		AddItem(quick, 0, 2, false)
	stack := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(upper, 0, 3, true).
		AddItem(tabs, 1, 0, false).
		AddItem(lower, 0, 2, false)
	pages := tview.NewPages().
		AddPage("main", stack, true, true).
		AddPage("palette", palette, true, false)
	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 0, false).
		AddItem(nav, 1, 0, false).
		AddItem(pages, 0, 1, true).
		AddItem(status, 1, 0, false)

	var frame TableSync
	watch := newWatchCtrl()
	var paint func()
	paint = func() {
		header.SetText(model.HeaderLine(time.Now()))
		nav.SetText(model.NavLine())
		tabs.SetText(model.TabLine())
		workspace.SetText(model.WorkspaceText())
		tree.SetText(model.TreeText())
		if model.Help {
			quick.SetText(ContextualHelp(model.Commanding))
		} else {
			quick.SetText(QuickHelp())
		}
		table.SetTitle(" " + model.TableTitle(time.Now()) + " ")
		table.SetBorder(true)
		columns := columnsFor(model.kind().Path)
		visible := model.Visible()
		plan := frame.Plan(columns, visible, model.Cursor)
		if plan.Full {
			table.Clear()
			for col, title := range columns {
				table.SetCell(0, col, tview.NewTableCell(title).SetSelectable(false).SetTextColor(tcell.ColorCadetBlue))
			}
		}
		for _, rowIndex := range plan.Paint {
			if rowIndex >= len(visible) {
				continue
			}
			for col, value := range cells(visible[rowIndex], columns) {
				cell := tview.NewTableCell(value)
				if columns[col] == "STATUS" {
					cell.SetTextColor(statusColor(value))
				}
				if rowIndex == model.Cursor {
					cell.SetBackgroundColor(tcell.NewRGBColor(8, 67, 122))
				}
				table.SetCell(rowIndex+1, col, cell)
			}
		}
		for i := 0; i < plan.Drop; i++ {
			if table.GetRowCount() > 1 {
				table.RemoveRow(table.GetRowCount() - 1)
			}
		}
		if len(visible) == 0 {
			table.SetCell(1, 0, tview.NewTableCell("No resources").SetTextColor(tcell.ColorGray))
		}
		inspector.SetText(model.InspectorText())
		status.SetText(model.StatusLine())
		if model.Palette {
			pages.ShowPage("palette")
		} else {
			pages.HidePage("palette")
		}
		watch.Update(model)
	}
	paint()
	go serveWatch(ctx, app, model, watch, paint)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		quit := false
		switch event.Key() {
		case tcell.KeyCtrlC:
			quit = model.OnKey("ctrlc")
		case tcell.KeyCtrlP:
			quit = model.OnKey("ctrlp")
		case tcell.KeyCtrlR:
			quit = model.OnKey("ctrlr")
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
		case tcell.KeyTab:
			model.OnKey("tab")
			paint()
			return nil
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

func statusColor(status string) tcell.Color {
	switch status {
	case "Running", "Available", "Ready", "Complete":
		return tcell.NewRGBColor(24, 213, 155)
	case "Pending", "Progressing", "Terminating":
		return tcell.NewRGBColor(240, 220, 85)
	case "CrashLoopBackOff", "Failed", "NotReady":
		return tcell.NewRGBColor(255, 92, 104)
	default:
		return tcell.ColorWhite
	}
}

func columnsFor(kind kubernetes.ResourceKindPath) []string {
	if kind == kubernetes.ResourcePods {
		return []string{"NAME", "READY", "STATUS", "RESTARTS", "AGE", "IP", "NODE", "CPU", "MEM"}
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
		case "CPU":
			out[i] = field(row, "cpu")
		case "MEM":
			out[i] = field(row, "mem")
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
