package workflowexport

import (
	"context"
	"fmt"
	"strings"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// MaxExportWorkflowCount is the max number of workflows allowed in one export.
const MaxExportWorkflowCount = 10000

// Layout selects export column layout. LayoutProject is frozen for in-project export.
type Layout int

const (
	// LayoutProject is the frozen in-project export layout (includeProjectName=false semantics).
	LayoutProject Layout = iota
	// LayoutGlobalSQLRelease is global export for workflow_type=sql_release.
	LayoutGlobalSQLRelease
	// LayoutGlobalDataExport is global export for workflow_type=data_export.
	LayoutGlobalDataExport
	// LayoutGlobalCommon is global export when workflow_type is empty (public columns).
	LayoutGlobalCommon
)

const (
	timeLayout        = "2006-01-02 15:04:05"
	projectFixedNodes = 4
	minAuditNodes     = 1
)

var workflowStepStateMap = map[string]*i18n.Message{
	model.WorkflowStepStateApprove: locale.WorkflowStepStateApprove,
	model.WorkflowStepStateReject:  locale.WorkflowStepStateReject,
}

var executeStateMap = map[string]*i18n.Message{
	model.TaskStatusExecuting:        locale.TaskStatusExecuting,
	model.TaskStatusExecuteSucceeded: locale.TaskStatusExecuteSucceeded,
	model.TaskStatusExecuteFailed:    locale.TaskStatusExecuteFailed,
	model.TaskStatusManuallyExecuted: locale.TaskStatusManuallyExecuted,
}

var projectNodeAuditorMsgs = []*i18n.Message{
	locale.WFExportNode1Auditor, locale.WFExportNode2Auditor, locale.WFExportNode3Auditor, locale.WFExportNode4Auditor,
}
var projectNodeAuditTimeMsgs = []*i18n.Message{
	locale.WFExportNode1AuditTime, locale.WFExportNode2AuditTime, locale.WFExportNode3AuditTime, locale.WFExportNode4AuditTime,
}
var projectNodeAuditResultMsgs = []*i18n.Message{
	locale.WFExportNode1AuditResult, locale.WFExportNode2AuditResult, locale.WFExportNode3AuditResult, locale.WFExportNode4AuditResult,
}

var unifiedStatusMsg = map[string]*i18n.Message{
	"pending_approval": locale.WFExportUnifiedStatusPendingApproval,
	"pending_action":   locale.WFExportUnifiedStatusPendingAction,
	"in_progress":      locale.WFExportUnifiedStatusInProgress,
	"exporting":        locale.WFExportUnifiedStatusExporting,
	"rejected":         locale.WFExportUnifiedStatusRejected,
	"cancelled":        locale.WFExportUnifiedStatusCancelled,
	"failed":           locale.WFExportUnifiedStatusFailed,
	"completed":        locale.WFExportUnifiedStatusCompleted,
	"unknown":          locale.WFExportUnifiedStatusUnknown,
}

// BuildHeader builds localized export headers for LayoutProject.
// When includeProjectName is true, "项目名称" is prepended (legacy global path; prefer BuildGlobal*).
func BuildHeader(ctx context.Context, includeProjectName bool) []string {
	header := buildProjectHeader(ctx)
	if includeProjectName {
		projectName := locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportProjectName)
		header = append([]string{projectName}, header...)
	}
	return header
}

// BuildHeaderForLayout builds headers for the given layout.
func BuildHeaderForLayout(ctx context.Context, layout Layout, maxAuditNodes int) []string {
	switch layout {
	case LayoutGlobalSQLRelease:
		return buildGlobalSQLReleaseHeader(ctx, normalizeMaxAuditNodes(maxAuditNodes))
	case LayoutGlobalDataExport:
		return buildGlobalDataExportHeader(ctx, normalizeMaxAuditNodes(maxAuditNodes))
	case LayoutGlobalCommon:
		return buildGlobalCommonHeader(ctx)
	default:
		return buildProjectHeader(ctx)
	}
}

// BuildRows assembles LayoutProject rows. When includeProjectName is true, project name is prepended.
func BuildRows(ctx context.Context, workflowIDs []string, includeProjectName bool, projectNameByUID map[string]string) ([][]string, error) {
	return buildProjectRows(ctx, workflowIDs, includeProjectName, projectNameByUID)
}

