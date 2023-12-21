package im

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"
)

var (
	approvalTableLayout = "[%v]"
	approvalTableRow    = "[{\"name\":\"数据源\",\"value\":\"%s\"},{\"name\":\"审核得分\",\"value\":\"%v\"},{\"name\":\"审核通过率\",\"value\":\"%v%%\"}]"
)

func CreateApprovalTemplate(imType string) {
	s := model.GetStorage()
	ims, err := s.GetAllIMConfig()
	if err != nil {
		log.NewEntry().Errorf("get im config error: %v", err)
		return
	}

	imTypeIM := make(map[string]model.IM)
	for _, im := range ims {
		imTypeIM[im.Type] = im
	}

	var im model.IM
	var ok bool
	if im, ok = imTypeIM[imType]; !ok {
		log.NewEntry().Errorf("im type %s not found", imType)
		return
	}

	switch im.Type {
	case model.ImTypeDingTalk:
		dingTalk := &dingding.DingTalk{
			Id:          im.ID,
			AppKey:      im.AppKey,
			AppSecret:   im.AppSecret,
			ProcessCode: im.ProcessCode,
		}

		if err := dingTalk.CreateApprovalTemplate(); err != nil {
			log.NewEntry().Errorf("create approval template error: %v", err)
			return
		}
	case model.ImTypeFeishuAudit:
		if err := CreateFeishuAuditTemplate(context.TODO(), im); err != nil {
			log.NewEntry().Errorf("create feishu audit template error: %v", err)
			return
		}
	}
}

