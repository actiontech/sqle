package mysql

import (
	"context"
	"fmt"
	"strings"

	dmsCommonSQLOp "github.com/actiontech/dms/pkg/dms-common/sql_op"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

// GetSQLOp 获取sql中涉及的对象操作，sql可以是单条语句，也可以是多条语句
// GetSQLOp 目前实现的对象的最小粒度只到表级别，考虑识别列对象的成本较高，暂时没有识别列对象
func (i *MysqlDriverImpl) GetSQLOp(ctx context.Context, sqls string) ([]*dmsCommonSQLOp.SQLObjectOps, error) {
	p := parser.New()
	stmts, _, err := p.PerfectParse(sqls, "", "")
	if err != nil {
		i.Logger().Errorf("parse sql failed, error: %v, sql: %s", err, sqls)
		return nil, err
	}

	ret := make([]*dmsCommonSQLOp.SQLObjectOps, 0, len(stmts))
	for _, stmt := range stmts {
		objectOps := dmsCommonSQLOp.NewSQLObjectOps(stmt.Text())
		err := parseStmtOpInfos(p, stmt, objectOps)
		if err != nil {
			i.Logger().Errorf("parse sql op failed, error: %v, sql: %s", err, stmt.Text())
			return nil, err
		}
		ret = append(ret, objectOps)
	}

	for i := range ret {
		ret[i].ObjectOps = SQLObjectOpsDuplicateRemoval(ret[i].ObjectOps)
	}

	return ret, nil
}

func parseStmtOpInfos(p *parser.Parser, stmt ast.StmtNode, ops *dmsCommonSQLOp.SQLObjectOps) error {

	switch s := stmt.(type) {
	case *ast.CreateTableStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/create-table.html
		// You must have the CREATE privilege for the table.
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.Table))
		if s.ReferTable != nil {
			ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpRead, s.ReferTable))
		}

	case *ast.DropTableStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/drop-table.html
		// You must have the DROP privilege for each table.
		for _, table := range s.Tables {
			ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpDelete, table))
		}

	case *ast.AlterTableStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
		// To use ALTER TABLE, you need ALTER, CREATE, and INSERT privileges for the table.
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.Table))

		// Renaming a table requires ALTER and DROP on the old table, ALTER, CREATE, and INSERT on the new table.
		for _, spec := range s.Specs {
			if spec.Tp == ast.AlterTableRenameTable {
				ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.Table))
				ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpDelete, s.Table))
				ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, spec.NewTable))
			}
		}

	case *ast.TruncateTableStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/truncate-table.html
		// It requires the DROP privilege.
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpDelete, s.Table))

	case *ast.RepairTableStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/repair-table.html
		// https://docs.pingcap.com/tidb/v5.4/sql-statement-admin/
		// TODO: TiDB的解析器(ADMIN REPAIR TABLE)不支持解析MySQL的REPAIR TABLE语句
		return fmt.Errorf("repair table not supported")

	case *ast.CreateDatabaseStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/create-database.html
		// To use this statement, you need the CREATE privilege for the database.
		ops.AddObjectOp(newDatabaseObject(dmsCommonSQLOp.SQLOpAddOrUpdate, s.Name))

	case *ast.AlterDatabaseStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/alter-database.html
		// This statement requires the ALTER privilege on the database.
		ops.AddObjectOp(newDatabaseObject(dmsCommonSQLOp.SQLOpAddOrUpdate, s.Name))

	case *ast.DropDatabaseStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/drop-database.html
		// To use DROP DATABASE, you need the DROP privilege on the database.
		ops.AddObjectOp(newDatabaseObject(dmsCommonSQLOp.SQLOpDelete, s.Name))

	case *ast.RenameTableStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/rename-table.html
		// You must have ALTER and DROP privileges for the original table, and CREATE and INSERT privileges for the new table
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.OldTable))
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpDelete, s.OldTable))
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.NewTable))

	case *ast.CreateViewStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/create-view.html
		// The CREATE VIEW statement requires the CREATE VIEW privilege for the view,
		// and some privilege for each column selected by the SELECT statement.
		// For columns used elsewhere in the SELECT statement, you must have the SELECT privilege.
		// If the OR REPLACE clause is present, you must also have the DROP privilege for the view.
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.ViewName))
		if s.OrReplace {
			ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpDelete, s.ViewName))
		}

		switch s := s.Select.(type) {
		case *ast.SelectStmt:
			o, err := getSelectStmtObjectOps(s)
			if nil != err {
				return fmt.Errorf("failed to parse create view, %v", err)
			}
			ops.AddObjectOp(o...)
		case *ast.UnionStmt:
			for _, selectStmt := range s.SelectList.Selects {
				o, err := getSelectStmtObjectOps(selectStmt)
				if nil != err {
					return fmt.Errorf("failed to parse create view, %v", err)
				}
				ops.AddObjectOp(o...)
			}
		default:
			return fmt.Errorf("failed to parse create view, not support select type: %T", s)
		}

	case *ast.CreateIndexStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/create-index.html
		// CREATE INDEX is mapped to an ALTER TABLE statement to create indexes.
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.Table))

	case *ast.DropIndexStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/drop-index.html
		// This statement is mapped to an ALTER TABLE statement to drop the index.
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.Table))

	case *ast.LockTablesStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/lock-tables.html
		// You must have the LOCK TABLES privilege, and the SELECT privilege for each object to be locked.
		for _, t := range s.TableLocks {
			ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAdmin, t.Table))
			ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpRead, t.Table))
		}

		// TODO: 视图处理
		// For view locking, LOCK TABLES adds all base tables used in the view to the set of tables to be locked and locks them automatically.
		// For tables underlying any view being locked, LOCK TABLES checks that the view definer (for SQL SECURITY DEFINER views) or invoker (for all views) has the proper privileges on the tables.

	case *ast.UnlockTablesStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/lock-tables.html
		// UNLOCK TABLES explicitly releases any table locks held by the current session.
		// TODO: Unlock table的操作针对的当前session，无法从sql语句中确定具体的表，暂且不处理
	case *ast.SelectStmt:

		o, err := getSelectStmtObjectOps(s)
		if nil != err {
			return fmt.Errorf("failed to parse select, %v", err)
		}
		ops.AddObjectOp(o...)

	case *ast.UnionStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/union.html

		for _, selectStmt := range s.SelectList.Selects {
			o, err := getSelectStmtObjectOps(selectStmt)
			if nil != err {
				return fmt.Errorf("failed to parse union, %v", err)
			}
			ops.AddObjectOp(o...)
		}

		// TiDB的解析器尚不支持MySQL8.0的一些新union语法，暂且不处理
		// 如(SELECT 1 UNION SELECT 1) UNION SELECT 1;

	case *ast.LoadDataStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/load-data.html
		// You must have the FILE privilege
		// For a LOCAL load operation, the client program reads a text file located on the client host.
		// Because the file contents are sent over the connection by the client to the server,
		// using LOCAL is a bit slower than when the server accesses the file directly.
		// On the other hand, you do not need the FILE privilege,
		// and the file can be located in any directory the client program can access.
		if !s.IsLocal {
			ops.AddObjectOp(newServerObject(dmsCommonSQLOp.SQLOpAdmin))
		}
		ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, s.Table))

	case *ast.InsertStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/insert.html
		// Inserting into a table requires the INSERT privilege for the table.
		// If the ON DUPLICATE KEY UPDATE clause is used and a duplicate key causes an UPDATE to be performed instead,
		// the statement requires the UPDATE privilege for the columns to be updated.
		// For columns that are read but not modified you need only the SELECT privilege
		// (such as for a column referenced only on the right hand side of an col_name=expr assignment in an ON DUPLICATE KEY UPDATE clause).
		tableNames := util.GetTables(s.Table.TableRefs)
		for _, tableName := range tableNames {
			ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpAddOrUpdate, tableName))
		}

		if s.Select != nil {
			if selectStmt, ok := s.Select.(*ast.SelectStmt); ok {
				o, err := getSelectStmtObjectOps(selectStmt)
				if nil != err {
					return fmt.Errorf("failed to parse insert, %v", err)
				}
				ops.AddObjectOp(o...)
			} else {
				return fmt.Errorf("failed to parse insert, not support select type: %T", s.Select)
			}
		}

	case *ast.DeleteStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/delete.html
		// You need the DELETE privilege on a table to delete rows from it.
		// You need only the SELECT privilege for any columns that are only read,
		// such as those named in the WHERE clause.

		var tableAlias []*tableAliasInfo
		if s.IsMultiTable {
			for _, table := range s.Tables.Tables {
				ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpDelete, table))
			}
			tableAlias = getTableAliasInfoFromTableNames(s.Tables.Tables)
		} else {
			tableNames := util.GetTables(s.TableRefs.TableRefs)
			for _, tableName := range tableNames {
				ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpDelete, tableName))
			}

			tableAlias = getTableAliasInfoFromJoin(s.TableRefs.TableRefs)
		}

		if s.Where != nil {
			ops.AddObjectOp(getTableObjectOpsInWhere(tableAlias, s.Where)...)
		}

	case *ast.UpdateStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/update.html
		// You need the UPDATE privilege only for columns referenced in an UPDATE that are actually updated.
		// You need only the SELECT privilege for any columns that are read but not modified.

		tableAlias := getTableAliasInfoFromJoin(s.TableRefs.TableRefs)
		// UPDATE items,month SET items.price=month.price
		// WHERE items.id=month.id;
		// 对 items 表有 更新和读取 ，对 month 表只有读取
		for _, list := range s.List {
			if list.Column != nil {
				ops.AddObjectOp(newTableObjectFromColumn(dmsCommonSQLOp.SQLOpAddOrUpdate, list.Column, tableAlias))
			}
			if c, ok := list.Expr.(*ast.ColumnNameExpr); ok {
				ops.AddObjectOp(newTableObjectFromColumn(dmsCommonSQLOp.SQLOpRead, c.Name, tableAlias))
			}
		}

		if s.Where != nil {
			ops.AddObjectOp(getTableObjectOpsInWhere(tableAlias, s.Where)...)
		}

	case *ast.ShowStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/show.html
		switch s.Tp {
		case ast.ShowEngines:
		case ast.ShowDatabases:
			// You see only those databases for which you have some kind of privilege,
			// unless you have the global SHOW DATABASES privilege.
		case ast.ShowTables, ast.ShowTableStatus:
		case ast.ShowColumns:
			// https://dev.mysql.com/doc/refman/8.0/en/show-columns.html
			dbName := s.Table.Schema.L
			if s.DBName != "" {
				dbName = s.DBName
			}
			ops.AddObjectOp(newTableObject(dmsCommonSQLOp.SQLOpRead, s.Table.Name.L, dbName))
		case ast.ShowWarnings, ast.ShowErrors:
		case ast.ShowCharset, ast.ShowCollation:
		case ast.ShowVariables, ast.ShowStatus:
		case ast.ShowCreateTable, ast.ShowCreateView:
			// https://dev.mysql.com/doc/refman/8.0/en/show-create-table.html
			ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpRead, s.Table))
		case ast.ShowCreateUser:
			// https://dev.mysql.com/doc/refman/8.0/en/show-create-user.html
			// The statement requires the SELECT privilege for the mysql system schema,
			// except to see information for the current user.
			ops.AddObjectOp(newMysqlDatabaseObject(dmsCommonSQLOp.SQLOpRead))
		case ast.ShowGrants:
			// https://dev.mysql.com/doc/refman/8.0/en/show-grants.html
			ops.AddObjectOp(newMysqlDatabaseObject(dmsCommonSQLOp.SQLOpRead))
		case ast.ShowTriggers:
		case ast.ShowProcedureStatus:
			// https://dev.mysql.com/doc/refman/8.0/en/show-procedure-status.html
			ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpRead))
		case ast.ShowIndex:
			ops.AddObjectOp(newTableObjectFromTableName(dmsCommonSQLOp.SQLOpRead, s.Table))
		case ast.ShowProcessList:
		case ast.ShowCreateDatabase:
			ops.AddObjectOp(newDatabaseObject(dmsCommonSQLOp.SQLOpRead, s.DBName))
		case ast.ShowEvents:
			// https://dev.mysql.com/doc/refman/5.7/en/show-events.html
			// It requires the EVENT privilege for the database from which the events are to be shown.
			ops.AddObjectOp(newDatabaseObject(dmsCommonSQLOp.SQLOpRead, s.DBName))
		case ast.ShowPlugins:
		case ast.ShowProfile, ast.ShowProfiles:
		case ast.ShowMasterStatus:
			// https://dev.mysql.com/doc/refman/5.7/en/show-master-status.html
			// It requires either the SUPER or REPLICATION CLIENT privilege.
			ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
		case ast.ShowPrivileges:
		default:
			return fmt.Errorf("failed to parse show, not support show type: %v", s.Tp)
		}
	case *ast.ExplainForStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/explain.html
		// EXPLAIN ... FOR CONNECTION also requires the PROCESS privilege if the specified connection belongs to a different user.
		// 由于无法判断FOR CONNECTION的连接是否属于当前用户，所以这里定义为整个实例的读
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpRead))
	case *ast.ExplainStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/explain.html
		// EXPLAIN requires the same privileges required to execute the explained statement.
		if err := parseStmtOpInfos(p, s.Stmt, ops); err != nil {
			return err
		}
		// TODO:
		// Additionally, EXPLAIN also requires the SHOW VIEW privilege for any explained view.
	case *ast.PrepareStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/sql-prepared-statements.html
		stmt, err := p.ParseOneStmt(s.SQLText, "", "")
		if nil != err {
			return err
		}
		if err := parseStmtOpInfos(p, stmt, ops); err != nil {
			return err
		}

		// TODO: prepare 可能使用上下文来定义语句，如下：
		// SET @s = 'SELECT SQRT(POW(?,2) + POW(?,2)) AS hypotenuse';
		// PREPARE stmt2 FROM @s;
		// 这种情况需要结合上下文来解析
	case *ast.ExecuteStmt:
	case *ast.DeallocateStmt:
	case *ast.BeginStmt:
	case *ast.CommitStmt:
	case *ast.RollbackStmt:
	case *ast.BinlogStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/binlog.html
		// To execute BINLOG statements when applying mysqlbinlog output,
		// a user account requires the BINLOG_ADMIN privilege (or the deprecated SUPER privilege),
		// or the REPLICATION_APPLIER privilege plus the appropriate privileges to execute each log event.
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
	case *ast.UseStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/use.html
		// This statement requires some privilege for the database or some object within it.
	case *ast.FlushStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/flush.html
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
	case *ast.KillStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/kill.html
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
	case *ast.SetStmt:
	case *ast.SetPwdStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/set-password.html
		if !s.User.CurrentUser {
			ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
		}
	case *ast.CreateUserStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/create-user.html
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
	case *ast.AlterUserStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/alter-user.html
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
	case *ast.AlterInstanceStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/alter-instance.html
		// TiDB的解析器仅支持ALTER INSTANCE RELOAD TLS语句
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
	case *ast.DropUserStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/drop-user.html
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
	case *ast.RevokeStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/revoke.html
		switch s.Level.Level {
		case ast.GrantLevelGlobal:
			ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpGrant))
		case ast.GrantLevelDB:
			ops.AddObjectOp(newDatabaseObject(dmsCommonSQLOp.SQLOpGrant, s.Level.DBName))
		case ast.GrantLevelTable:
			ops.AddObjectOp(newTableObject(dmsCommonSQLOp.SQLOpGrant, s.Level.TableName, s.Level.DBName))
		default:
			return fmt.Errorf("not support grant level: %v", s.Level.Level)
		}
	case *ast.GrantStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/grant.html
		switch s.Level.Level {
		case ast.GrantLevelGlobal:
			ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpGrant))
		case ast.GrantLevelDB:
			ops.AddObjectOp(newDatabaseObject(dmsCommonSQLOp.SQLOpGrant, s.Level.DBName))
		case ast.GrantLevelTable:
			ops.AddObjectOp(newTableObject(dmsCommonSQLOp.SQLOpGrant, s.Level.TableName, s.Level.DBName))
		default:
			return fmt.Errorf("not support grant level: %v", s.Level.Level)
		}
	case *ast.ShutdownStmt:
		// https://dev.mysql.com/doc/refman/8.0/en/shutdown.html
		ops.AddObjectOp(newInstanceObject(dmsCommonSQLOp.SQLOpAdmin))
	case *ast.UnparsedStmt:
		return fmt.Errorf("there is unparsed stmt: %s", stmt.Text())
	default:
		return fmt.Errorf("not support stmt type: %T", stmt)
	}
	return nil
}

