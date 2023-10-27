package index

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/log"
	indexoptimizer "github.com/actiontech/sqle/sqle/pkg/optimizer/index"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CanOptimize do some pre check on node.
func CanOptimize(l *logrus.Entry, ctx *session.Context, node ast.Node) bool {
	canNotOptimizeWarnf := "can not optimize node: %v, reason: %v"

	ss, ok := node.(*ast.SelectStmt)
	if !ok {
		l.Warnf(canNotOptimizeWarnf, node, "not select statement")
		return false
	}

	if ss.From == nil {
		l.Warnf(canNotOptimizeWarnf, node, "no from clause")
		return false
	}

	tne := util.TableNameExtractor{TableNames: map[string]*ast.TableName{}}
	ss.Accept(&tne)
	for name, ast := range tne.TableNames {
		exist, err := ctx.IsTableExistInDatabase(ast)
		if err != nil {
			l.Warnf(canNotOptimizeWarnf, node, err)
			return false
		}
		if !exist {
			l.Warnf(canNotOptimizeWarnf, node, fmt.Sprintf("table %s not exist", name))
			return false
		}
	}

	return true
}

const (
	defaultCalculateCardinalityMaxRow = 1000000
	defaultCompositeIndexMaxColumn    = 3
)

type Optimizer struct {
	*session.Context

	l *logrus.Entry

	// tables key is table name, use to match in execution plan.
	tables map[string] /*table name*/ *tableInSelect

	// optimizer options:
	calculateCardinalityMaxRow int
	compositeIndexMaxColumn    int
	createIndexStatement       func(string, ...string) string
}

func NewOptimizer(log *logrus.Entry, ctx *session.Context, opts ...optimizerOption) *Optimizer {
	log = log.WithField("optimizer", "index")

	optimizer := &Optimizer{
		Context:                    ctx,
		l:                          log,
		tables:                     make(map[string]*tableInSelect),
		createIndexStatement:       defaultCreateIndexStatement,
		compositeIndexMaxColumn:    defaultCompositeIndexMaxColumn,
		calculateCardinalityMaxRow: defaultCalculateCardinalityMaxRow,
	}

	for _, opt := range opts {
		opt.apply(optimizer)
	}

	return optimizer
}

type OptimizeResult struct {
	TableName      string
	IndexedColumns []string

	Reason string
}

// tableInSelect store the information of a table in select statement for later optimize.
// 1. when we find a table in single table select statement, we will store the select statement.
// 2. when we find a table in join statement, we will store the join on condition.
type tableInSelect struct {
	joinOnColumn   map[string]bool
	singleTableSel *ast.SelectStmt
}

// Optimize give index advice for the select statement.
func (o *Optimizer) Optimize(ctx context.Context, selectStmt *ast.SelectStmt) ([]*OptimizeResult, error) {
	// select 1; ...
	if selectStmt.From == nil {
		return nil, nil
	}

	o.parseTableFromSelectStmt(selectStmt)

	restoredSQL, err := restoreSelectStmt(selectStmt)
	if err != nil {
		return nil, err
	}

	executionPlan, err := o.GetExecutionPlan(restoredSQL)
	if err != nil {
		// TODO: explain will executed failure, if SQL is collect from MyBatis, it not executable SQL.
		log.NewEntry().Errorf("get execution plan failed, sql: %v, error: %v", restoredSQL, err)
		return nil, nil
	}

	executionPlan = removeDrivingTable(executionPlan)

	var needOptimizedTables []string
	for _, record := range executionPlan {
		if o.needOptimize(record) {
			needOptimizedTables = append(needOptimizedTables, record.Table)
		}
	}

	if len(needOptimizedTables) == 0 {
		return nil, nil
	}

	o.l.Infof("need optimize tables: %v", needOptimizedTables)

	var results []*OptimizeResult
	for _, tableName := range needOptimizedTables {
		table, ok := o.tables[tableName]
		if !ok {
			// given SQL: select * from t1 join t2, there is no join on condition,
			continue
		}

		var result *OptimizeResult
		if len(table.joinOnColumn) > 0 {
			result = o.optimizeJoinTable(tableName)
		} else {
			result, err = o.optimizeSingleTable(ctx, tableName, table.singleTableSel)
			if err != nil {
				return nil, errors.Wrapf(err, "optimize single table %s", tableName)
			}
		}
		if result != nil {
			results = append(results, result)
		}
	}

	return results, nil
}

