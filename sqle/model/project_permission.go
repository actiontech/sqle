package model

import "github.com/actiontech/sqle/sqle/errors"

// import (
// 	"github.com/actiontech/sqle/sqle/errors"
// 	"github.com/jinzhu/gorm"
// )

// /*

// instance permission.

// */

// var queryInstanceUserWithOp = `
// SELECT

// {{- template "select_fields" . -}}

// FROM instances
// {{- if .project_name }}
// LEFT JOIN projects ON instances.project_id = projects.id
// {{- end }}
// LEFT JOIN project_member_roles ON instances.id = project_member_roles.instance_id
// LEFT JOIN users ON project_member_roles.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// LEFT JOIN roles ON project_member_roles.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
// LEFT JOIN role_operations ON role_operations.role_id = roles.id
// WHERE
// instances.deleted_at IS NULL
// AND users.id = :user_id
// AND role_operations.op_code IN (:op_codes)

// {{- if .instance_ids }}
// AND instances.id IN (:instance_ids)
// {{- end }}

// {{- if .project_name }}
// AND projects.name = :project_name
// {{- end }}

// {{- if .db_type }}
// AND instances.db_type = :db_type
// {{- end }}
// GROUP BY instances.id

// UNION
// SELECT

// {{- template "select_fields" . -}}

// FROM instances
// {{- if .project_name }}
// LEFT JOIN projects ON instances.project_id = projects.id
// {{- end }}
// LEFT JOIN project_member_group_roles ON instances.id = project_member_group_roles.instance_id
// LEFT JOIN roles ON roles.id = project_member_group_roles.role_id AND roles.deleted_at IS NULL AND roles.stat = 0
// LEFT JOIN role_operations ON role_operations.role_id = roles.id
// LEFT JOIN user_groups ON project_member_group_roles.user_group_id = user_groups.id AND user_groups.deleted_at IS NULL AND user_groups.stat = 0
// JOIN user_group_users ON user_groups.id = user_group_users.user_group_id
// JOIN users ON users.id = user_group_users.user_id AND users.deleted_at IS NULL AND users.stat=0
// WHERE
// instances.deleted_at IS NULL
// AND users.id = :user_id
// AND role_operations.op_code IN (:op_codes)

// {{- if .instance_ids }}
// AND instances.id IN (:instance_ids)
// {{- end }}

// {{- if .project_name }}
// AND projects.name = :project_name
// {{- end }}

// {{- if .db_type }}
// AND instances.db_type = :db_type
// {{- end }}
// GROUP BY instances.id

// UNION
// SELECT

// {{- template "select_fields" . -}}

// FROM instances
// LEFT JOIN projects ON instances.project_id = projects.id
// LEFT JOIN project_manager ON project_manager.project_id = projects.id
// LEFT JOIN users ON project_manager.user_id = users.id AND users.deleted_at IS NULL AND users.stat = 0
// WHERE
// instances.deleted_at IS NULL
// AND users.id = :user_id

// {{- if .instance_ids }}
// AND instances.id IN (:instance_ids)
// {{- end }}

// {{- if .project_name }}
// AND projects.name = :project_name
// {{- end }}

// {{- if .db_type }}
// AND instances.db_type = :db_type
// {{- end }}
// GROUP BY instances.id

// `

// func (s *Storage) filterUserHasOpInstances(user *User, instanceIds []uint, ops []uint) ([]*Instance, error) {
// 	var instanceRecords []*Instance
// 	data := map[string]interface{}{
// 		"instance_ids": instanceIds,
// 		"op_codes":     ops,
// 		"user_id":      user.ID,
// 	}
// 	fields := `
// {{ define "select_fields" }}
// instances.id
// {{ end }}
// 	`
// 	err := s.getTemplateQueryResult(data, &instanceRecords, queryInstanceUserWithOp, fields)
// 	if err != nil {
// 		return nil, errors.ConnectStorageErrWrapper(err)
// 	}
// 	return instanceRecords, nil
// }

// func (s *Storage) CheckUserHasOpToInstances(user *User, instances []*Instance, ops []uint) (bool, error) {
// 	instanceIds := getDeduplicatedInstanceIds(instances)
// 	instanceRecords, err := s.filterUserHasOpInstances(user, instanceIds, ops)
// 	if err != nil {
// 		return false, err
// 	}
// 	return len(instanceRecords) == len(instanceIds), nil
// }

// func (s *Storage) CheckUserHasOpToAnyInstance(user *User, instances []*Instance, ops []uint) (bool, error) {
// 	instanceIds := getDeduplicatedInstanceIds(instances)
// 	instanceRecords, err := s.filterUserHasOpInstances(user, instanceIds, ops)
// 	if err != nil {
// 		return false, err
// 	}
// 	return len(instanceRecords) > 0, nil
// }

