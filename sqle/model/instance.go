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

func (s *Storage) CheckUserHasOpToInstance(user *User, instance *Instance, ops []uint) (bool, error) {
	query := `
SELECT instances.id
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
LEFT JOIN role_operations ON role_operations.role_id = roles.id
LEFT JOIN user_role ON user_role.role_id = roles.id
LEFT JOIN users ON users.id = user_role.user_id AND users.stat = 0
WHERE
instances.deleted_at IS NULL
AND instances.id = ?
AND users.id = ?
AND role_operations.op_code IN (?)
GROUP BY instances.id
UNION
SELECT instances.id
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
LEFT JOIN role_operations ON role_operations.role_id = roles.id
JOIN user_group_roles ON roles.id = user_group_roles.role_id
JOIN user_groups ON user_groups.id = user_group_roles.user_group_id AND user_groups.deleted_at IS NULL
JOIN user_group_users ON user_groups.id = user_group_users.user_group_id 
JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
WHERE 
instances.deleted_at IS NULL
AND instances.id = ?
AND users.id = ?
AND role_operations.op_code IN (?)
GROUP BY instances.id
`
	var instances []*Instance
	err := s.db.Raw(query, instance.ID, user.ID, ops, instance.ID, user.ID, ops).Scan(&instances).Error
	if err != nil {
		return false, errors.ConnectStorageErrWrapper(err)
	}
	return len(instances) > 0, nil
}

func (s *Storage) GetUserCanOpInstances(user *User, ops []uint) (instances []*Instance, err error) {
	query := `
SELECT instances.id
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
LEFT JOIN role_operations ON role_operations.role_id = roles.id
LEFT JOIN user_role ON user_role.role_id = roles.id
LEFT JOIN users ON users.id = user_role.user_id AND users.stat = 0
WHERE
instances.deleted_at IS NULL
AND users.id = ?
AND role_operations.op_code IN (?)
GROUP BY instances.id
UNION
SELECT instances.id
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
LEFT JOIN role_operations ON role_operations.role_id = roles.id
JOIN user_group_roles ON roles.id = user_group_roles.role_id
JOIN user_groups ON user_groups.id = user_group_roles.user_group_id AND user_groups.deleted_at IS NULL
JOIN user_group_users ON user_groups.id = user_group_users.user_group_id 
JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
WHERE 
instances.deleted_at IS NULL
AND users.id = ?
AND role_operations.op_code IN (?)
GROUP BY instances.id
`
	err = s.db.Raw(query, user.ID, ops, user.ID, ops).Scan(&instances).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}
	return
}

func getDbTypeQueryCond(dbType string) string {
	if dbType == "" {
		return ""
	}
	return `AND instances.db_type = ?`
}

func (s *Storage) GetAllInstanceTips(dbType string) (instances []*Instance, err error) {
	rawQuery := `
SELECT instances.name, instances.db_type
FROM instances
WHERE instances.deleted_at IS NULL
%s
GROUP BY instances.id
`

	query := fmt.Sprintf(rawQuery, getDbTypeQueryCond(dbType))
	if dbType == "" {
		err = s.db.Unscoped().Raw(query).Scan(&instances).Error
	} else {
		err = s.db.Unscoped().Raw(query, dbType).Scan(&instances).Error
	}
	return instances, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetAllInstanceCount() (int64, error) {
	var count int64
	return count, s.db.Model(&Instance{}).Count(&count).Error
}

func (s *Storage) GetInstanceTipsByUserViaRoles(
	user *User, dbType string) (instances []*Instance, err error) {

	rawQuery := `
SELECT instances.name, instances.db_type
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.stat = 0 AND roles.deleted_at IS NULL
LEFT JOIN user_role ON roles.id = user_role.role_id 
LEFT JOIN users ON users.id = user_role.user_id AND users.deleted_at IS NULL AND users.stat=0
WHERE instances.deleted_at IS NULL 
AND users.id = ?
%s
GROUP BY instances.id
UNION
SELECT instances.name, instances.db_type
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
JOIN user_group_roles ON roles.id = user_group_roles.role_id
JOIN user_groups ON user_groups.id = user_group_roles.user_group_id AND user_groups.stat = 0 AND user_groups.deleted_at IS NULL
JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
JOIN users ON users.id = user_group_users.user_id AND users.stat = 0 AND users.deleted_at IS NULL
WHERE instances.deleted_at IS NULL
%s
AND users.id = ?
GROUP BY instances.id
`

	dbTypeCond := getDbTypeQueryCond(dbType)

	query := fmt.Sprintf(rawQuery, dbTypeCond, dbTypeCond)

	if dbType == "" {
		err = s.db.Unscoped().Raw(query, user.ID, user.ID).Scan(&instances).Error
	} else {
		err = s.db.Unscoped().Raw(query, user.ID, dbType, user.ID, dbType).Scan(&instances).Error
	}

	return instances, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetInstanceTipsByUser(user *User, dbType string) (
	instances []*Instance, err error) {

	if IsDefaultAdminUser(user.Name) {
		return s.GetAllInstanceTips(dbType)
	}

	return s.GetInstanceTipsByUserViaRoles(user, dbType)
}

func getInstanceIDsFromInstances(instances []*Instance) (ids []uint) {
	ids = make([]uint, len(instances))
	for i := range instances {
		ids[i] = instances[i].ID
	}
	return ids
}

//SELECT instances.id
//FROM instances
//LEFT JOIN instance_role ON instance_role.instance_id = instances.id
//LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
//LEFT JOIN role_operations ON role_operations.role_id = roles.id
//LEFT JOIN user_role ON user_role.role_id = roles.id
//LEFT JOIN users ON users.id = user_role.user_id AND users.stat = 0
//WHERE
//instances.deleted_at IS NULL
//AND users.id = 5
//AND role_operations.op_code IN (20200)
//GROUP BY instances.id
//UNION
//SELECT instances.id
//FROM instances
//LEFT JOIN instance_role ON instance_role.instance_id = instances.id
//LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
//LEFT JOIN role_operations ON role_operations.role_id = roles.id
//JOIN user_group_roles ON roles.id = user_group_roles.role_id
//JOIN user_groups ON user_groups.id = user_group_roles.user_group_id AND user_groups.deleted_at IS NULL
//JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
//JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
//WHERE
//instances.deleted_at IS NULL
//AND users.id = 5
//AND role_operations.op_code IN (20200)
//GROUP BY instances.id

//SELECT instances.id
//FROM instances
//LEFT JOIN instance_role ON instance_role.instance_id = instances.id
//LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
//LEFT JOIN role_operations ON role_operations.role_id = roles.id
//LEFT JOIN user_role ON user_role.role_id = roles.id
//LEFT JOIN users ON users.id = user_role.user_id AND users.stat = 0
//WHERE
//instances.deleted_at IS NULL
//AND instances.id = 5
//AND users.id = 4
//AND role_operations.op_code IN (20200)
//GROUP BY instances.id
//UNION
//SELECT instances.id
//FROM instances
//LEFT JOIN instance_role ON instance_role.instance_id = instances.id
//LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
//LEFT JOIN role_operations ON role_operations.role_id = roles.id
//JOIN user_group_roles ON roles.id = user_group_roles.role_id
//JOIN user_groups ON user_groups.id = user_group_roles.user_group_id AND user_groups.deleted_at IS NULL
//JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
//JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
//WHERE
//instances.deleted_at IS NULL
//AND instances.id = 5
//AND users.id = 4
//AND role_operations.op_code IN (20200)
//GROUP BY instances.id