// BuildGlobalSQLReleaseExport builds header+rows for global sql_release (dynamic audit nodes, SQL last).
func BuildGlobalSQLReleaseExport(ctx context.Context, workflowIDs []string, projectNameByUID map[string]string) ([]string, [][]string, error) {
	workflows, err := loadWorkflowsForExport(ctx, workflowIDs)
	if err != nil {
		return nil, nil, err
	}
	maxN := maxAuditNodesFromWorkflows(workflows)
	header := buildGlobalSQLReleaseHeader(ctx, maxN)
	rows := make([][]string, 0)
	for _, workflow := range workflows {
		projectName := ""
		if projectNameByUID != nil {
			projectName = projectNameByUID[string(workflow.ProjectId)]
		}
		for _, instanceRecord := range workflow.Record.InstanceRecords {
			instanceName := instanceNameOf(instanceRecord)
			row := []string{
				projectName,
				workflow.WorkflowId,
				workflow.Subject,
				workflow.Desc,
				instanceName,
				workflow.Model.CreatedAt.Format(timeLayout),
				dms.GetUserNameWithDelTag(workflow.CreateUserId),
				locale.Bundle.LocalizeMsgByCtx(ctx, model.WorkflowStatus[workflow.Record.Status]),
			}
			row = append(row, getAuditListN(ctx, workflow, maxN)...)
			row = append(row,
				dms.GetUserNameWithDelTag(instanceRecord.ExecutionUserId),
				instanceRecord.Task.TaskExecStartAt(),
				instanceRecord.Task.TaskExecEndAt(),
				locale.Bundle.LocalizeMsgByCtx(ctx, executeStateMap[instanceRecord.Task.Status]),
				getExecuteSqlList(instanceRecord.Task.ExecuteSQLs),
			)
			rows = append(rows, row)
		}
	}
	return header, rows, nil
}

// DataExportExportRecord is the minimal record needed to build global data-export / common rows.
type DataExportExportRecord struct {
	WorkflowID     string
	WorkflowName   string
	Description    string
	ProjectName    string
	ProjectUID     string
	DBServiceNames []string
	CreatedAt      string
	CreatorName    string
	UnifiedStatus  string
}

// FromListDataExportWorkflow maps a DMS list item to export record.
func FromListDataExportWorkflow(r *dmsV1.ListDataExportWorkflow) DataExportExportRecord {
	out := DataExportExportRecord{
		WorkflowID:    r.WorkflowID,
		WorkflowName:  r.WorkflowName,
		Description:   r.Description,
		CreatedAt:     r.CreatedAt.Format(timeLayout),
		CreatorName:   r.Creater.Name,
		UnifiedStatus: string(mapDataExportStatus(r.Status)),
	}
	if r.ProjectInfo != nil {
		out.ProjectName = r.ProjectInfo.ProjectName
		out.ProjectUID = r.ProjectInfo.ProjectUid
	}
	for _, db := range r.DBServiceInfos {
		if db == nil {
			continue
		}
		out.DBServiceNames = append(out.DBServiceNames, db.DBServiceName)
	}
	return out
}

func mapDataExportStatus(status dmsV1.DataExportWorkflowStatus) string {
	switch status {
	case dmsV1.DataExportWorkflowStatusWaitForApprove:
		return "pending_approval"
	case dmsV1.DataExportWorkflowStatusWaitForExport:
		return "pending_action"
	case dmsV1.DataExportWorkflowStatusWaitForExporting:
		return "exporting"
	case dmsV1.DataExportWorkflowStatusRejected:
		return "rejected"
	case dmsV1.DataExportWorkflowStatusCancel:
		return "cancelled"
	case dmsV1.DataExportWorkflowStatusFailed:
		return "failed"
	case dmsV1.DataExportWorkflowStatusFinish:
		return "completed"
	default:
		return "unknown"
	}
}

// mapSQLWorkflowStatus mirrors dashboard.MapSQLWorkflowStatusToUnified for common-column export.
func mapSQLWorkflowStatus(native string) string {
	switch native {
	case model.WorkflowStatusWaitForAudit:
		return "pending_approval"
	case model.WorkflowStatusWaitForExecution:
		return "pending_action"
	case model.WorkflowStatusReject:
		return "rejected"
	case model.WorkflowStatusCancel:
		return "cancelled"
	case model.WorkflowStatusExecuting:
		return "in_progress"
	case model.WorkflowStatusExecFailed:
		return "failed"
	case model.WorkflowStatusFinish:
		return "completed"
	default:
		return "unknown"
	}
}

