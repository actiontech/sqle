package mysql

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/pingcap/parser/ast"
)

const (
	SchemaNotExistMessage              = "schema %s 不存在"
	SchemaExistMessage                 = "schema %s 已存在"
	TableNotExistMessage               = "表 %s 不存在"
	TableExistMessage                  = "表 %s 已存在"
	ColumnNotExistMessage              = "字段 %s 不存在"
	ColumnExistMessage                 = "字段 %s 已存在"
	ColumnIsAmbiguousMessage           = "字段 %s 指代不明"
	IndexNotExistMessage               = "索引 %s 不存在"
	IndexExistMessage                  = "索引 %s 已存在"
	DuplicateColumnsMessage            = "字段名 %s 重复"
	DuplicateIndexesMessage            = "索引名 %s 重复"
	MultiPrimaryKeyMessage             = "主键只能设置一个"
	KeyedColumnNotExistMessage         = "索引字段 %s 不存在"
	PrimaryKeyExistMessage             = "已经存在主键，不能再添加"
	PrimaryKeyNotExistMessage          = "当前没有主键，不能执行删除"
	ColumnsValuesNotMatchMessage       = "指定的值列数与字段列数不匹配"
	DuplicatePrimaryKeyedColumnMessage = "主键字段 %s 重复"
	DuplicateIndexedColumnMessage      = "索引 %s 字段 %s重复"
)

const CheckInvalidErrorFormat = "预检查失败: %v"

func (i *MysqlDriverImpl) CheckInvalid(node ast.Node) error {
	var err error
	switch stmt := node.(type) {
	case *ast.UseStmt:
		err = i.checkInvalidUse(stmt)
	case *ast.CreateTableStmt:
		err = i.checkInvalidCreateTable(stmt)
	case *ast.AlterTableStmt:
		err = i.checkInvalidAlterTable(stmt)
	case *ast.DropTableStmt:
		err = i.checkInvalidDropTable(stmt)
	case *ast.CreateDatabaseStmt:
		err = i.checkInvalidCreateDatabase(stmt)
	case *ast.DropDatabaseStmt:
		err = i.checkInvalidDropDatabase(stmt)
	case *ast.CreateIndexStmt:
		err = i.checkInvalidCreateIndex(stmt)
	case *ast.DropIndexStmt:
		err = i.checkInvalidDropIndex(stmt)
	case *ast.InsertStmt:
		err = i.checkInvalidInsert(stmt)
	case *ast.UpdateStmt:
		err = i.checkInvalidUpdate(stmt)
	case *ast.DeleteStmt:
		err = i.checkInvalidDelete(stmt)
	case *ast.SelectStmt:
		err = i.checkInvalidSelect(stmt)
	case *ast.UnparsedStmt:
		err = i.checkUnparsedStmt(stmt)
	}
	if err != nil {
		return fmt.Errorf(CheckInvalidErrorFormat, err)
	}
	return nil

}

func (i *MysqlDriverImpl) CheckExplain(node ast.Node) error {
	var err error
	switch node.(type) {
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.SelectStmt:
		if i.Ctx.GetHistorySQLInfo().HasDDL {
			return nil
		}
		_, err = i.Ctx.GetExecutionPlan(node.Text())
	}
	if err != nil {
		i.result.Add(driver.RuleLevelWarn, fmt.Sprintf(CheckInvalidErrorFormat, err))
	}
	return nil

}

func (i *MysqlDriverImpl) CheckInvalidOffline(node ast.Node) error {
	stmt, ok := node.(*ast.UnparsedStmt)
	if ok {
		return i.checkUnparsedStmt(stmt)
	}
	return nil
}

