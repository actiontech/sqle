package index

import "github.com/pingcap/parser/ast"

// SelectAST contain an abstract syntax tree of a single table SQL. It
// abstracts the syntax differences between different database.
type SelectAST interface {
	// EqualPredicateColumnsInWhere find the equal predicate column in where clause.
	//
	// For example, the SQL: select * from t where a = 1 and b = 2;
	// it returns []string{"a", "b"}.
	EqualPredicateColumnsInWhere() []string
	// EqualPredicateColumnsInWhere find the unequal predicate column in where clause.
	//
	// For example, the SQL: select * from t where a >= 1 and b <= 2;
	// it returns []string{"a", "b"}.
	UnequalPredicateColumnsInWhere() []string

	// ColumnsInOrderBy find the columns in order by clause.
	//
	// For example, the SQL: select * from t order by a desc, b;
	// it returns []string{"a desc", "b"}.
	ColumnsInOrderBy() []string

	// ColumnsInProjection find columns in select projection.
	//
	// For example, the SQL: select a, b from t;
	// it returns []string{"a", "b"}.
	//
	// If projection returns all columns, it returns nil.
	ColumnsInProjection() []string

	// GetCreateTableStmt get create table SQLs of the used table in select stmt
	// For example, the SQL: select a, b from t;
	// it returns []*ast.TableName{/* table info of t */}
	GetSelectedTables() []*ast.TableName
}
