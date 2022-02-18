package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/jinzhu/gorm"
)

// Instance is a table for database info
type Instance struct {
	Model
	// has created composite index: [id, name] by gorm#AddIndex
	Name               string `json:"name" gorm:"not null;index" example:""`
	DbType             string `json:"db_type" gorm:"column:db_type; not null" example:"mysql"`
	Host               string `json:"host" gorm:"column:db_host; not null" example:"10.10.10.10"`
	Port               string `json:"port" gorm:"column:db_port; not null" example:"3306"`
	User               string `json:"user" gorm:"column:db_user; not null" example:"root"`
	Password           string `json:"-" gorm:"-"`
	SecretPassword     string `json:"secret_password" gorm:"column:db_password; not null"`
	Desc               string `json:"desc" example:"this is a instance"`
	WorkflowTemplateId uint   `json:"workflow_template_id"`

	// relation table
	Roles            []*Role           `json:"-" gorm:"many2many:instance_role;"`
	RuleTemplates    []RuleTemplate    `json:"-" gorm:"many2many:instance_rule_template"`
	WorkflowTemplate *WorkflowTemplate `gorm:"foreignkey:WorkflowTemplateId"`
}

// BeforeSave is a hook implement gorm model before exec create
func (i *Instance) BeforeSave() error {
	return i.encryptPassword()
}

// AfterFind is a hook implement gorm model after query, ignore err if query from db
func (i *Instance) AfterFind() error {
	err := i.decryptPassword()
	if err != nil {
		log.NewEntry().Errorf("decrypt password for instance %d failed, error: %v", i.ID, err)
	}
	return nil
}

func (i *Instance) decryptPassword() error {
	if i == nil {
		return nil
	}
	if i.Password == "" {
		data, err := utils.AesDecrypt(i.SecretPassword)
		if err != nil {
			return err
		} else {
			i.Password = data
		}
	}
	return nil
}

func (i *Instance) encryptPassword() error {
	if i == nil {
		return nil
	}
	if i.SecretPassword == "" {
		data, err := utils.AesEncrypt(i.Password)
		if err != nil {
			return err
		}
		i.SecretPassword = data
	}
	return nil
}

func (s *Storage) GetInstanceById(id string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Preload("RuleTemplates").Where("id = ?", id).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstanceByName(name string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Where("name = ?", name).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstanceDetailByName(name string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Preload("Roles").Preload("WorkflowTemplate").Preload("RuleTemplates").
		Where("name = ?", name).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstancesByNames(names []string) ([]*Instance, error) {
	instances := []*Instance{}
	err := s.db.Where("name in (?)", names).Find(&instances).Error
	return instances, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateInstanceById(InstanceId uint, attrs ...interface{}) error {
	err := s.db.Table("instances").Where("id = ?", InstanceId).Update(attrs...).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateInstanceRuleTemplates(instance *Instance, ts ...*RuleTemplate) error {
	err := s.db.Model(instance).Association("RuleTemplates").Replace(ts).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) UpdateInstanceRoles(instance *Instance, rs ...*Role) error {
	err := s.db.Model(instance).Association("Roles").Replace(rs).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetUserInstanceTip(user *User, dbType string) (
	instances []*Instance, err error) {

	db := s.db.Model(&Instance{}).Select("instances.name, instances.db_type")
	if user.Name != DefaultAdminUser {
		// 1. get roles
		roles, err := s.GetRolesByUserID(int(user.ID))
		if err != nil {
			return nil, err
		}
		if len(roles) == 0 {
			return instances, nil
		}
		roleIDs := GetRoleIDsFromRoles(roles)

		// 2. get instances by roleIDs
		db = db.Joins("JOIN instance_role AS ir ON instances.id = ir.instance_id").
			Joins("JOIN roles ON ir.role_id = roles.id AND roles.deleted_at IS NULL").
			Where("roles.id IN (?)", roleIDs)
	}
	if dbType != "" {
		db = db.Where("instances.db_type = ?", dbType)
	}
	err = db.Scan(&instances).Error
	return instances, errors.New(errors.ConnectStorageError, err)
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
		return []string{}, errors.New(errors.ConnectStorageError, err)
	}
	names := make([]string, 0, len(instances))
	for _, instance := range instances {
		names = append(names, instance.Name)
	}
	return names, nil
}

func (s *Storage) GetInstancesByRoleIDs(roleIDs []uint) (insts []*Instance, err error) {

	if len(roleIDs) == 0 {
		return insts, nil
	}

	err = s.db.Unscoped().Model(&Instance{}).
		Joins("JOIN instance_role AS ir ON instances.id = ir.instance_id").
		Joins("JOIN roles ON ir.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0").
		Where("roles.id IN (?)", roleIDs).Group("instances.id").
		Find(&insts).Error
	if err != nil {
		return insts, errors.ConnectStorageErrWrapper(err)
	}

	return insts, nil
}

func GetInstanceIDsFromInst(insts []*Instance) (instIDs []uint) {
	instIDs = make([]uint, len(insts))
	for i := range insts {
		instIDs[i] = insts[i].ID
	}
	return
}
