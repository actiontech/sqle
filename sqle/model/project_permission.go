package model

import (
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
)

/*

instance permission.

*/

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

func (s *Storage) GetInstanceTipsByUserAndOperation(user *User, dbType, projectName string, opCode ...int) (
	instances []*Instance, err error) {

	isProjectManager, err := s.IsProjectManager(user.Name, projectName)
	if err != nil {
		return nil, err
	}

	if IsDefaultAdminUser(user.Name) || isProjectManager {
		return s.GetInstanceTipsByTypeAndTempID(dbType, 0, projectName)
	}
	return s.getInstanceTipsByUserAndOperation(user, dbType, projectName, opCode...)
}

func (s *Storage) getInstanceTipsByUserAndOperation(user *User, dbType string, projectName string, opCode ...int) (
	instances []*Instance, err error) {

	query1 := s.db.Table("instances").
		Select("instances.name, instances.db_type").
		Joins("LEFT JOIN instance_role ON instance_role.instance_id = instances.id").
		Joins("LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0").
		Joins("LEFT JOIN role_operations ON role_operations.role_id = roles.id").
		Joins("LEFT JOIN user_role ON user_role.role_id = roles.id").
		Joins("LEFT JOIN users ON users.id = user_role.user_id AND users.stat = 0").
		Joins("LEFT JOIN projects ON projects.id = instances.project_id").
		Where("instances.deleted_at IS NULL").
		Where("users.id = ?", user.ID).
		Where("role_operations.op_code IN (?)", opCode).
		Group("instances.id")

	query2 := s.db.Table("instances").
		Select("instances.name, instances.db_type").
		Joins("LEFT JOIN instance_role ON instance_role.instance_id = instances.id").
		Joins("LEFT JOIN roles ON roles.id = instance_role.role_id AND roles.deleted_at IS NULL AND roles.stat = 0").
		Joins("LEFT JOIN role_operations ON role_operations.role_id = roles.id").
		Joins("JOIN user_group_roles ON roles.id = user_group_roles.role_id").
		Joins("JOIN user_groups ON user_groups.id = user_group_roles.user_group_id AND user_groups.deleted_at IS NULL").
		Joins("JOIN user_group_users ON user_groups.id = user_group_users.user_group_id").
		Joins("JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0").
		Joins("LEFT JOIN projects ON projects.id = instances.project_id").
		Where("instances.deleted_at IS NULL").
		Where("users.id = ?", user.ID).
		Where("role_operations.op_code IN (?)", opCode).
		Group("instances.id")

	if projectName != "" {
		query1.Where("projects.name = ?", projectName)
		query2.Where("projects.name = ?", projectName)
	}

	if dbType != "" {
		query1.Where("AND instances.db_type = ?", dbType)
		query2.Where("AND instances.db_type = ?", dbType)
	}

	err = s.db.Raw("? UNION ?", query1.QueryExpr(), query2.QueryExpr()).Scan(&instances).Error

	return instances, errors.ConnectStorageErrWrapper(err)
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

func (s *Storage) UserCanAccessInstance(user *User, instance *Instance) (
	ok bool, err error) {

	isManager, err := s.IsProjectManagerByID(user.ID, instance.ProjectId)
	if err != nil {
		return false, err
	}

	if IsDefaultAdminUser(user.Name) || isManager {
		return true, nil
	}

	type countStruct struct {
		Count int `json:"count"`
	}

	query := `
SELECT COUNT(1) AS count
FROM instances
LEFT JOIN project_member_roles ON project_member_roles.instance_id = instances.id
LEFT JOIN users ON users.id = project_member_roles.user_id
WHERE instances.deleted_at IS NULL
AND users.stat = 0 
AND users.deleted_at IS NULL
AND instances.id = ?
AND users.id = ?
GROUP BY instances.id
UNION
SELECT instances.id
FROM instances
LEFT JOIN project_member_group_roles ON project_member_group_roles.instance_id = instances.id
JOIN user_group_users ON project_member_group_roles.user_group_id = user_group_users.user_group_id
JOIN users ON users.id = user_group_users.user_id
WHERE instances.deleted_at IS NULL
AND users.stat = 0 
AND users.deleted_at IS NULL
AND instances.id = ?
AND users.id = ?
GROUP BY instances.id
`
	var cnt countStruct
	err = s.db.Unscoped().Raw(query, instance.ID, user.ID, instance.ID, user.ID).Scan(&cnt).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return cnt.Count > 0, nil
}

/*

workflow permission.

*/

func (s *Storage) UserCanAccessWorkflow(user *User, workflow *Workflow) (bool, error) {
	query := `SELECT count(w.id) FROM workflows AS w
JOIN workflow_records AS wr ON w.workflow_record_id = wr.id AND w.id = ?
LEFT JOIN workflow_steps AS cur_ws ON wr.current_workflow_step_id = cur_ws.id
LEFT JOIN workflow_step_templates AS cur_wst ON cur_ws.workflow_step_template_id = cur_wst.id
LEFT JOIN workflow_step_user AS cur_wst_re_user ON cur_ws.id = cur_wst_re_user.workflow_step_id
LEFT JOIN users AS cur_ass_user ON cur_wst_re_user.user_id = cur_ass_user.id AND cur_ass_user.stat=0
LEFT JOIN workflow_steps AS op_ws ON w.id = op_ws.workflow_id AND op_ws.state != "initialized"
LEFT JOIN workflow_step_templates AS op_wst ON op_ws.workflow_step_template_id = op_wst.id
LEFT JOIN workflow_step_user AS op_wst_re_user ON op_ws.id = op_wst_re_user.workflow_step_id
LEFT JOIN users AS op_ass_user ON op_wst_re_user.user_id = op_ass_user.id AND op_ass_user.stat=0
where w.deleted_at IS NULL
AND (w.create_user_id = ? OR cur_ass_user.id = ? OR op_ass_user.id = ?)
`
	var count uint
	err := s.db.Raw(query, workflow.ID, user.ID, user.ID, user.ID).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return count > 0, nil
}

// GetUsersByOperationCode will return admin user if no qualified user is found, preventing the process from being stuck because no user can operate
func (s *Storage) GetUsersByOperationCode(instance *Instance, opCode ...int) (users []*User, err error) {
	names := []string{}
	err = s.db.Model(&User{}).Select("DISTINCT users.login_name").
		Joins("LEFT JOIN user_role ON users.id = user_role.user_id "+
			"LEFT JOIN user_group_users ON users.id = user_group_users.user_id "+
			"LEFT JOIN user_group_roles ON user_group_users.user_group_id = user_group_roles.role_id "+
			"LEFT JOIN role_operations ON ( user_role.role_id = role_operations.role_id OR user_group_roles.role_id = role_operations.role_id ) AND role_operations.deleted_at IS NULL "+
			"LEFT JOIN instance_role ON instance_role.role_id = role_operations.role_id ").
		Where("instance_role.instance_id = ?", instance.ID).
		Where("role_operations.op_code in (?)", opCode).
		Group("users.id").
		Pluck("login_name", &names).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}
	if len(names) == 0 {
		names = append(names, DefaultAdminUser)
	}
	return s.GetUsersByNames(names)
}

