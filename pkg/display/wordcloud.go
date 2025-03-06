package display

import (
	"fmt"
	"strings"

	"github.com/san-kum/diff-dance/pkg/diff"
)

func WordCloud(diffs []diff.Diff) {
	addedWords, removedWords := diff.WordDiff(diffs)

	fmt.Println("Added Words: ")
	printWordCounts(addedWords, green)

	fmt.Println("\nRemoved Words: ")
	printWordCounts(removedWords, green)
}

func printWordCounts(wordCounts []diff.WordCount, colorFunc func(string) string) {
	for _, wc := range wordCounts {
		// size := 10 + wc.Count*2

		padding := strings.Repeat(" ", wc.Count)
		fmt.Printf("%s%s (Count: %d)\n", colorFunc(fmt.Sprintf("%s%s", padding, wc.Word)), reset(), wc.Count)
	}
}

func reset() string {
	return "\033[0m"
}
