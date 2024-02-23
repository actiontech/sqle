//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"
	"database/sql"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/labstack/echo/v4"
	"github.com/actiontech/sqle/sqle/errors"
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
		"filter_task_instance_name":            req.FilterTaskInstanceName,
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

	buff := new(bytes.Buffer)
	buff.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	csvWriter := csv.NewWriter(buff)
	if err := csvWriter.Write([]string{
		"工单编号",
		"工单名称",
		"工单描述",
		"数据源",
		"创建时间",
		"创建人 ",
		"工单状态",
		"操作人",
		"工单执行时间",
		"具体执行SQL内容",
		"[节点1]审核人",
		"[节点1]审核时间",
		"[节点1]审核结果",
		"[节点2]审核人",
		"[节点2]审核时间",
		"[节点2]审核结果",
		"[节点3]审核人",
		"[节点3]审核时间",
		"[节点3]审核结果",
		"[节点4]审核人",
		"[节点4]审核时间",
		"[节点4]审核结果",
		"上线人",
		"上线开始时间",
		"上线结束时间",
		"上线结果",
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
				model.WorkflowStatus[workflow.Record.Status],
				dms.GetUserNameWithDelTag(instanceRecord.ExecutionUserId),
				instanceRecord.Task.TaskExecEndAt(),
				getExecuteSqlList(instanceRecord.Task.ExecuteSQLs),
			}
			exportWorkflowRecord = append(exportWorkflowRecord, getAuditAndExecuteList(workflow, instanceRecord)...)

			if err := csvWriter.Write(exportWorkflowRecord); err != nil {
				return controller.JSONBaseErrorReq(c, err)
			}
		}
	}

	csvWriter.Flush()

	fileName := fmt.Sprintf("%s_工单.csv", time.Now().Format("20060102150405"))
	c.Response().Header().Set(echo.HeaderContentDisposition, mime.FormatMediaType("attachment", map[string]string{
		"filename": fileName,
	}))

	return c.Blob(http.StatusOK, "text/csv", buff.Bytes())
}

var workflowStepStateMap = map[string]string{
	model.WorkflowStepStateApprove: "通过",
	model.WorkflowStepStateReject:  "驳回",
}

var executeStateMap = map[string]string{
	model.TaskStatusExecuting:        "正在上线",
	model.TaskStatusExecuteSucceeded: "上线成功",
	model.TaskStatusExecuteFailed:    "上线失败",
	model.TaskStatusManuallyExecuted: "手动上线",
}

// 获取审核和上线节点
func getAuditAndExecuteList(workflow *model.Workflow, instanceRecord *model.WorkflowInstanceRecord) (auditAndExecuteList []string) {
	// 审核节点
	auditAndExecuteList = append(auditAndExecuteList, getAuditList(workflow)...)
	// 上线节点
	auditAndExecuteList = append(auditAndExecuteList,
		dms.GetUserNameWithDelTag(instanceRecord.ExecutionUserId),
		instanceRecord.Task.TaskExecStartAt(),
		instanceRecord.Task.TaskExecEndAt(),
		executeStateMap[instanceRecord.Task.Status],
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

func getAuditList(workflow *model.Workflow) (workflowList []string) {
	auditNodeList := make([]string, 12) // 4个审核节点,每个节点有3个字段,最大3*4个字段
	stepSize := 3                       // 每个节点有3个字段
	for i, step := range workflow.AuditStepList() {
		stepIndex := i * stepSize
		auditNodeList[stepIndex] = dms.GetUserNameWithDelTag(step.OperationUserId)
		auditNodeList[stepIndex+1] = step.OperationTime()
		auditNodeList[stepIndex+2] = workflowStepStateMap[step.State]
	}
	return auditNodeList
}

func addDelUserTag(user *model.User) string {
	if user != nil {
		return utils.AddDelTag(user.DeletedAt, user.Name)
	}
	return ""
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