func CreateApprove(projectId, workflowId string) {
	newLog := log.NewEntry()
	s := model.GetStorage()
	workflow, err := dms.GetWorkflowDetailByWorkflowId(projectId, workflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
	if err != nil {
		newLog.Error("workflow not exist")
		return
	}

	user, err := dms.GetUser(context.TODO(), workflow.CreateUserId, dms.GetDMSServerAddress())
	if err != nil {
		newLog.Errorf("get user phone failed err: %v", err)
		return
	}
	if user.Phone == "" {
		newLog.Error("create user phone is empty")
		return
	}
	if workflow.CurrentStep() == nil {
		newLog.Infof("workflow %v has no current step, no need to create approve instance", workflow.WorkflowId)
	}

	if len(workflow.Record.Steps) == 1 || workflow.CurrentStep() == workflow.Record.Steps[len(workflow.Record.Steps)-1] {
		newLog.Infof("workflow %v only has one approve step or has been approved, no need to create approve instance", workflow.WorkflowId)
		return
	}

	assignUserIds := workflow.CurrentAssigneeUser()

	assignUsers, err := dms.GetUsers(context.TODO(), assignUserIds, controller.GetDMSServerAddress())
	if err != nil {
		newLog.Errorf("get user error: %v", err)
		return
	}

	ims, err := s.GetAllIMConfig()
	if err != nil {
		newLog.Errorf("get im config error: %v", err)
		return
	}

	for _, im := range ims {
		if !im.IsEnable {
			continue
		}

		systemVariables, err := s.GetAllSystemVariables()
		if err != nil {
			newLog.Errorf("get sqle url system variables error: %v", err)
			continue
		}

		sqleUrl := systemVariables[model.SystemVariableSqleUrl].Value
		workflowUrl := fmt.Sprintf("%v/project/%s/order/%s", sqleUrl, workflow.ProjectId, workflow.WorkflowId)
		if sqleUrl == "" {
			newLog.Errorf("sqle url is empty")
			workflowUrl = ""
		}

		switch im.Type {
		case model.ImTypeDingTalk:
			if len(workflow.Record.Steps) == 1 || workflow.CurrentStep() == workflow.Record.Steps[len(workflow.Record.Steps)-1] {
				newLog.Infof("workflow %v is the last step, no need to create approve instance", workflow.WorkflowId)
				return
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
				newLog.Errorf("get user error: %v", err)
				return
			}
			createUserId, err := dingTalk.GetUserIDByPhone(workflowCreateUser.Phone)
			if err != nil {
				newLog.Errorf("get origin user id by phone error: %v", err)
				continue
			}

			var userIds []*string
			for _, assignUser := range assignUsers {
				if user.Phone == "" {
					newLog.Infof("user %v phone is empty, skip", assignUser)
					continue
				}
				userId, err := dingTalk.GetUserIDByPhone(assignUser.Phone)
				if err != nil {
					newLog.Errorf("get user id by phone error: %v", err)
					continue
				}
				userIds = append(userIds, userId)
			}

			if err := dingTalk.CreateApprovalInstance(workflow.Subject, workflow.WorkflowId, createUserId, userIds, auditResult, string(workflow.ProjectId), workflow.Desc, workflowUrl); err != nil {
				newLog.Errorf("create dingtalk approval instance error: %v", err)
				continue
			}
		case model.ImTypeFeishuAudit:
			if err := CreateFeishuAuditInst(context.TODO(), im, workflow, assignUsers, workflowUrl); err != nil {
				newLog.Errorf("create feishu audit instance error: %v", err)
				continue
			}
		default:
			newLog.Errorf("im type %s not found", im.Type)
		}
	}
}

func UpdateApprove(workflowId string, user *model.User, status, reason string) {
	newLog := log.NewEntry()
	s := model.GetStorage()

	ims, err := s.GetAllIMConfig()
	if err != nil {
		newLog.Errorf("get im config error: %v", err)
		return
	}

	for _, im := range ims {
		if !im.IsEnable {
			continue
		}

		switch im.Type {
		case model.ImTypeDingTalk:
			dingTalk := &dingding.DingTalk{
				AppKey:    im.AppKey,
				AppSecret: im.AppSecret,
			}

			userID, err := dingTalk.GetUserIDByPhone(user.Phone)
			if err != nil {
				newLog.Errorf("get user id by phone error: %v", err)
				continue
			}

			if err := dingTalk.UpdateApprovalStatus(workflowId, status, *userID, reason); err != nil {
				newLog.Errorf("update approval status error: %v", err)
				continue
			}
		case model.ImTypeFeishuAudit:
			if err := UpdateFeishuAuditStatus(context.Background(), im, workflowId, user, status, reason); err != nil {
				newLog.Errorf("update feishu audit status error: %v", err)
				continue
			}
		}
	}
}

func BatchCancelApprove(workflowIds []string, user *model.User) {
	newLog := log.NewEntry()
	s := model.GetStorage()
	ims, err := s.GetAllIMConfig()
	if err != nil {
		newLog.Errorf("get im config error: %v", err)
		return
	}

	for _, im := range ims {
		if !im.IsEnable {
			continue
		}

		switch im.Type {
		case model.ImTypeDingTalk:
			dingTalk := &dingding.DingTalk{
				AppKey:    im.AppKey,
				AppSecret: im.AppSecret,
			}

			// batch update ding_talk_instances'status into canceled
			err = s.BatchUpdateStatusOfDingTalkInstance(workflowIds, model.ApproveStatusCancel)
			if err != nil {
				newLog.Errorf("batch update ding_talk_instances'status into canceled, error: %v", err)
				return
			}

			dingTalkInstList, err := s.GetDingTalkInstanceListByWorkflowIDs(workflowIds)
			if err != nil {
				newLog.Errorf("get dingtalk dingTalkInst list by workflow id slice error: %v", err)
				return
			}

			for _, dingTalkInst := range dingTalkInstList {
				inst := dingTalkInst
				// 如果在钉钉上已经同意或者拒绝<=>dingtalk instance的status不为initialized
				// 则只修改钉钉工单状态为取消，不调用取消钉钉工单的API
				if inst.Status != model.ApproveStatusInitialized {
					newLog.Infof("the dingtalk dingTalkInst cannot be canceled if its status is not initialized, workflow id: %v", dingTalkInst.WorkflowId)
					continue
				}

				go func() {
					if err := dingTalk.CancelApprovalInstance(inst.ApproveInstanceCode); err != nil {
						newLog.Errorf("cancel dingtalk approval instance error: %v,instant id: %v", err, inst.ID)
					}
				}()
			}
		case model.ImTypeFeishuAudit:
			err = CancelFeishuAuditInst(context.TODO(), im, workflowIds, user)
			if err != nil {
				newLog.Errorf("cancel feishu audit instance error: %v", err)
				return
			}
		default:
			newLog.Errorf("im type %s not found", im.Type)
		}
	}
}
