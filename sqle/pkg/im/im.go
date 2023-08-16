package im

import (
	"fmt"
	"strings"

	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"

	"golang.org/x/net/context"
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

	// TODO 工单无法直接关联用户，需要保留校验？
	// if workflow.CreateUser.Phone == "" {
	// 	newLog.Error("create user phone is empty")
	// 	return
	// }

	if len(workflow.Record.Steps) == 1 || workflow.CurrentStep() == workflow.Record.Steps[len(workflow.Record.Steps)-1] {
		newLog.Infof("workflow %v only has one approve step or has been approved, no need to create approve instance", workflow.ID)
		return
	}

	assignUsers := workflow.CurrentAssigneeUser()

	ims, err := s.GetAllIMConfig()
	if err != nil {
		newLog.Errorf("get im config error: %v", err)
		return
	}

	workflowCreateUser, err := dmsobject.GetUser(context.TODO(), workflow.CreateUserId, "")
	if err != nil {
		newLog.Errorf("get user error: %v", err)
		return
	}
	for _, im := range ims {
		switch im.Type {
		case model.ImTypeDingTalk:
			if !im.IsEnable {
				continue
			}

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

			createUserId, err := dingTalk.GetUserIDByPhone(workflowCreateUser.Phone)
			if err != nil {
				newLog.Errorf("get origin user id by phone error: %v", err)
				continue
			}

			var userIds []*string
			for _, assignUser := range assignUsers {
				// TODO 使用DMS提供的批量获取用户接口
				user, err := dmsobject.GetUser(context.TODO(), assignUser, controller.GetDMSServerAddress())
				if err != nil {
					newLog.Errorf("get user error: %v", err)
					return
				}
				if user.Phone == "" {
					newLog.Infof("user %v phone is empty, skip", assignUser)
					continue
				}

				userId, err := dingTalk.GetUserIDByPhone(user.Phone)
				if err != nil {
					newLog.Errorf("get user id by phone error: %v", err)
					continue
				}

				userIds = append(userIds, userId)
			}

			systemVariables, err := s.GetAllSystemVariables()
			if err != nil {
				newLog.Errorf("get sqle url system variables error: %v", err)
				continue
			}

			sqleUrl := systemVariables[model.SystemVariableSqleUrl].Value
			// dms-todo: 从 dms 获取 project name
			// workflowUrl := fmt.Sprintf("%v/project/%s/order/%s", sqleUrl, workflow.Project.Name, workflow.WorkflowId)
			// if sqleUrl == "" {
			// 	newLog.Errorf("sqle url is empty")
			// 	workflowUrl = ""
			// }

			// if err := dingTalk.CreateApprovalInstance(workflow.Subject, workflow.ID, createUserId, userIds, auditResult, workflow.Project.Name, workflow.Desc, workflowUrl); err != nil {
			// 	newLog.Errorf("create dingtalk approval instance error: %v", err)
			// 	continue
			// }

			workflowUrl := fmt.Sprintf("%v/project/%s/order/%s", sqleUrl, "todo", workflow.WorkflowId)
			if sqleUrl == "" {
				newLog.Errorf("sqle url is empty")
				workflowUrl = ""
			}

			if err := dingTalk.CreateApprovalInstance(workflow.Subject, workflow.ID, createUserId, userIds, auditResult, "todo", workflow.Desc, workflowUrl); err != nil {
				newLog.Errorf("create dingtalk approval instance error: %v", err)
				continue
			}
		default:
			newLog.Errorf("im type %s not found", im.Type)
		}
	}
}

func UpdateApprove(workflowId uint, phone, status, reason string) {
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
			if !im.IsEnable {
				continue
			}

			dingTalk := &dingding.DingTalk{
				AppKey:    im.AppKey,
				AppSecret: im.AppSecret,
			}

			userID, err := dingTalk.GetUserIDByPhone(phone)
			if err != nil {
				newLog.Errorf("get user id by phone error: %v", err)
				continue
			}

			if err := dingTalk.UpdateApprovalStatus(workflowId, status, *userID, reason); err != nil {
				newLog.Errorf("update approval status error: %v", err)
				continue
			}
		}
	}
}

func CancelApprove(workflowID uint) {
	newLog := log.NewEntry()
	s := model.GetStorage()
	dingTalkInst, exist, err := s.GetDingTalkInstanceByWorkflowID(workflowID)
	if err != nil {
		newLog.Errorf("get dingtalk instance by workflow step id error: %v", err)
		return
	}
	if !exist {
		newLog.Infof("dingtalk instance not exist, workflow id: %v", workflowID)
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
