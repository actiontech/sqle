package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	imPkg "github.com/actiontech/sqle/sqle/pkg/im"
	"github.com/actiontech/sqle/sqle/pkg/im/dingding"
	"github.com/sirupsen/logrus"
)

type DingTalkJob struct {
	BaseJob
}

func NewDingTalkJob(entry *logrus.Entry) ServerJob {
	entry = entry.WithField("job", "ding_talk")
	j := &DingTalkJob{}
	j.BaseJob = *NewBaseJob(entry, 60*time.Second, j.dingTalkRotation)
	return j
}

func (j *DingTalkJob) dingTalkRotation(entry *logrus.Entry) {
	st := model.GetStorage()

	ims, err := st.GetAllIMConfig()
	if err != nil {
		entry.Errorf("get all im config failed, error: %v", err)
	}

	for _, im := range ims {
		switch im.Type {
		case model.ImTypeDingTalk:
			d := &dingding.DingTalk{
				AppKey:    im.AppKey,
				AppSecret: im.AppSecret,
			}

			dingTalkInstances, err := st.GetDingTalkInstByStatus(model.ApproveStatusInitialized)
			if err != nil {
				entry.Errorf("get ding talk status error: %v", err)
				continue
			}

			for _, dingTalkInstance := range dingTalkInstances {
				approval, err := d.GetApprovalDetail(dingTalkInstance.ApproveInstanceCode)
				if err != nil {
					entry.Errorf("get ding talk approval detail error: %v", err)
					continue
				}

				switch *approval.Result {
				case model.ApproveStatusAgree:
					workflow, err := dms.GetWorkflowDetailByWorkflowId("", dingTalkInstance.WorkflowId, st.GetWorkflowDetailWithoutInstancesByWorkflowID)
					if err != nil {
						entry.Errorf("get workflow detail error: %v", err)
						continue
					}

					if workflow.Record.Status == model.WorkflowStatusCancel {
						entry.Errorf("workflow has canceled skip, id: %s", dingTalkInstance.WorkflowId)
						continue
					}

					nextStep := workflow.NextStep()

					userId := *approval.OperationRecords[1].UserId

					user, err := getUserByUserId(d, userId, workflow.CurrentStep().Assignees)
					if err != nil {
						entry.Errorf("get user by user id error: %v", err)
						continue
					}

					if err := ApproveWorkflowProcess(workflow, user, st); err != nil {
						entry.Errorf("approve workflow process error: %v", err)
						continue
					}

					dingTalkInstance.Status = model.ApproveStatusAgree
					if err := st.Save(&dingTalkInstance); err != nil {
						entry.Errorf("save ding talk instance error: %v", err)
						continue
					}

					if nextStep.Template.Typ != model.WorkflowStepTypeSQLExecute {
						imPkg.CreateApprove(string(workflow.ProjectId), workflow.WorkflowId)
					}

				case model.ApproveStatusRefuse:
					workflow, err := dms.GetWorkflowDetailByWorkflowId("", dingTalkInstance.WorkflowId, st.GetWorkflowDetailWithoutInstancesByWorkflowID)
					if err != nil {
						entry.Errorf("get workflow detail error: %v", err)
						continue
					}

					if workflow.Record.Status == model.WorkflowStatusCancel {
						entry.Errorf("workflow has canceled skip, id: %s", dingTalkInstance.WorkflowId)
						continue
					}

					var reason string
					if approval.OperationRecords[1] != nil && approval.OperationRecords[1].Remark != nil {
						reason = *approval.OperationRecords[1].Remark
					} else {
						reason = "审批拒绝"
					}

					userId := *approval.OperationRecords[1].UserId
					user, err := getUserByUserId(d, userId, workflow.CurrentStep().Assignees)
					if err != nil {
						entry.Errorf("get user by user id error: %v", err)
						continue
					}

					if err := RejectWorkflowProcess(workflow, reason, user, st); err != nil {
						entry.Errorf("reject workflow process error: %v", err)
						continue
					}

					dingTalkInstance.Status = model.ApproveStatusRefuse
					if err := st.Save(&dingTalkInstance); err != nil {
						entry.Errorf("save ding talk instance error: %v", err)
						continue
					}
				default:
					// ding talk rotation, no action
				}
			}
		}
	}
}

func getUserByUserId(d *dingding.DingTalk, userId string, assigneesUsers string) (*model.User, error) {
	userMaps, err := dms.GetMapUsers(context.TODO(), strings.Split(assigneesUsers, ","), dms.GetDMSServerAddress())
	if err != nil {
		return nil, err
	}
	phone, err := d.GetMobileByUserID(userId)
	if err != nil {
		return nil, fmt.Errorf("get user mobile error: %v", err)
	}

	for _, assigneeUser := range userMaps {
		if assigneeUser.Phone == phone {
			return assigneeUser, nil
		}
	}

	return nil, fmt.Errorf("user not found, phone: %s", phone)
}
