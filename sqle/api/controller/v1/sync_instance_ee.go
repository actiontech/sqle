//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/actiontech/sqle/sqle/api/controller"
	instSync "github.com/actiontech/sqle/sqle/pkg/sync"
	"github.com/labstack/echo/v4"
)

func createSyncInstanceTask(c echo.Context) error {
	req := new(CreateSyncInstanceTaskReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	ruleTemplate, exist, err := s.GetRuleTemplateDetailByNameAndProjectIds([]uint{model.ProjectIdForGlobalRuleTemplate}, req.GlobalRuleTemplate)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("rule template %s not exist", req.GlobalRuleTemplate))
	}

	syncTask := &model.SyncInstanceTask{
		Source:               req.Source,
		Version:              req.Version,
		URL:                  req.URL,
		DbType:               req.DbType,
		RuleTemplateID:       ruleTemplate.ID,
		SyncInstanceInterval: req.SyncInstanceInterval,
	}

	if err := s.Save(&syncTask); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instSync.ReloadInstance(context.Background(), "create new sync instance task")

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func updateSyncInstanceTask(c echo.Context) error {
	req := new(UpdateSyncInstanceTaskReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskId := c.Param("task_id")

	s := model.GetStorage()
	taskIdStr, err := strconv.Atoi(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	syncTask, exist, err := s.GetSyncTaskById(uint(taskIdStr))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("sync task %s not exist", taskId))
	}

	if req.Version != nil {
		syncTask.Version = *req.Version
	}

	if req.URL != nil {
		syncTask.URL = *req.URL
	}

	if req.GlobalRuleTemplate != nil {
		ruleTemplate, exist, err := s.GetGlobalAndProjectRuleTemplateByNameAndProjectId(*req.GlobalRuleTemplate, model.ProjectIdForGlobalRuleTemplate)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("rule template %s not exist", *req.GlobalRuleTemplate))
		}
		syncTask.RuleTemplateID = ruleTemplate.ID
	}

	if req.SyncInstanceInterval != nil {
		syncTask.SyncInstanceInterval = *req.SyncInstanceInterval
	}

	if err := s.Save(&syncTask); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	instSync.ReloadInstance(context.Background(), "update sync instance task")

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func deleteSyncInstanceTask(c echo.Context) error {
	return nil
}

func triggerSyncInstance(c echo.Context) error {
	return nil
}

func getSyncInstanceTaskList(c echo.Context) error {
	return nil
}

func getSyncInstanceTask(c echo.Context) error {
	return nil
}

var (
	syncTaskSourceList = []string{instSync.ActiontechDmp}
	// todo: 使用接口获取
	dmpSupportDbType = []string{driver.DriverTypeMySQL}
)

func getSyncTaskSourceTips(c echo.Context) error {
	m := make(map[string]struct{}, 0)

	drivers := driver.AllDrivers()
	for _, dbType := range drivers {
		m[dbType] = struct{}{}
	}

	var sourceList []SyncTaskTipsResV1
	for _, source := range syncTaskSourceList {
		var commonDbTypes []string

		// 外部平台和sqle共同支持的数据源
		switch source {
		case instSync.ActiontechDmp:
			for _, dbType := range dmpSupportDbType {
				if _, ok := m[dbType]; ok {
					commonDbTypes = append(commonDbTypes, dbType)
				}
			}
		default:
			continue
		}

		sourceList = append(sourceList, SyncTaskTipsResV1{
			Source:  source,
			DbTypes: commonDbTypes,
		})
	}

	return c.JSON(http.StatusOK, GetSyncTaskSourceTipsResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    sourceList,
	})
}
