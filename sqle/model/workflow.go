package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/jinzhu/gorm"
)

type WorkflowTemplate struct {
	Model
	Name                          string
	Desc                          string
	AllowSubmitWhenLessAuditLevel string

	Steps     []*WorkflowStepTemplate `json:"-" gorm:"foreignkey:workflowTemplateId"`
	Instances []*Instance             `gorm:"foreignkey:WorkflowTemplateId"`
}

const (
	WorkflowStepTypeSQLReview      = "sql_review"
	WorkflowStepTypeSQLExecute     = "sql_execute"
	WorkflowStepTypeCreateWorkflow = "create_workflow"
	WorkflowStepTypeUpdateWorkflow = "update_workflow"
)

type WorkflowStepTemplate struct {
	Model
	Number               uint   `gorm:"index; column:step_number"`
	WorkflowTemplateId   int    `gorm:"index"`
	Typ                  string `gorm:"column:type; not null"`
	Desc                 string
	ApprovedByAuthorized sql.NullBool `gorm:"column:approved_by_authorized"`

	Users []*User `gorm:"many2many:workflow_step_template_user"`
}

func (s *Storage) GetWorkflowTemplateByName(name string) (*WorkflowTemplate, bool, error) {
	workflowTemplate := &WorkflowTemplate{}
	err := s.db.Where("name = ?", name).First(workflowTemplate).Error
	if err == gorm.ErrRecordNotFound {
		return workflowTemplate, false, nil
	}
	return workflowTemplate, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetWorkflowTemplateById(id uint) (*WorkflowTemplate, bool, error) {
	workflowTemplate := &WorkflowTemplate{}
	err := s.db.Where("id = ?", id).First(workflowTemplate).Error
	if err == gorm.ErrRecordNotFound {
		return workflowTemplate, false, nil
	}
	return workflowTemplate, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetWorkflowStepsByTemplateId(id uint) ([]*WorkflowStepTemplate, error) {
	steps := []*WorkflowStepTemplate{}
	err := s.db.Preload("Users").Where("workflow_template_id = ?", id).Find(&steps).Error
	return steps, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetWorkflowStepsDetailByTemplateId(id uint) ([]*WorkflowStepTemplate, error) {
	steps := []*WorkflowStepTemplate{}
	err := s.db.Preload("Users").Where("workflow_template_id = ?", id).Find(&steps).Error
	return steps, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) SaveWorkflowTemplate(template *WorkflowTemplate) error {
	return s.TxExec(func(tx *sql.Tx) error {
		result, err := tx.Exec("INSERT INTO workflow_templates (name, `desc`, `allow_submit_when_less_audit_level`) values (?, ?, ?)",
			template.Name, template.Desc, template.AllowSubmitWhenLessAuditLevel)
		if err != nil {
			return err
		}
		templateId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		template.ID = uint(templateId)
		for _, step := range template.Steps {
			result, err = tx.Exec("INSERT INTO workflow_step_templates (step_number, workflow_template_id, type, `desc`, approved_by_authorized) values (?,?,?,?,?)",
				step.Number, templateId, step.Typ, step.Desc, step.ApprovedByAuthorized)
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
		_, err := tx.Exec("UPDATE workflow_step_templates SET workflow_template_id = NULL WHERE workflow_template_id = ?",
			templateId)
		if err != nil {
			return err
		}
		for _, step := range steps {
			result, err := tx.Exec("INSERT INTO workflow_step_templates (step_number, workflow_template_id, type, `desc`, approved_by_authorized) values (?,?,?,?,?)",
				step.Number, templateId, step.Typ, step.Desc, step.ApprovedByAuthorized)
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
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetWorkflowTemplateTip() ([]*WorkflowTemplate, error) {
	templates := []*WorkflowTemplate{}
	err := s.db.Select("name").Find(&templates).Error
	return templates, errors.New(errors.ConnectStorageError, err)
}

type Workflow struct {
	Model
	Subject          string
	Desc             string
	CreateUserId     uint
	WorkflowRecordId uint

	CreateUser    *User             `gorm:"foreignkey:CreateUserId"`
	Record        *WorkflowRecord   `gorm:"foreignkey:WorkflowRecordId"`
	RecordHistory []*WorkflowRecord `gorm:"many2many:workflow_record_history;"`
	Mode          string
}

const (
	WorkflowStatusWaitForAudit     = "wait_for_audit"
	WorkflowStatusWaitForExecution = "wait_for_execution"
	WorkflowStatusReject           = "rejected"
	WorkflowStatusCancel           = "canceled"
	WorkflowStatusExecuting        = "executing"
	WorkflowStatusExecFailed       = "exec_failed"
	WorkflowStatusFinish           = "finished"

	WorkflowModeSameSQLs      = "same_sqls"
	WorkflowModeDifferentSQLs = "different_sqls"
)

type WorkflowRecord struct {
	Model
	CurrentWorkflowStepId uint
	Status                string                    `gorm:"default:\"wait_for_audit\""`
	InstanceRecords       []*WorkflowInstanceRecord `gorm:"foreignkey:WorkflowRecordId"`

	// 当workflow只有部分数据源已上线时，current step仍处于"sql_execute"步骤
	CurrentStep *WorkflowStep   `gorm:"foreignkey:CurrentWorkflowStepId"`
	Steps       []*WorkflowStep `gorm:"foreignkey:WorkflowRecordId"`
}

type WorkflowInstanceRecord struct {
	Model
	TaskId           uint `gorm:"index"`
	WorkflowRecordId uint `gorm:"index; not null"`
	InstanceId       uint
	ScheduledAt      *time.Time
	ScheduleUserId   uint
	// 用于区分工单处于上线步骤时，某个数据源是否已上线，因为数据源可以分批上线
	IsSQLExecuted   bool
	ExecutionUserId uint
}

func (s *Storage) GetWorkInstanceRecordByTaskId(id string) (instanceRecord WorkflowInstanceRecord, err error) {
	return instanceRecord, s.db.Where("task_id = ?", id).First(&instanceRecord).Error
}

const (
	WorkflowStepStateInit    = "initialized"
	WorkflowStepStateApprove = "approved"
	WorkflowStepStateReject  = "rejected"
)

type WorkflowStep struct {
	Model
	OperationUserId        uint
	OperateAt              *time.Time
	WorkflowId             uint   `gorm:"index; not null"`
	WorkflowRecordId       uint   `gorm:"index; not null"`
	WorkflowStepTemplateId uint   `gorm:"index; not null"`
	State                  string `gorm:"default:\"initialized\""`
	Reason                 string

	Assignees     []*User               `gorm:"many2many:workflow_step_user"`
	Template      *WorkflowStepTemplate `gorm:"foreignkey:WorkflowStepTemplateId"`
	OperationUser *User                 `gorm:"foreignkey:OperationUserId"`
}

func generateWorkflowStepByTemplate(stepsTemplate []*WorkflowStepTemplate, allInspector []*User) []*WorkflowStep {
	steps := make([]*WorkflowStep, 0, len(stepsTemplate))
	for _, st := range stepsTemplate {
		step := &WorkflowStep{
			WorkflowStepTemplateId: st.ID,
			Assignees:              st.Users,
		}
		if st.ApprovedByAuthorized.Bool {
			step.Assignees = allInspector
		}
		steps = append(steps, step)
	}
	return steps
}

func (w *Workflow) cloneWorkflowStep() []*WorkflowStep {
	steps := make([]*WorkflowStep, 0, len(w.Record.Steps))
	for _, step := range w.Record.Steps {
		steps = append(steps, &WorkflowStep{
			WorkflowStepTemplateId: step.Template.ID,
			WorkflowId:             w.ID,
			Assignees:              step.Assignees,
		})
	}
	return steps
}

func (w *Workflow) CreateUserName() string {
	if w.CreateUser != nil {
		return w.CreateUser.Name
	}
	return ""
}

func (w *Workflow) CurrentStep() *WorkflowStep {
	return w.Record.CurrentStep
}

func (w *Workflow) CurrentAssigneeUser() []*User {
	currentStep := w.CurrentStep()
	if currentStep == nil {
		return []*User{}
	}
	return currentStep.Assignees
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
	for _, assUser := range w.CurrentStep().Assignees {
		if user.ID == assUser.ID {
			return true
		}
	}
	return false
}

// IsFirstRecord check the record is the first record in workflow;
// you must load record history first and then use it.
func (w *Workflow) IsFirstRecord(record *WorkflowRecord) bool {
	records := []*WorkflowRecord{}
	records = append(records, w.RecordHistory...)
	records = append(records, w.Record)
	if len(records) > 0 {
		return record == records[0]
	}
	return false
}

func (w *Workflow) GetTaskIds() []uint {
	taskIds := make([]uint, len(w.Record.InstanceRecords))
	for i, inst := range w.Record.InstanceRecords {
		taskIds[i] = inst.TaskId
	}
	return taskIds
}

func (s *Storage) CreateWorkflow(subject, desc string, user *User, tasks []*Task,
	stepTemplates []*WorkflowStepTemplate) error {
	if len(tasks) <= 0 {
		return errors.New(errors.DataConflict, fmt.Errorf("there is no task for creating workflow"))
	}

	workflowMode := WorkflowModeSameSQLs
	groupId := tasks[0].GroupId
	for _, task := range tasks {
		if task.GroupId != groupId {
			workflowMode = WorkflowModeDifferentSQLs
			break
		}
	}

	workflow := &Workflow{
		Subject:      subject,
		Desc:         desc,
		CreateUserId: user.ID,
		Mode:         workflowMode,
	}

	instanceRecords := make([]*WorkflowInstanceRecord, len(tasks))
	for i, task := range tasks {
		instanceRecords[i] = &WorkflowInstanceRecord{
			TaskId:     task.ID,
			InstanceId: task.InstanceId,
		}
	}
	record := &WorkflowRecord{
		InstanceRecords: instanceRecords,
	}

	allUsers := make([][]*User, len(tasks))
	for i, task := range tasks {
		users, err := s.GetUsersByOperationCode(task.Instance, OP_WORKFLOW_AUDIT)
		if err != nil {
			return err
		}
		allUsers[i] = users
	}

	canOptUsers := allUsers[0]
	for i := 1; i < len(allUsers); i++ {
		canOptUsers = getOverlapOfUsers(canOptUsers, allUsers[i])
	}

	steps := generateWorkflowStepByTemplate(stepTemplates, canOptUsers)

	tx := s.db.Begin()

	err := tx.Save(record).Error
	if err != nil {
		tx.Rollback()
		return errors.New(errors.ConnectStorageError, err)
	}

	workflow.WorkflowRecordId = record.ID
	err = tx.Save(workflow).Error
	if err != nil {
		tx.Rollback()
		return errors.New(errors.ConnectStorageError, err)
	}

	for _, step := range steps {
		currentStep := step
		currentStep.WorkflowRecordId = record.ID
		currentStep.WorkflowId = workflow.ID
		users := currentStep.Assignees
		currentStep.Assignees = nil
		err = tx.Save(currentStep).Error
		if err != nil {
			tx.Rollback()
			return errors.New(errors.ConnectStorageError, err)
		}
		err = tx.Model(currentStep).Association("Assignees").Replace(users).Error
		if err != nil {
			tx.Rollback()
			return errors.New(errors.ConnectStorageError, err)
		}
	}
	if len(steps) > 0 {
		err = tx.Model(record).Update("current_workflow_step_id", steps[0].ID).Error
		if err != nil {
			tx.Rollback()
			return errors.New(errors.ConnectStorageError, err)
		}
	}
	return errors.New(errors.ConnectStorageError, tx.Commit().Error)
}

func getOverlapOfUsers(users1, users2 []*User) []*User {
	var res []*User
	for _, user1 := range users1 {
		for _, user2 := range users2 {
			if user1.ID == user2.ID {
				res = append(res, user1)
			}
		}
	}
	return res
}

func (s *Storage) UpdateWorkflowRecord(w *Workflow, tasks []*Task) error {
	instanceRecords := make([]*WorkflowInstanceRecord, len(tasks))
	for i, task := range tasks {
		instanceRecords[i] = &WorkflowInstanceRecord{
			TaskId:     task.ID,
			InstanceId: task.InstanceId,
		}
	}

	record := &WorkflowRecord{
		InstanceRecords: instanceRecords,
	}
	steps := w.cloneWorkflowStep()

	tx := s.db.Begin()
	err := tx.Save(record).Error
	if err != nil {
		tx.Rollback()
		return errors.New(errors.ConnectStorageError, err)
	}

	for _, step := range steps {
		currentStep := step
		currentStep.WorkflowRecordId = record.ID
		users := currentStep.Assignees
		currentStep.Assignees = nil
		err = tx.Save(currentStep).Error
		if err != nil {
			tx.Rollback()
			return errors.New(errors.ConnectStorageError, err)
		}
		err = tx.Model(currentStep).Association("Assignees").Replace(users).Error
		if err != nil {
			tx.Rollback()
			return errors.New(errors.ConnectStorageError, err)
		}
	}
	if len(steps) > 0 {
		err = tx.Model(record).Update("current_workflow_step_id", steps[0].ID).Error
		if err != nil {
			tx.Rollback()
			return errors.New(errors.ConnectStorageError, err)
		}
	}
	// update record history
	err = tx.Exec("INSERT INTO workflow_record_history (workflow_record_id, workflow_id) value (?, ?)",
		w.Record.ID, w.ID).Error
	if err != nil {
		tx.Rollback()
		return errors.New(errors.ConnectStorageError, err)
	}

	// update workflow record to new
	if err := tx.Model(&Workflow{}).Where("id = ?", w.ID).
		Update("workflow_record_id", record.ID).Error; err != nil {
		tx.Rollback()
		return errors.New(errors.ConnectStorageError, err)
	}

	return errors.New(errors.ConnectStorageError, tx.Commit().Error)
}

func (s *Storage) UpdateWorkflowStatus(w *Workflow, operateStep *WorkflowStep, instanceRecords []*WorkflowInstanceRecord) error {
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

		if len(instanceRecords) <= 0 {
			return nil
		}
		for _, inst := range instanceRecords {
			_, err = tx.Exec("UPDATE workflow_instance_records SET is_sql_executed = ? WHERE id = ?",
				inst.IsSQLExecuted, inst.ID)
			if err != nil {
				return err
			}
			_, err = tx.Exec("UPDATE workflow_instance_records SET execution_user_id = ? WHERE id = ?",
				inst.ExecutionUserId, inst.ID)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Storage) UpdateWorkflowRecordByID(id uint, workFlow map[string]interface{}) error {
	return s.db.Model(&WorkflowRecord{}).Where("id = ?", id).Updates(workFlow).Error
}

func (s *Storage) UpdateInstanceRecordSchedule(ir *WorkflowInstanceRecord, userId uint, scheduleTime *time.Time) error {
	err := s.db.Model(&WorkflowInstanceRecord{}).Where("id = ?", ir.ID).Update(map[string]interface{}{
		"scheduled_at":     scheduleTime,
		"schedule_user_id": userId,
	}).Error
	return errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) getWorkflowStepsByRecordIds(ids []uint) ([]*WorkflowStep, error) {
	steps := []*WorkflowStep{}
	err := s.db.Where("workflow_record_id in (?)", ids).
		Preload("Assignees").
		Preload("OperationUser").Find(&steps).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	stepTemplateIds := make([]uint, 0, len(steps))
	for _, step := range steps {
		stepTemplateIds = append(stepTemplateIds, step.WorkflowStepTemplateId)
	}
	stepTemplates := []*WorkflowStepTemplate{}
	err = s.db.Where("id in (?)", stepTemplateIds).Find(&stepTemplates).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	for _, step := range steps {
		for _, stepTemplate := range stepTemplates {
			if step.WorkflowStepTemplateId == stepTemplate.ID {
				step.Template = stepTemplate
			}
		}
	}
	return steps, nil
}

func (s *Storage) getWorkflowInstanceRecordsByRecordId(id uint) ([]*WorkflowInstanceRecord, error) {
	instanceRecords := []*WorkflowInstanceRecord{}
	err := s.db.Where("workflow_record_id = ?", id).
		Find(&instanceRecords).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	return instanceRecords, nil
}

func (s *Storage) GetWorkflowDetailById(id string) (*Workflow, bool, error) {
	workflow := &Workflow{}
	err := s.db.Preload("CreateUser", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("Record").
		Where("id = ?", id).First(workflow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	if workflow.Record == nil {
		return nil, false, errors.New(errors.DataConflict, fmt.Errorf("workflow record not exist"))
	}

	instanceRecords, err := s.getWorkflowInstanceRecordsByRecordId(workflow.Record.ID)
	if err != nil {
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	workflow.Record.InstanceRecords = instanceRecords

	steps, err := s.getWorkflowStepsByRecordIds([]uint{workflow.Record.ID})
	if err != nil {
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	workflow.Record.Steps = steps
	for _, step := range steps {
		if step.ID == workflow.Record.CurrentWorkflowStepId {
			workflow.Record.CurrentStep = step
		}
	}
	return workflow, true, nil
}

func (s *Storage) GetWorkflowHistoryById(id string) ([]*WorkflowRecord, error) {
	records := []*WorkflowRecord{}
	err := s.db.Model(&WorkflowRecord{}).Select("workflow_records.*").
		Joins("JOIN workflow_record_history AS wrh ON workflow_records.id = wrh.workflow_record_id").
		Where("wrh.workflow_id = ?", id).Scan(&records).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	if len(records) == 0 {
		return records, nil
	}
	recordIds := make([]uint, 0, len(records))
	for _, record := range records {
		recordIds = append(recordIds, record.ID)
	}
	steps, err := s.getWorkflowStepsByRecordIds(recordIds)
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	for _, record := range records {
		record.Steps = []*WorkflowStep{}
		for _, step := range steps {
			if step.WorkflowRecordId == record.ID && step.State != WorkflowStepStateInit {
				record.Steps = append(record.Steps, step)
			}
		}
	}
	return records, nil
}

func (s *Storage) GetWorkflowRecordCountByTaskIds(ids []uint) (uint32, error) {
	var count uint32
	err := s.db.Model(&WorkflowInstanceRecord{}).Where("workflow_instance_records.task_id IN (?)", ids).Count(&count).Error
	if err != nil {
		return 0, errors.New(errors.ConnectStorageError, err)
	}
	return count, nil
}

func (s *Storage) GetWorkflowByTaskId(id uint) (*Workflow, bool, error) {
	workflow := &Workflow{}
	err := s.db.Model(&Workflow{}).Select("workflows.id").
		Joins("LEFT JOIN workflow_records AS wr ON "+
			"workflows.workflow_record_id = wr.id").
		Joins("LEFT JOIN workflow_record_history ON "+
			"workflows.id = workflow_record_history.workflow_id").
		Joins("LEFT JOIN workflow_records AS h_wr ON "+
			"workflow_record_history.workflow_record_id = h_wr.id").
		Joins("LEFT JOIN workflow_instance_records AS wir ON "+
			"wir.workflow_record_id = wr.id").
		Joins("LEFT JOIN workflow_instance_records AS h_wir ON "+
			"h_wir.workflow_record_id = workflow_record_history.workflow_record_id").
		Where("wir.task_id = ? OR h_wir.task_id = ? AND workflows.id IS NOT NULL", id, id).
		Limit(1).Group("workflows.id").Scan(workflow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errors.New(errors.ConnectStorageError, err)
	}
	return workflow, true, nil
}

func (s *Storage) GetLastWorkflow() (*Workflow, bool, error) {
	workflow := new(Workflow)
	err := s.db.Last(workflow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return workflow, true, errors.New(errors.ConnectStorageError, err)
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
		_, err = tx.Exec("DELETE FROM workflow_steps WHERE workflow_id = ?", workflow.ID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM workflow_record_history WHERE workflow_id = ?", workflow.ID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("DELETE FROM workflow_instance_records WHERE workflow_record_id = ?", workflow.WorkflowRecordId)
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
		Where("workflows.created_at < ? "+
			"AND (workflow_records.status = 'finished' "+
			"OR workflow_records.status = 'canceled' "+
			"OR workflow_records.status IS NULL)", start).
		Scan(&workflows).Error
	return workflows, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetNeedScheduledWorkflows() ([]*Workflow, error) {
	workflows := []*Workflow{}
	err := s.db.Model(&Workflow{}).Select("workflows.id, workflows.workflow_record_id").
		Joins("LEFT JOIN workflow_records ON workflows.workflow_record_id = workflow_records.id").
		Joins("LEFT JOIN workflow_instance_records ON workflow_records.id = workflow_instance_records.workflow_record_id").
		Where("workflow_records.status = 'wait_for_execution' "+
			"AND workflow_instance_records.scheduled_at IS NOT NULL "+
			"AND workflow_instance_records.scheduled_at <= ? "+
			"AND workflow_instance_records.is_sql_executed = false", time.Now()).
		Scan(&workflows).Error
	return workflows, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetWorkflowBySubject(subject string) (*Workflow, bool, error) {
	workflow := &Workflow{Subject: subject}
	err := s.db.Where(*workflow).First(workflow).Error
	if err == gorm.ErrRecordNotFound {
		return workflow, false, nil
	}
	return workflow, true, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) TaskWorkflowIsUnfinished(taskIds []uint) (bool, error) {
	count := 0
	err := s.db.Table("workflow_records").
		Joins("LEFT JOIN workflow_instance_records ON workflow_records.id = workflow_instance_records.workflow_record_id").
		Where("workflow_records.status = ? OR workflow_records.status = ?", WorkflowStatusWaitForAudit, WorkflowStatusWaitForExecution).
		Where("workflow_instance_records.task_id IN (?)", taskIds).
		Count(&count).Error
	return count > 0, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetInstancesByWorkflowID(workflowID uint) ([]*Instance, error) {
	query := `
SELECT instances.id ,instances.maintenance_period
FROM workflows AS w
LEFT JOIN workflow_records AS wr ON wr.id = w.workflow_record_id
LEFT JOIN workflow_instance_records AS wir ON wr.id = wir.workflow_record_id
LEFT JOIN instances ON instances.id = wir.instance_id
WHERE 
w.id = ?`
	instances := []*Instance{}
	err := s.db.Raw(query, workflowID).Scan(&instances).Error
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}
	return instances, err
}

// GetWorkFlowStepIdsHasAudit 返回走完所有审核流程的workflow_steps的id
// 返回以workflow_record_id为分组的倒数第二条记录的workflow_steps.id
// 如果存在多个工单审核流程，workflow_record_id为分组的倒数第二条记录仍然是判断审核流程是否结束的依据
// 如果不存在工单审核流程，LIMIT 1 offset 1 会将workflow过滤掉
// 每个workflow_record_id对应一个workflows表中的一条记录，返回的id数组可以作为工单数量统计的依据
func (s *Storage) GetWorkFlowStepIdsHasAudit() ([]uint, error) {
	workFlowStepsByIndexAndState, err := s.GetWorkFlowReverseStepsByIndexAndState(1, WorkflowStepStateApprove)
	if err != nil {
		return nil, errors.ConnectStorageErrWrapper(err)
	}

	ids := make([]uint, 0)
	for _, workflowStep := range workFlowStepsByIndexAndState {
		ids = append(ids, workflowStep.ID)
	}

	return ids, nil
}

func (s *Storage) GetDurationMinHasAudit(ids []uint) (int, error) {
	type minStruct struct {
		Min int `json:"min"`
	}

	var result minStruct
	err := s.db.Model(&Workflow{}).
		Select("sum(timestampdiff(minute, workflows.created_at, workflow_steps.operate_at)) as min").
		Joins("LEFT JOIN workflow_steps ON workflow_steps.workflow_record_id = workflows.workflow_record_id").
		Where("workflow_steps.id IN (?)", ids).Scan(&result).Error

	return result.Min, errors.ConnectStorageErrWrapper(err)
}

// WorkFlowStepsBO BO是business object的缩写，表示业务对象
type WorkFlowStepsBO struct {
	ID         uint
	OperateAt  *time.Time
	WorkflowId uint
}

// GetWorkFlowReverseStepsByIndexAndState 返回以workflow_id为分组的倒数第index个记录
func (s *Storage) GetWorkFlowReverseStepsByIndexAndState(index int, state string) ([]*WorkFlowStepsBO, error) {
	query := fmt.Sprintf(`SELECT id,operate_at,workflow_id
FROM workflow_steps a
WHERE a.id =
      (SELECT id
       FROM workflow_steps
       WHERE workflow_id = a.workflow_id
       ORDER BY id desc
       limit 1 offset %d)
  and a.state = '%s';`, index, state)

	workflowStepsBO := make([]*WorkFlowStepsBO, 0)
	return workflowStepsBO, s.db.Raw(query).Scan(&workflowStepsBO).Error
}

func (s *Storage) GetWorkflowCountByStepType(stepTypes []string) (int, error) {
	if len(stepTypes) == 0 {
		return 0, nil
	}

	var count int
	err := s.db.Table("workflows").
		Joins("left join workflow_records on workflows.workflow_record_id = workflow_records.id").
		Joins("left join workflow_steps on workflow_records.current_workflow_step_id = workflow_steps.id").
		Joins("left join workflow_step_templates on workflow_steps.workflow_step_template_id = workflow_step_templates.id ").
		Where("workflow_step_templates.type in (?)", stepTypes).
		Count(&count).Error

	return count, errors.New(errors.ConnectStorageError, err)
}

func (s *Storage) GetWorkflowCountByStatus(status string) (int, error) {
	var count int
	err := s.db.Table("workflows").
		Joins("left join workflow_records on workflows.workflow_record_id = workflow_records.id").
		Where("workflow_records.status = ?", status).
		Count(&count).Error
	if err != nil {
		return 0, errors.New(errors.ConnectStorageError, err)
	}

	return count, nil
}

// GetApprovedWorkflowCount 将会返回未被回收且审核流程全部通过的工单数, 包括上线成功(失败)的工单和等待上线的工单, 不包括关闭的工单
func (s *Storage) GetApprovedWorkflowCount() (int, error) {
	query := `
	select count(1) as count
from workflows 
left join workflow_records on workflows.workflow_record_id = workflow_records.id
left join workflow_steps on workflow_records.current_workflow_step_id = workflow_steps.id 
left join workflow_step_templates on workflow_steps.workflow_step_template_id = workflow_step_templates.id 
where 
workflow_records.status = 'finished'
or
workflow_step_templates.type = 'sql_execute';
`
	var count = struct {
		Count int `json:"count"`
	}{}
	return count.Count, errors.New(errors.ConnectStorageError, s.db.Raw(query).Scan(&count).Error)
}

func (s *Storage) GetAllWorkflowCount() (int, error) {
	var count int
	return count, errors.New(errors.ConnectStorageError, s.db.Model(&Workflow{}).Count(&count).Error)
}

func (s *Storage) GetWorkflowCountByTaskStatus(status []string) (int, error) {
	//if len(status) == 0 {
	//	return 0, nil
	//}
	//
	//var count int
	//err := s.db.Table("workflows").
	//	Joins("left join workflow_records on workflows.workflow_record_id = workflow_records.id").
	//	Joins("left join tasks on workflow_records.task_id = tasks.id").
	//	Where("tasks.status in (?)", status).
	//	Count(&count).Error
	//
	//return count, errors.New(errors.ConnectStorageError, err)
	// todo issue901
	return 0, nil
}

func (s *Storage) GetWorkFlowCountBetweenStartTimeAndEndTime(startTime, endTime time.Time) (int64, error) {
	var count int64
	return count, s.db.Model(&Workflow{}).Where("created_at BETWEEN ? and ?", startTime, endTime).Count(&count).Error
}

type DailyWorkflowCount struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func (s *Storage) GetWorkflowDailyCountBetweenStartTimeAndEndTime(startTime, endTime time.Time) ([]*DailyWorkflowCount, error) {
	var counts []*DailyWorkflowCount
	err := s.db.Table("workflows").
		Select("cast(created_at as date) as date, count(*) as count").
		Where("created_at BETWEEN ? and ?", startTime, endTime).
		Group("cast(created_at as date)").Find(&counts).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}
	return counts, nil
}

func (s *Storage) GetWaitExecInstancesCountByWorkflowId(workflowId uint) (int, error) {
	count := 0
	err := s.db.Table("workflows").
		Joins("LEFT JOIN workflow_records ON workflow_records.id = workflows.workflow_record_id").
		Joins("LEFT JOIN workflow_instance_records ON workflow_records.id = workflow_instance_records.workflow_record_id").
		Where("workflows.id = ?", workflowId).
		Where("workflow_instance_records.is_sql_executed = false").
		Count(&count).Error
	if err != nil {
		return 0, errors.New(errors.ConnectStorageError, err)
	}
	return count, nil
}

type WorkflowTasksSummaryDetail struct {
	WorkflowRecordStatus      string     `json:"workflow_record_status"`
	TaskId                    uint       `json:"task_id"`
	TaskExecStartAt           *time.Time `json:"task_exec_start_at"`
	TaskExecEndAt             *time.Time `json:"task_exec_end_at"`
	TaskPassRate              float64    `json:"task_pass_rate"`
	TaskScore                 int32      `json:"task_score"`
	TaskStatus                string     `json:"task_status"`
	InstanceName              string     `json:"instance_name"`
	InstanceDeletedAt         *time.Time `json:"instance_deleted_at"`
	InstanceMaintenancePeriod Periods    `json:"instance_maintenance_period" gorm:"text"`
	InstanceScheduledAt       *time.Time `json:"instance_scheduled_at"`
	ExecutionUserDeletedAt    *time.Time `json:"execution_user_deleted_at"`
	ExecutionUserName         string     `json:"execution_user_name"`
	CurrentStepAssigneeUsers  RowList    `json:"current_step_assignee_users"`
}

var workflowTasksSummaryQueryTpl = `
SELECT wr.status                                                     AS workflow_record_status,
       tasks.id                                                      AS task_id,
       tasks.exec_start_at                                           AS task_exec_start_at,
       tasks.exec_end_at                                             AS task_exec_end_at,
       tasks.pass_rate                                               AS task_pass_rate,
       tasks.score                                                   AS task_score,
       tasks.status                                                  AS task_status,
       inst.name                                                     AS instance_name,
       inst.deleted_at                                               AS instance_deleted_at,
       inst.maintenance_period                                       AS instance_maintenance_period,
       wir.scheduled_at                                              AS instance_scheduled_at,
       exec_user.deleted_at                                          AS execution_user_deleted_at,
       COALESCE(exec_user.login_name, '')                            AS execution_user_name,
       GROUP_CONCAT(DISTINCT COALESCE(curr_ass_user.login_name, '')) AS current_step_assignee_users

{{- template "body" . -}}
GROUP BY tasks.id, wir.id
`

var workflowTasksSummaryQueryBodyTpl = `
{{ define "body" }}
FROM workflow_instance_records AS wir
LEFT JOIN workflow_records AS wr ON wir.workflow_record_id = wr.id
LEFT JOIN workflows AS w ON w.workflow_record_id = wr.id
LEFT JOIN users AS exec_user ON wir.execution_user_id = exec_user.id
LEFT JOIN tasks ON wir.task_id = tasks.id
LEFT JOIN instances AS inst ON tasks.instance_id = inst.id
LEFT JOIN workflow_steps AS curr_ws ON wr.current_workflow_step_id = curr_ws.id
LEFT JOIN workflow_step_user AS curr_ws_user ON curr_ws.id = curr_ws_user.workflow_step_id
LEFT JOIN users AS curr_ass_user ON curr_ws_user.user_id = curr_ass_user.id

WHERE
w.deleted_at IS NULL
AND w.id = :workflow_id

{{ end }}
`

func (s *Storage) GetWorkflowTasksSummaryByReq(data map[string]interface{}) (
	result []*WorkflowTasksSummaryDetail, err error) {

	err = s.getListResult(workflowTasksSummaryQueryBodyTpl, workflowTasksSummaryQueryTpl, data, &result)
	if err != nil {
		return result, errors.New(errors.ConnectStorageError, err)
	}

	return result, nil
}

func (s *Storage) GetTasksByWorkFlowRecordID(id uint) ([]*Task, error) {
	var tasks []*Task
	err := s.db.Model(&WorkflowInstanceRecord{}).Select("tasks.status").
		Joins("left join tasks on tasks.id = workflow_instance_records.task_id").
		Where("workflow_instance_records.workflow_record_id = ?", id).Scan(&tasks).Error
	if err != nil {
		return nil, errors.New(errors.ConnectStorageError, err)
	}

	return tasks, nil
}
