//go:build release
// +build release

package license

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/labstack/echo/v4"
)

var std = &checker{
	mutex:      &sync.RWMutex{},
	permission: map[string]int64{},
}

type checker struct {
	mutex *sync.RWMutex
	// 各种限制
	permission map[string]int64
	// 机器信息
	hardwareInfo string
	// license中的机器信息
	limitInstallLocation string
	// 限制工作时长(天)
	WorkDurationDay int
	// 已工作时长（小时）
	timerHour int
}

func (c *checker) GetPermission(key string) int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.permission[key]
}

const (
	LimitTypeUser = "user"
)

// 当 skipSubsequentCheck==true 或 checkResult==false 或 err!=nil 时将会停止继续校验,并返回当前校验结果
type LicenseChecker func(c echo.Context) (skipSubsequentCheck bool, checkResult bool, err error)

// 有序检查列表
var LicenseCheckers = []LicenseChecker{
	// 跳过永远放行的接口
	func(c echo.Context) (skipSubsequentCheck bool, checkResult bool, err error) {
		path := c.Path()
		if c.Request().Method == http.MethodGet ||
			strings.TrimSuffix(path, "/") == "/v1/login" ||
			strings.TrimSuffix(path, "/") == "/v2/login" ||
			strings.TrimSuffix(path, "/") == "/v1/logout" ||
			strings.TrimSuffix(path, "/") == "/v1/configurations/license" ||
			strings.TrimSuffix(path, "/") == "/v1/configurations/license/check" ||
			strings.TrimSuffix(path, "/") == "/v1/configurations/license/info" {
			return true, true, nil
		}
		return false, true, nil
	},
	// 机器信息不匹配
	func(c echo.Context) (skipSubsequentCheck bool, checkResult bool, err error) {
		std.mutex.RLock()
		defer std.mutex.RUnlock()

		if std.hardwareInfo != std.limitInstallLocation {
			return true, false, nil
		}
		return false, true, nil
	},
	// license过期
	func(c echo.Context) (skipSubsequentCheck bool, checkResult bool, err error) {
		std.mutex.RLock()
		defer std.mutex.RUnlock()

		if std.timerHour/24 >= std.WorkDurationDay {
			return true, false, nil
		}
		return false, true, nil
	},
	// 添加用户
	func(c echo.Context) (skipSubsequentCheck bool, checkResult bool, err error) {
		s := model.GetStorage()
		if c.Request().Method == http.MethodPost && strings.TrimSuffix(c.Path(), "/") == "/v1/users" {
			count, err := s.GetAllUserCount()
			return true, count+1 <= std.GetPermission(LimitTypeUser), err
		}
		return false, true, nil
	},
	// 添加数据源
	func(c echo.Context) (skipSubsequentCheck bool, checkResult bool, err error) {
		return InstanceLicenseChecker(c)
	},
}

var InstanceLicenseChecker = func(c echo.Context) (skipSubsequentCheck bool, checkResult bool, err error) {
	if !(c.Request().Method == http.MethodPost &&
		(strings.TrimSuffix(c.Path(), "/") == "/v1/projects/:project_name/instances" ||
			strings.TrimSuffix(c.Path(), "/") == "/v2/projects/:project_name/instances")) {
		return false, true, nil
	}

	dbType, err := getDBTypeWithReq(c)
	if err != nil {
		return true, false, err
	}

	s := model.GetStorage()
	count, err := s.GetAllInstanceCount()
	if err != nil {
		return true, false, err
	}

	addedCount := false
	for _, dbCount := range count {
		if dbCount.DBType == dbType {
			dbCount.Count++
			addedCount = true
			break
		}
	}
	if !addedCount {
		count = append(count, &model.TypeCount{
			DBType: dbType,
			Count:  1,
		})
	}

	return true, calculateIdleCustomAmount(count) >= 0, nil
}

// 复制结构体是为了防止循环依赖
type CreateInstanceReqV1 struct {
	Name     string `json:"instance_name" form:"instance_name" example:"test" valid:"required,name"`
	DBType   string `json:"db_type" form:"db_type" example:"mysql"`
	User     string `json:"db_user" form:"db_user" example:"root" valid:"required"`
	Host     string `json:"db_host" form:"db_host" example:"10.10.10.10" valid:"required,ip_addr|uri|hostname|hostname_rfc1123"`
	Port     string `json:"db_port" form:"db_port" example:"3306" valid:"required,port"`
	Password string `json:"db_password" form:"db_password" example:"123456" valid:"required"`
}

func getDBTypeWithReq(c echo.Context) (dbType string, err error) {
	bodyBytes, _ := ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	req := new(CreateInstanceReqV1)
	err = controller.BindAndValidateReq(c, req)
	if err != nil {
		return "", err
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	dbType = fmt.Sprintf("%v", req.DBType)
	if dbType == "" {
		dbType = driver.DriverTypeMySQL
	}
	return dbType, err
}

const CustomTypeKey = "custom"

func calculateIdleCustomAmount(allCount []*model.TypeCount) int64 {
	var custom int64 = 0
	for _, count := range allCount {
		if t := count.Count - std.GetPermission(count.DBType); t > 0 {
			custom += t
		}
	}
	return std.GetPermission(CustomTypeKey) - custom
}

func Check(c echo.Context) (bool, error) {
	for _, licenseChecker := range LicenseCheckers {
		skipSubsequentCheck, checkResult, err := licenseChecker(c)
		if err != nil {
			return false, err
		}
		if !checkResult {
			return false, nil
		}
		if skipSubsequentCheck {
			return checkResult, nil
		}
	}

	return true, nil
}

func InitChecker(license string, workDurationHour int) error {
	info, err := CollectHardwareInfo()
	if err != nil {
		return fmt.Errorf("collect hardware info error: %v", err)
	}

	std.mutex.Lock()
	std.hardwareInfo = info
	std.timerHour = workDurationHour
	std.mutex.Unlock()

	go func() {
		for {
			time.Sleep(time.Hour)
			std.mutex.Lock()
			std.timerHour++
			std.mutex.Unlock()
			s := model.GetStorage()
			l, exist, err := s.GetLicense()
			if exist && err == nil {
				l.WorkDurationHour = std.timerHour
				_ = s.Save(l)
			}
		}
	}()

	return UpdateLicense(license)
}

func UpdateLicense(license string) error {
	std.mutex.Lock()
	defer std.mutex.Unlock()

	std.permission = map[string]int64{}
	std.limitInstallLocation = ""

	if license == "" {
		return nil
	}

	permission, info, err := DecodeLicense(license)
	if err != nil {
		return err
	}

	std.permission[LimitTypeUser] = int64(permission.UserCount)
	for n, i := range permission.NumberOfInstanceOfEachType {
		std.permission[n] = int64(i.Count)
	}
	std.limitInstallLocation = info
	std.WorkDurationDay = permission.WorkDurationDay

	return nil
}

func ResetWorkedDuration() {
	std.mutex.Lock()
	std.timerHour = 0
	std.mutex.Unlock()
}
