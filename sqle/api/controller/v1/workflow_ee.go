//go:build enterprise
// +build enterprise

package v1

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func exportWorkflowV1(c echo.Context) error {
	req := new(ExportWorkflowReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	projectName := c.Param("project_name")

	s := model.GetStorage()
	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrProjectNotExist(projectName))
	}

	user, err := controller.GetCurrentUser(c)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	if err := CheckIsProjectMember(user.Name, project.Name); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	data := map[string]interface{}{
		"filter_subject":                         req.FilterSubject,
		"filter_create_time_from":                req.FilterCreateTimeFrom,
		"filter_create_time_to":                  req.FilterCreateTimeTo,
		"filter_create_user_name":                req.FilterCreateUserName,
		"filter_task_execute_start_time_from":    req.FilterTaskExecuteStartTimeFrom,
		"filter_task_execute_start_time_to":      req.FilterTaskExecuteStartTimeTo,
		"filter_status":                          req.FilterStatus,
		"filter_current_step_assignee_user_name": req.FilterCurrentStepAssigneeUserName,
		"filter_task_instance_name":              req.FilterTaskInstanceName,
		"filter_project_name":                    project.Name,
		"current_user_id":                        user.ID,
		"check_user_can_access":                  CheckIsProjectManager(user.Name, project.Name) != nil,
	}

	idList, err := s.GetExportWorkflowIDListByReq(data, user)
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
		workflow, exist, err := s.GetWorkflowDetailById(id)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			log.NewEntry().Errorf("workflow not exist, id: %s", id)
			continue
		}

		var exportWorkflowRecord []string
		for _, instanceRecord := range workflow.Record.InstanceRecords {
			exportWorkflowRecord = []string{
				workflow.WorkflowId,
				workflow.Subject,
				workflow.Desc,
				instanceRecord.Instance.Name,
				workflow.Model.CreatedAt.Format("2006-01-02 15:04:05"),
				workflow.CreateUser.Name,
				model.WorkflowStatus[workflow.Record.Status],
				getUserNameList(workflow.CurrentAssigneeUser()),
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
		getUserNameList(workflow.FinalStep().Assignees),
		instanceRecord.Task.TaskExecStartAt(),
		instanceRecord.Task.TaskExecEndAt(),
		executeStateMap[instanceRecord.Task.Status],
	)
	return auditAndExecuteList
}

// 获取待操作人
func getUserNameList(users []*model.User) string {
	var names []string
	for _, user := range users {
		names = append(names, user.Name)
	}
	return strings.Join(names, ",")
}

func getAuditList(workflow *model.Workflow) (workflowList []string) {
	auditNodeList := make([]string, 12) // 4个审核节点,每个节点有3个字段,最大3*4个字段
	stepSize := 3                       // 每个节点有3个字段
	for i, step := range workflow.AuditStepList() {
		i2 := i * stepSize
		auditNodeList[i2] = getUserNameList(step.Assignees)
		auditNodeList[i2+1] = step.OperationTime()
		auditNodeList[i2+2] = workflowStepStateMap[step.State]
	}
	return auditNodeList
}
