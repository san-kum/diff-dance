package display

import (
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"

	"github.com/san-kum/diff-dance/pkg/diff"
)

func HTML(diffs []diff.Diff, w io.Writer) error {
	const tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>diff-dance</title>
<style>
body { font-family: monospace; }
.add { color: green; }
.remove { color: red; }
.context { color: black; }
.highlight { background-color: yellow; font-weight: bold; }
</style>
</head>
<body>
<pre>
%s
</pre>
</body>
</html>
  `
	var htmlBuilder strings.Builder

	for _, d := range diffs {
		var line string
		switch d.Type {
		case "add":
			line = fmt.Sprintf(`<span class="add">+ %s</span>`, html.EscapeString(d.Line))
		case "remove":
			line = fmt.Sprintf(`<span class="remove>- %s</span>`, html.EscapeString(d.Line))
		case "same":
			line = fmt.Sprintf(`<span class="context> %s</span>`, html.EscapeString(d.Line))
		}
		htmlBuilder.WriteString(line + "\n")
	}
	_, err := fmt.Fprintf(w, tmpl, htmlBuilder.String())
	return err
}

func HTMLWithHighlight(diffs []diff.Diff, w io.Writer, searchRegex *regexp.Regexp) error {
	const tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>diff-dance</title>
<style>
body { font-family: monospace; }
.add { color: green; }
.remove { color: red; }
.context { color: black; }
.highlight { background-color: yellow; font-weight: bold; }
</style>
</head>
<body>
<pre>
%s
</pre>
</body>
</html>`

	var htmlBuilder strings.Builder

	for _, d := range diffs {
		var line string
		switch d.Type {
		case "add":
			line = fmt.Sprintf(`<span class="add">+ %s</span>`, html.EscapeString(d.Line))
		case "remove":
			line = fmt.Sprintf(`<span class="remove>- %s</span>`, html.EscapeString(d.Line))
		case "same":
			line = fmt.Sprintf(`<span class="context> %s</span>`, html.EscapeString(d.Line))
		}

		if searchRegex != nil {
			line = searchRegex.ReplaceAllStringFunc(line, func(match string) string {
				return fmt.Sprintf(`<span class="highlight">%s</span>`, match)
			})
		}
		htmlBuilder.WriteString(line + "\n")
	}
	_, err := fmt.Fprintf(w, tmpl, htmlBuilder.String())
	return err

}

func HTMLDir(diffs []diff.DirectoryDiff, w io.Writer) error {
	const tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>diff-dance - Directory Diff</title>
<style>
body { font-family: monospace; }
.add { color: green; }
.remove { color: red; }
.change { color: orange; }
.same { color: black; }
.binary { color: magenta; }
.add_dir { color: blue; }
.remove_dir { color: blue; }
.same_dir { color: gray; }
</style>
</head>
<body>
<pre>
%s
</pre>
</body>
</html>`

	var htmlBuilder strings.Builder

	for _, d := range diffs {
		var line string
		switch d.Type {
		case "add":
			line = fmt.Sprintf(`<span class="add">+ %s</span>`, html.EscapeString(d.File2))
		case "remove":
			line = fmt.Sprintf(`<span class="remove">- %s</span>`, html.EscapeString(d.File1))
		case "change":
			if d.BinaryDiff { //Binary diff
				line = fmt.Sprintf(`<span class="binary">~ Binary files differ: %s &lt;-&gt; %s</span>`, html.EscapeString(d.File1), html.EscapeString(d.File2))
			} else { // Normal diff
				line = fmt.Sprintf(`<span class="change">~ %s</span>`+"\n", html.EscapeString(d.File1))
				for _, innerDiff := range d.Diffs { // Iterate the changes
					switch innerDiff.Type {
					case "add":
						line += fmt.Sprintf(`<span class="add">+ %s</span>`+"\n", html.EscapeString(innerDiff.Line))
					case "remove":
						line += fmt.Sprintf(`<span class="remove">- %s</span>`+"\n", html.EscapeString(innerDiff.Line))
					case "same":
						line += fmt.Sprintf(`<span class="context">  %s</span>`+"\n", html.EscapeString(innerDiff.Line))

					}
				}
			}
		case "same":
			line = fmt.Sprintf(`<span class="same">  %s</span>`, html.EscapeString(d.File1))
		case "add_dir":
			line = fmt.Sprintf(`<span class="add_dir">+ Directory added: %s</span>`, html.EscapeString(d.File2))
		case "remove_dir":
			line = fmt.Sprintf(`<span class="remove_dir">- Directory removed: %s</span>`, html.EscapeString(d.File1))
		case "same_dir":
			line = fmt.Sprintf(`<span class="same_dir">  Directory: %s</span>`, html.EscapeString(d.File1))
		}
		htmlBuilder.WriteString(line + "\n")
	}

	_, err := fmt.Fprintf(w, tmpl, htmlBuilder.String())
	return err
}
