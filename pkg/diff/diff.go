package diff

import (
	"bufio"
	"io"
)

type Diff struct {
	Line string
	Type string
}

func Files(file1, file2 io.Reader) ([]Diff, error) {
	lines1, err := readLines(file1)
	if err != nil {
		return nil, err
	}
	lines2, err := readLines(file2)
	if err != nil {
		return nil, err
	}

	return LineByLine(lines1, lines2), nil
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func LineByLine(lines1, lines2 []string) []Diff {
	var diffs []Diff
	i, j := 0, 0
	for i < len(lines1) || j < len(lines2) {
		if i < len(lines1) && j < len(lines2) && lines1[i] == lines2[j] {
			diffs = append(diffs, Diff{Line: lines1[i], Type: "same"})
			i++
			j++
		} else if i < len(lines1) && (j >= len(lines2) || lines1[i] < lines2[j]) { // KEY CHANGE
			diffs = append(diffs, Diff{Line: lines1[i], Type: "remove"})
			i++
		} else { // Simplified the else condition
			diffs = append(diffs, Diff{Line: lines2[j], Type: "add"})
			j++
		}
	}
	return diffs
}
