//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
)

func (h AfterCreateProject) Handle(ctx context.Context, currentUserId string, projectId string, extraParams string) error {
	s := model.GetStorage()
	// 添加默认模板
	td := model.DefaultWorkflowTemplate(projectId)
	err := s.SaveWorkflowTemplate(td)

	// 添加默认推送报告
	err = s.CreateDefaultReportPushConfigIfNotExist(projectId)
	if err != nil {
		return err
	}
	return nil
}
func (h BeforeDeleteProject) Handle(ctx context.Context, currentUserId string, projectId string, extraParams string) error {
	s := model.GetStorage()
	has, err := s.HasNotEndWorkflowByProjectId(projectId)
	if err != nil {
		return err
	}
	if has {
		return errors.New(errors.UserNotPermission, fmt.Errorf("there are unfinished work orders, and the current project cannot be archived"))
	}
	configs, err := s.GetReportPushConfigListInProject(projectId)
	if err != nil {
		return err
	}
	for _, config := range configs {
		if config.Enabled {
			return fmt.Errorf("current project has running push job for %v, you need to modify the configuration to stop it ", config.Type)
		}
	}
	instAuditPlans, err := s.GetAuditPlansByProjectId(projectId)
	if err != nil {
		return err
	}
	for _, instAP := range instAuditPlans {
		if instAP.ActiveStatus == model.ActiveStatusNormal {
			return fmt.Errorf("current project has running audit plan, you need to stop or delete it")
		}
	}

	return nil
}

func (h AfterDeleteProject) Handle(ctx context.Context, currentUserId string, projectId string, extraParams string) error {
	s := model.GetStorage()
	err := s.RemoveProjectRelateData(model.ProjectUID(projectId))
	if err != nil {
		return err
	}
	err = s.DeleteReportPushConfigInProject(projectId)
	if err != nil {
		return err
	}
	return nil
}

func (h BeforeArchiveProject) Handle(ctx context.Context, currentUserId string, projectId string, extraParams string) error {
	s := model.GetStorage()
	has, err := s.HasNotEndWorkflowByProjectId(projectId)
	if err != nil {
		return err
	}
	if has {
		return errors.New(errors.UserNotPermission, fmt.Errorf("there are unfinished work orders, and the current project cannot be archived"))
	}
	return nil
}
