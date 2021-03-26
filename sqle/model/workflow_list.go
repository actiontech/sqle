package model

import (
	"database/sql"
	"time"
)

type WorkflowListDetail struct {
	Id                      uint           `json:"workflow_id"`
	Subject                 string         `json:"subject"`
	Desc                    string         `json:"desc"`
	TaskPassRate            float64        `json:"task_pass_rate"`
	TaskInstance            string         `json:"task_instance_name"`
	TaskInstanceSchema      string         `json:"task_instance_schema"`
	TaskStatus              string         `json:"task_status"`
	CreateUser              string         `json:"create_user_name"`
	CreateTime              *time.Time     `json:"create_time"`
	CurrentStepType         sql.NullString `json:"current_step_type" enums:"sql_review, sql_execute"`
	CurrentStepAssigneeUser RowList        `json:"current_step_assignee_user_name_list"`
	Status                  string         `json:"status"`
}

var workflowsQueryTpl = `SELECT w.id AS workflow_id, w.subject, w.desc, wr.status,
tasks.status AS task_status, tasks.pass_rate AS task_pass_rate, tasks.status AS task_status, 
tasks.instance_schema AS task_instance_schema, inst.name AS task_instance_name, 
create_user.login_name AS create_user_name, w.created_at AS create_time, wst.type AS current_step_type, 
GROUP_CONCAT(DISTINCT COALESCE(ass_user.login_name,'')) AS current_step_assignee_user_name_list
FROM workflows AS w
LEFT JOIN tasks ON w.task_id = tasks.id
LEFT JOIN instances AS inst ON tasks.instance_id = inst.id
LEFT JOIN users AS create_user ON w.create_user_id = create_user.id
LEFT JOIN workflow_records AS wr ON w.workflow_record_id = wr.id
LEFT JOIN workflow_steps AS ws ON wr.current_workflow_step_id = ws.id
LEFT JOIN workflow_step_templates AS wst ON ws.workflow_step_template_id = wst.id
LEFT JOIN workflow_step_template_user AS wst_re_user ON wst.id = wst_re_user.workflow_step_template_id
LEFT JOIN users AS ass_user ON wst_re_user.user_id = ass_user.id
WHERE
w.id in (SELECT DISTINCT(w.id)

{{- template "body" . -}}
)
GROUP BY w.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var workflowsCountTpl = `SELECT COUNT(DISTINCT w.id)

{{- template "body" . -}}
`

var workflowsQueryBodyTpl = `
{{ define "body" }}
FROM workflows AS w
LEFT JOIN tasks ON w.task_id = tasks.id
LEFT JOIN instances AS inst ON tasks.instance_id = inst.id
LEFT JOIN users AS create_user ON w.create_user_id = create_user.id
LEFT JOIN workflow_records AS wr ON w.workflow_record_id = wr.id
LEFT JOIN workflow_steps AS ws ON wr.current_workflow_step_id = ws.id
LEFT JOIN workflow_step_templates AS wst ON ws.workflow_step_template_id = wst.id
LEFT JOIN workflow_step_template_user AS wst_re_user ON wst.id = wst_re_user.workflow_step_template_id
LEFT JOIN users AS ass_user ON wst_re_user.user_id = ass_user.id
WHERE
w.deleted_at IS NULL

{{- if .filter_create_user_name }}
AND create_user.login_name = :filter_create_user_name
{{- end }}

{{- if .filter_current_step_type }}
AND wst.type = :filter_current_step_type
{{- end }}

{{- if .filter_status }}
AND wr.status = :filter_status
{{- end }}

{{- if .filter_current_step_assignee_user_name }}
AND ass_user.login_name = :filter_current_step_assignee_user_name
{{- end }}

{{- if .filter_task_status }}
AND tasks.status = :filter_task_status
{{- end }}

{{- if .filter_task_instance_name }}
AND inst.name = :filter_task_instance_name
{{- end }}
{{- end }}
`

func (s *Storage) GetWorkflowsByReq(data map[string]interface{}) ([]*WorkflowListDetail, uint64, error) {
	result := []*WorkflowListDetail{}
	count, err := s.getListResult(workflowsQueryBodyTpl, workflowsQueryTpl, workflowsCountTpl, data, &result)
	return result, count, err
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

func (s *Storage) GetWorkflowTemplatesByReq(data map[string]interface{}) ([]*WorkflowTemplateDetail, uint64, error) {
	result := []*WorkflowTemplateDetail{}
	count, err := s.getListResult(workflowTemplatesQueryBodyTpl, workflowTemplatesQueryTpl,
		workflowTemplatesCountTpl, data, &result)
	return result, count, err
}
