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
	"github.com/actiontech/sqle/sqle/pkg/im/wechat"
	"github.com/sirupsen/logrus"
)

type WechatJob struct {
	BaseJob
}

func NewWechatJob(entry *logrus.Entry) ServerJob {
	w := new(WechatJob)
	w.BaseJob = *NewBaseJob(entry, 60*time.Second, w.wechatRotation)
	return w
}

func (w *WechatJob) wechatRotation(entry *logrus.Entry) {
	if err := sendWechatScheduledApprove(entry); err != nil {
		entry.Errorf("send wechat scheduled approve error: %v", err)
	}

	s := model.GetStorage()
	im, exist, err := s.GetImConfigByType(model.ImTypeWechatAudit)
	if err != nil {
		entry.Errorf("get wechat config by type error: %v", err)
		return
	}
	if !exist {
		return
	}

	if !im.IsEnable {
		entry.Infof("wechat config is disabled")
		return
	}
	records, err := s.GetWechatRecordByStatus(model.ApproveStatusInitialized)
	if err != nil {
		entry.Errorf("get wechat record by status error: %v", err)
		return
	}
	client := wechat.NewWechatClient(im.AppKey, im.AppSecret)
	for _, record := range records {
		instDetail, err := client.GetApprovalRecordDetail(context.TODO(), record.SpNo)
		if err != nil {
			entry.Errorf("get wechat approval record detail error: %v", err)
			continue
		}

		switch model.WechatOAStatus(instDetail.SpStatus) {
		case model.APPROVED:
			if err := s.WechatAgreeScheduledTask(record); err != nil {
				entry.Errorf("save wechat record error: %v", err)
				continue
			}
		case model.REJECTED:
			if err := s.WechatCancelScheduledTask(record); err != nil {
				entry.Errorf("save wechat record error: %v", err)
				continue
			}
			entry.Warnf("cancel scheduled task, workflow id:%v, instance id:%v", record.TaskId, record.Task.InstanceId)
		}
	}
}

func sendWechatScheduledApprove(entry *logrus.Entry) error {
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
		records, err := st.GetWechatRecordsByTaskIds(taskIds)
		if err != nil {
			return fmt.Errorf("get wechat record failed, taskIDs:%v, err:%v", taskIds, err)
		}

		needSendOATaskIds := []uint{}
		for _, r := range records {
			if r.SpNo == "" {
				needSendOATaskIds = append(needSendOATaskIds, r.TaskId)
			}
		}

		for _, taskId := range needSendOATaskIds {
			im.CreateScheduledApprove(taskId, string(w.ProjectId), w.WorkflowId, model.ImTypeWechatAudit)
		}
	}

	return nil
}
