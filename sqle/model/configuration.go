package model

import (
	"fmt"
	"strconv"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"

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
	data, err := utils.AesEncrypt(i.Password)
	if err != nil {
		return err
	}
	i.SecretPassword = data
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
		data, err := utils.AesDecrypt(i.SecretPassword)
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
	return smtpC, true, errors.New(errors.ConnectStorageError, err)
}

// WeChatConfiguration store WeChat configuration.
type WeChatConfiguration struct {
	Model
	EnableWeChatNotify bool   `json:"enable_wechat_notify" gorm:"not null"`
	CorpID             string `json:"corp_id" gorm:"not null"`
	CorpPassword       string `json:"-"`
	CorpSecret         string `json:"corp_secret" gorm:"not null"`
	AgentID            int    `json:"agent_id" gorm:"not null"`
	SafeEnabled        bool   `json:"safe_enabled" gorm:"not null"`
	ProxyIP            string `json:"proxy_ip"`
}

func (i *WeChatConfiguration) TableName() string {
	return fmt.Sprintf("%v_wechat", globalConfigurationTablePrefix)
}

// BeforeSave is a hook implement gorm model before exec create.
func (i *WeChatConfiguration) BeforeSave() error {
	return i.encryptPassword()
}

func (i *WeChatConfiguration) encryptPassword() error {
	if i == nil {
		return nil
	}
	data, err := utils.AesEncrypt(i.CorpPassword)
	if err != nil {
		return err
	}
	i.CorpSecret = data
	return nil
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db.
func (i *WeChatConfiguration) AfterFind() error {
	err := i.decryptPassword()
	if err != nil {
		log.NewEntry().Errorf("decrypt password for WeChat server configuration failed, error: %v", err)
	}
	return nil
}

func (i *WeChatConfiguration) decryptPassword() error {
	if i == nil {
		return nil
	}
	if i.CorpPassword == "" {
		data, err := utils.AesDecrypt(i.CorpSecret)
		if err != nil {
			return err
		}
		i.CorpPassword = data
	}
	return nil
}

func (s *Storage) GetWeChatConfiguration() (*WeChatConfiguration, bool, error) {
	wechatC := new(WeChatConfiguration)
	err := s.db.Last(wechatC).Error
	if err == gorm.ErrRecordNotFound {
		return wechatC, false, nil
	}
	return wechatC, true, errors.New(errors.ConnectStorageError, err)
}

// LDAPConfiguration store ldap server configuration.
type LDAPConfiguration struct {
	Model
	// whether the ldap is enabled
	Enable bool `json:"enable" gorm:"not null"`
	// ldap server's ip
	Host string `json:"host" gorm:"not null"`
	// ldap server's port
	Port string `json:"port" gorm:"not null"`
	// the DN of the ldap administrative user for verification
	ConnectDn string `json:"connect_dn" gorm:"not null"`
	// the password of the ldap administrative user for verification
	ConnectPassword string `json:"-" gorm:"-"`
	// the secret password of the ldap administrative user for verification
	ConnectSecretPassword string `json:"connect_secret_password" gorm:"not null"`
	// base dn used for ldap verification
	BaseDn string `json:"base_dn" gorm:"not null"`
	// the key corresponding to the user name in ldap
	UserNameRdnKey string `json:"ldap_user_name_rdn_key" gorm:"not null"`
	// the key corresponding to the user email in ldap
	UserEmailRdnKey string `json:"ldap_user_email_rdn_key" gorm:"not null"`
}

func (i *LDAPConfiguration) TableName() string {
	return fmt.Sprintf("%v_ldap", globalConfigurationTablePrefix)
}

// BeforeSave is a hook implement gorm model before exec create
func (i *LDAPConfiguration) BeforeSave() error {
	return i.encryptPassword()
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db
func (i *LDAPConfiguration) AfterFind() error {
	err := i.decryptPassword()
	if err != nil {
		log.NewEntry().Errorf("decrypt password for ldap administrative user failed, error: %v", err)
	}
	return nil
}

func (i *LDAPConfiguration) decryptPassword() error {
	if i == nil {
		return nil
	}
	if i.ConnectPassword == "" {
		data, err := utils.AesDecrypt(i.ConnectSecretPassword)
		if err != nil {
			return err
		} else {
			i.ConnectPassword = data
		}
	}
	return nil
}

func (i *LDAPConfiguration) encryptPassword() error {
	if i == nil {
		return nil
	}
	data, err := utils.AesEncrypt(i.ConnectPassword)
	if err != nil {
		return err
	}
	i.ConnectSecretPassword = data
	return nil
}

func (s *Storage) GetLDAPConfiguration() (*LDAPConfiguration, bool, error) {
	ldapC := new(LDAPConfiguration)
	err := s.db.Last(ldapC).Error
	if err == gorm.ErrRecordNotFound {
		return ldapC, false, nil
	}
	return ldapC, true, errors.New(errors.ConnectStorageError, err)
}

const (
	SystemVariableWorkflowExpiredHours = "system_variable_workflow_expired_hours"
)

// SystemVariable store misc K-V.
type SystemVariable struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"not null"`
}

func (s *Storage) GetWorkflowExpiredHoursOrDefault() (int64, error) {
	var svs []SystemVariable
	err := s.db.Find(&svs).Error
	if err != nil {
		return 0, errors.New(errors.ConnectStorageError, err)
	}

	for _, sv := range svs {
		if sv.Key == SystemVariableWorkflowExpiredHours {
			wfExpiredHs, err := strconv.ParseInt(sv.Value, 10, 64)
			if err != nil {
				return 0, err
			}
			return wfExpiredHs, nil
		}
	}

	return 30 * 24, nil
}

type License struct {
	Model
	Content string `json:"content" gorm:"type:text;"`
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
