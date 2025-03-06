package display

import (
	"fmt"

	"github.com/san-kum/diff-dance/pkg/diff"
)

func Terminal(diffs []diff.Diff) {
	for _, d := range diffs {
		switch d.Type {
		case "add":
			fmt.Println(green("+ " + d.Line))
		case "remove":
			fmt.Println(red("- " + d.Line))
		case "change":
			fmt.Println(yellow("~ " + d.Line))
			fmt.Println(" " + d.Line)
		default:
			fmt.Println(d.Line)
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
