package model

import (
	"actiontech.cloud/universe/sqle/v4/sqle/errors"
	"database/sql"
	"github.com/jinzhu/gorm"
	"time"
)

type WorkflowTemplate struct {
	Model
	Name string
	Desc string

	Steps     []*WorkflowStepTemplate `json:"-" gorm:"foreignkey:workflowTemplateId"`
	Instances []*Instance             `gorm:"foreignkey:"WorkflowTemplateId"`
}

const (
	WorkflowStepTypeSQLReview  = "sql_review"
	WorkflowStepTypeSQLExecute = "sql_execute"
	WorkflowStepTypeUnknown    = "unknown"
)

type WorkflowStepTemplate struct {
	Model
	Number             int    `gorm:"index; column:step_number"`
	WorkflowTemplateId int    `gorm:"index"`
	Typ                string `gorm:"column:type; not null"`
	Desc               string

	Users []*User `gorm:"many2many:workflow_step_template_user"`
}

type Workflow struct {
	Model
	Subject          string
	Desc             string
	TaskId           uint
	CreateUserId     uint
	WorkflowRecordId uint

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
	NextWorkflowStepId    uint
	Status                string `gorm:"default:\"on_process\"`

	CurrentStep *WorkflowStep   `gorm:"foreignkey:CurrentWorkflowStepId"`
	NextStep    *WorkflowStep   `gorm:"foreignkey:NextWorkflowStepId"`
	Steps       []*WorkflowStep `gorm:"foreignkey:WorkflowRecordId"`
}

const (
	WorkflowStepRecordStateInit    = "initialized"
	WorkflowStepRecordStateApprove = "approved"
	WorkflowStepRecordStateReject  = "rejected"
)

type WorkflowStep struct {
	Model
	OperationUserId        uint
	WorkflowRecordId       uint   `gorm:"index; not null"`
	WorkflowStepTemplateId uint   `gorm:"index; not null"`
	State                  string `gorm:"default:\"initialized\""`
	Reason                 string
	OperateAt              *time.Time

	WorkflowStepTemplate *WorkflowStepTemplate `gorm:"foreignkey:WorkflowStepTemplateId"`
	OperationUser        *User                 `gorm:"foreignkey:OperationUserId"`
}

func (w *Workflow) InitWorkflowRecord(stepsTemplate []*WorkflowStepTemplate) {
	steps := make([]*WorkflowStep, 0, len(stepsTemplate))
	for _, st := range stepsTemplate {
		step := &WorkflowStep{
			WorkflowStepTemplateId: st.ID,
		}
		steps = append(steps, step)
	}
	workflowRecord := &WorkflowRecord{
		Status:   WorkflowStepRecordStateInit,
		Steps:    steps,
		NextStep: steps[0],
	}
	w.Record = workflowRecord
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
		now := time.Now()
		result, err := tx.Exec("UPDATE workflow_step_templates SET deleted_at = ? WHERE workflow_template_id = ?",
			now, templateId)
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

func (s *Storage) SaveWorkflow(workflow *Workflow) error {
	err := s.Save(workflow)
	if err != nil {
		return errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	return nil
}

func (s *Storage) GetWorkflowById(id uint) (*Workflow, bool, error) {
	workflow := &Workflow{}
	err := s.db.Where("id = ?", id).First(workflow).Error
	if err == gorm.ErrRecordNotFound {
		return workflow, false, nil
	}
	return workflow, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetWorkflowRecordsById(id uint) (*WorkflowRecord, bool, error) {
	workflowStatus := &WorkflowRecord{}
	err := s.db.Preload("Steps").Where("id = ?", id).First(workflowStatus).Error
	if err == gorm.ErrRecordNotFound {
		return workflowStatus, false, nil
	}
	return workflowStatus, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

type WorkflowDetail struct {
	Id                    uint                  `json:"workflow_id"`
	WorkflowRecordId      uint                  `json:"workflow_record_id"`
	Subject               string                `json:"subject"`
	Desc                  string                `json:"desc"`
	TaskId                uint                  `json:"task_id"`
	CreateUserName        string                `json:"create_user_name"`
	CreateTime            *time.Time            `json:"create_time"`
	CurrentWorkflowStepId uint                  `json:"current_workflow_step_id"`
	NextWorkflowStepId    uint                  `json:"next_workflow_step_id"`
	Status                string                `json:"status"`
	Steps                 []*WorkflowStepDetail `json:"workflow_step_list"`
}

type WorkflowStepDetail struct {
	Id                uint    `json:"workflow_step_id"`
	Number            uint    `json:"number"`
	Type              string  `json:"type"`
	Desc              string  `json:"desc"`
	State             string  `json:"state"`
	Reason            string  `json:"reason"`
	OperationUser     string  `json:"operation_user_name"`
	OperationTime     string  `json:"operation_time"`
	AssigneeUserNames RowList `json:"assignee_user_names"`
}

var workflowQuery = `
select w.id AS workflow_id, w.task_id, w.workflow_record_id,
w.subject, w.desc, w.created_at AS create_time,
users.login_name AS create_user_name, wr.status,
wr.current_workflow_step_id, wr.next_workflow_step_id
FROM workflows w
LEFT JOIN users  ON w.create_user_id = users.id
LEFT JOIN workflow_records wr ON w.workflow_record_id=wr.id
WHERE w.id = ?
`

var workflowStepQuery = `
SELECT ws.id, ws.state, op_user.login_name operation_user_name,
ws.operate_at operation_time, wst.type, wst.desc, wst.step_number number,
GROUP_CONCAT(DISTINCT COALESCE(ass_user.login_name,'')) AS assignee_user_names
FROM workflow_steps AS ws
LEFT JOIN users AS op_user ON ws.operation_user_id = op_user.id 
LEFT JOIN workflow_step_templates AS wst ON ws.workflow_step_template_id = wst.id 
LEFT JOIN workflow_step_template_user AS wst_re_user ON wst.id = wst_re_user.workflow_step_template_id
LEFT JOIN users ass_user ON wst_re_user.user_id = ass_user.id
WHERE ws.workflow_record_id = ?
GROUP BY ws.id
`

func (s *Storage) GetWorkflowDetailById(id string) (*WorkflowDetail, bool, error) {
	workflowDetail := &WorkflowDetail{}
	err := s.db.Raw(workflowQuery, id).Scan(workflowDetail).Error
	if err == gorm.ErrRecordNotFound {
		return workflowDetail, false, nil
	}
	if err != nil {
		return workflowDetail, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}

	var workflowSteps []*WorkflowStepDetail
	err = s.db.Raw(workflowStepQuery, workflowDetail.WorkflowRecordId).Scan(&workflowSteps).Error
	if err == gorm.ErrRecordNotFound {
		return workflowDetail, false, nil
	}
	if err != nil {
		return workflowDetail, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	workflowDetail.Steps = workflowSteps
	return workflowDetail, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
