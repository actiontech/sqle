package mysql

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
	parser_driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/sirupsen/logrus"
)

const (
	MAX_INDEX_COLUMN               string = "composite_index_max_column"
	MAX_INDEX_COLUMN_DEFAULT_VALUE int    = 5
	threeStarIndexAdviceFormat     string = "索引建议 | 根据三星索引设计规范，建议对表%s添加%s索引：【%s】"
	prefixIndexAdviceFormat        string = "索引建议 | SQL使用了前缀模式匹配，数据量大时，可建立翻转函数索引"
	extremalIndexAdviceFormat      string = "索引建议 | SQL使用了最值函数，可以利用索引有序的性质快速找到最值，建议对表%s添加单列索引，参考列：%s"
	functionIndexAdviceFormatV80   string = "索引建议 | SQL使用了函数作为查询条件，在MySQL8.0.13以上的版本，可以创建函数索引，建议对表%s添加函数索引，参考列：%s"
	functionIndexAdviceFormatV57   string = "索引建议 | SQL使用了函数作为查询条件，在MySQL5.7以上的版本，可以在虚拟列上创建索引，建议对表%s添加虚拟列索引，参考列：%s"
	functionIndexAdviceFormatAll   string = "索引建议 | SQL使用了函数作为查询条件，在MySQL5.7以上的版本，可以在虚拟列上创建索引，在MySQL8.0.13以上的版本，可以创建函数索引，建议根据MySQL版本对表%s添加合适的索引，参考列：%s"
	joinIndexAdviceFormat          string = "索引建议 | SQL中字段%s为被驱动表%s上的关联字段，建议对表%s添加单列索引，参考列：%s"
)

type OptimizeResult struct {
	TableName      string
	IndexedColumns []string
	Reason         string
}

func optimize(log *logrus.Entry, ctx *session.Context, node ast.Node, params params.Params) []*OptimizeResult {
	if !canOptimize(log, ctx, node) {
		return nil
	}

	log = log.WithField("optimizer", "index")

	var optimizeResult []*OptimizeResult
	for _, meta := range AdvisorMetaList {
		optimizeResult = append(
			optimizeResult,
			meta.newFunction(ctx, log, node, params).GiveAdvices()...,
		)
	}
	return optimizeResult
}

func canOptimize(log *logrus.Entry, ctx *session.Context, node ast.Node) bool {
	canNotOptimizeWarnf := "can not optimize node: %v, reason: %v"
	if ctx == nil {
		return false
	}
	if node == nil {
		return false
	}
	selectStmt, ok := node.(*ast.SelectStmt)
	if !ok {
		log.Warnf(canNotOptimizeWarnf, node, "not select statement")
		return false
	}

	if selectStmt.From == nil {
		log.Warnf(canNotOptimizeWarnf, node, "no from clause")
		return false
	}

	extractor := util.TableNameExtractor{TableNames: map[string]*ast.TableName{}}
	selectStmt.Accept(&extractor)
	for name, ast := range extractor.TableNames {
		exist, err := ctx.IsTableExistInDatabase(ast)
		if err != nil {
			log.Warnf(canNotOptimizeWarnf, node, err)
			return false
		}
		if !exist {
			log.Warnf(canNotOptimizeWarnf, node, fmt.Sprintf("table %s not exist", name))
			return false
		}
	}
	executionPlans, err := ctx.GetExecutionPlan(node.Text())
	if err != nil {
		log.Errorf("get execution plan failed, sql: %v, error: %v", node.Text(), err)
		return false
	}
	for _, record := range executionPlans {
		if record.Type == executor.ExplainRecordAccessTypeAll || record.Type == executor.ExplainRecordAccessTypeIndex {
			return true
		}
	}
	return false
}

// CreateIndexAdvisor 基于SQL语句、SQL上下文、库表信息等生成创建索引的建议，在给出建议前需要指明优化建议针对的节点
type AdvisorMeta struct {
	advisorName string
	newFunction func(ctx *session.Context, log *logrus.Entry, originNode ast.Node, params params.Params) CreateIndexAdvisor
}

var AdvisorMetaList []AdvisorMeta = []AdvisorMeta{
	{
		advisorName: "prefix_index_advisor",
		newFunction: newPrefixIndexAdvisor,
	},
	{
		advisorName: "join_index_advisor",
		newFunction: newJoinIndexAdvisor,
	},
	{
		advisorName: "extremal_index_advisor",
		newFunction: newExtremalIndexAdvisor,
	},
	{
		advisorName: "function_index_advisor",
		newFunction: newFunctionIndexAdvisor,
	},
	{
		advisorName: "three_star_index_advisor",
		newFunction: newThreeStarIndexAdvisor,
	},
}

// 还原抽象语法树节点至SQL
func restore(node ast.Node) (sql string) {
	var buf strings.Builder
	rc := format.NewRestoreCtx(format.DefaultRestoreFlags, &buf)

	if err := node.Restore(rc); err != nil {
		return
	}
	sql = buf.String()
	return
}