// func (s *Storage) GetUserCanOpInstances(user *User, ops []uint) ([]*Instance, error) {
// 	var instances []*Instance
// 	data := map[string]interface{}{
// 		"op_codes": ops,
// 		"user_id":  user.ID,
// 	}
// 	fields := `
// {{ define "select_fields" }}
// instances.id, instances.name
// {{ end }}
// `
// 	err := s.getTemplateQueryResult(data, &instances, queryInstanceUserWithOp, fields)
// 	if err != nil {
// 		return instances, errors.ConnectStorageErrWrapper(err)
// 	}
// 	return instances, nil
// }

// func (s *Storage) GetUserCanOpInstancesFromProject(user *User, projectName string, ops []uint) ([]*Instance, error) {
// 	var instances []*Instance
// 	data := map[string]interface{}{
// 		"op_codes":     ops,
// 		"user_id":      user.ID,
// 		"project_name": projectName,
// 	}
// 	fields := `
// {{ define "select_fields" }}
// instances.id, instances.name
// {{ end }}
// `
// 	err := s.getTemplateQueryResult(data, &instances, queryInstanceUserWithOp, fields)
// 	if err != nil {
// 		return instances, errors.ConnectStorageErrWrapper(err)
// 	}
// 	return instances, nil
// }

// func (s *Storage) GetInstanceTipsByUserAndOperation(user *User, dbType, projectName string, opCode ...int) (
// 	instances []*Instance, err error) {

// 	isProjectManager, err := s.IsProjectManager(user.Name, projectName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if IsDefaultAdminUser(user.Name) || isProjectManager {
// 		return s.GetInstanceTipsByTypeAndTempID(dbType, 0, projectName)
// 	}
// 	return s.getInstanceTipsByUserAndOperation(user, dbType, projectName, opCode...)
// }

// func (s *Storage) getInstanceTipsByUserAndOperation(user *User, dbType string, projectName string, opCode ...int) ([]*Instance, error) {
// 	var instances []*Instance
// 	data := map[string]interface{}{
// 		"op_codes":     opCode,
// 		"user_id":      user.ID,
// 		"project_name": projectName,
// 		"db_type":      dbType,
// 	}
// 	fields := `
// {{ define "select_fields" }}
// instances.id, instances.name, instances.db_host as host, instances.db_port as port, instances.db_type
// {{ end }}
// 	`
// 	err := s.getTemplateQueryResult(data, &instances, queryInstanceUserWithOp, fields)
// 	if err != nil {
// 		return instances, errors.ConnectStorageErrWrapper(err)
// 	}
// 	return instances, nil
// }

// func (s *Storage) UserCanAccessInstance(user *User, instance *Instance) (
// 	ok bool, err error) {

// 	isManager, err := s.IsProjectManagerByID(user.ID, instance.ProjectId)
// 	if err != nil {
// 		return false, err
// 	}

// 	if IsDefaultAdminUser(user.Name) || isManager {
// 		return true, nil
// 	}

// 	type countStruct struct {
// 		Count int `json:"count"`
// 	}

// 	query := `
// SELECT COUNT(1) AS count
// FROM instances
// LEFT JOIN project_member_roles ON project_member_roles.instance_id = instances.id
// LEFT JOIN users ON users.id = project_member_roles.user_id
// WHERE instances.deleted_at IS NULL
// AND users.stat = 0
// AND users.deleted_at IS NULL
// AND instances.id = ?
// AND users.id = ?
// GROUP BY instances.id
// UNION
// SELECT instances.id
// FROM instances
// LEFT JOIN project_member_group_roles ON project_member_group_roles.instance_id = instances.id
// JOIN user_group_users ON project_member_group_roles.user_group_id = user_group_users.user_group_id
// JOIN users ON users.id = user_group_users.user_id
// WHERE instances.deleted_at IS NULL
// AND users.stat = 0
// AND users.deleted_at IS NULL
// AND instances.id = ?
// AND users.id = ?
// GROUP BY instances.id
// `
// 	var cnt countStruct
// 	err = s.db.Unscoped().Raw(query, instance.ID, user.ID, instance.ID, user.ID).Scan(&cnt).Error
// 	if err != nil {
// 		if gorm.IsRecordNotFoundError(err) {
// 			return false, nil
// 		}
// 		return false, errors.New(errors.ConnectStorageError, err)
// 	}
// 	return cnt.Count > 0, nil
// }

// func (s *Storage) GetWithOperationUserFromInstance(instance *Instance, opCode ...int) (users []*User, err error) {
// 	query := `
// 	SELECT users.id, users.login_name
// 	FROM users
// 	LEFT JOIN project_member_roles ON users.id = project_member_roles.user_id
// 	LEFT JOIN roles ON project_member_roles.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
// 	LEFT JOIN role_operations ON role_operations.role_id = roles.id
// 	WHERE
// 	users.deleted_at IS NULL
// 	AND users.stat = 0
// 	AND project_member_roles.instance_id = ?
// 	AND role_operations.op_code IN (?)