func getSelectStmtObjectOps(selectStmt *ast.SelectStmt) ([]*dmsCommonSQLOp.SQLObjectOp, error) {
	ret := make([]*dmsCommonSQLOp.SQLObjectOp, 0)
	if selectStmt.From != nil {
		tableNames := util.GetTables(selectStmt.From.TableRefs)
		for _, tableName := range tableNames {
			ret = append(ret, newTableObjectFromTableName(dmsCommonSQLOp.SQLOpRead, tableName))
		}
	}
	if selectStmt.SelectIntoOpt != nil {
		// https://dev.mysql.com/doc/refman/8.0/en/select-into.html
		// The SELECT ... INTO OUTFILE 'file_name' form of SELECT writes the selected rows to a file.
		// The file is created on the server host,
		// so you must have the FILE privilege to use this syntax.
		if selectStmt.SelectIntoOpt.Tp == ast.SelectIntoOutfile || selectStmt.SelectIntoOpt.Tp == ast.SelectIntoDumpfile {
			ret = append(ret, newServerObject(dmsCommonSQLOp.SQLOpAdmin))
		}
	}
	return ret, nil
}

func newTableObjectFromTableName(op dmsCommonSQLOp.SQLOp, table *ast.TableName) *dmsCommonSQLOp.SQLObjectOp {
	return &dmsCommonSQLOp.SQLObjectOp{
		Op:     op,
		Object: newTable(table.Name.L, table.Schema.L),
	}
}

