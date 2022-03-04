//go:build release
// +build release

package license

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

var std = checker{
	mutex:      &sync.RWMutex{},
	permission: map[LimitType]int64{},
}

type checker struct {
	mutex                *sync.RWMutex
	permission           map[LimitType]int64
	hardwareInfo         string
	limitInstallLocation string
	expireDate           time.Time
}

type LimitType string

const (
	LimitTypeUser     LimitType = "user"
	LimitTypeInstance LimitType = "instance"
)

func Check(c echo.Context) (bool, error) {
	path := c.Path()
	method := c.Request().Method
	s := model.GetStorage()

	std.mutex.RLock()
	if std.expireDate.Before(time.Now()) || std.hardwareInfo != std.limitInstallLocation {
		if method == http.MethodGet ||
			strings.TrimSuffix(path, "/") == "/v1/login" ||
			strings.TrimSuffix(path, "/") == "/v1/configurations/license" ||
			strings.TrimSuffix(path, "/") == "/v1/configurations/license/check" ||
			strings.TrimSuffix(path, "/") == "/v1/configurations/license/info" {
			std.mutex.RUnlock()
			return true, nil
		}
		std.mutex.RUnlock()
		return false, nil
	}

	{ // add user
		if method == http.MethodPost && strings.TrimSuffix(path, "/") == "/v1/users" {
			count, err := s.GetAllUserCount()
			std.mutex.RUnlock()
			return count+1 <= std.permission[LimitTypeUser], err
		}
	}
	{ // add instance
		if method == http.MethodPost && strings.TrimSuffix(path, "/") == "/v1/instances" {
			count, err := s.GetAllInstanceCount()
			std.mutex.RUnlock()
			return count+1 <= std.permission[LimitTypeInstance], err
		}
	}

	std.mutex.RUnlock()
	return true, nil
}

func InitChecker(license string) error {
	info, err := CollectHardwareInfo()
	if err != nil {
		return err
	}

	std.mutex.Lock()
	std.hardwareInfo = info
	std.mutex.Unlock()

	return UpdateLicense(license)
}

func UpdateLicense(license string) error {
	std.permission = map[LimitType]int64{}
	std.limitInstallLocation = ""

	if license == "" {
		return nil
	}

	permission, info, err := DecodeLicense(license)
	if err != nil {
		return err
	}
	data, err := time.Parse(ExpireDateFormat, permission.ExpireDate)
	if err != nil {
		return err
	}

	std.mutex.Lock()
	std.permission[LimitTypeUser] = int64(permission.UserCount)
	std.permission[LimitTypeInstance] = int64(permission.InstanceCount)
	std.limitInstallLocation = info
	std.expireDate = data
	std.mutex.Unlock()

	return nil
}