// BuildGlobalDataExportExport builds header+rows for global data_export (§4.4.3).
func BuildGlobalDataExportExport(ctx context.Context, records []DataExportExportRecord) ([]string, [][]string) {
	maxN := minAuditNodes
	header := buildGlobalDataExportHeader(ctx, maxN)
	rows := make([][]string, 0, len(records))
	for _, r := range records {
		row := []string{
			r.ProjectName,
			r.WorkflowID,
			r.WorkflowName,
			r.Description,
			strings.Join(r.DBServiceNames, ","),
			r.CreatedAt,
			r.CreatorName,
			localizeUnifiedStatus(ctx, r.UnifiedStatus),
		}
		row = append(row, emptyAuditNodes(maxN)...)
		row = append(row, "", "") // 导出执行时间、导出结果 — 列表无字段，留空
		rows = append(rows, row)
	}
	return header, rows
}

// CommonExportRow is one row for LayoutGlobalCommon (one workflow per row).
type CommonExportRow struct {
	ProjectName  string
	WorkflowType string
	WorkflowID   string
	WorkflowName string
	CreatorName  string
	CreatedAt    string
	Status       string
	DataSource   string
	SortKey      string
}

// BuildGlobalCommonExport builds public-column header+rows for empty workflow_type.
func BuildGlobalCommonExport(ctx context.Context, rowsIn []CommonExportRow) ([]string, [][]string) {
	header := buildGlobalCommonHeader(ctx)
	rows := make([][]string, 0, len(rowsIn))
	for _, r := range rowsIn {
		rows = append(rows, []string{
			r.ProjectName,
			r.WorkflowType,
			r.WorkflowID,
			r.WorkflowName,
			r.CreatorName,
			r.CreatedAt,
			r.Status,
			r.DataSource,
		})
	}
	return header, rows
}

// LocalizeWorkflowTypeSQLRelease returns localized "SQL上线工单".
func LocalizeWorkflowTypeSQLRelease(ctx context.Context) string {
	return locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportTypeSQLRelease)
}

// LocalizeWorkflowTypeDataExport returns localized "数据导出工单".
func LocalizeWorkflowTypeDataExport(ctx context.Context) string {
	return locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportTypeDataExport)
}

// LocalizeUnifiedStatus exposes unified status localization for callers assembling common rows.
func LocalizeUnifiedStatus(ctx context.Context, status string) string {
	return localizeUnifiedStatus(ctx, status)
}

// ExceedLimitError returns a user-readable business error when workflow count exceeds the limit.
func ExceedLimitError(count uint64) error {
	return errors.NewDataInvalidErr("导出工单数量超过上限 %d（当前 %d），请收窄筛选条件后重试", MaxExportWorkflowCount, count)
}

func buildProjectHeader(ctx context.Context) []string {
	header := []string{
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowNumber),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowName),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowDescription),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataSource),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportCreateTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportCreator),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportTaskOrderStatus),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportOperator),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportSQLContent),
	}
	for i := 0; i < projectFixedNodes; i++ {
		header = append(header,
			locale.Bundle.LocalizeMsgByCtx(ctx, projectNodeAuditorMsgs[i]),
			locale.Bundle.LocalizeMsgByCtx(ctx, projectNodeAuditTimeMsgs[i]),
			locale.Bundle.LocalizeMsgByCtx(ctx, projectNodeAuditResultMsgs[i]),
		)
	}
	header = append(header,
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutor),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionStartTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionEndTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionStatus),
	)
	return header
}

func buildGlobalSQLReleaseHeader(ctx context.Context, maxN int) []string {
	header := []string{
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportProjectName),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowNumber),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowName),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowDescription),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataSource),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportCreateTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportCreator),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportTaskOrderStatus),
	}
	header = append(header, buildDynamicAuditNodeHeaders(ctx, maxN)...)
	header = append(header,
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutor),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionStartTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionEndTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportExecutionStatus),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportSQLContent),
	)
	return header
}

func buildGlobalDataExportHeader(ctx context.Context, maxN int) []string {
	header := []string{
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportProjectName),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowNumber),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowName),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowDescription),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataSource),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportCreateTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportCreator),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportTaskOrderStatus),
	}
	header = append(header, buildDynamicAuditNodeHeaders(ctx, maxN)...)
	header = append(header,
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataExportExecTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataExportResult),
	)
	return header
}

func buildGlobalCommonHeader(ctx context.Context) []string {
	return []string{
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportProjectName),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowType),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowNumber),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportWorkflowName),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportCreator),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportCreateTime),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportTaskOrderStatus),
		locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportDataSource),
	}
}

