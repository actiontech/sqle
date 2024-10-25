package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00055 = "SQLE00055"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00055,
			Desc:       "对于MySQL的索引, 不建议创建冗余索引",
			Annotation: "MySQL需要单独维护重复的索引，冗余索引增加维护成本，影响更新性能",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeIndexOptimization,
		},
		Message: "已存在索引 %v , 索引 %v 为冗余索引",
		AllowOffline: false,
		Func:    RuleSQLE00055,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00055): "For table creation and index creation statements, creating redundant indexes is prohibited .".
You should follow the following logic:
1. For the "CREATE TABLE ..." statements, builds a list of index columns, which is used to record all the declared indexes and their columns, checking whether the columns of each index are redundant, meaning that the index columns are exactly the same, or have the same leftmost prefix.  If it does, report a violation.
2. For the  "CREATE INDEX ..." statements, builds a list of index columns to keep track of the existing indexes and their columns, and check that the new index is not redundant, meaning that the index field is the same as the old one or has the same leftmost prefix. If it does, report a violation.
3. For the  "ALTER TABLE ... ADD INDEX ..." statements, perform the same check as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00055(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table..."
		indexes := [][]string{}

		// get index column in column definition
		for _, col := range stmt.Cols {
			if util.IsColumnHasOption(col, ast.ColumnOptionUniqKey) || util.IsColumnPrimaryKey(col) {
				indexes = append(indexes, []string{util.GetColumnName(col)})
			}
		}

		// get index column in table constraint
		indexes = append(indexes, extractIndexesFromConstraints(util.GetTableConstraints(stmt.Constraints, util.GetIndexConstraintTypes()...))...)

		// check the index is not duplicated
		for redundant, source := range calculateIndexRedundant(indexes) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00055, indexes[source], indexes[redundant])
		}

	case *ast.CreateIndexStmt:
		// "create index..."

		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		indexes := extractIndexesFromConstraints(util.GetTableConstraints(createTableStmt.Constraints, util.GetIndexConstraintTypes()...))
		newIndex := extractIndexesFromIndexStmt(stmt.IndexPartSpecifications)

		// check the index is not duplicated
		for redundant, source := range calculateNewIndexRedundant(indexes, newIndex) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00055, indexes[source], newIndex[redundant])
		}

	case *ast.AlterTableStmt:
		// "alter table"

		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if nil != err {
			return err
		}

		indexes := extractIndexesFromConstraints(util.GetTableConstraints(createTableStmt.Constraints, util.GetIndexConstraintTypes()...))

		newIndex := [][]string{}
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// "alter table... add index..."
			indexCols := extractIndexesFromConstraints(util.GetTableConstraints([]*ast.Constraint{spec.Constraint}, util.GetIndexConstraintTypes()...))
			newIndex = append(newIndex, indexCols...)

		}
		// check the index is not duplicated
		for redundant, source := range calculateNewIndexRedundant(indexes, newIndex) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00055, indexes[source], newIndex[redundant])
		}
	}

	return nil
}

func extractIndexesFromConstraints(constraints []*ast.Constraint) [][]string {
	indexes := [][]string{}

	// Iterate over all constraints to extract index columns.
	for _, constraint := range constraints {
		indexCols := []string{}

		// Collect all column names that are part of the index represented by the constraint.
		for _, key := range constraint.Keys {
			colName := util.GetIndexColName(key)
			if colName != "" {
				indexCols = append(indexCols, colName)
			}
		}

		// Only append to indexes if indexCols is not empty, avoiding adding empty index definitions.
		if len(indexCols) > 0 {
			indexes = append(indexes, indexCols)
		}
	}

	return indexes
}

func extractIndexesFromIndexStmt(index []*ast.IndexPartSpecification) [][]string {
	// Initialize the slice that will hold the column names for the index.
	var indexCols []string

	// Iterate over each IndexPartSpecification to extract the column names.
	for _, col := range index {
		if colName := util.GetIndexColName(col); colName != "" {
			indexCols = append(indexCols, colName)
		}
	}

	// Return the indexes as a slice of a single index if indexCols is not empty.
	if len(indexCols) > 0 {
		return [][]string{indexCols}
	}
	return [][]string{} // Return an empty slice of indexes if no index columns were found.
}

// calculateIndexRedundant takes a slice of string slices representing indexes and their columns
// and returns a map where the key is the redundant index and the value is the corresponding
// original index that makes it redundant.
func calculateIndexRedundant(indexes [][]string) map[int]int {
	redundantIndexes := make(map[int]int)

	// Compare each index with every other index to check for redundancy.
	for i, indexColumns := range indexes {
		for j, otherIndexColumns := range indexes {
			// Skip comparing the same index.
			if i == j {
				continue
			}

			// Check if indexColumns are redundant with respect to otherIndexColumns.
			if isRedundant(indexColumns, otherIndexColumns) {
				// If index i is redundant and either it's not in the map yet, or it's in the map
				// but the current source index j has fewer columns (and thus is a 'stronger' source of redundancy),
				// then add/update the map with the index pair (i, j).
				if _, exists := redundantIndexes[i]; !exists || len(indexes[redundantIndexes[i]]) > len(otherIndexColumns) {
					redundantIndexes[i] = j
				}
			}
		}
	}

	return redundantIndexes
}

// calculateNewIndexRedundant takes a slice of string slices representing existing indexes and their columns
// and a slice of string slices representing new indexes' columns. It returns a map where the key is the index
// of the new index in `newIndexes` and the value is the index of the existing index in `indexes` that makes it redundant.
func calculateNewIndexRedundant(indexes [][]string, newIndexes [][]string) map[int]int {
	redundantIndexes := make(map[int]int)

	// Iterate over each new index
	for newIndexI, newIndexCols := range newIndexes {
		// Compare against each existing index
		for existingIndexI, existingIndexCols := range indexes {
			if isRedundant(newIndexCols, existingIndexCols) {
				// If the new index is redundant with respect to an existing index, add to the map.
				redundantIndexes[newIndexI] = existingIndexI
				// Since we only care about the first occurrence of redundancy, we can break here.
				break
			}
		}
	}

	return redundantIndexes
}

// isRedundant checks if the first index is redundant with respect to the second index,
// meaning all columns in the first index are a leftmost prefix of the second.
func isRedundant(index1, index2 []string) bool {
	for i, col := range index1 {
		// If we reach the end of index2 or the columns differ,
		// then index1 cannot be a leftmost prefix of index2.
		if i >= len(index2) || col != index2[i] {
			return false
		}
	}
	return true
}

// ==== Rule code end ====
