//go:build !enterprise

package auditreport

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/utils"
)

// ExportAuditReport CE 版统一导出入口。
// CE 版仅支持 CSV 和 HTML 两种格式。
// 请求 PDF 或 WORD 格式时返回错误提示，提醒用户需要企业版。
// 无效格式返回错误（REQ-6.3）。
func ExportAuditReport(format utils.ExportFormat, data *AuditReportData) (*utils.ExportDataResult, error) {
	switch format {
	case utils.CsvExportFormat:
		return NewCSVReportGenerator().Generate(data)
	case utils.ExportFormatHTML:
		gen, err := NewHTMLReportGenerator()
		if err != nil {
			return nil, err
		}
		return gen.Generate(data)
	case utils.ExportFormatPDF, utils.ExportFormatWORD:
		return nil, fmt.Errorf("export format %s is only supported in enterprise edition", format)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}