func buildDynamicAuditNodeHeaders(ctx context.Context, maxN int) []string {
	header := make([]string, 0, maxN*3)
	for k := 1; k <= maxN; k++ {
		header = append(header,
			fmt.Sprintf(locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportNodeAuditorTpl), k),
			fmt.Sprintf(locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportNodeAuditTimeTpl), k),
			fmt.Sprintf(locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportNodeAuditResultTpl), k),
		)
	}
	return header
}

func buildProjectRows(ctx context.Context, workflowIDs []string, includeProjectName bool, projectNameByUID map[string]string) ([][]string, error) {
	workflows, err := loadWorkflowsForExport(ctx, workflowIDs)
	if err != nil {
		return nil, err
	}
	rows := make([][]string, 0)
	for _, workflow := range workflows {
		projectName := ""
		if includeProjectName && projectNameByUID != nil {
			projectName = projectNameByUID[string(workflow.ProjectId)]
		}
		for _, instanceRecord := range workflow.Record.InstanceRecords {
			instanceName := instanceNameOf(instanceRecord)
			row := []string{
				workflow.WorkflowId,
				workflow.Subject,
				workflow.Desc,
				instanceName,
				workflow.Model.CreatedAt.Format(timeLayout),
				dms.GetUserNameWithDelTag(workflow.CreateUserId),
				locale.Bundle.LocalizeMsgByCtx(ctx, model.WorkflowStatus[workflow.Record.Status]),
				dms.GetUserNameWithDelTag(instanceRecord.ExecutionUserId),
				instanceRecord.Task.TaskExecEndAt(),
				getExecuteSqlList(instanceRecord.Task.ExecuteSQLs),
			}
			row = append(row, getAuditAndExecuteList(ctx, workflow, instanceRecord)...)
			if includeProjectName {
				row = append([]string{projectName}, row...)
			}
			rows = append(rows, row)
		}
	}
	return rows, nil
}

func loadWorkflowsForExport(ctx context.Context, workflowIDs []string) ([]*model.Workflow, error) {
	_ = ctx
	s := model.GetStorage()
	out := make([]*model.Workflow, 0, len(workflowIDs))
	for _, id := range workflowIDs {
		workflow, exist, err := s.GetWorkflowExportById(id)
		if err != nil {
			return nil, err
		}
		if !exist {
			log.NewEntry().Errorf("workflow not exist, id: %s", id)
			continue
		}

		instanceIds := make([]uint64, 0, len(workflow.Record.InstanceRecords))
		for _, item := range workflow.Record.InstanceRecords {
			instanceIds = append(instanceIds, item.InstanceId)
		}

		instances, err := dms.GetInstancesInProjectByIds(context.Background(), string(workflow.ProjectId), instanceIds)
		if err != nil {
			return nil, err
		}
		instanceMap := map[uint64]*model.Instance{}
		for _, instance := range instances {
			instanceMap[instance.ID] = instance
		}
		for i, item := range workflow.Record.InstanceRecords {
			if instance, ok := instanceMap[item.InstanceId]; ok {
				workflow.Record.InstanceRecords[i].Instance = instance
			}
		}
		out = append(out, workflow)
	}
	return out, nil
}

func maxAuditNodesFromWorkflows(workflows []*model.Workflow) int {
	maxN := 0
	for _, w := range workflows {
		if n := len(w.AuditStepList()); n > maxN {
			maxN = n
		}
	}
	return normalizeMaxAuditNodes(maxN)
}

func normalizeMaxAuditNodes(n int) int {
	if n < minAuditNodes {
		return minAuditNodes
	}
	return n
}

func instanceNameOf(instanceRecord *model.WorkflowInstanceRecord) string {
	if instanceRecord.Instance != nil {
		return utils.AddDelTag(nil, instanceRecord.Instance.Name)
	}
	return ""
}

func getAuditAndExecuteList(ctx context.Context, workflow *model.Workflow, instanceRecord *model.WorkflowInstanceRecord) (auditAndExecuteList []string) {
	auditAndExecuteList = append(auditAndExecuteList, getAuditListFixed(ctx, workflow)...)
	auditAndExecuteList = append(auditAndExecuteList,
		dms.GetUserNameWithDelTag(instanceRecord.ExecutionUserId),
		instanceRecord.Task.TaskExecStartAt(),
		instanceRecord.Task.TaskExecEndAt(),
		locale.Bundle.LocalizeMsgByCtx(ctx, executeStateMap[instanceRecord.Task.Status]),
	)
	return auditAndExecuteList
}