func newTableObject(op dmsCommonSQLOp.SQLOp, tableName, databaseName string) *dmsCommonSQLOp.SQLObjectOp {
	return &dmsCommonSQLOp.SQLObjectOp{
		Op:     op,
		Object: newTable(tableName, databaseName),
	}
}

type tableAliasInfo struct {
	tableName      string
	schemaName     string
	tableAliasName string
}

func newTableObjectFromColumn(op dmsCommonSQLOp.SQLOp, column *ast.ColumnName, tableAliasInfo []*tableAliasInfo) *dmsCommonSQLOp.SQLObjectOp {
	if column.Table.String() == "" && len(tableAliasInfo) == 1 {
		t := tableAliasInfo[0]
		return newTableObject(op, t.tableName, t.schemaName)
	}
	for _, t := range tableAliasInfo {
		if t.tableAliasName == column.Table.String() {
			return newTableObject(op, t.tableName, t.schemaName)
		}
	}
	return &dmsCommonSQLOp.SQLObjectOp{
		Op:     op,
		Object: newTable(column.Table.L, column.Schema.L),
	}
}

func newDatabaseObject(op dmsCommonSQLOp.SQLOp, databaseName string) *dmsCommonSQLOp.SQLObjectOp {
	return &dmsCommonSQLOp.SQLObjectOp{
		Op:     op,
		Object: newDatabase(databaseName),
	}
}

