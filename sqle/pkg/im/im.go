package im

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
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
		if err := CreateDingdingAuditTemplate(context.TODO(), im); err != nil {
			log.NewEntry().Errorf("create dingding audit template error: %v", err)
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
			if err := CreateDingdingAuditInst(context.TODO(), im, workflow, assignUsers, workflowUrl); err != nil {
				newLog.Errorf("create dingding audit instance error: %v", err)
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
			if err := UpdateDingdingAuditStatus(context.Background(), im, workflowId, user, status, reason); err != nil {
				newLog.Errorf("update dingding audit status error: %v", err)
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
			err = CancelDingdingAuditInst(context.TODO(), im, workflowIds, user)
			if err != nil {
				newLog.Errorf("cancel dingding audit instance error: %v", err)
				return
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
