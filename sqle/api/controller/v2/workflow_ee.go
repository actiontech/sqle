//go:build enterprise
// +build enterprise

package v2

import "github.com/actiontech/sqle/sqle/model"

const (
	NotifyTypeWechat = "Wechat"
	NotifyTypeFeishu = "feishu"
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
		err := s.UpdateWorkflowInstanceRecordById(curTaskRecord.ID, map[string]interface{}{"need_scheduled_task_notify": true})
		if err != nil {
			return err
		}
	case NotifyTypeFeishu:
		record := model.FeishuScheduledRecord{
			TaskId: curTaskRecord.TaskId,
		}
		if err := s.Save(&record); err != nil {
			return nil
		}
		err := s.UpdateWorkflowInstanceRecordById(curTaskRecord.ID, map[string]interface{}{"need_scheduled_task_notify": true})
		if err != nil {
			return err
		}
	}
	return nil
}
