package fillsql

import (
	"fmt"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/driver/mysql"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	v2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	pingcapMysql "github.com/pingcap/parser/mysql"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/sirupsen/logrus"
)

const (
	MybatisXMLCharDefaultValue  = "1"
	MybatisXMLIntDefaultValue   = 1
	MybatisXMLFloatDefaultValue = 1.0
	XMLFileExtension            = ".XML"
)

func FillingSQLWithParamMarker(sqlContent string, task *model.Task) (string, error) {
	l := log.NewEntry()
	driver := mysql.MysqlDriverImpl{}
	nodes, err := driver.ParseSql(sqlContent)
	if err != nil {
		return sqlContent, err
	}

	if task.Instance == nil {
		return sqlContent, nil
	}

	// sql分析是单条sql分析
	if len(nodes) == 0 {
		return sqlContent, nil
	}
	node := nodes[0]
	var tableRefs *ast.Join
	var where ast.ExprNode
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		tableRefs = stmt.From.TableRefs
		where = stmt.Where
	case *ast.UpdateStmt:
		tableRefs = stmt.TableRefs.TableRefs
		where = stmt.Where
	case *ast.DeleteStmt:
		tableRefs = stmt.TableRefs.TableRefs
		where = stmt.Where
	default:
		return sqlContent, nil
	}

	if where == nil {
		return sqlContent, nil
	}
	schema := task.Schema
	if task.Schema == "" {
		schema = getSchemaFromTableRefs(tableRefs)
	}

	dsn, err := newDSN(task.Instance, schema)
	if err != nil {
		return sqlContent, err
	}
	conn, err := executor.NewExecutor(log.NewEntry(), dsn, schema)
	if err != nil {
		return sqlContent, err
	}
	defer conn.Db.Close()
	ctx := session.NewContext(nil, session.WithExecutor(conn))
	ctx.SetCurrentSchema(schema)

	tableNameCreateTableStmtMap := ctx.GetTableNameCreateTableStmtMap(tableRefs)
	fillParamMarker(l, where, tableNameCreateTableStmtMap)
	return restore(node)
}

func newDSN(instance *model.Instance, database string) (*v2.DSN, error) {
	if instance == nil {
		return nil, fmt.Errorf("instance is nil")
	}

	return &v2.DSN{
		Host:             instance.Host,
		Port:             instance.Port,
		User:             instance.User,
		Password:         instance.Password,
		AdditionalParams: instance.AdditionalParams,
		DatabaseName:     database,
	}, nil
}

func fillParamMarker(l *logrus.Entry, where ast.ExprNode, tableCreateStmtMap map[string]*ast.CreateTableStmt) {
	util.ScanWhereStmt(func(expr ast.ExprNode) bool {
		switch stmt := expr.(type) {
		case *ast.BinaryOperationExpr:
			// where name=?, 解析器会将'?'解析为ParamMarkerExpr
			if column, ok := stmt.L.(*ast.ColumnNameExpr); ok {
				if _, ok := stmt.R.(*parserdriver.ParamMarkerExpr); !ok {
					return true
				}
				defaultValue, err := fillColumnDefaultValue(column, tableCreateStmtMap)
				if err != nil {
					l.Error(err)
				}
				if defaultValue == nil {
					return false
				}
				stmt.R = defaultValue
			} else if column, ok := stmt.R.(*ast.ColumnNameExpr); ok {
				// 存在列名在比较符号左侧的情况 `where ?=name`
				if _, ok := stmt.L.(*parserdriver.ParamMarkerExpr); !ok {
					return true
				}
				defaultValue, err := fillColumnDefaultValue(column, tableCreateStmtMap)
				if err != nil {
					l.Error(err)
				}
				if defaultValue == nil {
					return false
				}
				stmt.L = defaultValue
			}
		}
		return false
	}, where)
}

func fillColumnDefaultValue(column *ast.ColumnNameExpr, tableCreateStmtMap map[string]*ast.CreateTableStmt) (ast.ExprNode, error) {
	tableName := column.Name.Table.L
	columnName := column.Name.Name.L
	// table name为空，代表没有进行连表查询也没有使用别名 e.g: select * from users where id=?
	if tableName == "" {
		for k := range tableCreateStmtMap {
			tableName = k
		}
	}
	createTableStmt, ok := tableCreateStmtMap[tableName]
	if !ok {
		return nil, fmt.Errorf("fillXmlSql get create table sql failed, table:%v", tableName)
	}
	currentTime := time.Now()
	for _, col := range createTableStmt.Cols {
		if col.Name.Name.L != columnName {
			continue
		}
		var value interface{}
		switch col.Tp.Tp {
		case pingcapMysql.TypeVarchar, pingcapMysql.TypeString:
			value = MybatisXMLCharDefaultValue
		// int类型
		case pingcapMysql.TypeLong, pingcapMysql.TypeTiny, pingcapMysql.TypeShort,
			pingcapMysql.TypeInt24, pingcapMysql.TypeLonglong:
			value = MybatisXMLIntDefaultValue
		// 浮点类型
		case pingcapMysql.TypeNewDecimal, pingcapMysql.TypeFloat, pingcapMysql.TypeDouble:
			value = MybatisXMLFloatDefaultValue
		case pingcapMysql.TypeDatetime:
			value = currentTime.Format("2006-01-02 15:04:05")
		}
		if value != nil {
			defaultValue := &parserdriver.ValueExpr{}
			defaultValue.SetValue(value)
			return defaultValue, nil
		}
	}
	return nil, nil
}

// 还原抽象语法树节点至SQL
func restore(node ast.Node) (string, error) {
	var buf strings.Builder
	sql := ""
	rc := format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)

	if err := node.Restore(rc); err != nil {
		return "", err
	}
	sql = buf.String()
	return sql, nil
}

func getSchemaFromTableRefs(stmt *ast.Join) string {
	schema := ""
	if stmt == nil {
		return schema
	}
	if n := stmt.Left; n != nil {
		switch t := n.(type) {
		case *ast.TableSource:
			if tableName, ok := t.Source.(*ast.TableName); ok {
				schema = tableName.Schema.L
			}
		}
	}
	return schema
}