func (o *Optimizer) parseTableFromSelectNode(stmt *ast.SelectStmt) {
	if stmt.From == nil {
		return
	}
	joinNode := stmt.From.TableRefs
	if util.DoesNotJoinTables(joinNode) {
		// cache single table
		if joinNode.Left == nil {
			return
		}
		if table, ok := joinNode.Left.(*ast.TableSource); ok {
			// var name string := table.AsName.O
			if tableName, ok := table.Source.(*ast.TableName); ok {
				if tableName.Name.O != "" {
					tableInSelect := o.getTableInSelect(tableName.Name.O)
					tableInSelect.singleTableSel = stmt
				}
			}
			if table.AsName.O != "" {
				tableInSelect := o.getTableInSelect(table.AsName.O)
				tableInSelect.singleTableSel = stmt
			}
		}
	}
	o.parseTableFromJoinNode(joinNode)
}

func (o *Optimizer) parseTableFromJoinNode(joinNode *ast.Join) {
	// 深度遍历左子树类型为ast.Join的节点 一旦有节点是JOIN两表的节点，并且没有连接条件，则返回
	if leftNode, ok := joinNode.Left.(*ast.Join); ok {
		o.parseTableFromJoinNode(leftNode)
	}

	if util.IsJoinConditionInOnClause(joinNode) {
		columnNames := util.GetJoinedColumnNameExprInOnClause(joinNode)
		for _, columnName := range columnNames {
			if columnName.Name.Table.O == "" {
				/*
					unsupport sqls like
					ON (column_1 = column_2) should check column belongs to which table
				*/
				continue
			} else {
				/*
					support sqls like
					ON table_1.column_1=table_2.column_1
					ON t1.column_1 = COALESCE(a.c1, b.c1)
					ON (t1.column_1,t1.column_2 = t2.column_1,t2.column_2)
					ON table_1.column_1=table_2.column_1 AND table_1.column_2=table_2.column_2
				*/
				tableInSelect := o.getTableInSelect(columnName.Name.Table.O)
				tableInSelect.joinOnColumn[columnName.Name.Name.O] = true
			}
		}

	}
	if util.IsJoinConditionInUsingClause(joinNode) {
		left, right := util.GetJoinedTableName(joinNode)
		if left == nil || right == nil {
			return
		}
		leftInSelect := o.getTableInSelect(left.Name.O)
		rightInSelect := o.getTableInSelect(right.Name.O)
		for _, columnInUsing := range joinNode.Using {
			leftInSelect.joinOnColumn[columnInUsing.Name.O] = true
			rightInSelect.joinOnColumn[columnInUsing.Name.O] = true
		}
	}

}

func (o *Optimizer) getTableInSelect(tableName string) *tableInSelect {
	if o.tables[tableName] == nil {
		// if not exist, initialize it
		o.tables[tableName] = &tableInSelect{joinOnColumn: make(map[string]bool)}
	}
	return o.tables[tableName]
}

/*
traverse from select stmt and extract tables in it
1. if node is not a join node that join two tables, cache table into tableInSelect.singleTableSel do not fill the joinOnColumn.
2. if node is a join node that join two tables(using ON condition or Using condition) , cache join condition and only fill tableInSelect.joinOnColumn
*/
func (o *Optimizer) parseTableFromSelectStmt(selectStmt *ast.SelectStmt) {
	visitor := util.SelectStmtExtractor{}
	selectStmt.Accept(&visitor)

	for _, stmt := range visitor.SelectStmts {
		o.parseTableFromSelectNode(stmt)
	}
}

