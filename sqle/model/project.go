package model

import (
	"database/sql"
	"fmt"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/jinzhu/gorm"
)

const ProjectIdForGlobalRuleTemplate = 0

type Project struct {
	Model
	Name string
	Desc string

	CreateUserId uint  `gorm:"not null"`
	CreateUser   *User `gorm:"foreignkey:CreateUserId"`

	Managers   []*User      `gorm:"many2many:project_manager;"`
	Members    []*User      `gorm:"many2many:project_user;"`
	UserGroups []*UserGroup `gorm:"many2many:project_user_group;"`

	Workflows     []*Workflow     `gorm:"foreignkey:ProjectId"`
	AuditPlans    []*AuditPlan    `gorm:"foreignkey:ProjectId"`
	Instances     []*Instance     `gorm:"foreignkey:ProjectId"`
	SqlWhitelist  []*SqlWhitelist `gorm:"foreignkey:ProjectId"`
	RuleTemplates []*RuleTemplate `gorm:"foreignkey:ProjectId"`

	WorkflowTemplateId uint              `gorm:"not null"`
	WorkflowTemplate   *WorkflowTemplate `gorm:"foreignkey:WorkflowTemplateId"`
}

// IsProjectExist 用于判断当前是否存在项目, 而非某个项目是否存在
func (s *Storage) IsProjectExist() (bool, error) {
	var count uint
	err := s.db.Model(&Project{}).Count(&count).Error
	return count > 0, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) CreateProject(name string, desc string, createUserID uint) error {
	wt := &WorkflowTemplate{
		Name:                          fmt.Sprintf("%v-WorkflowTemplate", name),
		Desc:                          fmt.Sprintf("%v 默认模板", name),
		AllowSubmitWhenLessAuditLevel: string(driver.RuleLevelWarn),
		Steps: []*WorkflowStepTemplate{
			{
				Number: 1,
				Typ:    WorkflowStepTypeSQLReview,
				ApprovedByAuthorized: sql.NullBool{
					Bool:  true,
					Valid: true,
				},
			},
			{
				Number: 2,
				Typ:    WorkflowStepTypeSQLExecute,
				Users:  []*User{{Model: Model{ID: createUserID}}},
			},
		},
	}

	return s.TxExec(func(tx *sql.Tx) error {
		templateID, err := saveWorkflowTemplate(wt, tx)
		if err != nil {
			return err
		}

		result, err := tx.Exec("INSERT INTO projects (`name`, `desc`, `create_user_id`,`workflow_template_id`) values (?, ?, ?,?)", name, desc, createUserID, templateID)
		if err != nil {
			return err
		}
		projectID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO project_manager (`project_id`, `user_id`) VALUES (?, ?)", projectID, createUserID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("INSERT INTO project_user (`project_id`, `user_id`) VALUES (?, ?)", projectID, createUserID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Storage) UpdateProjectInfoByID(projectName string, attr map[string]interface{}) error {
	err := s.db.Table("projects").Where("name = ?", projectName).Update(attr).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetProjectByID(projectID uint) (*Project, bool, error) {
	p := &Project{}
	err := s.db.Model(&Project{}).Where("id = ?", projectID).Find(p).Error
	if err == gorm.ErrRecordNotFound {
		return p, false, nil
	}
	return p, true, errors.New(errors.ConnectStorageError, err)
}

func (s Storage) GetProjectByName(projectName string) (*Project, bool, error) {
	p := &Project{}
	err := s.db.Preload("CreateUser").Preload("Members").Preload("Managers").Preload("Instances").
		Where("name = ?", projectName).First(p).Error
	if err == gorm.ErrRecordNotFound {
		return p, false, nil
	}
	return p, true, errors.New(errors.ConnectStorageError, err)
}

func (s Storage) GetProjectTips(userName string) ([]*Project, error) {
	p := []*Project{}
	query := s.db.Table("projects").Select("name")

	var err error
	if userName != DefaultAdminUser {
		err = query.Joins("JOIN project_user on project_user.project_id = projects.id").
			Joins("JOIN users on users.id = project_user.user_id").
			Joins("JOIN project_user_group on project_user_group.project_id = projects.id").
			Joins("JOIN user_group_users on project_user_group.user_group_id = user_group_users.user_group_id").
			Joins("RIGHT JOIN users as u on u.id = user_group_users.user_id").
			Where("users.stat = 0").Where("u.stat = 0").
			Where("users.login_name = ? OR u.login_name = ?", userName, userName).Find(&p).Error
	} else {
		err = query.Find(&p).Error
	}

	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return p, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) IsUserInProject(userName, projectName string) (bool, error) {
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
AND users.deleted_at IS NULL
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

func (s *Storage) CheckUserIsMember(userName, projectName string) (bool, error) {
	query := `
SELECT EXISTS(
SELECT 1
FROM project_user
LEFT JOIN projects ON projects.id = project_user.project_id
LEFT JOIN users ON users.id = project_user.user_id
WHERE users.stat = 0
AND users.deleted_at IS NULL
AND users.login_name = ?
AND projects.name = ?
) AS exist
`
	var exist struct {
		Exist bool `json:"exist"`
	}
	err := s.db.Raw(query, userName, projectName).Find(&exist).Error
	return exist.Exist, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) CheckUserGroupIsMember(groupName, projectName string) (bool, error) {
	query := `
SELECT EXISTS(
SELECT 1
FROM project_user_group
LEFT JOIN projects ON projects.id = project_user_group.project_id
LEFT JOIN user_groups ON user_groups.id = project_user_group.user_group_id
WHERE user_groups.stat = 0
AND user_groups.deleted_at IS NULL
AND user_groups.name = ?
AND projects.name = ?
) AS exist
`
	var exist struct {
		Exist bool `json:"exist"`
	}
	err := s.db.Raw(query, groupName, projectName).Find(&exist).Error
	return exist.Exist, errors.New(errors.ConnectStorageError, err)
}

type ProjectAndInstance struct {
	InstanceName string `json:"instance_name"`
	ProjectName  string `json:"project_name"`
}

func (s *Storage) GetProjectNamesByInstanceIds(instanceIds []uint) (map[uint] /*instance id*/ ProjectAndInstance, error) {
	instanceIds = utils.RemoveDuplicateUint(instanceIds)
	type record struct {
		InstanceId   uint   `json:"instance_id"`
		InstanceName string `json:"instance_name"`
		ProjectName  string `json:"project_name"`
	}
	records := []record{}
	err := s.db.Table("instances").
		Joins("LEFT JOIN projects ON projects.id = instances.project_id").
		Select("instances.id AS instance_id, instances.name AS instance_name, projects.name AS project_name").
		Where("instances.id IN (?)", instanceIds).
		Find(&records).Error

	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}

	res := make(map[uint]ProjectAndInstance, len(records))
	for _, r := range records {
		res[r.InstanceId] = ProjectAndInstance{
			InstanceName: r.InstanceName,
			ProjectName:  r.ProjectName,
		}
	}

	return res, nil
}

func (s *Storage) AddMember(userName, projectName string, isManager bool, bindRole []BindRole) error {
	user, exist, err := s.GetUserByName(userName)
	if err != nil {
		return errors.ConnectStorageErrWrapper(err)
	}
	if !exist {
		return errors.ConnectStorageErrWrapper(fmt.Errorf("user not exist"))
	}

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return errors.ConnectStorageErrWrapper(err)
	}
	if !exist {
		return errors.ConnectStorageErrWrapper(fmt.Errorf("project not exist"))
	}

	return errors.New(errors.ConnectStorageError, s.db.Transaction(func(tx *gorm.DB) error {

		if err = tx.Exec("INSERT INTO project_user (project_id, user_id) VALUES (?,?)", project.ID, user.ID).Error; err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}

		if isManager {
			if err = tx.Exec("INSERT INTO project_manager (project_id, user_id) VALUES (?,?)", project.ID, user.ID).Error; err != nil {
				return errors.ConnectStorageErrWrapper(err)
			}
		}

		err = s.updateUserRoles(tx, user, projectName, bindRole)
		if err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}
		return nil
	}))
}

