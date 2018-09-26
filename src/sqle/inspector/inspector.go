package inspector

import (
	"errors"
	"sqle/storage"
)

func Inspect(task *storage.Task) ([]*storage.Sql, error) {
	switch task.Db.DbType {
	case storage.DB_TYPE_MYSQL:
		sqls, err := inspectMysql(task.ReqSql)
		if err != nil {
			return nil, err
		}
		return sqls, nil
	default:
		return nil, errors.New("db type is invalid")
	}
}

// InspectMysql support multi-sql, split by ";".
func inspectMysql(sql string) ([]*storage.Sql, error) {
	sqlMeta := []*storage.Sql{}
	return sqlMeta, nil
}

//
func GenerateRollbackSql(sql string) {

}
