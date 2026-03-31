package utils

import "fmt"

// CSVReportGenerator CSV 格式报告生成器
// 复用已有的 CSVBuilder 生成 CSV 报告，实现 ReportGenerator 接口。
type CSVReportGenerator struct{}

// NewCSVReportGenerator 创建并返回一个新的 CSVReportGenerator 实例
func NewCSVReportGenerator() *CSVReportGenerator {
	return &CSVReportGenerator{}
}

// Format 返回生成器支持的导出格式
func (g *CSVReportGenerator) Format() ExportFormat {
	return CsvExportFormat
}

// Generate 根据审核报告数据生成 CSV 格式的文件
//
// 参数：
//
//	data: 审核报告完整数据模型
//
// 返回：
//
//	*ExportDataResult: 包含 CSV 文件内容、ContentType 和文件名
//	error: 生成过程中的错误
func (g *CSVReportGenerator) Generate(data *AuditReportData) (*ExportDataResult, error) {
	builder := NewCSVBuilder()

	// 写入表头
	if err := builder.WriteHeader(data.CSVHeaders()); err != nil {
		return nil, fmt.Errorf("write csv header failed: %v", err)
	}

	// 写入数据行
	for _, sql := range data.SQLList {
		if err := builder.WriteRow(sql.ToCSVRow()); err != nil {
			return nil, fmt.Errorf("write csv row failed: %v", err)
		}
	}

	// 刷新缓冲区并获取内容
	content := builder.FlushAndGetBuffer().Bytes()
	if err := builder.Error(); err != nil {
		return nil, fmt.Errorf("csv builder error: %v", err)
	}

	return &ExportDataResult{
		Content:     content,
		ContentType: "text/csv",
		FileName:    fmt.Sprintf("SQL_audit_report_%s_%d.csv", data.InstanceName, data.TaskID),
	}, nil
}
