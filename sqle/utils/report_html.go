package utils

import (
	"bytes"
	"fmt"
	"html/template"
)

// HTMLReportGenerator HTML 格式报告生成器
// 使用 html/template 渲染嵌入的 HTML 模板，自动进行 HTML 转义防止 XSS。
// 实现 ReportGenerator 接口。
type HTMLReportGenerator struct {
	tmpl *template.Template
}

// NewHTMLReportGenerator 创建并返回一个新的 HTMLReportGenerator 实例。
// 在创建时解析嵌入的 HTML 模板，如果模板解析失败则返回错误。
func NewHTMLReportGenerator() (*HTMLReportGenerator, error) {
	templateContent, err := GetAuditReportHTMLTemplate()
	if err != nil {
		return nil, fmt.Errorf("read HTML template failed: %w", err)
	}

	tmpl, err := template.New("audit_report").Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("parse HTML template failed: %w", err)
	}

	return &HTMLReportGenerator{tmpl: tmpl}, nil
}

// Format 返回生成器支持的导出格式
func (g *HTMLReportGenerator) Format() ExportFormat {
	return ExportFormatHTML
}

// Generate 根据审核报告数据生成 HTML 格式的文件
//
// 参数：
//
//	data: 审核报告完整数据模型
//
// 返回：
//
//	*ExportDataResult: 包含 HTML 文件内容、ContentType 和文件名
//	error: 生成过程中的错误
func (g *HTMLReportGenerator) Generate(data *AuditReportData) (*ExportDataResult, error) {
	var buf bytes.Buffer
	if err := g.tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("render HTML report failed: %w", err)
	}

	return &ExportDataResult{
		Content:     buf.Bytes(),
		ContentType: "text/html",
		FileName:    fmt.Sprintf("SQL_audit_report_%s_%d.html", data.InstanceName, data.TaskID),
	}, nil
}
