package inspector

import (
	"github.com/pingcap/tidb/ast"
	"github.com/pingcap/tidb/mysql"
	"sqle/model"
	"strings"
)

func (i *Inspect) Advise(rules []model.Rule) error {
	i.initRulesFunc()
	defer i.closeDbConn()

	for _, commitSql := range i.Task.CommitSqls {
		currentSql := commitSql
		err := i.Add(&currentSql.Sql, func(sql *model.Sql) error {
			for _, rule := range rules {
				i.currentRule = rule
				if fn, ok := i.RulesFunc[rule.Name]; ok {
					if fn == nil {
						continue
					}
					for _, node := range sql.Stmts {
						err := fn(node, rule.Name)
						if err != nil {
							return err
						}
					}
				}
			}
			currentSql.InspectStatus = model.TASK_ACTION_DONE
			currentSql.InspectLevel = i.Results.level()
			currentSql.InspectResult = i.Results.message()

			// print osc

			// clean up results
			i.Results = newInspectResults()
			return nil
		})
		if err != nil {
			return err
		}
	}
	return i.Do()
}

func (i *Inspect) checkSelectAll(node ast.StmtNode, rule string) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		// check select all column
		if stmt.Fields != nil && stmt.Fields.Fields != nil {
			for _, field := range stmt.Fields.Fields {
				if field.WildCard != nil {
					i.addResult(DML_DISABE_SELECT_ALL_COLUMN)
				}
			}
		}
	}
	return nil
}

func (i *Inspect) checkSelectWhere(node ast.StmtNode, rule string) error {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		// where condition
		if stmt.Where == nil || !whereStmtHasOneColumn(stmt.Where) {
			i.addResult(DML_CHECK_INVALID_WHERE_CONDITION)
		}
	}
	return nil
}

func (i *Inspect) checkPrimaryKey(node ast.StmtNode, rule string) error {
	var hasPk = false
	var pkIsAutoIncrementBigIntUnsigned = false

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check primary key
		// TODO: tidb parser not support keyword for SERIAL; it is a alias for "BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE"
		/*
			match sql like:
			CREATE TABLE  tb1 (
			a1.id int(10) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
			);
		*/
		for _, col := range stmt.Cols {
			if IsAllInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
				hasPk = true
				if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) &&
					IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
					pkIsAutoIncrementBigIntUnsigned = true
				}
			}
		}
		/*
			match sql like:
			CREATE TABLE  tb1 (
			a1.id int(10) unsigned NOT NULL AUTO_INCREMENT,
			PRIMARY KEY (id)
			);
		*/
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				hasPk = true
				if len(constraint.Keys) == 1 {
					columnName := constraint.Keys[0].Column.Name.String()
					for _, col := range stmt.Cols {
						if col.Name.Name.String() == columnName {
							if col.Tp.Tp == mysql.TypeLonglong && mysql.HasUnsignedFlag(col.Tp.Flag) &&
								IsAllInOptions(col.Options, ast.ColumnOptionAutoIncrement) {
								pkIsAutoIncrementBigIntUnsigned = true
							}
						}
					}
				}
			}
		}
	default:
		return nil
	}

	if !hasPk {
		i.addResult(DDL_CHECK_PRIMARY_KEY_EXIST)
	}
	if hasPk && !pkIsAutoIncrementBigIntUnsigned {
		i.addResult(DDL_CHECK_PRIMARY_KEY_TYPE)
	}
	return nil
}

func (i *Inspect) checkMergeAlterTable(node ast.StmtNode, rule string) error {
	switch stmt := node.(type) {
	case *ast.AlterTableStmt:
		// merge alter table
		tableName := i.getTableName(stmt.Table)
		_, ok := i.alterTableStmts[tableName]
		if ok {
			i.addResult(DDL_CHECK_ALTER_TABLE_NEED_MERGE)
			i.alterTableStmts[tableName] = append(i.alterTableStmts[tableName], stmt)
		} else {
			i.alterTableStmts[tableName] = []*ast.AlterTableStmt{stmt}
		}
	}

	return nil
}

func (i *Inspect) checkEngineAndCharacterSet(node ast.StmtNode, rule string) error {
	var engine string
	var characterSet string
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, op := range stmt.Options {
			switch op.Tp {
			case ast.TableOptionEngine:
				engine = op.StrValue
			case ast.TableOptionCharset:
				characterSet = op.StrValue
			}
		}
	default:
		return nil
	}
	if strings.ToLower(engine) == "innodb" && strings.ToLower(characterSet) == "utf8mb4" {
		return nil
	}
	i.addResult(DDL_TABLE_USING_INNODB_UTF8MB4)
	return nil
}