func (o *Optimizer) optimizeSingleTable(ctx context.Context, tbl string, ss *ast.SelectStmt) (*OptimizeResult, error) {
	var (
		err            error
		optimizeResult *OptimizeResult
	)

	optimizeResult, err = o.doSpecifiedOptimization(ctx, ss)
	if err != nil {
		return nil, err
	}

	if optimizeResult == nil {
		optimizeResult, err = o.doGeneralOptimization(ctx, ss)
		if err != nil {
			return nil, err
		}
	}

	if optimizeResult == nil {
		return nil, nil
	}

	if len(optimizeResult.IndexedColumns) > o.compositeIndexMaxColumn {
		optimizeResult.IndexedColumns = optimizeResult.IndexedColumns[:o.compositeIndexMaxColumn]
	}

	needIndex, err := o.needIndex(optimizeResult.TableName, optimizeResult.IndexedColumns...)
	if err != nil {
		return nil, err
	}

	if !needIndex {
		return nil, nil
	}

	o.l.Infof("table:%s, indexed columns:%v, reason:%s", optimizeResult.TableName, optimizeResult.IndexedColumns, optimizeResult.Reason)

	if len(optimizeResult.IndexedColumns) > 1 {
		tableNameFromAST, err := extractTableNameFromAST(ss, tbl)
		if err != nil {
			return nil, errors.Wrap(err, "extract table name from AST")
		}

		rowCount, err := o.GetTableRowCount(tableNameFromAST)
		if err != nil {
			return nil, errors.Wrap(err, "get table row count when optimize")
		}

		if rowCount < o.calculateCardinalityMaxRow {
			optimizeResult.IndexedColumns, err = o.sortColumnsByCardinality(tableNameFromAST, optimizeResult.IndexedColumns)
			if err != nil {
				return nil, err
			}
		}
	}

	return optimizeResult, nil
}

func (o *Optimizer) optimizeJoinTable(tbl string) *OptimizeResult {
	table := o.getTableInSelect(tbl)
	indexColumns := make([]string, 0, len(table.joinOnColumn))
	for columnName := range table.joinOnColumn {
		indexColumns = append(indexColumns, columnName)
	}
	return &OptimizeResult{
		TableName:      tbl,
		IndexedColumns: indexColumns,
		Reason:         fmt.Sprintf("字段 %s 为被驱动表 %s 上的关联字段", indexColumns, tbl),
	}
}

// doSpecifiedOptimization optimize single table select.
func (o *Optimizer) doSpecifiedOptimization(ctx context.Context, selectStmt *ast.SelectStmt) (*OptimizeResult, error) {
	if where := selectStmt.Where; where != nil {
		if boe, ok := where.(*ast.BinaryOperationExpr); ok {
			// check function in select stmt
			if fce, ok := boe.L.(*ast.FuncCallExpr); ok {
				result, err := o.optimizeOnFunctionCallExpression(getTableNameFromSingleSelect(selectStmt), fce)
				if err != nil {
					return nil, err
				}
				if result != nil {
					return result, nil
				}
			}
		}

		// check where like 'mike%'
		if ple, ok := where.(*ast.PatternLikeExpr); ok {
			if cne, ok := ple.Expr.(*ast.ColumnNameExpr); ok {
				if ve, ok := ple.Pattern.(*driver.ValueExpr); ok {
					datum := ve.Datum.GetString()
					if !strings.HasPrefix(datum, "%") &&
						!strings.HasPrefix(datum, "_") {
						return &OptimizeResult{
							TableName:      getTableNameFromSingleSelect(selectStmt),
							IndexedColumns: []string{cne.Name.Name.L},
							Reason:         "为前缀模式匹配添加前缀索引",
						}, nil
					}
				}
			}
		}
	}

	if selectStmt.Where == nil {
		var cols []string
		for _, field := range selectStmt.Fields.Fields {
			if field.Expr != nil {
				afe, ok := field.Expr.(*ast.AggregateFuncExpr)
				if !ok {
					continue
				}
				if afe.F == ast.AggFuncMin ||
					afe.F == ast.AggFuncMax {
					cne, ok := afe.Args[0].(*ast.ColumnNameExpr)
					if ok {
						cols = append(cols, cne.Name.Name.L)
					}
				}
			}
		}
		if len(cols) > 0 {
			return &OptimizeResult{
				TableName:      getTableNameFromSingleSelect(selectStmt),
				IndexedColumns: cols,
				Reason:         "利用索引有序的性质快速找到最值",
			}, nil
		}
	}

	return nil, nil
}

