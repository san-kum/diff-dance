package diff

import (
	"bufio"
	"bytes"
	"fmt"
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

func FilesDiff(file1, file2 io.Reader) ([]Diff, bool, error) {
	// Check if files are likely binary.  If so, don't do line-by-line.
	isBin1, err1 := IsBinary(file1)
	isBin2, err2 := IsBinary(file2)
	//reset pointer
	if f, ok := file1.(io.Seeker); ok {
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, false, fmt.Errorf("seeking in file1: %w", err)
		}
	}
	if f, ok := file2.(io.Seeker); ok {
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, false, fmt.Errorf("seeking in file2: %w", err)
		}
	}
	if err1 != nil {
		return nil, false, fmt.Errorf("checking file1: %w", err1)
	}
	if err2 != nil {
		return nil, false, fmt.Errorf("checking file2: %w", err2)
	}
	if isBin1 || isBin2 {
		//Compare bytes to check if different.
		b1, err1 := io.ReadAll(file1)
		b2, err2 := io.ReadAll(file2)

		if err1 != nil {
			return nil, false, fmt.Errorf("reading file1: %w", err1)
		}
		if err2 != nil {
			return nil, false, fmt.Errorf("reading file2: %w", err2)
		}
		if !bytes.Equal(b1, b2) {
			return nil, true, nil // Indicate binary and different
		}
		return nil, false, nil //Binary and equal
	}

	lines1, err := readLines(file1)
	if err != nil {
		return nil, false, err
	}

	lines2, err := readLines(file2)
	if err != nil {
		return nil, false, err
	}

	return LineByLine(lines1, lines2), false, nil
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

func IsBinary(r io.Reader) (bool, error) {
	// Read a small chunk of the file.
	buf := make([]byte, 512) // Check first 512 bytes
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return false, fmt.Errorf("reading file for binary check: %w", err)
	}

	// Check for null bytes or non-printable ASCII characters.
	if n > 0 {
		if bytes.Contains(buf[:n], []byte{0}) {
			return true, nil
		}
		// Check for a high proportion of non-printable characters.
		nonPrintableCount := 0
		for _, b := range buf[:n] {
			if b < 32 && b != 9 && b != 10 && b != 13 { // Allow tab, newline, carriage return
				nonPrintableCount++
			}
		}
		if nonPrintableCount > n/5 { // >20% non-printable
			return true, nil
		}
	}

	return false, nil
}