func (i *Inspect) disableAddIndexForColumnsTypeBlob(node ast.StmtNode, rule string) error {
	indexColumns := map[string]struct{}{}
	isTypeBlobCols := map[string]bool{}
	indexDataTypeIsBlob := false
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				for _, col := range constraint.Keys {
					indexColumns[col.Column.Name.String()] = struct{}{}
				}
			}
		}
		for _, col := range stmt.Cols {
			if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					indexDataTypeIsBlob = true
					break
				}
			}
			if _, ok := indexColumns[col.Name.Name.String()]; ok {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					indexDataTypeIsBlob = true
					break
				}
			}
		}
	case *ast.AlterTableStmt:
		// collect index column
		for _, spec := range stmt.Specs {
			if spec.NewColumns == nil {
				continue
			}
			for _, col := range spec.NewColumns {
				if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) {
					indexColumns[col.Name.Name.String()] = struct{}{}
				}
			}
			if spec.Constraint != nil {
				switch spec.Constraint.Tp {
				case ast.ConstraintKey, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintIndex:
					for _, col := range spec.Constraint.Keys {
						indexColumns[col.Column.Name.String()] = struct{}{}
					}
				}
			}
		}
		if len(indexColumns) <= 0 {
			return nil
		}

		// collect columns type
		createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
		if err != nil {
			return err
		}
		if exist {
			for _, col := range createTableStmt.Cols {
				if MysqlDataTypeIsBlob(col.Tp.Tp) {
					isTypeBlobCols[col.Name.Name.String()] = true
				} else {
					isTypeBlobCols[col.Name.Name.String()] = false
				}
			}
		}
		for _, spec := range stmt.Specs {
			if spec.NewColumns != nil {
				for _, col := range spec.NewColumns {
					if MysqlDataTypeIsBlob(col.Tp.Tp) {
						isTypeBlobCols[col.Name.Name.String()] = true
					} else {
						isTypeBlobCols[col.Name.Name.String()] = false
					}
				}
			}
		}
		// check index columns string type
		for colName, _ := range indexColumns {
			if isTypeBlobCols[colName] {
				indexDataTypeIsBlob = true
				break
			}
		}
	case *ast.CreateIndexStmt:
		createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
		if err != nil || !exist {
			return err
		}
		for _, col := range createTableStmt.Cols {
			if HasOneInOptions(col.Options, ast.ColumnOptionUniqKey) && MysqlDataTypeIsBlob(col.Tp.Tp) {
				isTypeBlobCols[col.Name.Name.String()] = true
			} else {
				isTypeBlobCols[col.Name.Name.String()] = false
			}
		}
		for _, indexColumns := range stmt.IndexColNames {
			if isTypeBlobCols[indexColumns.Column.Name.String()] {
				indexDataTypeIsBlob = true
				break
			}
		}
	default:
		return nil
	}
	if indexDataTypeIsBlob {
		i.addResult(DDL_DISABLE_INDEX_DATA_TYPE_BLOB)
	}
	return nil
}

func (i *Inspect) checkNewObjectName(node ast.StmtNode, rule string) error {
	names := []string{}
	invalidNames := []string{}

	switch stmt := node.(type) {
	case *ast.CreateDatabaseStmt:
		// schema
		names = append(names, stmt.Name)
	case *ast.CreateTableStmt:

		// table
		names = append(names, stmt.Table.Name.String())

		// column
		for _, col := range stmt.Cols {
			names = append(names, col.Name.Name.String())
		}
		// index
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintUniqKey, ast.ConstraintKey, ast.ConstraintUniqIndex, ast.ConstraintIndex:
				names = append(names, constraint.Name)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableRenameTable:
				// rename table
				names = append(names, spec.NewTable.Name.String())
			case ast.AlterTableAddColumns:
				// new column
				for _, col := range spec.NewColumns {
					names = append(names, col.Name.Name.String())
				}
			case ast.AlterTableChangeColumn:
				// rename column
				for _, col := range spec.NewColumns {
					names = append(names, col.Name.Name.String())
				}
			case ast.AlterTableAddConstraint:
				// if spec.Constraint.Name not index name, it will be null
				names = append(names, spec.Constraint.Name)
			case ast.AlterTableRenameIndex:
				names = append(names, spec.ToKey.String())
			}
		}
	case *ast.CreateIndexStmt:
		names = append(names, stmt.IndexName)
	default:
		return nil
	}

	// check length
	for _, name := range names {
		if len(name) > 64 {
			i.addResult(DDL_CHECK_OBJECT_NAME_LENGTH)
			break
		}
	}
	// check keyword
	for _, name := range names {
		if IsMysqlReservedKeyword(name) {
			invalidNames = append(invalidNames, name)
		}
	}
	if len(invalidNames) > 0 {
		i.addResult(DDL_DISABLE_USING_KEYWORD, strings.Join(RemoveArrayRepeat(invalidNames), ", "))
	}
	return nil
}