func (o *Optimizer) optimizeOnFunctionCallExpression(tbl string, fce *ast.FuncCallExpr) (*OptimizeResult, error) {
	var cols []string
	for _, arg := range fce.Args {
		if cne, ok := arg.(*ast.ColumnNameExpr); ok {
			cols = append(cols, cne.Name.Name.L)
		}
	}
	if len(cols) == 0 {
		return nil, nil
	}

	var buf strings.Builder
	if err := fce.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)); err != nil {
		return nil, errors.Wrap(err, "restore func call expr when do specified optimization")
	}

	versionWithFlavor, err := o.GetSystemVariable("version")
	if err != nil {
		return nil, errors.Wrap(err, "get version when do specified optimization")
	}

	curVersion, err := semver.NewVersion(versionWithFlavor)
	if err != nil {
		return nil, errors.Wrap(err, "parse version when do specified optimization")
	}
	if curVersion.LessThan(semver.MustParse("5.7.0")) {
		return nil, nil
	}
	if curVersion.LessThan(semver.MustParse("8.0.13")) {
		return &OptimizeResult{
			TableName:      tbl,
			IndexedColumns: []string{buf.String()},
			Reason:         "MySQL5.7以上版本需要在虚拟列上创建索引",
		}, nil
	}

	return &OptimizeResult{
		TableName:      tbl,
		IndexedColumns: []string{buf.String()},
		Reason:         "MySQL8.0.13以上版本支持直接创建函数索引",
	}, nil
}

// doGeneralOptimization optimize single table select.
func (o *Optimizer) doGeneralOptimization(ctx context.Context, selectStmt *ast.SelectStmt) (*OptimizeResult, error) {
	generalOptimizer := indexoptimizer.NewOptimizer(o.Context)

	restoredSQL, err := restoreSelectStmt(selectStmt)
	if err != nil {
		return nil, err
	}

	sa, err := newSelectAST(restoredSQL)
	if err != nil {
		return nil, err
	}

	indexedColumns, err := generalOptimizer.Optimize(sa)
	if err != nil {
		return nil, err
	}

	if len(indexedColumns) == 0 {
		return nil, nil
	}

	o.l.Infof("general optimize result: %v(index columns)", indexedColumns)

	return &OptimizeResult{
		TableName:      getTableNameFromSingleSelect(selectStmt),
		IndexedColumns: indexedColumns,
		Reason:         "三星索引建议",
	}, nil
}

type cardinality struct {
	columnName  string
	cardinality int
}

type cardinalities []cardinality

func (c cardinalities) Len() int {
	return len(c)
}

func (c cardinalities) Less(i, j int) bool {
	return c[i].cardinality > c[j].cardinality
}

func (c cardinalities) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (o *Optimizer) sortColumnsByCardinality(tn *ast.TableName, indexedColumns []string) (sortedColumns []string, err error) {
	if tn == nil {
		return nil, errors.New("table ast not found when sort columns by cardinality")
	}

	cardinalitySlice := make(cardinalities, len(indexedColumns))
	for i, column := range indexedColumns {
		c, err := o.GetColumnCardinality(tn, column)
		if err != nil {
			return nil, errors.Wrap(err, "get column cardinality")
		}
		cardinalitySlice[i] = cardinality{
			columnName:  column,
			cardinality: c,
		}
	}

	o.l.Debugf("table %s column cardinalities(before sort): %+v", tn.Name, cardinalitySlice)
	sort.Sort(cardinalitySlice)
	o.l.Debugf("table %s column cardinalities(after sort): %+v", tn.Name, cardinalitySlice)

	for _, c := range cardinalitySlice {
		sortedColumns = append(sortedColumns, c.columnName)
	}

	return sortedColumns, nil
}

// needOptimize check table need optimize index of table or not.
//
// Optimize means that:
// 1. When SQL do not use index, we can create index for the select statement.
// 2. When SQL use index, but the index is not suitable, we should optimize the index.
//
// We do it by check MySQL execution plan's access_type field.
// ref: https://dev.mysql.com/doc/refman/5.7/en/explain-output.html#explain-join-types
func (o *Optimizer) needOptimize(record *executor.ExplainRecord) bool {

	// Full table scan: select * from t1 where common_column = 'a'
	// This SQL will scan all rows of table t1.
	if record.Type == executor.ExplainRecordAccessTypeAll {
		return true
	}

	// Index-only scan: select key_part2 from t1 where key_part3 = 'a'
	// This SQL will scan all rows of index idx_composite. It's a little better than previous case.
	if record.Type == executor.ExplainRecordAccessTypeIndex {
		return true
	}

	return false
}

