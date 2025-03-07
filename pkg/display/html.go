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
