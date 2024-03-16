package util

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/tidb/types"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/sirupsen/logrus"
)

var ErrUnsupportedSqlType = errors.New("unsupported sql type")

func GetAffectedRowNum(ctx context.Context, originSql string, conn *executor.Executor) (int64, error) {
	node, err := ParseOneSql(originSql)
	if err != nil {
		return 0, err
	}

	var newNode ast.Node
	var affectRowSql string
	var cannotConvert bool

	// 语法规则文档
	// select: https://dev.mysql.com/doc/refman/8.0/en/select.html
	// insert: https://dev.mysql.com/doc/refman/8.0/en/insert.html
	// update: https://dev.mysql.com/doc/refman/8.0/en/update.html
	// delete: https://dev.mysql.com/doc/refman/8.0/en/delete.html
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		isGroupByAndHavingBothExist := stmt.GroupBy != nil && stmt.Having != nil
		if stmt.GroupBy != nil || isGroupByAndHavingBothExist || stmt.Limit != nil {
			cannotConvert = true
		}

		newNode = getSelectNodeFromSelect(stmt)
	case *ast.InsertStmt:
		// 普通的insert语句，insert into t1 (name) values ('name1'), ('name2')
		isCommonInsert := stmt.Lists != nil && stmt.Select == nil
		// 包含子查询的insert语句，insert into t1 (name) select name from t2
		isSelectInsert := stmt.Select != nil && stmt.Lists == nil
		if isSelectInsert {
			newNode = getSelectNodeFromSelect(stmt.Select.(*ast.SelectStmt))
		} else if isCommonInsert {
			return int64(len(stmt.Lists)), nil
		} else {
			return 0, ErrUnsupportedSqlType
		}
	case *ast.UpdateStmt:
		newNode = getSelectNodeFromUpdate(stmt)
	case *ast.DeleteStmt:
		newNode = getSelectNodeFromDelete(stmt)
	default:
		return 0, ErrUnsupportedSqlType
	}

	// 1. 存在group by或者group by和having都存在的select语句，无法转换为select count语句
	// 2. SELECT COUNT(1) FROM test LIMIT 10,10 类型的SQL结果集为空
	// 已上两种情况,使用子查询 select count(*) from (输入的sql) as t的方式来获取影响行数
	if cannotConvert {
		// 将select语句中的查询字段替换为数字1
		// https://github.com/actiontech/sqle/issues/2175
		newSql, err := useIntReplaceSelectFields(node)
		if err != nil {
			log.NewEntry().Errorf("replace select fields failed, err: %v", err)
			newSql = originSql
		}
		// 移除后缀分号，避免sql语法错误
		trimSuffix := strings.TrimRight(newSql, ";")
		affectRowSql = fmt.Sprintf("select count(*) from (%s) as t", trimSuffix)
	} else {
		sqlBuilder := new(strings.Builder)
		err = newNode.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, sqlBuilder))
		if err != nil {
			return 0, err
		}

		affectRowSql = sqlBuilder.String()
	}

	// 验证sql语法是否正确，select 字段是否有且仅有 count(*)
	// 避免在客户机器上执行不符合预期的sql语句
	err = checkSql(affectRowSql)
	if err != nil {
		return 0, fmt.Errorf("check sql(%v) failed, origin sql(%v), err: %v", affectRowSql, originSql, err)
	}

	_, row, err := conn.Db.QueryWithContext(ctx, affectRowSql)
	if err != nil {
		return 0, err
	}

	// 如果下发的 SELECT COUNT(1) 的SQL，返回的结果集为空, 则返回0
	// 例: SELECT COUNT(1) FROM test LIMIT 10,10 结果集为空
	if len(row) == 0 {
		log.NewEntry().Errorf("affected row sql(%v) result row count is 0", affectRowSql)
		return 0, nil
	}

	if len(row) < 1 {
		return 0, fmt.Errorf("affected row sql(%v) result row count(%v) less than 1", affectRowSql, len(row))
	}

	affectCount, err := strconv.ParseInt(row[0][0].String, 10, 64)
	if err != nil {
		return 0, err
	}

	return affectCount, nil
}

func useIntReplaceSelectFields(node ast.StmtNode) (string, error) {
	stmt, ok := node.(*ast.SelectStmt)
	if !ok {
		return "", errors.New("pass parameter is not select node")
	}
	newValue := &driver.ValueExpr{}
	newValue.SetInt64(1)
	selectFields := &ast.SelectField{Expr: newValue}
	stmt.Fields.Fields = []*ast.SelectField{selectFields}

	sqlBuilder := new(strings.Builder)
	err := stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, sqlBuilder))
	if err != nil {
		return "", err
	}
	affectRowSql := sqlBuilder.String()
	return affectRowSql, nil
}

func getSelectNodeFromDelete(stmt *ast.DeleteStmt) *ast.SelectStmt {
	newSelect := newSelectWithCount()

	if stmt.TableRefs != nil {
		newSelect.From = stmt.TableRefs
	}

	if stmt.Where != nil {
		newSelect.Where = stmt.Where
	}

	if stmt.Order != nil {
		newSelect.OrderBy = stmt.Order
	}

	if stmt.Limit != nil {
		newSelect.Limit = stmt.Limit
	}

	return newSelect
}

