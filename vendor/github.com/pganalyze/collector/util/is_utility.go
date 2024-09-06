package util

import (
	pg_query "github.com/pganalyze/pg_query_go/v4"
)

// IsUtilityStmt determines whether each statement in the query text is a
// utility statement or a standard SELECT/INSERT/UPDATE/DELETE statement.
func IsUtilityStmt(query string) ([]bool, error) {
	var result []bool
	parseResult, err := pg_query.Parse(query)
	if err != nil {
		return nil, err
	}
	for _, rawStmt := range parseResult.Stmts {
		stmt := rawStmt.Stmt.Node
		var isUtility bool
		switch stmt.(type) {
		case *pg_query.Node_SelectStmt, *pg_query.Node_InsertStmt, *pg_query.Node_UpdateStmt, *pg_query.Node_DeleteStmt:
			isUtility = false
		default:
			isUtility = true
		}
		result = append(result, isUtility)
	}
	return result, nil
}
