package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/san-kum/diff-dance/pkg/diff"
	"github.com/san-kum/diff-dance/pkg/display"
	"github.com/san-kum/diff-dance/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diff-dance",
	Short: "Visualize file and directory differences creatively",
	Long: `diff-dance provides various ways to visualize differences
between files and directories, going beyond traditional line-by-line diffs.`,
	Run: diffDance,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("file1", "1", "", "Path to the first file or directory")
	rootCmd.Flags().StringP("file2", "2", "", "Path to the second file or directory")
	rootCmd.Flags().Bool("heatmap", false, "Generate a heatmap visualization")
	rootCmd.Flags().Bool("wordcloud", false, "Generate a word cloud visualization")
	rootCmd.Flags().Bool("structural", false, "Show structural changes (for code)")
	rootCmd.Flags().Bool("interactive", false, "Enable interactive navigation")
	rootCmd.Flags().String("format", "terminal", "Output format (terminal, html)")

	rootCmd.MarkFlagRequired("file1")
	rootCmd.MarkFlagRequired("file2")
}
func diffDance(cmd *cobra.Command, args []string) {
	file1Path, _ := cmd.Flags().GetString("file1")
	file2Path, _ := cmd.Flags().GetString("file2")
	heatmap, _ := cmd.Flags().GetBool("heatmap")
	wordcloud, _ := cmd.Flags().GetBool("wordcloud")
	structural, _ := cmd.Flags().GetBool("structural")
	interactive, _ := cmd.Flags().GetBool("interactive")
	format, _ := cmd.Flags().GetString("format")

	// Check if paths are directories.
	info1, err := os.Stat(file1Path)
	if err != nil {
		fmt.Printf("Error stating file1: %v\n", err)
		os.Exit(1)
	}
	info2, err := os.Stat(file2Path)
	if err != nil {
		fmt.Printf("Error stating file2: %v\n", err)
		os.Exit(1)
	}

	// Handle directory diffs.
	if info1.IsDir() && info2.IsDir() {
		dirDiffs, err := diff.DirectoryDiffs(file1Path, file2Path) // Call DirectoryDiffs
		if err != nil {
			fmt.Printf("Error diffing directories: %v\n", err)
			os.Exit(1)
		}

		// Handle the format for directories
		switch {
		case interactive: //If interactive
			// Interactive mode for directory diffs not supported yet
			fmt.Fprintf(os.Stderr, "Interactive mode for directories is not implemented yet")
		case format == "html": //If HTML
			err = display.HTMLDir(dirDiffs, os.Stdout)
			if err != nil {
				fmt.Printf("Error displaying HTML: %v", err)
			}
		default: //Terminal
			display.TerminalDir(dirDiffs, os.Stdout)

		}
		return // Important: Return after handling directory diff
	} else if info1.IsDir() || info2.IsDir() { // One is file, the other is directory
		fmt.Println("Cannot compare a file with a directory.")
		os.Exit(1)
	}

	// --- Handle file diffs (existing logic) ---
	file1, err := os.Open(file1Path)
	if err != nil {
		fmt.Printf("Error opening file1: %v\n")
		os.Exit(1)
	}
	defer file1.Close()

	file2, err := os.Open(file2Path)
	if err != nil {
		fmt.Printf("Error opening file2: %v\n", err)
		os.Exit(1)
	}
	defer file2.Close()

	// We seek, so we can use the files multiple times
	file1.Seek(0, 0)
	file2.Seek(0, 0)
	diffs, err := diff.Files(file1, file2) // Calculate diffs
	if err != nil {
		fmt.Printf("Error diffing files: %v\n", err)
		os.Exit(1)
	}

	switch {
	case heatmap:
		// Read lines from files *before* calculating diffs (for heatmap).
		file1Lines, err := utils.ReadLines(file1)
		if err != nil {
			fmt.Printf("Error reading lines from file1: %v\n", err)
			os.Exit(1)
		}
		file2Lines, err := utils.ReadLines(file2)
		if err != nil {
			fmt.Printf("Error reading lines from file2: %v\n", err)
			os.Exit(1)
		}
		//We reset the file pointers to calculate diffs again.
		file1.Seek(0, 0)
		file2.Seek(0, 0)
		display.Heatmap(diffs, file1Lines, file2Lines, os.Stdout)
	case wordcloud:
		display.WordCloud(diffs, os.Stdout)
	case structural:
		if filepath.Ext(file1Path) == ".go" && filepath.Ext(file2Path) == ".go" {
			structuralDiffs, err := diff.StructuralDiffs(file1, file2)
			if err != nil {
				fmt.Printf("Error calculating structural diff: %v\n", err)
				os.Exit(1)
			}
			display.Structural(structuralDiffs, os.Stdout)

		} else {
			fmt.Println("Structural diff is only supported for Go files (.go).")
		}
	case interactive:
		display.Interactive(file1Path, file2Path) // Pass file *paths*

	default: //Handle format here, so we print in terminal or HTML
		if format == "html" {
			// Check if a search term was provided via environment variable (for interactive mode).
			searchTerm := os.Getenv("DIFF_DANCE_SEARCH")
			var searchRegex *regexp.Regexp
			if searchTerm != "" {
				searchRegex, err = regexp.Compile(`\b` + regexp.QuoteMeta(searchTerm) + `\b`)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Invalid search term: %v\n", err) //Error
					// Don't exit, just proceed without highlighting.
				}
			}
			if searchRegex != nil { // If we have search
				err = display.HTMLWithHighlight(diffs, os.Stdout, searchRegex)
			} else {
				err = display.HTML(diffs, os.Stdout) // Use standard output
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating HTML: %v\n", err) //Error
				os.Exit(1)
			}
		} else {
			display.Terminal(diffs, os.Stdout) //Use standard output
		}
	}
}

