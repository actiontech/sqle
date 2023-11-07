package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/model"

	"github.com/sirupsen/logrus"
)

const (
	SqlAuditTaskExpiredTime = 3 * 24 // 3 days
)

type CleanJob struct {
	BaseJob
}

func NewCleanJob(entry *logrus.Entry) ServerJob {
	entry = entry.WithField("job", "clean")
	j := &CleanJob{}
	j.BaseJob = *NewBaseJob(entry, 1*time.Hour, j.job)
	return j
}

func (j *CleanJob) job(entry *logrus.Entry) {
	j.CleanExpiredWorkflows(entry)
	j.CleanExpiredTasks(entry)
	j.CleanExpiredOperationLog(entry)
}

func (j *CleanJob) CleanExpiredWorkflows(entry *logrus.Entry) {
	st := model.GetStorage()

	expiredHours, err := st.GetWorkflowExpiredHoursOrDefault()
	if err != nil {
		entry.Errorf("get workflow expired hours error: %v", err)
		return
	}

	start := time.Now().Add(time.Duration(-expiredHours * int64(time.Hour)))
	workflows, err := st.GetExpiredWorkflows(start)
	if err != nil {
		entry.Errorf("get workflows from storage error: %v", err)
		return
	}
	hasDeletedWorkflowIds := make([]string, 0, len(workflows))
	for _, workflow := range workflows {
		err := st.DeleteWorkflow(workflow)
		if err != nil {
			entry.Errorf("clean workflow %s error: %s", workflow.WorkflowId, err)
			break
		}
		hasDeletedWorkflowIds = append(hasDeletedWorkflowIds, workflow.WorkflowId)
	}
	if len(hasDeletedWorkflowIds) > 0 {
		entry.Infof("clean workflow [%s] success", strings.Join(hasDeletedWorkflowIds, ", "))
	}
}

func (j *CleanJob) CleanExpiredTasks(entry *logrus.Entry) {
	st := model.GetStorage()
	start := time.Now().Add(-SqlAuditTaskExpiredTime * time.Hour)
	tasks, err := st.GetExpiredTasks(start)
	if err != nil {
		entry.Errorf("get tasks from storage error: %v", err)
		return
	}
	hasDeletedTaskIds := make([]string, 0, len(tasks))
	for _, task := range tasks {
		err := st.DeleteTask(task)
		if err != nil {
			entry.Errorf("clean task %d error: %s", task.ID, err)
			break
		}
		hasDeletedTaskIds = append(hasDeletedTaskIds, strconv.FormatUint(uint64(task.ID), 10))
	}
	if len(hasDeletedTaskIds) > 0 {
		entry.Infof("clean task [%s] success", strings.Join(hasDeletedTaskIds, ", "))
	}
}

func (j *CleanJob) CleanExpiredOperationLog(entry *logrus.Entry) {

	st := model.GetStorage()
	operationRecordExpiredHours := getOperationRecordExpiredHours(st, j.entry)
	start := time.Now().Add(-time.Duration(operationRecordExpiredHours) * time.Hour)
	idList, err := st.GetExpiredOperationRecordIDListByStartTime(start)

	if err != nil {
		entry.Errorf("get expired operation record id list error: %v", err)
		return
	}

	if len(idList) > 0 {
		if err := st.DeleteExpiredOperationRecordByIDList(idList); err != nil {
			entry.Errorf("delete expired operation record error: %v", err)
			return
		}

		entry.Infof("delete expired operation record succeeded, count: %d id: %s", len(idList), strings.Join(idList, ","))
	}
}

func getOperationRecordExpiredHours(
	s *model.Storage, entry *logrus.Entry) (operationRecordExpiredHours int) {

	operationRecordExpiredHours = model.DefaultOperationRecordExpiredHours
	systemVariables, err := s.GetAllSystemVariables()
	if err != nil {
		entry.Warnf("get system variables failed, err: %s", err.Error())
		return operationRecordExpiredHours
	}
	strVal := systemVariables[model.SystemVariableOperationRecordExpiredHours].Value
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		entry.Warnf(
			"get system variables operation_record_expired_hours failed, err: %s",
			err.Error())
		return operationRecordExpiredHours
	}
	operationRecordExpiredHours = intVal

	return operationRecordExpiredHours
}
