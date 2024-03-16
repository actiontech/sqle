package notification

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	v1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
)

const ManuallyAudit = "manually-audit"

type webHookRequestBody struct {
	Event     string           `json:"event"`
	Action    string           `json:"action"`
	Timestamp string           `json:"timestamp"` // time.RFC3339
	Payload   *httpBodyPayload `json:"payload"`
}

type workflowPayload struct {
	ProjectName     string `json:"project_name"`
	ProjectUID      string `json:"project_uid"`
	WorkflowID      string `json:"workflow_id"`
	WorkflowSubject string `json:"workflow_subject"`
	WorkflowStatus  string `json:"workflow_status"`

	ThirdPartyUserInfo string          `json:"third_party_user_info"`
	CurrentStepID      uint            `json:"current_step_info"`
	WorkflowTaskID     uint            `json:"workflow_task_id"`
	InstanceInfo       []InstanceInfo  `json:"instanceInfo"`
	WorkflowDesc       string          `json:"workflow_desc"`
	SqlTypeMap         map[string]bool `json:"sql_type_map"` // DDL DML DQL
}

type InstanceInfo struct {
	Host   string `json:"host"`
	Schema string `json:"schema"`
	Port   string `json:"port"`
	Desc   string `json:"desc"`
}

type httpBodyPayload struct {
	Workflow *workflowPayload `json:"workflow"`
}

// func TestWorkflowConfig() (err error) {
// 	return workflowSendRequest("create",
// 		"test_project", "1658637666259832832", "test_workflow", "wait_for_audit")
// }

func getProjectNameByID(ProjectId string) (string, error) {
	var projectName string
	ret, _, err := dmsobject.ListProjects(context.TODO(), dms.GetDMSServerAddress(), v1.ListProjectReq{
		PageSize:    1,
		PageIndex:   1,
		FilterByUID: ProjectId,
	})
	if err != nil {
		return projectName, err
	}
	if len(ret) > 0 {
		projectName = ret[0].Name
	}
	return projectName, nil
}

func workflowSendRequest(action string, workflow *model.Workflow) (err error) {
	user, err := dms.GetUser(context.TODO(), workflow.CreateUserId, dms.GetDMSServerAddress())
	if err != nil {
		return err
	}
	projectName, err := getProjectNameByID(string(workflow.ProjectId))
	if err != nil {
		return err
	}
	currentStepID := uint(0)
	if workflow.CurrentStep() != nil {
		currentStepID = workflow.CurrentStep().ID
	}
	reqBody := &webHookRequestBody{
		Event:     "workflow",
		Action:    action,
		Timestamp: time.Now().Format(time.RFC3339),
		Payload: &httpBodyPayload{
			Workflow: &workflowPayload{
				ProjectName:        projectName,
				ProjectUID:         string(workflow.ProjectId),
				WorkflowID:         workflow.WorkflowId,
				WorkflowSubject:    workflow.Subject,
				WorkflowStatus:     workflow.Record.Status,
				ThirdPartyUserInfo: user.ThirdPartyUserInfo,
				CurrentStepID:      currentStepID,
				WorkflowDesc:       workflow.Desc,
				SqlTypeMap: map[string]bool{
					driverV2.SQLTypeDDL: false,
					driverV2.SQLTypeDML: false,
					driverV2.SQLTypeDQL: false,
				},
			},
		},
	}
	for _, record := range workflow.Record.InstanceRecords {
		if record.Instance == nil {
			continue
		}
		info := InstanceInfo{
			Host: record.Instance.Host,
			Port: record.Instance.Port,
			Desc: record.Instance.Desc,
		}
		if record.Task == nil {
			reqBody.Payload.Workflow.InstanceInfo = append(reqBody.Payload.Workflow.InstanceInfo, info)
			continue
		}
		info.Schema = record.Task.Schema
		reqBody.Payload.Workflow.WorkflowTaskID = record.Task.ID
		for _, executeSql := range record.Task.ExecuteSQLs {
			if executeSql.SQLType != "" {
				reqBody.Payload.Workflow.SqlTypeMap[executeSql.SQLType] = true
			}
		}
		reqBody.Payload.Workflow.InstanceInfo = append(reqBody.Payload.Workflow.InstanceInfo, info)
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	return dmsobject.WebHookSendMessage(context.TODO(), controller.GetDMSServerAddress(), &v1.WebHookSendMessageReq{
		WebHookMessage: &v1.WebHooksMessage{
			Message:          string(b),
			TriggerEventType: v1.TriggerEventTypeWorkflow,
		},
	})

}

type webHookAuditPlanRequestBody struct {
	Event     string                `json:"event"`
	Action    string                `json:"action"`
	Timestamp string                `json:"timestamp"` // time.RFC3339
	Payload   *AuditPlanBodyPayload `json:"payload"`
}

type AuditPlanBodyPayload struct {
	AuditPlan *AuditPlanPayload `json:"audit_plan"`
}

type AuditPlanPayload struct {
	ProjectId        string  `json:"project_id"`        // 项目id
	ProjectName      string  `json:"project_name"`      // 项目名称
	ReportId         string  `json:"report_id"`         // 扫描报告id
	AuditPlanName    string  `json:"audit_plan_name"`   // 扫描任务名称
	AuditCreateTime  string  `json:"audit_create_time"` // 扫描任务触发时间
	AuditType        string  `json:"audit_type"`        // 扫描任务类型
	InstanceName     string  `json:"instance_name"`     // 数据源名称
	InstanceDatabase string  `json:"instance_database"` // 数据库名称
	Score            int32   `json:"score"`             // 审核得分
	PassRate         float64 `json:"pass_rate"`         // 审核通过率
	AuditLevel       string  `json:"audit_level"`       // 审核结果等级

	SQLEUrl string `json:"sqle_url"` // sqle地址
}

func auditPlanSendRequest(auditPlan *model.AuditPlan, report *model.AuditPlanReportV2, config AuditPlanNotifyConfig) (err error) {
	var s string
	if config.SQLEUrl != nil {
		s = *config.SQLEUrl
	}

	projectName, err := getProjectNameByID(string(auditPlan.ProjectId))
	if err != nil {
		return err
	}

	reqBody := &webHookAuditPlanRequestBody{
		Event:     "auditplan",
		Action:    ManuallyAudit,
		Timestamp: time.Now().Format(time.RFC3339),
		Payload: &AuditPlanBodyPayload{
			AuditPlan: &AuditPlanPayload{
				ProjectId:        string(auditPlan.ProjectId),
				ProjectName:      projectName,
				ReportId:         strconv.Itoa(int(report.ID)),
				AuditPlanName:    auditPlan.Name,
				AuditCreateTime:  auditPlan.CreatedAt.String(),
				AuditType:        auditPlan.Type,
				InstanceName:     auditPlan.InstanceName,
				InstanceDatabase: auditPlan.InstanceDatabase,
				Score:            report.Score,
				PassRate:         report.PassRate,
				AuditLevel:       report.AuditLevel,
				SQLEUrl:          s,
			},
		},
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	return dmsobject.WebHookSendMessage(context.TODO(), controller.GetDMSServerAddress(), &v1.WebHookSendMessageReq{
		WebHookMessage: &v1.WebHooksMessage{
			Message:          string(b),
			TriggerEventType: v1.TriggerEventAuditPlan,
		},
	})
}
