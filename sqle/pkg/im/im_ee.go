//go:build enterprise
// +build enterprise

package im

import (
	"context"
	"errors"
	e "errors"
	"fmt"
	"strings"

	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"
	"github.com/actiontech/sqle/sqle/pkg/im/feishu"
	"github.com/actiontech/sqle/sqle/pkg/im/wechat"
	larkContact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

type OaType string

const (
	WorkflowAudit        OaType = "SQLE工单审核"
	ScheduledTaskConfirm OaType = "定时上线审核"
)

var (
	approvalTableLayout = "[%v]"
	approvalTableRow    = "[{\"name\":\"数据源\",\"value\":\"%s\"},{\"name\":\"审核得分\",\"value\":\"%v\"},{\"name\":\"审核通过率\",\"value\":\"%v%%\"}]"
)

var FeishuAuditResultLayout = `
      [
        {
          "id": "7",
          "type": "input",
          "value": "%s"
        },
        {
          "id": "8",
          "type": "input",
          "value": "%v"
        },
        {
          "id": "9",
          "type": "input",
          "value": "%v%%"
        }
      ]
`

func CreateFeishuAuditTemplate(ctx context.Context, im model.IM) error {
	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	approvalCode, err := client.CreateApprovalTemplate(ctx)
	if err != nil {
		return err
	}

	s := model.GetStorage()
	if err := s.UpdateImConfigById(im.ID, map[string]interface{}{
		"process_code": *approvalCode,
	}); err != nil {
		return err
	}

	return nil
}

func CreateFeishuAuditInst(ctx context.Context, im model.IM, workflow *model.Workflow, assignUsers []*model.User, url string) error {
	createUser, err := dms.GetUser(ctx, workflow.CreateUserId, dms.GetDMSServerAddress())
	if err != nil {
		return err
	}
	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	originUser, err := client.GetFeishuUserIdList([]*model.User{createUser}, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}
	if len(originUser) == 0 {
		return nil
	}

	assignUserIDs, err := client.GetFeishuUserIdList(assignUsers, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}

	var tableRows []string
	for _, record := range workflow.Record.InstanceRecords {
		tableRow := fmt.Sprintf(FeishuAuditResultLayout, record.Instance.Name, record.Task.Score, record.Task.PassRate*100)
		tableRows = append(tableRows, tableRow)
	}
	auditResult := strings.Join(tableRows, ",")

	approvalInstCode, err := client.CreateApprovalInstance(ctx, im.ProcessCode, workflow.Subject, originUser[0],
		assignUserIDs, auditResult, string(workflow.ProjectId), workflow.Desc, url, string(WorkflowAudit))
	if err != nil {
		return err
	}

	instDetail, err := client.GetApprovalInstDetail(ctx, *approvalInstCode)
	if err != nil {
		return err
	}

	s := model.GetStorage()
	feishuInst := &model.FeishuInstance{
		ApproveInstanceCode: *approvalInstCode,
		WorkflowId:          workflow.WorkflowId,
		TaskID:              *instDetail.TaskList[0].Id,
	}

	if err := s.Save(&feishuInst); err != nil {
		return err
	}

	return nil
}

func CreateFeishuScheduledRecord(ctx context.Context, im model.IM, workflow *model.Workflow, assignUsers []*model.User, url string, taskId uint) error {
	s := model.GetStorage()

	createUser, err := dms.GetUser(ctx, workflow.CreateUserId, dms.GetDMSServerAddress())
	if err != nil {
		return err
	}
	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	originUser, err := client.GetFeishuUserIdList([]*model.User{createUser}, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}
	if len(originUser) == 0 {
		return nil
	}

	assignUserIDs, err := client.GetFeishuUserIdList(assignUsers, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}

	tableRows := []string{}
	for _, record := range workflow.Record.InstanceRecords {
		if record.TaskId == taskId {
			tableRow := fmt.Sprintf(FeishuAuditResultLayout, record.Instance.Name, record.Task.Score, record.Task.PassRate*100)
			tableRows = append(tableRows, tableRow)
		}
	}

	auditResult := strings.Join(tableRows, ",")

	approvalInstCode, err := client.CreateApprovalInstance(ctx, im.ProcessCode, workflow.Subject, originUser[0],
		assignUserIDs, auditResult, string(workflow.ProjectId), workflow.Desc, url, string(ScheduledTaskConfirm))
	if err != nil {
		return err
	}

	if approvalInstCode == nil {
		return errors.New("feishu send scheduled task error")
	}
	err = s.UpdateFeishuScheduledByTaskId(taskId, map[string]interface{}{"approve_instance_code": *approvalInstCode})
	if err != nil {
		return err
	}

	return nil
}

func UpdateFeishuAuditStatus(ctx context.Context, im model.IM, workflowId string, user *model.User, status string, reason string) error {
	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	userId, err := client.GetFeishuUserIdList([]*model.User{user}, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}
	if len(userId) == 0 {
		return fmt.Errorf("user %s has no associated feishu account", user.Name)
	}

	s := model.GetStorage()
	feishuInst, exist, err := s.GetFeishuInstanceByWorkflowID(workflowId)
	if err != nil {
		return err
	}
	if !exist {
		return e.New("feishu instance not found")
	}

	switch status {
	case model.ApproveStatusAgree:
		err = client.ApproveApproval(ctx, im.ProcessCode, feishuInst.ApproveInstanceCode, userId[0], feishuInst.TaskID)
		if err != nil {
			return err
		}
	case model.WorkflowStatusReject:
		err = client.RejectApproval(ctx, im.ProcessCode, feishuInst.ApproveInstanceCode, userId[0], feishuInst.TaskID, reason)
		if err != nil {
			return err
		}
	default:
		return e.New("invalid approve status")
	}

	feishuInst.Status = status
	if err := s.Save(&feishuInst); err != nil {
		return err
	}

	return nil
}

func CancelFeishuAuditInst(ctx context.Context, im model.IM, workflowIDs []string, user *model.User) error {
	s := model.GetStorage()
	err := s.BatchUpdateStatusOfFeishuInstance(workflowIDs, model.WorkflowStatusCancel)
	if err != nil {
		return err
	}

	feishuInstList, err := s.GetFeishuInstanceListByWorkflowIDs(workflowIDs)
	if err != nil {
		return err
	}

	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	userIdList, err := client.GetFeishuUserIdList([]*model.User{user}, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}

	for _, feishuInst := range feishuInstList {
		inst := feishuInst
		if inst.Status != model.FeishuAuditStatusInitialized {
			log.NewEntry().Infof("feishu approval instance %v status is %v, skip cancel", inst.ApproveInstanceCode, inst.Status)
			continue
		}

		go func() {
			err = client.CancelApproval(ctx, im.ProcessCode, inst.ApproveInstanceCode, userIdList[0])
			if err != nil {
				log.NewEntry().Errorf("cancel feishu approval instance %v error: %v", inst.ApproveInstanceCode, err)
			}
		}()
	}

	return nil
}

func CreateDingdingAuditTemplate(ctx context.Context, im model.IM) error {
	dingTalk := &dingding.DingTalk{
		Id:          im.ID,
		AppKey:      im.AppKey,
		AppSecret:   im.AppSecret,
		ProcessCode: im.ProcessCode,
	}

	err := dingTalk.CreateApprovalTemplate()
	return err
}

func CreateDingdingAuditInst(ctx context.Context, im model.IM, workflow *model.Workflow, assignUsers []*model.User, url string) error {
	if len(workflow.Record.Steps) == 1 || workflow.CurrentStep() == workflow.Record.Steps[len(workflow.Record.Steps)-1] {
		return fmt.Errorf("workflow %v is the last step, no need to create approve instance", workflow.WorkflowId)
	}

	// if workflow.CreateUser.Phone == "" {
	// 	newLog.Error("create user phone is empty")
	// 	return
	// }

	var tableRows []string
	for _, record := range workflow.Record.InstanceRecords {
		tableRow := fmt.Sprintf(approvalTableRow, record.Instance.Name, record.Task.Score, record.Task.PassRate*100)
		tableRows = append(tableRows, tableRow)
	}
	tableRowJoins := strings.Join(tableRows, ",")
	auditResult := fmt.Sprintf(approvalTableLayout, tableRowJoins)

	dingTalk := &dingding.DingTalk{
		Id:          im.ID,
		AppKey:      im.AppKey,
		AppSecret:   im.AppSecret,
		ProcessCode: im.ProcessCode,
	}
	workflowCreateUser, err := dmsobject.GetUser(context.TODO(), workflow.CreateUserId, dms.GetDMSServerAddress())
	if err != nil {
		return fmt.Errorf("get user error: %v", err)
	}
	createUserId, err := dingTalk.GetUserIDByPhone(workflowCreateUser.Phone)
	if err != nil {
		return fmt.Errorf("get origin user id by phone error: %v", err)
	}

	var userIds []*string
	for _, assignUser := range assignUsers {
		if assignUser.Phone == "" {
			log.NewEntry().Infof("user %v phone is empty, skip", assignUser)
			continue
		}
		userId, err := dingTalk.GetUserIDByPhone(assignUser.Phone)
		if err != nil {
			log.NewEntry().Errorf("get user id by phone error: %v", err)
			continue
		}
		userIds = append(userIds, userId)
	}

	if err := dingTalk.CreateApprovalInstance(workflow.Subject, workflow.WorkflowId, createUserId, userIds, auditResult, string(workflow.ProjectId), workflow.Desc, url); err != nil {
		return fmt.Errorf("create dingtalk approval instance error: %v", err)
	}
	return nil
}

func UpdateDingdingAuditStatus(ctx context.Context, im model.IM, workflowId string, user *model.User, status string, reason string) error {
	dingTalk := &dingding.DingTalk{
		AppKey:    im.AppKey,
		AppSecret: im.AppSecret,
	}

	userID, err := dingTalk.GetUserIDByPhone(user.Phone)
	if err != nil {
		return fmt.Errorf("get user id by phone error: %v", err)
	}

	if err := dingTalk.UpdateApprovalStatus(workflowId, status, *userID, reason); err != nil {
		return fmt.Errorf("update approval status error: %v", err)
	}
	return nil
}

func CancelDingdingAuditInst(ctx context.Context, im model.IM, workflowIDs []string, user *model.User) error {
	dingTalk := &dingding.DingTalk{
		AppKey:    im.AppKey,
		AppSecret: im.AppSecret,
	}

	// batch update ding_talk_instances'status into canceled
	s := model.GetStorage()
	err := s.BatchUpdateStatusOfDingTalkInstance(workflowIDs, model.ApproveStatusCancel)
	if err != nil {
		return fmt.Errorf("batch update ding_talk_instances'status into canceled, error: %v", err)
	}

	dingTalkInstList, err := s.GetDingTalkInstanceListByWorkflowIDs(workflowIDs)
	if err != nil {
		return fmt.Errorf("get dingtalk dingTalkInst list by workflow id slice error: %v", err)
	}

	for _, dingTalkInst := range dingTalkInstList {
		inst := dingTalkInst
		// 如果在钉钉上已经同意或者拒绝<=>dingtalk instance的status不为initialized
		// 则只修改钉钉工单状态为取消，不调用取消钉钉工单的API
		if inst.Status != model.ApproveStatusInitialized {
			log.NewEntry().Infof("the dingtalk dingTalkInst cannot be canceled if its status is not initialized, workflow id: %v", dingTalkInst.WorkflowId)
			continue
		}

		go func() {
			if err := dingTalk.CancelApprovalInstance(inst.ApproveInstanceCode); err != nil {
				log.NewEntry().Errorf("cancel dingtalk approval instance error: %v,instant id: %v", err, inst.ID)
			}
		}()
	}
	return nil
}

func CreateWechatAuditTemplate(ctx context.Context, im model.IM) error {
	client := wechat.NewWechatClient(im.AppKey, im.AppSecret)

	approvalCode, err := client.CreateApprovalTemplate(ctx)
	if err != nil {
		return err
	}

	s := model.GetStorage()
	if err := s.UpdateImConfigById(im.ID, map[string]interface{}{
		"process_code": *approvalCode,
	}); err != nil {
		return err
	}
	return nil
}

func CreateWechatAuditRecord(ctx context.Context, im model.IM, workflow *model.Workflow, assignUsers []*model.User, url string, taskId uint) error {
	s := model.GetStorage()
	l := log.NewEntry()
	createUser, err := dms.GetUser(ctx, workflow.CreateUserId, dms.GetDMSServerAddress())
	if err != nil {
		return err
	}
	client := wechat.NewWechatClient(im.AppKey, im.AppSecret)
	if err != nil {
		return err
	}
	if createUser.WeChatID == "" {
		return e.New(fmt.Sprintf("workflow creator not hava wecharID, user id:%v", createUser.ID))
	}

	assignUserWechatIDs := []string{}
	for _, user := range assignUsers {
		if user.WeChatID == "" {
			l.Warnf("assign user not hava wecharID, user id:%v", user.ID)
			continue
		}
		assignUserWechatIDs = append(assignUserWechatIDs, user.WeChatID)
	}

	if len(assignUserWechatIDs) == 0 {
		return e.New("assignUsers not hava wechatID")
	}

	insRecord, err := s.GetWorkInstanceRecordByTaskId(fmt.Sprint(taskId))
	if err != nil {
		return err
	}

	spNo, err := client.CreateApprovalInstance(ctx, im.ProcessCode, workflow.Subject, createUser.WeChatID,
		assignUserWechatIDs, string(workflow.ProjectId), url, string(ScheduledTaskConfirm), []*model.WorkflowInstanceRecord{&insRecord})
	if err != nil {
		return err
	}
	err = s.UpdateWechatRecordByTaskId(taskId, map[string]interface{}{"sp_no": spNo})
	if err != nil {
		return err
	}

	return nil
}

func CreateScheduledApprove(taskId uint, projectId, workflowId string, ImType string) {
	newLog := log.NewEntry()
	s := model.GetStorage()

	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectId, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		newLog.Errorf("get workflow error: %v", err)
	}
	assignUserIds := workflow.CurrentAssigneeUser()

	assignUsers, err := dms.GetUsers(context.TODO(), assignUserIds, controller.GetDMSServerAddress())
	if err != nil {
		newLog.Errorf("get user error: %v", err)
		return
	}

	im, exist, err := s.GetImConfigByType(ImType)
	if err != nil {
		newLog.Errorf("get im config error: %v", err)
		return
	}
	if !exist {
		newLog.Errorf("im does not exist, im type: %v", ImType)
		return
	}

	if !im.IsEnable {
		return
	}

	systemVariables, err := s.GetAllSystemVariables()
	if err != nil {
		newLog.Errorf("get sqle url system variables error: %v", err)
		return
	}

	sqleUrl := systemVariables[model.SystemVariableSqleUrl].Value
	workflowUrl := fmt.Sprintf("%v/project/%s/exec-workflow/%s", sqleUrl, workflow.ProjectId, workflow.WorkflowId)
	if sqleUrl == "" {
		newLog.Errorf("sqle url is empty")
		workflowUrl = ""
	}

	switch im.Type {
	case model.ImTypeWechatAudit:
		if err := CreateWechatAuditRecord(context.TODO(), *im, workflow, assignUsers, workflowUrl, taskId); err != nil {
			newLog.Errorf("create wechat scheduled audit error: %v", err)
			return
		}
	case model.ImTypeFeishuAudit:
		if err := CreateFeishuScheduledRecord(context.TODO(), *im, workflow, assignUsers, workflowUrl, taskId); err != nil {
			newLog.Errorf("create feishu scheduled audit error: %v", err)
			return
		}
	}
}