// needIndex check need add index on tbl.columns or not.
func (o *Optimizer) needIndex(tbl string, columns ...string) (bool, error) {
	table, ok := o.tables[tbl]
	if !ok {
		return false, fmt.Errorf("table %s not found when check index", tbl)
	}

	if table.singleTableSel == nil {
		return false, fmt.Errorf("table %s do not have select statement when check index", tbl)
	}

	tableNameFromAST, err := extractTableNameFromAST(table.singleTableSel, tbl)
	if err != nil {
		return false, fmt.Errorf("extract table name from AST failed when check index: %v", err)
	}

	cts, exist, err := o.GetCreateTableStmt(tableNameFromAST)
	if err != nil {
		return false, errors.Wrap(err, "get create table statement when check index")
	}
	if !exist {
		return false, fmt.Errorf("table %s not found on session context when check index", tbl)
	}

	for _, index := range util.ExtractIndexFromCreateTableStmt(cts) {
		if reflect.DeepEqual(index, columns) {
			return false, nil
		}
		if strings.HasPrefix(strings.Join(index, ","), strings.Join(columns, ",")) {
			return false, nil
		}
	}
	return true, nil
}

type optimizerOption func(*Optimizer)

func (oo optimizerOption) apply(o *Optimizer) {
	oo(o)
}

func WithCalculateCardinalityMaxRow(row int) optimizerOption {
	return func(o *Optimizer) {
		o.calculateCardinalityMaxRow = row
	}
}

func WithCompositeIndexMaxColumn(column int) optimizerOption {
	return func(o *Optimizer) {
		o.compositeIndexMaxColumn = column
	}
}

func WithCreateIndexStatement(f func(tableName string, columns ...string) string) optimizerOption {
	return func(o *Optimizer) {
		o.createIndexStatement = f
	}
}

func defaultCreateIndexStatement(tableName string, columns ...string) string {
	indexName := fmt.Sprintf("idx_%s_%s", tableName, strings.Join(columns, "_"))

	return fmt.Sprintf("CREATE INDEX %s ON %s (%s)",
		indexName,
		tableName,
		strings.Join(columns, ", "))
}

func restoreSelectStmt(ss *ast.SelectStmt) (string, error) {
	var buf strings.Builder
	rc := format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)

	if err := ss.Restore(rc); err != nil {
		return "", errors.Wrap(err, "restore select statement")
	}

	return buf.String(), nil
}

func extractTableNameFromAST(ss *ast.SelectStmt, tbl string) (*ast.TableName, error) {
	sourceExtractor := util.TableSourceExtractor{TableSources: map[string]*ast.TableSource{}}
	ss.Accept(&sourceExtractor)

	for _, ts := range sourceExtractor.TableSources {
		tableName, ok := (ts.Source).(*ast.TableName)
		if !ok {
			return nil, errors.New("tableName not found")
		}

		if tableName.Name.O == tbl || tbl == ts.AsName.O {
			return tableName, nil
		}
	}

	return nil, fmt.Errorf("table %s not found in select statement", tbl)
}

func getTableNameFromSingleSelect(ss *ast.SelectStmt) string {
	if ss.From.TableRefs.Left == nil {
		return ""
	}
	return ss.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.O
}

// removeDrivingTable remove driving table from execution plan.
//
// Index is not silver bullet, we only give advice on driven table.
// Such as : select * from t1, t2 where t1.id = t2.id;
// There are two records in execution plan, the first one is driving table, the second one is driven table.
func removeDrivingTable(records []*executor.ExplainRecord) []*executor.ExplainRecord {
	var result []*executor.ExplainRecord

	if len(records) == 0 || len(records) == 1 {
		return records
	}

	i, j := 0, 1
	for j < len(records) {
		if records[i].Id == records[j].Id {
			result = append(result, records[j])
		} else {
			i = j
		}
		j++
	}

	return result
}
