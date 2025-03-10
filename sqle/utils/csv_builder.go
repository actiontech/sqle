package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

/*
	该文件用于维护SQLE中创建CSV文件的构建器
*/

const (
	// MaxNumberPerCell Excel 单元格字符数限制
	MaxNumberPerCell int = 32767
)

var ErrCSVColumnCountNotMatch = fmt.Errorf("csv column count not match")

// CSVBuilder 用于构建CSV文件，支持Excel兼容性处理
//
// 主要功能：
//   - 自动添加UTF-8 BOM头
//   - 处理Excel单元格字符数限制
//   - 提供简单的API进行CSV文件构建
//
// 参考：https://support.microsoft.com/search/results?query=excel-specifications-and-limits
//
// 使用示例：
//
//	builder := NewCSVBuilder()
//	builder.WriteHeader([]string{"Name", "Age"}) // 若设置了表头，则每一行的长度需要与表头长度匹配
//	builder.WriteRows([][]string{{"John", "30"}, {"Alice", "25"}})
//	buffer := builder.FlushAndGetBuffer()
type CSVBuilder struct {
	columnCount uint
	buffer      *bytes.Buffer
	csvWriter   *csv.Writer
}

// NewCSVBuilder 创建并返回一个新的CSVBuilder实例
//
// 返回：
//
//	*CSVBuilder: 新的CSVBuilder实例
func NewCSVBuilder() *CSVBuilder {
	buffer := new(bytes.Buffer)
	// 写入 UTF-8 BOM 有助于Excel正确识别文件的编码格式
	buffer.WriteString("\xEF\xBB\xBF")
	csvWriter := csv.NewWriter(buffer)
	return &CSVBuilder{
		buffer:    buffer,
		csvWriter: csvWriter,
	}
}

// WriteHeader 写入CSV文件的表头
//
// 参数：
//
//	header: 表头字符串数组
//
// 返回：
//
//	error: 如果写入失败则返回错误
func (b *CSVBuilder) WriteHeader(header []string) error {
	b.columnCount = uint(len(header))
	return b.csvWriter.Write(header)
}

// WriteRows 写入多行数据到CSV文件
//
// 参数：
//
//	rows: 二维字符串数组，每行代表一条记录
//
// 返回：
//
//	error: 如果写入失败则返回错误
func (b *CSVBuilder) WriteRows(rows [][]string) error {
	for _, row := range rows {
		err := b.WriteRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteRow 写入单行数据到CSV文件
//
// 参数：
//
//	row: 字符串数组，代表一条记录
//
// 返回：
//
//	error: 如果列数不匹配或写入失败则返回错误
func (b *CSVBuilder) WriteRow(row []string) error {
	// 检查列数是否匹配
	if b.columnCount > 0 && b.columnCount != uint(len(row)) {
		return ErrCSVColumnCountNotMatch
	}
	// 超过最大字符数，则截断并标记 Excel 单元格
	for idx, cell := range row {
		if len(cell) > MaxNumberPerCell {
			row[idx] = TruncateAndMarkForExcelCell(cell)
		}
	}
	return b.csvWriter.Write(row)
}

// FlushAndGetBuffer 刷新缓冲区并返回包含CSV内容的缓冲区
//
// 返回：
//
//	*bytes.Buffer: 包含CSV文件内容的缓冲区
func (b *CSVBuilder) FlushAndGetBuffer() *bytes.Buffer {
	b.csvWriter.Flush()
	return b.buffer
}

// Error 返回CSV写入过程中遇到的任何错误
//
// 返回：
//
// 	error: 如果存在错误则返回错误，否则返回nil
func (b *CSVBuilder) Error() error {
	return b.csvWriter.Error()
}
