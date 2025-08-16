package processor

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/strukturag/libheif/go/heif"
)

type Config struct {
	OutputFormat string
	MaxWidth     int
	MaxHeight    int
	Quality      int
	OutputDir    string
}

func ProcessImage(inputPath string, config Config) error {
	// Load image
	img, err := loadImage(inputPath)
	if err != nil {
		return err
	}

	// Resize image
	img = resizeImage(img, config.MaxWidth, config.MaxHeight)

	// Generate output path
	outputPath := generateOutputPath(inputPath, config.OutputDir, config.OutputFormat)

	// Save image
	return saveImage(img, outputPath, config.OutputFormat, config.Quality)
}

// ProcessImageWithSameFormat processes image and keeps the same format
func ProcessImageWithSameFormat(inputPath string, config Config) error {
	// Load image
	img, err := loadImage(inputPath)
	if err != nil {
		return err
	}

	// Resize image
	img = resizeImage(img, config.MaxWidth, config.MaxHeight)

	// Generate output path with same format
	outputPath := generateOutputPathWithSameFormat(inputPath, config.OutputDir)

	// Get original format
	ext := filepath.Ext(strings.ToLower(inputPath))
	format := "jpg" // default
	if ext == ".png" {
		format = "png"
	}

	// Save image
	return saveImage(img, outputPath, format, config.Quality)
}

func loadImage(path string) (image.Image, error) {
	ext := filepath.Ext(strings.ToLower(path))
	if ext == ".heic" || ext == ".heif" {
		// Handle HEIC/HEIF format
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Create a new context
		ctx, err := heif.NewContext()
		if err != nil {
			return nil, err
		}

		// Read the file into the context
		err = ctx.ReadFromFile(path)
		if err != nil {
			return nil, err
		}

		// Get the primary image handle
		hdl, err := ctx.GetPrimaryImageHandle()
		if err != nil {
			return nil, err
		}

		// Decode the image
		img, err := hdl.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
		if err != nil {
			return nil, err
		}

		// Convert to go image
		return img.GetImage()
	}

	// Handle other common formats
	return imaging.Open(path)
}

func resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// If the image is already smaller than the maximum size, no adjustment is made
	if width <= maxWidth && height <= maxHeight {
		return img
	}

	// Calculate scaling ratio
	scale := 1.0
	widthScale := float64(maxWidth) / float64(width)
	heightScale := float64(maxHeight) / float64(height)

	// Use the smaller scaling ratio to maintain aspect ratio
	if widthScale < heightScale {
		scale = widthScale
	} else {
		scale = heightScale
	}

	newWidth := int(float64(width) * scale)
	newHeight := int(float64(height) * scale)

	return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
}

func generateOutputPath(inputPath, outputDir, format string) string {
	filename := filepath.Base(inputPath)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	var newExt string
	switch format {
	case "jpg":
		newExt = ".jpg"
	case "png":
		newExt = ".png"
	default:
		newExt = ".jpg"
	}

	return filepath.Join(outputDir, name+newExt)
}

// Generate output path keeping the same format
func generateOutputPathWithSameFormat(inputPath, outputDir string) string {
	filename := filepath.Base(inputPath)
	return filepath.Join(outputDir, filename)
}

func saveImage(img image.Image, path, format string, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "jpg":
		options := &jpeg.Options{Quality: quality}
		return jpeg.Encode(file, img, options)
	case "png":
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		return encoder.Encode(file, img)
	default:
		options := &jpeg.Options{Quality: quality}
		return jpeg.Encode(file, img, options)
	}
}
