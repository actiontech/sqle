package mysql

import (
	"fmt"
	"strings"

	"actiontech.cloud/sqle/sqle/sqle/model"
	"actiontech.cloud/sqle/sqle/sqle/utils"

	"github.com/pingcap/parser/ast"
)

func (i *Inspect) Advise(rules []model.Rule, wl []model.SqlWhitelist) error {
	i.Logger().Info("start advise sql")
	err := i.advise(rules, wl)
	if err != nil {
		i.Logger().Error("advise sql failed")
	} else {
		i.Logger().Info("advise sql finish")
	}
	return err
}

func (i *Inspect) advise(rules []model.Rule, wl []model.SqlWhitelist) error {
	for _, commitSql := range i.Task.ExecuteSQLs {
		currentSql := commitSql
		if err := i.Add(&currentSql.BaseSQL, func(node ast.Node) (err error) {
			lowerCaseTableNames, err := i.getSystemVariable(SysVarLowerCaseTableNames)
			if err != nil {
				return err
			}

			sqlFP, err := Fingerprint(currentSql.BaseSQL.Content, lowerCaseTableNames == "0")
			if err != nil {
				return err
			}

			var whitelistMatch bool
			for _, sqlInWL := range wl {
				if sqlInWL.MatchType == model.SQLWhitelistFPMatch {
					whitelistFP, err := Fingerprint(sqlInWL.Value, lowerCaseTableNames == "0")
					if err != nil {
						return err
					}
					if whitelistFP == sqlFP {
						whitelistMatch = true
						break
					}
				} else {
					if sqlInWL.CapitalizedValue == strings.ToUpper(currentSql.BaseSQL.Content) {
						whitelistMatch = true
						break
					}
				}
			}
			if whitelistMatch {
				var results InspectResults
				results.add(model.RuleLevelNormal, "白名单")
				currentSql.AuditStatus = model.SQLAuditStatusFinished
				currentSql.AuditLevel = results.level()
				currentSql.AuditResult = results.message()
			} else {
				results, err := i.CheckInvalid(node)
				if err != nil {
					return err
				}
				if results.level() == model.RuleLevelError {
					i.HasInvalidSql = true
					i.Logger().Warnf("sql %s invalid, %s", node.Text(), results.message())
				}
				i.Results = results
				if rules != nil {
					for _, rule := range rules {
						i.currentRule = rule
						handler, ok := RuleHandlerMap[rule.Name]
						if !ok || handler.Func == nil {
							continue
						}
						err := handler.Func(rule, i, node)
						if err != nil {
							return err
						}
					}
				}
				currentSql.AuditStatus = model.SQLAuditStatusFinished
				currentSql.AuditLevel = i.Results.level()
				currentSql.AuditResult = i.Results.message()
				// clean up results
				i.Results = newInspectResults()

				// print osc
				oscCommandLine, err := i.generateOSCCommandLine(node)
				if err != nil {
					return err
				}
				if oscCommandLine != "" {
					results := newInspectResults()
					if currentSql.AuditResult != "" {
						results.add(currentSql.AuditLevel, currentSql.AuditResult)
					}
					results.add(model.RuleLevelNotice, fmt.Sprintf("[osc]%s", oscCommandLine))
					currentSql.AuditLevel = results.level()
					currentSql.AuditResult = results.message()
				}
			}

			currentSql.AuditFingerprint = utils.Md5String(string(append([]byte(currentSql.AuditResult), []byte(sqlFP)...)))

			i.Logger().Infof("sql=%s, level=%s, result=%s", currentSql.Content, currentSql.AuditLevel, currentSql.AuditResult)
			return nil
		}); err != nil {
			i.Logger().Error("add commit sql to task failed")
			return err
		}
	}
	return i.Do()
}

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