/*
------------------------------------------------------------------
create table ...
------------------------------------------------------------------
1. schema must exist;
2. table can't exist if SQL has not "IF NOT EXISTS";
3. column name can't duplicated;
4. primary key can only be set once;
5. index name can't be duplicated;
6. index column must exist;
7. index column can't duplicated, "index idx_1(id,id)" is invalid
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidCreateTable(stmt *ast.CreateTableStmt) error {
	schemaName := i.Ctx.GetSchemaName(stmt.Table)
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage, schemaName)
	} else {
		tableExist, err := i.Ctx.IsTableExist(stmt.Table)
		if err != nil {
			return err
		}
		if tableExist && !stmt.IfNotExists {
			i.result.Add(driver.RuleLevelError, TableExistMessage,
				i.getTableName(stmt.Table))
		}
		if stmt.ReferTable != nil {
			referTableExist, err := i.Ctx.IsTableExist(stmt.ReferTable)
			if err != nil {
				return err
			}
			if !referTableExist {
				i.result.Add(driver.RuleLevelError, TableNotExistMessage,
					i.getTableName(stmt.ReferTable))
			}
		}
	}
	colsName := []string{}
	colsNameMap := map[string]struct{}{}
	pkCounter := 0
	for _, col := range stmt.Cols {
		colName := col.Name.Name.L
		colsName = append(colsName, colName)
		colsNameMap[colName] = struct{}{}
		if util.HasOneInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
			pkCounter += 1
		}
	}
	indexesName := []string{}
	keyColsName := []string{}
	for _, constraint := range stmt.Constraints {
		switch constraint.Tp {
		case ast.ConstraintPrimaryKey:
			pkCounter += 1
			names := []string{}
			for _, col := range constraint.Keys {
				colName := col.Column.Name.L
				names = append(names, colName)
				keyColsName = append(keyColsName, colName)
			}
			duplicateName := utils.GetDuplicate(names)
			if len(duplicateName) > 0 {
				i.result.Add(driver.RuleLevelError, DuplicatePrimaryKeyedColumnMessage,
					strings.Join(duplicateName, ","))
			}
		case ast.ConstraintIndex, ast.ConstraintUniq, ast.ConstraintFulltext:
			constraintName := constraint.Name
			if constraintName != "" {
				indexesName = append(indexesName, constraint.Name)
			} else {
				constraintName = "(匿名)"
			}
			names := []string{}
			for _, col := range constraint.Keys {
				colName := col.Column.Name.L
				names = append(names, colName)
				keyColsName = append(keyColsName, colName)
			}
			duplicateName := utils.GetDuplicate(names)
			if len(duplicateName) > 0 {
				i.result.Add(driver.RuleLevelError, DuplicateIndexedColumnMessage, constraintName,
					strings.Join(duplicateName, ","))
			}
		}
	}
	if d := utils.GetDuplicate(colsName); len(d) > 0 {
		i.result.Add(driver.RuleLevelError, DuplicateColumnsMessage,
			strings.Join(d, ","))
	}

	if d := utils.GetDuplicate(indexesName); len(d) > 0 {
		i.result.Add(driver.RuleLevelError, DuplicateIndexesMessage,
			strings.Join(d, ","))
	}

	if pkCounter > 1 {
		i.result.Add(driver.RuleLevelError, MultiPrimaryKeyMessage)
	}
	notExistKeyColsName := []string{}
	for _, colName := range keyColsName {
		if _, ok := colsNameMap[colName]; !ok {
			notExistKeyColsName = append(notExistKeyColsName, colName)
		}
	}
	if len(notExistKeyColsName) > 0 {
		i.result.Add(driver.RuleLevelError, KeyedColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(notExistKeyColsName), ","))
	}
	return nil
}

/*
------------------------------------------------------------------
alter table ...
------------------------------------------------------------------
1. schema must exist;
2. table must exist;
3. add/update column, name can't duplicated;
4. delete column, name must exist;
5. add/update pk, pk can only be set once;
6. delete pk, pk must exist;
7. add/update index, name can't be duplicated;
8. delete index, name must exist;
9. index column must exist;
10. index column can't duplicated.
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidAlterTable(stmt *ast.AlterTableStmt) error {
	schemaName := i.Ctx.GetSchemaName(stmt.Table)
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.Ctx.GetCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		i.result.Add(driver.RuleLevelError, TableNotExistMessage,
			i.getTableName(stmt.Table))
		return nil
	}

	hasPk := false
	colNameMap := map[string]struct{}{}
	indexNameMap := map[string]struct{}{}
	for _, col := range createTableStmt.Cols {
		colNameMap[col.Name.Name.L] = struct{}{}
		if util.HasOneInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
			hasPk = true
		}
	}
	for _, constraint := range createTableStmt.Constraints {
		switch constraint.Tp {
		case ast.ConstraintPrimaryKey:
			hasPk = true
		default:
			if constraint.Name != "" {
				indexNameMap[constraint.Name] = struct{}{}
			}
		}
	}

	needNotExistsColsName := []string{}
	needExistsColsName := []string{}

	// check drop column
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropColumn) {
		oldColName := spec.OldColumnName.Name.L
		if _, ok := colNameMap[oldColName]; !ok {
			needExistsColsName = append(needExistsColsName, oldColName)
		} else {
			delete(colNameMap, oldColName)
		}
	}

	// check change column
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableChangeColumn) {
		oldColName := spec.OldColumnName.Name.L
		if _, ok := colNameMap[oldColName]; !ok {
			needExistsColsName = append(needExistsColsName, oldColName)
		}
		for _, col := range spec.NewColumns {
			newColName := col.Name.Name.L
			if newColName == oldColName {
				continue
			}
			if _, ok := colNameMap[newColName]; ok {
				needNotExistsColsName = append(needNotExistsColsName, newColName)
			} else {
				if newColName != oldColName {
					delete(colNameMap, oldColName)
					colNameMap[newColName] = struct{}{}
				}
			}
		}
	}

	// check add column
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns) {
		for _, col := range spec.NewColumns {
			colName := col.Name.Name.L
			if _, ok := colNameMap[colName]; ok {
				needNotExistsColsName = append(needNotExistsColsName, colName)
			} else {
				colNameMap[colName] = struct{}{}
				if hasPk && util.HasOneInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
					i.result.Add(driver.RuleLevelError, PrimaryKeyExistMessage)
				} else {
					hasPk = true
				}
			}
		}
	}

	// check alter column
	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAlterColumn) {
		for _, col := range spec.NewColumns {
			colName := col.Name.Name.L
			if _, ok := colNameMap[colName]; !ok {
				needExistsColsName = append(needExistsColsName, colName)
			}
		}
	}

	needNotExistsIndexesName := []string{}
	needExistsIndexesName := []string{}
	needExistsKeyColsName := []string{}

	if len(util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropPrimaryKey)) > 0 && !hasPk {
		// primary key not exist, can not drop primary key
		i.result.Add(driver.RuleLevelError, PrimaryKeyNotExistMessage)
	}

	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropIndex) {
		indexName := strings.ToLower(spec.Name)
		if _, ok := indexNameMap[indexName]; !ok {
			needExistsIndexesName = append(needExistsIndexesName, indexName)
		}
	}

	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableRenameIndex) {
		oldIndexName := spec.FromKey.String()
		newIndexName := spec.ToKey.String()
		_, oldIndexExist := indexNameMap[oldIndexName]
		if !oldIndexExist {
			needExistsIndexesName = append(needExistsIndexesName, oldIndexName)
		}
		_, newIndexExist := indexNameMap[newIndexName]
		if newIndexExist {
			needNotExistsIndexesName = append(needNotExistsIndexesName, newIndexName)
		}
		if oldIndexExist && !newIndexExist {
			delete(indexNameMap, oldIndexName)
			indexNameMap[newIndexName] = struct{}{}
		}
	}

	for _, spec := range util.GetAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
		switch spec.Constraint.Tp {
		case ast.ConstraintPrimaryKey:
			if hasPk {
				// primary key has exist, can not add primary key
				i.result.Add(driver.RuleLevelError, PrimaryKeyExistMessage)
			} else {
				hasPk = true
			}
			names := []string{}
			for _, col := range spec.Constraint.Keys {
				colName := col.Column.Name.L
				names = append(names, colName)
				if _, ok := colNameMap[colName]; !ok {
					needExistsKeyColsName = append(needExistsKeyColsName, colName)
				}
			}
			duplicateColumn := utils.GetDuplicate(names)
			if len(duplicateColumn) > 0 {
				i.result.Add(driver.RuleLevelError, DuplicatePrimaryKeyedColumnMessage,
					strings.Join(duplicateColumn, ","))
			}
		case ast.ConstraintUniq, ast.ConstraintIndex, ast.ConstraintFulltext:
			indexName := strings.ToLower(spec.Constraint.Name)
			if indexName != "" {
				if _, ok := indexNameMap[indexName]; ok {
					needNotExistsIndexesName = append(needNotExistsIndexesName, indexName)
				} else {
					indexNameMap[indexName] = struct{}{}
				}
			} else {
				indexName = "(匿名)"
			}
			names := []string{}
			for _, col := range spec.Constraint.Keys {
				colName := col.Column.Name.L
				names = append(names, colName)
				if _, ok := colNameMap[colName]; !ok {
					needExistsKeyColsName = append(needExistsKeyColsName, colName)
				}
			}
			duplicateColumn := utils.GetDuplicate(names)
			if len(duplicateColumn) > 0 {
				i.result.Add(driver.RuleLevelError, DuplicateIndexedColumnMessage, indexName,
					strings.Join(duplicateColumn, ","))
			}
		}
	}

	if len(needExistsColsName) > 0 {
		i.result.Add(driver.RuleLevelError, ColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsColsName), ","))
	}
	if len(needNotExistsColsName) > 0 {
		i.result.Add(driver.RuleLevelError, ColumnExistMessage,
			strings.Join(utils.RemoveDuplicate(needNotExistsColsName), ","))
	}
	if len(needExistsIndexesName) > 0 {
		i.result.Add(driver.RuleLevelError, IndexNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsIndexesName), ","))
	}
	if len(needNotExistsIndexesName) > 0 {
		i.result.Add(driver.RuleLevelError, IndexExistMessage,
			strings.Join(utils.RemoveDuplicate(needNotExistsIndexesName), ","))
	}
	if len(needExistsKeyColsName) > 0 {
		i.result.Add(driver.RuleLevelError, KeyedColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsKeyColsName), ","))
	}
	return nil
}

/*
------------------------------------------------------------------
drop table ...
------------------------------------------------------------------
1. schema must exist;
2. table must exist if SQL has not "IF EXISTS".
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidDropTable(stmt *ast.DropTableStmt) error {
	if stmt.IfExists {
		return nil
	}
	needExistsSchemasName := []string{}
	needExistsTablesName := []string{}
	for _, table := range stmt.Tables {
		schemaName := i.Ctx.GetSchemaName(table)
		schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if !schemaExist {
			needExistsSchemasName = append(needExistsSchemasName, schemaName)
		} else {
			tableExist, err := i.Ctx.IsTableExist(table)
			if err != nil {
				return err
			}
			if !tableExist {
				needExistsTablesName = append(needExistsTablesName, i.getTableName(table))
			}
		}
	}
	if len(needExistsSchemasName) > 0 {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		i.result.Add(driver.RuleLevelError, TableNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsTablesName), ","))
	}
	return nil
}

/*
------------------------------------------------------------------
use database ...
------------------------------------------------------------------
1. schema must exist.
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidUse(stmt *ast.UseStmt) error {
	schemaExist, err := i.Ctx.IsSchemaExist(stmt.DBName)
	if err != nil {
		return err
	}
	if !schemaExist {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage, stmt.DBName)
	}
	return nil
}

/*
------------------------------------------------------------------
create database ...
------------------------------------------------------------------
1. schema can't exist if SQL has not "IF NOT EXISTS".
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidCreateDatabase(stmt *ast.CreateDatabaseStmt) error {
	if stmt.IfNotExists {
		return nil
	}
	schemaName := stmt.Name
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if schemaExist {
		i.result.Add(driver.RuleLevelError, SchemaExistMessage, schemaName)
	}
	return nil
}

/*
------------------------------------------------------------------
drop database ...
------------------------------------------------------------------
1. schema must exist if SQL has not "IF EXISTS".
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidDropDatabase(stmt *ast.DropDatabaseStmt) error {
	if stmt.IfExists {
		return nil
	}
	schemaName := stmt.Name
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage, schemaName)
	}
	return nil
}

/*
------------------------------------------------------------------
create index ...
------------------------------------------------------------------
1. schema must exist;
2. table must exist;
3. index name can't be duplicated;
4. index column name can't be duplicated.
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidCreateIndex(stmt *ast.CreateIndexStmt) error {
	schemaName := i.Ctx.GetSchemaName(stmt.Table)
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.Ctx.GetCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		i.result.Add(driver.RuleLevelError, TableNotExistMessage,
			i.getTableName(stmt.Table))
		return nil
	}
	colNameMap := map[string]struct{}{}
	indexNameMap := map[string]struct{}{}
	for _, col := range createTableStmt.Cols {
		colNameMap[col.Name.Name.L] = struct{}{}
	}
	for _, constraint := range createTableStmt.Constraints {
		if constraint.Name != "" {
			indexNameMap[constraint.Name] = struct{}{}
		}
	}
	if _, ok := indexNameMap[stmt.IndexName]; ok {
		i.result.Add(driver.RuleLevelError, IndexExistMessage, stmt.IndexName)
	}
	keyColsName := []string{}
	keyColNeedExist := []string{}
	for _, col := range stmt.IndexPartSpecifications {
		colName := col.Column.Name.L
		keyColsName = append(keyColsName, colName)
		if _, ok := colNameMap[col.Column.Name.L]; !ok {
			keyColNeedExist = append(keyColNeedExist, colName)
		}
	}
	duplicateName := utils.GetDuplicate(keyColsName)
	if len(duplicateName) > 0 {
		i.result.Add(driver.RuleLevelError, DuplicateIndexedColumnMessage, stmt.IndexName,
			strings.Join(duplicateName, ","))
	}

	if len(keyColNeedExist) > 0 {
		i.result.Add(driver.RuleLevelError, KeyedColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(keyColNeedExist), ","))
	}
	return nil
}

/*
------------------------------------------------------------------
drop index ...
------------------------------------------------------------------
1. schema must exist;
2. table must exist;
3. index name must exist if SQL has not "IF EXISTS".
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidDropIndex(stmt *ast.DropIndexStmt) error {
	if stmt.IfExists {
		return nil
	}
	schemaName := i.Ctx.GetSchemaName(stmt.Table)
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.Ctx.GetCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		i.result.Add(driver.RuleLevelError, TableNotExistMessage,
			i.getTableName(stmt.Table))
		return nil
	}
	indexNameMap := map[string]struct{}{}
	for _, constraint := range createTableStmt.Constraints {
		if constraint.Name != "" {
			indexNameMap[constraint.Name] = struct{}{}
		}
	}
	if _, ok := indexNameMap[stmt.IndexName]; !ok {
		i.result.Add(driver.RuleLevelError, IndexNotExistMessage, stmt.IndexName)
	}
	return nil
}

/*
------------------------------------------------------------------
insert into ... values ...
------------------------------------------------------------------
1. schema must exist;
2. table must exist;
3. column must exist;
4. value length must match column length.
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidInsert(stmt *ast.InsertStmt) error {
	tables := util.GetTables(stmt.Table.TableRefs)
	table := tables[0]
	schemaName := i.Ctx.GetSchemaName(table)
	schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.Ctx.GetCreateTableStmt(table)
	if err != nil {
		return err
	}
	if !tableExist {
		i.result.Add(driver.RuleLevelError, TableNotExistMessage,
			i.getTableName(table))
		return nil
	}

	colNameMap := map[string]struct{}{}
	for _, col := range createTableStmt.Cols {
		colNameMap[col.Name.Name.L] = struct{}{}
	}

	insertColsName := []string{}
	if stmt.Columns != nil {
		for _, col := range stmt.Columns {
			insertColsName = append(insertColsName, col.Name.L)
		}

	} else if stmt.Setlist != nil {
		for _, set := range stmt.Setlist {
			insertColsName = append(insertColsName, set.Column.Name.L)
		}
	} else {
		for _, col := range createTableStmt.Cols {
			insertColsName = append(insertColsName, col.Name.Name.L)
		}
	}
	if d := utils.GetDuplicate(insertColsName); len(d) > 0 {
		i.result.Add(driver.RuleLevelError, DuplicateColumnsMessage, strings.Join(d, ","))
	}

	needExistColsName := []string{}
	for _, colName := range insertColsName {
		if _, ok := colNameMap[colName]; !ok {
			needExistColsName = append(needExistColsName, colName)
		}
	}
	if len(needExistColsName) > 0 {
		i.result.Add(driver.RuleLevelError, ColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistColsName), ","))
	}

	if stmt.Lists != nil {
		for _, list := range stmt.Lists {
			if len(list) != len(insertColsName) {
				i.result.Add(driver.RuleLevelError, ColumnsValuesNotMatchMessage)
				break
			}
		}
	}
	return nil
}

/*
------------------------------------------------------------------
update ... set  ... where ...
------------------------------------------------------------------
1. schema must exist;
2. table must exist;
3. field column ("set column = ...") must exist;
4. where column ("where column = ...") must exist.
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidUpdate(stmt *ast.UpdateStmt) error {
	tables := []*ast.TableName{}
	tableAlias := map[*ast.TableName]string{}
	tableSources := util.GetTableSources(stmt.TableRefs.TableRefs)
	var hasSelectStmtTableSource bool
	for _, tableSource := range tableSources {
		switch source := tableSource.Source.(type) {
		case *ast.TableName:
			table := source
			tables = append(tables, table)
			alias := tableSource.AsName.String()
			if alias != "" {
				tableAlias[table] = alias
			}
		case *ast.SelectStmt:
			hasSelectStmtTableSource = true
		case *ast.UnionStmt:
			continue
		}
	}
	needExistsSchemasName := []string{}
	needExistsTablesName := []string{}
	for _, table := range tables {
		schemaName := i.Ctx.GetSchemaName(table)
		schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if !schemaExist {
			needExistsSchemasName = append(needExistsSchemasName, schemaName)
		} else {
			tableExist, err := i.Ctx.IsTableExist(table)
			if err != nil {
				return err
			}
			if !tableExist {
				needExistsTablesName = append(needExistsTablesName, i.getTableName(table))
			}
		}
	}
	if len(needExistsSchemasName) > 0 {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		i.result.Add(driver.RuleLevelError, TableNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsTablesName), ","))
	}

	if len(needExistsSchemasName) > 0 || len(needExistsTablesName) > 0 {
		return nil
	}

	tc := util.NewTableChecker()
	for _, table := range tables {
		schemaName := table.Schema.String()
		if schemaName == "" {
			schemaName = i.Ctx.CurrentSchema()
		}
		tableName := table.Name.String()
		if alias, ok := tableAlias[table]; ok {
			tableName = alias
		}
		createStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
		if err != nil || !exist {
			return err
		}
		tc.Add(schemaName, tableName, createStmt)
	}

	// https://github.com/actiontech/sqle/issues/708
	// If the updated table contains subquery, do not check fields in set and where clause.
	if hasSelectStmtTableSource {
		return nil
	}

	needExistColsName := []string{}
	ambiguousColsName := []string{}
	for _, list := range stmt.List {
		col := list.Column
		colExists, colIsAmbiguous := tc.CheckColumnByName(col)
		if colIsAmbiguous {
			ambiguousColsName = append(ambiguousColsName, col.String())
			continue
		}
		if !colExists {
			needExistColsName = append(needExistColsName, col.String())
		}
	}

	util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
		switch x := expr.(type) {
		case *ast.ColumnNameExpr:
			col := x.Name
			colExists, colIsAmbiguous := tc.CheckColumnByName(col)
			if colIsAmbiguous {
				ambiguousColsName = append(ambiguousColsName, col.String())
			}
			if !colExists {
				needExistColsName = append(needExistColsName, col.String())
			}
		}
		return false
	}, stmt.Where)

	if len(needExistColsName) > 0 {
		i.result.Add(driver.RuleLevelError, ColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistColsName), ","))
	}

	if len(ambiguousColsName) > 0 {
		i.result.Add(driver.RuleLevelError, ColumnIsAmbiguousMessage,
			strings.Join(utils.RemoveDuplicate(ambiguousColsName), ","))
	}
	return nil
}

/*
------------------------------------------------------------------
delete from ... where ...
------------------------------------------------------------------
1. schema must exist;
2. table must exist;
3. where column ("where column = ...") must exist.
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidDelete(stmt *ast.DeleteStmt) error {
	tables := util.GetTables(stmt.TableRefs.TableRefs)
	needExistsSchemasName := []string{}
	needExistsTablesName := []string{}
	for _, table := range tables {
		schemaName := i.Ctx.GetSchemaName(table)
		schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if !schemaExist {
			needExistsSchemasName = append(needExistsSchemasName, schemaName)
		} else {
			tableExist, err := i.Ctx.IsTableExist(table)
			if err != nil {
				return err
			}
			if !tableExist {
				needExistsTablesName = append(needExistsTablesName, i.getTableName(table))
			}
		}
	}
	if len(needExistsSchemasName) > 0 {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		i.result.Add(driver.RuleLevelError, TableNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsTablesName), ","))
	}
	if len(needExistsSchemasName) > 0 || len(needExistsTablesName) > 0 {
		return nil
	}

	tc := util.NewTableChecker()
	for _, table := range tables {
		schemaName := table.Schema.String()
		if schemaName == "" {
			schemaName = i.Ctx.CurrentSchema()
		}
		tableName := table.Name.String()
		createStmt, exist, err := i.Ctx.GetCreateTableStmt(table)
		if err != nil || !exist {
			return err
		}
		tc.Add(schemaName, tableName, createStmt)
	}

	needExistColsName := []string{}
	ambiguousColsName := []string{}
	util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
		switch x := expr.(type) {
		case *ast.ColumnNameExpr:
			col := x.Name
			colExists, colIsAmbiguous := tc.CheckColumnByName(col)
			if colIsAmbiguous {
				ambiguousColsName = append(ambiguousColsName, col.String())
			}
			if !colExists {
				needExistColsName = append(needExistColsName, col.String())
			}
		}
		return false
	}, stmt.Where)

	if len(needExistColsName) > 0 {
		i.result.Add(driver.RuleLevelError, ColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistColsName), ","))
	}

	if len(ambiguousColsName) > 0 {
		i.result.Add(driver.RuleLevelError, ColumnIsAmbiguousMessage,
			strings.Join(utils.RemoveDuplicate(ambiguousColsName), ","))
	}
	return nil
}

/*
------------------------------------------------------------------
select ... from ...
------------------------------------------------------------------
1. schema must exist;
2. table must exist.
------------------------------------------------------------------
*/
func (i *MysqlDriverImpl) checkInvalidSelect(stmt *ast.SelectStmt) error {
	if stmt.From == nil {
		return nil
	}
	tables := []*ast.TableName{}
	tableSources := util.GetTableSources(stmt.From.TableRefs)
	// not select from table statement
	if len(tableSources) < 1 {
		return nil
	}
	for _, tableSource := range tableSources {
		switch source := tableSource.Source.(type) {
		case *ast.TableName:
			table := source
			tables = append(tables, table)
		case *ast.SelectStmt, *ast.UnionStmt:
			continue
		}
	}
	needExistsSchemasName := []string{}
	needExistsTablesName := []string{}
	for _, table := range tables {
		schemaName := i.Ctx.GetSchemaName(table)
		schemaExist, err := i.Ctx.IsSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if !schemaExist {
			needExistsSchemasName = append(needExistsSchemasName, schemaName)
		} else {
			tableExist, err := i.Ctx.IsTableExist(table)
			if err != nil {
				return err
			}
			if !tableExist {
				needExistsTablesName = append(needExistsTablesName, i.getTableName(table))
			}
		}
	}
	if len(needExistsSchemasName) > 0 {
		i.result.Add(driver.RuleLevelError, SchemaNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		i.result.Add(driver.RuleLevelError, TableNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsTablesName), ","))
	}
	return nil
}

// checkUnparsedStmt might add more check in future.
func (i *MysqlDriverImpl) checkUnparsedStmt(stmt *ast.UnparsedStmt) error {
	i.result.Add(driver.RuleLevelWarn, "语法错误或者解析器不支持，请人工确认SQL正确性")
	return nil
}
