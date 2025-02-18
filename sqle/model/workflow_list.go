package model

import (
	"database/sql"
	"time"
)

type WorkflowListDetail struct {
	ProjectId                  string         `json:"project_id"`
	Subject                    string         `json:"subject"`
	WorkflowId                 string         `json:"workflow_id"`
	Desc                       string         `json:"desc"`
	CreateUser                 sql.NullString `json:"create_user_id"`
	CreateUserDeletedAt        *time.Time     `json:"create_user_deleted_at"`
	CreateTime                 *time.Time     `json:"create_time"`
	CurrentStepType            sql.NullString `json:"current_step_type" enums:"sql_review,sql_execute"`
	CurrentStepAssigneeUserIds sql.NullString `json:"current_step_assignee_user_id_list"`
	Status                     string         `json:"status"`
	TaskInstanceType           RowList        `json:"task_instance_type"` // unused
	SqlVersionNames            RowList        `json:"versions"`
	InstanceIds                RowList        `json:"instance_ids"`
}

var workflowsQueryTpl = `
SELECT 
	   w.project_id,
       w.subject,
       w.workflow_id,
       w.desc,
       w.create_user_id,
	   CAST("" AS DATETIME)											 AS create_user_deleted_at,
       w.created_at                                                  AS create_time,
       curr_wst.type                                                 AS current_step_type,
       curr_ws.assignees											 AS current_step_assignee_user_id_list,
       wr.status,
	   GROUP_CONCAT(DISTINCT wir.instance_id SEPARATOR ',') AS instance_ids,
	   GROUP_CONCAT(DISTINCT versions.version SEPARATOR ',') AS versions
{{- template "body" . -}}
GROUP BY w.id
{{- if .filter_instance_id }}
HAVING 
    FIND_IN_SET(:filter_instance_id, instance_ids) > 0
{{- end }}

ORDER BY w.id DESC
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var workflowsCountTpl = `SELECT COUNT(DISTINCT w.id)

{{- template "body" . -}}
`

var exportWorkflowIDListTpl = `
SELECT w.workflow_id AS workflow_id
{{- template "body" . -}}
GROUP BY w.id
ORDER BY w.id DESC
`

var workflowsQueryBodyTpl = `
{{ define "body" }}
FROM workflows w
LEFT JOIN workflow_records AS wr ON w.workflow_record_id = wr.id
LEFT JOIN workflow_instance_records wir on wir.workflow_record_id = wr.id
LEFT JOIN tasks ON wir.task_id = tasks.id
LEFT JOIN workflow_steps AS curr_ws ON wr.current_workflow_step_id = curr_ws.id
LEFT JOIN workflow_step_templates AS curr_wst ON curr_ws.workflow_step_template_id = curr_wst.id
LEFT JOIN workflow_version_stages AS stages ON stages.workflow_id = w.workflow_id
LEFT JOIN sql_versions AS versions ON stages.sql_version_id = versions.id AND versions.deleted_at IS NULL

{{- if .check_user_can_access }}
LEFT JOIN workflow_steps AS all_ws ON w.workflow_id = all_ws.workflow_id AND all_ws.state !='initialized'
LEFT JOIN workflow_step_templates AS all_wst ON all_ws.workflow_step_template_id = all_wst.id
{{- end }}
WHERE
w.deleted_at IS NULL 

{{- if .check_user_can_access }}
AND (
w.create_user_id = :current_user_id 
OR curr_ws.assignees REGEXP :current_user_id
OR all_ws.operation_user_id = :current_user_id
OR wir.execution_assignees = :current_user_id
OR IF(wr.status = 'wait_for_execution'
				, wir.execution_assignees REGEXP :current_user_id
				, '')

{{- if .viewable_instance_ids }} 
OR tasks.instance_id IN ( {{ .viewable_instance_ids }})
{{- end }}

)
{{- end }}

{{- if .filter_subject }}
AND w.subject = :filter_subject
{{- end }}

{{- if .filter_create_time_from }}
AND w.created_at > :filter_create_time_from
{{- end }}

{{- if .filter_create_time_to }}
AND w.created_at < :filter_create_time_to
{{- end }}

{{- if .filter_task_execute_start_time_from }}
AND tasks.exec_start_at > :filter_task_execute_start_time_from
{{- end }}

{{- if .filter_task_execute_start_time_to }}
AND tasks.exec_start_at < :filter_task_execute_start_time_to
{{- end }}

{{- if .filter_create_user_id }}
AND w.create_user_id = :filter_create_user_id
{{- end }}

{{- if .filter_current_step_type }}
AND curr_wst.type = :filter_current_step_type
{{- end }}

