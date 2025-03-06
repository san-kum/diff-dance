package display

import (
	"fmt"

	"github.com/san-kum/diff-dance/pkg/diff"
)

func Structural(diffs []diff.StructuralDiff) {
	for _, d := range diffs {
		switch d.Type {
		case "add_func":
			fmt.Printf("Added function: %s%s()\033[0m\n", green("+ "), d.FuncName)
		case "remove_func":
			fmt.Printf("Removed function: %s%s()\033[0m\n", red("- "), d.FuncName)
		case "change_func_sig":
			fmt.Printf("Changed function signature: %s\n", yellow(d.FuncName))
			fmt.Printf("  Old: %s\n", red(d.OldSig))
			fmt.Printf("  New: %s\n", green(d.NewSig))
		}
	}
}
