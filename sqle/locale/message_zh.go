package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// rule template
var (
	DefaultRuleTemplatesDesc = &i18n.Message{ID: "DefaultRuleTemplatesDesc", Other: "默认规则模板"}
	DefaultTemplatesDesc     = &i18n.Message{ID: "DefaultTemplatesDesc", Other: "%s 默认模板"}
)

// task
var (
	TaskStatusExecuting        = &i18n.Message{ID: "TaskStatusExecuting", Other: "正在上线"}
	TaskStatusExecuteSucceeded = &i18n.Message{ID: "TaskStatusExecuteSucceeded", Other: "上线成功"}
	TaskStatusExecuteFailed    = &i18n.Message{ID: "TaskStatusExecuteFailed", Other: "上线失败"}
	TaskStatusManuallyExecuted = &i18n.Message{ID: "TaskStatusManuallyExecuted", Other: "手动上线"}

	FileOrderMethodPrefixNumAsc = &i18n.Message{ID: "FileOrderMethodPrefixNumAsc", Other: "文件名前缀数字升序"}
	FileOrderMethodSuffixNumAsc = &i18n.Message{ID: "FileOrderMethodSuffixNumAsc", Other: "文件名后缀数字升序"}

	SQLAuditStatusInitialized = &i18n.Message{ID: "SQLAuditStatusInitialized", Other: "未审核"}
	SQLAuditStatusDoing       = &i18n.Message{ID: "SQLAuditStatusDoing", Other: "正在审核"}
	SQLAuditStatusFinished    = &i18n.Message{ID: "SQLAuditStatusFinished", Other: "审核完成"}
	SQLAuditStatusUnknown     = &i18n.Message{ID: "SQLAuditStatusUnknown", Other: "未知状态"}

	SQLAuditResultDescPass = &i18n.Message{ID: "SQLAuditResultDescPass", Other: "审核通过"}

	SQLExecuteStatusInitialized      = &i18n.Message{ID: "SQLExecuteStatusInitialized", Other: "准备执行"}
	SQLExecuteStatusDoing            = &i18n.Message{ID: "SQLExecuteStatusDoing", Other: "正在执行"}
	SQLExecuteStatusFailed           = &i18n.Message{ID: "SQLExecuteStatusFailed", Other: "执行失败"}
	SQLExecuteStatusSucceeded        = &i18n.Message{ID: "SQLExecuteStatusSucceeded", Other: "执行成功"}
	SQLExecuteStatusManuallyExecuted = &i18n.Message{ID: "SQLExecuteStatusManuallyExecuted", Other: "人工执行"}
	SQLExecuteStatusUnknown          = &i18n.Message{ID: "SQLExecuteStatusUnknown", Other: "未知"}

	TaskSQLReportIndex       = &i18n.Message{ID: "TaskSQLReportIndex", Other: "序号"}
	TaskSQLReportSQL         = &i18n.Message{ID: "TaskSQLReportSQL", Other: "SQL"}
	TaskSQLReportAuditStatus = &i18n.Message{ID: "TaskSQLReportAuditStatus", Other: "SQL审核状态"}
	TaskSQLReportAuditResult = &i18n.Message{ID: "TaskSQLReportAuditResult", Other: "SQL审核结果"}
	TaskSQLReportExecStatus  = &i18n.Message{ID: "TaskSQLReportExecStatus", Other: "SQL执行状态"}
	TaskSQLReportExecResult  = &i18n.Message{ID: "TaskSQLReportExecResult", Other: "SQL执行结果"}
	TaskSQLReportRollbackSQL = &i18n.Message{ID: "TaskSQLReportRollbackSQL", Other: "SQL对应的回滚语句"}
	TaskSQLReportDescription = &i18n.Message{ID: "TaskSQLReportDescription", Other: "SQL描述"}
)

