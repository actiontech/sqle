package model

import (
	"database/sql"
	"time"

	"github.com/actiontech/sqle/sqle/utils"
)

type WorkflowListDetail struct {
	ProjectName             string         `json:"project_name"`
	Id                      uint           `json:"workflow_id"`
	Subject                 string         `json:"subject"`
	Desc                    string         `json:"desc"`
	CreateUser              sql.NullString `json:"create_user_name"`
	CreateUserDeletedAt     *time.Time     `json:"create_user_deleted_at"`
	CreateTime              *time.Time     `json:"create_time"`
	CurrentStepType         sql.NullString `json:"current_step_type" enums:"sql_review,sql_execute"`
	CurrentStepAssigneeUser RowList        `json:"current_step_assignee_user_name_list"`
	TaskStatus              RowList        `json:"task_status"`
	Status                  string         `json:"status"`
	TaskInstanceType        RowList        `json:"task_instance_type"`
}

var workflowsQueryTpl = `
SELECT p.name 														 AS project_name,
       w.id                                                          AS workflow_id,
       w.subject,
       w.desc,
       create_user.login_name                                        AS create_user_name,
	   create_user.deleted_at                                        AS create_user_deleted_at,
       w.created_at                                                  AS create_time,
       curr_wst.type                                                 AS current_step_type,
       GROUP_CONCAT(DISTINCT COALESCE(curr_ass_user.login_name, '')) AS current_step_assignee_user_name_list,
       GROUP_CONCAT(tasks.status)                                    AS task_status,
       wr.status,																							
	   GROUP_CONCAT(inst.db_type)                                    AS task_instance_type
{{- template "body" . -}}
GROUP BY p.id,w.id
ORDER BY w.id DESC
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var workflowsCountTpl = `SELECT COUNT(DISTINCT w.id)

{{- template "body" . -}}
`

var workflowsQueryBodyTpl = `
{{ define "body" }}
FROM projects p
LEFT JOIN workflows AS w ON w.project_id = p.id
LEFT JOIN users AS create_user ON w.create_user_id = create_user.id
LEFT JOIN workflow_records AS wr ON w.workflow_record_id = wr.id
LEFT JOIN workflow_instance_records wir on wir.workflow_record_id = wr.id
LEFT JOIN tasks ON wir.task_id = tasks.id
LEFT JOIN instances AS inst ON tasks.instance_id = inst.id
LEFT JOIN workflow_steps AS curr_ws ON wr.current_workflow_step_id = curr_ws.id
LEFT JOIN workflow_step_templates AS curr_wst ON curr_ws.workflow_step_template_id = curr_wst.id
LEFT JOIN workflow_step_user AS curr_wst_re_user ON curr_ws.id = curr_wst_re_user.workflow_step_id
LEFT JOIN users AS curr_ass_user ON curr_wst_re_user.user_id = curr_ass_user.id

{{- if .check_user_can_access }}
LEFT JOIN workflow_steps AS all_ws ON w.id = all_ws.workflow_id AND all_ws.state !='initialized'
LEFT JOIN workflow_step_templates AS all_wst ON all_ws.workflow_step_template_id = all_wst.id
LEFT JOIN workflow_step_user AS all_wst_re_user ON all_ws.id = all_wst_re_user.workflow_step_id
LEFT JOIN users AS all_ass_user ON all_wst_re_user.user_id = all_ass_user.id
{{- end }}
WHERE
w.deleted_at IS NULL 

{{- if .check_user_can_access }}
AND (
w.create_user_id = :current_user_id 
OR curr_ass_user.id = :current_user_id
OR all_ass_user.id = :current_user_id

{{- if .viewable_instance_ids }} 
OR inst.id IN ( {{ .viewable_instance_ids }})
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

{{- if .filter_create_user_name }}
AND create_user.login_name = :filter_create_user_name
{{- end }}

{{- if .filter_current_step_type }}
AND curr_wst.type = :filter_current_step_type
{{- end }}

{{- if .filter_status }}
AND wr.status = :filter_status
{{- end }}

{{- if .filter_current_step_assignee_user_name }}
AND curr_ass_user.login_name = :filter_current_step_assignee_user_name
{{- end }}

{{- if .filter_task_status }}
AND tasks.status = :filter_task_status
{{- end }}

{{- if .filter_task_instance_name }}
AND inst.name = :filter_task_instance_name
{{- end }}

{{- if .filter_project_name }}
AND p.name = :filter_project_name
{{- end }}

{{ end }}

`

func (s *Storage) GetWorkflowsByReq(data map[string]interface{}, user *User) (
	result []*WorkflowListDetail, count uint64, err error) {

	var ids []uint
	{
		instances, err := s.GetUserCanOpInstances(user, []uint{OP_WORKFLOW_VIEW_OTHERS})
		if err != nil {
			return result, 0, err
		}
		ids = getInstanceIDsFromInstances(instances)
	}

	if len(ids) > 0 {
		data["viewable_instance_ids"] = utils.JoinUintSliceToString(ids, ", ")
	}

	err = s.getListResult(workflowsQueryBodyTpl, workflowsQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}

	count, err = s.getCountResult(workflowsQueryBodyTpl, workflowsCountTpl, data)

	return result, count, err
}

func (s *Storage) GetWorkflowCountByReq(data map[string]interface{}) (uint64, error) {
	return s.getCountResult(workflowsQueryBodyTpl, workflowsCountTpl, data)
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