func (i *Inspect) CheckInvalid(node ast.Node) (*InspectResults, error) {
	results := newInspectResults()
	var err error
	switch stmt := node.(type) {
	case *ast.UseStmt:
		err = i.checkInvalidUse(stmt, results)
	case *ast.CreateTableStmt:
		err = i.checkInvalidCreateTable(stmt, results)
	case *ast.AlterTableStmt:
		err = i.checkInvalidAlterTable(stmt, results)
	case *ast.DropTableStmt:
		err = i.checkInvalidDropTable(stmt, results)
	case *ast.CreateDatabaseStmt:
		err = i.checkInvalidCreateDatabase(stmt, results)
	case *ast.DropDatabaseStmt:
		err = i.checkInvalidDropDatabase(stmt, results)
	case *ast.CreateIndexStmt:
		err = i.checkInvalidCreateIndex(stmt, results)
	case *ast.DropIndexStmt:
		err = i.checkInvalidDropIndex(stmt, results)
	case *ast.InsertStmt:
		err = i.checkInvalidInsert(stmt, results)
	case *ast.UpdateStmt:
		err = i.checkInvalidUpdate(stmt, results)
	case *ast.DeleteStmt:
		err = i.checkInvalidDelete(stmt, results)
	case *ast.SelectStmt:
		err = i.checkInvalidSelect(stmt, results)
	case *ast.UnparsedStmt:
		err = i.checkUnparsedStmt(stmt, results)
	}
	return results, err
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
func (i *Inspect) checkInvalidCreateTable(stmt *ast.CreateTableStmt, results *InspectResults) error {
	schemaName := i.getSchemaName(stmt.Table)
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		results.add(model.RuleLevelError, SchemaNotExistMessage, schemaName)
	} else {
		tableExist, err := i.isTableExist(stmt.Table)
		if err != nil {
			return err
		}
		if tableExist && !stmt.IfNotExists {
			results.add(model.RuleLevelError, TableExistMessage,
				i.getTableName(stmt.Table))
		}
		if stmt.ReferTable != nil {
			referTableExist, err := i.isTableExist(stmt.ReferTable)
			if err != nil {
				return err
			}
			if !referTableExist {
				results.add(model.RuleLevelError, TableNotExistMessage,
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
		if HasOneInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
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
				results.add(model.RuleLevelError, DuplicatePrimaryKeyedColumnMessage,
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
				results.add(model.RuleLevelError, DuplicateIndexedColumnMessage, constraintName,
					strings.Join(duplicateName, ","))
			}
		}
	}
	if d := utils.GetDuplicate(colsName); len(d) > 0 {
		results.add(model.RuleLevelError, DuplicateColumnsMessage,
			strings.Join(d, ","))
	}

	if d := utils.GetDuplicate(indexesName); len(d) > 0 {
		results.add(model.RuleLevelError, DuplicateIndexesMessage,
			strings.Join(d, ","))
	}

	if pkCounter > 1 {
		results.add(model.RuleLevelError, MultiPrimaryKeyMessage)
	}
	notExistKeyColsName := []string{}
	for _, colName := range keyColsName {
		if _, ok := colsNameMap[colName]; !ok {
			notExistKeyColsName = append(notExistKeyColsName, colName)
		}
	}
	if len(notExistKeyColsName) > 0 {
		results.add(model.RuleLevelError, KeyedColumnNotExistMessage,
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
func (i *Inspect) checkInvalidAlterTable(stmt *ast.AlterTableStmt, results *InspectResults) error {
	schemaName := i.getSchemaName(stmt.Table)
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		results.add(model.RuleLevelError, SchemaNotExistMessage, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.getCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		results.add(model.RuleLevelError, TableNotExistMessage,
			i.getTableName(stmt.Table))
		return nil
	}

	hasPk := false
	colNameMap := map[string]struct{}{}
	indexNameMap := map[string]struct{}{}
	for _, col := range createTableStmt.Cols {
		colNameMap[col.Name.Name.L] = struct{}{}
		if HasOneInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropColumn) {
		oldColName := spec.OldColumnName.Name.L
		if _, ok := colNameMap[oldColName]; !ok {
			needExistsColsName = append(needExistsColsName, oldColName)
		} else {
			delete(colNameMap, oldColName)
		}
	}

	// check change column
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableChangeColumn) {
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
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns) {
		for _, col := range spec.NewColumns {
			colName := col.Name.Name.L
			if _, ok := colNameMap[colName]; ok {
				needNotExistsColsName = append(needNotExistsColsName, colName)
			} else {
				colNameMap[colName] = struct{}{}
				if hasPk && HasOneInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
					results.add(model.RuleLevelError, PrimaryKeyExistMessage)
				} else {
					hasPk = true
				}
			}
		}
	}

	// check alter column
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAlterColumn) {
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

	if len(getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropPrimaryKey)) > 0 && !hasPk {
		// primary key not exist, can not drop primary key
		results.add(model.RuleLevelError, PrimaryKeyNotExistMessage)
	}

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropIndex) {
		indexName := strings.ToLower(spec.Name)
		if _, ok := indexNameMap[indexName]; !ok {
			needExistsIndexesName = append(needExistsIndexesName, indexName)
		}
	}

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableRenameIndex) {
		oldIndexName := spec.FromKey.String()
		newIndexName := spec.ToKey.String()
		_, oldIndexExist := indexNameMap[oldIndexName]
		if !oldIndexExist {
			needExistsIndexesName = append(needExistsIndexesName, oldIndexName)
		}
		_, newIndexExist := indexNameMap[newIndexName]
		if newIndexExist {
			needNotExistsIndexesName = append(needNotExistsIndexesName)
		}
		if oldIndexExist && !newIndexExist {
			delete(indexNameMap, oldIndexName)
			indexNameMap[newIndexName] = struct{}{}
		}
	}

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
		switch spec.Constraint.Tp {
		case ast.ConstraintPrimaryKey:
			if hasPk {
				// primary key has exist, can not add primary key
				results.add(model.RuleLevelError, PrimaryKeyExistMessage)
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
				results.add(model.RuleLevelError, DuplicatePrimaryKeyedColumnMessage,
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
				results.add(model.RuleLevelError, DuplicateIndexedColumnMessage, indexName,
					strings.Join(duplicateColumn, ","))
			}
		}
	}

	if len(needExistsColsName) > 0 {
		results.add(model.RuleLevelError, ColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsColsName), ","))
	}
	if len(needNotExistsColsName) > 0 {
		results.add(model.RuleLevelError, ColumnExistMessage,
			strings.Join(utils.RemoveDuplicate(needNotExistsColsName), ","))
	}
	if len(needExistsIndexesName) > 0 {
		results.add(model.RuleLevelError, IndexNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsIndexesName), ","))
	}
	if len(needNotExistsIndexesName) > 0 {
		results.add(model.RuleLevelError, IndexExistMessage,
			strings.Join(utils.RemoveDuplicate(needNotExistsIndexesName), ","))
	}
	if len(needExistsKeyColsName) > 0 {
		results.add(model.RuleLevelError, KeyedColumnNotExistMessage,
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
func (i *Inspect) checkInvalidDropTable(stmt *ast.DropTableStmt, results *InspectResults) error {
	if stmt.IfExists {
		return nil
	}
	needExistsSchemasName := []string{}
	needExistsTablesName := []string{}
	for _, table := range stmt.Tables {
		schemaName := i.getSchemaName(table)
		schemaExist, err := i.isSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if !schemaExist {
			needExistsSchemasName = append(needExistsSchemasName, schemaName)
		} else {
			tableExist, err := i.isTableExist(table)
			if err != nil {
				return err
			}
			if !tableExist {
				needExistsTablesName = append(needExistsTablesName, i.getTableName(table))
			}
		}
	}
	if len(needExistsSchemasName) > 0 {
		results.add(model.RuleLevelError, SchemaNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		results.add(model.RuleLevelError, TableNotExistMessage,
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
func (i *Inspect) checkInvalidUse(stmt *ast.UseStmt, results *InspectResults) error {
	schemaExist, err := i.isSchemaExist(stmt.DBName)
	if err != nil {
		return err
	}
	if !schemaExist {
		results.add(model.RuleLevelError, SchemaNotExistMessage, stmt.DBName)
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
func (i *Inspect) checkInvalidCreateDatabase(stmt *ast.CreateDatabaseStmt,
	results *InspectResults) error {
	if stmt.IfNotExists {
		return nil
	}
	schemaName := stmt.Name
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if schemaExist {
		results.add(model.RuleLevelError, SchemaExistMessage, schemaName)
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
func (i *Inspect) checkInvalidDropDatabase(stmt *ast.DropDatabaseStmt,
	results *InspectResults) error {
	if stmt.IfExists {
		return nil
	}
	schemaName := stmt.Name
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		results.add(model.RuleLevelError, SchemaNotExistMessage, schemaName)
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
func (i *Inspect) checkInvalidCreateIndex(stmt *ast.CreateIndexStmt,
	results *InspectResults) error {
	schemaName := i.getSchemaName(stmt.Table)
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		results.add(model.RuleLevelError, SchemaNotExistMessage, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.getCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		results.add(model.RuleLevelError, TableNotExistMessage,
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
		results.add(model.RuleLevelError, IndexExistMessage, stmt.IndexName)
	}
	keyColsName := []string{}
	keyColNeedExist := []string{}
	for _, col := range stmt.IndexColNames {
		colName := col.Column.Name.L
		keyColsName = append(keyColsName, colName)
		if _, ok := colNameMap[col.Column.Name.L]; !ok {
			keyColNeedExist = append(keyColNeedExist, colName)
		}
	}
	duplicateName := utils.GetDuplicate(keyColsName)
	if len(duplicateName) > 0 {
		results.add(model.RuleLevelError, DuplicateIndexedColumnMessage, stmt.IndexName,
			strings.Join(duplicateName, ","))
	}

	if len(keyColNeedExist) > 0 {
		results.add(model.RuleLevelError, KeyedColumnNotExistMessage,
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
func (i *Inspect) checkInvalidDropIndex(stmt *ast.DropIndexStmt,
	results *InspectResults) error {
	if stmt.IfExists {
		return nil
	}
	schemaName := i.getSchemaName(stmt.Table)
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		results.add(model.RuleLevelError, SchemaNotExistMessage, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.getCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		results.add(model.RuleLevelError, TableNotExistMessage,
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
		results.add(model.RuleLevelError, IndexNotExistMessage, stmt.IndexName)
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
func (i *Inspect) checkInvalidInsert(stmt *ast.InsertStmt, results *InspectResults) error {
	tables := getTables(stmt.Table.TableRefs)
	table := tables[0]
	schemaName := i.getSchemaName(table)
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		results.add(model.RuleLevelError, SchemaNotExistMessage, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.getCreateTableStmt(table)
	if err != nil {
		return err
	}
	if !tableExist {
		results.add(model.RuleLevelError, TableNotExistMessage,
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
		results.add(model.RuleLevelError, DuplicateColumnsMessage, strings.Join(d, ","))
	}

	needExistColsName := []string{}
	for _, colName := range insertColsName {
		if _, ok := colNameMap[colName]; !ok {
			needExistColsName = append(needExistColsName, colName)
		}
	}
	if len(needExistColsName) > 0 {
		results.add(model.RuleLevelError, ColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistColsName), ","))
	}

	if stmt.Lists != nil {
		for _, list := range stmt.Lists {
			if len(list) != len(insertColsName) {
				results.add(model.RuleLevelError, ColumnsValuesNotMatchMessage)
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
func (i *Inspect) checkInvalidUpdate(stmt *ast.UpdateStmt, results *InspectResults) error {
	tables := []*ast.TableName{}
	tableAlias := map[*ast.TableName]string{}
	tableSources := getTableSources(stmt.TableRefs.TableRefs)
	for _, tableSource := range tableSources {
		switch source := tableSource.Source.(type) {
		case *ast.TableName:
			table := source
			tables = append(tables, table)
			alias := tableSource.AsName.String()
			if alias != "" {
				tableAlias[table] = alias
			}
		case *ast.SelectStmt, *ast.UnionStmt:
			continue
		}
	}
	needExistsSchemasName := []string{}
	needExistsTablesName := []string{}
	for _, table := range tables {
		schemaName := i.getSchemaName(table)
		schemaExist, err := i.isSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if !schemaExist {
			needExistsSchemasName = append(needExistsSchemasName, schemaName)
		} else {
			tableExist, err := i.isTableExist(table)
			if err != nil {
				return err
			}
			if !tableExist {
				needExistsTablesName = append(needExistsTablesName, i.getTableName(table))
			}
		}
	}
	if len(needExistsSchemasName) > 0 {
		results.add(model.RuleLevelError, SchemaNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		results.add(model.RuleLevelError, TableNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsTablesName), ","))
	}

	if len(needExistsSchemasName) > 0 || len(needExistsTablesName) > 0 {
		return nil
	}

	tc := newTableChecker()
	for _, table := range tables {
		schemaName := table.Schema.String()
		if schemaName == "" {
			schemaName = i.Ctx.currentSchema
		}
		tableName := table.Name.String()
		if alias, ok := tableAlias[table]; ok {
			tableName = alias
		}
		createStmt, exist, err := i.getCreateTableStmt(table)
		if err != nil || !exist {
			return err
		}
		tc.add(schemaName, tableName, createStmt)
	}

	needExistColsName := []string{}
	ambiguousColsName := []string{}
	for _, list := range stmt.List {
		col := list.Column
		colExists, colIsAmbiguous := tc.checkColumnByName(col)
		if colIsAmbiguous {
			ambiguousColsName = append(ambiguousColsName, col.String())
			continue
		}
		if !colExists {
			needExistColsName = append(needExistColsName, col.String())
		}
	}

	scanWhereStmt(func(expr ast.ExprNode) (skip bool) {
		switch x := expr.(type) {
		case *ast.ColumnNameExpr:
			col := x.Name
			colExists, colIsAmbiguous := tc.checkColumnByName(col)
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
		results.add(model.RuleLevelError, ColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistColsName), ","))
	}

	if len(ambiguousColsName) > 0 {
		results.add(model.RuleLevelError, ColumnIsAmbiguousMessage,
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
func (i *Inspect) checkInvalidDelete(stmt *ast.DeleteStmt, results *InspectResults) error {
	tables := getTables(stmt.TableRefs.TableRefs)
	needExistsSchemasName := []string{}
	needExistsTablesName := []string{}
	for _, table := range tables {
		schemaName := i.getSchemaName(table)
		schemaExist, err := i.isSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if !schemaExist {
			needExistsSchemasName = append(needExistsSchemasName, schemaName)
		} else {
			tableExist, err := i.isTableExist(table)
			if err != nil {
				return err
			}
			if !tableExist {
				needExistsTablesName = append(needExistsTablesName, i.getTableName(table))
			}
		}
	}
	if len(needExistsSchemasName) > 0 {
		results.add(model.RuleLevelError, SchemaNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		results.add(model.RuleLevelError, TableNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsTablesName), ","))
	}
	if len(needExistsSchemasName) > 0 || len(needExistsTablesName) > 0 {
		return nil
	}

	tc := newTableChecker()
	for _, table := range tables {
		schemaName := table.Schema.String()
		if schemaName == "" {
			schemaName = i.Ctx.currentSchema
		}
		tableName := table.Name.String()
		createStmt, exist, err := i.getCreateTableStmt(table)
		if err != nil || !exist {
			return err
		}
		tc.add(schemaName, tableName, createStmt)
	}

	needExistColsName := []string{}
	ambiguousColsName := []string{}
	scanWhereStmt(func(expr ast.ExprNode) (skip bool) {
		switch x := expr.(type) {
		case *ast.ColumnNameExpr:
			col := x.Name
			colExists, colIsAmbiguous := tc.checkColumnByName(col)
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
		results.add(model.RuleLevelError, ColumnNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistColsName), ","))
	}

	if len(ambiguousColsName) > 0 {
		results.add(model.RuleLevelError, ColumnIsAmbiguousMessage,
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
func (i *Inspect) checkInvalidSelect(stmt *ast.SelectStmt, results *InspectResults) error {
	if stmt.From == nil {
		return nil
	}
	tables := []*ast.TableName{}
	tableSources := getTableSources(stmt.From.TableRefs)
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
		schemaName := i.getSchemaName(table)
		schemaExist, err := i.isSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if !schemaExist {
			needExistsSchemasName = append(needExistsSchemasName, schemaName)
		} else {
			tableExist, err := i.isTableExist(table)
			if err != nil {
				return err
			}
			if !tableExist {
				needExistsTablesName = append(needExistsTablesName, i.getTableName(table))
			}
		}
	}
	if len(needExistsSchemasName) > 0 {
		results.add(model.RuleLevelError, SchemaNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		results.add(model.RuleLevelError, TableNotExistMessage,
			strings.Join(utils.RemoveDuplicate(needExistsTablesName), ","))
	}
	return nil
}

// checkUnparsedStmt might add more check in future.
func (i *Inspect) checkUnparsedStmt(stmt *ast.UnparsedStmt, results *InspectResults) error {
	results.add(model.RuleLevelError, "语法错误或者解析器不支持")
	return nil
}
