package processor

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"
)

func TestGetImageFormat(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"test.jpg", "jpg"},
		{"test.jpeg", "jpg"},
		{"test.JPG", "jpg"},
		{"test.png", "png"},
		{"test.PNG", "png"},
		{"test.bmp", "bmp"},
		{"test.tiff", "tiff"},
		{"test.tif", "tiff"},
		{"test.heic", "jpg"},
		{"test.heif", "jpg"},
		{"test.unknown", "jpg"},
		{"test", "jpg"},
	}

	for _, test := range tests {
		result := getImageFormat(test.path)
		if result != test.expected {
			t.Errorf("getImageFormat(%s) = %s, expected %s", test.path, result, test.expected)
		}
	}
}

func TestResizeImage(t *testing.T) {
	// Create a test image (100x50)
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))

	tests := []struct {
		name       string
		maxWidth   int
		maxHeight  int
		expectSize bool
		expWidth   int
		expHeight  int
	}{
		{
			name:       "No resize needed",
			maxWidth:   200,
			maxHeight:  100,
			expectSize: true,
			expWidth:   100,
			expHeight:  50,
		},
		{
			name:       "Width limited",
			maxWidth:   50,
			maxHeight:  100,
			expectSize: true,
			expWidth:   50,
			expHeight:  25,
		},
		{
			name:       "Height limited",
			maxWidth:   200,
			maxHeight:  25,
			expectSize: true,
			expWidth:   50,
			expHeight:  25,
		},
		{
			name:       "Both limited",
			maxWidth:   40,
			maxHeight:  40,
			expectSize: true,
			expWidth:   40,
			expHeight:  20,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := resizeImage(img, test.maxWidth, test.maxHeight)

			bounds := result.Bounds()
			width := bounds.Max.X - bounds.Min.X
			height := bounds.Max.Y - bounds.Min.Y

			if test.expectSize {
				if width != test.expWidth || height != test.expHeight {
					t.Errorf("resizeImage() size = %dx%d, expected %dx%d", width, height, test.expWidth, test.expHeight)
				}
			}
		})
	}
}

func TestGenerateOutputPath(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		inputPath string
		outputDir string
		format    string
		expected  string
	}{
		{
			name:      "JPG format",
			inputPath: "/path/to/image.jpg",
			outputDir: tempDir,
			format:    "jpg",
			expected:  filepath.Join(tempDir, "image.jpg"),
		},
		{
			name:      "PNG format conversion",
			inputPath: "/path/to/image.jpg",
			outputDir: tempDir,
			format:    "png",
			expected:  filepath.Join(tempDir, "image.png"),
		},
		{
			name:      "Complex path",
			inputPath: "/very/complex/path/to/my/image.jpeg",
			outputDir: tempDir,
			format:    "jpg",
			expected:  filepath.Join(tempDir, "image.jpg"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := generateOutputPath(test.inputPath, test.outputDir, test.format)
			if result != test.expected {
				t.Errorf("generateOutputPath() = %s, expected %s", result, test.expected)
			}
		})
	}
}

func TestGenerateOutputPathWithSameFormat(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		inputPath string
		outputDir string
		expected  string
	}{
		{
			name:      "Simple path",
			inputPath: "image.jpg",
			outputDir: tempDir,
			expected:  filepath.Join(tempDir, "image.jpg"),
		},
		{
			name:      "Complex path",
			inputPath: "/path/to/image.png",
			outputDir: tempDir,
			expected:  filepath.Join(tempDir, "image.png"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := generateOutputPathWithSameFormat(test.inputPath, test.outputDir)
			if result != test.expected {
				t.Errorf("generateOutputPathWithSameFormat() = %s, expected %s", result, test.expected)
			}
		})
	}
}

func TestSaveImage(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	// Fill with some color
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	tempDir := t.TempDir()

	tests := []struct {
		name   string
		format string
		path   string
	}{
		{
			name:   "Save JPG",
			format: "jpg",
			path:   filepath.Join(tempDir, "test.jpg"),
		},
		{
			name:   "Save PNG",
			format: "png",
			path:   filepath.Join(tempDir, "test.png"),
		},
		{
			name:   "Default format",
			format: "unknown",
			path:   filepath.Join(tempDir, "test.jpg"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := saveImage(img, test.path, test.format, 90)
			if err != nil {
				t.Errorf("saveImage() error = %v", err)
				return
			}

			// Check if file exists
			if _, err := os.Stat(test.path); os.IsNotExist(err) {
				t.Errorf("saveImage() file was not created: %s", test.path)
			}
		})
	}
}

func TestProcessImage(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{0, 255, 0, 255})
		}
	}

	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.png")
	outputPath := filepath.Join(tempDir, "output")

	// Save test image
	err := imaging.Save(img, inputPath)
	if err != nil {
		t.Fatalf("Failed to save test image: %v", err)
	}

	config := Config{
		OutputFormat: "jpg",
		MaxWidth:     50,
		MaxHeight:    50,
		Quality:      90,
		OutputDir:    outputPath,
	}

	// Create output directory
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	err = ProcessImage(inputPath, config)
	if err != nil {
		t.Errorf("ProcessImage() error = %v", err)
		return
	}

	// Check if output file exists
	expectedOutput := filepath.Join(outputPath, "input.jpg")
	if _, err := os.Stat(expectedOutput); os.IsNotExist(err) {
		t.Errorf("ProcessImage() output file was not created: %s", expectedOutput)
	}
}

func TestProcessImageWithSameFormat(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "input.png")
	outputPath := filepath.Join(tempDir, "output")

	// Save test image
	err := imaging.Save(img, inputPath)
	if err != nil {
		t.Fatalf("Failed to save test image: %v", err)
	}

	config := Config{
		OutputFormat: "jpg", // This should be ignored
		MaxWidth:     50,
		MaxHeight:    50,
		Quality:      90,
		OutputDir:    outputPath,
	}

	// Create output directory
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	err = ProcessImageWithSameFormat(inputPath, config)
	if err != nil {
		t.Errorf("ProcessImageWithSameFormat() error = %v", err)
		return
	}

	// Check if output file exists with same format
	expectedOutput := filepath.Join(outputPath, "input.png")
	if _, err := os.Stat(expectedOutput); os.IsNotExist(err) {
		t.Errorf("ProcessImageWithSameFormat() output file was not created: %s", expectedOutput)
	}
}

func TestLoadImageInvalidPath(t *testing.T) {
	_, err := loadImage("/nonexistent/path/image.jpg")
	if err == nil {
		t.Error("loadImage() expected error for invalid path, got nil")
	}
}
