package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"actiontech.cloud/universe/ucommon/v4/util"

	"actiontech.cloud/sqle/sqle/sqle/errors"
	"actiontech.cloud/sqle/sqle/sqle/log"
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
	// has created composite index: [id, name] by gorm#AddIndex
	Name               string       `json:"name" gorm:"not null;index" example:""`
	DbType             string       `json:"db_type" gorm:"column:db_type; not null" example:"mysql"`
	Host               string       `json:"host" gorm:"column:db_host; not null" example:"10.10.10.10"`
	Port               string       `json:"port" gorm:"column:db_port; not null" example:"3306"`
	User               string       `json:"user" gorm:"column:db_user; not null" example:"root"`
	Password           string       `json:"-" gorm:"-"`
	SecretPassword     string       `json:"secret_password" gorm:"column:db_password; not null"`
	Desc               string       `json:"desc" example:"this is a instance"`
	WorkflowTemplateId uint         `json:"workflow_template_id"`
	MycatConfig        *MycatConfig `json:"-" gorm:"-"`
	MycatConfigJson    string       `json:"-" gorm:"type:text;column:mycat_config"`

	// relation table
	Roles            []*Role           `json:"-" gorm:"many2many:instance_role;"`
	RuleTemplates    []RuleTemplate    `json:"-" gorm:"many2many:instance_rule_template"`
	WorkflowTemplate *WorkflowTemplate `gorm:"foreignkey:WorkflowTemplateId"`
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

func (s *Storage) GetInstanceById(id string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Preload("RuleTemplates").Where("id = ?", id).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetInstanceByName(name string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Where("name = ?", name).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetInstanceDetailByName(name string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Preload("Roles").Preload("WorkflowTemplate").Preload("RuleTemplates").
		Where("name = ?", name).First(instance).Error
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

func (s *Storage) UpdateInstanceById(InstanceId uint, attrs ...interface{}) error {
	err := s.db.Table("instances").Where("id = ?", InstanceId).Update(attrs...).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) CheckInstanceBindCount(ruleTemplates []string, instances ...*Instance) error {
	//需求: 限制一个数据源只能绑定一个规则模版
	//创建/更新 数据源：	保证len(ruleTemplates)=1即可
	//创建/更新 规则模版：1.数据源已绑定模版数量=0
	//					2.数据源已绑定模版名==更新的模版名
	//					3.其余情况 return err
	if len(ruleTemplates) > 1 {
		return errors.New(errors.DataExist, fmt.Errorf("an instance can only bind one rule template"))
	}

	for _, inst := range instances {
		var associationRT []RuleTemplate
		count := s.db.Model(inst).Association("RuleTemplates").Find(&associationRT).Count()
		if count > 1 {
			return errors.New(errors.DataExist, fmt.Errorf("an instance can only bind one rule template"))
		}
		if count == 1 && associationRT[0].Name != ruleTemplates[0] {
			return errors.New(errors.DataExist, fmt.Errorf("an instance can only bind one rule template"))
		}
	}
	return nil
}

func (s *Storage) UpdateInstanceRuleTemplates(instance *Instance, ts ...*RuleTemplate) error {
	err := s.db.Model(instance).Association("RuleTemplates").Replace(ts).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateInstanceRoles(instance *Instance, rs ...*Role) error {
	err := s.db.Model(instance).Association("Roles").Replace(rs).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetUserInstanceTip(user *User) ([]*Instance, error) {
	instances := []*Instance{}
	db := s.db.Model(&Instance{}).Select("instances.name")
	if user.Name != DefaultAdminUser {
		db = db.Joins("JOIN instance_role AS ir ON instances.id = ir.instance_id").
			Joins("JOIN user_role AS ur ON ir.role_id = ur.role_id").
			Joins("JOIN users ON ur.user_id = users.id AND users.id = ?", user.ID)
	}
	err := db.Scan(&instances).Error
	return instances, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetAndCheckInstanceExist(instanceNames []string) (instances []*Instance, err error) {
	instances, err = s.GetInstancesByNames(instanceNames)
	if err != nil {
		return instances, err
	}
	existInstanceNames := map[string]struct{}{}
	for _, instance := range instances {
		existInstanceNames[instance.Name] = struct{}{}
	}
	notExistInstanceNames := []string{}
	for _, instanceName := range instanceNames {
		if _, ok := existInstanceNames[instanceName]; !ok {
			notExistInstanceNames = append(notExistInstanceNames, instanceName)
		}
	}
	if len(notExistInstanceNames) > 0 {
		return instances, errors.New(errors.DataNotExist,
			fmt.Errorf("instance %s not exist", strings.Join(notExistInstanceNames, ", ")))
	}
	return instances, nil
}

func (s *Storage) GetInstanceNamesByWorkflowTemplateId(id uint) ([]string, error) {
	var instances []*Instance
	err := s.db.Select("name").Where("workflow_template_id = ?", id).Find(&instances).Error
	if err != nil {
		return []string{}, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	names := make([]string, 0, len(instances))
	for _, instance := range instances {
		names = append(names, instance.Name)
	}
	return names, nil
}
