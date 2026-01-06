package utils

import (
	"fmt"
	"unicode/utf8"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// ExcelBuilder 用于构建Excel文件，支持自适应列宽
type ExcelBuilder struct {
	file      *excelize.File
	sheetName string
	rowIndex  int
	colWidths map[int]float64 // 列索引 -> 最大宽度
}

// NewExcelBuilder 创建并返回一个新的ExcelBuilder实例
//
// 返回：
//
//	*ExcelBuilder: 新的ExcelBuilder实例
//	error: 如果创建失败则返回错误
func NewExcelBuilder() (*ExcelBuilder, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// 删除默认的Sheet1，创建新的工作表
	f.NewSheet(sheetName)

	// 删除默认的Sheet1
	f.DeleteSheet("Sheet1")

	return &ExcelBuilder{
		file:      f,
		sheetName: sheetName,
		rowIndex:  1,
		colWidths: make(map[int]float64),
	}, nil
}

// WriteHeader 写入Excel文件的表头
//
// 参数：
//
//	header: 表头字符串数组
//
// 返回：
//
//	error: 如果写入失败则返回错误
func (b *ExcelBuilder) WriteHeader(header []string) error {
	// excelize v1使用JSON字符串定义样式
	styleJSON := `{"font":{"bold":true},"fill":{"type":"pattern","color":["#E0E0E0"],"pattern":1}}`
	styleID, err := b.file.NewStyle(styleJSON)
	if err != nil {
		return fmt.Errorf("create header style failed: %v", err)
	}

	for col, value := range header {
		cellName := getCellName(col+1, b.rowIndex)

		b.file.SetCellValue(b.sheetName, cellName, value)

		// 应用表头样式
		b.file.SetCellStyle(b.sheetName, cellName, cellName, styleID)

		// 更新列宽
		b.updateColumnWidth(col, value)
	}

	b.rowIndex++
	return nil
}

// WriteRow 写入单行数据到Excel文件
//
// 参数：
//
//	row: 字符串数组，代表一条记录
//
// 返回：
//
//	error: 如果写入失败则返回错误
func (b *ExcelBuilder) WriteRow(row []string) error {
	for col, value := range row {
		cellName := getCellName(col+1, b.rowIndex)

		// 处理超长字符串
		truncatedValue := TruncateAndMarkForExcelCell(value)

		b.file.SetCellValue(b.sheetName, cellName, truncatedValue)

		// 更新列宽
		b.updateColumnWidth(col, truncatedValue)
	}

	b.rowIndex++
	return nil
}

// WriteRows 写入多行数据到Excel文件
//
// 参数：
//
//	rows: 二维字符串数组，每行代表一条记录
//
// 返回：
//
//	error: 如果写入失败则返回错误
func (b *ExcelBuilder) WriteRows(rows [][]string) error {
	for _, row := range rows {
		if err := b.WriteRow(row); err != nil {
			return err
		}
	}
	return nil
}

// updateColumnWidth 更新列的宽度，基于单元格内容
func (b *ExcelBuilder) updateColumnWidth(col int, value string) {
	// 计算字符串宽度（考虑中文字符）
	width := calculateStringWidth(value)

	// 添加一些边距（约2个字符宽度）
	width += 2

	// 设置最小宽度为10，最大宽度为100
	if width < 10 {
		width = 10
	} else if width > 100 {
		width = 100
	}

	// 更新最大宽度
	if currentWidth, exists := b.colWidths[col]; !exists || width > currentWidth {
		b.colWidths[col] = width
	}
}

// calculateStringWidth 计算字符串的显示宽度
// 中文字符按2个字符宽度计算，其他字符按1个字符宽度计算
func calculateStringWidth(s string) float64 {
	var width float64
	for _, r := range s {
		if utf8.RuneLen(r) > 1 {
			// 中文字符或其他多字节字符，按2个字符宽度计算
			width += 2
		} else {
			// ASCII字符，按1个字符宽度计算
			width += 1
		}
	}
	return width
}

// SetColumnWidths 设置所有列的宽度为自适应后的宽度
func (b *ExcelBuilder) SetColumnWidths() error {
	for col, width := range b.colWidths {
		colName := getColumnName(col + 1)

		b.file.SetColWidth(b.sheetName, colName, colName, width)
	}
	return nil
}

// getCellName 将列号和行号转换为Excel单元格名称（如 A1, B2）
func getCellName(col, row int) string {
	colName := getColumnName(col)
	return fmt.Sprintf("%s%d", colName, row)
}

// getColumnName 将列号转换为Excel列名（如 1->A, 2->B, 27->AA）
func getColumnName(col int) string {
	result := ""
	for col > 0 {
		col--
		result = string(rune('A'+col%26)) + result
		col /= 26
	}
	return result
}

// FlushAndGetBuffer 刷新缓冲区并返回包含Excel内容的字节数组
//
// 返回：
//
//	[]byte: 包含Excel文件内容的字节数组
//	error: 如果处理失败则返回错误
func (b *ExcelBuilder) FlushAndGetBuffer() ([]byte, error) {
	// 设置所有列的宽度
	if err := b.SetColumnWidths(); err != nil {
		return nil, fmt.Errorf("set column widths failed: %v", err)
	}

	// 设置活动工作表
	sheetIndex := b.file.GetSheetIndex(b.sheetName)
	if sheetIndex > 0 {
		b.file.SetActiveSheet(sheetIndex)
	}

	// 写入到缓冲区
	buffer, err := b.file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("write to buffer failed: %v", err)
	}

	return buffer.Bytes(), nil
}

// Close 关闭Excel文件并释放资源（excelize v1不需要显式关闭）
func (b *ExcelBuilder) Close() error {
	// excelize v1的File结构体没有Close方法，这里不需要做任何操作
	return nil
}