func newMysqlDatabaseObject(op dmsCommonSQLOp.SQLOp) *dmsCommonSQLOp.SQLObjectOp {
	return newDatabaseObject(op, "mysql")
}

func newInstanceObject(op dmsCommonSQLOp.SQLOp) *dmsCommonSQLOp.SQLObjectOp {
	return &dmsCommonSQLOp.SQLObjectOp{
		Op:     op,
		Object: newInstance(),
	}
}

func newServerObject(op dmsCommonSQLOp.SQLOp) *dmsCommonSQLOp.SQLObjectOp {
	return &dmsCommonSQLOp.SQLObjectOp{
		Op:     op,
		Object: newServer(),
	}
}

func newTable(tableName, databaseName string) *dmsCommonSQLOp.SQLObject {
	return &dmsCommonSQLOp.SQLObject{
		Type:         dmsCommonSQLOp.SQLObjectTypeTable,
		DatabaseName: databaseName,
		SchemaName:   "", // MySQL中schema与database是同一个概念，已设置databaseName，schemaName不存在意义，所以设置为空
		TableName:    tableName,
	}
}

func newDatabase(databaseName string) *dmsCommonSQLOp.SQLObject {
	return &dmsCommonSQLOp.SQLObject{
		Type:         dmsCommonSQLOp.SQLObjectTypeDatabase,
		DatabaseName: databaseName,
		SchemaName:   "", // MySQL中schema与database是同一个概念，已设置databaseName，schemaName不存在意义，所以设置为空
		TableName:    "",
	}
}

