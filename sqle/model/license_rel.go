//go:build release
// +build release

package model

import (
	"fmt"

	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/license"
	"github.com/actiontech/sqle/sqle/utils"

	"gorm.io/gorm"
)

var relAse = utils.NewEncryptor(config.RelAesKey)

func init() {
	autoMigrateList = append(autoMigrateList, &License{})
}

type License struct {
	Model
	WorkDurationHour int              `json:"work_duration_hour"`
	Content          *license.License `json:"content_secret" gorm:"column:content_secret;type:text"`
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