// 	UNION
// 	SELECT users.id, users.login_name
// 	FROM users
// 	LEFT JOIN user_group_users ON users.id = user_group_users.user_id
// 	LEFT JOIN user_groups ON user_group_users.user_group_id = user_groups.id AND user_groups.stat = 0
// 	LEFT JOIN project_member_group_roles ON user_groups.id = project_member_group_roles.user_group_id
// 	LEFT JOIN roles ON project_member_group_roles.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
// 	LEFT JOIN role_operations ON role_operations.role_id = roles.id
// 	WHERE
// 	users.deleted_at IS NULL
// 	AND users.stat = 0
// 	AND project_member_group_roles.instance_id = ?
// 	AND role_operations.op_code IN (?)
// 	`
// 	err = s.db.Raw(query, instance.ID, opCode, instance.ID, opCode).Scan(&users).Error
// 	if err != nil {
// 		return nil, errors.ConnectStorageErrWrapper(err)
// 	}
// 	return
// }

// /*

// workflow permission.

// */

// func (s *Storage) UserCanAccessWorkflow(user *User, workflow *Workflow) (bool, error) {
// 	query := `SELECT count(w.id) FROM workflows AS w
// JOIN workflow_records AS wr ON w.workflow_record_id = wr.id AND w.id = ?
// LEFT JOIN workflow_steps AS cur_ws ON wr.current_workflow_step_id = cur_ws.id
// LEFT JOIN workflow_step_templates AS cur_wst ON cur_ws.workflow_step_template_id = cur_wst.id
// LEFT JOIN workflow_step_user AS cur_wst_re_user ON cur_ws.id = cur_wst_re_user.workflow_step_id
// LEFT JOIN users AS cur_ass_user ON cur_wst_re_user.user_id = cur_ass_user.id AND cur_ass_user.stat=0
// LEFT JOIN workflow_steps AS op_ws ON w.id = op_ws.workflow_id AND op_ws.state != "initialized"
// LEFT JOIN workflow_step_templates AS op_wst ON op_ws.workflow_step_template_id = op_wst.id
// LEFT JOIN workflow_step_user AS op_wst_re_user ON op_ws.id = op_wst_re_user.workflow_step_id
// LEFT JOIN users AS op_ass_user ON op_wst_re_user.user_id = op_ass_user.id AND op_ass_user.stat=0
// where w.deleted_at IS NULL
// AND (w.create_user_id = ? OR cur_ass_user.id = ? OR op_ass_user.id = ?)
// `
// 	var count uint
// 	err := s.db.Raw(query, workflow.ID, user.ID, user.ID, user.ID).Count(&count).Error
// 	if err != nil {
// 		return false, errors.New(errors.ConnectStorageError, err)
// 	}
// 	return count > 0, nil
// }

// /*

// workflow permission. TODO DMS权限接口替换掉UserCanAccessWorkflow

//    去除了与用户表关联关系（TODO 用户不可用时需要判断）
//  LEFT JOIN users AS cur_ass_user ON cur_wst_re_user.user_id = cur_ass_user.id AND cur_ass_user.stat=0
//  LEFT JOIN users AS op_ass_user ON op_wst_re_user.user_id = op_ass_user.id AND op_ass_user.stat=0
//  LEFT JOIN workflow_step_user AS cur_wst_re_user ON cur_ws.id = cur_wst_re_user.workflow_step_id
//  LEFT JOIN workflow_step_user AS op_wst_re_user ON op_ws.id = op_wst_re_user.workflow_step_id
// */

func (s *Storage) UserCanAccessWorkflow(userId string, workflow *Workflow) (bool, error) {
	query := `SELECT count(w.id) FROM workflows AS w
JOIN workflow_records AS wr ON w.workflow_record_id = wr.id AND w.workflow_id = ?
LEFT JOIN workflow_steps AS cur_ws ON wr.current_workflow_step_id = cur_ws.id
LEFT JOIN workflow_step_templates AS cur_wst ON cur_ws.workflow_step_template_id = cur_wst.id
LEFT JOIN workflow_steps AS op_ws ON w.id = op_ws.workflow_id AND op_ws.state != "initialized"
LEFT JOIN workflow_step_templates AS op_wst ON op_ws.workflow_step_template_id = op_wst.id
where w.deleted_at IS NULL
AND (w.create_user_id = ? OR cur_ws.assignees REGEXP ? OR op_ws.assignees REGEXP ?)
`
	var count uint
	err := s.db.Raw(query, workflow.WorkflowId, userId, userId, userId).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ConnectStorageError, err)
	}
	return count > 0, nil
}

