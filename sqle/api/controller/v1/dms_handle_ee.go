//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
)

func (h AfterCreateProject) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	s := model.GetStorage()
	td := model.DefaultWorkflowTemplate(dataResourceId)
	return s.SaveWorkflowTemplate(td)
}
func (h BeforeDeleteProject) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	s := model.GetStorage()
	has, err := s.HasNotEndWorkflowByProjectId(dataResourceId)
	if err != nil {
		return err
	}
	if has {
		return errors.New(errors.UserNotPermission, fmt.Errorf("there are unfinished work orders, and the current project cannot be archived"))
	}
	return nil
}

func (h AfterDeleteProject) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	s := model.GetStorage()
	return s.RemoveProjectRelateData(model.ProjectUID(dataResourceId))
}

func (h BeforeArvhiveProject) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	s := model.GetStorage()
	has, err := s.HasNotEndWorkflowByProjectId(dataResourceId)
	if err != nil {
		return err
	}
	if has {
		return errors.New(errors.UserNotPermission, fmt.Errorf("there are unfinished work orders, and the current project cannot be archived"))
	}
	return nil
}
