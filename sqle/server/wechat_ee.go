//go:build enterprise
// +build enterprise

package server

import (
	"context"
	"time"

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
			record.OaResult = model.ApproveStatusAgree
			if err := s.Save(&record); err != nil {
				entry.Errorf("save wechat record error: %v", err)
				continue
			}
		case model.REJECTED:
			if err := s.WechatCancelScheduledTask(record); err != nil {
				entry.Errorf("save wechat record error: %v", err)
				continue
			}
		}
	}
}

func sendWechatScheduledApprove(entry *logrus.Entry) error {
	s := model.GetStorage()

	workflowWithRecords, err := getWorkflowWithScheduledRecords(entry)
	if err != nil {
		return err
	}
	for _, wfWithRecords := range workflowWithRecords {
		w := wfWithRecords.Workflow

		for _, record := range wfWithRecords.NeedScheduledRecords {
			if !record.NeedScheduledTaskNotify {
				continue
			}
			wechatScheduledRecord, err := s.GetWechatRecordByTaskId(record.TaskId)
			if err != nil {
				entry.Errorf("get wechat scheduled record error: %v", err)
			}
			// 审批工单编号为空，代表未发送审批，需要发送oa审批
			if wechatScheduledRecord.SpNo == "" {
				im.CreateScheduledApprove(record.TaskId, string(w.ProjectId), w.WorkflowId)
			}
		}
	}
	return nil
}