/*

audit plan permission.

*/

func (s *Storage) CheckUserCanCreateAuditPlan(user *User, instName, dbType string) (bool, error) {
	if user.Name == DefaultAdminUser {
		return true, nil
	}
	instances, err := s.GetUserCanOpInstances(user, []uint{OP_AUDIT_PLAN_SAVE})
	if err != nil {
		return false, err
	}
	for _, instance := range instances {
		if instName == instance.Name {
			return true, nil
		}
	}
	return false, nil
}

/*

project permission.

*/

func (s *Storage) CheckUserCanUpdateProject(projectName string, userID uint) (bool, error) {
	user, exist, err := s.GetUserByID(userID)
	if err != nil || !exist {
		return false, err
	}

	if user.Name == DefaultAdminUser {
		return true, nil
	}

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil || !exist {
		return false, err
	}

	for _, manager := range project.Managers {
		if manager.ID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) IsProjectManager(userName string, projectName string) (bool, error) {
	var count uint

	err := s.db.Table("project_manager").
		Joins("LEFT JOIN projects ON projects.id = project_manager.project_id").
		Joins("LEFT JOIN users ON project_manager.user_id = users.id").
		Where("users.login_name = ?", userName).
		Where("users.stat = 0").
		Where("projects.name = ?", projectName).
		Where("users.deleted_at IS NULL").
		Where("projects.deleted_at IS NULL").
		Count(&count).Error

	return count > 0, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) IsProjectManagerByID(userID, projectID uint) (bool, error) {
	var count uint

	err := s.db.Table("project_manager").
		Joins("projects ON projects.id = project_manager.project_id").
		Joins("JOIN users ON project_manager.user_id = users.id").
		Where("users.id = ?", userID).
		Where("users.stats = 0").
		Where("project_manager.project_id = ?", projectID).
		Where("users.deleted_at IS NULL").
		Where("projects.deleted_at IS NULL").
		Count(&count).Error

	return count > 0, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) IsProjectMember(userName, projectName string) (bool, error) {
	query := `
SELECT EXISTS(
SELECT users.login_name 
FROM users
JOIN project_user on project_user.user_id = users.id
JOIN projects on project_user.project_id = projects.id
JOIN user_group_users on users.id = user_group_users.user_id 
JOIN project_user_group on user_group_users.user_group_id = project_user_group.user_group_id
JOIN projects as p on project_user_group.project_id = p.id
WHERE users.stat = 0
AND( 
	projects.name = ?
OR
	p.name = ?
)) AS exist
`
	var exist struct {
		Exist bool `json:"exist"`
	}
	err := s.db.Raw(query, userName, projectName).Find(&exist).Error
	return exist.Exist, errors.New(errors.ConnectStorageError, err)
}
