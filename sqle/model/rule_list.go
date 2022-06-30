package model

type RuleTemplateDetail struct {
	Name          string  `json:"name"`
	Desc          string  `json:"desc"`
	DBType        string  `json:"db_type"`
	RuleNames     RowList `json:"rule_names"`
	InstanceNames RowList `json:"instance_names"`
}

var ruleTemplatesQueryTpl = `SELECT rt1.name, rt1.desc, rt1.db_type,
GROUP_CONCAT(DISTINCT COALESCE(rules.name,'')) AS rule_names,
GROUP_CONCAT(DISTINCT COALESCE(instances.name,'')) AS instance_names
FROM rule_templates AS rt1
LEFT JOIN rule_template_rule ON rt1.id = rule_template_rule.rule_template_id
LEFT JOIN rules ON rule_template_rule.rule_name = rules.name
LEFT JOIN instance_rule_template ON rt1.id = instance_rule_template.rule_template_id
LEFT JOIN instances ON instance_rule_template.instance_id = instances.id AND instances.deleted_at IS NULL
WHERE EXISTS (SELECT DISTINCT(rule_templates.id)

{{- template "body" . }}
AND rt1.id = rule_templates.id
)
GROUP BY rt1.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var ruleTemplatesCountTpl = `SELECT COUNT(DISTINCT rule_templates.id)

{{- template "body" . -}}
`

var ruleTemplatesQueryBodyTpl = `
{{ define "body" }}
FROM rule_templates
LEFT JOIN rule_template_rule ON rule_templates.id = rule_template_rule.rule_template_id
LEFT JOIN rules ON rule_template_rule.rule_name = rules.name
LEFT JOIN instance_rule_template ON rule_templates.id = instance_rule_template.rule_template_id
LEFT JOIN instances ON instance_rule_template.instance_id = instances.id AND instances.deleted_at IS NULL
WHERE
rule_templates.deleted_at IS NULL

{{- if .filter_instance_name }}
AND instances.name = :filter_instance_name
{{- end }}
{{- end }}
`

func (s *Storage) GetRuleTemplatesByReq(data map[string]interface{}) (
	result []*RuleTemplateDetail, count uint64, err error) {

	err = s.getListResult(ruleTemplatesQueryBodyTpl, ruleTemplatesQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(ruleTemplatesQueryBodyTpl, ruleTemplatesCountTpl, data)
	return result, count, err
}
