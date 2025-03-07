package diff

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type DirectoryDiff struct {
	File1      string
	File2      string
	Type       string
	Diffs      []Diff
	BinaryDiff bool
}

func DirectoryDiffs(dir1, dir2 string) ([]DirectoryDiff, error) {
	var diffs []DirectoryDiff

	fileMap1 := make(map[string]bool)

	err := filepath.Walk(dir1, func(path1 string, info1 os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		realPath, err := filepath.Rel(dir1, path1)
		if err != nil {
			return err
		}
		if realPath == "." {
			return nil
		}

		fileMap1[realPath] = true

		path2 := filepath.Join(dir2, realPath)
		info2, err := os.Stat(path2)

		if err != nil {
			if info1.IsDir() {
				diffs = append(diffs, DirectoryDiff{File1: realPath, File2: "", Type: "remove_dir"})
			} else {
				diffs = append(diffs, DirectoryDiff{File1: realPath, File2: "", Type: "remove"})
			}
			return nil
		}
		if info1.IsDir() && info2.IsDir() {
			diffs = append(diffs, DirectoryDiff{File1: realPath, File2: realPath, Type: "same_dir"})
			return nil
		}
		if !info1.IsDir() && !info2.IsDir() {
			file1, err := os.Open(path1)
			if err != nil {
				return fmt.Errorf("opening file1: %w", err)
			}
			defer file1.Close()

			file2, err := os.Open(path2)
			if err != nil {
				return fmt.Errorf("opening file2: %w", err)
			}
			defer file2.Close()

			fileDiffs, binDiff, err := FilesDiff(file1, file2)
			if err != nil {
				return fmt.Errorf("diffing files: %w", err)
			}
			if binDiff {
				diffs = append(diffs, DirectoryDiff{File1: realPath, File2: realPath, Type: "change", BinaryDiff: true})

			} else if len(fileDiffs) > 0 {
				diffs = append(diffs, DirectoryDiff{File1: realPath, File2: realPath, Type: "change", Diffs: fileDiffs})
			} else {
				diffs = append(diffs, DirectoryDiff{File1: realPath, File2: realPath, Type: "same"})
			}
		} else {
			diffs = append(diffs, DirectoryDiff{File1: realPath, File2: "", Type: "remove"})
			diffs = append(diffs, DirectoryDiff{File1: "", File2: realPath, Type: "add"})

		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking dir1: %w", err)
	}

	err = filepath.Walk(dir2, func(path2 string, info2 os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		realPath, err := filepath.Rel(dir2, path2)
		if err != nil {
			return err
		}
		if realPath == "." {
			return nil
		}
		if _, ok := fileMap1[realPath]; !ok {
			if info2.IsDir() {
				diffs = append(diffs, DirectoryDiff{File1: "", File2: realPath, Type: "add_dir"})
			} else {
				diffs = append(diffs, DirectoryDiff{File1: "", File2: realPath, Type: "add"})
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking dir2: %w", err)
	}

	sort.Slice(diffs, func(i, j int) bool {
		if (diffs[i].Type == "add_dir" || diffs[i].Type == "remove_dir" || diffs[i].Type == "same_dir") &&
			!(diffs[j].Type == "add_dir" || diffs[j].Type == "remove_dir" || diffs[j].Type == "same_dir") {
			return true
		}
		if !(diffs[i].Type == "add_dir" || diffs[i].Type == "remove_dir" || diffs[i].Type == "same_dir") &&
			(diffs[j].Type == "add_dir" || diffs[j].Type == "remove_dir" || diffs[j].Type == "same_dir") {
			return false
		}

		return diffs[i].File1 < diffs[j].File1 || diffs[i].File2 < diffs[j].File2
	})

	return diffs, nil
}
