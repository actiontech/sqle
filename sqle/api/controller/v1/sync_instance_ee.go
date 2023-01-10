//go:build enterprise
// +build enterprise

package v1

import (
	"net/http"

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

var syncTaskSourceList = []string{"actiontech-dmp"}

func getSyncTaskSourceList(c echo.Context) error {
	var sourceList []SyncTaskSourceListResV1
	for _, source := range syncTaskSourceList {
		sourceList = append(sourceList, SyncTaskSourceListResV1{
			Source: source,
		})
	}

	return c.JSON(http.StatusOK, GetSyncTaskSourceListResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    sourceList,
	})
}