// workflow
var (
	WorkflowStepStateApprove = &i18n.Message{ID: "WorkflowStepStateApprove", Other: "通过"}
	WorkflowStepStateReject  = &i18n.Message{ID: "WorkflowStepStateReject", Other: "驳回"}

	WorkflowStatusWaitForAudit     = &i18n.Message{ID: "WorkflowStatusWaitForAudit", Other: "待审核"}
	WorkflowStatusWaitForExecution = &i18n.Message{ID: "WorkflowStatusWaitForExecution", Other: "待上线"}
	WorkflowStatusReject           = &i18n.Message{ID: "WorkflowStatusReject", Other: "已驳回"}
	WorkflowStatusCancel           = &i18n.Message{ID: "WorkflowStatusCancel", Other: "已关闭"}
	WorkflowStatusExecuting        = &i18n.Message{ID: "WorkflowStatusExecuting", Other: "正在上线"}
	WorkflowStatusExecFailed       = &i18n.Message{ID: "WorkflowStatusExecFailed", Other: "上线失败"}
	WorkflowStatusFinish           = &i18n.Message{ID: "WorkflowStatusFinish", Other: "上线成功"}

	WFExportWorkflowNumber      = &i18n.Message{ID: "ExportWorkflowNumber", Other: "工单编号"}
	WFExportWorkflowName        = &i18n.Message{ID: "ExportWorkflowName", Other: "工单名称"}
	WFExportWorkflowDescription = &i18n.Message{ID: "ExportWorkflowDescription", Other: "工单描述"}
	WFExportDataSource          = &i18n.Message{ID: "ExportDataSource", Other: "数据源"}
	WFExportCreateTime          = &i18n.Message{ID: "ExportCreateTime", Other: "创建时间"}
	WFExportCreator             = &i18n.Message{ID: "ExportCreator", Other: "创建人"}
	WFExportTaskOrderStatus     = &i18n.Message{ID: "ExportTaskOrderStatus", Other: "工单状态"}
	WFExportOperator            = &i18n.Message{ID: "ExportOperator", Other: "操作人"}
	WFExportExecutionTime       = &i18n.Message{ID: "ExportExecutionTime", Other: "工单执行时间"}
	WFExportSQLContent          = &i18n.Message{ID: "ExportSQLContent", Other: "具体执行SQL内容"}
	WFExportNode1Auditor        = &i18n.Message{ID: "ExportNode1Auditor", Other: "[节点1]审核人"}
	WFExportNode1AuditTime      = &i18n.Message{ID: "ExportNode1AuditTime", Other: "[节点1]审核时间"}
	WFExportNode1AuditResult    = &i18n.Message{ID: "ExportNode1AuditResult", Other: "[节点1]审核结果"}
	WFExportNode2Auditor        = &i18n.Message{ID: "ExportNode2Auditor", Other: "[节点2]审核人"}
	WFExportNode2AuditTime      = &i18n.Message{ID: "ExportNode2AuditTime", Other: "[节点2]审核时间"}
	WFExportNode2AuditResult    = &i18n.Message{ID: "ExportNode2AuditResult", Other: "[节点2]审核结果"}
	WFExportNode3Auditor        = &i18n.Message{ID: "ExportNode3Auditor", Other: "[节点3]审核人"}
	WFExportNode3AuditTime      = &i18n.Message{ID: "ExportNode3AuditTime", Other: "[节点3]审核时间"}
	WFExportNode3AuditResult    = &i18n.Message{ID: "ExportNode3AuditResult", Other: "[节点3]审核结果"}
	WFExportNode4Auditor        = &i18n.Message{ID: "ExportNode4Auditor", Other: "[节点4]审核人"}
	WFExportNode4AuditTime      = &i18n.Message{ID: "ExportNode4AuditTime", Other: "[节点4]审核时间"}
	WFExportNode4AuditResult    = &i18n.Message{ID: "ExportNode4AuditResult", Other: "[节点4]审核结果"}
	WFExportExecutor            = &i18n.Message{ID: "ExportExecutor", Other: "上线人"}
	WFExportExecutionStartTime  = &i18n.Message{ID: "ExportExecutionStartTime", Other: "上线开始时间"}
	WFExportExecutionEndTime    = &i18n.Message{ID: "ExportExecutionEndTime", Other: "上线结束时间"}
	WFExportExecutionStatus     = &i18n.Message{ID: "ExportExecutionStatus", Other: "上线结果"}
)

