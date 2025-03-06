package utils

import (
	"bufio"
	"io"
)

func ReadLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
