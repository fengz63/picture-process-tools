package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"picture-resize-tools/pkg/processor"
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Start batch processing images",
	Run:   runProcess,
}

func init() {
	rootCmd.AddCommand(processCmd)
}

// validateInputs validates command line inputs
func validateInputs() error {
	// Validate output format
	if outputFormat != "jpg" && outputFormat != "png" {
		return fmt.Errorf("output format must be jpg or png, got: %s", outputFormat)
	}

	// Validate input directory exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		return fmt.Errorf("input directory does not exist: %s", inputDir)
	}

	// Validate quality range
	if quality < 1 || quality > 100 {
		return fmt.Errorf("quality must be between 1 and 100, got: %d", quality)
	}

	// Validate dimensions
	if maxWidth <= 0 || maxHeight <= 0 {
		return fmt.Errorf("maximum dimensions must be positive, got: %dx%d", maxWidth, maxHeight)
	}

	// Validate worker count
	if workers <= 0 {
		return fmt.Errorf("worker count must be positive, got: %d", workers)
	}

	return nil
}

func runProcess(cmd *cobra.Command, args []string) {
	// Validate inputs
	if err := validateInputs(); err != nil {
		fmt.Printf("Input validation failed: %v\n", err)
		os.Exit(1)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory '%s': %v\n", outputDir, err)
		os.Exit(1)
	}

	// Get all image files
	imageFiles, err := getImageFiles(inputDir, recursive)
	if err != nil {
		fmt.Printf("Failed to scan image files: %v\n", err)
		os.Exit(1)
	}

	if len(imageFiles) == 0 {
		fmt.Println("No image files found")
		return
	}

	// Separate HEIC and regular images
	heicFiles, regularFiles := separateImageFiles(imageFiles)

	fmt.Printf("Found %d image files (%d HEIC, %d regular), starting processing...\n", len(imageFiles), len(heicFiles), len(regularFiles))

	// Configure processor
	config := processor.Config{
		OutputFormat: outputFormat,
		MaxWidth:     maxWidth,
		MaxHeight:    maxHeight,
		Quality:      quality,
		OutputDir:    outputDir,
	}

	// If there are HEIC files, process all images with format conversion
	if len(heicFiles) > 0 {
		fmt.Println("HEIC files found, processing all images with format conversion...")
		processImagesConcurrently(imageFiles, config)
	} else {
		// No HEIC files, only resize regular images and keep original format
		fmt.Println("No HEIC files found, only resizing regular images and keeping original format...")
		processImagesWithSameFormat(regularFiles, config)
	}

	fmt.Println("All images processed!")
}

func getImageFiles(dir string, recursive bool) ([]string, error) {
	var files []string
	exts := map[string]bool{
		".heic": true, ".heif": true,
		".jpg": true, ".jpeg": true,
		".png": true, ".bmp": true,
		".tiff": true, ".tif": true,
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(strings.ToLower(path))
			if exts[ext] {
				files = append(files, path)
			}
		} else if !recursive && path != dir {
			return filepath.SkipDir
		}
		return nil
	}

	err := filepath.Walk(dir, walkFunc)
	return files, err
}

// Separate HEIC and regular image files
func separateImageFiles(files []string) ([]string, []string) {
	heicFiles := []string{}
	regularFiles := []string{}

	for _, file := range files {
		ext := filepath.Ext(strings.ToLower(file))
		if ext == ".heic" || ext == ".heif" {
			heicFiles = append(heicFiles, file)
		} else {
			regularFiles = append(regularFiles, file)
		}
	}

	return heicFiles, regularFiles
}

// Process images concurrently
func processImagesConcurrently(files []string, config processor.Config) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, workers) // Use configurable worker count

	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := processor.ProcessImage(filePath, config); err != nil {
				fmt.Printf("Processing failed %s: %v\n", filePath, err)
			} else {
				fmt.Printf("Processing completed: %s\n", filepath.Base(filePath))
			}
		}(file)
	}

	wg.Wait()
}

// Process images concurrently while keeping the same format
func processImagesWithSameFormat(files []string, config processor.Config) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, workers) // Use configurable worker count

	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := processor.ProcessImageWithSameFormat(filePath, config); err != nil {
				fmt.Printf("Processing failed %s: %v\n", filePath, err)
			} else {
				fmt.Printf("Processing completed: %s\n", filepath.Base(filePath))
			}
		}(file)
	}

	wg.Wait()
}