// audit plan
var (
	APExportTaskName         = &i18n.Message{ID: "APExportTaskName", Other: "扫描任务名称"}
	APExportGenerationTime   = &i18n.Message{ID: "APExportGenerationTime", Other: "报告生成时间"}
	APExportResultRating     = &i18n.Message{ID: "APExportResultRating", Other: "审核结果评分"}
	APExportApprovalRate     = &i18n.Message{ID: "APExportApprovalRate", Other: "审核通过率"}
	APExportBelongingProject = &i18n.Message{ID: "APExportBelongingProject", Other: "所属项目"}
	APExportCreator          = &i18n.Message{ID: "APExportCreator", Other: "扫描任务创建人"}
	APExportType             = &i18n.Message{ID: "APExportType", Other: "扫描任务类型"}
	APExportDbType           = &i18n.Message{ID: "APExportDbType", Other: "数据库类型"}
	APExportDatabase         = &i18n.Message{ID: "APExportDatabase", Other: "审核的数据库"}

	APExportNumber      = &i18n.Message{ID: "APExportNumber", Other: "编号"}
	APExportAuditResult = &i18n.Message{ID: "APExportAuditResult", Other: "审核结果"}
)

// sql audit record
var (
	AuditRecordTagFull      = &i18n.Message{ID: "AuditRecordTagFull", Other: "全量"}
	AuditRecordTagIncrement = &i18n.Message{ID: "AuditRecordTagIncrement", Other: "增量"}
)

// sql manager
var (
	SMExportTotalSQLCount     = &i18n.Message{ID: "SMExportTotalSQLCount", Other: "SQL总数"}
	SMExportProblemSQLCount   = &i18n.Message{ID: "SMExportProblemSQLCount", Other: "问题SQL数"}
	SMExportOptimizedSQLCount = &i18n.Message{ID: "SMExportOptimizedSQLCount", Other: "已优化SQL数"}

	SMExportSQLFingerprint = &i18n.Message{ID: "SMExportSQLFingerprint", Other: "SQL指纹"}
	SMExportSQL            = &i18n.Message{ID: "SMExportSQL", Other: "SQL"}
	SMExportSource         = &i18n.Message{ID: "SMExportSource", Other: "来源"}
	SMExportDataSource     = &i18n.Message{ID: "SMExportDataSource", Other: "数据源"}
	SMExportSCHEMA         = &i18n.Message{ID: "SMExportSCHEMA", Other: "SCHEMA"}
	SMExportAuditResult    = &i18n.Message{ID: "SMExportAuditResult", Other: "审核结果"}
	SMExportEndpoint       = &i18n.Message{ID: "SMExportEndpoint", Other: "端点信息"}
	SMExportPersonInCharge = &i18n.Message{ID: "SMExportPersonInCharge", Other: "负责人"}
	SMExportState          = &i18n.Message{ID: "SMExportState", Other: "状态"}
	SMExportRemarks        = &i18n.Message{ID: "SMExportRemarks", Other: "备注"}

	SQLManageSourceSqlAuditRecord = &i18n.Message{ID: "SQLManageSourceSqlAuditRecord", Other: "SQL审核"}
	SQLManageSourceAuditPlan      = &i18n.Message{ID: "SQLManageSourceAuditPlan", Other: "智能扫描"}
	SQLManageStatusUnhandled      = &i18n.Message{ID: "SQLManageStatusUnhandled", Other: "未处理"}
	SQLManageStatusSolved         = &i18n.Message{ID: "SQLManageStatusSolved", Other: "已解决"}
	SQLManageStatusIgnored        = &i18n.Message{ID: "SQLManageStatusIgnored", Other: "已忽略"}
	SQLManageStatusManualAudited  = &i18n.Message{ID: "SQLManageStatusManualAudited", Other: "已人工审核"}
)