func (i *Inspect) checkForeignKey(node ast.StmtNode, rule string) error {
	hasFk := false

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintForeignKey {
				hasFk = true
				break
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Constraint != nil && spec.Constraint.Tp == ast.ConstraintForeignKey {
				hasFk = true
				break
			}
		}
	default:
		return nil
	}
	if hasFk {
		i.addResult(DDL_DISABLE_FOREIGN_KEY)
	}
	return nil
}

func (i *Inspect) checkIndex(node ast.StmtNode, rule string) error {
	indexCounter := 0
	compositeIndexMax := 0

	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check index
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				indexCounter++
				if compositeIndexMax < len(constraint.Keys) {
					compositeIndexMax = len(constraint.Keys)
				}
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			if spec.Constraint == nil {
				continue
			}
			switch spec.Constraint.Tp {
			case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
				indexCounter++
				if compositeIndexMax < len(spec.Constraint.Keys) {
					compositeIndexMax = len(spec.Constraint.Keys)
				}
			}
		}
		createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
		if err != nil {
			return err
		}
		if exist {
			for _, constraint := range createTableStmt.Constraints {
				switch constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					indexCounter++
				}
			}
		}

	case *ast.CreateIndexStmt:
		indexCounter++
		if compositeIndexMax < len(stmt.IndexColNames) {
			compositeIndexMax = len(stmt.IndexColNames)
		}
		createTableStmt, exist, err := i.getCreateTableStmt(i.getTableName(stmt.Table))
		if err != nil {
			return err
		}
		if exist {
			for _, constraint := range createTableStmt.Constraints {
				switch constraint.Tp {
				case ast.ConstraintIndex, ast.ConstraintUniqIndex, ast.ConstraintKey, ast.ConstraintUniqKey:
					indexCounter++
				}
			}
		}
	default:
		return nil
	}
	if indexCounter > 5 {
		i.addResult(DDL_CHECK_INDEX_COUNT)
	}
	if compositeIndexMax > 5 {
		i.addResult(DDL_CHECK_COMPOSITE_INDEX_MAX)
	}
	return nil
}

func (i *Inspect) checkStringType(node ast.StmtNode, rule string) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// if char length >20 using varchar.
		for _, col := range stmt.Cols {
			if col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
				i.addResult(DDL_CHECK_TYPE_CHAR_LENGTH)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			for _, col := range spec.NewColumns {
				if col.Tp.Tp == mysql.TypeString && col.Tp.Flen > 20 {
					i.addResult(DDL_CHECK_TYPE_CHAR_LENGTH)
				}
			}
		}
	default:
		return nil
	}
	return nil
}

func (i *Inspect) checkObjectExist(node ast.StmtNode, rule string) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check schema
		schema := i.getSchemaName(stmt.Table)
		tableName := i.getTableName(stmt.Table)
		exist, err := i.isSchemaExist(schema)
		if err != nil {
			return err
		}
		if !exist {
			// if schema not exist, table must not exist
			return nil

		} else {
			// check table if schema exist
			exist, err = i.isTableExist(tableName)
			if err != nil {
				return err
			}
			if exist {
				i.addResult(TABLE_EXIST, tableName)
			}
		}
	case *ast.CreateDatabaseStmt:
		schemaName := stmt.Name
		exist, err := i.isSchemaExist(schemaName)
		if err != nil {
			return err
		}
		if exist {
			i.addResult(SCHEMA_EXIST, schemaName)
		}
	}
	return nil
}

func (i *Inspect) checkObjectNotExist(node ast.StmtNode, rule string) error {
	var tablesName = []string{}
	var schemasName = []string{}

	switch stmt := node.(type) {
	case *ast.UseStmt:
		schemasName = append(schemasName, stmt.DBName)

	case *ast.CreateTableStmt:
		schemasName = append(schemasName, i.getSchemaName(stmt.Table))

	case *ast.AlterTableStmt:
		schemasName = append(schemasName, i.getSchemaName(stmt.Table))
		tablesName = append(tablesName, i.getTableName(stmt.Table))

	case *ast.SelectStmt:
		for _, table := range getTables(stmt.From.TableRefs) {
			schemasName = append(schemasName, i.getSchemaName(table))
			tablesName = append(tablesName, i.getTableName(table))
		}
	case *ast.InsertStmt:
		for _, table := range getTables(stmt.Table.TableRefs) {
			schemasName = append(schemasName, i.getSchemaName(table))
			tablesName = append(tablesName, i.getTableName(table))
		}

	case *ast.DeleteStmt:
		if stmt.Tables != nil && stmt.Tables.Tables != nil {
			for _, table := range stmt.Tables.Tables {
				schemasName = append(schemasName, i.getSchemaName(table))
				tablesName = append(tablesName, i.getTableName(table))
			}
		}
		for _, table := range getTables(stmt.TableRefs.TableRefs) {
			schemasName = append(schemasName, i.getSchemaName(table))
			tablesName = append(tablesName, i.getTableName(table))
		}

	case *ast.UpdateStmt:
		for _, table := range getTables(stmt.TableRefs.TableRefs) {
			schemasName = append(schemasName, i.getSchemaName(table))
			tablesName = append(tablesName, i.getTableName(table))
		}
	}

	notExistSchemas := []string{}
	for _, schema := range schemasName {
		exist, err := i.isSchemaExist(schema)
		if err != nil {
			return err
		}
		if !exist {
			notExistSchemas = append(notExistSchemas, schema)
		}
	}
	if len(notExistSchemas) > 0 {
		i.addResult(SCHEMA_NOT_EXIST, strings.Join(RemoveArrayRepeat(notExistSchemas), ", "))
	}

	notExistTables := []string{}
	for _, table := range tablesName {
		exist, err := i.isTableExist(table)
		if err != nil {
			return err
		}
		if !exist {
			notExistTables = append(notExistTables, table)
		}
	}
	if len(notExistTables) > 0 {
		i.addResult(TABLE_NOT_EXIST, strings.Join(RemoveArrayRepeat(notExistTables), ", "))
	}
	return nil
}

