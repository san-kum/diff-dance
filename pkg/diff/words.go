package diff

import (
	"regexp"
	"sort"
	"strings"
)

type WordCount struct {
	Word  string
	Count int
}

func WordDiff(diffs []Diff) (addedWords, removedWords []WordCount) {
	addedCounts := make(map[string]int)
	removedCounts := make(map[string]int)

	for _, d := range diffs {
		switch d.Type {
		case "add":
			for _, word := range extractWords(d.Line) {
				addedCounts[word]++
			}
		case "remove":
			for _, word := range extractWords(d.Line) {
				removedCounts[word]++
			}
		}
	}
	addedWords = sortWordCounts(addedCounts)
	removedWords = sortWordCounts(removedCounts)
	return
}

func extractWords(line string) []string {
	re := regexp.MustCompile(`\b\w+\b`)
	words := re.FindAllString(strings.ToLower(line), -1)
	return words
}

func sortWordCounts(counts map[string]int) []WordCount {
	var wordCounts []WordCount
	for word, count := range counts {
		wordCounts = append(wordCounts, WordCount{Word: word, Count: count})
	}
	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].Count > wordCounts[j].Count
	})

	return wordCounts
}
