package dms

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

// DataExportWorkflowDetail is a minimal view of GET .../data_export_workflows/{uid}.
type DataExportWorkflowDetail struct {
	WorkflowUID string
	Tasks       []DataExportTaskRef
	Steps       []DataExportWorkflowStep
}

// DataExportTaskRef holds a task uid from the workflow record.
type DataExportTaskRef struct {
	UID string
}

// DataExportWorkflowStep is one workflow step from DMS detail.
type DataExportWorkflowStep struct {
	Number            uint64
	Type              string
	OperationUserName string
	OperationTime     *time.Time
	State             string
}

// DataExportTaskDetail is one task from BatchGetDataExportTask.
type DataExportTaskDetail struct {
	TaskUID         string
	Status          string
	ExportStartTime *time.Time
	ExportEndTime   *time.Time
}

// DataExportTaskSQL is one SQL row from ListDataExportTaskSQLs.
type DataExportTaskSQL struct {
	SQL          string
	ExportResult string
}

type getDataExportWorkflowReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		WorkflowUID    string `json:"workflow_uid"`
		WorkflowRecord struct {
			Tasks []struct {
				UID string `json:"task_uid"`
			} `json:"tasks"`
			Steps []struct {
				Number        uint64 `json:"number"`
				Type          string `json:"type"`
				OperationUser struct {
					Name string `json:"name"`
				} `json:"operation_user"`
				OperationTime *time.Time `json:"operation_time"`
				State         string     `json:"state"`
			} `json:"workflow_step_list"`
		} `json:"workflow_record"`
		WorkflowRecordHistory []struct {
			Steps []struct {
				Number        uint64 `json:"number"`
				Type          string `json:"type"`
				OperationUser struct {
					Name string `json:"name"`
				} `json:"operation_user"`
				OperationTime *time.Time `json:"operation_time"`
				State         string     `json:"state"`
			} `json:"workflow_step_list"`
		} `json:"workflow_record_history"`
	} `json:"data"`
}

type batchGetDataExportTaskReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		TaskUID         string     `json:"task_uid"`
		Status          string     `json:"status"`
		ExportStartTime *time.Time `json:"export_start_time"`
		ExportEndTime   *time.Time `json:"export_end_time"`
	} `json:"data"`
}

type listDataExportTaskSQLsReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		SQL          string `json:"sql"`
		ExportResult string `json:"export_result"`
	} `json:"data"`
	Total int64 `json:"total_nums"`
}

// GetDataExportWorkflowDetail calls existing project-scoped DMS API (no new DMS contract).
func GetDataExportWorkflowDetail(ctx context.Context, dmsAddr, projectUID, workflowUID string) (*DataExportWorkflowDetail, error) {
	if dmsAddr == "" || projectUID == "" || workflowUID == "" {
		return nil, nil
	}
	header := map[string]string{"Authorization": pkgHttp.DefaultDMSToken}
	reply := &getDataExportWorkflowReply{}
	urlStr := fmt.Sprintf("%s/v1/dms/projects/%s/data_export_workflows/%s", dmsAddr, projectUID, workflowUID)
	if err := pkgHttp.Get(ctx, urlStr, header, nil, reply); err != nil {
		return nil, fmt.Errorf("get data export workflow: %v", err)
	}
	if reply.Code != 0 {
		return nil, fmt.Errorf("get data export workflow code(%v): %v", reply.Code, reply.Message)
	}
	if reply.Data == nil {
		return nil, nil
	}
	out := &DataExportWorkflowDetail{WorkflowUID: reply.Data.WorkflowUID}
	for _, t := range reply.Data.WorkflowRecord.Tasks {
		if t.UID != "" {
			out.Tasks = append(out.Tasks, DataExportTaskRef{UID: t.UID})
		}
	}
	appendSteps := func(steps []struct {
		Number        uint64 `json:"number"`
		Type          string `json:"type"`
		OperationUser struct {
			Name string `json:"name"`
		} `json:"operation_user"`
		OperationTime *time.Time `json:"operation_time"`
		State         string     `json:"state"`
	}) {
		for _, s := range steps {
			out.Steps = append(out.Steps, DataExportWorkflowStep{
				Number:            s.Number,
				Type:              s.Type,
				OperationUserName: s.OperationUser.Name,
				OperationTime:     s.OperationTime,
				State:             s.State,
			})
		}
	}
	appendSteps(reply.Data.WorkflowRecord.Steps)
	// Prefer current record steps; if empty, fall back to latest history record steps.
	if len(out.Steps) == 0 {
		for i := len(reply.Data.WorkflowRecordHistory) - 1; i >= 0; i-- {
			appendSteps(reply.Data.WorkflowRecordHistory[i].Steps)
			if len(out.Steps) > 0 {
				break
			}
		}
	}
	return out, nil
}

