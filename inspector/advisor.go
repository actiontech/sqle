package inspector

import (
	"fmt"
	"github.com/pingcap/tidb/ast"
	"sqle/model"
	"strings"
)

func (i *Inspect) Advise(rules []model.Rule) error {
	defer i.closeDbConn()
	i.Logger().Info("start advise sql")
	for _, commitSql := range i.Task.CommitSqls {
		currentSql := commitSql
		err := i.Add(&currentSql.Sql, func(sql *model.Sql) error {
			for _, node := range sql.Stmts {
				err := i.CheckInvalid(node, i.Results)
				if err != nil {
					return err
				}
			}
			for _, rule := range rules {
				i.currentRule = rule
				if handler, ok := RuleHandlerMap[rule.Name]; ok {
					if handler.Func == nil {
						continue
					}
					for _, node := range sql.Stmts {
						err := handler.Func(i, node)
						if err != nil {
							return err
						}
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
				if currentSql.InspectResult != "" {
					currentSql.InspectResult += "\n"
				}
				currentSql.InspectResult = fmt.Sprintf("%s[osc]%s",
					currentSql.InspectResult, oscCommandLine)
			}

			return nil
		})
		if err != nil {
			i.Logger().Error("add commit sql to task failed")
			return err
		}
	}
	err := i.Do()
	if err != nil {
		i.Logger().Error("advise sql failed")
	} else {
		i.Logger().Info("advise sql finish")
	}
	return err
}

var (
	SCHEMA_NOT_EXIST_MSG        = "schema %s 不存在"
	SCHEMA_EXIST_MSG            = "schema %s 已存在"
	TABLE_NOT_EXIST_MSG         = "表 %s 不存在"
	TABLE_EXIST_MSG             = "表 %s 已存在"
	COLUMN_NOT_EXIST_MSG        = "字段 %s 不存在"
	COLUMN_EXIST_MSG            = "字段 %s 已存在"
	INDEX_NOT_EXIST_MSG         = "索引 %s 不存在"
	INDEX_EXIST_MSG             = "索引 %s 已存在"
	DUPLICATE_COLUMN_ERROR_MSG  = "字段名 %s 重复"
	DUPLICATE_INDEX_ERROR_MSG   = "索引名 %s 重复"
	PRIMARY_KEY_MULTI_ERROR_MSG = "主键只能设置一个"
	KEY_COLUMN_NOT_EXIST_MSG    = "索引字段 %s 不存在"
	PRIMARY_KEY_EXIST_MSG       = "已经存在主键，不能再添加"
	PRIMARY_KEY_NOT_EXIST_MSG   = "当前没有主键，不能执行删除"
)

func (i *Inspect) CheckInvalid(node ast.Node, results *InspectResults) error {
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
	}
	return err
}

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
		case ast.ConstraintIndex, ast.ConstraintUniq, ast.ConstraintFulltext:
			if constraint.Name != "" {
				indexesName = append(indexesName, constraint.Name)
			}
			for _, col := range constraint.Keys {
				keyColsName = append(keyColsName, col.Column.Name.L)
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
		} else {
			delete(colNameMap, oldColName)
		}
		for _, col := range spec.NewColumns {
			newColName := col.Name.Name.L
			if newColName == oldColName {
				continue
			}
			if _, ok := colNameMap[newColName]; ok {
				needNotExistsColsName = append(needNotExistsColsName, newColName)
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

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
		switch spec.Constraint.Tp {
		case ast.ConstraintPrimaryKey:
			if hasPk {
				// primary key has exist, can not add primary key
				results.add(model.RULE_LEVEL_ERROR, PRIMARY_KEY_EXIST_MSG)
			}
		case ast.ConstraintUniq, ast.ConstraintIndex, ast.ConstraintFulltext:
			indexName := strings.ToLower(spec.Constraint.Name)
			if _, ok := indexNameMap[indexName]; ok {
				needNotExistsIndexesName = append(needNotExistsIndexesName, indexName)
			}
			for _, col := range spec.Constraint.Keys {
				colName := col.Column.Name.L
				if _, ok := colNameMap[colName]; !ok {
					needExistsKeyColsName = append(needExistsKeyColsName, colName)
				}
			}
		}
	}

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
	keyColNeedExist := []string{}
	for _, col := range stmt.IndexColNames {
		colName := col.Column.Name.L
		if _, ok := colNameMap[col.Column.Name.L]; !ok {
			keyColNeedExist = append(keyColNeedExist, colName)
		}
	}
	if len(keyColNeedExist) > 0 {
		results.add(model.RULE_LEVEL_ERROR, KEY_COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(keyColNeedExist), ","))
	}
	return nil
}

func (i *Inspect) checkInvalidDropIndex(stmt *ast.DropIndexStmt,
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
	if _, ok := indexNameMap[stmt.IndexName]; !ok {
		results.add(model.RULE_LEVEL_ERROR, INDEX_NOT_EXIST_MSG, stmt.IndexName)
	}
	return nil
}

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
	needExistColsName := []string{}
	if stmt.Lists != nil {
		for _, col := range stmt.Columns {
			insertColsName = append(insertColsName, col.Name.L)
		}
	} else if stmt.Setlist != nil {
		for _, set := range stmt.Setlist {
			insertColsName = append(insertColsName, set.Column.Name.L)
		}
	}
	for _, colName := range insertColsName {
		if _, ok := colNameMap[colName]; !ok {
			needExistColsName = append(needExistColsName, colName)
		}
	}
	if len(needExistColsName) > 0 {
		results.add(model.RULE_LEVEL_ERROR, COLUMN_NOT_EXIST_MSG,
			strings.Join(removeDuplicate(needExistColsName), ","))
	}
	if d := getDuplicate(insertColsName); len(d) > 0 {
		results.add(model.RULE_LEVEL_ERROR, DUPLICATE_COLUMN_ERROR_MSG, strings.Join(d, ","))
	}
	return nil
}

func (i *Inspect) checkInvalidUpdate(stmt *ast.UpdateStmt, results *InspectResults) error {
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
	if stmt.MultipleTable {
		return nil
	}
	return nil
}

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
	return nil
}
