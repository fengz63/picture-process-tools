# 图片批量处理工具 (Picture Resize Tools)

一个使用 Golang 编写的批量图片处理工具，支持 JPG/PNG 格式转换和智能缩放。

## 功能特性

- ✅ 支持 JPG/PNG/BMP/TIFF 格式批量处理
- ✅ 可导出为 JPG 或 PNG 格式
- ✅ 智能缩放保持宽高比
- ✅ 可配置最大分辨率
- ✅ 并发处理提高效率
- ✅ 递归处理子目录
- ✅ 可扩展的模块化设计

## 安装与使用

### 1. 安装依赖
```bash
go mod tidy
```

### 2. 构建程序
```bash
go build -o picture-resize-tools
```

### 3. 使用示例

#### 基本使用
```bash
# 处理当前目录下所有图片，输出到 ./output 目录
./picture-resize-tools process

# 指定输入和输出目录
./picture-resize-tools process -i ./photos -o ./processed

# 输出为 PNG 格式
./picture-resize-tools process -f png

# 设置最大宽度为 1920，质量为 90
./picture-resize-tools process -W 1920 -H 1920 -q 90

# 递归处理子目录
./picture-resize-tools process -r
```

#### 完整参数说明

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| --input | -i | "." | 输入目录路径 |
| --output | -o | "./output" | 输出目录路径 |
| --format | -f | "jpg" | 输出格式 (jpg, png) |
| --width | -W | 1920 | 最大宽度 |
| --height | -H | 1920 | 最大高度 |
| --quality | -q | 85 | 输出质量 (1-100) |
| --recursive | -r | false | 递归处理子目录 |