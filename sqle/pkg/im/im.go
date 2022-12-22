package im

import (
	"bytes"
	"fmt"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"
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
			Id:        im.ID,
			AppKey:    im.AppKey,
			AppSecret: im.AppSecret,
		}

		if err := dingTalk.CreateApprovalTemplate(); err != nil {
			log.NewEntry().Errorf("create approval template error: %v", err)
			return
		}
	}
}

func CreateApprove(id string) {
	newLog := log.NewEntry()
	s := model.GetStorage()
	workflow, exist, err := s.GetWorkflowDetailById(id)
	if err != nil {
		newLog.Error("get workflow detail error: ", err)
		return
	}
	if !exist {
		newLog.Error("workflow not exist")
		return
	}

	var buff bytes.Buffer
	for _, record := range workflow.Record.InstanceRecords {
		buff.WriteString(fmt.Sprintf("数据源%s审核评分%d;", record.Instance.Name, record.Task.Score))
	}

	if workflow.CreateUser.Phone == "" {
		newLog.Error("create user phone is empty")
		return
	}

	if len(workflow.Record.Steps) == 1 || workflow.CurrentStep() == workflow.Record.Steps[len(workflow.Record.Steps)-1] {
		newLog.Infof("workflow %v only has one approve step or has been approved, no need to create approve instance", workflow.ID)
		return
	}

	users := workflow.CurrentAssigneeUser()

	ims, err := s.GetAllIMConfig()
	if err != nil {
		newLog.Errorf("get im config error: %v", err)
		return
	}

	for _, im := range ims {
		switch im.Type {
		case model.ImTypeDingTalk:
			if !im.IsEnable {
				continue
			}

			dingTalk := &dingding.DingTalk{
				Id:          im.ID,
				AppKey:      im.AppKey,
				AppSecret:   im.AppSecret,
				ProcessCode: im.ProcessCode,
			}

			createUserId, err := dingTalk.GetUserIDByPhone(workflow.CreateUser.Phone)
			if err != nil {
				newLog.Errorf("get origin user id by phone error: %v", err)
				return
			}

			var userIds []*string
			for _, user := range users {
				if user.Phone == "" {
					newLog.Infof("user %v phone is empty, skip", user.ID)
					continue
				}

				userId, err := dingTalk.GetUserIDByPhone(user.Phone)
				if err != nil {
					newLog.Errorf("get user id by phone error: %v", err)
					continue
				}

				userIds = append(userIds, userId)
			}

			if err := dingTalk.CreateApprovalInstance(workflow.Subject, workflow.ID, workflow.CurrentStep().ID, createUserId, userIds, buff.String()); err != nil {
				newLog.Errorf("create dingtalk approval instance error: %v", err)
				return
			}
		default:
			newLog.Errorf("im type %s not found", im.Type)
		}
	}
}

func UpdateApprove(workflowId, stepId uint, phone, status, reason string) {
	newLog := log.NewEntry()
	s := model.GetStorage()

	ims, err := s.GetAllIMConfig()
	if err != nil {
		newLog.Errorf("get im config error: %v", err)
		return
	}

	for _, im := range ims {
		switch im.Type {
		case model.ImTypeDingTalk:
			dingTalk := &dingding.DingTalk{
				AppKey:    im.AppKey,
				AppSecret: im.AppSecret,
			}

			userID, err := dingTalk.GetUserIDByPhone(phone)
			if err != nil {
				newLog.Errorf("get user id by phone error: %v", err)
				continue
			}

			if err := dingTalk.UpdateApprovalStatus(workflowId, stepId, status, *userID, reason); err != nil {
				newLog.Errorf("update approval status error: %v", err)
				continue
			}
		}
	}
}

func CancelApprove(workflowID, workflowStepId uint) {
	newLog := log.NewEntry()
	s := model.GetStorage()
	dingTalkInst, exist, err := s.GetDingTalkInstanceByWorkflowStepID(workflowID, workflowStepId)
	if err != nil {
		newLog.Errorf("get dingtalk instance by workflow step id error: %v", err)
		return
	}
	if !exist {
		newLog.Infof("workflow step %v not exist", workflowStepId)
		return
	}

	ims, err := s.GetAllIMConfig()
	if err != nil {
		newLog.Errorf("get im config error: %v", err)
		return
	}

	for _, im := range ims {
		switch im.Type {
		case model.ImTypeDingTalk:
			if !im.IsEnable {
				continue
			}

			dingTalk := &dingding.DingTalk{
				AppKey:    im.AppKey,
				AppSecret: im.AppSecret,
			}

			if err := dingTalk.CancelApprovalInstance(dingTalkInst.ApproveInstanceCode); err != nil {
				newLog.Errorf("cancel dingtalk approval instance error: %v", err)
				return
			}
		default:
			newLog.Errorf("im type %s not found", im.Type)
		}
	}
}
