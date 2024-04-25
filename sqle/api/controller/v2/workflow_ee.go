//go:build enterprise
// +build enterprise

package v2

import "github.com/actiontech/sqle/sqle/model"

const (
	NotifyTypeWechat = "wechat"
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
	case NotifyTypeFeishu:
		record := model.FeishuScheduledRecord{
			TaskId: curTaskRecord.TaskId,
		}
		if err := s.Save(&record); err != nil {
			return nil
		}
	default:
		return nil
	}
	err := s.UpdateWorkflowInstanceRecordById(curTaskRecord.ID, map[string]interface{}{"need_scheduled_task_notify": true})
	if err != nil {
		return err
	}
	return nil
}

func cancelNotify(taskId uint) error {
	s := model.GetStorage()

	// wechat
	{
		records, err := s.GetWechatRecordsByTaskIds([]uint{taskId})
		if err != nil {
			return err
		}
		if len(records) > 0 {
			return s.DeleteWechatRecordByTaskId(taskId)
		}
	}
	// feishu
	{
		records, err := s.GetFeishuRecordsByTaskIds([]uint{taskId})
		if err != nil {
			return err
		}
		if len(records) > 0 {
			return s.DeleteFeishuRecordByTaskId(taskId)
		}
	}
	return nil
}
