//go:build enterprise
// +build enterprise

package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"

	"github.com/pingcap/parser/ast"
	parserMysql "github.com/pingcap/parser/mysql"
)

func (i *MysqlDriverImpl) GetOriginalRow(node ast.Node) (rollbackSql []string, unableRollbackReason i18nPkg.I18nStr, err error) {
	switch stmt := node.(type) {
	// row
	case *ast.DeleteStmt:
		return i.GetOriginalRowDelete(stmt)
	case *ast.UpdateStmt:
		return i.GetOriginalRowUpdate(stmt)
	case *ast.InsertStmt:
		return i.generateInsertRollbackSqls(stmt)
	// table
	case *ast.DropTableStmt:
		return i.generateDropTableRollbackSqls(stmt)
	case *ast.AlterTableStmt:
		return i.generateAlterTableRollbackSqls(stmt)
	case *ast.CreateTableStmt:
		return i.generateCreateTableRollbackSqls(stmt)
	// index
	case *ast.CreateIndexStmt:
		return i.generateCreateIndexRollbackSqls(stmt)
	case *ast.DropIndexStmt:
		return i.generateDropIndexRollbackSqls(stmt)
	// database
	case *ast.CreateDatabaseStmt:
		return i.generateCreateSchemaRollbackSqls(stmt)
	case *ast.DropDatabaseStmt:
		return i.generateDropDatabaseRollbackSqls(stmt)
	// other
	case *ast.UnparsedStmt:
		return []string{}, i18nPkg.ConvertStr2I18nAsDefaultLang("无法正常解析该SQL，无法进行备份"), nil
	default:
		return []string{}, i18nPkg.ConvertStr2I18nAsDefaultLang("暂不支持，该SQL的行备份"), nil
	}
}

// generateDeleteRollbackSql generate insert SQL for delete.
func (i *MysqlDriverImpl) GetOriginalRowDelete(stmt *ast.DeleteStmt) ([]string, i18nPkg.I18nStr, error) {
	tables := util.GetTables(stmt.TableRefs.TableRefs)
	if len(tables) == 0 {
		return []string{}, nil, fmt.Errorf("can not extract table from sql")
	}
	originRowInsertSql, err := i.getOriginalRowReplaceIntoSql(tables[0], stmt.Where, stmt.Order)
	if err != nil {
		return []string{}, nil, err
	}
	return originRowInsertSql, nil, nil
}

// generateDeleteRollbackSql generate insert SQL for delete.
func (i *MysqlDriverImpl) GetOriginalRowUpdate(stmt *ast.UpdateStmt) ([]string, i18nPkg.I18nStr, error) {
	tables := util.GetTables(stmt.TableRefs.TableRefs)
	if len(tables) == 0 {
		return []string{}, nil, fmt.Errorf("can not extract table from sql")
	}
	originRowInsertSql, err := i.getOriginalRowReplaceIntoSql(tables[0], stmt.Where, stmt.Order)
	if err != nil {
		return []string{}, nil, err
	}
	return originRowInsertSql, nil, nil
}

func (i *MysqlDriverImpl) GetCreateTableClause(tables []*ast.TableName) ([]string, i18nPkg.I18nStr, error) {
	createTableClauses := make([]string, 0, len(tables))
	for _, table := range tables {
		createTableClause, err := i.getCreateTableStmt(table)
		if err != nil {
			return []string{}, nil, err
		}
		createTableClauses = append(createTableClauses, createTableClause)
	}
	return createTableClauses, nil, nil
}

func (i *MysqlDriverImpl) getCreateTableStmt(table *ast.TableName) (string, error) {
	stmt, tableExist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil {
		return "", err
	}
	// if table not exist, can not rollback it.
	if !tableExist {
		return "", fmt.Errorf("table not exist")
	}
	return stmt.Text() + ";\n", nil
}

func (i *MysqlDriverImpl) getOriginalRowReplaceIntoSql(table *ast.TableName, whereClause ast.ExprNode, orderClause *ast.OrderByClause) ([]string, error) {
	createTableStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil || !exist {
		return []string{}, err
	}

	records, err := i.getRecords(table, "", whereClause, orderClause, 10000)
	if err != nil {
		return []string{}, err
	}

	columnNames := []string{}
	colNameDefMap := make(map[string]*ast.ColumnDef)
	for _, col := range createTableStmt.Cols {
		columnNames = append(columnNames, col.Name.Name.String())
		colNameDefMap[col.Name.Name.String()] = col
	}
	originRowInsertSqls := make([]string, 0, len(records))
	for _, record := range records {
		if len(record) != len(columnNames) {
			return []string{}, nil
		}
		values := []string{}
		for _, name := range columnNames {
			value := "NULL"
			if record[name].Valid {
				colDef := colNameDefMap[name]
				if parserMysql.HasBinaryFlag(colDef.Tp.Flag) {
					hexStr := getHexStrFromBytesStr(record[name].String)
					value = fmt.Sprintf("X'%s'", hexStr)
				} else {
					value = fmt.Sprintf("'%s'", record[name].String)
				}
				if colDef.Tp.Tp == parserMysql.TypeTimestamp {
					rowTime, err := time.Parse(time.RFC3339, record[name].String)
					if err != nil {
						value = fmt.Sprintf("'%s'", record[name].String)
					} else {
						value = fmt.Sprintf("'%s'", rowTime.Format("2006-01-02 15:04:05"))
					}
				}
			}
			values = append(values, value)
		}
		originRowInsertSqls = append(originRowInsertSqls,
			fmt.Sprintf("REPLACE INTO %s (`%s`) VALUES %s;",
				i.getTableNameWithQuote(table),
				strings.Join(columnNames, "`, `"),
				fmt.Sprintf("(%s)", strings.Join(values, ", ")),
			))
	}
	return originRowInsertSqls, nil
}
