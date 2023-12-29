package model

// import "database/sql"

// type UserDetail struct {
// 	Id             int
// 	Name           string `json:"login_name"`
// 	Email          string
// 	Phone          sql.NullString `json:"phone"`
// 	WeChatID       sql.NullString `json:"wechat_id"`
// 	LoginType      string         `json:"user_authentication_type"`
// 	Stat           int            `json:"stat"`
// 	RoleNames      RowList        `json:"role_names"`
// 	UserGroupNames RowList        `json:"user_group_names"`
// }

// func (u *UserDetail) IsDisabled() bool {
// 	return u.Stat == Disabled
// }

// var usersQueryTpl = `SELECT
// users.id, users.login_name, users.email, users.phone, users.wechat_id,
// users.user_authentication_type, users.stat,
// GROUP_CONCAT(DISTINCT COALESCE(user_groups.name,'')) AS user_group_names
// FROM users
// LEFT JOIN user_group_users ON users.id = user_group_users.user_id
// LEFT JOIN user_groups ON user_group_users.user_group_id = user_groups.id AND user_groups.deleted_at IS NULL
// WHERE
// users.id in (SELECT DISTINCT(users.id)
// {{- template "body" . -}}
// )
// GROUP BY users.id
// {{- if .limit }}
// LIMIT :limit OFFSET :offset
// {{- end -}}
// `

// var usersCountTpl = `SELECT COUNT(DISTINCT users.id)
// {{- template "body" . -}}
// `

// var usersQueryBodyTpl = `
// {{ define "body" }}
// FROM users
// WHERE
// users.deleted_at IS NULL
// {{- if .filter_user_name }}
// AND users.login_name = :filter_user_name
// {{- end }}
// {{- end }}
// `

// func (s *Storage) GetUsersByReq(data map[string]interface{}) (
// 	result []*UserDetail, count uint64, err error) {

// 	err = s.getListResult(usersQueryBodyTpl, usersQueryTpl, data, &result)
// 	if err != nil {
// 		return result, 0, err
// 	}
// 	count, err = s.getCountResult(usersQueryBodyTpl, usersCountTpl, data)
// 	return result, count, err
// }
