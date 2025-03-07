package display

import (
	"fmt"
	"io"
	"strings"

	"github.com/san-kum/diff-dance/pkg/diff"
)

func WordCloud(diffs []diff.Diff, w io.Writer) {
	addedWords, removedWords := diff.WordDiff(diffs)

	fmt.Fprintln(w, "Added Words: ")
	printWordCounts(addedWords, green, w)

	fmt.Fprintln(w, "\nRemoved Words: ")
	printWordCounts(removedWords, green, w)
}

func printWordCounts(wordCounts []diff.WordCount, colorFunc func(string) string, w io.Writer) {
	for _, wc := range wordCounts {
		// size := 10 + wc.Count*2

		padding := strings.Repeat(" ", wc.Count)
		fmt.Fprintf(w, "%s%s (Count: %d)\n", colorFunc(fmt.Sprintf("%s%s", padding, wc.Word)), reset(), wc.Count)
	}
}

func reset() string {
	return "\033[0m"
}
