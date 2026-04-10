package auditreport

import (
	"time"

	"github.com/actiontech/sqle/sqle/utils"
)

// AuditReportData 审核报告完整数据模型
type AuditReportData struct {
	// 元信息
	TaskID       uint64    `json:"task_id"`
	Title        string    `json:"title"`         // 报告标题 (i18n)
	InstanceName string    `json:"instance_name"` // 数据源名称
	Schema       string    `json:"schema"`
	GeneratedAt  time.Time `json:"generated_at"` // 报告生成时间
	Lang         string    `json:"lang"`         // 语言: zh-CN / en-US
	LogoBase64   string    `json:"logo_base64"`  // Logo 图片 base64

	// 审核概要
	Summary AuditSummary `json:"summary"`

	// 审核结果统计
	Statistics AuditStatistics `json:"statistics"`

	// SQL 列表
	SQLList     []AuditSQLItem `json:"sql_list"`     // 全部 SQL
	ProblemSQLs []AuditSQLItem `json:"problem_sqls"` // 问题 SQL（AuditLevel != normal）

	// 国际化标签
	Labels ReportLabels `json:"labels"`
}

// AuditSummary 审核概要
type AuditSummary struct {
	AuditTime    string  `json:"audit_time"`
	InstanceName string  `json:"instance_name"`
	Schema       string  `json:"schema"`
	TotalSQL     int     `json:"total_sql"`
	PassRate     float64 `json:"pass_rate"`
	Score        int32   `json:"score"`
	AuditLevel   string  `json:"audit_level"`
}

// AuditStatistics 审核结果统计
type AuditStatistics struct {
	LevelDistribution []LevelCount `json:"level_distribution"` // 按等级分布
	RuleHits          []RuleHit    `json:"rule_hits"`            // 规则命中统计
}

// LevelCount 等级统计
type LevelCount struct {
	Level string `json:"level"` // normal/notice/warn/error
	Count int    `json:"count"`
}

// RuleHit 规则命中统计
type RuleHit struct {
	RuleName string `json:"rule_name"`
	HitCount int    `json:"hit_count"`
}

// AuditSQLItem 单条 SQL 审核结果
type AuditSQLItem struct {
	Number      uint   `json:"number"`
	SQL         string `json:"sql"`
	AuditLevel  string `json:"audit_level"`
	AuditStatus string `json:"audit_status"`
	AuditResult string `json:"audit_result"` // 审核结果描述
	ExecStatus  string `json:"exec_status"`
	ExecResult  string `json:"exec_result"`
	RollbackSQL string `json:"rollback_sql"`
	Description string `json:"description"`
	// HTML/PDF/WORD 报告扩展字段
	RuleName   string `json:"rule_name"`   // 触发的规则名称
	Suggestion string `json:"suggestion"` // 优化建议
}

// ReportLabels 报告中的国际化标签
type ReportLabels struct {
	AuditSummary      string `json:"audit_summary"`
	ResultStatistics  string `json:"result_statistics"`
	ProblemSQLList    string `json:"problem_sql_list"`
	RuleHitStatistics string `json:"rule_hit_statistics"`
	AuditTime         string `json:"audit_time"`
	DataSource        string `json:"data_source"`
	Schema            string `json:"schema"`
	TotalSQL          string `json:"total_sql"`
	PassRate          string `json:"pass_rate"`
	Score             string `json:"score"`
	AuditLevel        string `json:"audit_level"`
	Number            string `json:"number"`
	SQL               string `json:"sql"`
	AuditStatus       string `json:"audit_status"`
	AuditResult       string `json:"audit_result"`
	ExecStatus        string `json:"exec_status"`
	ExecResult        string `json:"exec_result"`
	RollbackSQL       string `json:"rollback_sql"`
	RuleName          string `json:"rule_name"`
	Description       string `json:"description"`
	Suggestion        string `json:"suggestion"`
	Count             string `json:"count"`
	HitCount          string `json:"hit_count"`
}

// ReportGenerator 报告生成器接口
type ReportGenerator interface {
	// Generate 根据报告数据生成指定格式的文件
	Generate(data *AuditReportData) (*utils.ExportDataResult, error)
	// ReportType 返回生成器对应的导出格式
	ReportType() utils.ExportFormat
}
