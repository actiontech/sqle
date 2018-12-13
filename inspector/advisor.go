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

func (i *Inspect) CheckInvalid() {

}

func (i *Inspect) checkInvalidCreateTable(stmt *ast.CreateTableStmt) error {
	schemaName := i.getSchemaName(stmt.Table)
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		// add result, database not exist
	} else {
		tableExist, err := i.isTableExist(stmt.Table)
		if err != nil {
			return err
		}
		if tableExist {
			// add result, table is exist
		}
		referTableExist, err := i.isTableExist(stmt.ReferTable)
		if err != nil {
			return err
		}
		if referTableExist {
			// add result, table is exist
		}
	}
	colsName := []string{}
	colsNameMap := map[string]struct{}{}
	pkCounter := 0
	for _, col := range stmt.Cols {
		colsName = append(colsName, col.Name.Name.L)
		if HasOneInOptions(col.Options, ast.ColumnOptionPrimaryKey) {
			pkCounter += 1
		}
	}
	indexsName := []string{}
	keyColsName := []string{}
	for _, constraint := range stmt.Constraints {
		switch constraint.Tp {
		case ast.ConstraintPrimaryKey:
			pkCounter += 1
		case ast.ConstraintIndex, ast.ConstraintUniq, ast.ConstraintFulltext:
			if constraint.Name != "" {
				indexsName = append(indexsName, constraint.Name)
			}
			for _, col := range constraint.Keys {
				keyColsName = append(keyColsName, col.Column.Name.L)
			}
		}
	}
	if d := getDuplicate(colsName); len(d) > 0 {
		// add result, duplicate column
	}

	if d := getDuplicate(indexsName); len(d) > 0 {
		// add result, duplicate column name
	}

	if pkCounter > 1 {
		// add result, multiple primary key
	}
	notExistKeyColumns := []string{}
	for _, colName := range keyColsName {
		if _, ok := colsNameMap[colName]; !ok {
			notExistKeyColumns = append(notExistKeyColumns, colName)
		}
	}
	if len(notExistKeyColumns) > 0 {
		// add result, key column doesn't exist in table
	}
	return nil
}

func (i *Inspect) checkInvalidAlterTable(stmt *ast.AlterTableStmt) error {
	schemaName := i.getSchemaName(stmt.Table)
	schemaExist, err := i.isSchemaExist(schemaName)
	if err != nil {
		return err
	}
	if !schemaExist {
		// add result, database not exist
	} else {
		tableExist, err := i.isTableExist(stmt.Table)
		if err != nil {
			return err
		}
		if !tableExist {
			// add result, table not exist
		}
	}
	createTableStmt, _, err := i.getCreateTableStmt(stmt.Table)
	if err != nil {
		return err
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

	columnNeedNotExist := []string{}
	columnNeedExist := []string{}
	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddColumns) {
		for _, col := range spec.NewColumns {
			colName := col.Name.Name.L
			if _, ok := colNameMap[colName]; ok {
				columnNeedNotExist = append(columnNeedNotExist, colName)
			}
		}
	}

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAlterColumn) {
		for _, col := range spec.NewColumns {
			colName := col.Name.Name.L
			if _, ok := colNameMap[colName]; !ok {
				columnNeedExist = append(columnNeedExist, colName)
			}
		}
	}

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableChangeColumn) {
		oldColumnName := spec.OldColumnName.Name.L
		if _, ok := colNameMap[oldColumnName]; !ok {
			columnNeedExist = append(columnNeedExist, oldColumnName)
		}
		for _, col := range spec.NewColumns {
			newColName := col.Name.Name.L
			if newColName == oldColumnName {
				continue
			}
			if _, ok := colNameMap[newColName]; ok {
				columnNeedNotExist = append(columnNeedNotExist, newColName)
			}
		}
	}

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropColumn) {
		oldColName := spec.OldColumnName.Name.L
		if _, ok := colNameMap[oldColName]; !ok {
			columnNeedExist = append(columnNeedExist, oldColName)
		}
	}

	indexNeedNotExist := []string{}
	indexNeedExist := []string{}

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableAddConstraint) {
		switch spec.Constraint.Tp {
		// add primary key
		case ast.ConstraintPrimaryKey:
			if hasPk {
				// add result, pk has exist
			}
		case ast.ConstraintUniq, ast.ConstraintIndex, ast.ConstraintFulltext:
			indexName := strings.ToLower(spec.Constraint.Name)
			if _, ok := indexNameMap[indexName]; ok {
				indexNeedNotExist = append(indexNeedNotExist, indexName)
			}
		}
	}

	if len(getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropPrimaryKey)) > 0 && !hasPk {
		// add result, pk not exist
	}

	for _, spec := range getAlterTableSpecByTp(stmt.Specs, ast.AlterTableDropIndex) {
		indexName := strings.ToLower(spec.Name)
		if _, ok := indexNameMap[indexName]; !ok {
			indexNeedExist = append(indexNeedExist, indexName)
		}
	}

	if len(columnNeedExist) > 0 {
		// add result columns not exist
	}
	if len(columnNeedNotExist) > 0 {
		// add result columns has exist
	}
	if len(indexNeedExist) > 0 {
		// add result index not exist
	}
	if len(indexNeedNotExist) > 0 {
		// add result index has exist
	}
	return nil
}
