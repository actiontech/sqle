package model

import (
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type WorkflowTemplate struct {
	Model
	Name string
	Desc string

	Steps     []*WorkflowStepTemplate `json:"-" gorm:"foreignkey:workflowTemplateId"`
	Instances []*Instance             `gorm:"foreignkey:WorkflowTemplateId"`
}

const (
	WorkflowStepTypeSQLReview  = "sql_review"
	WorkflowStepTypeSQLExecute = "sql_execute"
	WorkflowStepTypeUnknown    = "unknown"
)

type WorkflowStepTemplate struct {
	Model
	Number             uint   `gorm:"index; column:step_number"`
	WorkflowTemplateId int    `gorm:"index"`
	Typ                string `gorm:"column:type; not null"`
	Desc               string

	Users []*User `gorm:"many2many:workflow_step_template_user"`
}

func (s *Storage) GetWorkflowTemplateByName(name string) (*WorkflowTemplate, bool, error) {
	workflowTemplate := &WorkflowTemplate{}
	err := s.db.Where("name = ?", name).First(workflowTemplate).Error
	if err == gorm.ErrRecordNotFound {
		return workflowTemplate, false, nil
	}
	return workflowTemplate, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetWorkflowTemplateById(id uint) (*WorkflowTemplate, bool, error) {
	workflowTemplate := &WorkflowTemplate{}
	err := s.db.Where("id = ?", id).First(workflowTemplate).Error
	if err == gorm.ErrRecordNotFound {
		return workflowTemplate, false, nil
	}
	return workflowTemplate, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetWorkflowStepsByTemplateId(id uint) ([]*WorkflowStepTemplate, error) {
	steps := []*WorkflowStepTemplate{}
	err := s.db.Where("workflow_template_id = ?", id).Find(&steps).Error
	return steps, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetWorkflowStepsDetailByTemplateId(id uint) ([]*WorkflowStepTemplate, error) {
	steps := []*WorkflowStepTemplate{}
	err := s.db.Preload("Users").Where("workflow_template_id = ?", id).Find(&steps).Error
	return steps, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) SaveWorkflowTemplate(template *WorkflowTemplate) error {
	return s.TxExec(func(tx *sql.Tx) error {
		result, err := tx.Exec("INSERT INTO workflow_templates (name, `desc`) values (?, ?)",
			template.Name, template.Desc)
		if err != nil {
			return err
		}
		templateId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		template.ID = uint(templateId)
		for _, step := range template.Steps {
			result, err = tx.Exec("INSERT INTO workflow_step_templates (step_number, workflow_template_id, type, `desc`) values (?,?,?,?)",
				step.Number, templateId, step.Typ, step.Desc)
			if err != nil {
				return err
			}
			stepId, err := result.LastInsertId()
			if err != nil {
				return err
			}
			step.ID = uint(stepId)
			for _, user := range step.Users {
				_, err = tx.Exec("INSERT INTO workflow_step_template_user (workflow_step_template_id, user_id) values (?,?)",
					stepId, user.ID)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *Storage) UpdateWorkflowTemplateSteps(templateId uint, steps []*WorkflowStepTemplate) error {
	return s.TxExec(func(tx *sql.Tx) error {
		result, err := tx.Exec("UPDATE workflow_step_templates SET workflow_template_id = NULL WHERE workflow_template_id = ?",
			templateId)
		if err != nil {
			return err
		}
		for _, step := range steps {
			result, err = tx.Exec("INSERT INTO workflow_step_templates (step_number, workflow_template_id, type, `desc`) values (?,?,?,?)",
				step.Number, templateId, step.Typ, step.Desc)
			if err != nil {
				return err
			}
			stepId, err := result.LastInsertId()
			if err != nil {
				return err
			}
			step.ID = uint(stepId)
			for _, user := range step.Users {
				_, err = tx.Exec("INSERT INTO workflow_step_template_user (workflow_step_template_id, user_id) values (?,?)",
					stepId, user.ID)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *Storage) UpdateWorkflowTemplateInstances(workflowTemplate *WorkflowTemplate,
	instances ...*Instance) error {
	err := s.db.Model(workflowTemplate).Association("Instances").Replace(instances).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetWorkflowTemplateTip() ([]*WorkflowTemplate, error) {
	templates := []*WorkflowTemplate{}
	err := s.db.Select("name").Find(&templates).Error
	return templates, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

type Workflow struct {
	Model
	Subject          string
	Desc             string
	CreateUserId     uint
	WorkflowRecordId uint
	TaskId           uint

	CreateUser *User           `gorm:"foreignkey:CreateUserId"`
	Record     *WorkflowRecord `gorm:"foreignkey:WorkflowRecordId"`
}

const (
	WorkflowStatusRunning = "on_process"
	WorkflowStatusFinish  = "finished"
	WorkflowStatusReject  = "rejected"
	WorkflowStatusCancel  = "canceled"
)

type WorkflowRecord struct {
	Model
	CurrentWorkflowStepId uint
	Status                string `gorm:"default:\"on_process\""`

	CurrentStep *WorkflowStep   `gorm:"foreignkey:CurrentWorkflowStepId"`
	Steps       []*WorkflowStep `gorm:"foreignkey:WorkflowRecordId"`
}

const (
	WorkflowStepStateInit    = "initialized"
	WorkflowStepStateApprove = "approved"
	WorkflowStepStateReject  = "rejected"
)

const (
	WorkflowStepActionApprove = "approve"
	WorkflowStepActionReject  = "reject"
)

type WorkflowStep struct {
	Model
	OperationUserId        uint
	OperateAt              *time.Time
	WorkflowRecordId       uint   `gorm:"index; not null"`
	WorkflowStepTemplateId uint   `gorm:"index; not null"`
	State                  string `gorm:"default:\"initialized\""`
	Reason                 string

	Template      *WorkflowStepTemplate `gorm:"foreignkey:WorkflowStepTemplateId"`
	OperationUser *User                 `gorm:"foreignkey:OperationUserId"`
}

func (w *Workflow) InitWorkflowStepByTemplate(stepsTemplate []*WorkflowStepTemplate) {
	steps := make([]*WorkflowStep, 0, len(stepsTemplate))
	for _, st := range stepsTemplate {
		step := &WorkflowStep{
			WorkflowStepTemplateId: st.ID,
		}
		steps = append(steps, step)
	}
	workflowRecord := &WorkflowRecord{
		Steps:       steps,
		CurrentStep: steps[0],
	}
	w.Record = workflowRecord
}

func (w *Workflow) CurrentStep() *WorkflowStep {
	return w.Record.CurrentStep
}

func (w *Workflow) NextStep() *WorkflowStep {
	var nextIndex int
	for i, step := range w.Record.Steps {
		if step.ID == w.Record.CurrentWorkflowStepId {
			nextIndex = i + 1
			break
		}
	}
	if nextIndex <= len(w.Record.Steps)-1 {
		return w.Record.Steps[nextIndex]
	}
	return nil
}

func (w *Workflow) FinalStep() *WorkflowStep {
	return w.Record.Steps[len(w.Record.Steps)-1]
}

func (w *Workflow) IsOperationUser(user *User) bool {
	if w.CurrentStep() == nil {
		return false
	}
	for _, assUser := range w.CurrentStep().Template.Users {
		if user.ID == assUser.ID {
			return true
		}
	}
	return false
}

func (s *Storage) UpdateWorkflowStatus(w *Workflow, operateStep *WorkflowStep) error {
	return s.TxExec(func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE workflow_records SET status = ?, current_workflow_step_id = ? WHERE id = ?",
			w.Record.Status, w.Record.CurrentWorkflowStepId, w.Record.ID)
		if err != nil {
			return err
		}
		if operateStep == nil {
			return nil
		}
		_, err = tx.Exec("UPDATE workflow_steps SET operation_user_id = ?, operate_at = ?, state = ?, reason = ? WHERE id = ?",
			operateStep.OperationUserId, operateStep.OperateAt, operateStep.State, operateStep.Reason, operateStep.ID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) SaveWorkflow(workflow *Workflow) error {
	err := s.Save(workflow)
	if err != nil {
		return errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	return nil
}

func (s *Storage) GetWorkflowDetailById(id string) (*Workflow, bool, error) {
	workflow := &Workflow{}
	err := s.db.Preload("CreateUser").Preload("Record").
		Where("id = ?", id).First(workflow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	if workflow.Record == nil {
		return nil, false, errors.New(errors.DataConflict, fmt.Errorf("workflow record not exist"))
	}
	steps := []*WorkflowStep{}
	err = s.db.Where("workflow_record_id = ?", workflow.Record.ID).
		Preload("OperationUser").Find(&steps).Error
	if err != nil {
		return nil, false, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	workflow.Record.Steps = steps
	stepTemplateIds := make([]uint, 0, len(steps))
	for _, step := range steps {
		stepTemplateIds = append(stepTemplateIds, step.WorkflowStepTemplateId)
	}
	stepTemplates := []*WorkflowStepTemplate{}
	err = s.db.Preload("Users").Where("id in (?)", stepTemplateIds).Find(&stepTemplates).Error
	if err != nil {
		return nil, false, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	for _, step := range steps {
		for _, stepTemplate := range stepTemplates {
			if step.WorkflowStepTemplateId == stepTemplate.ID {
				step.Template = stepTemplate
			}
		}
		if step.ID == workflow.Record.CurrentWorkflowStepId {
			workflow.Record.CurrentStep = step
		}
	}
	return workflow, true, nil
}

func (s *Storage) GetWorkflowByTaskId(id string) (*Workflow, bool, error) {
	workflow := &Workflow{}
	err := s.db.Where("task_id = ?", id).First(workflow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	return workflow, true, nil
}

func (s *Storage) DeleteWorkflow(workflow *Workflow) error {
	return s.TxExec(func(tx *sql.Tx) error {
		_, err := tx.Exec("DELETE FROM workflows WHERE id = ?", workflow.ID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM workflow_records WHERE id = ?", workflow.WorkflowRecordId)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM workflow_steps WHERE workflow_record_id = ?", workflow.WorkflowRecordId)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) GetExpiredWorkflows(start time.Time) ([]*Workflow, error) {
	workflows := []*Workflow{}
	err := s.db.Model(&Workflow{}).Select("workflows.id, workflows.workflow_record_id").
		Joins("LEFT JOIN workflow_records ON workflows.workflow_record_id = workflow_records.id").
		Where("workflows.created_at < ?", start).
		Where("workflow_records.status = \"finish\"").
		Or("workflow_records.status = \"canceled\"").
		Or("workflow_records.status IS NULL").
		Scan(&workflows).Error

	return workflows, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
