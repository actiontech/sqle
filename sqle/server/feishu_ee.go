//go:build enterprise
// +build enterprise

package server

import (
	"context"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im"
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
	im, exist, err := s.GetImConfigByType(model.ImTypeFeishuAudit)
	if err != nil {
		entry.Errorf("get im config by type error: %v", err)
		return
	}
	if !exist {
		return
	}

	if !im.IsEnable {
		entry.Infof("im config is disabled")
		return
	}

	// todo：临时将发送和获取定时任务审批流程的接口放在feishuRotation中
	// 后续会定制方案拆分不同任务的轮询或者以并发的形式发送接口
	// https://github.com/actiontech/sqle-ee/issues/1478
	err = sendFeishuScheduledApprove(entry)
	if err != nil {
		entry.Errorf("send feishu scheduled approve error: %v", err)
	}

	err = updateFeishuScheduledTask(entry, im)
	if err != nil {
		entry.Errorf("update feishu scheduled approve error: %v", err)
	}

	instList, err := s.GetFeishuInstByStatus(model.FeishuAuditStatusInitialized)
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
		case model.FeishuAuditStatusApprove:
			workflow, err := dms.GetWorkflowDetailByWorkflowId("", inst.WorkflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
			if err != nil {
				entry.Errorf("get workflow detail error: %v", err)
				continue
			}

			nextStep := workflow.NextStep()

			openId := instDetail.TaskList[0].OpenId
			assigneesUsers, err := dms.GetUsers(context.Background(), workflow.CurrentAssigneeUser(), dms.GetDMSServerAddress())
			if err != nil {
				entry.Errorf("get user by user id error: %v", err)
				continue
			}
			user, err := getSqleUserByFeishuUserID(client, *openId, assigneesUsers)
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
				if _, err := ExecuteTasksProcess(workflow.WorkflowId, string(workflow.ProjectId), user); err != nil {
					entry.Errorf("execute workflow process error: %v", err)
					continue
				}
			} else {
				entry.Errorf("workflow status error, status: %s", workflow.Record.Status)
				continue
			}

			inst.Status = model.FeishuAuditStatusApprove
			if err := s.Save(&inst); err != nil {
				entry.Errorf("save feishu instance error: %v", err)
				continue
			}

			if nextStep != nil {
				imPkg.CreateApprove(string(workflow.ProjectId), workflow.WorkflowId)
			}
		case model.FeishuAuditStatusRejected:
			workflow, err := dms.GetWorkflowDetailByWorkflowId("", inst.WorkflowId, s.GetWorkflowDetailWithoutInstancesByWorkflowID)
			if err != nil {
				entry.Errorf("get workflow detail error: %v", err)
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
			assigneesUsers, err := dms.GetUsers(context.Background(), workflow.CurrentAssigneeUser(), dms.GetDMSServerAddress())
			if err != nil {
				entry.Errorf("get user by user id error: %v", err)
				continue
			}
			user, err := getSqleUserByFeishuUserID(client, *openId, assigneesUsers)
			if err != nil {
				entry.Errorf("get user by user id error: %v", err)
				continue
			}

			if err := RejectWorkflowProcess(workflow, reason, user, s); err != nil {
				entry.Errorf("reject workflow process error: %v", err)
				continue
			}

			inst.Status = model.FeishuAuditStatusRejected
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

func sendFeishuScheduledApprove(entry *logrus.Entry) error {
	st := model.GetStorage()
	workflows, err := st.GetNeedScheduledWorkflows()
	if err != nil {
		return fmt.Errorf("get need scheduled workflows from storage error: %v", err)
	}

	for _, workflow := range workflows {
		w, err := dms.GetWorkflowDetailByWorkflowId(string(workflow.ProjectId), workflow.WorkflowId, st.GetWorkflowDetailWithoutInstancesByWorkflowID)
		if err != nil {
			return fmt.Errorf("get workflow from storage error: %v", err)
		}
		taskIds, err := w.GetNeedSendOATaskIds(entry)
		if err != nil {
			return err
		}
		records, err := st.GetFeishuRecordsByTaskIds(taskIds)
		if err != nil {
			return fmt.Errorf("get feishu record failed, taskIDs:%v, err:%v", taskIds, err)
		}

		needSendOATaskIds := []uint{}
		for _, r := range records {
			if r.ApproveInstanceCode == "" {
				needSendOATaskIds = append(needSendOATaskIds, r.TaskId)
			}
		}

		for _, taskId := range needSendOATaskIds {
			im.CreateScheduledApprove(taskId, string(w.ProjectId), w.WorkflowId, model.ImTypeFeishuAudit)
		}
	}

	return nil
}

func updateFeishuScheduledTask(entry *logrus.Entry, im *model.IM) error {
	s := model.GetStorage()

	records, err := s.GetFeishuScheduledByStatus(model.FeishuAuditStatusInitialized)
	if err != nil {
		return fmt.Errorf("get feishu record by status error: %v", err)
	}
	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	for _, record := range records {
		instDetail, err := client.GetApprovalInstDetail(context.TODO(), record.ApproveInstanceCode)
		if err != nil {
			entry.Errorf("get feishu approval record detail error: %v", err)
			continue
		}

		switch *instDetail.Status {
		case model.FeishuAuditStatusApprove:
			if err := s.FeishuAgreeScheduledTask(record); err != nil {
				entry.Errorf("save feishu record error: %v", err)
				continue
			}
		case model.FeishuAuditStatusRejected:
			if err := s.FeishuCancelScheduledTask(record); err != nil {
				entry.Errorf("save feishu record error: %v", err)
				continue
			}
			entry.Warnf("cancel scheduled task, workflow id:%v, instance id:%v", record.TaskId, record.Task.InstanceId)
		}
	}
	return nil
}