func getSelectNodeFromUpdate(stmt *ast.UpdateStmt) *ast.SelectStmt {
	newSelect := newSelectWithCount()

	if stmt.TableRefs != nil {
		newSelect.From = stmt.TableRefs
	}

	if stmt.Where != nil {
		newSelect.Where = stmt.Where
	}

	if stmt.Order != nil {
		newSelect.OrderBy = stmt.Order
	}

	if stmt.Limit != nil {
		newSelect.Limit = stmt.Limit
	}

	return newSelect
}

func getSelectNodeFromSelect(stmt *ast.SelectStmt) *ast.SelectStmt {
	newSelect := newSelectWithCount()

	// todo: hint
	// todo: union
	if stmt.From != nil {
		newSelect.From = stmt.From
	}

	if stmt.Where != nil {
		newSelect.Where = stmt.Where
	}

	if stmt.OrderBy != nil {
		newSelect.OrderBy = stmt.OrderBy
	}

	if stmt.Limit != nil {
		newSelect.Limit = stmt.Limit
	}

	return newSelect
}

func newSelectWithCount() *ast.SelectStmt {
	newSelect := new(ast.SelectStmt)
	a := new(ast.SelectStmtOpts)
	a.SQLCache = true
	newSelect.SelectStmtOpts = a

	newSelect.Fields = getCountFieldList()
	return newSelect
}

// getCountFieldList
// 获取count(*)函数的字段列表
func getCountFieldList() *ast.FieldList {
	datum := new(types.Datum)
	datum.SetInt64(1)

	return &ast.FieldList{
		Fields: []*ast.SelectField{
			{
				Expr: &ast.AggregateFuncExpr{
					F: ast.AggFuncCount,
					Args: []ast.ExprNode{
						&driver.ValueExpr{
							Datum: *datum,
						},
					},
				},
			},
		},
	}
}

func checkSql(affectRowSql string) error {
	node, err := ParseOneSql(affectRowSql)
	if err != nil {
		return err
	}

	fieldExtractor := new(SelectFieldExtractor)
	node.Accept(fieldExtractor)

	if !fieldExtractor.IsSelectOnlyIncludeCountFunc {
		return errors.New("affectRowSql error")
	}

	return nil
}

func KillProcess(ctx context.Context, killSQL string, killConn *executor.Executor, logEntry *logrus.Entry) error {
	killFunc := func() error {
		_, err := killConn.Db.Exec(killSQL)
		return err
	}
	err := utils.AsyncCallTimeout(ctx, killFunc)
	if err != nil {
		err = fmt.Errorf("exec sql(%v) failed, err: %v", killSQL, err)
		return err
	}
	logEntry.Infof("exec sql(%v) successfully", killSQL)
	return nil
}

func IsGeometryColumn(col *ast.ColumnDef) bool {
	switch col.Tp.Tp {
	case mysql.TypeGeometry, mysql.TypePoint, mysql.TypeLineString, mysql.TypePolygon,
		mysql.TypeMultiPoint, mysql.TypeMultiLineString, mysql.TypeMultiPolygon, mysql.TypeGeometryCollection:
		return true
	}
	return false
}

// TODO: 暂时使用正则表达式匹配event，后续会修改语法树进行匹配event
func IsEventSQL(sql string) bool {
	createPattern := `^CREATE\s+(DEFINER\s?=.+?)?EVENT`
	createRe := regexp.MustCompile(createPattern)
	alterPattern := `^ALTER\s+(DEFINER\s?=.+?)?EVENT`
	alterRe := regexp.MustCompile(alterPattern)

	sql = strings.ToUpper(strings.TrimSpace(sql))
	if createRe.MatchString(sql) {
		return true
	} else {
		return alterRe.MatchString(sql)
	}
}

func GetTableNameFromTableSource(tableSource *ast.TableSource) string {
	if tableSource == nil {
		return ""
	}
	if tableSource.AsName.L != "" {
		return tableSource.AsName.L
	}
	if tableName, ok := tableSource.Source.(*ast.TableName); ok {
		return tableName.Name.L
	}
	return ""
}

/*
IsIndex

	判断单列或多列是否属于索引切片中的索引：
		1. 单列：满足单列索引或多列索引的第一列，则返回true
		2. 多列：满足N列是M列索引的前N列（M>=N），则返回true
		3. 否则返回false
*/
func IsIndex(columnMap map[string] /*column name*/ struct{}, constraints []*ast.Constraint) bool {
	if len(columnMap) == 0 {
		return false
	}
	for _, constraint := range constraints {
		if len(columnMap) > len(constraint.Keys) {
			// 若符合索引的列数小于关联列的列数 一定不满足多列索引
			continue
		}
		var matchCount int
		for _, key := range constraint.Keys {
			if _, ok := columnMap[key.Column.Name.L]; ok {
				matchCount++
			} else {
				break
			}
		}
		if matchCount == len(columnMap) {
			return true
		}
	}
	return false
}
