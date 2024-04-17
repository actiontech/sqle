//go:build enterprise
// +build enterprise

package v2

import "github.com/actiontech/sqle/sqle/model"

const (
	NotifyTypeWechat = "Wechat"
)

func createNotifyRecord(notifyType string, curTaskRecord *model.WorkflowInstanceRecord) error {
	s := model.GetStorage()
	switch notifyType {
	case NotifyTypeWechat:
		record := model.WechatRecord{
			TaskId: curTaskRecord.TaskId,
		}
		if err := s.Save(&record); err != nil {
			return nil
		}
		err := s.UpdateWorkflowInstanceRecordById(curTaskRecord.ID, map[string]interface{}{"need_scheduled_task_notify": true, "send_oa_im_type": "wechat"})
		if err != nil {
			return err
		}
	}
	return nil
}
