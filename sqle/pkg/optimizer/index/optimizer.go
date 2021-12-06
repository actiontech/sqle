// Package index implements the general index optimization algorithm.
//
// The main idea of the algorithm is to find the best index
// for the query by three-star index algorithm. For more information,
// please refer to the Chapter 4 of the book 《Relational database index design and the optimizers》.
//
// A single table SELECT using a three-star index normally needs only one random disk drive read and
// a scan of a thin slice of an index.
//
// The three-star index algorithm described in the following steps:
//	1. An index deserves the first star if the index rows
//	   relevant to the SELECT are next to each other or at least as
//	   close to each other as possible. This minimizes the thickness
//	   of the index slice that must be scanned.
//	2. The second star is given if the index rows are in the right
//	   order for the SELECT.
//	3. If the index rows contain all the columns referred to
//	   by the SELECT the index is given the third star.
//	   This eliminates table access: The access path is index only.
//
// Suppose we have a SQL:
// 	SELECT CNO, FNAME
//	FROM CUST
//	WHERE LNAME = :LNAME AND CITY = :CITY
//	ORDER BY FNAME
//
// To Qualify for the first star:
// Pick the columns from all equal predicates (WHERE COL = . . .). Make these
// the first columns of the index—in any order. For above SQL, the three-star
// index will begin with columns LNAME, CITY or CITY, LNAME. In both cases the
// index slice that must be scanned will be as thin as possible.
//
// To Qualify for the second star:
// Add the ORDER BY columns. Do not change the order of these columns, but
// ignore columns that were already picked in step 1. For example, if above SQL
// had redundant columns in the ORDER BY, say ORDER BY LNAME, FNAME or ORDER BY
// FNAME, CITY, only FNAME would be added in this step. When FNAME is the third
// index column, the result table will be in the right order without sorting. The
// first FETCH call will return the row with the smallest FNAME value.
//
// To Qualify for the third star:
// Add all the remaining columns from the SELECT statement. The order of the columns
// added in this step has no impact on the performance of the SELECT, but the cost of
// updates should be reduced by placing volatile columns at the end. Now the index
// contains all the columns required for an index-only access path.
package index

const (
	defaultCompositeIndexMaxColumn = 3
)

// Optimizer give best index advice for the single table query.
type Optimizer struct {
	compositeIndexMaxColumn int
}

// NewOptimizer creates a new optimizer.
func NewOptimizer(opts ...optimizerOption) *Optimizer {
	optimizer := &Optimizer{
		compositeIndexMaxColumn: defaultCompositeIndexMaxColumn,
	}

	for _, opt := range opts {
		opt.apply(optimizer)
	}

	return optimizer
}

// Optimize try him best to give three-star index advice for ast.
func (o *Optimizer) Optimize(ast SelectAST) (columns []string, err error) {
	columns = append(ast.EqualPredicateColumnsInWhere(), ast.ColumnsInOrderBy()...)

	columns = append(columns, ast.ColumnsInProjection()...)

	if len(columns) > o.compositeIndexMaxColumn {
		columns = columns[:o.compositeIndexMaxColumn]
	}

	return
}

type optimizerOption func(*Optimizer)

func (opt optimizerOption) apply(o *Optimizer) {
	opt(o)
}

// WithCompositeIndexMaxColumn sets the max column number of the composite index.
func WithCompositeIndexMaxColumn(column int) optimizerOption {
	return func(o *Optimizer) {
		o.compositeIndexMaxColumn = column
	}
}
