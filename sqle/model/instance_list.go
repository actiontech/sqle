package model

type InstanceDetail struct {
	Name                 string `json:"name"`
	Desc                 string `json:"desc"`
	Host                 string `json:"db_host"`
	Port                 string `json:"db_port"`
	User                 string `json:"db_user"`
	WorkflowTemplateName string `json:"workflow_template_name"`
	RoleNames            string `json:"role_names"`          // is a role name list, separated by commas.
	RuleTemplateNames    string `json:"rule_template_names"` // is a rule template name list, separated by commas.
}

var instancesQueryTpl = `SELECT instances.name, instances.desc, 
instances.db_host, instances.db_port, instances.db_user,
GROUP_CONCAT(DISTINCT COALESCE(roles.name,'')) AS role_names,
GROUP_CONCAT(DISTINCT COALESCE(rule_templates.name,'')) AS rule_template_names
FROM instances
LEFT JOIN instance_role ON instances.id = instance_role.instance_id
LEFT JOIN roles ON instance_role.role_id = roles.id
LEFT JOIN instance_rule_template ON instances.id = instance_rule_template.instance_id
LEFT JOIN rule_templates ON instance_rule_template.rule_template_id = rule_templates.id
WHERE
instances.id in (SELECT DISTINCT(instances.id)

{{- template "body" . -}}
)
GROUP BY instances.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var instancesCountTpl = `SELECT COUNT(DISTINCT instances.id)

{{- template "body" . -}}
`

var instancesQueryBodyTpl = `
{{ define "body" }}
FROM instances
LEFT JOIN instance_role ON instances.id = instance_role.instance_id
LEFT JOIN roles ON instance_role.role_id = roles.id
LEFT JOIN instance_rule_template ON instances.id = instance_rule_template.instance_id
LEFT JOIN rule_templates ON instance_rule_template.rule_template_id = rule_templates.id
WHERE
instances.deleted_at IS NULL

{{- if .filter_instance_name }}
AND instances.name = :filter_instance_name
{{- end }}

{{- if .filter_db_host }}
AND instances.db_host = :filter_db_host
{{- end }}

{{- if .filter_db_port }}
AND instances.db_port = :filter_db_port
{{- end }}

{{- if .filter_db_user }}
AND instances.db_user = :filter_db_user
{{- end }}

{{- if .filter_role_name }}
AND roles.name = :filter_role_name
{{- end }}

{{- if .filter_rule_template_name }}
AND rule_templates.name = :filter_rule_template_name
{{- end }}
{{- end }}
`

func (s *Storage) GetInstancesByReq(data map[string]interface{}) ([]*InstanceDetail, uint64, error) {
	result := []*InstanceDetail{}
	count, err := s.getListResult(instancesQueryBodyTpl, instancesQueryTpl, instancesCountTpl, data, &result)
	return result, count, err
}
