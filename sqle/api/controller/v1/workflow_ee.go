//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func exportWorkflowV1(c echo.Context) error {
	req := new(ExportWorkflowReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	up, err := dms.NewUserPermission(user.GetIDStr(), projectUid)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := map[string]interface{}{
		"filter_subject":                       req.FilterSubject,
		"filter_create_time_from":              req.FilterCreateTimeFrom,
		"filter_create_time_to":                req.FilterCreateTimeTo,
		"filter_create_user_id":                req.FilterCreateUserID,
		"filter_task_execute_start_time_from":  req.FilterTaskExecuteStartTimeFrom,
		"filter_task_execute_start_time_to":    req.FilterTaskExecuteStartTimeTo,
		"filter_status":                        req.FilterStatus,
		"filter_current_step_assignee_user_id": req.FilterCurrentStepAssigneeUserId,
		"filter_task_instance_id":              req.FilterTaskInstanceId,
		"filter_project_id":                    projectUid,
		"current_user_id":                      user.ID,
		"check_user_can_access":                !up.IsAdmin(),
	}

	if req.FuzzyKeyword != "" {
		data["fuzzy_keyword"] = fmt.Sprintf("%%%s%%", req.FuzzyKeyword)
	}
	if !up.IsAdmin() {
		data["viewable_instance_ids"] = strings.Join(up.GetInstancesByOP(dmsV1.OpPermissionTypeViewOthersWorkflow), ",")
	}

	idList, err := s.GetExportWorkflowIDListByReq(data, nil)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	ctx := c.Request().Context()
	buff := new(bytes.Buffer)
	buff.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	csvWriter := csv.NewWriter(buff)
	if err := csvWriter.Write([]string{
		locale.ShouldLocalizeMsg(ctx, locale.WFExportWorkflowNumber),      // "工单编号",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportWorkflowName),        // "工单名称",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportWorkflowDescription), // "工单描述",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportDataSource),          // "数据源",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportCreateTime),          // "创建时间",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportCreator),             // "创建人 ",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportTaskOrderStatus),     // "工单状态",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportOperator),            // "操作人",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportExecutionTime),       // "工单执行时间",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportSQLContent),          // "具体执行SQL内容",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode1Auditor),        // "[节点1]审核人",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode1AuditTime),      // "[节点1]审核时间",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode1AuditResult),    // "[节点1]审核结果",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode2Auditor),        // "[节点2]审核人",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode2AuditTime),      // "[节点2]审核时间",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode2AuditResult),    // "[节点2]审核结果",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode3Auditor),        // "[节点3]审核人",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode3AuditTime),      // "[节点3]审核时间",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode3AuditResult),    // "[节点3]审核结果",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode4Auditor),        // "[节点4]审核人",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode4AuditTime),      // "[节点4]审核时间",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportNode4AuditResult),    // "[节点4]审核结果",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportExecutor),            // "上线人",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportExecutionStartTime),  // "上线开始时间",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportExecutionEndTime),    // "上线结束时间",
		locale.ShouldLocalizeMsg(ctx, locale.WFExportExecutionStatus),     // "上线结果",
	}); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	for _, id := range idList {
		workflow, exist, err := s.GetWorkflowExportById(id)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
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
			return controller.JSONBaseErrorReq(c, err)
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

		var exportWorkflowRecord []string
		for _, instanceRecord := range workflow.Record.InstanceRecords {
			exportWorkflowRecord = []string{
				workflow.WorkflowId,
				workflow.Subject,
				workflow.Desc,
				utils.AddDelTag(nil, instanceRecord.Instance.Name),
				workflow.Model.CreatedAt.Format("2006-01-02 15:04:05"),
				dms.GetUserNameWithDelTag(workflow.CreateUserId),
				locale.ShouldLocalizeMsg(ctx, model.WorkflowStatus[workflow.Record.Status]),
				dms.GetUserNameWithDelTag(instanceRecord.ExecutionUserId),
				instanceRecord.Task.TaskExecEndAt(),
				getExecuteSqlList(instanceRecord.Task.ExecuteSQLs),
			}
			exportWorkflowRecord = append(exportWorkflowRecord, getAuditAndExecuteList(ctx, workflow, instanceRecord)...)

			if err := csvWriter.Write(exportWorkflowRecord); err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
		}
	}

	csvWriter.Flush()

	fileName := fmt.Sprintf("%s_workflow.csv", time.Now().Format("20060102150405"))
	c.Response().Header().Set(echo.HeaderContentDisposition, mime.FormatMediaType("attachment", map[string]string{
		"filename": fileName,
	}))

	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}

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

// 获取审核和上线节点
func getAuditAndExecuteList(ctx context.Context, workflow *model.Workflow, instanceRecord *model.WorkflowInstanceRecord) (auditAndExecuteList []string) {
	// 审核节点
	auditAndExecuteList = append(auditAndExecuteList, getAuditList(ctx, workflow)...)
	// 上线节点
	auditAndExecuteList = append(auditAndExecuteList,
		dms.GetUserNameWithDelTag(instanceRecord.ExecutionUserId),
		instanceRecord.Task.TaskExecStartAt(),
		instanceRecord.Task.TaskExecEndAt(),
		locale.ShouldLocalizeMsg(ctx, executeStateMap[instanceRecord.Task.Status]),
	)
	return auditAndExecuteList
}

