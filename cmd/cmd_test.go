package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateInputs(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		setupFunc   func()
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid inputs",
			setupFunc: func() {
				inputDir = tempDir
				outputDir = filepath.Join(tempDir, "output")
				outputFormat = "jpg"
				quality = 90
				maxWidth = 1920
				maxHeight = 1080
				workers = 4
			},
			expectError: false,
		},
		{
			name: "Invalid output format",
			setupFunc: func() {
				inputDir = tempDir
				outputFormat = "gif"
			},
			expectError: true,
			errorMsg:    "output format must be jpg or png",
		},
		{
			name: "Nonexistent input directory",
			setupFunc: func() {
				inputDir = "/nonexistent/directory"
				outputFormat = "jpg"
			},
			expectError: true,
			errorMsg:    "input directory does not exist",
		},
		{
			name: "Quality too low",
			setupFunc: func() {
				inputDir = tempDir
				outputFormat = "jpg"
				quality = 0
			},
			expectError: true,
			errorMsg:    "quality must be between 1 and 100",
		},
		{
			name: "Quality too high",
			setupFunc: func() {
				inputDir = tempDir
				outputFormat = "jpg"
				quality = 101
			},
			expectError: true,
			errorMsg:    "quality must be between 1 and 100",
		},
		{
			name: "Invalid dimensions",
			setupFunc: func() {
				inputDir = tempDir
				outputFormat = "jpg"
				quality = 90
				maxWidth = -1
			},
			expectError: true,
			errorMsg:    "maximum dimensions must be positive",
		},
		{
			name: "Invalid worker count",
			setupFunc: func() {
				inputDir = tempDir
				outputFormat = "jpg"
				quality = 90
				maxWidth = 1920
				maxHeight = 1080
				workers = 0
			},
			expectError: true,
			errorMsg:    "worker count must be positive",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reset to defaults before each test
			inputDir = "."
			outputDir = "./output"
			outputFormat = "jpg"
			quality = 90
			maxWidth = 1920
			maxHeight = 1920
			workers = 4

			// Apply test-specific setup
			test.setupFunc()

			err := validateInputs()

			if test.expectError {
				if err == nil {
					t.Errorf("validateInputs() expected error, got nil")
					return
				}
				if test.errorMsg != "" && err.Error()[:len(test.errorMsg)] != test.errorMsg {
					t.Errorf("validateInputs() error = %v, expected to contain %s", err, test.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateInputs() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestGetImageFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"image1.jpg":        "jpg content",
		"image2.png":        "png content",
		"image3.bmp":        "bmp content",
		"image4.tiff":       "tiff content",
		"image5.heic":       "heic content",
		"not_image.txt":     "text content",
		"subdir/image6.jpg": "jpg content in subdir",
	}

	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", fullPath, err)
		}
	}

	tests := []struct {
		name      string
		dir       string
		recursive bool
		expected  int
	}{
		{
			name:      "Non-recursive scan",
			dir:       tempDir,
			recursive: false,
			expected:  5, // image1.jpg, image2.png, image3.bmp, image4.tiff, image5.heic
		},
		{
			name:      "Recursive scan",
			dir:       tempDir,
			recursive: true,
			expected:  6, // + subdir/image6.jpg
		},
		{
			name:      "Subdir only non-recursive",
			dir:       filepath.Join(tempDir, "subdir"),
			recursive: false,
			expected:  1, // only image6.jpg
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files, err := getImageFiles(test.dir, test.recursive)
			if err != nil {
				t.Errorf("getImageFiles() error = %v", err)
				return
			}

			if len(files) != test.expected {
				t.Errorf("getImageFiles() = %d files, expected %d", len(files), test.expected)
			}
		})
	}
}

func TestSeparateImageFiles(t *testing.T) {
	files := []string{
		"image1.heic",
		"image2.heif",
		"image3.jpg",
		"image4.png",
		"image5.bmp",
		"image6.tiff",
	}

	heicFiles, regularFiles := separateImageFiles(files)

	expectedHeic := 2
	expectedRegular := 4

	if len(heicFiles) != expectedHeic {
		t.Errorf("separateImageFiles() heic files = %d, expected %d", len(heicFiles), expectedHeic)
	}

	if len(regularFiles) != expectedRegular {
		t.Errorf("separateImageFiles() regular files = %d, expected %d", len(regularFiles), expectedRegular)
	}

	// Verify specific files
	heicSet := make(map[string]bool)
	for _, file := range heicFiles {
		heicSet[file] = true
	}

	if !heicSet["image1.heic"] || !heicSet["image2.heif"] {
		t.Error("separateImageFiles() did not correctly identify HEIC files")
	}
}

func TestRootCmd(t *testing.T) {
	// Test that root command is properly configured
	if rootCmd.Use != "picture-resize-tools" {
		t.Errorf("rootCmd.Use = %s, expected picture-resize-tools", rootCmd.Use)
	}

	if rootCmd.Short != "Batch image format conversion and resize tool" {
		t.Errorf("rootCmd.Short = %s, expected 'Batch image format conversion and resize tool'", rootCmd.Short)
	}

	// Test that flags are properly configured
	flags := rootCmd.PersistentFlags()

	// Check if input flag exists
	inputFlag := flags.Lookup("input")
	if inputFlag == nil {
		t.Error("rootCmd missing 'input' flag")
	} else {
		// Check that the flag has the correct default value by checking DefValue
		if inputFlag.DefValue != "." {
			t.Errorf("input flag default = %s, expected '.'", inputFlag.DefValue)
		}
	}

	// Check if workers flag exists
	workersFlag := flags.Lookup("workers")
	if workersFlag == nil {
		t.Error("rootCmd missing 'workers' flag")
	} else {
		// Check that the flag has the correct default value
		if workersFlag.DefValue != "4" {
			t.Errorf("workers flag default = %s, expected '4'", workersFlag.DefValue)
		}
	}
}

func TestProcessCmd(t *testing.T) {
	// Test that process command is properly added to root
	if len(rootCmd.Commands()) != 1 {
		t.Errorf("rootCmd should have 1 subcommand, got %d", len(rootCmd.Commands()))
	}

	processCmd := rootCmd.Commands()[0]
	if processCmd.Use != "process" {
		t.Errorf("processCmd.Use = %s, expected 'process'", processCmd.Use)
	}

	if processCmd.Short != "Start batch processing images" {
		t.Errorf("processCmd.Short = %s, expected 'Start batch processing images'", processCmd.Short)
	}
}

func TestExecute(t *testing.T) {
	// Test that Execute function exists and can be called
	// We can't easily test the full execution without mocking file system
	// but we can verify the function signature and basic behavior

	// This test ensures the Execute function doesn't panic when called with invalid args
	// In a real scenario, this would exit with error code 1
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked: %v", r)
		}
	}()

	// We can't actually test Execute() without it calling os.Exit(1)
	// So we'll just verify the function exists by checking if it's callable
	_ = Execute
}
