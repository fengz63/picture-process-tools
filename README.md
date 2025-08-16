# Picture Resize Tools

A batch image processing tool written in Golang that supports JPG/PNG format conversion and intelligent resizing.

## Features

- ✅ Supports batch processing of JPG/PNG/BMP/TIFF formats
- ✅ Can export to JPG or PNG format
- ✅ Intelligent resizing maintains aspect ratio
- ✅ Configurable maximum resolution
- ✅ Concurrent processing for improved efficiency
- ✅ Recursive processing of subdirectories
- ✅ Extensible modular design

## Installation and Usage

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Build the Program
```bash
go build -o picture-resize-tools
```

### 3. Usage Examples

#### Basic Usage
```bash
# Process all images in the current directory, output to ./output directory
./picture-resize-tools process

# Specify input and output directories
./picture-resize-tools process -i ./photos -o ./processed

# Export as PNG format
./picture-resize-tools process -f png

# Set maximum width to 1920, quality to 90
./picture-resize-tools process -W 1920 -H 1920 -q 90

# Recursively process subdirectories
./picture-resize-tools process -r
```

#### Complete Parameter Description

| Parameter | Short | Default | Description |
|-----------|-------|---------|-------------|
| input     | -i    | .       | Input directory |
| output    | -o    | ./output| Output directory |
| format    | -f    | jpg     | Output format (jpg/png) |
| maxWidth  | -W    | 1920    | Maximum width |
| maxHeight | -H    | 1920    | Maximum height |
| quality   | -q    | 90      | JPEG quality (1-100) |
| recursive | -r    | false   | Recursively process subdirectories |
| workers   | -w    | 4       | Number of concurrent workers |

## Language

[中文版 README](README_zh.md)

## License

MIT