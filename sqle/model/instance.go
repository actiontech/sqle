package model

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/jinzhu/gorm"
)

// Instance is a table for database info
type Instance struct {
	Model
	// has created composite index: [id, name] by gorm#AddIndex
	Name               string         `json:"name" gorm:"not null;index" example:""`
	DbType             string         `json:"db_type" gorm:"column:db_type; not null" example:"mysql"`
	Host               string         `json:"host" gorm:"column:db_host; not null" example:"10.10.10.10"`
	Port               string         `json:"port" gorm:"column:db_port; not null" example:"3306"`
	User               string         `json:"user" gorm:"column:db_user; not null" example:"root"`
	Password           string         `json:"-" gorm:"-"`
	SecretPassword     string         `json:"secret_password" gorm:"column:db_password; not null"`
	Desc               string         `json:"desc" example:"this is a instance"`
	WorkflowTemplateId uint           `json:"workflow_template_id"`
	AdditionalParams   params.Params  `json:"additional_params" gorm:"type:text"`
	MaintenancePeriod  Periods        `json:"maintenance_period" gorm:"type:text"`
	SqlQueryConfig     SqlQueryConfig `json:"sql_query_config" gorm:"type:varchar(255); default:'{\"max_pre_query_rows\":100,\"query_timeout_second\":10}'"`

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

func (s *Storage) GetInstancesByType(dbType string) ([]*Instance, error) {
	instances := []*Instance{}
	err := s.db.Where("db_type = ?", dbType).Find(&instances).Error
	return instances, errors.New(errors.ConnectStorageError, err)
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

var checkUserHasOpToInstancesQuery = `
SELECT instances.id
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
LEFT JOIN role_operations ON role_operations.role_id = roles.id
LEFT JOIN user_role ON user_role.role_id = roles.id
LEFT JOIN users ON users.id = user_role.user_id AND users.stat = 0
WHERE
instances.deleted_at IS NULL
AND instances.id IN (?)
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
AND instances.id IN (?)
AND users.id = ?
AND role_operations.op_code IN (?)
GROUP BY instances.id
`

func getDeduplicatedInstanceIds(instances []*Instance) []uint {
	instanceIds := make([]uint, len(instances))
	for i, inst := range instances {
		instanceIds[i] = inst.ID
	}
	instanceIds = utils.RemoveDuplicateUint(instanceIds)
	return instanceIds
}

func (s *Storage) CheckUserHasOpToInstances(user *User, instances []*Instance, ops []uint) (bool, error) {
	instanceIds := getDeduplicatedInstanceIds(instances)
	var instanceRecords []*Instance
	err := s.db.Raw(checkUserHasOpToInstancesQuery, instanceIds, user.ID, ops, instanceIds, user.ID, ops).Scan(&instanceRecords).Error
	if err != nil {
		return false, errors.ConnectStorageErrWrapper(err)
	}
	return len(instanceRecords) == len(instanceIds), nil
}

func (s *Storage) CheckUserHasOpToAnyInstance(user *User, instances []*Instance, ops []uint) (bool, error) {
	instanceIds := getDeduplicatedInstanceIds(instances)
	var instanceRecords []*Instance
	err := s.db.Raw(checkUserHasOpToInstancesQuery, instanceIds, user.ID, ops, instanceIds, user.ID, ops).Scan(&instanceRecords).Error
	if err != nil {
		return false, errors.ConnectStorageErrWrapper(err)
	}
	return len(instanceRecords) > 0, nil
}

func (s *Storage) GetUserCanOpInstances(user *User, ops []uint) (instances []*Instance, err error) {
	query := `
SELECT instances.id, instances.name
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
SELECT instances.id, instances.name
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

func (s *Storage) GetInstancesTipsByUserAndTempId(user *User, dbType string, tempID uint32) ([]*Instance, error) {
	if IsDefaultAdminUser(user.Name) {
		return s.GetInstanceTipsByTempID(dbType, tempID)
	}

	return s.GetInstanceTipsByUserViaRolesAndTempID(user, dbType, tempID)
}

func (s *Storage) GetInstanceTipsByTempID(dbType string, tempID uint32) (instances []*Instance, err error) {
	query := s.db.Model(&Instance{}).Select("name, db_type, workflow_template_id").Group("id")

	if dbType != "" {
		query = query.Where("db_type = ?", dbType)
	}

	if tempID != 0 {
		query = query.Where("workflow_template_id = ?", tempID)
	}

	return instances, query.Scan(&instances).Error
}

func (s *Storage) GetAllInstanceCountByType(dbTypes ...string) (int64, error) {
	var count int64
	return count, s.db.Model(&Instance{}).Where("db_type in (?)", dbTypes).Count(&count).Error
}

type TypeCount struct {
	DBType string `json:"db_type"`
	Count  int64  `json:"count"`
}

func (s *Storage) GetAllInstanceCount() ([]*TypeCount, error) {
	var counts []*TypeCount
	return counts, s.db.Table("instances").Select("db_type, count(*) as count").Where("deleted_at is NULL").Group("db_type").Find(&counts).Error
}

func (s *Storage) GetInstanceTipsByUserViaRolesAndTempID(user *User, dbType string, tempID uint32) (instances []*Instance, err error) {

	queryByRole := s.db.Model(&Instance{}).Select("instances.name, instances.db_type , instances.workflow_template_id").
		Joins("LEFT JOIN instance_role ON instance_role.instance_id = instances.id").
		Joins("LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.stat = 0 AND roles.deleted_at IS NULL").
		Joins("LEFT JOIN user_role ON roles.id = user_role.role_id").
		Joins("LEFT JOIN users ON users.id = user_role.user_id AND users.deleted_at IS NULL AND users.stat=0").
		Where("users.id = ?", user.ID).
		Group("instances.id")

	queryByGroup := s.db.Model(&Instance{}).Select("instances.name, instances.db_type , instances.workflow_template_id").
		Joins("LEFT JOIN instance_role ON instance_role.instance_id = instances.id").
		Joins("LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0").
		Joins("JOIN user_group_roles ON roles.id = user_group_roles.role_id").
		Joins("JOIN user_groups ON user_groups.id = user_group_roles.user_group_id AND user_groups.stat = 0 AND user_groups.deleted_at IS NULL").
		Joins("JOIN user_group_users ON user_groups.id = user_group_users.user_group_id").
		Joins("JOIN users ON users.id = user_group_users.user_id AND users.stat = 0 AND users.deleted_at IS NULL").
		Where("users.id = ?", user.ID).
		Group("instances.id")

	if dbType != "" {
		queryByRole = queryByRole.Where("db_type = ?", dbType)
		queryByGroup = queryByGroup.Where("db_type = ?", dbType)
	}

	if tempID != 0 {
		queryByRole = queryByRole.Where("workflow_template_id = ?", tempID)
		queryByGroup = queryByGroup.Where("workflow_template_id = ?", tempID)
	}

	var intsByRole, instByGroup []*Instance
	if err := queryByRole.Scan(&intsByRole).Error; err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	if err := queryByGroup.Scan(&instByGroup).Error; err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	instances = append(instances, intsByRole...)
	instances = append(instances, instByGroup...)

	return instances, nil
}

func (s *Storage) GetInstanceTipsByUser(user *User, dbType string) (
	instances []*Instance, err error) {

	if IsDefaultAdminUser(user.Name) {
		return s.GetInstanceTipsByTempID(dbType, 0)
	}

	return s.GetInstanceTipsByUserViaRolesAndTempID(user, dbType, 0)
}

func (s *Storage) GetInstanceTipsByUserAndOperation(user *User, dbType string, opCode ...int) (
	instances []*Instance, err error) {

	if IsDefaultAdminUser(user.Name) {
		return s.GetInstanceTipsByTempID(dbType, 0)
	}
	return s.getInstanceTipsByUserAndOperation(user, dbType, opCode...)
}

func (s *Storage) getInstanceTipsByUserAndOperation(user *User, dbType string, opCode ...int) (
	instances []*Instance, err error) {
	query := `
SELECT instances.name, instances.db_type
FROM instances
LEFT JOIN instance_role ON instance_role.instance_id = instances.id
LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
LEFT JOIN role_operations ON role_operations.role_id = roles.id
LEFT JOIN user_role ON user_role.role_id = roles.id
LEFT JOIN users ON users.id = user_role.user_id AND users.stat = 0
WHERE
instances.deleted_at IS NULL
%s
AND users.id = ?
AND role_operations.op_code IN (?)
GROUP BY instances.id
UNION
SELECT instances.name, instances.db_type
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
%s
AND users.id = ?
AND role_operations.op_code IN (?)
GROUP BY instances.id
`
	dbTypeCond := getDbTypeQueryCond(dbType)
	query = fmt.Sprintf(query, dbTypeCond, dbTypeCond)
	if dbType == "" {
		err = s.db.Raw(query, user.ID, opCode, user.ID, opCode).Scan(&instances).Error
	} else {
		err = s.db.Raw(query, dbType, user.ID, opCode, dbType, user.ID, opCode).Scan(&instances).Error
	}
	return instances, errors.ConnectStorageErrWrapper(err)
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
