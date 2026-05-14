package auditreport

import (
	"fmt"
	"strings"
)

// ExportFormat 审核报告等业务的导出格式（与 utils 中通用表格导出的 csv/excel 区分）。
type ExportFormat string

const (
	CsvExportFormat   ExportFormat = "csv"
	ExcelExportFormat ExportFormat = "excel"
	ExportFormatHTML  ExportFormat = "html"
	ExportFormatPDF   ExportFormat = "pdf"
	ExportFormatWORD  ExportFormat = "word"
)

// NormalizeExportFormatStr 规范化导出格式查询参数。
// 空字符串默认返回 CSV（向后兼容）；无效格式返回错误。
func NormalizeExportFormatStr(format string) (ExportFormat, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "html":
		return ExportFormatHTML, nil
	case "pdf":
		return ExportFormatPDF, nil
	case "word", "docx":
		return ExportFormatWORD, nil
	case "excel", "xlsx":
		return ExcelExportFormat, nil
	case "csv", "":
		return CsvExportFormat, nil
	default:
		return "", fmt.Errorf("unsupported export format: %s", format)
	}
}
