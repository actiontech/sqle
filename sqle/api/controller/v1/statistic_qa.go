//go:build !release
// +build !release

package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func getLicenseUsageV1(c echo.Context) error {
	s := model.GetStorage()
	usersTotal, err := s.GetAllUserCount()
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("get users count failed: %v", err))
	}

	instCounts, err := s.GetAllInstanceCount()
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("get instance count failed: %v", err))
	}

	instUsage := make([]LicenseUsageItem, len(instCounts))
	for i, count := range instCounts {
		instUsage[i] = LicenseUsageItem{
			ResourceType:     count.DBType,
			ResourceTypeDesc: count.DBType,
			Used:             uint(count.Count),
			Limit:            0,
			IsLimited:        false,
		}
	}

	return c.JSON(http.StatusOK, &GetLicenseUsageResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &LicenseUsageV1{
			UsersUsage: LicenseUsageItem{
				ResourceType:     "user",
				ResourceTypeDesc: "用户",
				Used:             uint(usersTotal),
				Limit:            0,
				IsLimited:        false,
			},
			InstancesUsage: instUsage,
		},
	})
}