func getExecuteSqlList(executeSQLList []*model.ExecuteSQL) string {
	var stringBuilder strings.Builder
	for _, executeSQL := range executeSQLList {
		stringBuilder.WriteString(executeSQL.Content)
		stringBuilder.WriteString("\n")
	}
	return stringBuilder.String()
}

func getAuditListFixed(ctx context.Context, workflow *model.Workflow) (workflowList []string) {
	auditNodeList := make([]string, projectFixedNodes*3)
	stepSize := 3
	for i, step := range workflow.AuditStepList() {
		if i >= projectFixedNodes {
			break
		}
		stepIndex := i * stepSize
		auditNodeList[stepIndex] = dms.GetUserNameWithDelTag(step.OperationUserId)
		auditNodeList[stepIndex+1] = step.OperationTime()
		auditNodeList[stepIndex+2] = locale.Bundle.LocalizeMsgByCtx(ctx, workflowStepStateMap[step.State])
	}
	return auditNodeList
}

func getAuditListN(ctx context.Context, workflow *model.Workflow, maxN int) []string {
	out := make([]string, maxN*3)
	steps := workflow.AuditStepList()
	for i := 0; i < maxN && i < len(steps); i++ {
		step := steps[i]
		base := i * 3
		out[base] = dms.GetUserNameWithDelTag(step.OperationUserId)
		out[base+1] = step.OperationTime()
		out[base+2] = locale.Bundle.LocalizeMsgByCtx(ctx, workflowStepStateMap[step.State])
	}
	return out
}

func emptyAuditNodes(maxN int) []string {
	return make([]string, maxN*3)
}

func localizeUnifiedStatus(ctx context.Context, status string) string {
	if msg, ok := unifiedStatusMsg[status]; ok {
		return locale.Bundle.LocalizeMsgByCtx(ctx, msg)
	}
	return locale.Bundle.LocalizeMsgByCtx(ctx, locale.WFExportUnifiedStatusUnknown)
}

// BuildCommonRowsFromSQLRelease builds common-layout rows from SQL-release workflow IDs (one row per workflow).
func BuildCommonRowsFromSQLRelease(ctx context.Context, workflowIDs []string, projectNameByUID map[string]string) ([]CommonExportRow, error) {
	workflows, err := loadWorkflowsForExport(ctx, workflowIDs)
	if err != nil {
		return nil, err
	}
	typeLabel := LocalizeWorkflowTypeSQLRelease(ctx)
	out := make([]CommonExportRow, 0, len(workflows))
	for _, workflow := range workflows {
		names := make([]string, 0, len(workflow.Record.InstanceRecords))
		for _, ir := range workflow.Record.InstanceRecords {
			if n := instanceNameOf(ir); n != "" {
				names = append(names, n)
			}
		}
		projectName := ""
		if projectNameByUID != nil {
			projectName = projectNameByUID[string(workflow.ProjectId)]
		}
		created := workflow.Model.CreatedAt.Format(timeLayout)
		out = append(out, CommonExportRow{
			ProjectName:  projectName,
			WorkflowType: typeLabel,
			WorkflowID:   workflow.WorkflowId,
			WorkflowName: workflow.Subject,
			CreatorName:  dms.GetUserNameWithDelTag(workflow.CreateUserId),
			CreatedAt:    created,
			Status:       localizeUnifiedStatus(ctx, mapSQLWorkflowStatus(workflow.Record.Status)),
			DataSource:   strings.Join(names, ","),
			SortKey:      created,
		})
	}
	return out, nil
}

// BuildCommonRowsFromDataExport builds common-layout rows from data-export records.
func BuildCommonRowsFromDataExport(ctx context.Context, records []DataExportExportRecord) []CommonExportRow {
	typeLabel := LocalizeWorkflowTypeDataExport(ctx)
	out := make([]CommonExportRow, 0, len(records))
	for _, r := range records {
		out = append(out, CommonExportRow{
			ProjectName:  r.ProjectName,
			WorkflowType: typeLabel,
			WorkflowID:   r.WorkflowID,
			WorkflowName: r.WorkflowName,
			CreatorName:  r.CreatorName,
			CreatedAt:    r.CreatedAt,
			Status:       localizeUnifiedStatus(ctx, r.UnifiedStatus),
			DataSource:   strings.Join(r.DBServiceNames, ","),
			SortKey:      r.CreatedAt,
		})
	}
	return out
}
