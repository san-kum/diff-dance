package diff

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"sort"
	"strings"
)

type StructuralDiff struct {
	Type     string
	FuncName string
	OldSig   string
	NewSig   string
}

func StructuralDiffs(file1, file2 io.Reader) ([]StructuralDiff, error) {
	fset := token.NewFileSet()

	// Parse file1
	f1, err := parser.ParseFile(fset, "file1.go", file1, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing file1: %w", err)
	}

	// Parse file2
	f2, err := parser.ParseFile(fset, "file2.go", file2, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing file2: %w", err)
	}

	funcs1 := extractFuncs(f1)
	funcs2 := extractFuncs(f2)

	return compareFuncs(funcs1, funcs2), nil
}

type FuncDecl struct {
	Name string
	Sig  string
}

func extractFuncs(f *ast.File) []FuncDecl {
	var funcs []FuncDecl

	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			funcs = append(funcs, FuncDecl{
				Name: fn.Name.Name,
				Sig:  formatFuncSignature(fn),
			})
		}
	}
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].Name < funcs[j].Name
	})
	return funcs
}

func formatFuncSignature(fn *ast.FuncDecl) string {
	var sig strings.Builder

	sig.WriteString(fn.Name.Name)
	sig.WriteString("(")

	if fn.Type.Params != nil {
		for i, param := range fn.Type.Params.List {
			if i > 0 {
				sig.WriteString(", ")
			}
			for j, name := range param.Names {
				if j > 0 {
					sig.WriteString(", ")
				}
				sig.WriteString(name.Name)
			}
			sig.WriteString(", ")
			sig.WriteString(typeToString(param.Type))
		}
	}
	sig.WriteString(")")

	// Return types

	if fn.Type.Results != nil {
		sig.WriteString(" (")
		for i, result := range fn.Type.Results.List {
			if i > 0 {
				sig.WriteString(", ")
			}
			if len(result.Names) > 0 {
				for j, name := range result.Names {
					if j > 0 {
						sig.WriteString(", ")
					}
					sig.WriteString(name.Name)
				}
				sig.WriteString(" ")
			}
			sig.WriteString(typeToString(result.Type))
		}
		sig.WriteString(")")
	}
	return sig.String()
}

func compareFuncs(funcs1, funcs2 []FuncDecl) []StructuralDiff {

	var diffs []StructuralDiff
	i, j := 0, 0

	for i < len(funcs1) || j < len(funcs2) {
		if i < len(funcs1) && j < len(funcs2) && funcs1[i].Name == funcs2[j].Name {
			if funcs1[i].Sig != funcs2[j].Sig {
				diffs = append(diffs, StructuralDiff{
					Type:     "change_func_sig",
					FuncName: funcs1[i].Name,
					OldSig:   funcs1[i].Sig,
					NewSig:   funcs2[j].Sig,
				})
			}
			i++
			j++
		} else if i < len(funcs1) && (j >= len(funcs2) || funcs1[i].Name < funcs2[j].Name) {
			diffs = append(diffs, StructuralDiff{
				Type:     "remove_func",
				FuncName: funcs1[i].Name,
				OldSig:   funcs1[i].Sig,
			})
			i++
		} else {
			diffs = append(diffs, StructuralDiff{
				Type:     "add_func",
				FuncName: funcs2[j].Name,
				NewSig:   funcs2[j].Sig,
			})
			j++
		}
	}
	return diffs
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "[" + typeToString(t.Len) + "]" + typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType: // Channels
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + typeToString(t.Value)
		case ast.RECV:
			return "<-chan " + typeToString(t.Value)
		default:
			return "chan "
		}
	case *ast.FuncType:
		return "func"
	case *ast.Ellipsis:
		return "..." + typeToString(t.Elt)
	case *ast.BasicLit:
		return t.Value
	default:
		return fmt.Sprintf("unknown_type: %T", t)
	}
}
