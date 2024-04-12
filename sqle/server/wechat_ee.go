//go:build enterprise
// +build enterprise

package server

import (
	"context"
	"time"

	"github.com/actiontech/sqle/sqle/model"
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

// 当前企微轮询只为二次确认工单定时上线功能
// https://github.com/actiontech/sqle-ee/issues/1441
func (w *WechatJob) wechatRotation(entry *logrus.Entry) {
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
