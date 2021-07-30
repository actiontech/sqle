package model

type UserDetail struct {
	Id        int
	Name      string `json:"login_name"`
	Email     string
	RoleNames RowList `json:"role_names"`
}

var usersQueryTpl = `SELECT users.id, users.login_name, users.email, GROUP_CONCAT(DISTINCT COALESCE(roles.name,'')) AS role_names
FROM users 
LEFT JOIN user_role ON users.id = user_role.user_id
LEFT JOIN roles ON user_role.role_id = roles.id AND roles.deleted_at IS NULL
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
LEFT JOIN roles ON user_role.role_id = roles.id AND roles.deleted_at IS NULL
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

func (s *Storage) GetUsersByReq(data map[string]interface{}) (
	result []*UserDetail, count uint64, err error) {

	err = s.getListResult(usersQueryBodyTpl, usersQueryTpl, data, &result)
	if err != nil {
		return result, 0, err
	}
	count, err = s.getCountResult(usersQueryBodyTpl, usersCountTpl, data)
	return result, count, err
}
