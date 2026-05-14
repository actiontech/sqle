package auditreport

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/utils"
)

// CSVReportGenerator CSV 格式报告生成器
// 复用已有的 CSVBuilder 生成 CSV 报告，实现 ReportGenerator 接口。
type CSVReportGenerator struct{}

// NewCSVReportGenerator 创建并返回一个新的 CSVReportGenerator 实例
func NewCSVReportGenerator() *CSVReportGenerator {
	return &CSVReportGenerator{}
}

// ReportType 返回生成器支持的导出格式
func (g *CSVReportGenerator) ReportType() ExportFormat {
	return CsvExportFormat
}

// CSVHeaders 返回 CSV 报告的表头列表
func (g *CSVReportGenerator) CSVHeaders(data *AuditReportData) []string {
	return []string{
		data.Labels.Number,
		data.Labels.SQL,
		data.Labels.AuditStatus,
		data.Labels.AuditResult,
		data.Labels.ExecStatus,
		data.Labels.ExecResult,
		data.Labels.RollbackSQL,
		data.Labels.Description,
	}
}

// ToCSVRow 将单条审核 SQL 转为 CSV 行
func (g *CSVReportGenerator) ToCSVRow(item *AuditSQLItem) []string {
	return []string{
		fmt.Sprintf("%d", item.Number),
		item.SQL,
		item.AuditStatus,
		item.AuditResult,
		item.ExecStatus,
		item.ExecResult,
		item.RollbackSQL,
		item.Description,
	}
}

// Generate 根据审核报告数据生成 CSV 格式的文件
func (g *CSVReportGenerator) Generate(data *AuditReportData) (*utils.ExportDataResult, error) {
	builder := utils.NewCSVBuilder()

	if err := builder.WriteHeader(g.CSVHeaders(data)); err != nil {
		return nil, fmt.Errorf("write csv header failed: %v", err)
	}

	for i := range data.SQLList {
		if err := builder.WriteRow(g.ToCSVRow(&data.SQLList[i])); err != nil {
			return nil, fmt.Errorf("write csv row failed: %v", err)
		}
	}

	content := builder.FlushAndGetBuffer().Bytes()
	if err := builder.Error(); err != nil {
		return nil, fmt.Errorf("csv builder error: %v", err)
	}

	return &utils.ExportDataResult{
		Content:     content,
		ContentType: "text/csv",
		FileName:    fmt.Sprintf("SQL_audit_report_%s_%d.csv", data.InstanceName, data.TaskID),
	}, nil
}
