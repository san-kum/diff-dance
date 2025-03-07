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

func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}
func green(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func yellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}