func (s *Storage) AddProjectManager(userName, projectName string) error {
	sql := `
INSERT INTO project_manager 
SELECT projects.id AS project_id , users.id AS user_id
FROM projects
JOIN users
WHERE projects.name = ?
AND users.login_name = ?
LIMIT 1
`

	return errors.ConnectStorageErrWrapper(s.db.Exec(sql, projectName, userName).Error)
}

func (s *Storage) RemoveMember(userName, projectName string) error {
	sql := `
DELETE project_user, project_manager, project_member_role  
FROM project_user
LEFT JOIN project_manager ON project_user.project_id = project_manager.project_id 
	AND project_user.user_id = project_manager.user_id
LEFT JOIN projects ON project_user.project_id = projects.id
LEFT JOIN users ON project_user.user_id = users.id
LEFT JOIN project_member_role ON project_member_role.user_id = users.id
WHERE 
users.login_name = ?
AND
projects.name = ?
`

	return errors.ConnectStorageErrWrapper(s.db.Exec(sql, projectName, userName).Error)
}

func (s *Storage) AddMemberGroup(groupName, projectName string, bindRole []BindRole) error {
	group, exist, err := s.GetUserGroupByName(groupName)
	if err != nil {
		return errors.ConnectStorageErrWrapper(err)
	}
	if !exist {
		return errors.ConnectStorageErrWrapper(fmt.Errorf("user group not exist"))
	}

	project, exist, err := s.GetProjectByName(projectName)
	if err != nil {
		return errors.ConnectStorageErrWrapper(err)
	}
	if !exist {
		return errors.ConnectStorageErrWrapper(fmt.Errorf("project not exist"))
	}

	return errors.New(errors.ConnectStorageError, s.db.Transaction(func(tx *gorm.DB) error {

		if err = tx.Exec("INSERT INTO project_user_group (project_id, user_group_id) VALUES (?,?)", project.ID, group.ID).Error; err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}

		err = s.updateUserGroupRoles(tx, group, projectName, bindRole)
		if err != nil {
			return errors.ConnectStorageErrWrapper(err)
		}
		return nil
	}))
}

