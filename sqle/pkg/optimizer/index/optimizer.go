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
//  1. An index deserves the first star if the index rows
//     relevant to the SELECT are next to each other or at least as
//     close to each other as possible. This minimizes the thickness
//     of the index slice that must be scanned.
//  2. The second star is given if the index rows are in the right
//     order for the SELECT.
//  3. If the index rows contain all the columns referred to
//     by the SELECT the index is given the third star.
//     This eliminates table access: The access path is index only.
//
// Suppose we have a SQL:
//
//	SELECT CNO, FNAME
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

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/utils"
)

// Optimizer give best index advice for the single table query.
type Optimizer struct {
	*session.Context
}

// NewOptimizer creates a new optimizer.
func NewOptimizer(ctx *session.Context, opts ...optimizerOption) *Optimizer {
	optimizer := &Optimizer{
		ctx,
	}

	for _, opt := range opts {
		opt.apply(optimizer)
	}

	return optimizer
}

/*
GiveThreeStarAdvice 三星索引建议

	三星索引要求:
	1. 第一颗星:取出所有等值谓词中的列，作为索引开头的最开始的列
	2. 第二颗星:添加排序列到索引的列中
	3. 第三颗星:将查询语句剩余的列全部加入到索引中

	其他要求:
	1. 最后添加范围查询列
	2. 每个星级添加的列按照索引区分度由高到低排序

	注意:
	1. 若索引列数达到索引列数阈值，依次舍弃第三颗星和第二颗星
	2. 不支持根据索引区分度阈值舍弃等值列和排序列
	3. 给出的三星索引列数不在这里限制 在外层有限制
*/
func (o *Optimizer) GiveThreeStarAdvice(ast SelectAST) (columns []string, err error) {
	// 排序后的等值谓词中的列
	equalColumnInWhere := ast.EqualPredicateColumnsInWhere()
	// 排序后的排序列
	columnInOrderBy := ast.ColumnsInOrderBy()
	columns = append(equalColumnInWhere, columnInOrderBy...)
	// 排序后的SELECT中所有列
	columnInProjection := ast.ColumnsInProjection()
	columns = append(columns, columnInProjection...)
	// 排序后的范围查询列
	unequalColumnInWhere := ast.UnequalPredicateColumnsInWhere()
	if len(unequalColumnInWhere) > 0 {
		columns = append(columns, unequalColumnInWhere[0])
	}

	columns = utils.RemoveDuplicate(columns)
	return columns, nil
}

type optimizerOption func(*Optimizer)

func (opt optimizerOption) apply(o *Optimizer) {
	opt(o)
}
