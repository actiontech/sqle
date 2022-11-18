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
// NOTE: related model:
// - ProjectMemberRole, ProjectMemberGroupRole
type Instance struct {
	Model
	ProjectId uint `gorm:"index; not null"`
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

func (i *Instance) Fingerprint() string {
	return fmt.Sprintf(`
{
    "id":"%v",
    "host":"%v",
    "port":"%v",
    "user":"%v",
    "secret_password":"%v",
    "additional_params":"%v"
}
`, i.ID, i.Host, i.Port, i.User, i.SecretPassword, i.AdditionalParams)
}

func (s *Storage) GetInstanceById(id string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Preload("RuleTemplates").Where("id = ?", id).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstancesByIds(ids []uint) (instances []*Instance, err error) {
	distinctIds := utils.RemoveDuplicateUint(ids)
	return instances, s.db.Model(&Instance{}).Where("id in (?)", distinctIds).Scan(&instances).Error
}

func (s *Storage) GetAllInstance() ([]*Instance, error) {
	i := []*Instance{}
	err := s.db.Preload("RuleTemplates").Find(&i).Error
	return i, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstanceByNameAndProjectName(instName, projectName string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Joins("JOIN projects on projects.id = instances.project_id").Where("projects.name = ?", projectName).Where("instances.name = ?", instName).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstanceByNameAndProjectID(instName string, projectID uint) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Where("instances.project_id = ?", projectID).Where("instances.name = ?", instName).First(instance).Error
	if err == gorm.ErrRecordNotFound {
		return instance, false, nil
	}
	return instance, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstanceDetailByNameAndProjectName(instName string, projectName string) (*Instance, bool, error) {
	instance := &Instance{}
	err := s.db.Preload("WorkflowTemplate").Preload("RuleTemplates").
		Joins("JOIN projects on projects.id = instances.project_id").
		Where("instances.name = ?", instName).Where("projects.name = ?", projectName).First(instance).Error
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

func (s *Storage) GetInstancesByNamesAndProjectName(instNames []string, projectName string) ([]*Instance, error) {
	instances := []*Instance{}
	err := s.db.Joins("JOIN projects on projects.id = instances.project_id").
		Where("projects.name = ?", projectName).
		Where("instances.name in (?)", instNames).
		Find(&instances).Error
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

func (s *Storage) GetAndCheckInstanceExist(instanceNames []string, projectName string) (instances []*Instance, err error) {
	instances, err = s.GetInstancesByNamesAndProjectName(instanceNames, projectName)
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

func getDeduplicatedInstanceIds(instances []*Instance) []uint {
	instanceIds := make([]uint, len(instances))
	for i, inst := range instances {
		instanceIds[i] = inst.ID
	}
	instanceIds = utils.RemoveDuplicateUint(instanceIds)
	return instanceIds
}

func (s *Storage) GetInstanceTipsByTypeAndTempID(dbType string, tempID uint32, projectName string) (instances []*Instance, err error) {
	query := s.db.Model(&Instance{}).Select("instances.name, instances.db_type, instances.workflow_template_id").Group("instances.id")

	if dbType != "" {
		query = query.Where("db_type = ?", dbType)
	}

	if tempID != 0 {
		query = query.Where("workflow_template_id = ?", tempID)
	}

	if projectName != "" {
		query = query.Joins("LEFT JOIN projects ON projects.id = instances.project_id").Where("projects.name = ?", projectName)
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

func (s *Storage) GetInstanceTipsByUser(user *User, dbType string, projectName string) (
	instances []*Instance, err error) {

	isMember, err := s.IsProjectMember(user.Name, projectName)
	if err != nil {
		return nil, err
	}
	if IsDefaultAdminUser(user.Name) || isMember {
		return s.GetInstanceTipsByTypeAndTempID(dbType, 0, projectName)
	}

	return []*Instance{}, nil
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

func (s *Storage) CheckInstancesExist(projectName string, instNames []string) (bool, error) {
	instNames = utils.RemoveDuplicate(instNames)

	var count int
	err := s.db.Model(&Instance{}).
		Joins("JOIN projects ON projects.id = instances.project_id").
		Where("projects.name = ?", projectName).
		Where("instances.name in (?)", instNames).
		Count(&count).Error
	return len(instNames) == count, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) getInstanceBindCacheByNames(instNames []string, projectName string) (map[string /*inst name*/ ]uint /*inst id*/, error) {
	instNames = utils.RemoveDuplicate(instNames)

	insts, err := s.GetInstancesByNamesAndProjectName(instNames, projectName)
	if err != nil {
		return nil, err
	}

	if len(insts) != len(instNames) {
		return nil, errors.NewDataNotExistErr("some instances don't exist")
	}

	instCache := map[string /*inst name*/ ]uint /*inst id*/ {}

	for _, inst := range insts {
		instCache[inst.Name] = inst.ID
	}

	return instCache, nil
}