func (s *Storage) RemoveMemberGroup(groupName, projectName string) error {
	sql := `
DELETE project_user_group, project_member_group_role 
FROM project_user_group
LEFT JOIN projects ON project_user_group.project_id = projects.id
LEFT JOIN user_groups ON project_user_group.user_group_id = user_groups.id
LEFT JOIN project_member_group_role ON project_member_group_role.user_group_id = user_groups.id
WHERE 
user_groups.name = ?
AND
projects.name = ?
`

	return errors.ConnectStorageErrWrapper(s.db.Exec(sql, projectName, groupName).Error)
}

type GetMemberGroupFilter struct {
	FilterProjectName   *string
	FilterUserGroupName *string
	FilterInstanceName  *string
	Limit               *uint32
	Offset              *uint32
}

func generateMemberGroupQueryCriteria(query *gorm.DB, filter GetMemberGroupFilter) (*gorm.DB, error) {
	if filter.FilterProjectName == nil {
		return nil, errors.ConnectStorageErrWrapper(fmt.Errorf("project name cannot be empty"))
	}

	query = query.Model(&UserGroup{}).
		Joins("LEFT JOIN project_user_group ON project_user_group.user_group_id = user_groups.id").
		Joins("LEFT JOIN projects ON projects.id = project_user_group.project_id").
		Joins("LEFT JOIN instances ON instances.project_id = projects.id").
		Where("user_groups.stat = 0").
		Where("user_groups.deleted_at IS NULL").
		Where("instances.deleted_at IS NULL").
		Where("projects.name = ?", *filter.FilterProjectName)

	if filter.Limit != nil {
		query = query.Limit(*filter.Limit).Offset(*filter.Offset)
	}
	if filter.FilterUserGroupName != nil {
		query = query.Where("user_groups.name = ?", *filter.FilterUserGroupName)
	}
	if filter.FilterInstanceName != nil {
		query = query.Where("instances.name = ?", *filter.FilterInstanceName)
	}
	return query, nil
}