// 获取上线sql
func getExecuteSqlList(executeSQLList []*model.ExecuteSQL) string {
	var stringBuilder strings.Builder
	for _, executeSQL := range executeSQLList {
		stringBuilder.WriteString(executeSQL.Content)
		stringBuilder.WriteString("\n")
	}
	return stringBuilder.String()
}

func getAuditList(ctx context.Context, workflow *model.Workflow) (workflowList []string) {
	auditNodeList := make([]string, 12) // 4个审核节点,每个节点有3个字段,最大3*4个字段
	stepSize := 3                       // 每个节点有3个字段
	for i, step := range workflow.AuditStepList() {
		stepIndex := i * stepSize
		auditNodeList[stepIndex] = dms.GetUserNameWithDelTag(step.OperationUserId)
		auditNodeList[stepIndex+1] = step.OperationTime()
		auditNodeList[stepIndex+2] = locale.ShouldLocalizeMsg(ctx, workflowStepStateMap[step.State])
	}
	return auditNodeList
}

func getWorkflowTemplate(c echo.Context) error {
	s := model.GetStorage()

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	var td *model.WorkflowTemplate

	template, exist, err := s.GetWorkflowTemplateByProjectId(model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		td = model.DefaultWorkflowTemplate(projectUid)
		err = s.SaveWorkflowTemplate(td)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		td.Desc = fmt.Sprintf(locale.ShouldLocalizeMsg(c.Request().Context(), locale.DefaultTemplatesDesc), projectUid)
	} else {
		td, err = getWorkflowTemplateDetailByTemplate(template)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	return c.JSON(http.StatusOK, &GetWorkflowTemplateResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    convertWorkflowTemplateToRes(td),
	})
}

func getWorkflowTemplateDetailByTemplate(template *model.WorkflowTemplate) (*model.WorkflowTemplate, error) {
	s := model.GetStorage()
	steps, err := s.GetWorkflowStepsDetailByTemplateId(template.ID)
	if err != nil {
		return nil, err
	}
	template.Steps = steps
	return template, nil
}

func validWorkflowTemplateReq(steps []*WorkFlowStepTemplateReqV1) error {
	if len(steps) == 0 {
		return fmt.Errorf("workflow steps cannot be empty")
	}
	if len(steps) > 5 {
		return fmt.Errorf("workflow steps length must be less than 6")
	}

	for i, step := range steps {
		isLastStep := i == len(steps)-1
		if isLastStep && step.Type != model.WorkflowStepTypeSQLExecute {
			return fmt.Errorf("the last workflow step type must be sql_execute")
		}
		if !isLastStep && step.Type == model.WorkflowStepTypeSQLExecute {
			return fmt.Errorf("workflow step type sql_execute just be used in last step")
		}
		if len(step.Users) == 0 && !step.ApprovedByAuthorized && !step.ExecuteByAuthorized {
			return fmt.Errorf("the assignee is empty for step %s", step.Desc)
		}
		if len(step.Users) > 3 {
			return fmt.Errorf("the assignee for step cannot be more than 3")
		}
	}
	return nil
}