{{- if .filter_status }}
AND wr.status IN (:filter_status)
{{- end }}

{{- if .filter_current_step_assignee_user_id }}
AND (curr_ws.assignees REGEXP :filter_current_step_assignee_user_id
OR IF(wr.status = 'wait_for_execution'
                , wir.execution_assignees REGEXP :filter_current_step_assignee_user_id
                , '')
)
{{- end }}

{{- if .filter_task_status }}
AND tasks.status = :filter_task_status
{{- end }}

{{- if .filter_task_instance_id }}
AND tasks.instance_id = :filter_task_instance_id
{{- end }}

{{- if .filter_workflow_id }}
AND w.workflow_id = :filter_workflow_id
{{- end }}

{{- if .filter_sql_version_id }}
AND versions.id = :filter_sql_version_id
{{- end }}

{{- if .filter_project_id }}
AND w.project_id = :filter_project_id
{{- end }}
{{- if .filter_status_list }}
AND wr.status IN ( 
	{{ range $index, $element := .filter_status_list }}
		{{ if $index }},{{ end }}"{{ $element }}"
	{{ end }}
 )
{{- end }}

{{- if .filter_project_id_list }}
AND w.project_id IN ( 
	{{ range $index, $element := .filter_project_id_list }}
		{{ if $index }},{{ end }}"{{ $element }}"
	{{ end }}
 )
{{- end }}

{{- if .fuzzy_keyword }}
AND (w.subject like :fuzzy_keyword or w.workflow_id like :fuzzy_keyword or w.desc like :fuzzy_keyword)
{{- end }}

{{ end }}

`

func (s *Storage) GetWorkflowsByReq(data map[string]interface{}) (
	result []*WorkflowListDetail, count uint64, err error) {

	err = s.getListResult(workflowsQueryBodyTpl, workflowsQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}

	count, err = s.getCountResult(workflowsQueryBodyTpl, workflowsCountTpl, data)

	return result, count, err
}

func (s *Storage) GetGlobalWorkflowTotalNum(data map[string]interface{}) (count uint64, err error) {
	return s.getCountResult(workflowsQueryBodyTpl, workflowsCountTpl, data)
}

func (s *Storage) GetWorkflowCountByReq(data map[string]interface{}) (uint64, error) {
	return s.getCountResult(workflowsQueryBodyTpl, workflowsCountTpl, data)
}

// func (s *Storage) GetWorkflowTotalByProjectId(projectId string) (uint64, error) {
// 	data := map[string]interface{}{
// 		"filter_project_id": projectId,
// 	}
// 	return s.GetWorkflowCountByReq(data)
// }

// dms-todo: using project id as name, 临时方案
var projectWorkflowCountTpl = `
SELECT w.project_id AS project_name, COUNT(DISTINCT w.id) AS workflow_count
{{- template "body" . -}}
GROUP BY p.name
`

type ProjectWorkflowCount struct {
	ProjectName   string `json:"project_name"`
	WorkflowCount int    `json:"workflow_count"`
}

func (s *Storage) GetWorkflowCountForDashboardProjectTipsByReq(data map[string]interface{}) (
	result []*ProjectWorkflowCount, err error) {
	err = s.getTemplateQueryResult(data, &result, workflowsQueryBodyTpl, projectWorkflowCountTpl)
	if err != nil {
		return result, err
	}
	return result, nil
}

type WorkflowTemplateDetail struct {
	Name string `json:"workflow_template_name"`
	Desc string `json:"desc"`
}

var workflowTemplatesQueryTpl = `SELECT workflow_templates.name AS workflow_template_name, workflow_templates.desc

{{- template "body" . -}}

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var workflowTemplatesCountTpl = `SELECT COUNT(*)

{{- template "body" . -}}
`

var workflowTemplatesQueryBodyTpl = `
{{ define "body" }}
FROM workflow_templates
WHERE
workflow_templates.deleted_at IS NULL
{{- end }}
`

func (s *Storage) GetWorkflowTemplatesByReq(data map[string]interface{}) (
	result []*WorkflowTemplateDetail, count uint64, err error) {

	err = s.getListResult(workflowTemplatesQueryBodyTpl, workflowTemplatesQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(workflowTemplatesQueryBodyTpl, workflowTemplatesCountTpl, data)
	return result, count, err
}

func (s *Storage) GetExportWorkflowIDListByReq(data map[string]interface{}, user *User) (idList []string, err error) {
	err = s.getListResult(workflowsQueryBodyTpl, exportWorkflowIDListTpl, data, &idList)
	if err != nil {
		return idList, err
	}

	return idList, nil
}
