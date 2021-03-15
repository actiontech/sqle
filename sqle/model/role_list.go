package model

import (
	"bytes"
	"text/template"
)

type RoleDetail struct {
	Id            int
	Name          string `json:"name"`
	Desc          string
	UserNames     string `json:"user_names"`     // is a user name list, separated by commas.
	InstanceNames string `json:"instance_names"` // is a instance name list, separated by commas.
}

var rolesQueryTpl = `SELECT roles.id, roles.name, roles.desc,
GROUP_CONCAT(DISTINCT COALESCE(users.login_name,'')) AS user_names,
GROUP_CONCAT(DISTINCT COALESCE(instances.name,'')) AS instance_names
FROM roles
LEFT JOIN user_role ON roles.id = user_role.role_id
LEFT JOIN users ON user_role.user_id = users.id
LEFT JOIN instance_role ON roles.id = instance_role.role_id
LEFT JOIN instances ON instance_role.instance_id = instances.id
WHERE
roles.id in (SELECT DISTINCT(roles.id)

{{- template "body" . -}}
)
GROUP BY roles.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var rolesCountTpl = `SELECT COUNT(DISTINCT roles.id)

{{- template "body" . -}}
`

var rolesQueryBodyTpl = `
{{ define "body" }}
FROM roles 
LEFT JOIN user_role ON roles.id = user_role.role_id
LEFT JOIN users ON user_role.user_id = users.id
LEFT JOIN instance_role ON roles.id = instance_role.role_id
LEFT JOIN instances ON instance_role.instance_id = instances.id
WHERE
roles.deleted_at IS NULL

{{- if .filter_role_name }}
AND roles.name = :filter_role_name
{{- end }}

{{- if .filter_user_name }}
AND users.login_name = :filter_user_name
{{- end }}

{{- if .filter_instance_name }}
AND instances.name = :filter_instance_name
{{- end }}
{{- end }}
`

func (s *Storage) GetRolesByReq(data map[string]interface{}) ([]*RoleDetail, uint64, error) {
	result := []*RoleDetail{}
	count, err := s.getListResult(rolesQueryBodyTpl, rolesQueryTpl, rolesCountTpl, data, &result)
	return result, count, err
}

func getSelectQuery(bodyTpl, queryTpl string, data interface{}) (string, error) {
	var buff bytes.Buffer
	tpl, err := template.New("getQuery").Parse(bodyTpl)
	if err != nil {
		return "", err
	}
	_, err = tpl.Parse(queryTpl)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func getCountQuery(bodyTpl, countTpl string, data interface{}) (string, error) {
	var buff bytes.Buffer
	tpl, err := template.New("getCount").Parse(bodyTpl)
	if err != nil {
		return "", err
	}
	_, err = tpl.Parse(countTpl)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (s *Storage) getListResult(bodyTpl, queryTpl, countTpl string, data map[string]interface{},
	result interface{}) (uint64, error) {
	selectQuery, err := getSelectQuery(bodyTpl, queryTpl, data)
	if err != nil {
		return 0, err
	}
	countQuery, err := getCountQuery(bodyTpl, countTpl, data)
	if err != nil {
		return 0, err
	}

	sqlxDb := GetSqlxDb()
	nstmtTasksQuery, err := sqlxDb.PrepareNamed(selectQuery)
	if err != nil {
		return 0, err
	}

	err = nstmtTasksQuery.Select(result, data)
	if err != nil {
		return 0, err
	}

	nstmtCountQuery, err := sqlxDb.PrepareNamed(countQuery)
	if err != nil {
		return 0, err
	}
	var count uint64
	err = nstmtCountQuery.Get(&count, data)
	if err != nil {
		return 0, err
	}
	return count, nil
}
