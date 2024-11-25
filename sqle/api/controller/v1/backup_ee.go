//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/labstack/echo/v4"
)

func getBackupSqlList(c echo.Context) error {
	req := new(BackupSqlListReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	var backupService server.BackupService
	limit, offset := controller.GetLimitAndOffset(req.PageIndex, req.PageSize)
	sqlList, count, err := backupService.GetBackupSqlList(c.Request().Context(), c.Param("workflow_id"), req.FilterInstanceId, req.FilterExecStatus, limit, offset)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, &BackupSqlListRes{
		BaseRes:   controller.NewBaseReq(nil),
		Data:      toApiBackupSqlData(sqlList),
		TotalNums: count,
	})
}

func toApiBackupSqlData(sqlList []*server.BackupSqlData) []*BackupSqlData {
	apiSqlList := make([]*BackupSqlData, 0, len(sqlList))
	for _, sql := range sqlList {
		apiSqlList = append(apiSqlList, &BackupSqlData{
			ExecOrder:      sql.ExecOrder,
			ExecSqlID:      sql.ExecSqlID,
			OriginSQL:      sql.OriginSQL,
			OriginTaskId:   sql.OriginTaskId,
			BackupSqls:     sql.BackupSqls,
			BackupStrategy: sql.BackupStrategy,
			InstanceName:   sql.InstanceName,
			InstanceId:     fmt.Sprintf("%d", sql.InstanceId),
			ExecStatus:     sql.ExecStatus,
			Description:    sql.Description,
		})
	}
	return apiSqlList
}

func updateSqlBackupStrategy(c echo.Context) error{
	req := new(UpdateSqlBackupStrategyReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return err
	}
	taskId := c.Param("task_id")
	sqlId := c.Param("sql_id")
	sqlIdInt, err := strconv.Atoi(sqlId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	task, err := getTaskById(c.Request().Context(), taskId)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = CheckCurrentUserCanOpTask(c, task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	backupService := server.BackupService{}
	err = backupService.CanUpdateStrategyForTask(task)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	_, err = backupService.CheckSqlsTasksMappingRelationship([]uint{task.ID}, []uint{uint(sqlIdInt)})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	err = backupService.UpdateBackupStrategyForSql(sqlId, req.Strategy)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}