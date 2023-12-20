//go:build !enterprise
// +build !enterprise

package im

import (
	"context"
	e "errors"

	"github.com/actiontech/sqle/sqle/model"
)

var ErrCommunityEditionNotSupportFeishuAudit = e.New("community edition not support feishu audit")

func CreateFeishuAuditTemplate(ctx context.Context, im model.IM) error {
	return ErrCommunityEditionNotSupportFeishuAudit
}

func CreateFeishuAuditInst(ctx context.Context, im model.IM, workflow *model.Workflow, assignUsers []*model.User, url string) error {
	return ErrCommunityEditionNotSupportFeishuAudit
}

func UpdateFeishuAuditStatus(ctx context.Context, im model.IM, workflowId string, user *model.User, status string, reason string) error {
	return ErrCommunityEditionNotSupportFeishuAudit
}

func CancelFeishuAuditInst(ctx context.Context, im model.IM, workflowIDs []string, user *model.User) error {
	return ErrCommunityEditionNotSupportFeishuAudit
}
