//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	e "errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/actiontech/sqle/sqle/api/controller"
	instSync "github.com/actiontech/sqle/sqle/pkg/sync_task"
	"github.com/labstack/echo/v4"
)

var (
	ErrSyncInstanceTaskNotExist = func(taskId int) error {
		return errors.New(errors.DataNotExist, fmt.Errorf("sync instance task [%s] not exist", taskId))
	}
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
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func updateSyncInstanceTask(c echo.Context) error {
	req := new(UpdateSyncInstanceTaskReqV1)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	taskId := c.Param("task_id")

	s := model.GetStorage()
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	syncTask, exist, err := s.GetSyncInstanceTaskById(uint(taskIdInt))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrSyncInstanceTaskNotExist(taskIdInt))
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
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func deleteSyncInstanceTask(c echo.Context) error {
	taskIdStr := c.Param("task_id")

	s := model.GetStorage()
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	syncTask, exist, err := s.GetSyncInstanceTaskById(uint(taskId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrSyncInstanceTaskNotExist(taskId))
	}

	if err := s.Delete(&syncTask); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

const TriggerSyncInstanceTimeout = 30 * time.Second

func triggerSyncInstance(c echo.Context) error {
	taskIdStr := c.Param("task_id")
	s := model.GetStorage()

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	task, exist, err := s.GetSyncInstanceTaskById(uint(taskId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrSyncInstanceTaskNotExist(taskId))
	}

	l := log.Logger().WithField("action", "trigger_sync_instance_task")
	syncInstanceTaskEntity := instSync.NewSyncInstanceTask(l, uint(taskId), task.Source, task.URL, task.Version, task.DbType, task.RuleTemplate.Name)

	ctx, cancel := context.WithTimeout(c.Request().Context(), TriggerSyncInstanceTimeout)
	defer cancel()

	syncFunc := syncInstanceTaskEntity.GetSyncInstanceTaskFunc(ctx)
	syncFunc()

	if ctx.Err() != nil && e.Is(ctx.Err(), context.DeadlineExceeded) {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("sync instance task timeout: %v,timeout configuration: %v seconds", ctx.Err(), TriggerSyncInstanceTimeout))
	}

	if getSyncTaskStatus(s, task.ID) == model.SyncInstanceStatusFailed {
		return controller.JSONBaseErrorReq(c, e.New("sync instance task failed"))
	}

	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func getSyncTaskStatus(s *model.Storage, taskId uint) string {
	syncTask, exist, err := s.GetSyncInstanceTaskById(taskId)
	if err != nil || !exist {
		return model.SyncInstanceStatusFailed
	}
	return syncTask.LastSyncStatus
}

func getSyncInstanceTaskList(c echo.Context) error {
	s := model.GetStorage()

	tasks, err := s.GetAllSyncInstanceTasks()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	tasksRes := make([]InstanceTaskResV1, len(tasks))
	for i, t := range tasks {
		tasksRes[i].ID = int(t.ID)
		tasksRes[i].Source = t.Source
		tasksRes[i].Version = t.Version
		tasksRes[i].URL = t.URL
		tasksRes[i].DbType = t.DbType
		tasksRes[i].LastSyncStatus = t.LastSyncStatus
		tasksRes[i].LastSyncSuccessTime = t.LastSyncSuccessTime
	}
	return c.JSON(http.StatusOK, GetSyncInstanceTaskListResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    tasksRes,
	})
}

func getSyncInstanceTask(c echo.Context) error {
	taskIdStr := c.Param("task_id")
	s := model.GetStorage()

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, err))
	}

	task, exist, err := s.GetSyncInstanceTaskById(uint(taskId))
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrSyncInstanceTaskNotExist(taskId))
	}
	taskRes := InstanceTaskDetailResV1{
		ID:                   int(task.ID),
		Source:               task.Source,
		Version:              task.Version,
		URL:                  task.URL,
		DbType:               task.DbType,
		RuleTemplate:         task.RuleTemplate.Name,
		SyncInstanceInterval: task.SyncInstanceInterval,
	}
	return c.JSON(http.StatusOK, GetSyncInstanceTaskResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    taskRes,
	})
}

var (
	syncTaskSourceList = []string{model.SyncTaskSourceActiontechDmp}
	// todo: 使用接口获取
	dmpSupportDbType = []string{driverV2.DriverTypeMySQL}
)

func getSyncTaskSourceTips(c echo.Context) error {
	m := make(map[string]struct{}, 0)

	drivers := driver.GetPluginManager().AllDrivers()
	for _, dbType := range drivers {
		m[dbType] = struct{}{}
	}

	var sourceList []SyncTaskTipsResV1
	for _, source := range syncTaskSourceList {
		var commonDbTypes []string

		// 外部平台和sqle共同支持的数据源
		switch source {
		case model.SyncTaskSourceActiontechDmp:
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
