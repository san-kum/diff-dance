package cmd

import (
	"fmt"
	"os"

	"github.com/san-kum/diff-dance/pkg/diff"
	"github.com/san-kum/diff-dance/pkg/display"
	"github.com/san-kum/diff-dance/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diff-dance",
	Short: "Visualize file and directory differences creatively",
	Long:  `diff-dance provides various ways to visualize differences between files and directories, going beyond traditional line-by-line diffs.`,
	Run:   diffDance,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("file1", "1", "", "Path to first file")
	rootCmd.Flags().StringP("file2", "2", "", "Path to second file")
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

	file1, err := os.Open(file1Path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file1.Close()
	file2, err := os.Open(file2Path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file2.Close()

	file1Lines, err := utils.ReadLines(file1)
	if err != nil {
		fmt.Printf("Error reading lines from file 1: %v", err)
		os.Exit(1)
	}
	file2Lines, err := utils.ReadLines(file2)
	if err != nil {
		fmt.Printf("Error reading lines from file 2: %v", err)
		os.Exit(1)
	}
	file1.Seek(0, 0)
	file2.Seek(0, 0)

	diffs, err := diff.Files(file1, file2)
	if err != nil {
		fmt.Printf("Error diffing files: %v\n", err)
		os.Exit(1)
	}
	switch {
	case heatmap:
		display.HeatMap(diffs, file1Lines, file2Lines)
	case wordcloud:
		display.WordCloud(diffs)
	case structural:
		fmt.Println("TBD")
	case interactive:
		fmt.Println("TBD")
	default:
		display.Terminal(diffs)
	}
	if format != "terminal" {
		fmt.Println("HTML/Image output (not implemented)")
	}
}
