package model

import (
	"bytes"
	"text/template"
)

type UserDetail struct {
	Id        int
	Name      string `json:"login_name"`
	Email     string
	RoleNames string `json:"role_names"` // is a role name list, separated by commas.
}

var usersQueryTpl = `SELECT users.id, users.login_name, users.email, GROUP_CONCAT(DISTINCT COALESCE(roles.name,'')) AS role_names
FROM users 
LEFT JOIN user_role ON users.id = user_role.user_id
LEFT JOIN roles ON user_role.role_id = roles.id
WHERE
users.id in (SELECT DISTINCT(users.id)

{{- template "body" . -}}
)
GROUP BY users.id
{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var usersCountTpl = `SELECT COUNT(DISTINCT users.id)

{{- template "body" . -}}
`

var usersQueryBodyTpl = `
{{ define "body" }}
FROM users 
LEFT JOIN user_role ON users.id = user_role.user_id
LEFT JOIN roles ON user_role.role_id = roles.id
WHERE
users.deleted_at IS NULL

{{- if .filter_user_name }}
AND users.login_name = :filter_user_name
{{- end }}

{{- if .filter_role_name }}
AND roles.name = :filter_role_name
{{- end }}
{{- end }}
`

func getUsersQuery(data interface{}) (string, error) {
	var buff bytes.Buffer
	tpl, err := template.New("getUsersByReq").Parse(usersQueryBodyTpl)
	if err != nil {
		return "", err
	}
	_, err = tpl.Parse(usersQueryTpl)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func getUsersCountQuery(data interface{}) (string, error) {
	var buff bytes.Buffer
	tpl, err := template.New("getUsersByReq").Parse(usersQueryBodyTpl)
	if err != nil {
		return "", err
	}
	_, err = tpl.Parse(usersCountTpl)
	if err != nil {
		return "", err
	}
	err = tpl.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (s *Storage) GetUsersByReq(data map[string]interface{}) ([]*UserDetail, uint64, error) {
	tasksQuery, err := getUsersQuery(data)
	if err != nil {
		return nil, 0, err
	}
	tasksCountQuery, err := getUsersCountQuery(data)
	if err != nil {
		return nil, 0, err
	}

	sqlxDb := GetSqlxDb()
	nstmtTasksQuery, err := sqlxDb.PrepareNamed(tasksQuery)
	if err != nil {
		return nil, 0, err
	}
	users := []*UserDetail{}

	err = nstmtTasksQuery.Select(&users, data)
	if err != nil {
		return nil, 0, err
	}

	nstmtTasksCountQuery, err := sqlxDb.PrepareNamed(tasksCountQuery)
	if err != nil {
		return nil, 0, err
	}
	var count uint64
	err = nstmtTasksCountQuery.Get(&count, data)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}
