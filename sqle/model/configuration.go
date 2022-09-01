package model

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/jinzhu/gorm"
)

const globalConfigurationTablePrefix = "global_configuration"

// SMTPConfiguration store SMTP server configuration.
type SMTPConfiguration struct {
	Model
	EnableSMTPNotify sql.NullBool `json:"enable_smtp_notify" gorm:"default:true"`
	Host             string       `json:"smtp_host" gorm:"column:smtp_host; not null"`
	Port             string       `json:"smtp_port" gorm:"column:smtp_port; not null"`
	Username         string       `json:"smtp_username" gorm:"column:smtp_username; not null"`
	Password         string       `json:"-" gorm:"-"`
	SecretPassword   string       `json:"secret_smtp_password" gorm:"column:secret_smtp_password; not null"`
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
	EnableWeChatNotify  bool   `json:"enable_wechat_notify" gorm:"not null"`
	CorpID              string `json:"corp_id" gorm:"not null"`
	CorpSecret          string `json:"-" gorm:"-"`
	EncryptedCorpSecret string `json:"encrypted_corp_secret" gorm:"not null"`
	AgentID             int    `json:"agent_id" gorm:"not null"`
	SafeEnabled         bool   `json:"safe_enabled" gorm:"not null"`
	ProxyIP             string `json:"proxy_ip"`
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
	data, err := utils.AesEncrypt(i.CorpSecret)
	if err != nil {
		return err
	}
	i.EncryptedCorpSecret = data
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
	if i.CorpSecret == "" {
		data, err := utils.AesDecrypt(i.EncryptedCorpSecret)
		if err != nil {
			return err
		}
		i.CorpSecret = data
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
	// whether the ssl is enabled
	EnableSSL bool `json:"enable_ssl" gorm:"not null"`
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

// Oauth2Configuration store ldap server configuration.
type Oauth2Configuration struct {
	Model
	EnableOauth2    bool   `json:"enable_oauth2" gorm:"column:enable_oauth2"`
	ClientID        string `json:"client_id" gorm:"column:client_id"`
	ClientKey       string `json:"-" gorm:"-"`
	ClientSecret    string `json:"client_secret" gorm:"client_secret"`
	ClientHost      string `json:"client_host" gorm:"column:client_host"`
	ServerAuthUrl   string `json:"server_auth_url" gorm:"column:server_auth_url"`
	ServerTokenUrl  string `json:"server_token_url" gorm:"column:server_token_url"`
	ServerUserIdUrl string `json:"server_user_id_url" gorm:"column:server_user_id_url"`
	Scopes          string `json:"scopes" gorm:"column:scopes"`
	AccessTokenTag  string `json:"access_token_tag" gorm:"column:access_token_tag"`
	UserIdTag       string `json:"user_id_tag" gorm:"column:user_id_tag"`
	LoginTip        string `json:"login_tip" gorm:"column:login_tip; default:'使用第三方账户登录'"`
}

func (i *Oauth2Configuration) GetScopes() []string {
	return strings.Split(i.Scopes, ",")
}

func (i *Oauth2Configuration) SetScopes(s []string) {
	i.Scopes = strings.Join(s, ",")
}

func (i *Oauth2Configuration) TableName() string {
	return fmt.Sprintf("%v_oauth2", globalConfigurationTablePrefix)
}

// BeforeSave is a hook implement gorm model before exec create
func (i *Oauth2Configuration) BeforeSave() error {
	return i.encryptPassword()
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db
func (i *Oauth2Configuration) AfterFind() error {
	err := i.decryptPassword()
	if err != nil {
		log.NewEntry().Errorf("decrypt password for ldap administrative user failed, error: %v", err)
	}
	return nil
}

func (i *Oauth2Configuration) decryptPassword() error {
	if i == nil {
		return nil
	}
	if i.ClientKey == "" {
		data, err := utils.AesDecrypt(i.ClientSecret)
		if err != nil {
			return err
		} else {
			i.ClientKey = data
		}
	}
	return nil
}

func (i *Oauth2Configuration) encryptPassword() error {
	if i == nil {
		return nil
	}
	data, err := utils.AesEncrypt(i.ClientKey)
	if err != nil {
		return err
	}
	i.ClientSecret = data
	return nil
}

func (s *Storage) GetOauth2Configuration() (*Oauth2Configuration, bool, error) {
	oauth2C := new(Oauth2Configuration)
	err := s.db.Last(oauth2C).Error
	if err == gorm.ErrRecordNotFound {
		return oauth2C, false, nil
	}
	return oauth2C, true, errors.New(errors.ConnectStorageError, err)
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
