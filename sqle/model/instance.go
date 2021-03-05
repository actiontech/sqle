package model

import (
	"encoding/json"

	"actiontech.cloud/universe/ucommon/v4/util"

	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"actiontech.cloud/universe/sqle/v4/sqle/log"
	"github.com/jinzhu/gorm"
)

const (
	DB_TYPE_MYSQL     = "mysql"
	DB_TYPE_MYCAT     = "mycat"
	DB_TYPE_SQLSERVER = "sqlserver"
)

// Instance is a table for database info
type Instance struct {
	Model
	Name               string         `json:"name" gorm:"not null;index" example:""`
	DbType             string         `json:"db_type" gorm:"not null" example:"mysql"`
	Host               string         `json:"host" gorm:"not null" example:"10.10.10.10"`
	Port               string         `json:"port" gorm:"not null" example:"3306"`
	User               string         `json:"user" gorm:"not null" example:"root"`
	Password           string         `json:"-" gorm:"-"`
	SecretPassword     string         `json:"secret_password" gorm:"not null;column:password"`
	Desc               string         `json:"desc" example:"this is a instance"`
	RuleTemplates      []RuleTemplate `json:"-" gorm:"many2many:instance_rule_template"`
	MycatConfig        *MycatConfig   `json:"-" gorm:"-"`
	MycatConfigJson    string         `json:"-" gorm:"type:text;column:mycat_config"`
	WorkflowTemplateId int            `json:"workflow_template_id"`
}

// BeforeSave is a hook implement gorm model before exec create
func (i *Instance) BeforeSave() error {
	err := i.encryptPassword()
	if err != nil {
		return err
	}
	return i.marshalMycatConfig()
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db
func (i *Instance) AfterFind() error {
	err := i.decryptPassword()
	if err != nil {
		log.NewEntry().Errorf("decrypt password for instance %d failed, error: %v", i.ID, err)
	}
	err = i.unmarshalMycatConfig()
	if err != nil {
		log.NewEntry().Errorf("unmarshal mycat config for instance %d failed, error: %v", i.ID, err)
	}
	return nil
}

func (i *Instance) decryptPassword() error {
	if i == nil {
		return nil
	}
	if i.Password == "" {
		data, err := util.AesDecrypt(i.SecretPassword)
		if err != nil {
			return err
		} else {
			i.Password = string(data)
		}
	}
	return nil
}

func (i *Instance) encryptPassword() error {
	if i == nil {
		return nil
	}
	if i.SecretPassword == "" {
		data, err := util.AesEncrypt(i.Password)
		if err != nil {
			return err
		}
		i.SecretPassword = string(data)
	}
	return nil
}

func (i *Instance) unmarshalMycatConfig() error {
	if i == nil {
		return nil
	}
	if i.MycatConfigJson == "" {
		return nil
	}
	if i.MycatConfig == nil {
		i.MycatConfig = &MycatConfig{}
	}
	err := json.Unmarshal([]byte(i.MycatConfigJson), i.MycatConfig)
	if err != nil {
		return err
	}
	for _, dataHost := range i.MycatConfig.DataHosts {
		password, err := util.AesDecrypt(string(dataHost.Password))
		if err != nil {
			return err
		}
		dataHost.Password = util.Password(password)
	}
	return nil
}

func (i *Instance) marshalMycatConfig() error {
	if i == nil {
		return nil
	}
	if i.MycatConfig == nil {
		return nil
	}
	data, err := json.Marshal(i.MycatConfig)
	if err != nil {
		return err
	}
	i.MycatConfigJson = string(data)
	return nil
}

// InstanceDetail use for http request and swagger docs;
// it is same as Instance, but display RuleTemplates in json format.
type InstanceDetail struct {
	Instance
	RuleTemplates []RuleTemplate `json:"rule_template_list"`
	MycatConfig   *MycatConfig   `json:"mycat_config,omitempty"`
}

func (i *Instance) Detail() InstanceDetail {
	data := InstanceDetail{
		Instance:      *i,
		RuleTemplates: i.RuleTemplates,
		MycatConfig:   i.MycatConfig,
	}
	if i.RuleTemplates == nil {
		data.RuleTemplates = []RuleTemplate{}
	}
	return data
}

func (s *Storage) GetInstById(id string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Preload("RuleTemplates").Where("id = ?", id).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetInstByName(name string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Preload("RuleTemplates").Where("name = ?", name).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetInstancesByNames(names []string) ([]*Instance, error) {
	instances := []*Instance{}
	err := s.db.Where("name in (?)", names).Find(&instances).Error
	return instances, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateInstanceById(InstanceId string, attrs ...interface{}) error {
	err := s.db.Table("instances").Where("id = ?", InstanceId).Update(attrs...).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetInstances() ([]Instance, error) {
	instances := []Instance{}
	err := s.db.Find(&instances).Error
	return instances, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateInstRuleTemplate(inst *Instance, ts ...RuleTemplate) error {
	err := s.db.Model(inst).Association("RuleTemplates").Replace(ts).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
