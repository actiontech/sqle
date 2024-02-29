// sqlop 包定义了SQL权限操作的相关结构体，目前用于事前权限校验
// 背景见：https://github.com/actiontech/dms-ee/issues/125

package sqlop

import "encoding/json"

func EncodingSQLObjectOps(s *SQLObjectOps) (string, error) {
	jsonStr, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(jsonStr), nil
}

func DecodingSQLObjectOps(s string) (*SQLObjectOps, error) {
	var sqlObjectOps SQLObjectOps
	err := json.Unmarshal([]byte(s), &sqlObjectOps)
	if err != nil {
		return nil, err
	}
	return &sqlObjectOps, nil
}

type SQLObjectOps struct {
	// ObjectOps 表示sql中涉及的对象及对对象的操作
	ObjectOps []*SQLObjectOp
	Sql       SQLInfo
}

func NewSQLObjectOps(sql string) *SQLObjectOps {
	return &SQLObjectOps{
		ObjectOps: []*SQLObjectOp{},
		Sql:       SQLInfo{Sql: sql},
	}
}

func (s *SQLObjectOps) AddObjectOp(o ...*SQLObjectOp) {
	s.ObjectOps = append(s.ObjectOps, o...)
}

type SQLObjectOp struct {
	Op     SQLOp      // 对象操作
	Object *SQLObject // 对象
}

type SQLObject struct {
	// Type 表示对象的类型
	Type SQLObjectType
	// DatabaseName 表示对象所在的database，若对象不属于database，或无法从Sql中解析出当前database，则为空字符串
	DatabaseName string
	// SchemaName 表示对象所在的schema，若对象不属于schema，或无法从Sql中解析出当前schema，则为空字符串
	// 对于一些数据库类型，如PostgreSQL，可能存在schema的概念，此时SchemaName字段应该被使用
	// 对于一些数据库类型，如MySQL，可能不存在schema的概念，或schema的概念与database的概念相同，此时SchemaName字段应该为空字符串
	SchemaName string
	// TableName 表示对象的表名，如果对象不是表，则为空字符串
	TableName string
}

type SQLObjectType string

const (
	SQLObjectTypeTable    SQLObjectType = "Table"
	SQLObjectTypeSchema   SQLObjectType = "Schema"
	SQLObjectTypeDatabase SQLObjectType = "Database"
	SQLObjectTypeInstance SQLObjectType = "Instance"
	SQLObjectTypeServer   SQLObjectType = "Server"
)

type SQLInfo struct {
	Sql string
}
type SQLOp string

const (
	// 增或改操作
	SQLOpAddOrUpdate SQLOp = "AddOrUpdate"
	// 读取操作
	SQLOpRead SQLOp = "Read"
	// 删除操作
	SQLOpDelete SQLOp = "Delete"
	// 授权操作
	SQLOpGrant SQLOp = "Grant"
	// 高权限操作，如锁表、导出表到文件等
	SQLOpAdmin SQLOp = "Admin"
)
