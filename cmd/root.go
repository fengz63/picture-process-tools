package cmd

import (
	"github.com/spf13/cobra"
)

var (
	inputDir     string
	outputDir    string
	outputFormat string
	maxWidth     int
	maxHeight    int
	quality      int
	recursive    bool
)

var rootCmd = &cobra.Command{
	Use:   "picture-resize-tools",
	Short: "Batch image format conversion and resize tool",
	Long: `Supports batch conversion of JPG/PNG/BMP/TIFF formats, export to JPG/PNG format,
intelligent resize maintains aspect ratio, maximum side resize to specified resolution`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&inputDir, "input", "i", ".", "Input directory path")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "./output", "Output directory path")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "jpg", "Output format (jpg, png)")
	rootCmd.PersistentFlags().IntVarP(&maxWidth, "width", "W", 1920, "Maximum width")
	rootCmd.PersistentFlags().IntVarP(&maxHeight, "height", "H", 1920, "Maximum height")
	rootCmd.PersistentFlags().IntVarP(&quality, "quality", "q", 85, "Output quality (1-100)")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "Recursively process subdirectories")
}
