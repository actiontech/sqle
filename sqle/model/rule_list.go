package model

type RuleTemplateDetail struct {
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	DBType string `json:"db_type"`
	// InstanceIds   RowList `json:"instance_ids"`
	// InstanceNames RowList `json:"instance_names"`
}

var ruleTemplatesQueryTpl = `SELECT rule_templates.name, rule_templates.desc, rule_templates.db_type
{{- template "body" . }}

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var ruleTemplatesCountTpl = `SELECT COUNT(DISTINCT id)

{{- template "body" . -}}
`

var ruleTemplatesQueryBodyTpl = `
{{ define "body" }}
FROM rule_templates
WHERE
deleted_at IS NULL
AND project_id = :project_id
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
