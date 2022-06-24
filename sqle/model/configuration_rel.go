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
			separate := strings.Index(data, "~~")
			if separate == -1 {
				i.Content = data
			} else {
				i.WorkDurationHour, _ = strconv.Atoi(data[:separate])
				i.Content = data[separate+1:]
			}
		}
	}
	return nil
}

func (i *License) encryptPassword() error {
	if i == nil {
		return nil
	}
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
