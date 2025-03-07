package display

import (
	"fmt"
	"io"

	"github.com/san-kum/diff-dance/pkg/diff"
)

func Structural(diffs []diff.StructuralDiff, w io.Writer) {
	for _, d := range diffs {
		switch d.Type {
		case "add_func":
			fmt.Fprintf(w, "Added function: %s%s()\033[0m\n", green("+ "), d.FuncName)
		case "remove_func":
			fmt.Fprintf(w, "Removed function: %s%s()\033[0m\n", red("- "), d.FuncName)
		case "change_func_sig":
			fmt.Fprintf(w, "Changed function signature: %s\n", yellow(d.FuncName))
			fmt.Fprintf(w, "  Old: %s\n", red(d.OldSig))
			fmt.Fprintf(w, "  New: %s\n", green(d.NewSig))
		}
	}
}
