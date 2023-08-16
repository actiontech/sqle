package model

import (
	"bytes"
	"text/template"
)

// type RoleDetail struct {
// 	Id            int
// 	Name          string  `json:"name"`
// 	Desc          string  `json:"desc"`
// 	UserNames     RowList `json:"user_names"`
// 	InstanceNames RowList `json:"instance_names"`

// 	// New fields: Stat, UserGroupNames, OperationsCodes
// 	// Issue: https://github.com/actiontech/sqle/issues/228
// 	// Version: >= sqle-v1.2202.0
// 	Stat            int     `json:"stat"`
// 	UserGroupNames  RowList `json:"user_group_names"`
// 	OperationsCodes RowList `json:"operations_codes"`
// }

// func (rd *RoleDetail) IsDisabled() bool {
// 	return rd.Stat == Disabled
// }

// var rolesQueryTpl = `SELECT roles.id, roles.name, roles.desc, roles.stat,
// GROUP_CONCAT(DISTINCT COALESCE(role_operations.op_code,'')) AS operations_codes
// FROM roles
// LEFT JOIN role_operations ON role_operations.role_id = roles.id AND role_operations.deleted_at IS NULL
// WHERE
// roles.id in (SELECT DISTINCT(roles.id)

// {{- template "body" . -}}
// )
// GROUP BY roles.id
// {{- if .limit }}
// LIMIT :limit OFFSET :offset
// {{- end -}}
// `

// var rolesCountTpl = `SELECT COUNT(DISTINCT roles.id)

// {{- template "body" . -}}
// `

// var rolesQueryBodyTpl = `
// {{ define "body" }}
// FROM roles
// LEFT JOIN role_operations ON role_operations.role_id = roles.id AND role_operations.deleted_at IS NULL
// WHERE
// roles.deleted_at IS NULL

// {{- if .filter_role_name }}
// AND roles.name = :filter_role_name
// {{- end }}

// {{- end }}
// `

// func (s *Storage) GetRolesByReq(data map[string]interface{}) (
// 	result []*RoleDetail, count uint64, err error) {

// 	err = s.getListResult(rolesQueryBodyTpl, rolesQueryTpl, data, &result)
// 	if err != nil {
// 		return result, 0, err
// 	}
// 	count, err = s.getCountResult(rolesQueryBodyTpl, rolesCountTpl, data)
// 	return result, count, err
// }

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

func (s *Storage) getListResult(bodyTpl, queryTpl string, data map[string]interface{},
	result interface{}) error {
	selectQuery, err := getSelectQuery(bodyTpl, queryTpl, data)
	if err != nil {
		return err
	}

	sqlxDb := GetSqlxDb()
	nameStmtTasksQuery, err := sqlxDb.PrepareNamed(selectQuery)
	if err != nil {
		return err
	}
	err = nameStmtTasksQuery.Select(result, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) getCountResult(bodyTpl, countTpl string, data map[string]interface{}) (uint64, error) {
	sqlxDb := GetSqlxDb()
	countQuery, err := getCountQuery(bodyTpl, countTpl, data)
	if err != nil {
		return 0, err
	}
	nameStmtCountQuery, err := sqlxDb.PrepareNamed(countQuery)
	if err != nil {
		return 0, err
	}
	var count uint64
	err = nameStmtCountQuery.Get(&count, data)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// func GetRoleIDsFromRoles(roles []*Role) (roleIDs []uint) {

// 	roleIDs = make([]uint, len(roles))

// 	for i := range roles {
// 		roleIDs[i] = roles[i].ID
// 	}

// 	return roleIDs
// }
