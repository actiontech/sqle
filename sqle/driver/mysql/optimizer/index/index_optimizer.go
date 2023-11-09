package index

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/sirupsen/logrus"
)

func restore(node ast.Node) (sql string) {
	var buf strings.Builder
	rc := format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)

	if err := node.Restore(rc); err != nil {
		return
	}
	sql = buf.String()
	return
}

func NewIndexOptimizer(log *logrus.Entry, ctx *session.Context, node ast.Node, opts ...option) *IndexOptimizer {
	log = log.WithField("optimizer", "index")

	optimizer := &IndexOptimizer{
		OriginNode:                 node,
		SqlContext:                 ctx,
		Log:                        log,
		compositeIndexMaxColumn:    defaultCompositeIndexMaxColumn,
		calculateCardinalityMaxRow: defaultCalculateCardinalityMaxRow,
		Results:                    make([]*OptimizeResult, 0),
	}

	for _, opt := range opts {
		opt.apply(optimizer)
	}

	return optimizer
}

type option func(*IndexOptimizer)

func (oo option) apply(o *IndexOptimizer) {
	oo(o)
}
func WithMaxRow(row int) option {
	return func(o *IndexOptimizer) {
		o.calculateCardinalityMaxRow = row
	}
}

func WithIndexMaxColumn(column int) option {
	return func(o *IndexOptimizer) {
		o.compositeIndexMaxColumn = column
	}
}

func WithStatement(f func(tableName string, columns ...string) string) optimizerOption {
	return func(o *Optimizer) {
		o.createIndexStatement = f
	}
}

type IndexOptimizer struct {
	SqlContext       *session.Context
	Log              *logrus.Entry
	CurrentTableName string
	OriginNode       ast.Node
	Results          []*OptimizeResult
	// optimizer options:
	calculateCardinalityMaxRow int
	compositeIndexMaxColumn    int
}

func (opt *IndexOptimizer) Optimize() {
	if !CanOptimize(opt.Log, opt.SqlContext, opt.OriginNode) {
		return
	}

	tables := opt.getTableShouldBeOptimize()
	if len(tables) == 0 {
		return
	}

	for _, tableName := range tables {
		opt.CurrentTableName = tableName
		opt.OriginNode.Accept(opt)
	}
}

func (opt *IndexOptimizer) getTableShouldBeOptimize() []string {

	executionPlans, err := opt.SqlContext.GetExecutionPlan(opt.OriginNode.Text())
	if err != nil {
		// explain will executed failure, if SQL is collect from MyBatis, it not executable SQL.
		opt.Log.Errorf("get execution plan failed, sql: %v, error: %v", opt.OriginNode.Text(), err)
		return nil
	}
	executionPlans = removeDrivingTable(executionPlans)
	var needOptimizedTables []string
	for _, record := range executionPlans {
		if record.Type == executor.ExplainRecordAccessTypeAll || record.Type == executor.ExplainRecordAccessTypeIndex {
			needOptimizedTables = append(needOptimizedTables, record.Table)
		}
	}
	if len(needOptimizedTables) == 0 {
		return nil
	}

	return needOptimizedTables
}

func (opt *IndexOptimizer) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	if advisor := opt.NewCreateIndexAdvisor(in); advisor != nil {
		advise := advisor.GiveAdvice()
		if advise != nil {
			opt.Results = append(opt.Results, advise)
		}
	}
	return in, false
}

func (v *IndexOptimizer) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (opt *IndexOptimizer) NewCreateIndexAdvisor(in ast.Node) CreateIndexAdvisor {
	switch stmt := in.(type) {
	case *ast.Join:
		return JoinIndexingAdvisor{
			SqlContext: opt.SqlContext,
			node:       stmt,
			tableName:  opt.CurrentTableName,
		}
	case *ast.FuncCallExpr:
		return FunctionIndexingAdvisor{
			SqlContext: opt.SqlContext,
			node:       stmt,
			tableName:  opt.CurrentTableName,
		}
	case *ast.PatternLikeExpr:
		return PrefixIndexingAdvisor{
			SqlContext: opt.SqlContext,
			node:       stmt,
			tableName:  opt.CurrentTableName,
		}
	case *ast.AggregateFuncExpr:
		return AggregateIndexingAdvisor{
			SqlContext: opt.SqlContext,
			node:       stmt,
			tableName:  opt.CurrentTableName,
		}
	case *ast.SelectStmt:
		if in == opt.OriginNode {
			// 三星索引建议只对SQL原语句对应的节点生效
			return ThreeStarAdvisor{
				SqlContext: opt.SqlContext,
				node:       stmt,
				tableName:  opt.CurrentTableName,
			}
		}
	}
	return nil
}

type CreateIndexAdvisor interface {
	GiveAdvice() *OptimizeResult
}
type ThreeStarAdvisor struct {
	SqlContext *session.Context
	node       *ast.SelectStmt
	tableName  string
}

// TODO 实现具体逻辑
func (a ThreeStarAdvisor) GiveAdvice() *OptimizeResult {
	var columns []string
	return &OptimizeResult{
		TableName:      a.tableName,
		IndexedColumns: columns,
		Reason:         "三星索引建议",
	}
}