func getDrivingTableInfo(originNode ast.Node, sqlContext *session.Context) (*ast.TableSource, *ast.CreateTableStmt, error) {
	executionPlans, err := sqlContext.GetExecutionPlan(originNode.Text())
	if err != nil {
		return nil, nil, err
	}
	var tableSource *ast.TableSource
	var createTable *ast.CreateTableStmt
	extractor := util.TableSourceExtractor{TableSources: map[string]*ast.TableSource{}}
	originNode.Accept(&extractor)
	var ok bool
	if len(executionPlans) > 0 {
		tableSource, ok = extractor.TableSources[strings.ToLower(executionPlans[0].Table)]
		if !ok {
			return nil, nil, fmt.Errorf("get driving table source failed")
		}
	}
	tableName, ok := tableSource.Source.(*ast.TableName)
	if !ok {
		return nil, nil, fmt.Errorf("driving tableSource.Source is not ast.TableName")
	}
	createTable, ok, err = sqlContext.GetCreateTableStmt(tableName)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, fmt.Errorf("driving table CreateTableStmt is not exist")
	}
	return tableSource, createTable, nil
}

type columnMap map[string] /*column name or alias name*/ struct{}

func (c columnMap) add(col *ast.ColumnNameExpr) {
	c[col.Name.Name.L] = struct{}{}
}

func (c columnMap) delete(col *ast.ColumnNameExpr) {
	delete(c, col.Name.Name.L)
}

type columnWithSelectivity struct {
	columnName  *ast.ColumnNameExpr
	selectivity float64
}

type columnsWithSelectivity []columnWithSelectivity

type columnGroup struct {
	columnsWithSelectivity
	columnMap
}

func newColumnGroup() columnGroup {
	return columnGroup{
		columnsWithSelectivity: make(columnsWithSelectivity, 0),
		columnMap:              make(columnMap),
	}
}

func (g columnGroup) len() int {
	return len(g.columnMap)
}

func (g *columnGroup) sort() {
	sort.Sort(g.columnsWithSelectivity)
}

func (g columnGroup) has(col *ast.ColumnNameExpr) bool {
	_, exist := g.columnMap[col.Name.Name.L]
	return exist
}

func (g *columnGroup) add(col *ast.ColumnNameExpr, selectivity float64) {
	if _, exist := g.columnMap[col.Name.Name.L]; exist {
		// 不重复添加
		return
	}
	g.columnMap.add(col)
	g.columnsWithSelectivity = append(g.columnsWithSelectivity, columnWithSelectivity{
		columnName:  col,
		selectivity: selectivity,
	})
}

func (g *columnGroup) replace(newColumn columnWithSelectivity, index int) {
	oldColumn := g.columnsWithSelectivity[index]
	g.columnsWithSelectivity[index] = newColumn
	g.columnMap.delete(oldColumn.columnName)
	g.columnMap.add(newColumn.columnName)
}

func (g *columnGroup) delete(col *ast.ColumnNameExpr) {
	g.columnMap.delete(col)
	var index int
	for index = range g.columnsWithSelectivity {
		if g.columnsWithSelectivity[index].columnName == col {
			break
		}
	}
	g.columnsWithSelectivity = append(g.columnsWithSelectivity[:index], g.columnsWithSelectivity[index:]...)
}

func (g columnGroup) stringSlice() []string {
	indexedColumn := make([]string, 0, len(g.columnMap))
	for _, column := range g.columnsWithSelectivity {
		indexedColumn = append(indexedColumn, column.columnName.Name.Name.L)
	}
	return indexedColumn
}

func (c columnsWithSelectivity) Len() int {
	return len(c)
}

func (c columnsWithSelectivity) Less(i, j int) bool {
	return c[i].selectivity > c[j].selectivity
}

func (c columnsWithSelectivity) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type CreateIndexAdvisor interface {
	GiveAdvices() []*OptimizeResult
}

type threeStarIndexAdvisor struct {
	sqlContext             *session.Context
	log                    *logrus.Entry
	drivingTableSource     *ast.TableSource     // 驱动表的TableSource
	drivingTableColumn     *columnInSelect      // 经过解析的驱动表的列
	drivingTableCreateStmt *ast.CreateTableStmt // 驱动表的建表语句
	originNode             ast.Node             // 原SQL的节点
	maxColumns             int                  // 复合索引列的上限数量
	minSelecivity          float64              // 列选择性阈值下限
	usingOr                bool                 // 使用了或条件
	possibleColumns        columnGroup          // SQL语句中可能作为索引的备选列
	columnLastAdd          columnGroup          // 最后添加的列，例如：非等值列
	columnShouldNotAdd     columnGroup          // 不该添加的列，例如：类型不适合作为索引的列、单列主键列
	adviceColumns          columnGroup
	selectivityMap         map[string]float64
}

type columnInSelect struct {
	equalColumnInWhere   columnGroup
	unequalColumnInWhere columnGroup
	columnInOrderBy      columnGroup
	columnInFieldList    columnGroup
}

func newThreeStarIndexAdvisor(ctx *session.Context, log *logrus.Entry, originNode ast.Node, params params.Params) CreateIndexAdvisor {
	maxColumns := params.GetParam(MAX_INDEX_COLUMN).Int()
	if maxColumns == 0 {
		maxColumns = MAX_INDEX_COLUMN_DEFAULT_VALUE
	}
	return &threeStarIndexAdvisor{
		sqlContext:    ctx,
		log:           log,
		originNode:    originNode,
		maxColumns:    maxColumns,
		minSelecivity: 2,
		drivingTableColumn: &columnInSelect{
			equalColumnInWhere:   newColumnGroup(),
			unequalColumnInWhere: newColumnGroup(),
			columnInOrderBy:      newColumnGroup(),
			columnInFieldList:    newColumnGroup(),
		},
		columnLastAdd:      newColumnGroup(),
		possibleColumns:    newColumnGroup(),
		columnShouldNotAdd: newColumnGroup(),
		adviceColumns:      newColumnGroup(),
	}
}