func newInstance() *dmsCommonSQLOp.SQLObject {
	return &dmsCommonSQLOp.SQLObject{
		Type:         dmsCommonSQLOp.SQLObjectTypeInstance,
		DatabaseName: "",
		SchemaName:   "",
		TableName:    "",
	}
}

func newServer() *dmsCommonSQLOp.SQLObject {
	return &dmsCommonSQLOp.SQLObject{
		Type:         dmsCommonSQLOp.SQLObjectTypeServer,
		DatabaseName: "",
		SchemaName:   "",
		TableName:    "",
	}
}

func getTableAliasInfoFromTableNames(tableNames []*ast.TableName) []*tableAliasInfo {
	tableAlias := make([]*tableAliasInfo, 0)
	for _, tableName := range tableNames {
		tableAlias = append(tableAlias, &tableAliasInfo{
			tableAliasName: "",
			tableName:      tableName.Name.L,
			schemaName:     tableName.Schema.L,
		})
	}
	return tableAlias
}

func getTableAliasInfoFromJoin(stmt *ast.Join) []*tableAliasInfo {
	tableAlias := make([]*tableAliasInfo, 0)
	tableSources := util.GetTableSources(stmt)
	for _, tableSource := range tableSources {
		switch source := tableSource.Source.(type) {
		case *ast.TableName:
			tableAlias = append(tableAlias, &tableAliasInfo{
				tableAliasName: tableSource.AsName.String(),
				tableName:      source.Name.L,
				schemaName:     source.Schema.L,
			})
		default:
		}
	}
	return tableAlias
}

