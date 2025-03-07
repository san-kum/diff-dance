package display

import (
	"fmt"
	"io"
	"strings"

	"github.com/san-kum/diff-dance/pkg/diff"
	"github.com/san-kum/diff-dance/pkg/utils"
)

func Heatmap(diffs []diff.Diff, file1Lines, file2Lines []string, w io.Writer) {
	maxLength := utils.Max(len(file1Lines), len(file2Lines))

	heat := make([]int, maxLength)

	for _, d := range diffs {
		switch d.Type {
		case "add":
			for i, line := range file2Lines {
				if line == d.Line {
					heat[i]++
					break
				}
			}

		case "remove":
			for i, line := range file1Lines {
				if line == d.Line {
					heat[i]++
					break
				}
			}
		}
	}

	for i := 0; i < maxLength; i++ {
		heatColor := heatColor(heat[i])

		var line1, line2 string
		if i < len(file1Lines) {
			line1 = file1Lines[i]
		}
		if i < len(file2Lines) {
			line2 = file2Lines[i]
		}

		if strings.Contains(line1, "\t") {
			line1 = strings.ReplaceAll(line1, "\t", "    ")
		}
		if strings.Contains(line2, "\t") {
			line2 = strings.ReplaceAll(line2, "\t", "    ")
		}

		displayLine := line1
		if displayLine == "" {
			displayLine = line2
		}

		fmt.Fprintf(w, "%s%s\033[0m\n", heatColor, displayLine)
	}
}

func heatColor(heat int) string {
	switch {
	case heat == 0:
		return "\033[48;5;232m" // Darkest gray for no changes
	case heat == 1:
		return "\033[48;5;235m" // Slightly lighter gray
	case heat == 2:
		return "\033[48;5;238m" // Lighter gray
	case heat == 3:
		return "\033[48;5;241m" // Even lighter
	case heat > 3:
		return "\033[48;5;196m" // Red for high change density
	default:
		return "" // Should not happen, but safe default
	}
}
