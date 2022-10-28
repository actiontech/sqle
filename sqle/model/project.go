package model

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
