//go:build release
// +build release

package v1

import (
	"encoding/json"
	e "errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/license"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/cluster"

	"github.com/labstack/echo/v4"
)

var ErrNoLicenseRequired = errors.New(errors.ErrAccessDeniedError, e.New("sqle-ce no license required"))

const (
	HardwareInfoFileName = "collected.infos"
	LicenseFileParamKey  = "license_file"
)

func getLicense(c echo.Context) error {
	s := model.GetStorage()
	l, exist, err := s.GetLicense()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist || l.Content == nil {
		return c.JSON(http.StatusOK, GetLicenseResV1{
			BaseRes: controller.NewBaseReq(nil),
		})
	}
	content, err := l.Content.LicenseContent.Encode()
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	items := generateLicenseItems(&l.Content.LicenseContent)

	items = append(items, LicenseItem{
		Description: "已运行时长(天)",
		Name:        "duration of running",
		Limit:       strconv.Itoa(l.WorkDurationHour / 24),
	}, LicenseItem{
		Description: "预计到期时间",
		Name:        "estimated maturity",
		// 这个时间要展示给人看, 展示成RFC3339不够友好, 也不需要展示精确的时间, 所以展示成自定义时间格式
		Limit: time.Now().Add(time.Hour * time.Duration(l.Content.Permission.WorkDurationDay*24-l.WorkDurationHour)).Format("2006-01-02"),
	})

	return c.JSON(http.StatusOK, GetLicenseResV1{
		BaseRes: controller.NewBaseReq(nil),
		Content: content,
		License: items,
	})

}

func getSQLELicenseInfo(c echo.Context) error {
	var data []byte
	if cluster.IsClusterMode {
		s := model.GetStorage()
		nodes, err := s.GetClusterNodes()
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}
		var clusterHardwareSigns = []license.ClusterHardwareSign{}
		for _, node := range nodes {
			if node.ServerId != "" && node.HardwareSign != "" {
				clusterHardwareSigns = append(clusterHardwareSigns, license.ClusterHardwareSign{
					Id:   node.ServerId,
					Sign: node.HardwareSign,
				})
			}
		}
		data, err = json.Marshal(clusterHardwareSigns)
		if err != nil {
			return controller.JSONBaseErrorReq(c, errors.New(errors.DataConflict, err))
		}
	} else {
		hardwareSign, err := license.CollectHardwareInfo()
		if err != nil {
			return controller.JSONBaseErrorReq(c, license.ErrCollectLicenseInfo)
		}
		data = []byte(hardwareSign)
	}

	c.Response().Header().Set(echo.HeaderContentDisposition,
		mime.FormatMediaType("attachment", map[string]string{"filename": HardwareInfoFileName}))

	return c.Blob(http.StatusOK, echo.MIMETextPlain, []byte(data))
}

func setLicense(c echo.Context) error {
	_, file, exist, err := controller.ReadFile(c, LicenseFileParamKey)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, license.ErrLicenseEmpty))
	}

	l := &license.License{}
	l.WorkDurationHour = 0
	err = l.Decode(file)
	if err != nil {
		return controller.JSONBaseErrorReq(c, license.ErrInvalidLicense)
	}

	collected, err := license.CollectHardwareInfo()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, license.ErrCollectLicenseInfo))
	}
	err = l.CheckHardwareSignIsMatch(collected)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, err))
	}

	s := model.GetStorage()
	err = s.Delete(&model.License{})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	err = s.Save(&model.License{Content: l, WorkDurationHour: 0})
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, controller.NewBaseReq(nil))
}

func checkLicense(c echo.Context) error {
	_, file, exist, err := controller.ReadFile(c, LicenseFileParamKey)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataNotExist, license.ErrLicenseEmpty))
	}

	l := &license.License{}
	err = l.Decode(file)
	if err != nil {
		return controller.JSONBaseErrorReq(c, license.ErrInvalidLicense)
	}

	collected, err := license.CollectHardwareInfo()
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, license.ErrCollectLicenseInfo))
	}
	err = l.CheckHardwareSignIsMatch(collected)
	if err != nil {
		return controller.JSONBaseErrorReq(c, errors.New(errors.DataParseFail, err))
	}

	items := generateLicenseItems(&l.LicenseContent)

	return c.JSON(http.StatusOK, GetLicenseResV1{
		BaseRes: controller.NewBaseReq(nil),
		Content: file,
		License: items,
	})

}

func generateLicenseItems(l *license.LicenseContent) []LicenseItem {
	items := []LicenseItem{}

	for n, i := range l.Permission.NumberOfInstanceOfEachType {
		items = append(items, LicenseItem{
			Description: fmt.Sprintf("[%v]类型实例数", n),
			Name:        n,
			Limit:       strconv.Itoa(i.Count),
		})
	}

	items = append(items, LicenseItem{
		Description: "用户数",
		Name:        "user",
		Limit:       strconv.Itoa(l.Permission.UserCount),
	})

	if l.HardwareSign != "" {
		items = append(items, LicenseItem{
			Description: "机器信息",
			Name:        "info",
			Limit:       l.HardwareSign,
		})
	}
	if len(l.ClusterHardwareSigns) > 0 {
		for _, s := range l.ClusterHardwareSigns {
			items = append(items, LicenseItem{
				Description: fmt.Sprintf("节点[%s]机器信息", s.Id),
				Name:        fmt.Sprintf("node_%s_info", s.Id),
				Limit:       s.Sign,
			})
		}
	}
	items = append(items, []LicenseItem{
		{
			Description: "SQLE版本",
			Name:        "version",
			Limit:       l.Permission.Version,
		}, {
			Description: "授权运行时长(天)",
			Name:        "work duration day",
			Limit:       strconv.Itoa(l.Permission.WorkDurationDay),
		},
	}...)

	return items
}
