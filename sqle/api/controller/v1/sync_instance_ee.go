//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/driver"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

func createSyncInstanceTask(c echo.Context) error {
	return nil
}

func updateSyncInstanceTask(c echo.Context) error {
	return nil
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

const ActiontechDmp = "actiontech-dmp"

var (
	syncTaskSourceList = []string{ActiontechDmp}
	// todo: 使用接口获取
	dmpSupportDbType = []string{driver.DriverTypeMySQL}
)

func getSyncTaskSourceTips(c echo.Context) error {
	m := make(map[string]struct{}, 0)

	additionalParams := driver.AllAdditionalParams()
	for dbType := range additionalParams {
		m[dbType] = struct{}{}
	}

	var sourceList []SyncTaskTipsResV1
	for _, source := range syncTaskSourceList {
		var commonDbTypes []string

		// 外部平台和sqle共同支持的数据源
		switch source {
		case ActiontechDmp:
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
