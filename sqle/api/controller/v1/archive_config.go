package v1

import "fmt"

// archiveConfig 压缩包处理的统一安全配置
type archiveConfig struct {
	MaxTotalSize    int64 // 解压后文件总大小上限（字节）
	MaxFileCount    int   // 压缩包内文件数量上限
	MaxNestingDepth int   // 嵌套压缩包最大解压深度
}

// defaultArchiveConfig 默认的压缩包安全配置
var defaultArchiveConfig = archiveConfig{
	MaxTotalSize:    10 * 1024 * 1024, // 10MB
	MaxFileCount:    1000,
	MaxNestingDepth: 1,
}

// checkSize 检查解压后文件总大小是否超出限制
func (c *archiveConfig) checkSize(currentTotal int64, entrySize int64) error {
	if currentTotal+entrySize > c.MaxTotalSize {
		return fmt.Errorf("archive total size exceeds limit %d bytes", c.MaxTotalSize)
	}
	return nil
}

// checkFileCount 检查压缩包内文件数量是否超出限制
func (c *archiveConfig) checkFileCount(count int) error {
	if count > c.MaxFileCount {
		return fmt.Errorf("archive file count exceeds limit %d", c.MaxFileCount)
	}
	return nil
}

// supportedArchiveExts 支持的压缩包文件扩展名
var supportedArchiveExts = map[string]bool{
	".zip": true,
	".rar": true,
	".7z":  true,
}

// supportedTextExts 支持的文本文件扩展名（用于压缩包内部文件处理和单文件上传）
// 注意：.xlsx 为第二期功能，暂不加入
var supportedTextExts = map[string]bool{
	".sql":  true,
	".txt":  true,
	".java": true,
}
