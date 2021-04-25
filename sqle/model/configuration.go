package model

import (
	"fmt"

	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/log"
	"actiontech.cloud/universe/ucommon/v4/util"

	"github.com/jinzhu/gorm"
)

const globalConfigurationTablePrefix = "global_configuration"

// SMTPConfiguration store SMTP server configuration.
type SMTPConfiguration struct {
	Model
	Host           string `json:"smtp_host" gorm:"column:smtp_host; not null"`
	Port           string `json:"smtp_port" gorm:"column:smtp_port; not null"`
	Username       string `json:"smtp_username" gorm:"column:smtp_username; not null"`
	Password       string `json:"-"`
	SecretPassword string `json:"secret_smtp_password" gorm:"column:secret_smtp_password; not null"`
}

func (i *SMTPConfiguration) TableName() string {
	return fmt.Sprintf("%v_smtp", globalConfigurationTablePrefix)
}

// BeforeSave is a hook implement gorm model before exec create.
func (i *SMTPConfiguration) BeforeSave() error {
	return i.encryptPassword()
}

func (i *SMTPConfiguration) encryptPassword() error {
	if i == nil {
		return nil
	}
	if i.SecretPassword == "" {
		data, err := util.AesEncrypt(i.Password)
		if err != nil {
			return err
		}
		i.SecretPassword = data
	}
	return nil
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db.
func (i *SMTPConfiguration) AfterFind() error {
	err := i.decryptPassword()
	if err != nil {
		log.NewEntry().Errorf("decrypt password for SMTP server configuration failed, error: %v", err)
	}
	return nil
}

func (i *SMTPConfiguration) decryptPassword() error {
	if i == nil {
		return nil
	}
	if i.Password == "" {
		data, err := util.AesDecrypt(i.SecretPassword)
		if err != nil {
			return err
		}
		i.Password = data
	}
	return nil
}

func (s *Storage) GetSMTPConfiguration() (*SMTPConfiguration, bool, error) {
	smtpC := new(SMTPConfiguration)
	err := s.db.Last(smtpC).Error
	if err == gorm.ErrRecordNotFound {
		return smtpC, false, nil
	}
	return smtpC, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
