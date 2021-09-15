package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"

	"github.com/sirupsen/logrus"
)

const (
	SqlAuditTaskExpiredTime = 3 * 24 // 3 days
)

func (s *Sqled) cleanLoop() {
	tick := time.Tick(1 * time.Hour)
	entry := log.NewEntry().WithField("type", "cron")
	s.CleanExpiredWorkflows(entry)
	s.CleanExpiredTasks(entry)
	for {
		select {
		case <-s.exit:
			return
		case <-tick:
			s.CleanExpiredWorkflows(entry)
			s.CleanExpiredTasks(entry)
		}
	}
}

func (s *Sqled) CleanExpiredWorkflows(entry *logrus.Entry) {
	st := model.GetStorage()

	expiredHours, err := st.GetWorkflowExpiredHoursOrDefault()
	if err != nil {
		entry.Errorf("get workflow expired hours error: %v", err)
		return
	}

	start := time.Now().Add(-expiredHours * time.Hour)
	workflows, err := st.GetExpiredWorkflows(start)
	if err != nil {
		entry.Errorf("get workflows from storage error: %v", err)
		return
	}
	hasDeletedWorkflowIds := make([]string, 0, len(workflows))
	for _, workflow := range workflows {
		err := st.DeleteWorkflow(workflow)
		if err != nil {
			entry.Errorf("clean workflow %d error: %s", workflow.ID, err)
			break
		}
		hasDeletedWorkflowIds = append(hasDeletedWorkflowIds, strconv.FormatUint(uint64(workflow.ID), 10))
	}
	if len(hasDeletedWorkflowIds) > 0 {
		entry.Infof("clean workflow [%s] success", strings.Join(hasDeletedWorkflowIds, ", "))
	}
}

func (s *Sqled) CleanExpiredTasks(entry *logrus.Entry) {
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
