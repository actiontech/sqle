package auditreport

import "embed"

//go:embed templates/audit_report.html
var auditReportTemplateFS embed.FS

const auditReportHTMLTemplatePath = "templates/audit_report.html"

// GetAuditReportHTMLTemplate reads the embedded HTML template and returns its content as a string.
func GetAuditReportHTMLTemplate() (string, error) {
	content, err := auditReportTemplateFS.ReadFile(auditReportHTMLTemplatePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
