package server

import (
	"context"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/dms"
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
	// j.CleanExpiredWorkflows(entry) /* 不再自动销毁工单（目前没有使用场景）*/
	j.CleanExpiredTasks(entry)
	j.CleanExpiredOperationLog(entry)
}

func (j *CleanJob) CleanExpiredWorkflows(entry *logrus.Entry) {
	st := model.GetStorage()

	expiredHours, err := dms.GetWorkflowExpiredHoursOrDefault()
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
	operationRecordExpiredHours := getOperationRecordExpiredHours(j.entry)
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

	operationRecordExpiredHours = dms.DefaultOperationRecordExpiredHours
	systemVariables, err := dmsobject.GetSystemVariables(context.TODO(), dms.GetDMSServerAddress())
	if err != nil || systemVariables.Code != 0 {
		entry.Warnf("get system variables failed, err: %s", err.Error())
		return operationRecordExpiredHours
	}

	return systemVariables.Data.OperationRecordExpiredHours
}

type CleanJobForAllNodes struct {
	BaseJob
}

func NewCleanJobForAllNodes(entry *logrus.Entry) ServerJob {
	entry = entry.WithField("job", "clean_for_all_nodes")
	j := &CleanJobForAllNodes{}
	j.BaseJob = *NewBaseJob(entry, 1*time.Hour, j.job)
	return j
}

func (j *CleanJobForAllNodes) job(entry *logrus.Entry) {
	j.CleanUpExpiredFiles(entry)
}

func (j *CleanJobForAllNodes) CleanUpExpiredFiles(entry *logrus.Entry) {

	s := model.GetStorage()
	var files []model.AuditFile
	var err error
	// get expired file with no workflow (expired time 24h) in this machine
	files, err = s.GetExpiredFileWithNoWorkflow()
	if err != nil {
		entry.Errorf("get expired file with no workflow error: %v", err)
		return
	}
	// get expired file (expired time 7*24h) in this machine
	expiredFiles, err := s.GetExpiredFile()
	if err != nil {
		entry.Errorf("get expired files error: %v", err)
		return
	}

	files = append(files, expiredFiles...)
	var filePath string
	for _, file := range files {
		filePath = model.DefaultFilePath(file.UniqueName)
		// if file exist delete file
		if _, err = os.Stat(filePath); err == nil {
			err = os.Remove(filePath)
			if err != nil {
				entry.Warnf("remove audit file failed %v", err)
				continue
			}
			err = s.Delete(&file)
			if err != nil {
				entry.Warnf("remove audit file record failed %v", err)
				continue
			}
			entry.Infof("delete files with no workflow success, file path: %s", filePath)
			continue
		}
		// if file is not eixt delete file record
		if os.IsNotExist(err) {
			entry.Infof("file %s not exist, delete files record", filePath)
			err = s.Delete(&file)
			if err != nil {
				entry.Warnf("while file not exist, removing file record failed %v", err)
				continue
			}
		} else {
			entry.Errorf("when read stat of file %s, unexpected err %v", filePath, err)
		}
	}
}