func (i *Inspect) checkIfNotExist(node ast.StmtNode, rule string) error {
	switch stmt := node.(type) {
	case *ast.CreateTableStmt:
		// check `if not exists`
		if !stmt.IfNotExists {
			i.addResult(DDL_CREATE_TABLE_NOT_EXIST)
		}
	}
	return nil
}

func (i *Inspect) disableDropStmt(node ast.StmtNode, rule string) error {
	// specific check
	switch node.(type) {
	case *ast.DropDatabaseStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	case *ast.DropTableStmt:
		i.addResult(DDL_DISABLE_DROP_STATEMENT)
	}
	return nil
}

func (i *Inspect) checkMycatShardingColumn(node ast.StmtNode, rule string) error {
	if i.Task.Instance.DbType != model.DB_TYPE_MYCAT {
		return nil
	}
	config := i.Task.Instance.MycatConfig
	hasShardingColumn := false
	switch stmt := node.(type) {
	case *ast.InsertStmt:
		tables := getTables(stmt.Table.TableRefs)
		// tables must be one on InsertIntoStmt in parser.go
		if len(tables) != 1 {
			return nil
		}
		table := tables[0]
		schema, ok := config.AlgorithmSchemas[i.getSchemaName(table)]
		if !ok {
			return nil
		}
		tableName := table.Name.String()
		if schema.AlgorithmTables == nil {
			return nil
		}
		at, ok := schema.AlgorithmTables[tableName]
		if !ok {
			return nil
		}
		shardingCoulmn := at.ShardingColumn
		if stmt.Columns != nil {
			for _, column := range stmt.Columns {
				if column.Name.L == strings.ToLower(shardingCoulmn) {
					hasShardingColumn = true
				}
			}
		}
		if stmt.Setlist != nil {
			for _, set := range stmt.Setlist {
				if set.Column.Name.L == strings.ToLower(shardingCoulmn) {
					hasShardingColumn = true
				}
			}
		}
	case *ast.UpdateStmt:
		tables := getTables(stmt.TableRefs.TableRefs)
		// multi table related update not supported on mycat
		if len(tables) != 1 {
			return nil
		}
		table := tables[0]
		schema, ok := config.AlgorithmSchemas[i.getSchemaName(table)]
		if !ok {
			return nil
		}
		tableName := table.Name.String()
		if schema.AlgorithmTables == nil {
			return nil
		}
		at, ok := schema.AlgorithmTables[tableName]
		if !ok {
			return nil
		}
		shardingCoulmn := at.ShardingColumn
		hasShardingColumn = whereStmtHasSpecificColumn(stmt.Where, shardingCoulmn)
	case *ast.DeleteStmt:
		// not support multi table related delete
		if stmt.IsMultiTable {
			return nil
		}
		tables := getTables(stmt.TableRefs.TableRefs)
		if len(tables) != 1 {
			return nil
		}
		table := tables[0]
		schema, ok := config.AlgorithmSchemas[i.getSchemaName(table)]
		if !ok {
			return nil
		}
		tableName := table.Name.String()
		if schema.AlgorithmTables == nil {
			return nil
		}
		at, ok := schema.AlgorithmTables[tableName]
		if !ok {
			return nil
		}
		shardingCoulmn := at.ShardingColumn
		hasShardingColumn = whereStmtHasSpecificColumn(stmt.Where, shardingCoulmn)
	}
	if !hasShardingColumn {
		i.addResult(DML_MYCAT_MUST_USING_SHARDING_CLOUNM)
	}
	return nil
}