func (s *Storage) GetMemberGroups(filter GetMemberGroupFilter) ([]*UserGroup, error) {
	group := []*UserGroup{}
	query, err := generateMemberGroupQueryCriteria(s.db, filter)
	if err != nil {
		return nil, err
	}
	err = query.Scan(&group).Error
	return group, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetMemberGroupCount(filter GetMemberGroupFilter) (uint64, error) {
	// count要查总数, 清掉limit
	filter.Limit = nil
	filter.Offset = nil

	var count uint64
	query, err := generateMemberGroupQueryCriteria(s.db, filter)
	if err != nil {
		return 0, err
	}
	err = query.Count(&count).Error
	return count, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetMemberGroupByGroupName(projectName, groupName string) (*UserGroup, error) {
	group := &UserGroup{}
	err := s.db.Joins("LEFT JOIN project_user_group ON project_user_group.user_group_id = user_groups.id").
		Joins("LEFT JOIN projects ON project_user_group.project_id = projects.id").
		Where("projects.name = ?", projectName).
		Where("user_groups.name = ?", groupName).
		Find(group).Error
	return group, errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) RemoveMemberFromAllProjectByUserID(userID uint) error {
	sql := `
DELETE project_user, project_manager, project_member_role 
FROM project_user
LEFT JOIN project_manager ON project_user.project_id = project_manager.project_id 
	AND project_user.user_id = project_manager.user_id
LEFT JOIN project_member_role ON project_member_role.user_id = ?
LEFT JOIN users ON project_user.user_id = users.id
WHERE 
users.id = ?
`

	return errors.ConnectStorageErrWrapper(s.db.Exec(sql, userID, userID).Error)
}

func (s *Storage) RemoveMemberGroupFromAllProjectByUserGroupID(userGroupID uint) error {
	sql := `
DELETE project_user_group, project_member_group_role 
FROM project_user_group
LEFT JOIN projects ON project_user_group.project_id = projects.id
LEFT JOIN user_groups ON project_user_group.user_group_id = user_groups.id
LEFT JOIN project_member_group_role ON project_member_group_role.user_group_id = ?
WHERE
user_groups.id = ?
`

	return errors.ConnectStorageErrWrapper(s.db.Exec(sql, userGroupID, userGroupID).Error)
}

// 检查用户是否是某一个项目的最后一个管理员
func (s *Storage) IsLastProjectManagerOfAnyProjectByUserID(userID uint) (bool, error) {

	sql := `
SELECT count(1) AS count
FROM project_manager
WHERE project_manager.user_id = ?
GROUP BY project_manager.project_id
`

	var count []*struct {
		Count int `json:"count"`
	}
	err := s.db.Raw(sql, userID).Scan(&count).Error
	if err != nil {
		return true, errors.ConnectStorageErrWrapper(err)
	}

	for _, c := range count {
		if c.Count == 1 {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) CheckUserHasManagementPermission(userID uint, code []uint) (bool, error) {
	code = utils.RemoveDuplicateUint(code)

	user, _, err := s.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	if user.Name == DefaultAdminUser {
		return true, nil
	}

	var count int
	err = s.db.Model(&ManagementPermission{}).
		Where("user_id = ?", userID).
		Where("permission_code in (?)", code).
		Count(&count).Error

	return count == len(code), errors.ConnectStorageErrWrapper(err)
}

func (s *Storage) GetManagedProjects(userID uint) ([]*Project, error) {
	p := []*Project{}
	err := s.db.Joins("LEFT JOIN project_manager ON project_manager.project_id = projects.id").
		Where("project_manager.user_id = ?", userID).
		Find(&p).Error
	return p, errors.ConnectStorageErrWrapper(err)
}
