//go:build release
// +build release

package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/jinzhu/gorm"
)

var relAse = utils.NewEncryptor(config.RelAesKey)

func init() {
	autoMigrateList = append(autoMigrateList, &License{})
}

type License struct {
	Model
	Content          string `json:"-" gorm:"-"`
	WorkDurationHour int    `json:"work_duration_hour"`
	ContentSecret    string `json:"content_secret" gorm:"type:text"`
}

// BeforeSave is a hook implement gorm model before exec create
func (i *License) BeforeSave() error {
	return i.encryptPassword()
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db
func (i *License) AfterFind() error {
	err := i.decryptPassword()
	if err != nil {
		log.NewEntry().Errorf("decrypt license failed, error: %v", err)
	}
	return nil
}

func (i *License) decryptPassword() error {
	if i == nil {
		return nil
	}
	if i.Content == "" {
		data, err := relAse.AesDecrypt(i.ContentSecret)
		if err != nil {
			return err
		} else {
			/*
				保存license文本时拼入了运行时长, 取出时需要再次切开,拼入详见下方encryptPassword()
			*/
			separate := strings.Index(data, "~~")
			if separate == -1 {
				i.Content = data
			} else {
				i.WorkDurationHour, _ = strconv.Atoi(data[:separate])
				i.Content = data[separate+2:]
			}
		}
	}
	return nil
}

func (i *License) encryptPassword() error {
	if i == nil {
		return nil
	}
	/*
		拼接后的文本类似于↓
		10~~This license is for: &{WorkDurationDay:60 Version:演示环境 UserCount:10 NumberOfInstanceOfEachType:map[custom:{DBType:custom Count:3} mysql:{DBType:mysql Count:3}]};;1_XBm2N8t7coUEuhg7J5V8o9AYlhUfq2AmndctDHCxz9u~GyOKyJW0e~sVDuQVbkaKzAZQvpsGBqB~liD7svsTvbzD3ZHfdvEtSPkoYSnk2nxrYJLrW0wmzTVIicDWg1Dp2MICEK9T09Od3Xn1u4XWO7e182mzrHqncLOGKXJKlSrCsL_kWY6o6w8pWKL1Xdzduyq4uLdXuL9E6oOzyUMF3rYlnOhvoOwdoE;;9S~ViK_ZoRx8045cLM5pTZXCCpDEY_yxjfaLYGBMMOKyWpgc

		`~~`前面的10代表sqle已运行时长,`~~`后面的内容表示license原文, 拼接运行时长因为sqle license没有记录到期时间, 需要防止用户手动修改数据库中的数据重置license时长
	*/
	data, err := relAse.AesEncrypt(fmt.Sprintf("%v~~%v", i.WorkDurationHour, i.Content))
	if err != nil {
		return err
	}
	i.ContentSecret = data
	return nil
}

func (l *License) TableName() string {
	return fmt.Sprintf("%v_license", globalConfigurationTablePrefix)
}

func (s *Storage) GetLicense() (*License, bool, error) {
	license := new(License)
	err := s.db.Last(license).Error
	if err == gorm.ErrRecordNotFound {
		return license, false, nil
	}
	return license, true, errors.New(errors.ConnectStorageError, err)
}
