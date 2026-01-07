package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// rule template
var (
	DefaultRuleTemplatesDesc   = &i18n.Message{ID: "DefaultRuleTemplatesDesc", Other: "默认规则模板"}
	DefaultTemplatesDesc       = &i18n.Message{ID: "DefaultTemplatesDesc", Other: "%s 默认模板"}
	RuleTemplateName           = &i18n.Message{ID: "RuleTemplateName", Other: "规则模板名"}
	RuleTemplateDesc           = &i18n.Message{ID: "RuleTemplateDesc", Other: "规则模板描述"}
	RuleTemplateInstType       = &i18n.Message{ID: "RuleTemplateInstType", Other: "数据源类型"}
	RuleTemplateRuleName       = &i18n.Message{ID: "RuleTemplateRuleName", Other: "规则名"}
	RuleTemplateRuleVersion    = &i18n.Message{ID: "RuleTemplateRuleVersion", Other: "规则版本"}
	RuleTemplateRuleDesc       = &i18n.Message{ID: "RuleTemplateRuleDesc", Other: "规则描述"}
	RuleTemplateRuleAnnotation = &i18n.Message{ID: "RuleTemplateRuleAnnotation", Other: "规则注解"}
	RuleTemplateRuleLevel      = &i18n.Message{ID: "RuleTemplateRuleLevel", Other: "规则等级"}
	RuleTemplateRuleCategory   = &i18n.Message{ID: "RuleTemplateRuleCategory", Other: "规则分类"}
	RuleTemplateRuleParam      = &i18n.Message{ID: "RuleTemplateRuleParam", Other: "规则参数"}
	RuleTemplateRuleErr        = &i18n.Message{ID: "RuleTemplateRuleErr", Other: "问题"}
)