func updateWorkflowTemplate(c echo.Context) error {
	req := new(UpdateWorkflowTemplateReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}

	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()

	workflowTemplate, exist, err := s.GetWorkflowTemplateByProjectId(model.ProjectUID(projectUid))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist,
			fmt.Errorf("workflow template is not exist")))
	}

	if req.Steps != nil {
		err = validWorkflowTemplateReq(req.Steps)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
		}

		// dms-todo: 校验step.Users用户是否存在

		steps := make([]*model.WorkflowStepTemplate, 0, len(req.Steps))
		for i, step := range req.Steps {
			s := &model.WorkflowStepTemplate{
				Number: uint(i + 1),
				ApprovedByAuthorized: sql.NullBool{
					Bool:  step.ApprovedByAuthorized,
					Valid: true,
				},
				ExecuteByAuthorized: sql.NullBool{
					Bool:  step.ExecuteByAuthorized,
					Valid: true,
				},
				Typ:  step.Type,
				Desc: step.Desc,
			}
			s.Users = strings.Join(step.Users, ",")
			steps = append(steps, s)
		}
		err = s.UpdateWorkflowTemplateSteps(workflowTemplate.ID, steps)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
	}

	if req.Desc != nil {
		workflowTemplate.Desc = *req.Desc
	}

	if req.AllowSubmitWhenLessAuditLevel != nil {
		workflowTemplate.AllowSubmitWhenLessAuditLevel = *req.AllowSubmitWhenLessAuditLevel
	}

	err = s.Save(workflowTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func updateSqlFileOrderByWorkflow(c echo.Context) error {
	req := new(UpdateSqlFileOrderV1Req)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	taskIdStr := c.Param("task_id")

	isCan, err := checkTaskCanExecAndUserHasPermission(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !isCan {
		return controller.JSONBaseErrorReq(c,
			fmt.Errorf("task has not reached the execution step or you are not allow to adjust the order, task id:%s", taskIdStr))
	}

	originSortedFiles, err := getOriginSortedFiles(taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = checkParamsIsValid(req, originSortedFiles)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	newOrderFiles := reorderFiles(req.FilesToSort, originSortedFiles)

	s := model.GetStorage()
	if err := s.BatchSaveFileRecords(newOrderFiles); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func checkTaskCanExecAndUserHasPermission(c echo.Context) (bool, error) {
	s := model.GetStorage()
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), c.Param("project_name"), true)
	taskIdStr := c.Param("task_id")
	workflowID := c.Param("workflow_id")
	taskID, err := strconv.Atoi(taskIdStr)
	if err != nil {
		return false, err
	}

	user, err := controller.GetCurrentUser(c, dms.GetUser)
	if err != nil {
		return false, err
	}

	var workflow *model.Workflow
	{
		workflow, err = dms.GetWorkflowDetailByWorkflowId(projectUid, workflowID, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
		if err != nil {
			return false, err
		}
	}
	err = PrepareForTaskExecution(c, string(workflow.ProjectId), workflow, user, taskID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getOriginSortedFiles(taskId string) ([]*model.AuditFile, error) {
	s := model.GetStorage()
	sortedFiles := []*model.AuditFile{}
	files, err := s.GetFileByTaskId(taskId)
	if err != nil {
		return sortedFiles, err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ExecOrder < files[j].ExecOrder
	})
	return files, nil
}

/*
参数校验内容：
1. 新索引是否越界
2. 传参中的文件ID是否存在
3. 传参中的文件ID和索引是否存在重复传递
*/
func checkParamsIsValid(req *UpdateSqlFileOrderV1Req, originFiles []*model.AuditFile) error {
	originFileIdSet := make(map[uint]struct{})
	for _, file := range originFiles {
		originFileIdSet[file.ID] = struct{}{}
	}

	filesLength := len(originFiles)
	paramFileIdSet := make(map[uint]struct{})
	paramNewIndexSet := make(map[uint]struct{})
	for _, fileToSort := range req.FilesToSort {
		newIndex := fileToSort.NewIndex
		fileId := fileToSort.FileID
		if int(newIndex) >= filesLength {
			return fmt.Errorf("new index setting is too long, file id:%d, new index:%d, length of files:%d", fileId, newIndex, filesLength)
		}
		if newIndex <= 0 {
			return fmt.Errorf("new index must be greater than 0, file id:%d, new index:%d", fileId, newIndex)
		}

		_, exist := originFileIdSet[fileId]
		if !exist {
			return fmt.Errorf("file id is not exist, file id:%d", fileId)
		}

		_, exist = paramFileIdSet[fileId]
		if exist {
			return fmt.Errorf("duplicate file IDs found, file id:%d", fileId)
		} else {
			paramFileIdSet[fileId] = struct{}{}
		}
		_, exist = paramNewIndexSet[newIndex]
		if exist {
			return fmt.Errorf("duplicate indexes found, new index:%d", newIndex)
		} else {
			paramNewIndexSet[newIndex] = struct{}{}
		}
	}
	return nil
}

// originSortedFiles索引0的位置是一个不允许排序的文件，代表的是zip包
// 重新排序时，不允许传递newindex为0的数据
func reorderFiles(filesToSort []FileToSort, originSortedFiles []*model.AuditFile) []*model.AuditFile {
	auditFileById := make(map[uint]*model.AuditFile)
	for _, file := range originSortedFiles {
		auditFileById[file.ID] = file
	}

	newOrderFiles := make([]*model.AuditFile, len(originSortedFiles))
	for _, fileToSort := range filesToSort {
		auditFile := auditFileById[fileToSort.FileID]
		newOrderFiles[fileToSort.NewIndex] = auditFile
	}

	unadjustedFiles := removeFileFromOriginFiles(filesToSort, originSortedFiles)

	fillUnadjustedFiles(newOrderFiles, unadjustedFiles)

	return newOrderFiles
}

func removeFileFromOriginFiles(fileNewIndexes []FileToSort, originFiles []*model.AuditFile) []*model.AuditFile {
	files := []*model.AuditFile{}

	newIndexFileIds := make(map[uint]struct{})
	for _, file := range fileNewIndexes {
		newIndexFileIds[file.FileID] = struct{}{}
	}

	for _, file := range originFiles {
		if _, exist := newIndexFileIds[file.ID]; !exist {
			files = append(files, file)
		}
	}
	return files
}

func fillUnadjustedFiles(orderedFiles, unadjustedFiles []*model.AuditFile) {
	index := 0
	for i := 0; i < len(orderedFiles); i++ {
		if orderedFiles[i] == nil {
			orderedFiles[i] = unadjustedFiles[index]
			index++
		}
	}

	for i, file := range orderedFiles {
		file.ExecOrder = uint(i)
	}
}