func getTableObjectOpsInWhere(tableAliasInfo []*tableAliasInfo, where ast.ExprNode) []*dmsCommonSQLOp.SQLObjectOp {
	c := &tableObjectsInWhere{
		tableAliasInfo: tableAliasInfo,
	}
	where.Accept(c)
	return c.tables
}

type tableObjectsInWhere struct {
	tables         []*dmsCommonSQLOp.SQLObjectOp
	tableAliasInfo []*tableAliasInfo
}

func (c *tableObjectsInWhere) Enter(in ast.Node) (ast.Node, bool) {
	if cn, ok := in.(*ast.ColumnName); ok {
		c.tables = append(c.tables, newTableObjectFromColumn(dmsCommonSQLOp.SQLOpRead, cn, c.tableAliasInfo))
	}
	return in, false
}

func (c *tableObjectsInWhere) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func SQLObjectOpsDuplicateRemoval(ops []*dmsCommonSQLOp.SQLObjectOp) []*dmsCommonSQLOp.SQLObjectOp {
	m := make(map[string]*dmsCommonSQLOp.SQLObjectOp)
	for _, o := range ops {
		m[SQLObjectOpFingerPrint(o)] = o
	}
	ret := make([]*dmsCommonSQLOp.SQLObjectOp, 0)
	for _, o := range m {
		ret = append(ret, o)
	}
	return ret
}

func SQLObjectOpsFingerPrint(ops []*dmsCommonSQLOp.SQLObjectOps) string {
	s := make([]string, len(ops))
	for i := range ops {
		s[i] = SQLObjectOpFingerPrints(ops[i].ObjectOps)
		s[i] = fmt.Sprintf("%s %s", s[i], ops[i].Sql.Sql)
	}
	return strings.Join(s, "\n")
}

func SQLObjectOpFingerPrints(ops []*dmsCommonSQLOp.SQLObjectOp) string {
	s := make([]string, len(ops))
	for i := range ops {
		s[i] = SQLObjectOpFingerPrint(ops[i])
	}
	return strings.Join(s, ";")
}

func SQLObjectOpFingerPrint(op *dmsCommonSQLOp.SQLObjectOp) string {
	return fmt.Sprintf("%s %s", SQLObjectFingerPrint(op.Object), op.Op)
}

func SQLObjectFingerPrint(obj *dmsCommonSQLOp.SQLObject) string {
	switch obj.Type {
	case dmsCommonSQLOp.SQLObjectTypeInstance:
		return "*.*"
	case dmsCommonSQLOp.SQLObjectTypeDatabase:
		return fmt.Sprintf("%s.*", obj.DatabaseName)
	case dmsCommonSQLOp.SQLObjectTypeTable:
		return fmt.Sprintf("%s.%s", obj.DatabaseName, obj.TableName)
	default:
		return "unknown"
	}
}
