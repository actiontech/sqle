//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/model"
)

func (h AfterCreateNamespace) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	return nil
}
func (h BeforeDeleteNamespace) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
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

func (h AfterDeleteNamespace) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
	s := model.GetStorage()
	return s.RemoveProjectRelateData(model.ProjectUID(dataResourceId))
}

func (h BeforeArvhiveNamespace) Hanle(ctx context.Context, currentUserId string, dataResourceId string) error {
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
