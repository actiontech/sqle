//go:build !enterprise
// +build !enterprise

package im

import (
	"context"
	e "errors"

	"github.com/actiontech/sqle/sqle/model"
)

var ErrCommunityEditionNotSupportFeishuApproval = e.New("community edition not support feishu approval")

func CreateFeishuApprovalTemplate(ctx context.Context, im model.IM) error {
	return ErrCommunityEditionNotSupportFeishuApproval
}

func CreateFeishuApprovalInst(ctx context.Context, im model.IM, workflow *model.Workflow, assignUsers []*model.User, url string) error {
	return ErrCommunityEditionNotSupportFeishuApproval
}

func UpdateFeishuApprovalStatus(ctx context.Context, im model.IM, workflowId uint, user *model.User, status string, reason string) error {
	return ErrCommunityEditionNotSupportFeishuApproval
}

func CancelFeishuApprovalInst(ctx context.Context, im model.IM, workflowID uint, user *model.User) error {
	return ErrCommunityEditionNotSupportFeishuApproval
}
