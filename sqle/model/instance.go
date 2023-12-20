package model

import (
	"fmt"

	dmsCommonAes "github.com/actiontech/dms/pkg/dms-common/pkg/aes"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

const InstanceSourceSQLE string = "SQLE"

// Instance is a table for database info
// NOTE: related model:
// - ProjectMemberRole, ProjectMemberGroupRole
type Instance struct {
	ID               uint64 `json:"id"`
	RuleTemplateId   uint64 `json:"rule_template_id"`
	RuleTemplateName string `json:"rule_template_name"`
	ProjectId        string `gorm:"index; not null"`
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
	Source             string         `json:"source" gorm:"not null"`
	SyncInstanceTaskID uint           `json:"sync_instance_task_id"`

	// relation table
	RuleTemplates    []RuleTemplate    `json:"-" gorm:"many2many:instance_rule_template"`
	WorkflowTemplate *WorkflowTemplate `gorm:"foreignkey:WorkflowTemplateId"`
}

func (i *Instance) GetIDStr() string {
	return fmt.Sprintf("%d", i.ID)
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
		data, err := dmsCommonAes.AesDecrypt(i.SecretPassword)
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
		data, err := dmsCommonAes.AesEncrypt(i.Password)
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
    "password":"%v",
    "additional_params":"%v"
}
`, i.ID, i.Host, i.Port, i.User, i.Password, i.AdditionalParams)
}

// func (s *Storage) GetInstancesFromActiveProjectByIds(ids []uint) ([]*Instance, error) {
// 	instances := []*Instance{}
// 	err := s.db.Joins("LEFT JOIN projects ON projects.id = instances.project_id").
// 		Where("instances.id IN (?)", ids).
// 		Where("projects.status = ?", ProjectStatusActive).
// 		Find(&instances).Error
// 	return instances, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) GetAllInstance() ([]*Instance, error) {
// 	i := []*Instance{}
// 	err := s.db.Preload("RuleTemplates").Find(&i).Error
// 	return i, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) GetInstanceByNameAndProjectName(instName, projectName string) (*Instance, bool, error) {
// 	instance := &Instance{}
// 	err := s.db.Joins("JOIN projects on projects.id = instances.project_id").Where("projects.name = ?", projectName).Where("instances.name = ?", instName).First(instance).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return instance, false, nil
// 	}
// 	return instance, true, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) GetInstances(instIds []string) ([]*Instance, bool, error) {
// 	instances := []*Instance{}
// 	err := s.db.Where("instances.id in (?)", instIds).Find(&instances).Error
// 	if err == gorm.ErrRecordNotFound {
// 		return instances, false, nil
// 	}
// 	return instances, true, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) DeleteInstance(instance *Instance) error {
// 	err := s.DeleteRoleByInstanceID(instance.ID)
// 	if err != nil {
// 		return errors.ConnectStorageErrWrapper(err)
// 	}
// 	err = s.Delete(instance)
// 	if err != nil {
// 		return errors.ConnectStorageErrWrapper(err)
// 	}
// 	return nil
// }

// func (s *Storage) UpdateInstanceById(InstanceId uint, attrs ...interface{}) error {
// 	err := s.db.Table("instances").Where("id = ?", InstanceId).Update(attrs...).Error
// 	return errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) UpdateInstanceRoles(instance *Instance, rs ...*Role) error {
// 	err := s.db.Model(instance).Association("Roles").Replace(rs).Error
// 	return errors.New(errors.ConnectStorageError, err)
// }

type TypeCount struct {
	DBType string `json:"db_type"`
	Count  int64  `json:"count"`
}

// func (s *Storage) GetInstanceTipsByUser(user *User, dbType string, projectName string) (
// 	instances []*Instance, err error) {

// 	isMember, err := s.IsProjectMember(user.Name, projectName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if IsDefaultAdminUser(user.Name) || isMember {
// 		return s.GetInstanceTipsByTypeAndTempID(dbType, 0, projectName)
// 	}

// 	return []*Instance{}, nil
// }

/*func GetInstanceIDsFromInstances(instances []*Instance) (ids []uint) {
	ids = make([]uint, len(instances))
	for i := range instances {
		ids[i] = instances[i].ID
	}
	return ids
}*/

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

// type SyncTaskInstance struct {
// 	Instances            []*Instance
// 	RuleTemplate         *RuleTemplate
// 	NeedDeletedInstances []*Instance
// }

// func (s *Storage) BatchUpdateSyncTask(instTemplate *SyncTaskInstance) error {
// 	err := s.Tx(func(tx *gorm.DB) error {
// 		for _, instance := range instTemplate.Instances {
// 			if err := tx.Save(instance).Error; err != nil {
// 				return err
// 			}

// 			if err := s.UpdateInstanceRuleTemplates(instance, instTemplate.RuleTemplate); err != nil {
// 				return err
// 			}
// 		}

// 		for _, instance := range instTemplate.NeedDeletedInstances {
// 			if err := s.DeleteInstance(instance); err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return errors.ConnectStorageErrWrapper(err)
// 	}

// 	return nil
// }

// GetSqlExecutionFailCount 获取sql上线失败统计
func (s *Storage) GetSqlExecutionFailCount() ([]SqlExecutionCount, error) {
	var sqlExecutionFailCount []SqlExecutionCount

	err := s.db.Model(&Workflow{}).Select("t.instance_id, count(*) as count").
		Joins("left join workflow_records wr on workflows.workflow_record_id = wr.id").
		Joins("left join workflow_instance_records wir on wr.id = wir.workflow_record_id").
		Joins("left join tasks t on wir.task_id = t.id").
		Where("t.status = ?", TaskStatusExecuteFailed).
		Where("t.exec_start_at is not null").
		Where("t.exec_end_at is not null").
		Group("t.instance_id").
		Scan(&sqlExecutionFailCount).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	return sqlExecutionFailCount, nil
}

// GetSqlExecutionTotalCount 获取sql上线总数统计
// 上线总数(根据数据源划分)是指：正在上线,上线成功,上线失败 task的总数
func (s *Storage) GetSqlExecutionTotalCount() ([]SqlExecutionCount, error) {
	var sqlExecutionTotalCount []SqlExecutionCount

	err := s.db.Model(&Workflow{}).Select("t.instance_id, count(*) as count").
		Joins("left join workflow_records wr on workflows.workflow_record_id = wr.id").
		Joins("left join workflow_instance_records wir on wr.id = wir.workflow_record_id").
		Joins("left join tasks t on wir.task_id = t.id").
		Where("t.status not in (?)", []string{TaskStatusInit, TaskStatusAudited}).
		Where("t.exec_start_at is not null").
		Group("t.instance_id").
		Scan(&sqlExecutionTotalCount).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	return sqlExecutionTotalCount, nil
}

type InstanceWorkFlowStatusCount struct {
	DbType       string `json:"db_type"`
	InstanceName string `json:"instance_name"`
	StatusCount  uint   `json:"status_count"`
	InstanceId   uint64 `json:"instance_id"`
}

func (s *Storage) GetInstanceWorkFlowStatusCountByProject(instances []*Instance, queryStatus []string) ([]*InstanceWorkFlowStatusCount, error) {
	var instanceWorkFlowStatusCount []*InstanceWorkFlowStatusCount

	instanceIds := make([]uint64, 0)
	var instanceMap = make(map[uint64]InstanceWorkFlowStatusCount)
	for _, instance := range instances {
		instanceMap[instance.ID] = InstanceWorkFlowStatusCount{
			DbType:       instance.DbType,
			InstanceName: instance.Name,
		}
		instanceIds = append(instanceIds, instance.ID)
	}

	err := s.db.Model(&WorkflowInstanceRecord{}).
		Select("workflow_instance_records.instance_id, count(case when workflow_records.status in (?) then 1 else null end) status_count", queryStatus).
		Joins("join workflow_records on workflow_instance_records.workflow_record_id=workflow_records.id").
		Joins("join workflows on workflow_records.id=workflows.workflow_record_id").
		Where("workflow_instance_records.instance_id in (?)", instanceIds).
		Group("workflow_instance_records.instance_id").
		Scan(&instanceWorkFlowStatusCount).Error

	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	for i, item := range instanceWorkFlowStatusCount {
		if instance, ok := instanceMap[item.InstanceId]; ok {
			instanceWorkFlowStatusCount[i].DbType = instance.DbType
			instanceWorkFlowStatusCount[i].InstanceName = instance.InstanceName
		}
	}

	return instanceWorkFlowStatusCount, nil
}
