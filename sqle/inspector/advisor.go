package inspector

import (
	"fmt"
	"strings"

	"actiontech.cloud/universe/sqle/v4/sqle/utils"

	"actiontech.cloud/universe/sqle/v4/sqle/model"

	"github.com/pingcap/parser/ast"
)

func (i *Inspect) Advise(rules []model.Rule) error {
	i.Logger().Info("start advise sql")
	err := i.advise(rules, model.GetSqlWhitelistMD5Map())
	if err != nil {
		i.Logger().Error("advise sql failed")
	} else {
		i.Logger().Info("advise sql finish")
	}
	return err
}

func (i *Inspect) advise(rules []model.Rule, sqlWhiltelistMD5Map map[string]struct{}) error {
	err := i.adviseRelateTask(i.RelateTasks)
	if err != nil {
		return err
	}

	for _, commitSql := range i.Task.CommitSqls {
		currentSql := commitSql
		err := i.Add(&currentSql.Sql, func(sql *model.Sql) error {
			if len(sql.Stmts) <= 0 {
				return nil
			}
			if _, ok := sqlWhiltelistMD5Map[utils.Md5String(strings.ToUpper(sql.Content))]; ok {
				currentSql.InspectStatus = model.TASK_ACTION_DONE
				currentSql.InspectLevel = model.RULE_LEVEL_NORMAL
				currentSql.InspectResult = "白名单"
			} else {
				node := sql.Stmts[0]

				var err error
				if currentSql.FingerPrint, err = i.Fingerprint(sql.Content); err != nil {
					i.Logger().Warnf("sql %s generate fingerprint failed, error: %v", sql.Content, err)
				}

				results, err := i.CheckInvalid(node)
				if err != nil {
					return err
				}
				if results.level() == model.RULE_LEVEL_ERROR {
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
				currentSql.InspectStatus = model.TASK_ACTION_DONE
				currentSql.InspectLevel = i.Results.level()
				currentSql.InspectResult = i.Results.message()
				// clean up results
				i.Results = newInspectResults()

				// print osc
				oscCommandLine, err := i.generateOSCCommandLine(sql.Stmts[0])
				if err != nil {
					return err
				}
				if oscCommandLine != "" {
					results := newInspectResults()
					if currentSql.InspectResult != "" {
						results.add(currentSql.InspectLevel, currentSql.InspectResult)
					}
					results.add(model.RULE_LEVEL_NOTICE, fmt.Sprintf("[osc]%s", oscCommandLine))
					currentSql.InspectLevel = results.level()
					currentSql.InspectResult = results.message()
				}
			}

			i.Logger().Infof("sql=%s, level=%s, result=%s",
				currentSql.Content, currentSql.InspectLevel, currentSql.InspectResult)
			return nil
		})
		if err != nil {
			i.Logger().Error("add commit sql to task failed")
			return err
		}
	}
	return i.Do()
}

func (i *Inspect) adviseRelateTask(relateTasks []model.Task) error {
	if relateTasks == nil || len(relateTasks) == 0 {
		return nil
	}
	taskIdList := []string{}
	for _, task := range relateTasks {
		taskIdList = append(taskIdList, fmt.Sprintf("%d", task.ID))
	}
	i.Logger().Infof("relate advise tasks: %s", strings.Join(taskIdList, ", "))
	currentCtx := NewContext(i.Context())
	for _, task := range relateTasks {
		ri := NewInspect(i.Logger(), currentCtx, &task, nil, nil)
		err := ri.advise(nil, nil)
		if err != nil {
			return err
		}
		if ri.SqlInvalid() {
			i.Logger().Warnf("relate tasks failed, because task %d invalid in tasks", task.ID)
			return i.adviseRelateTask(relateTasks[1:])
		}
	}

	i.Logger().Infof("relate tasks success")

	i.Ctx = currentCtx
	return nil
}

var (
	SCHEMA_NOT_EXIST_MSG             = "schema %s 不存在"
	SCHEMA_EXIST_MSG                 = "schema %s 已存在"
	TABLE_NOT_EXIST_MSG              = "表 %s 不存在"
	TABLE_EXIST_MSG                  = "表 %s 已存在"
	COLUMN_NOT_EXIST_MSG             = "字段 %s 不存在"
	COLUMN_EXIST_MSG                 = "字段 %s 已存在"
	COLUMN_IS_AMBIGUOUS              = "字段 %s 指代不明"
	INDEX_NOT_EXIST_MSG              = "索引 %s 不存在"
	INDEX_EXIST_MSG                  = "索引 %s 已存在"
	DUPLICATE_COLUMN_ERROR_MSG       = "字段名 %s 重复"
	DUPLICATE_INDEX_ERROR_MSG        = "索引名 %s 重复"
	PRIMARY_KEY_MULTI_ERROR_MSG      = "主键只能设置一个"
	KEY_COLUMN_NOT_EXIST_MSG         = "索引字段 %s 不存在"
	PRIMARY_KEY_EXIST_MSG            = "已经存在主键，不能再添加"
	PRIMARY_KEY_NOT_EXIST_MSG        = "当前没有主键，不能执行删除"
	NOT_MATCH_VALUES_AND_COLUMNS     = "指定的值列数与字段列数不匹配"
	DUPLICATE_PRIMARY_KEY_COLUMN_MSG = "主键字段 %s 重复"
	DUPLICATE_INDEX_COLUMN_MSG       = "索引 %s 字段 %s重复"
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, schemaName)
	} else {
		tableExist, err := i.isTableExist(stmt.Table)
		if err != nil {
			return err
		}
		if tableExist && !stmt.IfNotExists {
			results.add(model.RULE_LEVEL_ERROR, TABLE_EXIST_MSG,
				i.getTableName(stmt.Table))
		}
		if stmt.ReferTable != nil {
			referTableExist, err := i.isTableExist(stmt.ReferTable)
			if err != nil {
				return err
			}
			if !referTableExist {
				results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
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
			duplicateName := getDuplicate(names)
			if len(duplicateName) > 0 {
				results.add(model.RULE_LEVEL_ERROR, DUPLICATE_PRIMARY_KEY_COLUMN_MSG,
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
			duplicateName := getDuplicate(names)
			if len(duplicateName) > 0 {
				results.add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, constraintName,
					strings.Join(duplicateName, ","))
			}
		}
	}
	if d := getDuplicate(colsName); len(d) > 0 {
		results.add(model.RULE_LEVEL_ERROR, DUPLICATE_COLUMN_ERROR_MSG,
			strings.Join(d, ","))
	}

	if d := getDuplicate(indexesName); len(d) > 0 {
		results.add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_ERROR_MSG,
			strings.Join(d, ","))
	}

	if pkCounter > 1 {
		results.add(model.RULE_LEVEL_ERROR, PRIMARY_KEY_MULTI_ERROR_MSG)
	}
	notExistKeyColsName := []string{}
	for _, colName := range keyColsName {
		if _, ok := colsNameMap[colName]; !ok {
			notExistKeyColsName = append(notExistKeyColsName, colName)
		}
	}
	if len(notExistKeyColsName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(notExistKeyColsName), ","))
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.getCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
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
					results.add(model.RULE_LEVEL_ERROR, PRIMARY_KEY_EXIST_MSG)
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
		results.add(model.RULE_LEVEL_ERROR, PRIMARY_KEY_NOT_EXIST_MSG)
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
				results.add(model.RULE_LEVEL_ERROR, PRIMARY_KEY_EXIST_MSG)
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
			duplicateColumn := getDuplicate(names)
			if len(duplicateColumn) > 0 {
				results.add(model.RULE_LEVEL_ERROR, DUPLICATE_PRIMARY_KEY_COLUMN_MSG,
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
			duplicateColumn := getDuplicate(names)
			if len(duplicateColumn) > 0 {
				results.add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, indexName,
					strings.Join(duplicateColumn, ","))
			}
		}
	}

	if len(needExistsColsName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsColsName), ","))
	}
	if len(needNotExistsColsName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, COLUMN_EXIST_MSG,
			strings.Join(removeDuplicate(needNotExistsColsName), ","))
	}
	if len(needExistsIndexesName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, INDEX_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsIndexesName), ","))
	}
	if len(needNotExistsIndexesName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, INDEX_EXIST_MSG,
			strings.Join(removeDuplicate(needNotExistsIndexesName), ","))
	}
	if len(needExistsKeyColsName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsKeyColsName), ","))
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsTablesName), ","))
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, stmt.DBName)
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_EXIST_MSG, schemaName)
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, schemaName)
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.getCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
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
		results.add(model.RULE_LEVEL_ERROR, INDEX_EXIST_MSG, stmt.IndexName)
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
	duplicateName := getDuplicate(keyColsName)
	if len(duplicateName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, DUPLICATE_INDEX_COLUMN_MSG, stmt.IndexName,
			strings.Join(duplicateName, ","))
	}

	if len(keyColNeedExist) > 0 {
		results.add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(keyColNeedExist), ","))
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.getCreateTableStmt(stmt.Table)
	if err != nil {
		return err
	}
	if !tableExist {
		results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
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
		results.add(model.RULE_LEVEL_ERROR, INDEX_NOT_EXIST_MSG, stmt.IndexName)
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG, schemaName)
		return nil
	}
	createTableStmt, tableExist, err := i.getCreateTableStmt(table)
	if err != nil {
		return err
	}
	if !tableExist {
		results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
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
	if d := getDuplicate(insertColsName); len(d) > 0 {
		results.add(model.RULE_LEVEL_ERROR, DUPLICATE_COLUMN_ERROR_MSG, strings.Join(d, ","))
	}

	needExistColsName := []string{}
	for _, colName := range insertColsName {
		if _, ok := colNameMap[colName]; !ok {
			needExistColsName = append(needExistColsName, colName)
		}
	}
	if len(needExistColsName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistColsName), ","))
	}

	if stmt.Lists != nil {
		for _, list := range stmt.Lists {
			if len(list) != len(insertColsName) {
				results.add(model.RULE_LEVEL_ERROR, NOT_MATCH_VALUES_AND_COLUMNS)
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsTablesName), ","))
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
		results.add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistColsName), ","))
	}

	if len(ambiguousColsName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, COLUMN_IS_AMBIGUOUS,
			strings.Join(removeDuplicate(ambiguousColsName), ","))
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsTablesName), ","))
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
		results.add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistColsName), ","))
	}

	if len(ambiguousColsName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, COLUMN_IS_AMBIGUOUS,
			strings.Join(removeDuplicate(ambiguousColsName), ","))
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
		return fmt.Errorf("failed sql :%v  only support select from table", stmt.Text())
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
		results.add(model.RULE_LEVEL_ERROR, SCHEMA_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsSchemasName), ","))
	}
	if len(needExistsTablesName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, TABLE_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistsTablesName), ","))
	}
	return nil
}
