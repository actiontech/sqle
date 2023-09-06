//go:build enterprise
// +build enterprise

package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/model"
	imPkg "github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/pkg/im/feishu"
	larkContact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/sirupsen/logrus"
)

type FeishuJob struct {
	BaseJob
}

func NewFeishuJob(entry *logrus.Entry) ServerJob {
	entry = entry.WithField("job", "feishu")
	f := new(FeishuJob)
	f.BaseJob = *NewBaseJob(entry, 60*time.Second, f.feishuRotation)
	return f
}

func (j *FeishuJob) feishuRotation(entry *logrus.Entry) {
	s := model.GetStorage()
	im, exist, err := s.GetImConfigByType(model.ImTypeFeishuApproval)
	if err != nil {
		entry.Errorf("get im config by type error: %v", err)
		return
	}
	if !exist {
		entry.Errorf("im config not exist")
		return
	}

	if !im.IsEnable {
		entry.Infof("im config is disabled")
		return
	}

	instList, err := s.GetFeishuInstByStatus(model.FeishuApproveStatusInitialized)
	if err != nil {
		entry.Errorf("get feishu instance by status error: %v", err)
		return
	}

	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	for _, inst := range instList {
		instDetail, err := client.GetApprovalInstDetail(context.TODO(), inst.ApproveInstanceCode)
		if err != nil {
			entry.Errorf("get feishu approval instance detail error: %v", err)
			continue
		}

		switch *instDetail.Status {
		case model.FeishuApproveStatusApprove:
			workflow, exist, err := s.GetWorkflowDetailById(strconv.Itoa(int(inst.WorkflowId)))
			if err != nil {
				entry.Errorf("get workflow detail error: %v", err)
				continue
			}
			if !exist {
				entry.Errorf("workflow not exist, id: %d", inst.WorkflowId)
				continue
			}

			nextStep := workflow.NextStep()

			openId := instDetail.TaskList[0].OpenId
			user, err := getSqleUserByFeishuUserID(client, *openId, workflow.CurrentAssigneeUser())
			if err != nil {
				entry.Errorf("get user by user id error: %v", err)
				continue
			}

			if workflow.Record.Status == model.WorkflowStatusWaitForAudit {
				if err := ApproveWorkflowProcess(workflow, user, s); err != nil {
					entry.Errorf("approve workflow process error: %v", err)
					continue
				}
			} else if workflow.Record.Status == model.WorkflowStatusWaitForExecution {
				if err := ExecuteTasksProcess(strconv.Itoa(int(workflow.ID)), workflow.Project.Name, user); err != nil {
					entry.Errorf("execute workflow process error: %v", err)
					continue
				}
			} else {
				entry.Errorf("workflow status error, status: %s", workflow.Record.Status)
				continue
			}

			inst.Status = model.FeishuApproveStatusApprove
			if err := s.Save(&inst); err != nil {
				entry.Errorf("save feishu instance error: %v", err)
				continue
			}

			if nextStep != nil {
				imPkg.CreateApprove(strconv.Itoa(int(workflow.ID)))
			}
		case model.FeishuApproveStatusRejected:
			workflow, exist, err := s.GetWorkflowDetailById(strconv.Itoa(int(inst.WorkflowId)))
			if err != nil {
				entry.Errorf("get workflow detail error: %v", err)
				continue
			}
			if !exist {
				entry.Errorf("workflow not exist, id: %d", inst.WorkflowId)
				continue
			}

			var reason string
			timeline := instDetail.Timeline
			if timeline != nil && len(timeline) >= 2 && timeline[1].Comment != nil {
				reason = *timeline[1].Comment
			} else {
				reason = "审批拒绝"
			}

			openId := instDetail.TaskList[0].OpenId
			user, err := getSqleUserByFeishuUserID(client, *openId, workflow.CurrentAssigneeUser())
			if err != nil {
				entry.Errorf("get user by user id error: %v", err)
				continue
			}

			if err := RejectWorkflowProcess(workflow, reason, user, s); err != nil {
				entry.Errorf("reject workflow process error: %v", err)
				continue
			}

			inst.Status = model.FeishuApproveStatusRejected
			if err := s.Save(&inst); err != nil {
				entry.Errorf("save feishu instance error: %v", err)
				continue
			}
		}
	}
}

func getSqleUserByFeishuUserID(client *feishu.FeishuClient, userId string, assignees []*model.User) (*model.User, error) {
	userInfo, err := client.GetFeishuUserInfo(userId, larkContact.UserIdTypeOpenId)
	if err != nil {
		return nil, err
	}

	for _, assignee := range assignees {
		emailEqual := userInfo.Email != nil && assignee.Email == *userInfo.Email
		// 移除前三位区号
		mobileEqual := userInfo.Mobile != nil && assignee.Phone == (*userInfo.Mobile)[3:]
		if emailEqual || mobileEqual {
			return assignee, nil
		}
	}

	return nil, fmt.Errorf("user not found, userId: %s", userId)
}