// rule
var (
	RuleLevelError  = &i18n.Message{ID: "RuleLevelError", Other: "错误"}
	RuleLevelWarn   = &i18n.Message{ID: "RuleLevelWarn", Other: "警告"}
	RuleLevelNotice = &i18n.Message{ID: "RuleLevelNotice", Other: "提示"}
	RuleLevelNormal = &i18n.Message{ID: "RuleLevelNormal", Other: "常规"}
	WordIs          = &i18n.Message{ID: "WordIs", Other: "为"}
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

// sql version
var (
	SqlVersionInvalidStatusReason = &i18n.Message{ID: "SqlVersionInvalidStatusReason", Other: "执行失败：在该工单绑定的SQL版本的阶段上，存在状态为%v的工单，其SQL版本id为%v"}
	SqlVersionExecFailedReason    = &i18n.Message{ID: "SqlVersionExecFailedReason", Other: "工单：%s 上线失败并停止继续上线"}
	SqlVersionReleaseFailedReason = &i18n.Message{ID: "SqlVersionReleaseFailedReason", Other: "工单：%s 发布失败并停止继续发布"}
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
	SQLManageStatusSent           = &i18n.Message{ID: "SQLManageStatusSent", Other: "已推送到其他平台"}
)

// license
var (
	LicenseInstanceNum       = &i18n.Message{ID: "LicenseInstanceNum", Other: "实例数"}
	LicenseUserNum           = &i18n.Message{ID: "LicenseUserNum", Other: "用户数"}
	LicenseWorkDurationDay   = &i18n.Message{ID: "LicenseWorkDurationDay", Other: "授权运行时长(天)"}
	LicenseUnlimited         = &i18n.Message{ID: "LicenseUnlimited", Other: "无限制"}
	LicenseDurationOfRunning = &i18n.Message{ID: "LicenseDurationOfRunning", Other: "已运行时长(天)"}
	LicenseEstimatedMaturity = &i18n.Message{ID: "LicenseEstimatedMaturity", Other: "预计到期时间"}
	LicenseInstanceNumOfType = &i18n.Message{ID: "LicenseInstanceNumOfType", Other: "[%v]类型实例数"}
	LicenseMachineInfo       = &i18n.Message{ID: "LicenseMachineInfo", Other: "机器信息"}
	LicenseMachineInfoOfNode = &i18n.Message{ID: "LicenseMachineInfoOfNode", Other: "节点[%s]机器信息"}
	LicenseSQLEVersion       = &i18n.Message{ID: "LicenseSQLEVersion", Other: "SQLE版本"}
)

// statistic
var (
	StatisticResourceTypeUser = &i18n.Message{ID: "StatisticResourceTypeUser", Other: "用户"}
)

// configuration
var (
	ConfigTestAudit         = &i18n.Message{ID: "ConfigTestAudit", Other: "测试审批"}
	ConfigFeishuTestContent = &i18n.Message{ID: "ConfigFeishuTestContent", Other: "这是一条测试审批,用来测试SQLE飞书审批功能是否正常"}
	ConfigCoding            = &i18n.Message{ID: "ConfigCodingTest", Other: "这是一条测试信息，用来测试SQLE推送事项到Coding平台功能是否正常"}
)

// operation_record
var (
	OprAddRuleTemplateWithName  = &i18n.Message{ID: "OprAddRuleTemplateWithName", Other: "添加规则模板，模板名：%v"}
	OprEditRuleTemplateWithName = &i18n.Message{ID: "OprEditRuleTemplateWithName", Other: "编辑规则模板，模板名：%v"}
	OprDelRuleTemplateWithName  = &i18n.Message{ID: "OprDelRuleTemplateWithName", Other: "删除规则模板，模板名：%v"}

	OprAddGlobalRuleTemplateWithName  = &i18n.Message{ID: "OprAddGlobalRuleTemplateWithName", Other: "创建全局规则模板，模板名：%v"}
	OprEditGlobalRuleTemplateWithName = &i18n.Message{ID: "OprEditGlobalRuleTemplateWithName", Other: "编辑全局规则模板，模板名：%v"}
	OprDelGlobalRuleTemplateWithName  = &i18n.Message{ID: "OprDelGlobalRuleTemplateWithName", Other: "删除全局规则模板，模板名：%v"}

	OprAddAuditPlanWithName  = &i18n.Message{ID: "OprAddAuditPlanWithName", Other: "创建智能扫描任务，任务名：%v"}
	OprEditAuditPlanWithName = &i18n.Message{ID: "OprEditAuditPlanWithName", Other: "编辑智能扫描任务，任务名：%v"}
	OprDelAuditPlanWithName  = &i18n.Message{ID: "OprDelAuditPlanWithName", Other: "删除智能扫描任务，任务名：%v"}

	OprAddSchedulingWorkflowWithNameAndDB = &i18n.Message{ID: "OprAddSchedulingWorkflowWithNameAndDB", Other: "设置定时上线，工单名称：%v, 数据源名: %v"}
	OprDelSchedulingWorkflowWithNameAndDB = &i18n.Message{ID: "OprDelSchedulingWorkflowWithNameAndDB", Other: "取消定时上线，工单名称：%v, 数据源名: %v"}

	OprBatchExecutingWorkflowWithName = &i18n.Message{ID: "OprBatchExecutingWorkflowWithName", Other: "上线工单，工单名称：%v"}
	OprExecutingWorkflowWithNameAndDB = &i18n.Message{ID: "OprExecutingWorkflowWithNameAndDB", Other: "上线工单的单个数据源, 工单名称：%v, 数据源名: %v"}
	OprBatchCancelingWorkflowWithName = &i18n.Message{ID: "OprBatchCancelingWorkflowWithName", Other: "批量取消工单，工单名称：%v"}
	OprCancelingWorkflowWithName      = &i18n.Message{ID: "OprCancelingWorkflowWithName", Other: "取消工单，工单名称：%v"}
	OprApprovingWorkflowWithName      = &i18n.Message{ID: "OprApprovingWorkflowWithName", Other: "审核通过工单，工单名称：%v"}
	OprRejectingWorkflowWithName      = &i18n.Message{ID: "OprRejectingWorkflowWithName", Other: "驳回工单，工单名称：%v"}
	OprCreatingWorkflowWithName       = &i18n.Message{ID: "OprCreatingWorkflowWithName", Other: "创建工单，工单名：%v"}

	OprEditProcedureTemplate = &i18n.Message{ID: "OprEditProcedureTemplate", Other: "编辑流程模板"}
	OprEditDingConfig        = &i18n.Message{ID: "OprEditDingConfig", Other: "修改钉钉配置"}
	OprEditGlobalConfig      = &i18n.Message{ID: "OprEditGlobalConfig", Other: "修改全局配置"}

	OprUpdateFilesOrderWithOrderAndName = &i18n.Message{ID: "OprUpdateFilesOrderWithOrderAndName", Other: "文件上线顺序调整：%s，工单名称：%s"}

	OprActionCreateProject               = &i18n.Message{ID: "OprActionCreateProject", Other: "创建项目"}
	OprActionDeleteProject               = &i18n.Message{ID: "OprActionDeleteProject", Other: "删除项目"}
	OprActionUpdateProject               = &i18n.Message{ID: "OprActionUpdateProject", Other: "编辑项目"}
	OprActionArchiveProject              = &i18n.Message{ID: "OprActionArchiveProject", Other: "冻结项目"}
	OprActionUnarchiveProject            = &i18n.Message{ID: "OprActionUnarchiveProject", Other: "取消冻结项目"}
	OprActionCreateInstance              = &i18n.Message{ID: "OprActionCreateInstance", Other: "创建数据源"}
	OprActionUpdateInstance              = &i18n.Message{ID: "OprActionUpdateInstance", Other: "编辑数据源"}
	OprActionDeleteInstance              = &i18n.Message{ID: "OprActionDeleteInstance", Other: "删除数据源"}
	OprActionCreateProjectRuleTemplate   = &i18n.Message{ID: "OprActionCreateProjectRuleTemplate", Other: "添加规则模版"}
	OprActionDeleteProjectRuleTemplate   = &i18n.Message{ID: "OprActionDeleteProjectRuleTemplate", Other: "删除规则模版"}
	OprActionUpdateProjectRuleTemplate   = &i18n.Message{ID: "OprActionUpdateProjectRuleTemplate", Other: "编辑规则模版"}
	OprActionUpdateWorkflowTemplate      = &i18n.Message{ID: "OprActionUpdateWorkflowTemplate", Other: "编辑流程模版"}
	OprActionCreateAuditPlan             = &i18n.Message{ID: "OprActionCreateAuditPlan", Other: "创建智能扫描任务"}
	OprActionDeleteAuditPlan             = &i18n.Message{ID: "OprActionDeleteAuditPlan", Other: "删除智能扫描任务"}
	OprActionUpdateAuditPlan             = &i18n.Message{ID: "OprActionUpdateAuditPlan", Other: "编辑智能扫描任务"}
	OprActionCreateWorkflow              = &i18n.Message{ID: "OprActionCreateWorkflow", Other: "创建工单"}
	OprActionCancelWorkflow              = &i18n.Message{ID: "OprActionCancelWorkflow", Other: "关闭工单"}
	OprActionApproveWorkflow             = &i18n.Message{ID: "OprActionApproveWorkflow", Other: "审核通过工单"}
	OprActionRejectWorkflow              = &i18n.Message{ID: "OprActionRejectWorkflow", Other: "驳回工单"}
	OprActionExecuteWorkflow             = &i18n.Message{ID: "OprActionExecuteWorkflow", Other: "上线工单"}
	OprActionScheduleWorkflow            = &i18n.Message{ID: "OprActionScheduleWorkflow", Other: "定时上线"}
	OprActionCreateUser                  = &i18n.Message{ID: "OprActionCreateUser", Other: "创建用户"}
	OprActionUpdateUser                  = &i18n.Message{ID: "OprActionUpdateUser", Other: "编辑用户"}
	OprActionDeleteUser                  = &i18n.Message{ID: "OprActionDeleteUser", Other: "删除用户"}
	OprActionCreateGlobalRuleTemplate    = &i18n.Message{ID: "OprActionCreateGlobalRuleTemplate", Other: "创建全局规则模版"}
	OprActionUpdateGlobalRuleTemplate    = &i18n.Message{ID: "OprActionUpdateGlobalRuleTemplate", Other: "编辑全局规则模版"}
	OprActionDeleteGlobalRuleTemplate    = &i18n.Message{ID: "OprActionDeleteGlobalRuleTemplate", Other: "删除全局规则模版"}
	OprActionUpdateDingTalkConfiguration = &i18n.Message{ID: "OprActionUpdateDingTalkConfiguration", Other: "修改钉钉配置"}
	OprActionUpdateSMTPConfiguration     = &i18n.Message{ID: "OprActionUpdateSMTPConfiguration", Other: "修改SMTP配置"}
	OprActionUpdateWechatConfiguration   = &i18n.Message{ID: "OprActionUpdateWechatConfiguration", Other: "修改微信配置"}
	OprActionUpdateSystemVariables       = &i18n.Message{ID: "OprActionUpdateSystemVariables", Other: "修改系统变量"}
	OprActionUpdateLDAPConfiguration     = &i18n.Message{ID: "OprActionUpdateLDAPConfiguration", Other: "修改LDAP配置"}
	OprActionUpdateOAuth2Configuration   = &i18n.Message{ID: "OprActionUpdateOAuth2Configuration", Other: "修改OAuth2配置"}
	OprActionCreateMember                = &i18n.Message{ID: "OprActionCreateMember", Other: "添加成员"}
	OprActionCreateMemberGroup           = &i18n.Message{ID: "OprActionCreateMemberGroup", Other: "添加成员组"}
	OprActionDeleteMember                = &i18n.Message{ID: "OprActionDeleteMember", Other: "删除成员"}
	OprActionDeleteMemberGroup           = &i18n.Message{ID: "OprActionDeleteMemberGroup", Other: "删除成员组"}
	OprActionUpdateMember                = &i18n.Message{ID: "OprActionUpdateMember", Other: "编辑成员"}
	OprActionUpdateMemberGroup           = &i18n.Message{ID: "OprActionUpdateMemberGroup", Other: "编辑成员组"}

	OprOperationTime        = &i18n.Message{ID: "OprOperationTime", Other: "操作时间"}
	OprOperationProjectName = &i18n.Message{ID: "OprOperationProjectName", Other: "项目"}
	OprOperationUserName    = &i18n.Message{ID: "OprOperationUserName", Other: "操作人"}
	OprOperationAction      = &i18n.Message{ID: "OprOperationAction", Other: "操作对象"}
	OprOperationContent     = &i18n.Message{ID: "OprOperationContent", Other: "操作内容"}
	OprOperationStatus      = &i18n.Message{ID: "OprOperationStatus", Other: "状态"}

	OprTypeProject             = &i18n.Message{ID: "OprTypeProject", Other: "项目"}
	OprTypeInstance            = &i18n.Message{ID: "OprTypeInstance", Other: "数据源"}
	OprTypeProjectRuleTemplate = &i18n.Message{ID: "OprTypeProjectRuleTemplate", Other: "项目规则模板"}
	OprTypeWorkflowTemplate    = &i18n.Message{ID: "OprTypeWorkflowTemplate", Other: "流程模板"}
	OprTypeAuditPlan           = &i18n.Message{ID: "OprTypeAuditPlan", Other: "智能扫描任务"}
	OprTypeWorkflow            = &i18n.Message{ID: "OprTypeWorkflow", Other: "工单"}
	OprTypeGlobalUser          = &i18n.Message{ID: "OprTypeGlobalUser", Other: "平台用户"}
	OprTypeGlobalRuleTemplate  = &i18n.Message{ID: "OprTypeGlobalRuleTemplate", Other: "全局规则模板"}
	OprTypeSystemConfiguration = &i18n.Message{ID: "OprTypeSystemConfiguration", Other: "系统配置"}
	OprTypeProjectMember       = &i18n.Message{ID: "OprTypeProjectMember", Other: "项目成员"}

	OprStatusSucceeded = &i18n.Message{ID: "OprStatusSucceeded", Other: "成功"}
	OprStatusFailed    = &i18n.Message{ID: "OprStatusFailed", Other: "失败"}
)

// operation
var (
	OpWorkflowViewOthers  = &i18n.Message{ID: "OpWorkflowViewOthers", Other: "查看他人创建的工单"}
	OpWorkflowSave        = &i18n.Message{ID: "OpWorkflowSave", Other: "创建/编辑工单"}
	OpWorkflowAudit       = &i18n.Message{ID: "OpWorkflowAudit", Other: "审核/驳回工单"}
	OpWorkflowExecute     = &i18n.Message{ID: "OpWorkflowExecute", Other: "上线工单"}
	OpAuditPlanViewOthers = &i18n.Message{ID: "OpAuditPlanViewOthers", Other: "查看他人创建的扫描任务"}
	OpAuditPlanSave       = &i18n.Message{ID: "OpAuditPlanSave", Other: "创建扫描任务"}
	OpSqlQueryQuery       = &i18n.Message{ID: "OpSqlQueryQuery", Other: "SQL查询"}
	OpUnknown             = &i18n.Message{ID: "OpUnknown", Other: "未知动作"}
)

// audit plan
var (
	ApAuditResult    = &i18n.Message{ID: "ApAuditResult", Other: "审核结果"}
	ApSQLStatement   = &i18n.Message{ID: "ApSQLStatement", Other: "SQL语句"}
	ApPriority       = &i18n.Message{ID: "ApPriority", Other: "优先级"}
	ApSchema         = &i18n.Message{ID: "ApSchema", Other: "schema"}
	ApRuleName       = &i18n.Message{ID: "ApRuleName", Other: "审核规则"}
	ApSQLFingerprint = &i18n.Message{ID: "ApSQLFingerprint", Other: "SQL指纹"}
	ApLastSQL        = &i18n.Message{ID: "ApLastSQL", Other: "最后一次匹配到该指纹的语句"}
	ApNum            = &i18n.Message{ID: "ApNum", Other: "数量"}
	ApLastMatchTime  = &i18n.Message{ID: "ApLastMatchTime", Other: "最后匹配时间"}
	ApQueryTimeAvg   = &i18n.Message{ID: "ApQueryTimeAvg", Other: "平均执行时间"}
	ApQueryTimeMax   = &i18n.Message{ID: "ApQueryTimeMax", Other: "最长执行时间"}
	ApRowExaminedAvg = &i18n.Message{ID: "ApRowExaminedAvg", Other: "平均扫描行数"}

	ApDatabase                          = &i18n.Message{ID: "ApDatabase", Other: "数据库"}
	ApMetricNameLockType                = &i18n.Message{ID: "ApMetricNameLockType", Other: "锁类型"}
	ApMetricNameLockMode                = &i18n.Message{ID: "ApMetricNameLockMode", Other: "锁模式"}
	ApMetricEngine                      = &i18n.Message{ID: "ApMetricEngine", Other: "引擎"}
	ApMetricNameTable                   = &i18n.Message{ID: "ApMetricNameTable", Other: "表名"}
	ApMetricNameGrantedLockId           = &i18n.Message{ID: "ApMetricNameGrantedLockId", Other: "持有锁ID"}
	ApMetricNameWaitingLockId           = &i18n.Message{ID: "ApMetricNameWaitingLockId", Other: "等待锁ID"}
	ApMetricNameGrantedLockTrxId        = &i18n.Message{ID: "ApMetricNameGrantedLockTrxId", Other: "持有锁事务ID"}
	ApMetricNameWaitingLockTrxId        = &i18n.Message{ID: "ApMetricNameWaitingLockTrxId", Other: "等待锁事务ID"}
	ApMetricNameTrxStarted              = &i18n.Message{ID: "ApMetricNameTransactionStarted", Other: "持有锁事务开始时间"}
	ApMetricNameTrxWaitStarted          = &i18n.Message{ID: "ApMetricNameTrxWaitStarted", Other: "等待锁事务开始时间"}
	ApMetricNameGrantedLockConnectionId = &i18n.Message{ID: "ApMetricNameGrantedLockConnectionId", Other: "持有锁连接ID"}
	ApMetricNameWaitingLockConnectionId = &i18n.Message{ID: "ApMetricNameWaitingLockConnectionId", Other: "等待锁连接ID"}
	ApMetricNameGrantedLockSql          = &i18n.Message{ID: "ApMetricNameGrantedLockSql", Other: "持有锁SQL"}
	ApMetricNameWaitingLockSql          = &i18n.Message{ID: "ApMetricNameWaitingLockSql", Other: "等待锁SQL"}
	ApMetricNameDBUser                  = &i18n.Message{ID: "ApMetricNameDBUser", Other: "用户"}
	ApMetricUserClientIP                = &i18n.Message{ID: "ApMetricUserClientIP", Other: "客户端IP"}
	ApMetricNameHost                    = &i18n.Message{ID: "ApMetricNameHost", Other: "主机"}
	ApMetricNameMetaName                = &i18n.Message{ID: "ApMetricNameMetaName", Other: "对象名称"}
	ApMetricNameMetaType                = &i18n.Message{ID: "ApMetricNameMetaType", Other: "对象类型"}
	ApMetricNameQueryTimeTotal          = &i18n.Message{ID: "ApMetricNameQueryTimeTotal", Other: "总执行时间(s)"}
	ApMetricNameQueryTimeAvg            = &i18n.Message{ID: "ApMetricNameQueryTimeAvg", Other: "平均执行时间(s)"}
	ApMetricNameQueryTimeTotalMS        = &i18n.Message{ID: "ApMetricNameQueryTimeTotalMS", Other: "总执行时间(ms)"}
	ApMetricNameQueryTimeTotalUS        = &i18n.Message{ID: "ApMetricNameQueryTimeTotalUS", Other: "总执行时间(μs)"}
	ApMetricNameQueryTimeAvgMS          = &i18n.Message{ID: "ApMetricNameQueryTimeAvgMS", Other: "平均执行时间(ms)"}
	ApMetricNameCounter                 = &i18n.Message{ID: "ApMetricNameCounter", Other: "执行次数"}
	ApMetricNameCPUTimeAvg              = &i18n.Message{ID: "ApMetricNameCPUTimeAvg", Other: "平均 CPU 时间(μs)"}
	ApMetricNameLockWaitTimeTotal       = &i18n.Message{ID: "ApMetricNameLockWaitTimeTotal", Other: "锁等待时间(ms)"}
	ApMetricNameLockWaitCounter         = &i18n.Message{ID: "ApMetricNameLockWaitCounter", Other: "锁等待次数"}
	ApMetricNameLockWaitTimeAvg         = &i18n.Message{ID: "ApMetricNameLockWaitTimeAvg", Other: "平均锁等待时间(ms)"}
	ApMetricNameLockWaitTimeMax         = &i18n.Message{ID: "ApMetricNameLockWaitTimeMax", Other: "锁最大等待时间(ms)"}
	ApMetricNameActiveWaitTimeTotal     = &i18n.Message{ID: "ApMetricNameActiveWaitTimeTotal", Other: "活动等待总时间(ms)"}
	ApMetricNameActiveTimeTotal         = &i18n.Message{ID: "ApMetricNameActiveTimeTotal", Other: "活动总时间(ms)"}
	ApMetricNameLastReceiveTimestamp    = &i18n.Message{ID: "ApMetricNameLastReceiveTimestamp", Other: "最后一次匹配到该指纹的时间"}
	ApMetricNameCPUTimeTotal            = &i18n.Message{ID: "ApMetricNameCPUTimeTotal", Other: "CPU时间占用(s)"}
	ApMetricNamePhyReadPageTotal        = &i18n.Message{ID: "ApMetricNamePhyReadPageTotal", Other: "物理读页数"}
	ApMetricNameLogicReadPageTotal      = &i18n.Message{ID: "ApMetricNameLogicReadPageTotal", Other: "逻辑读页数"}
	ApMetricNameQueryTimeMax            = &i18n.Message{ID: "ApMetricNameQueryTimeMax", Other: "最长执行时间(s)"}
	ApMetricNameQueryTimeMaxMS          = &i18n.Message{ID: "ApMetricNameQueryTimeMaxMS", Other: "最长执行时间(ms)"}
	ApMetricNameRowExaminedAvg          = &i18n.Message{ID: "ApMetricNameRowExaminedAvg", Other: "平均扫描行数"}
	ApMetricNameDiskReadTotal           = &i18n.Message{ID: "ApMetricNameDiskReadTotal", Other: "物理读次数"}
	ApMetricNameBufferGetCounter        = &i18n.Message{ID: "ApMetricNameBufferGetCounter", Other: "逻辑读次数"}
	ApMetricNameUserIOWaitTimeTotal     = &i18n.Message{ID: "ApMetricNameUserIOWaitTimeTotal", Other: "I/O等待时间(s)"}
	ApMetricNameIoWaitTimeAvg           = &i18n.Message{ID: "ApMetricNameIoWaitTimeAvg", Other: "平均IO等待时间(毫秒)"}
	ApMetricNameBufferReadAvg           = &i18n.Message{ID: "ApMetricNameBufferReadAvg", Other: "平均逻辑读次数"}
	ApMetricNameDiskReadAvg             = &i18n.Message{ID: "ApMetricNameDiskReadAvg", Other: "平均物理读次数"}
	ApMetricNameFirstQueryAt            = &i18n.Message{ID: "ApMetricNameFirstQueryAt", Other: "首次执行时间"}
	ApMetricNameLastQueryAt             = &i18n.Message{ID: "ApMetricNameLastQueryAt", Other: "最后执行时间"}
	ApMetricNameMaxQueryTime            = &i18n.Message{ID: "ApMetricNameMaxQueryTime", Other: "最长执行时间(s)"}
	ApMetricNameRowsAffectedMax         = &i18n.Message{ID: "ApMetricNameRowsAffectedMax", Other: "最大影响行数"}
	ApMetricNameRowsAffectedAvg         = &i18n.Message{ID: "ApMetricNameRowsAffectedAvg", Other: "平均影响行数"}
	ApMetricNameChecksum                = &i18n.Message{ID: "ApMetricNameChecksum", Other: "校验和"}
	ApMetricNameNoIndexUsedTotal        = &i18n.Message{ID: "ApMetricNameNoIndexUsedTotal", Other: "累计未使用索引次数"}

	ApMetricNameCounterMoreThan        = &i18n.Message{ID: "ApMetricNameCounterMoreThan", Other: "出现次数 > "}
	ApMetricNameQueryTimeAvgMoreThan   = &i18n.Message{ID: "ApMetricNameQueryTimeAvgMoreThan", Other: "平均执行时间(s) > "}
	ApMetricNameRowExaminedAvgMoreThan = &i18n.Message{ID: "ApMetricNameRowExaminedAvgMoreThan", Other: "平均扫描行数 > "}

	ApMetricNameInstance    = &i18n.Message{ID: "ApMetricNameInstance", Other: "节点地址"}
	ApMetricNameMemMax      = &i18n.Message{ID: "ApMetricNameMemMax", Other: "使用的最大内存空间"}
	ApMetricNameDiskMax     = &i18n.Message{ID: "ApMetricNameDiskMax", Other: "使用的最大硬盘空间"}
	ApMetricNameTenantName  = &i18n.Message{ID: "ApMetricNameTenantName", Other: "租户名称"}
	ApMetricNameRequestTime = &i18n.Message{ID: "ApMetricNameRequestTime", Other: "请求时间"}

	ApMetaCustom                       = &i18n.Message{ID: "ApMetaCustom", Other: "自定义"}
	ApMetaMySQLSchemaMeta              = &i18n.Message{ID: "ApMetaMySQLSchemaMeta", Other: "库表元数据"}
	ApMetaMySQLProcesslist             = &i18n.Message{ID: "ApMetaMySQLProcesslist", Other: "processlist 列表"}
	ApMetaAliRdsMySQLSlowLog           = &i18n.Message{ID: "ApMetaAliRdsMySQLSlowLog", Other: "阿里RDS MySQL慢日志"}
	ApMetaAliRdsMySQLAuditLog          = &i18n.Message{ID: "ApMetaAliRdsMySQLAuditLog", Other: "阿里RDS MySQL审计日志"}
	ApMetaBaiduRdsMySQLSlowLog         = &i18n.Message{ID: "ApMetaBaiduRdsMySQLSlowLog", Other: "百度云RDS MySQL慢日志"}
	ApMetaHuaweiRdsMySQLSlowLog        = &i18n.Message{ID: "ApMetaHuaweiRdsMySQLSlowLog", Other: "华为云RDS MySQL慢日志"}
	ApMetaOracleTopSQL                 = &i18n.Message{ID: "ApMetaOracleTopSQL", Other: "Oracle TOP SQL"}
	ApMetaAllAppExtract                = &i18n.Message{ID: "ApMetaAllAppExtract", Other: "应用程序SQL抓取"}
	ApMetaTiDBAuditLog                 = &i18n.Message{ID: "ApMetaTiDBAuditLog", Other: "TiDB审计日志"}
	ApMetaSlowLog                      = &i18n.Message{ID: "ApMetaSlowLog", Other: "慢日志"}
	ApMetaMDBSlowLog                   = &i18n.Message{ID: "ApMetaMDBSlowLog", Other: "慢日志（监控库）"}
	ApMetaTopSQL                       = &i18n.Message{ID: "ApMetaTopSQL", Other: "Top SQL"}
	ApMetaDB2TopSQL                    = &i18n.Message{ID: "ApMetaDB2TopSQL", Other: "DB2 Top SQL"}
	ApMetaSchemaMeta                   = &i18n.Message{ID: "ApMetaSchemaMeta", Other: "库表元数据"}
	ApMetaDistributedLock              = &i18n.Message{ID: "ApMetaDistributedLock", Other: "分布式锁"}
	ApMetaDmTopSQL                     = &i18n.Message{ID: "ApMetaDmTopSQL", Other: "DM TOP SQL"}
	ApMetaObForOracleTopSQL            = &i18n.Message{ID: "ApMetaObForOracleTopSQL", Other: "OceanBase For Oracle TOP SQL"}
	ApMetaOceanBaseForMySQLFullCollect = &i18n.Message{ID: "ApMetaOceanBaseForMySQLFullCollect", Other: "全量采集"}
	ApMetaPostgreSQLTopSQL             = &i18n.Message{ID: "ApMetaPostgreSQLTopSQL", Other: "TOP SQL"}
	ApMetaGoldenDBTopSQL               = &i18n.Message{ID: "ApMetaGoldenDBTopSQL", Other: "GoldenDB TOP SQL"}
	ApMetaTiDBTopSQL                   = &i18n.Message{ID: "ApMetaTiDBTopSQL", Other: "TiDB TOP SQL"}
	ApMetaMySQLTopSQL                  = &i18n.Message{ID: "ApMetaMySQLTopSQL", Other: "MySQL TOP SQL"}
	ApMetricQueryTimeAvg               = &i18n.Message{ID: "ApMetricQueryTimeAvg", Other: "平均查询时间(s)"}
	ApMetricRowExaminedAvg             = &i18n.Message{ID: "ApMetricRowExaminedAvg", Other: "平均扫描行数"}
	ApMetaPerformanceCollect           = &i18n.Message{ID: "ApMetaPerformanceCollect", Other: "数据源性能指标"}
	ApMetaPerformanceCollectTips       = &i18n.Message{ID: "ApMetaPerformanceCollectTips", Other: "性能指标采集将产生较大性能开销,请谨慎开启。开启后,系统将持续采集该数据源的性能数据(如QPS、连接数等)，并生成性能趋势图表，体现在性能洞察页面。"}
	ApMetaCollectTime                  = &i18n.Message{ID: "ApMetaCollectTime", Other: "采集时间"}
	ApMetaThreadsConnected             = &i18n.Message{ID: "ApMetaThreadsConnected", Other: "线程数"}
	ApMetaQPS                          = &i18n.Message{ID: "ApMetaQueries", Other: "QPS"}
	ApMetricNameFullTableScanCount     = &i18n.Message{ID: "ApMetricNameFullTableScanCount", Other: "全表扫描次数"}

	ApPriorityHigh = &i18n.Message{ID: "ApPriorityHigh", Other: "高优先级"}

	ParamCollectIntervalMinute           = &i18n.Message{ID: "ParamCollectIntervalMinute", Other: "采集周期（分钟）"}
	ParamTopN                            = &i18n.Message{ID: "ParamTopN", Other: "Top N"}
	ParamIndicator                       = &i18n.Message{ID: "ParamIndicator", Other: "关注指标"}
	ParamCollectIntervalMinuteMySQL      = &i18n.Message{ID: "ParamCollectIntervalMinuteMySQL", Other: "采集周期（分钟，仅对 mysql.slow_log 有效）"}
	ParamSlowLogCollectInput             = &i18n.Message{ID: "ParamSlowLogCollectInput", Other: "采集来源"}
	ParamFirstSqlsScrappedHours          = &i18n.Message{ID: "ParamFirstSqlsScrappedHours", Other: "启动任务时拉取慢日志时间范围(单位:小时，仅对 mysql.slow_log 有效)"}
	ParamCollectIntervalMinuteOracle     = &i18n.Message{ID: "ParamCollectIntervalMinuteOracle", Other: "采集周期（分钟）"}
	ParamOrderByColumn                   = &i18n.Message{ID: "ParamOrderByColumn", Other: "V$SQLAREA中的排序字段"}
	ParamOrderByColumnGeneric            = &i18n.Message{ID: "ParamOrderByColumnGeneric", Other: "排序字段"}
	ParamCollectIntervalSecond           = &i18n.Message{ID: "ParamCollectIntervalSecond", Other: "采集周期（秒）"}
	ParamSQLMinSecond                    = &i18n.Message{ID: "ParamSQLMinSecond", Other: "SQL 最小执行时间（秒）"}
	ParamCollectView                     = &i18n.Message{ID: "ParamCollectView", Other: "是否采集视图信息"}
	ParamDBInstanceId                    = &i18n.Message{ID: "ParamDBInstanceId", Other: "实例ID"}
	ParamAccessKeyId                     = &i18n.Message{ID: "ParamAccessKeyId", Other: "Access Key ID"}
	ParamAccessKeySecret                 = &i18n.Message{ID: "ParamAccessKeySecret", Other: "Access Key Secret"}
	ParamFirstCollectDurationWithMaxDays = &i18n.Message{ID: "ParamFirstCollectDurationWithMaxDays", Other: "启动任务时拉取日志时间范围(单位:小时,最大%d天)"}
	ParamRdsPath                         = &i18n.Message{ID: "ParamRdsPath", Other: "RDS Open API地址"}
	ParamProjectId                       = &i18n.Message{ID: "ParamProjectId", Other: "项目ID"}
	ParamRegion                          = &i18n.Message{ID: "ParamRegion", Other: "当前RDS实例所在的地区（示例：cn-east-2）"}
	ParamTimeSpan                        = &i18n.Message{ID: "ParamTimeSpan", Other: "时间跨度（小时）"}
	ParamInstance                        = &i18n.Message{ID: "ParamInstance", Other: "节点地址（0 代表所有节点）"}
	ParamKpiType                         = &i18n.Message{ID: "ParamkpiType", Other: "指标"}

	EnumKpiTypeQueryTime          = &i18n.Message{ID: "EnumkpiTypeQueryTime", Other: "执行时间"}
	EnumKpiTypeMemMax             = &i18n.Message{ID: "EnumKpiTypeMemMax", Other: "使用的最大内存空间"}
	EnumKpiTypeDiskMax            = &i18n.Message{ID: "EnumKpiTypeDiskMax", Other: "使用的最大硬盘空间"}
	EnumKpiTypeExecuteCount       = &i18n.Message{ID: "EnumKpiTypeExecuteCount", Other: "执行次数"}
	EnumKpiTypeFullTableScanCount = &i18n.Message{ID: "EnumKpiTypeFullTableScan", Other: "全表扫描次数"}
	EnumKpiTypeLockWaitTotal      = &i18n.Message{ID: "EnumKpiTypeLockWaitTotal", Other: "累计锁等待时间"}
	EnumKpiTypeNoIndexUsedTotal   = &i18n.Message{ID: "EnumKpiTypeNoIndexUsedTotal", Other: "累计未使用索引次数"}

	EnumSlowLogFileSource  = &i18n.Message{ID: "EnumSlowLogFileSource", Other: "从slow.log 文件采集,需要适配scanner"}
	EnumSlowLogTableSource = &i18n.Message{ID: "EnumSlowLogTableSource", Other: "从mysql.slow_log 表采集"}

	OperatorGreaterThan = &i18n.Message{ID: "OperatorGreaterThan", Other: "大于"}
	OperatorEqualTo     = &i18n.Message{ID: "OperatorEqualTo", Other: "等于"}
	OperatorLessThan    = &i18n.Message{ID: "OperatorLessThan", Other: "小于"}

	OperationParamAuditLevel = &i18n.Message{ID: "OperationParamAuditLevel", Other: "触发审核级别"}
)

var (
	PipelineCmdUsage = &i18n.Message{ID: "PipelineCmdUsage", Other: "#使用方法#\n1. 确保运行该命令的用户具有scannerd的执行权限。\n2. 在scannerd文件所在目录执行启动命令。\n#启动命令#\n"}
)

// notification
var (
	NotifyWorkflowStepTypeSQLExecute       = &i18n.Message{ID: "WorkflowStepTypeSQLExecute", Other: "上线"}
	NotifyWorkflowStepTypeSQLAudit         = &i18n.Message{ID: "WorkflowStepTypeSQLAudit", Other: "审批"}
	NotifyWorkflowNotifyTypeWaiting        = &i18n.Message{ID: "WorkflowNotifyTypeWaiting", Other: "SQL工单待%s"}
	NotifyWorkflowNotifyTypeReject         = &i18n.Message{ID: "WorkflowNotifyTypeReject", Other: "SQL工单已被驳回"}
	NotifyWorkflowNotifyTypeExecuteSuccess = &i18n.Message{ID: "WorkflowNotifyTypeExecuteSuccess", Other: "SQL工单上线成功"}
	NotifyWorkflowNotifyTypeExecuteFail    = &i18n.Message{ID: "WorkflowNotifyTypeExecuteFail", Other: "SQL工单上线失败"}
	NotifyWorkflowNotifyTypeComplete       = &i18n.Message{ID: "WorkflowNotifyTypeComplete", Other: "SQL工单标记为人工上线"}
	NotifyWorkflowNotifyTypeCancel         = &i18n.Message{ID: "WorkflowNotifyTypeCancel", Other: "SQL工单已关闭"}
	NotifyWorkflowNotifyTypeDefault        = &i18n.Message{ID: "WorkflowNotifyTypeDefault", Other: "SQL工单未知请求"}

	NotifyAuditPlanSubject  = &i18n.Message{ID: "NotifyAuditPlanSubject", Other: "SQLE扫描任务[%v]扫描结果[%v]"}
	NotifyAuditPlanBody     = &i18n.Message{ID: "NotifyAuditPlanBody", Other: "\n- 扫描任务: %v\n- 审核时间: %v\n- 审核类型: %v\n- 数据源: %v\n- 数据库名: %v\n- 审核得分: %v\n- 审核通过率：%v\n- 审核结果等级: %v%v"}
	NotifyAuditPlanBodyLink = &i18n.Message{ID: "NotifyAuditPlanBodyLink", Other: "\n- 扫描任务链接: %v"}

	NotifyManageRecordSubject    = &i18n.Message{ID: "NotifyManageRecordSubject", Other: "SQL管控记录"}
	NotifyManageRecordBodyLink   = &i18n.Message{ID: "NotifyManageRecordBodyLink", Other: "\n- SQL管控记录链接: %v\n"}
	NotifyManageRecordBodyRecord = &i18n.Message{ID: "NotifyManageRecordBodyRecord", Other: "- SQL ID: %v\n- 所在数据源名称: %v\n- 环境属性: %v\n- SQL: %v\n- 触发规则级别: %v\n- SQL审核建议: %v\n================================"}
	NotifyManageRecordBodyTime   = &i18n.Message{ID: "NotifyManageRecordBodyTime", Other: "记录时间周期: %v - %v"}
	NotifyManageRecordBodyProj   = &i18n.Message{ID: "NotifyManageRecordBodyProj", Other: "所属项目: %v"}

	NotifyWorkflowBodyHead              = &i18n.Message{ID: "NotifyWorkflowBodyHead", Other: "\n📋 工单主题: %v\n📍 所属项目: %v\n🆔 工单ID: %v\n📝 工单描述: %v\n👤 申请人: %v\n⏰ 创建时间: %v\n"}
	NotifyWorkflowBodyWorkFlowErr       = &i18n.Message{ID: "NotifyWorkflowBodyWorkFlowErr", Other: "❌ 读取工单任务内容失败，请通过SQLE界面确认工单状态"}
	NotifyWorkflowBodyLink              = &i18n.Message{ID: "NotifyWorkflowBodyLink", Other: "🔗 工单链接: %v"}
	NotifyWorkflowBodyConfigUrl         = &i18n.Message{ID: "NotifyWorkflowBodyConfigUrl", Other: "请在系统设置-全局配置中补充全局url"}
	NotifyWorkflowBodyInstanceErr       = &i18n.Message{ID: "NotifyWorkflowBodyInstanceErr", Other: "❌ 获取数据源实例失败: %v\n"}
	NotifyWorkflowBodyInstanceAndSchema = &i18n.Message{ID: "NotifyWorkflowBodyInstanceAndSchema", Other: "🗄️ 数据源: %v\n📊 Schema: %v\n"}
	NotifyWorkflowBodyStartEnd          = &i18n.Message{ID: "NotifyWorkflowBodyStartEnd", Other: "▶️ 上线开始时间: %v\n◀️ 上线结束时间: %v\n"}
	NotifyWorkflowBodyReason            = &i18n.Message{ID: "NotifyWorkflowBodyReason", Other: "❌ 驳回原因: %v\n"}
	NotifyWorkflowBodyReport            = &i18n.Message{ID: "NotifyWorkflowBodyReport", Other: "✅ 工单审核得分: %v\n"}
	NotifyWorkflowBodyCancel            = &i18n.Message{ID: "NotifyWorkflowBodyCancel", Other: "🚫 工单已关闭\n"}
	NotifyWorkflowBodyComplete          = &i18n.Message{ID: "NotifyWorkflowBodyComplete", Other: "✅ 工单已标记为人工上线\n"}
)
