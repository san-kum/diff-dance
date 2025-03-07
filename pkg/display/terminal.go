package display

import (
	"fmt"
	"io"

	"github.com/san-kum/diff-dance/pkg/diff"
)

func Terminal(diffs []diff.Diff, w io.Writer) {
	for _, d := range diffs {
		switch d.Type {
		case "add":
			fmt.Fprintln(w, green("+ "+d.Line))
		case "remove":
			fmt.Fprintln(w, red("- "+d.Line))
		case "change":
			fmt.Fprintln(w, yellow("~ "+d.Line))
			fmt.Fprintln(w, " "+d.Line)
		default:
			fmt.Fprintln(w, d.Line)
		}
	}
}

func TerminalDir(diffs []diff.DirectoryDiff, w io.Writer) {
	for _, d := range diffs {
		switch d.Type {
		case "add":
			fmt.Fprintf(w, "%s %s\n", green("+"), d.File2)
		case "remove":
			fmt.Fprintf(w, "%s %s\n", red("-"), d.File1)
		case "change":
			if d.BinaryDiff {
				fmt.Fprintf(w, "%s %s\n", yellow("~"), fmt.Sprintf("Binary files differ: %s <-> %s", d.File1, d.File2))
			} else {
				fmt.Fprintf(w, "%s %s\n", yellow("~"), fmt.Sprintf("File: %s", d.File1))
				Terminal(d.Diffs, w) //Recursive call to show changes
			}
		case "same":
			fmt.Fprintf(w, "  %s\n", d.File1)
		case "add_dir":
			fmt.Fprintf(w, "%s %s\n", green("+"), fmt.Sprintf("Directory added: %s", d.File2))
		case "remove_dir":
			fmt.Fprintf(w, "%s %s\n", red("-"), fmt.Sprintf("Directory removed: %s", d.File1))
		case "same_dir":
			fmt.Fprintf(w, "  %s\n", fmt.Sprintf("Directory: %s", d.File1))
		}
	}
}

func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}
func green(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func yellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}