/*
三星索引建议

	三星索引要求:
	1. 第一颗星:取出所有等值谓词中的列，作为索引开头的最开始的列
	2. 第二颗星:添加排序列到索引的列中
	3. 第三颗星:将查询语句剩余的列全部加入到索引中

	其他要求:
	1. 最后添加范围查询列，并仅添加一列
	2. 每个星级添加的列按照索引区分度由高到低排序
*/
func (a *threeStarIndexAdvisor) GiveAdvices() []*OptimizeResult {
	if a.originNode == nil {
		return nil
	}
	err := a.loadEssentials()
	if err != nil {
		a.log.Logger.Warnf("when three star index advisor load essentials failed, err:%v", err)
		return nil
	}
	err = a.extractColumnInSelect()
	if err != nil {
		a.log.Logger.Warnf("extract column in select failed, sql:%s,err:%v", restore(a.originNode), err)
		return nil
	}

	if !a.canGiveAdvice() {
		return nil
	}

	a.sortColumnBySelectivity()
	a.giveAdvice()
	if a.adviceColumns.len() == 0 {
		return nil
	}
	if util.IsIndex(a.adviceColumns.columnMap, a.drivingTableCreateStmt.Constraints) {
		return nil
	}
	tableName := util.GetTableNameFromTableSource(a.drivingTableSource)
	indexColumns := a.adviceColumns.stringSlice()
	var indexType string = "复合"
	if len(indexColumns) == 1 {
		indexType = "单列"
	}
	return []*OptimizeResult{{
		TableName:      tableName,
		IndexedColumns: indexColumns,
		Reason:         fmt.Sprintf(threeStarIndexAdviceFormat, tableName, indexType, strings.Join(indexColumns, "，")),
	}}
}

func (a *threeStarIndexAdvisor) loadEssentials() (err error) {
	a.drivingTableSource, a.drivingTableCreateStmt, err = getDrivingTableInfo(a.originNode, a.sqlContext)
	if err != nil {
		return err
	}
	columnInTable := make([]string, 0, len(a.drivingTableCreateStmt.Cols))
	for _, column := range a.drivingTableCreateStmt.Cols {
		columnInTable = append(columnInTable, column.Name.Name.L)
	}
	tableName, ok := a.drivingTableSource.Source.(*ast.TableName)
	if !ok {
		return fmt.Errorf("in three star advisor driving tableSource.Source is not ast.TableName,sql %s", restore(a.originNode))
	}
	a.selectivityMap, err = a.sqlContext.GetSelectivityOfColumns(tableName, columnInTable)
	if err != nil {
		return err
	}
	for _, value := range a.selectivityMap {
		if value <= 0 {
			return fmt.Errorf("get selectivity of columns failed, value of selectivity Less than or equal to zero,sql %s", restore(a.originNode))
		}
	}
	return nil
}

