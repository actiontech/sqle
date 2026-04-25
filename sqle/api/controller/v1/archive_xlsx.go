package v1

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

// getSqlsFromXlsx 从 XLSX 文件中按模板约定提取 SQL。
// 使用 github.com/xuri/excelize/v2 库解析 XLSX 文件，读取第一个 Sheet，
// 第一行为表头，查找列名包含 "SQL"（不区分大小写）的列索引，逐行读取该列内容。
// 函数签名与设计文档 3.1.4 节对齐。
func getSqlsFromXlsx(c echo.Context) (string, bool, error) {
	file, err := c.FormFile(InputSQLFileName)
	if err == http.ErrMissingFile {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}

	f, err := file.Open()
	if err != nil {
		return "", false, err
	}
	defer f.Close()

	sqlContent, err := processXlsxContent(f)
	if err != nil {
		return "", false, err
	}

	return sqlContent, true, nil
}

// processXlsxContent 从 io.Reader 中读取 XLSX 内容并提取 SQL。
// 该函数封装了 XLSX 解析的核心逻辑，独立于 echo.Context，便于单元测试。
// 处理流程：
//  1. 使用 excelize.OpenReader() 打开文件
//  2. 获取第一个 Sheet
//  3. 第一行为表头，查找列名包含 "SQL"（不区分大小写）的列索引
//  4. 若未找到 SQL 列，返回错误
//  5. 逐行读取该列内容，跳过空行
//  6. 拼接所有 SQL 用 ";\n" 连接返回
func processXlsxContent(r io.Reader) (string, error) {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return "", fmt.Errorf("open xlsx file failed: %v", err)
	}
	defer xlsx.Close()

	// 获取第一个 Sheet 名称
	sheetList := xlsx.GetSheetList()
	if len(sheetList) == 0 {
		return "", fmt.Errorf("xlsx file has no sheets")
	}
	sheetName := sheetList[0]

	// 读取所有行
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return "", fmt.Errorf("read xlsx rows failed: %v", err)
	}

	// 空文件（无数据行）
	if len(rows) == 0 {
		return "", nil
	}

	// 第一行为表头，查找列名包含 "SQL"（不区分大小写）的列索引
	headerRow := rows[0]
	sqlColIdx := -1
	for i, cellValue := range headerRow {
		if strings.Contains(strings.ToLower(cellValue), "sql") {
			sqlColIdx = i
			break
		}
	}

	if sqlColIdx == -1 {
		return "", fmt.Errorf("no column containing \"SQL\" found in the header row")
	}

	// 逐行读取 SQL 列内容，跳过空行
	sqls := make([]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		// 行的列数可能小于 SQL 列索引（excelize 对尾部空列会截断）
		if sqlColIdx >= len(row) {
			continue
		}
		cellValue := strings.TrimSpace(row[sqlColIdx])
		if cellValue == "" {
			continue
		}
		sqls = append(sqls, cellValue)
	}

	return strings.Join(sqls, ";\n"), nil
}
