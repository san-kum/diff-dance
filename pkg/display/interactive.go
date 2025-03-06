package display

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/san-kum/diff-dance/pkg/diff"
	"github.com/san-kum/diff-dance/pkg/utils"
)

func Interactive(file1Path, file2Path string) {
	app := tview.NewApplication()

	// --- Shared Variables ---
	var (
		diffs            []diff.Diff
		file1Lines       []string
		file2Lines       []string
		searchRegex      *regexp.Regexp
		searchText       string
		currentHighlight int
	)

	// Create a TextView to display the diff.
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		}).SetWrap(false)

	// --- Helper Functions ---
	showDetails := func(index int) {
		var detailText string
		if index >= 0 && index < len(diffs) {
			d := diffs[index]
			switch d.Type {
			case "add":
				detailText = fmt.Sprintf("Added:\n%s", d.Line)
			case "remove":
				detailText = fmt.Sprintf("Removed:\n%s", d.Line)
			case "same":
				line1 := ""
				if index < len(file1Lines) {
					line1 = file1Lines[index]
				}
				line2 := ""
				if index < len(file2Lines) {
					line2 = file2Lines[index]
				}
				detailText = fmt.Sprintf("Context:\nFile 1: %s\nFile 2: %s", line1, line2)
			}
		}

		modal := tview.NewModal().
			SetText(detailText).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.SetRoot(textView, true)
				textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
			})

		app.SetRoot(modal, false)
	}

	// performSearch now returns a boolean indicating if a match was found.
	performSearch := func(text string, next bool) bool {
		if text == "" {
			return false // Nothing to search
		}

		// Compile regex only if the search text has changed.
		if text != searchText {
			var err error
			searchRegex, err = regexp.Compile(text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid regular expression: %v\n", err)
				return false
			}
			searchText = text
		}

		start := currentHighlight
		if next {
			start++ // Start searching *after* the current highlight
		}

		// Search forward from the starting position.
		for i := start; i < len(diffs); i++ {
			if searchRegex.MatchString(diffs[i].Line) {
				currentHighlight = i
				textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
				return true // Match found
			}
		}

		// If not found, wrap around to the beginning.
		for i := 0; i < start; i++ {
			if searchRegex.MatchString(diffs[i].Line) {
				currentHighlight = i
				textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
				return true // Match found
			}
		}

		return false // No match found
	}

	setNextHighlight := func() {
		currentHighlight++
		if currentHighlight >= len(diffs) {
			currentHighlight = len(diffs) - 1
		}
		textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
	}

	setPreviousHighlight := func() {
		currentHighlight--
		if currentHighlight < 0 {
			currentHighlight = 0
		}
		textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
	}

	// --- Key Input Handling ---

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyEnter:
			showDetails(currentHighlight)
			return nil
		case tcell.KeyUp:
			setPreviousHighlight()
			return nil
		case tcell.KeyDown:
			setNextHighlight()
			return nil
		case tcell.KeyPgUp:
			_, _, _, height := textView.GetInnerRect()
			scrollUpBy := height / 2
			if scrollUpBy <= 0 {
				scrollUpBy = 1
			}
			for i := 0; i < scrollUpBy; i++ {
				setPreviousHighlight()
			}
			return nil
		case tcell.KeyPgDn:
			_, _, _, height := textView.GetInnerRect()
			scrollDownBy := height / 2
			if scrollDownBy <= 0 {
				scrollDownBy = 1
			}
			for i := 0; i < scrollDownBy; i++ {
				setNextHighlight()
			}
			return nil
		case tcell.KeyHome:
			currentHighlight = 0
			textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
			return nil
		case tcell.KeyEnd:
			currentHighlight = len(diffs) - 1
			textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
			return nil
		case tcell.KeyRune:
			if event.Rune() == '/' {
				var inputField *tview.InputField // Declare outside the closure
				inputField = tview.NewInputField().
					SetLabel("Search: ").
					SetFieldWidth(30).
					SetDoneFunc(func(key tcell.Key) {
						if key == tcell.KeyEnter {
							if !performSearch(inputField.GetText(), false) { // If not found
								// Show a "Not found" message using a modal.
								modal := tview.NewModal().
									SetText("Not found").
									AddButtons([]string{"OK"}).
									SetDoneFunc(func(buttonIndex int, buttonLabel string) {
										app.SetRoot(textView, true) // Put focus back in text
										textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
									})
								app.SetRoot(modal, false)

							} else {
								app.SetRoot(textView, true)
							}

						} else if key == tcell.KeyEscape {
							app.SetRoot(textView, true)
						}
					})

				app.SetRoot(inputField, false)
				app.SetFocus(inputField)
				return nil
			} else if event.Rune() == 'n' && searchText != "" {
				if !performSearch(searchText, true) { //Pass search text
					//Show not found modal
					modal := tview.NewModal().
						SetText("Not found").
						AddButtons([]string{"OK"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							app.SetRoot(textView, true)
							textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()
						})
					app.SetRoot(modal, false)
				}
				return nil
			}
		}
		return event
	})

	// --- Load and Diff Files ---
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

	file1Lines, err = utils.ReadLines(file1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading lines from file1: %v\n", err)
		os.Exit(1)
	}
	file2Lines, err = utils.ReadLines(file2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading lines from file2: %v\n", err)
		os.Exit(1)
	}

	file1.Seek(0, 0)
	file2.Seek(0, 0)

	diffs, err = diff.Files(file1, file2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error diffing files: %v\n", err)
		os.Exit(1)
	}

	diffText := buildInteractiveDiffText(diffs)
	textView.SetText(diffText)
	textView.Highlight(strconv.Itoa(currentHighlight)).ScrollToHighlight()

	// --- Set Up UI and Run ---

	app.SetRoot(textView, true).SetFocus(textView)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}

// buildInteractiveDiffText builds the diff text with color tags and region markers.
func buildInteractiveDiffText(diffs []diff.Diff) string {
	var builder strings.Builder
	for i, d := range diffs {
		regionTag := strconv.Itoa(i)
		switch d.Type {
		case "add":
			builder.WriteString(fmt.Sprintf(`["%s"][green]+ %s[white][""]`, regionTag, d.Line))
		case "remove":
			builder.WriteString(fmt.Sprintf(`["%s"][red]- %s[white][""]`, regionTag, d.Line))
		case "same":
			builder.WriteString(fmt.Sprintf(`["%s"]  %s[""]`, regionTag, d.Line))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
