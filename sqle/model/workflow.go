package model

type WorkflowTemplate struct {
	Model
	Name      string
	Desc      string
	FlowSteps []*WorkflowStep `json:"-" gorm:"foreignkey:WorkflowTemplateId"`
}

type WorkflowStep struct {
	Id                 int
	StepNumber         int `gorm:"index"`
	WorkflowTemplateId int `gorm:"index"`
	Typ                string
	Desc               string
	Users              []User `json:"-" gorm:"many2many:Workflow_step_user;"`
}

const (
	WorkflowStateRunning = 0
	WorkflowStateFinish  = 1
	WorkflowStateCancel  = 2
)

type Workflow struct {
	Model
	WorkflowTemplateId int
	Template           *WorkflowTemplate `gorm:"foreignkey:WorkflowTemplateId"`
	CreateUserId       int
	CreateUser         User `gorm:"foreignkey:CreateUserId"`
	TaskId             int
	CurrentStepNumber  int
	State              int
	Records            []*WorkflowRecord `gorm:"foreignkey:WorkflowId"`
}

const (
	WorkflowRecordStateAccept = 1
	WorkflowRecordStateReject = 2
)

type WorkflowRecord struct {
	Model
	WorkflowId      int `gorm:"index"`
	WorkflowStepId  int
	OperationUserId int
	OperationUser   User `gorm:"foreignkey:OperationUserId"`
	State           int
	Reason          string
}