// // GetCanAuditWorkflowUsers will return admin user if no qualified user is found, preventing the process from being stuck because no user can operate
// func (s *Storage) GetCanAuditWorkflowUsers(instance *Instance) (users []*User, err error) {
// 	users, err = s.GetWithOperationUserFromInstance(instance, OP_WORKFLOW_AUDIT)
// 	if err != nil {
// 		return
// 	}
// 	if len(users) != 0 {
// 		return
// 	}
// 	return s.GetUsersByNames([]string{DefaultAdminUser})
// }

// // GetCanExecuteWorkflowUsers will return admin user if no qualified user is found, preventing the process from being stuck because no user can operate
// func (s *Storage) GetCanExecuteWorkflowUsers(instance *Instance) (users []*User, err error) {
// 	users, err = s.GetWithOperationUserFromInstance(instance, OP_WORKFLOW_EXECUTE)
// 	if err != nil {
// 		return
// 	}
// 	if len(users) != 0 {
// 		return
// 	}
// 	return s.GetUsersByNames([]string{DefaultAdminUser})
// }

// /*

// audit plan permission.

// */

// func (s *Storage) CheckUserCanCreateAuditPlan(user *User, projectName, instName string) (bool, error) {
// 	if IsDefaultAdminUser(user.Name) {
// 		return true, nil
// 	}

// 	isManager, err := s.IsProjectManager(user.Name, projectName)
// 	if err != nil {
// 		return false, err
// 	}

// 	if isManager {
// 		return true, nil
// 	}

// 	// todo: check it in db, don't get all instances.
// 	instances, err := s.GetUserCanOpInstancesFromProject(user, projectName, []uint{OP_AUDIT_PLAN_SAVE})
// 	if err != nil {
// 		return false, err
// 	}
// 	for _, instance := range instances {
// 		if instName == instance.Name {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// /*

// project permission.

// */

// func (s *Storage) CheckUserCanUpdateProject(projectName string, userID uint) (bool, error) {
// 	user, exist, err := s.GetUserByID(userID)
// 	if err != nil || !exist {
// 		return false, err
// 	}

// 	if user.Name == DefaultAdminUser {
// 		return true, nil
// 	}

// 	project, exist, err := s.GetProjectByName(projectName)
// 	if err != nil || !exist {
// 		return false, err
// 	}

// 	for _, manager := range project.Managers {
// 		if manager.ID == userID {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// func (s *Storage) IsProjectManager(userName string, projectName string) (bool, error) {
// 	var count uint

// 	err := s.db.Table("project_manager").
// 		Joins("LEFT JOIN projects ON projects.id = project_manager.project_id").
// 		Joins("LEFT JOIN users ON project_manager.user_id = users.id").
// 		Where("users.login_name = ?", userName).
// 		Where("users.stat = 0").
// 		Where("projects.name = ?", projectName).
// 		Where("users.deleted_at IS NULL").
// 		Where("projects.deleted_at IS NULL").
// 		Count(&count).Error

// 	return count > 0, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) IsProjectManagerByID(userID, projectID uint) (bool, error) {
// 	var count uint

// 	err := s.db.Table("project_manager").
// 		Joins("LEFT JOIN projects ON projects.id = project_manager.project_id").
// 		Joins("LEFT JOIN users ON project_manager.user_id = users.id").
// 		Where("users.id = ?", userID).
// 		Where("users.stat = 0").
// 		Where("project_manager.project_id = ?", projectID).
// 		Where("users.deleted_at IS NULL").
// 		Where("projects.deleted_at IS NULL").
// 		Count(&count).Error

// 	return count > 0, errors.New(errors.ConnectStorageError, err)
// }

// func (s *Storage) IsProjectMember(userName, projectName string) (bool, error) {
// 	query := `
// SELECT EXISTS(
// SELECT users.login_name
// FROM users
// LEFT JOIN project_user on project_user.user_id = users.id
// LEFT JOIN projects on project_user.project_id = projects.id
// LEFT JOIN user_group_users on users.id = user_group_users.user_id
// LEFT JOIN project_user_group on user_group_users.user_group_id = project_user_group.user_group_id
// LEFT JOIN projects as p on project_user_group.project_id = p.id
// WHERE users.stat = 0
// AND users.login_name = ?
// AND(
// 	projects.name = ?
// OR
// 	p.name = ?
// )) AS exist
// `
// 	var exist struct {
// 		Exist bool `json:"exist"`
// 	}
// 	err := s.db.Raw(query, userName, projectName, projectName).Find(&exist).Error
// 	return exist.Exist, errors.New(errors.ConnectStorageError, err)
// }
