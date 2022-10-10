//go:build release && enterprise
// +build release,enterprise

package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/license"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

func getLicenseUsageV1(c echo.Context) error {
	s := model.GetStorage()
	l, exist, err := s.GetLicense()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return c.JSON(http.StatusOK, &GetLicenseUsageResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: &LicenseUsageV1{
				// only admin user is valid if there is no license
				UsersUsage: LicenseUsageItem{
					ResourceType:     "user",
					ResourceTypeDesc: "用户",
					Used:             1,
					Limit:            1,
					IsLimited:        true,
				},
				InstancesUsage: nil,
			},
		})
	}

	permission, _, err := license.DecodeLicense(l.Content)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataInvalid, license.ErrInvalidLicense))
	}

	usersTotal, err := s.GetAllUserCount()
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("get users count failed: %v", err))
	}

	instCounts, err := s.GetAllInstanceCount()
	if err != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("get instance count failed: %v", err))
	}

	var instUsage []LicenseUsageItem
	var customUsed uint
	customDbType := "custom"
	for dbType, limit := range permission.NumberOfInstanceOfEachType {
		if dbType == customDbType {
			// we will handle it after loop if dbType is custom
			continue
		}
		var usedCount uint
		for _, used := range instCounts {
			if used.DBType != dbType {
				continue
			}
			usedCount = uint(used.Count)
			break
		}
		if usedCount > uint(limit.Count) {
			customUsed += usedCount - uint(limit.Count)
		}

		instUsage = append(instUsage, LicenseUsageItem{
			ResourceType:     dbType,
			ResourceTypeDesc: dbType,
			Used:             usedCount,
			Limit:            uint(limit.Count),
			IsLimited:        true,
		})
	}

	// count custom type
	limit, ok := permission.NumberOfInstanceOfEachType[customDbType]
	if ok {
		instUsage = append(instUsage, LicenseUsageItem{
			ResourceType:     customDbType,
			ResourceTypeDesc: customDbType,
			Used:             customUsed,
			Limit:            uint(limit.Count),
			IsLimited:        true,
		})
	}

	return c.JSON(http.StatusOK, &GetLicenseUsageResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: &LicenseUsageV1{
			UsersUsage: LicenseUsageItem{
				ResourceType:     "user",
				ResourceTypeDesc: "用户",
				Used:             uint(usersTotal),
				Limit:            uint(permission.UserCount),
				IsLimited:        true,
			},
			InstancesUsage: instUsage,
		},
	})
}