/*
获取SELECT语句中:

	1 SELECT中的裸的列
	2 WHERE等值条件中，列=值的筛选列，其中列属于驱动表
	3 WHERE不等值条件中，列(非等)值的范围删选列，其中列属于驱动表
	4 ORDER BY中裸的列，其中列属于驱动表
*/
func (a *threeStarIndexAdvisor) extractColumnInSelect() error {
	selectStmt, ok := a.originNode.(*ast.SelectStmt)
	if !ok {
		return fmt.Errorf("in three star advisor, type of current node is not ast.SelectStmt")
	}
	var selecivity float64
	if selectStmt.Where != nil {
		// 访问Where子句，解析并存储属于驱动表等值列和非等值列
		selectStmt.Where.Accept(a)
		for _, col := range a.drivingTableColumn.equalColumnInWhere.columnsWithSelectivity {
			selecivity = a.selectivityMap[col.columnName.Name.Name.L]
			if selecivity <= a.minSelecivity {
				a.columnShouldNotAdd.add(col.columnName, selecivity)
			} else {
				a.possibleColumns.add(col.columnName, selecivity)
			}
		}
		for _, col := range a.drivingTableColumn.unequalColumnInWhere.columnsWithSelectivity {
			selecivity = a.selectivityMap[col.columnName.Name.Name.L]
			a.columnLastAdd.add(col.columnName, selecivity)
			if selecivity <= a.minSelecivity {
				a.columnShouldNotAdd.add(col.columnName, selecivity)
			} else {
				a.possibleColumns.add(col.columnName, selecivity)
			}
		}
	}
	if selectStmt.OrderBy != nil {
		// 遍历Order By的对象切片，存储其中属于驱动表的裸列
		for _, item := range selectStmt.OrderBy.Items {
			if col, ok := item.Expr.(*ast.ColumnNameExpr); ok {
				if !a.isColumnInDrivingTable(col) {
					continue
				}
				selecivity = a.selectivityMap[col.Name.Name.L]
				a.drivingTableColumn.columnInOrderBy.add(col, selecivity)
				if selecivity <= a.minSelecivity {
					a.columnShouldNotAdd.add(col, selecivity)
				} else {
					a.possibleColumns.add(col, selecivity)
				}
			}
		}
	}
	if selectStmt.Fields != nil {
		// 遍历Select子句，存储其中属于驱动表的裸列
		for _, field := range selectStmt.Fields.Fields {
			// Expr is not nil, WildCard will be nil.
			if field.WildCard != nil && field.WildCard.Table.String() == "" && field.WildCard.Schema.String() == "" {
				// 如果是*则添加所有列
				var col *ast.ColumnNameExpr
				for _, coldef := range a.drivingTableCreateStmt.Cols {
					col = &ast.ColumnNameExpr{Name: coldef.Name}
					selecivity = a.selectivityMap[coldef.Name.Name.L]
					a.drivingTableColumn.columnInFieldList.add(col, selecivity)
					if selecivity <= a.minSelecivity {
						a.columnShouldNotAdd.add(col, selecivity)
					} else {
						a.possibleColumns.add(col, selecivity)
					}
				}
				continue
			}
			if col, ok := field.Expr.(*ast.ColumnNameExpr); ok {
				if !a.isColumnInDrivingTable(col) {
					continue
				}
				selecivity = a.selectivityMap[col.Name.Name.L]
				a.drivingTableColumn.columnInFieldList.add(col, selecivity)
				if selecivity <= a.minSelecivity {
					a.columnShouldNotAdd.add(col, selecivity)
				} else {
					a.possibleColumns.add(col, selecivity)
				}
			}
		}
	}
	if a.drivingTableCreateStmt != nil {
		// 遍历建表语句，加入备选列中不应该添加到索引中的列
		// 1. 遍历索引
		// 若主键在SQL中备选的所有列中，并且SQL的排序会根据主键的排序走，此时主键不添加到索引中
		var primaryColumn *ast.ColumnName // 先只考虑单列主键
		for _, constraint := range a.drivingTableCreateStmt.Constraints {
			if constraint.Tp == ast.ConstraintPrimaryKey {
				primaryColumn = constraint.Keys[0].Column
			}
		}
		// SQL结果的排序根据主键的顺序排序
		if primaryColumn != nil {
			primaryColumnName := primaryColumn.Name.L
			// 当主键不在备选列中时，不考虑该情况
			if !a.possibleColumns.has(&ast.ColumnNameExpr{Name: primaryColumn}) {
				return nil
			}
			var orderByPrimaryKey bool
			// 当没有Order By时，按照主键排序，覆盖索引可以不包含主键
			if a.drivingTableColumn.columnInOrderBy.len() == 0 {
				orderByPrimaryKey = true
			}
			// 当Order By主键的时候，按照主键排序，覆盖索引可以不包含主键
			if a.drivingTableColumn.columnInOrderBy.len() == 1 && a.drivingTableColumn.columnInOrderBy.columnsWithSelectivity[0].columnName.Name.Name.L == primaryColumnName {
				orderByPrimaryKey = true

			}
			if orderByPrimaryKey {
				a.possibleColumns.delete(&ast.ColumnNameExpr{Name: primaryColumn})
			}
		}
		// 2. 遍历列的类型
		// 把类型不适合作为索引的列添加到columnShouldNotAdd中，目前有Blob/Text和JSON
		for _, columnDefine := range a.drivingTableCreateStmt.Cols {
			if !a.possibleColumns.has(&ast.ColumnNameExpr{Name: columnDefine.Name}) {
				continue
			}
			switch columnDefine.Tp.Tp {
			case mysql.TypeBlob, mysql.TypeLongBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeJSON:
				selecivity = a.selectivityMap[columnDefine.Name.Name.L]
				a.columnShouldNotAdd.add(&ast.ColumnNameExpr{Name: columnDefine.Name}, selecivity)
				a.possibleColumns.delete(&ast.ColumnNameExpr{Name: columnDefine.Name})
			}

		}
	}
	return nil
}

func (a *threeStarIndexAdvisor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch currentNode := in.(type) {
	case *ast.PatternInExpr:
		col, ok := currentNode.Expr.(*ast.ColumnNameExpr)
		if ok {
			if !a.isColumnInDrivingTable(col) {
				return in, false
			}
			a.drivingTableColumn.unequalColumnInWhere.add(
				col,
				a.selectivityMap[col.Name.Name.L],
			)
		}
	case *ast.BinaryOperationExpr:
		switch currentNode.Op {
		case opcode.EQ:
			if _, ok := currentNode.R.(*parser_driver.ValueExpr); !ok {
				return in, false
			}
			if col, ok := currentNode.L.(*ast.ColumnNameExpr); ok {
				if !a.isColumnInDrivingTable(col) {
					return in, false
				}
				a.drivingTableColumn.equalColumnInWhere.add(col, a.selectivityMap[col.Name.Name.L])
			}
		case opcode.LogicOr:
			a.usingOr = true
		case opcode.GE, opcode.GT, opcode.LE, opcode.LT, opcode.NE:
			if _, ok := currentNode.R.(*parser_driver.ValueExpr); !ok {
				return in, false
			}
			if col, ok := currentNode.L.(*ast.ColumnNameExpr); ok {
				if !a.isColumnInDrivingTable(col) {
					return in, false
				}
				a.drivingTableColumn.unequalColumnInWhere.add(
					col,
					a.selectivityMap[col.Name.Name.L],
				)
			}
		}
	}
	return in, false
}

