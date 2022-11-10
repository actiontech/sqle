package model

import (
	"database/sql"
	"fmt"

	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/utils"
)

type InstanceDetail struct {
	Name                 string         `json:"name"`
	DbType               string         `json:"db_type"`
	Desc                 string         `json:"desc"`
	Host                 string         `json:"db_host"`
	Port                 string         `json:"db_port"`
	User                 string         `json:"db_user"`
	MaintenancePeriod    Periods        `json:"maintenance_period" gorm:"text"`
	WorkflowTemplateName sql.NullString `json:"workflow_template_name"`
	RoleNames            RowList        `json:"role_names"`
	RuleTemplateNames    RowList        `json:"rule_template_names"`
	SqlQueryConfig       SqlQueryConfig `json:"sql_query_config"`
}

var instancesQueryTpl = `SELECT inst.name, inst.db_type, inst.desc, inst.db_host,
inst.db_port, inst.db_user, inst.maintenance_period, inst.sql_query_config, wt.name AS workflow_template_name,
GROUP_CONCAT(DISTINCT COALESCE(roles.name,'')) AS role_names,
GROUP_CONCAT(DISTINCT COALESCE(rt.name,'')) AS rule_template_names
FROM instances AS inst
LEFT JOIN instance_role AS ir ON inst.id = ir.instance_id
LEFT JOIN roles ON ir.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
LEFT JOIN instance_rule_template AS inst_rt ON inst.id = inst_rt.instance_id
LEFT JOIN rule_templates AS rt ON inst_rt.rule_template_id = rt.id AND rt.deleted_at IS NULL
LEFT JOIN workflow_templates AS wt ON inst.workflow_template_id = wt.id AND wt.deleted_at IS NULL 
WHERE
inst.id in (SELECT DISTINCT(inst.id)

{{- template "body" . -}}
)
GROUP BY inst.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var instancesCountTpl = `SELECT COUNT(DISTINCT inst.id)

{{- template "body" . -}}
`

var instancesQueryBodyTpl = `
{{ define "body" }}
FROM instances AS inst
LEFT JOIN instance_role AS ir ON inst.id = ir.instance_id
LEFT JOIN roles ON ir.role_id = roles.id AND roles.deleted_at IS NULL AND roles.stat = 0
LEFT JOIN instance_rule_template AS inst_rt ON inst.id = inst_rt.instance_id
LEFT JOIN rule_templates AS rt ON inst_rt.rule_template_id = rt.id AND rt.deleted_at IS NULL
LEFT JOIN workflow_templates AS wt ON inst.workflow_template_id = wt.id AND wt.deleted_at IS NULL
LEFT JOIN projects AS p ON inst.project_id = projects.id

WHERE inst.deleted_at IS NULL

AND p.name = :filter_project_name

{{- if .filter_instance_name }}
AND inst.name = :filter_instance_name
{{- end }}

{{- if .filter_db_host }}
AND inst.db_host = :filter_db_host
{{- end }}

{{- if .filter_db_port }}
AND inst.db_port = :filter_db_port
{{- end }}

{{- if .filter_db_user }}
AND inst.db_user = :filter_db_user
{{- end }}

{{- if .filter_db_type }}
AND inst.db_type = :filter_db_type
{{- end }}

{{- if .filter_rule_template_name }}
AND rt.name = :filter_rule_template_name
{{- end }}

{{- if .check_user_can_access }}

AND roles.id IN  ( {{ .role_id_list }} )

{{- end }}

{{- end }}
`

func (s *Storage) GetInstancesByReq(data map[string]interface{}, user *User) (
	result []*InstanceDetail, count uint64, err error) {

	if data["filter_project_name"] == "" {
		return nil, 0, errors.New(errors.DataInvalid, fmt.Errorf("project name can not be empty"))
	}

	if !IsDefaultAdminUser(user.Name) {
		roles, err := s.GetRolesByUserID(int(user.ID))
		if err != nil {
			return result, count, err
		}
		if len(roles) == 0 {
			return result, count, nil
		}
		roleIDs := GetRoleIDsFromRoles(roles)
		data["role_id_list"] = utils.JoinUintSliceToString(roleIDs, ", ")
	}
	err = s.getListResult(instancesQueryBodyTpl, instancesQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}

	count, err = s.getCountResult(instancesQueryBodyTpl, instancesCountTpl, data)
	return result, count, err

}
