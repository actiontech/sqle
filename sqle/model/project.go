package model

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/errors"

	"github.com/jinzhu/gorm"
)

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
		templateIDStr, err := saveWorkflowTemplate(wt, tx)
		if err != nil {
			return err
		}
		// 这里不会报错, templateIDStr是向数据库写入数据后返回的被修改行ID, 永远是个数字
		templateID, _ := strconv.Atoi(templateIDStr)

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

func (s *Storage) CheckUserCanUpdateProject(projectID uint, userID uint) (bool, error) {
	user, exist, err := s.GetUserByID(userID)
	if err != nil || !exist {
		return false, err
	}

	if user.Name == DefaultAdminUser {
		return true, nil
	}

	project, exist, err := s.GetProjectByID(projectID)
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

func (s *Storage) UpdateProjectInfoByID(projectID uint, attr map[string]interface{}) error {
	err := s.db.Table("projects").Where("id = ?", projectID).Update(attr).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetProjectByID(projectID uint) (*Project, bool, error) {
	p := &Project{}
	err := s.db.Model(&Project{}).Preload("Managers").Where("id = ?", projectID).Find(p).Error
	if err == gorm.ErrRecordNotFound {
		return p, false, nil
	}
	return p, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) ListProject(limit, offset uint32) (p []*Project, totalNum uint64, err error) {
	query := s.db.Model(&Project{})

	err = query.Preload("CreateUser").Limit(limit).Offset(offset).Scan(&p).Error
	if err != nil {
		return nil, 0, errors.New(errors.ConnectStorageError, err)
	}
	err = query.Count(&totalNum).Error
	if err != nil {
		return nil, 0, errors.New(errors.ConnectStorageError, err)
	}

	return
}

func (s *Storage) IsProjectManager(userID uint, projectName string) (bool, error) {
	var count uint

	err := s.db.Table("project_manager").
		Joins("projects ON projects.id = project_manager.project_id").
		Where("project_manager.user_id = ?", userID).
		Where("projects.name = ?", projectName).
		Count(&count).Error

	return count > 0, errors.New(errors.ConnectStorageError, err)
}
