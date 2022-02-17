package model

import (
	"database/sql"
	"fmt"

	"github.com/actiontech/sqle/sqle/utils"
)

type InstanceDetail struct {
	Name                 string         `json:"name"`
	Desc                 string         `json:"desc"`
	Host                 string         `json:"db_host"`
	Port                 string         `json:"db_port"`
	User                 string         `json:"db_user"`
	WorkflowTemplateName sql.NullString `json:"workflow_template_name"`
	RoleNames            RowList        `json:"role_names"`
	RuleTemplateNames    RowList        `json:"rule_template_names"`
}

var instancesQueryTpl = `SELECT inst.name, inst.desc, inst.db_host,
inst.db_port, inst.db_user, wt.name AS workflow_template_name,
GROUP_CONCAT(DISTINCT COALESCE(roles.name,'')) AS role_names,
GROUP_CONCAT(DISTINCT COALESCE(rt.name,'')) AS rule_template_names
FROM instances AS inst
LEFT JOIN instance_role AS ir ON inst.id = ir.instance_id
LEFT JOIN roles ON ir.role_id = roles.id AND roles.deleted_at IS NULL
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
LEFT JOIN roles ON ir.role_id = roles.id AND roles.deleted_at IS NULL
LEFT JOIN instance_rule_template AS inst_rt ON inst.id = inst_rt.instance_id
LEFT JOIN rule_templates AS rt ON inst_rt.rule_template_id = rt.id AND rt.deleted_at IS NULL
LEFT JOIN workflow_templates AS wt ON inst.workflow_template_id = wt.id AND wt.deleted_at IS NULL

WHERE inst.deleted_at IS NULL

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

{{- if .filter_role_name }}
AND roles.name = :filter_role_name
{{- end }}

{{- if .filter_rule_template_name }}
AND rt.name = :filter_rule_template_name
{{- end }}

{{- if .filter_workflow_template_name }}
AND wt.name = :filter_workflow_template_name
{{- end }}

{{- if .check_user_can_access }}

AND roles.id IN  ( {{ .role_id_list }} )

{{- end }}

{{- end }}
`

func (s *Storage) GetInstancesByReq(data map[string]interface{}, user *User) (
	result []*InstanceDetail, count uint64, err error) {

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

func (s *Storage) GetUserAvailableInstancesByUserID(userID int) (
	availabeInsts []*Instance, err error) {

	roles, err := s.GetRolesByUserID(userID)
	if err != nil {
		return availabeInsts, err
	}

	if len(roles) == 0 { // user has no roles
		return availabeInsts, nil
	}

	roleIDs := GetRoleIDsFromRoles(roles)

	return s.GetInstancesByRoleIDs(roleIDs)

}

var instancesQueryByRoleIDs = `
SELECT inst.id, inst.name, inst.db_type, inst.db_host, inst.db_port, inst.db_user, inst.desc 
FROM instances AS inst
LEFT JOIN instance_role AS ir ON inst.id = ir.instance_id
LEFT JOIN roles ON ir.role_id = roles.id AND roles.deleted_at IS NULL
WHERE roles.id in (%s)
`

func (s *Storage) GetInstancesByRoleIDs(roleIds []uint) (insts []*Instance, err error) {

	if len(roleIds) == 0 {
		return insts, nil
	}

	query := fmt.Sprintf(instancesQueryByRoleIDs, utils.JoinUintSliceToString(roleIds, ", "))
	if err := s.db.Unscoped().Raw(query).Find(&insts).Error; err != nil {
		return insts, err
	}

	return insts, err
}

func GetInstanceIDsFromInstances(insts []*Instance) (ids []uint) {
	ids = make([]uint, len(insts))
	for i := range insts {
		ids[i] = insts[i].ID
	}
	return ids
}