func (a *threeStarIndexAdvisor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (a *threeStarIndexAdvisor) canGiveAdvice() bool {
	// 若没有备选的列，则无法给出三星索引建议
	if a.possibleColumns.len() == 0 {
		return false
	}
	// 如果使用了OR并且无法给出覆盖索引，基于三星索引建议对表加复合索引，该SQL也不会走索引
	// 这种情况应参考创建多个索引触发索引合并优化
	// 参考文档：https://dev.mysql.com/doc/refman/8.0/en/index-merge-optimization.html
	if a.usingOr && !a.canGiveCoverIndex() {
		return false
	}
	return true
}

func (a *threeStarIndexAdvisor) giveAdvice() {
	// 加入WHERE等值条件中的列
	for _, column := range a.drivingTableColumn.equalColumnInWhere.columnsWithSelectivity {
		if a.adviceColumns.len() == a.maxColumns {
			break
		}
		if a.shouldSkipColumn(column) {
			continue
		}
		a.adviceColumns.add(column.columnName, column.selectivity)
	}
	// 添加一个排序列
	if a.canAddOrderColumn() {
		for _, column := range a.drivingTableColumn.columnInOrderBy.columnsWithSelectivity {
			if a.shouldSkipColumn(column) {
				continue
			}
			if a.adviceColumns.len() < a.maxColumns {
				a.adviceColumns.add(column.columnName, column.selectivity)
			} else if a.adviceColumns.columnsWithSelectivity[a.adviceColumns.len()-1].selectivity < column.selectivity {
				// 当建议的列已满，此时如果建议的列的最后一列的区分度小于排序列的区分度，则替换该列为排序列
				a.adviceColumns.replace(column, len(a.adviceColumns.columnMap)-1)
			}
			break
		}
	}
	// 如果能够形成覆盖索引，则添加SELECT中的剩余列
	if a.canGiveCoverIndex() {
		for _, column := range a.drivingTableColumn.columnInFieldList.columnsWithSelectivity {
			if len(a.adviceColumns.columnMap) == a.maxColumns {
				break
			}
			if a.shouldSkipColumn(column) {
				continue
			}
			a.adviceColumns.add(column.columnName, column.selectivity)
		}
	}
	// 最后添加一列WHERE中的非等值列
	if a.canAddUnequalColumn() {
		if a.drivingTableColumn.unequalColumnInWhere.len() > 0 {
			a.adviceColumns.add(
				a.drivingTableColumn.unequalColumnInWhere.columnsWithSelectivity[0].columnName,
				a.drivingTableColumn.unequalColumnInWhere.columnsWithSelectivity[0].selectivity,
			)
		}
	}
}

func (a *threeStarIndexAdvisor) sortColumnBySelectivity() {
	a.drivingTableColumn.columnInOrderBy.sort()
	a.drivingTableColumn.columnInFieldList.sort()
	a.drivingTableColumn.equalColumnInWhere.sort()
	a.drivingTableColumn.unequalColumnInWhere.sort()
}

func (a threeStarIndexAdvisor) isColumnInDrivingTable(column *ast.ColumnNameExpr) bool {
	if column.Name.Table.L == "" {
		// 没有表名，说明只有一张表
		return true
	}
	return column.Name.Table.L == util.GetTableNameFromTableSource(a.drivingTableSource)
}

func (a threeStarIndexAdvisor) shouldSkipColumn(column columnWithSelectivity) bool {
	if !a.possibleColumns.has(column.columnName) {
		// 跳过非备选列
		return true
	}
	if a.adviceColumns.has(column.columnName) {
		// 跳过已有列
		return true
	}
	if a.columnLastAdd.has(column.columnName) {
		// 跳过最后添加的列
		return true
	}
	if a.columnShouldNotAdd.has(column.columnName) {
		// 跳过不建议添加的列
		return true
	}
	return false
}

func (a threeStarIndexAdvisor) canAddOrderColumn() bool {
	originNode, ok := a.originNode.(*ast.SelectStmt)
	if !ok {
		return false
	}
	if originNode.OrderBy == nil {
		return false
	}
	if len(originNode.OrderBy.Items) == 0 {
		return false
	}
	// 如果有多个不同方向的排序，则不将排序列放到索引中
	var firstOrder bool = originNode.OrderBy.Items[0].Desc
	for _, col := range originNode.OrderBy.Items {
		if col.Desc != firstOrder {
			return false
		}
	}
	// 如果排序列已在索引建议中则不添加
	for _, column := range a.drivingTableColumn.columnInOrderBy.columnsWithSelectivity {
		if a.columnLastAdd.has(column.columnName) {
			continue
		}
		if a.adviceColumns.has(column.columnName) {
			return false
		}
	}
	return true
}

func (a threeStarIndexAdvisor) canAddUnequalColumn() bool {
	return a.adviceColumns.len() < a.maxColumns
}

func (a threeStarIndexAdvisor) canGiveCoverIndex() bool {
	// 非等值列大于1时，覆盖索引走不到索引的最后一列，不添加覆盖索引
	if a.drivingTableColumn.unequalColumnInWhere.len() > 1 {
		return false
	}
	// 一旦SQL中有列不适合作为索引的列，则无法给出覆盖索引
	if a.columnShouldNotAdd.len() > 0 {
		return false
	}
	// 当备选列大于索引列的上限时，覆盖索引不满足该限制，不添加覆盖索引
	if a.possibleColumns.len() > a.maxColumns {
		return false
	}
	return true
}

/*
joinIndexAdvisor

	1 驱动表不一定是From后的第一个表，要根据ExecutionPlan得出
	2 该规则仅检查被驱动表是否需要添加相应索引
	3 ast.JoinStmt是根据SQL语句来组织语法树的 因此被驱动表可能是抽象语法树的左节点和右节点其一 如果右节点是驱动表 那左节点就是被驱动表 如果右节点不是驱动表 则右节点为被驱动表 左节点可能是驱动表也可能是ast.JoinStmt
*/
type joinIndexAdvisor struct {
	sqlContext        *session.Context
	log               *logrus.Entry
	originNode        ast.Node
	currentNode       *ast.Join
	drivenTableSource map[string] /*table name*/ *ast.TableSource
	advices           []*OptimizeResult
}

func newJoinIndexAdvisor(ctx *session.Context, log *logrus.Entry, originNode ast.Node, params params.Params) CreateIndexAdvisor {
	return &joinIndexAdvisor{
		sqlContext:        ctx,
		log:               log,
		originNode:        originNode,
		drivenTableSource: make(map[string]*ast.TableSource),
	}
}

func (a *joinIndexAdvisor) GiveAdvices() []*OptimizeResult {
	if a.originNode == nil {
		return nil
	}
	err := a.loadEssentials()
	if err != nil {
		a.log.Logger.Warnf("when join index advisor load essentials failed, err:%v", err)
		return nil
	}
	a.originNode.Accept(a)
	return a.advices
}

func (a *joinIndexAdvisor) loadEssentials() error {
	executionPlans, err := a.sqlContext.GetExecutionPlan(a.originNode.Text())
	if err != nil {
		return err
	}
	extractor := util.TableSourceExtractor{TableSources: map[string]*ast.TableSource{}}
	a.originNode.Accept(&extractor)
	for id, record := range executionPlans {
		if id == 0 {
			continue
		}
		if record.Type == executor.ExplainRecordAccessTypeAll || record.Type == executor.ExplainRecordAccessTypeIndex {
			recordTableNameLow := strings.ToLower(record.Table)
			tableSource, ok := extractor.TableSources[recordTableNameLow]
			if !ok {
				continue
			}
			a.drivenTableSource[recordTableNameLow] = tableSource
		}
	}
	return nil
}

func (a *joinIndexAdvisor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch currentNode := in.(type) {
	case *ast.Join:
		a.currentNode = currentNode
		a.giveAdvice()
	}
	return in, false
}

func (v *joinIndexAdvisor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (a *joinIndexAdvisor) giveAdvice() {
	indexColumnMap := make(columnMap)
	drivenTableName := a.getDrivenTableName()
	if drivenTableName == "" {
		return
	}
	// 在ON和USING中找被驱动表的列
	if a.currentNode.On != nil {
		bo, ok := a.currentNode.On.Expr.(*ast.BinaryOperationExpr)
		if ok {
			leftName, ok := bo.L.(*ast.ColumnNameExpr)
			if ok && leftName.Name.Table.L == drivenTableName {
				indexColumnMap.add(&ast.ColumnNameExpr{Name: leftName.Name})
			}
			rightName, ok := bo.R.(*ast.ColumnNameExpr)
			if ok && rightName.Name.Table.L == drivenTableName {
				indexColumnMap.add(&ast.ColumnNameExpr{Name: rightName.Name})
			}
		}
	}
	if len(a.currentNode.Using) > 0 {
		// https://dev.mysql.com/doc/refman/8.0/en/join.html
		for _, column := range a.currentNode.Using {
			indexColumnMap.add(&ast.ColumnNameExpr{Name: column})
		}
	}
	if len(indexColumnMap) == 0 {
		return
	}

	tableSource := a.drivenTableSource[drivenTableName]
	tableName, ok := tableSource.Source.(*ast.TableName)
	if !ok {
		a.log.Warn("in join index advisor driven tableSource.Source is not ast.TableName")
		return
	}
	createTable, exist, err := a.sqlContext.GetCreateTableStmt(tableName)
	if err != nil {
		a.log.Warnf("join index advisor get create table statement failed,err %v", err)
		return
	}
	if !exist {
		a.log.Warnf("join index advisor get create table statement failed,table not exist %s", drivenTableName)
		return
	}
	if util.IsIndex(indexColumnMap, createTable.Constraints) {
		return
	}
	indexColumn := make([]string, 0, len(indexColumnMap))
	for column := range indexColumnMap {
		indexColumn = append(indexColumn, column)
	}
	a.advices = append(a.advices, &OptimizeResult{
		TableName:      drivenTableName,
		IndexedColumns: indexColumn,
		Reason:         fmt.Sprintf(joinIndexAdviceFormat, strings.Join(indexColumn, "，"), drivenTableName, drivenTableName, strings.Join(indexColumn, "，")),
	})
}

// 获取到Join节点左右节点中被驱动表的名称
func (a joinIndexAdvisor) getDrivenTableName() string {

	if tableSource, ok := a.currentNode.Right.(*ast.TableSource); ok {
		if tableSource.AsName.L != "" {
			if _, ok := a.drivenTableSource[tableSource.AsName.L]; ok {
				return tableSource.AsName.L
			}
		}
		if tableName, ok := tableSource.Source.(*ast.TableName); ok {
			if _, ok := a.drivenTableSource[tableName.Name.L]; ok {
				return tableName.Name.L
			}
		}
	}
	if tableSource, ok := a.currentNode.Left.(*ast.TableSource); ok {
		if tableSource.AsName.L != "" {
			if _, ok := a.drivenTableSource[tableSource.AsName.L]; ok {
				return tableSource.AsName.L
			}
		}
		if tableName, ok := tableSource.Source.(*ast.TableName); ok {
			if _, ok := a.drivenTableSource[tableName.Name.L]; ok {
				return tableName.Name.L
			}
		}
	}
	return ""
}

/*
functionIndexAdvisor 函数索引 虚拟列索引建议者

	触发条件:
		1. 判断WHERE子句的等值条件中是否使用了函数
		2. 如果使用函数，根据MySQL版本给出函数索引或虚拟列索引的建议

https://dev.mysql.com/doc/refman/8.0/en/create-index.html#create-index-functional-key-parts
*/
type functionIndexAdvisor struct {
	sqlContext         *session.Context
	log                *logrus.Entry
	originNode         ast.Node
	currentNode        *ast.BinaryOperationExpr
	drivingTableSource *ast.TableSource // 驱动表的TableSource
	advices            []*OptimizeResult
}

func newFunctionIndexAdvisor(ctx *session.Context, log *logrus.Entry, originNode ast.Node, params params.Params) CreateIndexAdvisor {
	return &functionIndexAdvisor{
		sqlContext: ctx,
		log:        log,
		originNode: originNode,
	}
}

func (a *functionIndexAdvisor) GiveAdvices() []*OptimizeResult {
	node, ok := a.originNode.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	if node.Where == nil {
		return nil
	}
	err := a.loadEssentials()
	if err != nil {
		a.log.Logger.Warnf("when function index advisor load essentials failed, err:%v", err)
		return nil
	}
	node.Where.Accept(a)
	return a.advices
}

func (a *functionIndexAdvisor) loadEssentials() (err error) {
	a.drivingTableSource, _, err = getDrivingTableInfo(a.originNode, a.sqlContext)
	if err != nil {
		return err
	}
	return nil
}

func (a *functionIndexAdvisor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch currentNode := in.(type) {
	case *ast.BinaryOperationExpr:
		if currentNode.Op == opcode.EQ {
			a.currentNode = currentNode
			a.giveAdvice()
		}
	}
	return in, false
}

func (v *functionIndexAdvisor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (a *functionIndexAdvisor) giveAdvice() {

	if _, ok := a.currentNode.L.(*ast.FuncCallExpr); !ok {
		if _, ok := a.currentNode.R.(*ast.FuncCallExpr); !ok {
			return
		}
	}
	columnNameVisitor := util.ColumnNameVisitor{}
	a.currentNode.L.Accept(&columnNameVisitor)
	if len(columnNameVisitor.ColumnNameList) == 0 {
		return
	}
	var curVersion *semver.Version
	versionWithFlavor, err := a.sqlContext.GetSystemVariable("version")
	if err != nil {
		a.log.Logger.Warnf("when function index advisor get system version failed %v", err)
	} else {
		curVersion, err = semver.NewVersion(versionWithFlavor)
		if err != nil {
			a.log.Logger.Warnf("when function index advisor parse version %s, failed %v", versionWithFlavor, err)

		} else {
			if curVersion.LessThan(semver.MustParse("5.7.0")) {
				return
			}
		}
	}

	columns := make([]string, 0, len(columnNameVisitor.ColumnNameList))
	var tableName string = columnNameVisitor.ColumnNameList[0].Name.Table.L
	for _, columnName := range columnNameVisitor.ColumnNameList {
		columns = append(columns, columnName.Name.Name.L)
	}
	if tableName == "" {
		tableName = util.GetTableNameFromTableSource(a.drivingTableSource)
	}
	if curVersion != nil && curVersion.LessThan(semver.MustParse("8.0.13")) {
		a.advices = append(a.advices, &OptimizeResult{
			TableName:      tableName,
			IndexedColumns: columns,
			Reason:         fmt.Sprintf(functionIndexAdviceFormatV57, tableName, strings.Join(columns, "，")),
		})
		return
	}
	if curVersion != nil && curVersion.GreaterThan(semver.MustParse("8.0.12")) {
		a.advices = append(a.advices, &OptimizeResult{
			TableName:      tableName,
			IndexedColumns: columns,
			Reason:         fmt.Sprintf(functionIndexAdviceFormatV80, tableName, strings.Join(columns, "，")),
		})
		return
	}
	// 某些版本解析会出错，例如"8.0.35-0<system_name>0.22.04.1"
	a.advices = append(a.advices, &OptimizeResult{
		TableName:      tableName,
		IndexedColumns: columns,
		Reason:         fmt.Sprintf(functionIndexAdviceFormatAll, tableName, strings.Join(columns, "，")),
	})
}

/*
extremalIndexAdvisor 极值索引建议者

	触发条件:
		1. WHERE等值条件中使用了聚合函数：*ast.AggregateFuncExpr
		2. 检查聚合函数是MAX或者MIN

https://dev.mysql.com/doc/refman/8.0/en/aggregate-functions.html#function_max
https://dev.mysql.com/doc/refman/8.0/en/aggregate-functions.html#function_min
*/
type extremalIndexAdvisor struct {
	sqlContext             *session.Context
	log                    *logrus.Entry
	originNode             ast.Node
	currentNode            *ast.SelectField
	drivingTableCreateStmt *ast.CreateTableStmt // 驱动表的建表语句
	drivingTableSource     *ast.TableSource     // 驱动表的TableSource
	advices                []*OptimizeResult
}

func newExtremalIndexAdvisor(ctx *session.Context, log *logrus.Entry, originNode ast.Node, params params.Params) CreateIndexAdvisor {
	return &extremalIndexAdvisor{
		sqlContext: ctx,
		log:        log,
		originNode: originNode,
	}
}

func (a *extremalIndexAdvisor) GiveAdvices() []*OptimizeResult {
	if a.originNode == nil {
		return nil
	}
	err := a.loadEssentials()
	if err != nil {
		a.log.Logger.Warnf("when extremal index advisor load essentials failed, err:%v", err)
		return nil
	}
	a.originNode.Accept(a)
	return a.advices
}

func (a *extremalIndexAdvisor) loadEssentials() (err error) {
	a.drivingTableSource, a.drivingTableCreateStmt, err = getDrivingTableInfo(a.originNode, a.sqlContext)
	if err != nil {
		return err
	}
	return nil
}

func (a *extremalIndexAdvisor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch currentNode := in.(type) {
	case *ast.SelectField:
		a.currentNode = currentNode
		a.giveAdvice()
	}
	return in, false
}

func (v *extremalIndexAdvisor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (a *extremalIndexAdvisor) giveAdvice() {
	var indexColumn string
	var node *ast.AggregateFuncExpr
	var ok bool
	if node, ok = a.currentNode.Expr.(*ast.AggregateFuncExpr); !ok {
		return
	}
	if len(node.Args) == 0 {
		return
	}
	var column *ast.ColumnNameExpr
	if strings.ToLower(node.F) == ast.AggFuncMin || strings.ToLower(node.F) == ast.AggFuncMax {
		if column, ok = node.Args[0].(*ast.ColumnNameExpr); ok {
			indexColumn = column.Name.Name.L
		}
	} else {
		return
	}
	var tableName string
	if column.Name.Table.L != "" {
		tableName = column.Name.Table.L
	} else {
		tableName = util.GetTableNameFromTableSource(a.drivingTableSource)

	}
	if util.IsIndex(map[string]struct{}{indexColumn: {}}, a.drivingTableCreateStmt.Constraints) {
		return
	}
	a.advices = append(a.advices, &OptimizeResult{
		TableName:      tableName,
		IndexedColumns: []string{indexColumn},
		Reason:         fmt.Sprintf(extremalIndexAdviceFormat, tableName, indexColumn),
	})
}

/*
prefixIndexAdvisor 前缀索引建议者

	触发条件:
		1. WHERE语句中等值条件包含Like子句
		2. Like子句使用了前缀匹配
*/
type prefixIndexAdvisor struct {
	sqlContext         *session.Context
	log                *logrus.Entry
	originNode         ast.Node
	currentNode        *ast.PatternLikeExpr
	drivingTableSource *ast.TableSource // 驱动表的TableSource
	advices            []*OptimizeResult
}

func newPrefixIndexAdvisor(ctx *session.Context, log *logrus.Entry, originNode ast.Node, params params.Params) CreateIndexAdvisor {
	return &prefixIndexAdvisor{
		sqlContext: ctx,
		log:        log,
		originNode: originNode,
	}
}

func (a *prefixIndexAdvisor) GiveAdvices() []*OptimizeResult {
	node, ok := a.originNode.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	if node.Where == nil {
		return nil
	}
	err := a.loadEssentials()
	if err != nil {
		a.log.Logger.Warnf("when prefix index advisor load essentials failed, err:%v", err)
		return nil
	}
	node.Where.Accept(a)
	return a.advices
}

func (a *prefixIndexAdvisor) loadEssentials() error {
	executionPlans, err := a.sqlContext.GetExecutionPlan(a.originNode.Text())
	if err != nil {
		return err
	}
	extractor := util.TableSourceExtractor{TableSources: map[string]*ast.TableSource{}}
	a.originNode.Accept(&extractor)
	if len(executionPlans) > 0 {
		tableSource, ok := extractor.TableSources[strings.ToLower(executionPlans[0].Table)]
		if !ok {
			return fmt.Errorf("get driving table source failed")
		}
		a.drivingTableSource = tableSource
	}
	return nil
}

func (a *prefixIndexAdvisor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch currentNode := in.(type) {
	case *ast.PatternLikeExpr:
		a.currentNode = currentNode
		a.giveAdvice()
	}
	return in, false
}

func (v *prefixIndexAdvisor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func (a *prefixIndexAdvisor) giveAdvice() {
	if !util.CheckWhereLeftFuzzySearch(a.currentNode) {
		return
	}
	column, ok := a.currentNode.Expr.(*ast.ColumnNameExpr)
	if !ok {
		return
	}
	var tableName string
	if column.Name.Table.L != "" {
		tableName = column.Name.Table.L
	} else {
		tableName = util.GetTableNameFromTableSource(a.drivingTableSource)
	}
	a.advices = append(a.advices, &OptimizeResult{
		TableName:      tableName,
		IndexedColumns: []string{column.Name.Name.L},
		Reason:         prefixIndexAdviceFormat,
	})
}