// BatchGetDataExportTasks calls existing BatchGetDataExportTask API.
func BatchGetDataExportTasks(ctx context.Context, dmsAddr, projectUID string, taskUIDs []string) ([]DataExportTaskDetail, error) {
	if dmsAddr == "" || projectUID == "" || len(taskUIDs) == 0 {
		return nil, nil
	}
	header := map[string]string{"Authorization": pkgHttp.DefaultDMSToken}
	reply := &batchGetDataExportTaskReply{}
	q := url.Values{}
	q.Set("data_export_task_uids", strings.Join(taskUIDs, ","))
	urlStr := fmt.Sprintf("%s/v1/dms/projects/%s/data_export_tasks?%s", dmsAddr, projectUID, q.Encode())
	if err := pkgHttp.Get(ctx, urlStr, header, nil, reply); err != nil {
		return nil, fmt.Errorf("batch get data export tasks: %v", err)
	}
	if reply.Code != 0 {
		return nil, fmt.Errorf("batch get data export tasks code(%v): %v", reply.Code, reply.Message)
	}
	out := make([]DataExportTaskDetail, 0, len(reply.Data))
	for _, t := range reply.Data {
		out = append(out, DataExportTaskDetail{
			TaskUID:         t.TaskUID,
			Status:          t.Status,
			ExportStartTime: t.ExportStartTime,
			ExportEndTime:   t.ExportEndTime,
		})
	}
	return out, nil
}

// ListDataExportTaskSQLsAll pages through existing ListDataExportTaskSQLs API.
func ListDataExportTaskSQLsAll(ctx context.Context, dmsAddr, projectUID, taskUID string) ([]DataExportTaskSQL, error) {
	if dmsAddr == "" || projectUID == "" || taskUID == "" {
		return nil, nil
	}
	header := map[string]string{"Authorization": pkgHttp.DefaultDMSToken}
	const pageSize uint32 = 100
	out := make([]DataExportTaskSQL, 0)
	for page := uint32(1); ; page++ {
		reply := &listDataExportTaskSQLsReply{}
		urlStr := fmt.Sprintf(
			"%s/v1/dms/projects/%s/data_export_tasks/%s/data_export_task_sqls?page_size=%d&page_index=%d",
			dmsAddr, projectUID, taskUID, pageSize, page,
		)
		if err := pkgHttp.Get(ctx, urlStr, header, nil, reply); err != nil {
			return nil, fmt.Errorf("list data export task sqls: %v", err)
		}
		if reply.Code != 0 {
			return nil, fmt.Errorf("list data export task sqls code(%v): %v", reply.Code, reply.Message)
		}
		if len(reply.Data) == 0 {
			break
		}
		for _, row := range reply.Data {
			out = append(out, DataExportTaskSQL{SQL: row.SQL, ExportResult: row.ExportResult})
		}
		if uint64(len(out)) >= uint64(reply.Total) || len(reply.Data) < int(pageSize) {
			break
		}
	}
	return out, nil
}