type JoinIndexingAdvisor struct {
	SqlContext *session.Context
	node       *ast.Join
	tableName  string
}

func (a JoinIndexingAdvisor) GiveAdvice() *OptimizeResult {
	var indexColumn []string
	// TODO 找到JOIN Right上的Table 为被驱动表
	// TODO 处理多表JOIN
	var tableName string
	table, ok := a.node.Right.(*ast.TableSource)
	if ok {
		if table.AsName.L != "" {
			tableName = table.AsName.L
		} else if tb, ok := table.Source.(*ast.TableName); ok {
			tableName = tb.Name.L
		}
	}
	if a.node.On != nil {
		bo, ok := a.node.On.Expr.(*ast.BinaryOperationExpr)
		if ok {
			leftName, ok := bo.L.(*ast.ColumnNameExpr)
			if ok && leftName.Name.Table.L != tableName {
				indexColumn = append(indexColumn, leftName.Name.String())
			}
			rightName, ok := bo.R.(*ast.ColumnNameExpr)
			if ok && rightName.Name.Table.L != tableName {
				indexColumn = append(indexColumn, rightName.Name.String())
			}
		}
	}
	if len(a.node.Using) > 0 {
		// https://dev.mysql.com/doc/refman/8.0/en/join.html
		// TODO 根据SQLContext中的表结构，确定USING的连接的是哪两张表
		for _, column := range a.node.Using {
			indexColumn = append(indexColumn, column.Name.L)
		}
	}
	if len(indexColumn) == 0 {
		return nil
	}
	return &OptimizeResult{
		TableName:      a.tableName,
		IndexedColumns: indexColumn,
		Reason:         fmt.Sprintf("字段 %s 为被驱动表 %s 上的关联字段", indexColumn, tableName),
	}
}

type FunctionIndexingAdvisor struct {
	SqlContext *session.Context
	node       *ast.FuncCallExpr
	tableName  string
}

func (a FunctionIndexingAdvisor) GiveAdvice() *OptimizeResult {

	// versionWithFlavor, err := a.SqlContext.GetSystemVariable("version")
	// if err != nil {
	// 	// LOG HERE
	// 	return nil
	// }
	versionWithFlavor := "5.7.1"
	curVersion, err := semver.NewVersion(versionWithFlavor)
	if err != nil {
		return nil
	}
	if curVersion.LessThan(semver.MustParse("5.7.0")) {
		return nil
	}

	columnNameVisitor := util.ColumnNameVisitor{}
	a.node.Accept(&columnNameVisitor)
	if len(columnNameVisitor.ColumnNameList) == 0 {
		return nil
	}
	columns := make([]string, 0, len(columnNameVisitor.ColumnNameList))
	for _, columnName := range columnNameVisitor.ColumnNameList {
		columns = append(columns, columnName.Name.Name.L)
	}
	if curVersion.LessThan(semver.MustParse("8.0.13")) {
		return &OptimizeResult{
			TableName:      a.tableName,
			IndexedColumns: columns,
			Reason:         "MySQL5.7以上版本需要在虚拟列上创建索引",
		}
	}

	return &OptimizeResult{
		TableName:      a.tableName,
		IndexedColumns: columns,
		Reason:         "MySQL8.0.13以上版本支持直接创建函数索引",
	}
}

type AggregateIndexingAdvisor struct {
	SqlContext *session.Context
	node       *ast.AggregateFuncExpr
	tableName  string
}

// https://dev.mysql.com/doc/refman/8.0/en/aggregate-functions.html#function_max
// https://dev.mysql.com/doc/refman/8.0/en/aggregate-functions.html#function_min
func (a AggregateIndexingAdvisor) GiveAdvice() *OptimizeResult {
	var indexColumns []string
	if strings.ToLower(a.node.F) == ast.AggFuncMin || strings.ToLower(a.node.F) == ast.AggFuncMax {
		column, ok := a.node.Args[0].(*ast.ColumnNameExpr)
		if ok {
			indexColumns = append(indexColumns, column.Name.Name.L)
		}
	} else {
		return nil
	}
	return &OptimizeResult{
		TableName:      a.tableName,
		IndexedColumns: indexColumns,
		Reason:         fmt.Sprintf("索引建议 | 对于SQL:%s 可以利用索引有序的性质快速找到最值", restore(a.node)),
	}
}

type PrefixIndexingAdvisor struct {
	SqlContext *session.Context
	node       *ast.PatternLikeExpr
	tableName  string
}

func (a PrefixIndexingAdvisor) GiveAdvice() *OptimizeResult {
	if !util.CheckWhereFuzzySearch(a.node) {
		return nil
	}
	column, ok := a.node.Expr.(*ast.ColumnNameExpr)
	if !ok {
		return nil
	}
	sql := restore(a.node)
	return &OptimizeResult{
		TableName:      a.tableName,
		IndexedColumns: []string{column.Name.Name.L},
		Reason:         fmt.Sprintf("索引建议 | SQL使用了前缀模式匹配：%s ", sql),
	}
}
