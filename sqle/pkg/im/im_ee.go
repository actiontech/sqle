//go:build enterprise
// +build enterprise

package im

import (
	"context"
	e "errors"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im/feishu"
	larkContact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

var FeishuAuditResultLayout = `
      [
        {
          "id": "6",
          "type": "input",
          "value": "%s"
        },
        {
          "id": "7",
          "type": "input",
          "value": "%v"
        },
        {
          "id": "8",
          "type": "input",
          "value": "%v%%"
        }
      ]
`

func CreateFeishuAuditTemplate(ctx context.Context, im model.IM) error {
	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	approvalCode, err := client.CreateApprovalTemplate(ctx)
	if err != nil {
		return err
	}

	s := model.GetStorage()
	if err := s.UpdateImConfigById(im.ID, map[string]interface{}{
		"process_code": *approvalCode,
	}); err != nil {
		return err
	}

	return nil
}

func CreateFeishuAuditInst(ctx context.Context, im model.IM, workflow *model.Workflow, assignUsers []*model.User, url string) error {
	createUser, err := dms.GetUser(ctx, workflow.CreateUserId, dms.GetDMSServerAddress())
	if err != nil {
		return err
	}
	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	originUser, err := client.GetFeishuUserIdList([]*model.User{createUser}, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}
	if len(originUser) == 0 {
		return nil
	}

	assignUserIDs, err := client.GetFeishuUserIdList(assignUsers, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}

	var tableRows []string
	for _, record := range workflow.Record.InstanceRecords {
		tableRow := fmt.Sprintf(FeishuAuditResultLayout, record.Instance.Name, record.Task.Score, record.Task.PassRate*100)
		tableRows = append(tableRows, tableRow)
	}
	auditResult := strings.Join(tableRows, ",")

	approvalInstCode, err := client.CreateApprovalInstance(ctx, im.ProcessCode, workflow.Subject, originUser[0],
		assignUserIDs, auditResult, string(workflow.ProjectId), workflow.Desc, url)
	if err != nil {
		return err
	}

	instDetail, err := client.GetApprovalInstDetail(ctx, *approvalInstCode)
	if err != nil {
		return err
	}

	s := model.GetStorage()
	feishuInst := &model.FeishuInstance{
		ApproveInstanceCode: *approvalInstCode,
		WorkflowId:          workflow.WorkflowId,
		TaskID:              *instDetail.TaskList[0].Id,
	}

	if err := s.Save(&feishuInst); err != nil {
		return err
	}

	return nil
}

func UpdateFeishuAuditStatus(ctx context.Context, im model.IM, workflowId string, user *model.User, status string, reason string) error {
	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	userId, err := client.GetFeishuUserIdList([]*model.User{user}, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}
	if len(userId) == 0 {
		return fmt.Errorf("user %s has no associated feishu account", user.Name)
	}

	s := model.GetStorage()
	feishuInst, exist, err := s.GetFeishuInstanceByWorkflowID(workflowId)
	if err != nil {
		return err
	}
	if !exist {
		return e.New("feishu instance not found")
	}

	switch status {
	case model.ApproveStatusAgree:
		err = client.ApproveApproval(ctx, im.ProcessCode, feishuInst.ApproveInstanceCode, userId[0], feishuInst.TaskID)
		if err != nil {
			return err
		}
	case model.WorkflowStatusReject:
		err = client.RejectApproval(ctx, im.ProcessCode, feishuInst.ApproveInstanceCode, userId[0], feishuInst.TaskID, reason)
		if err != nil {
			return err
		}
	default:
		return e.New("invalid approve status")
	}

	feishuInst.Status = status
	if err := s.Save(&feishuInst); err != nil {
		return err
	}

	return nil
}

func CancelFeishuAuditInst(ctx context.Context, im model.IM, workflowIDs []string, user *model.User) error {
	s := model.GetStorage()
	err := s.BatchUpdateStatusOfFeishuInstance(workflowIDs, model.WorkflowStatusCancel)
	if err != nil {
		return err
	}

	feishuInstList, err := s.GetFeishuInstanceListByWorkflowIDs(workflowIDs)
	if err != nil {
		return err
	}

	client := feishu.NewFeishuClient(im.AppKey, im.AppSecret)
	userIdList, err := client.GetFeishuUserIdList([]*model.User{user}, larkContact.UserIdTypeOpenId)
	if err != nil {
		return err
	}

	for _, feishuInst := range feishuInstList {
		inst := feishuInst
		if inst.Status != model.FeishuAuditStatusInitialized {
			log.NewEntry().Infof("feishu approval instance %v status is %v, skip cancel", inst.ApproveInstanceCode, inst.Status)
			continue
		}

		go func() {
			err = client.CancelApproval(ctx, im.ProcessCode, inst.ApproveInstanceCode, userIdList[0])
			if err != nil {
				log.NewEntry().Errorf("cancel feishu approval instance %v error: %v", inst.ApproveInstanceCode, err)
			}
		}()
	}

	return nil
}
