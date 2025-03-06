package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/san-kum/diff-dance/pkg/diff"
)

func Interactive(file1Path, file2Path string) {
	app := tview.NewApplication()

	textView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetChangedFunc(func() {
		app.Draw()
	})

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyEnter:
			return nil
		}
		return event
	})

	// Load and diff the files
	file1, err := os.Open(file1Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file1: %v\n", err)
		os.Exit(1)
	}
	defer file1.Close()

	file2, err := os.Open(file2Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file2: %v\n", err)
		os.Exit(1)
	}
	defer file2.Close()

	diffs, err := diff.Files(file1, file2)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error diffing files: %v\n", err)
		os.Exit(1)
	}

	diffText := buildInteractiveDiffText(diffs)
	textView.SetText(diffText)
	textView.Highlight("0")

	flex := tview.NewFlex().AddItem(textView, 0, 1, true)

	app.SetRoot(flex, true).SetFocus(textView)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}

func buildInteractiveDiffText(diffs []diff.Diff) string {
	var builder strings.Builder
	HighlightIndex := 0

	for _, d := range diffs {
		switch d.Type {
		case "add":
			builder.WriteString(fmt.Sprintf(`["%d"][green]+ %s[white][""]`, HighlightIndex, d.Line))
		case "remove":
			builder.WriteString(fmt.Sprintf(`["%d"][red]- %s[white][""]`, HighlightIndex, d.Line))
		case "same":
			builder.WriteString(fmt.Sprintf(`["%d"] %s[""]`, HighlightIndex, d.Line))
		}
		builder.WriteString("\n")
		HighlightIndex++
	}
	return builder.String()
}
