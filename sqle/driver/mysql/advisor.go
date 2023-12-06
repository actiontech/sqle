package mysql

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
	parser_driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/sirupsen/logrus"
)

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

/*
SELECT语句建议者

	1 组合了基于SELECT语句给出建议的Advisor
	2 遍历SQL语句中所有SELECT语句，使用具体的Advisor给出添加索引的建议
*/
type selectAdvisorVisitor struct {
	log                  *logrus.Entry
	joinIndexAdvisor     CreateIndexAdvisor
	threeStarAdvisor     CreateIndexAdvisor
	extremalIndexAdvisor CreateIndexAdvisor
	whereAdvisorVisitor
	advices []*OptimizeResult
}

func (a *selectAdvisorVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch currentNode := in.(type) {
	case *ast.SelectField:
		if a.extremalIndexAdvisor != nil {
			advice := a.extremalIndexAdvisor.setCurrentNode(currentNode).GiveAdvice()
			if advice != nil {
				a.advices = append(a.advices, advice)
			}
		}
	case *ast.Join:
		if a.joinIndexAdvisor != nil {
			advice := a.joinIndexAdvisor.setCurrentNode(currentNode).GiveAdvice()
			if advice != nil {
				a.advices = append(a.advices, advice)
			}
		}
	case *ast.SelectStmt:
		if currentNode.Where != nil {
			// 检查Where子句
			currentNode.Where.Accept(&a.whereAdvisorVisitor)
			a.advices = append(a.advices, a.whereAdvisorVisitor.advices...)
		}

		// 检查Select语句
		if a.threeStarAdvisor != nil {
			advice := a.threeStarAdvisor.setCurrentNode(currentNode).GiveAdvice()
			if advice != nil {
				a.advices = append(a.advices, advice)
			}
		}
	}
	return in, false
}

func (v *selectAdvisorVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

/*
WHERE子句建议者

	1 组合了基于WHERE子句给出建议的Advisor
	2 遍历WHERE子句，使用具体的Advisor给出建议
*/
type whereAdvisorVisitor struct {
	functionIndexAdvisor CreateIndexAdvisor
	prefixIndexAdvisor   CreateIndexAdvisor
	advices              []*OptimizeResult
}

func (opt *whereAdvisorVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	var advice *OptimizeResult
	switch currentNode := in.(type) {
	case *ast.BinaryOperationExpr:
		if currentNode.Op == opcode.EQ {
			if opt.functionIndexAdvisor != nil {
				advice = opt.functionIndexAdvisor.setCurrentNode(currentNode).GiveAdvice()
				if advice != nil {
					opt.advices = append(opt.advices, advice)
				}
			}
		}
	case *ast.PatternLikeExpr:
		if opt.prefixIndexAdvisor != nil {
			advice = opt.prefixIndexAdvisor.setCurrentNode(currentNode).GiveAdvice()
			if advice != nil {
				opt.advices = append(opt.advices, advice)
			}
		}
	}
	return in, false
}

func (v *whereAdvisorVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

// CreateIndexAdvisor 基于SQL语句、SQL上下文、库表信息等生成创建索引的建议，在给出建议前需要指明优化建议针对的节点
type CreateIndexAdvisor interface {
	GiveAdvice() *OptimizeResult
	setCurrentNode(interface{}) CreateIndexAdvisor
}

type threeStarAdvisor struct {
	sqlContext             *session.Context
	log                    *logrus.Entry
	drivingTableSource     *ast.TableSource        // 驱动表的TableSource
	drivingTableColumn     *columnInSelect         // 经过解析的驱动表的列
	drivingTableCreateStmt *ast.CreateTableStmt    // 驱动表的建表语句
	originNode             ast.Node                // 原SQL的节点
	currentNode            ast.Node                // 当前节点
	maxColumns             int                     // 复合索引列的上限数量
	possibleColumns        columnMap               // SQL语句中可能作为索引的备选列
	columnLastAdd          columnMap               // 最后添加的列，例如：非等值列
	columnShouldNotAdd     columnMap               // 不该添加的列，例如：类型不适合作为索引的列、单列主键列
	columnHasAdded         columnMap               // 已经添加的列
	adviceColumns          []columnWithSelectivity // 给出建议的列
}

type columnMap map[string] /*column name or alias name*/ struct{}

func (c columnMap) add(col *ast.ColumnNameExpr) {
	c[col.Name.Name.L] = struct{}{}
}
func (c columnMap) delete(col *ast.ColumnNameExpr) {
	delete(c, col.Name.Name.L)
}

type columnInSelect struct {
	equalColumnInWhere   columnsWithSelectivity
	unequalColumnInWhere columnsWithSelectivity
	columnInOrderBy      columnsWithSelectivity
	columnInFieldList    columnsWithSelectivity
}

func newThreeStarAdvisor(ctx *session.Context, log *logrus.Entry, drivingTableSource *ast.TableSource, drivingTableCreateStmt *ast.CreateTableStmt, originNode ast.Node, maxColumns int) CreateIndexAdvisor {
	return &threeStarAdvisor{
		sqlContext:             ctx,
		log:                    log,
		drivingTableSource:     drivingTableSource,
		drivingTableCreateStmt: drivingTableCreateStmt,
		originNode:             originNode,
		maxColumns:             maxColumns,
		drivingTableColumn:     &columnInSelect{},
		columnLastAdd:          make(columnMap),
		possibleColumns:        make(columnMap),
		columnHasAdded:         make(columnMap, maxColumns),
		columnShouldNotAdd:     make(columnMap),
		adviceColumns:          make([]columnWithSelectivity, 0),
	}
}

func (a *threeStarAdvisor) setCurrentNode(node interface{}) CreateIndexAdvisor {
	a.currentNode, _ = node.(ast.Node)
	return a
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
func (a *threeStarAdvisor) GiveAdvice() *OptimizeResult {
	if !a.canGiveAdvice() {
		return nil
	}
	err := a.extractColumnInSelect()
	if err != nil {
		a.log.Logger.Warnf("extract column in select failed, sql:%s,err:%v", restore(a.originNode), err)
		return nil
	}
	if len(a.possibleColumns) == 0 {
		return nil
	}
	err = a.fillColumnWithSelectivity()
	if err != nil {
		a.log.Logger.Warnf("fill column with selectivity failed, sql:%s,err:%v", restore(a.originNode), err)
		return nil
	}
	a.sortColumnBySelectivity()
	a.giveAdvice()
	if util.IsIndex(a.columnHasAdded, a.drivingTableCreateStmt.Constraints) {
		return nil
	}
	return &OptimizeResult{
		TableName:      util.GetTableNameFromTableSource(a.drivingTableSource),
		IndexedColumns: a.indexColumns(),
		Reason:         fmt.Sprintf("索引建议 | SQL：%s 中，根据三星索引设计规范", restore(a.currentNode)),
	}
}

func (a threeStarAdvisor) indexColumns() []string {
	indexedColumn := make([]string, 0, len(a.adviceColumns))
	for _, column := range a.adviceColumns {
		indexedColumn = append(indexedColumn, column.columnName.Name.Name.L)
	}
	return indexedColumn
}

func (a *threeStarAdvisor) giveAdvice() {
	// 加入WHERE等值条件中的列
	for _, column := range a.drivingTableColumn.equalColumnInWhere {
		if len(a.adviceColumns) == a.maxColumns {
			break
		}
		if a.shouldSkipColumn(column) {
			continue
		}
		a.addColumn(column)
	}
	// 添加一个排序列
	if a.canAddOrderColumn() {
		for _, column := range a.drivingTableColumn.columnInOrderBy {
			if a.shouldSkipColumn(column) {
				continue
			}
			if len(a.adviceColumns) < a.maxColumns {
				a.addColumn(column)
			} else if a.adviceColumns[len(a.adviceColumns)-1].selectivity < column.selectivity {
				// 当建议的列已满，此时如果建议的列的最后一列的区分度小于排序列的区分度，则替换该列为排序列
				a.replaceColumn(column, len(a.adviceColumns)-1)
			}
			break
		}
	}
	// 如果能够形成覆盖索引，则添加SELECT中的剩余列
	if a.canGiveCoverIndex() {
		for _, column := range a.drivingTableColumn.columnInFieldList {
			if len(a.adviceColumns) == a.maxColumns {
				break
			}
			if a.shouldSkipColumn(column) {
				continue
			}
			a.addColumn(column)
		}
	}
	// 最后添加一列WHERE中的非等值列
	if a.canAddUnequalColumn() {
		if len(a.drivingTableColumn.unequalColumnInWhere) > 0 {
			a.addColumn(a.drivingTableColumn.unequalColumnInWhere[0])
		}
	}
}

func (a threeStarAdvisor) shouldSkipColumn(column columnWithSelectivity) bool {
	columnName := column.columnName.Name.Name.L
	if _, exist := a.possibleColumns[columnName]; !exist {
		// 跳过非备选列
		return true
	}
	if _, exist := a.columnHasAdded[columnName]; exist {
		// 跳过已有列
		return true
	}
	if _, exist := a.columnLastAdd[columnName]; exist {
		// 跳过最后添加的列
		return true
	}
	if _, exist := a.columnShouldNotAdd[columnName]; exist {
		// 跳过不建议添加的列
		return true
	}
	return false
}

func (a *threeStarAdvisor) addColumn(column columnWithSelectivity) {
	a.adviceColumns = append(a.adviceColumns, column)
	a.columnHasAdded.add(column.columnName)
}

func (a *threeStarAdvisor) replaceColumn(newColumn columnWithSelectivity, index int) {
	oldColumn := a.adviceColumns[index]
	a.adviceColumns[index] = newColumn
	a.columnHasAdded.delete(oldColumn.columnName)
	a.columnHasAdded.add(newColumn.columnName)
}

func (a threeStarAdvisor) canAddOrderColumn() bool {
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
	for _, column := range a.drivingTableColumn.columnInOrderBy {
		columnName := column.columnName.Name.Name.L
		if _, exist := a.columnLastAdd[columnName]; exist {
			continue
		}
		if _, exist := a.columnHasAdded[columnName]; exist {
			return false
		}
	}
	return true
}

func (a threeStarAdvisor) canAddUnequalColumn() bool {
	return len(a.adviceColumns) < a.maxColumns
}

func (a threeStarAdvisor) canGiveCoverIndex() bool {
	// 非等值列大于1时，覆盖索引走不到索引的最后一列，不添加覆盖索引
	if len(a.drivingTableColumn.unequalColumnInWhere) > 1 {
		return false
	}
	// 当备选列大于索引列的上限时，覆盖索引不满足该限制，不添加覆盖索引
	if len(a.possibleColumns) > a.maxColumns {
		return false
	}
	return true
}

func (a threeStarAdvisor) canGiveAdvice() bool {

	if a.currentNode != a.originNode {
		// 三星索引建议只对SQL原语句对应的节点生效
		return false
	}
	// 三星索引建议只对驱动表给出索引建议 不对被驱动表给出
	return a.drivingTableSource != nil
}

func (a threeStarAdvisor) isColumnInDrivingTable(column *ast.ColumnNameExpr) bool {
	if column.Name.Table.L == "" {
		// 没有表名，说明只有一张表
		return true
	}
	return column.Name.Table.L == util.GetTableNameFromTableSource(a.drivingTableSource)
}

func (a *threeStarAdvisor) sortColumnBySelectivity() {
	sort.Sort(a.drivingTableColumn.columnInOrderBy)
	sort.Sort(a.drivingTableColumn.columnInFieldList)
	sort.Sort(a.drivingTableColumn.unequalColumnInWhere)
	sort.Sort(a.drivingTableColumn.equalColumnInWhere)
}

func (a *threeStarAdvisor) fillColumnWithSelectivity() error {
	tableName, ok := a.drivingTableSource.Source.(*ast.TableName)
	if !ok {
		return fmt.Errorf("in three star advisor driving tableSource.Source is not ast.TableName")
	}
	columnNames := make([]string, 0, len(a.possibleColumns))
	for key := range a.possibleColumns {
		columnNames = append(columnNames, key)
	}

	selectivityMap, err := a.sqlContext.GetSelectivityOfColumns(tableName, columnNames)
	if err != nil {
		return err
	}
	// 填充驱动表中各列的列区分度
	for i := range a.drivingTableColumn.equalColumnInWhere {
		a.drivingTableColumn.equalColumnInWhere[i].selectivity = selectivityMap[a.drivingTableColumn.equalColumnInWhere[i].columnName.Name.Name.L]
	}
	for i := range a.drivingTableColumn.columnInFieldList {
		a.drivingTableColumn.columnInFieldList[i].selectivity = selectivityMap[a.drivingTableColumn.columnInFieldList[i].columnName.Name.Name.L]
	}
	for i := range a.drivingTableColumn.columnInOrderBy {
		a.drivingTableColumn.columnInOrderBy[i].selectivity = selectivityMap[a.drivingTableColumn.columnInOrderBy[i].columnName.Name.Name.L]
	}
	for i := range a.drivingTableColumn.unequalColumnInWhere {
		a.drivingTableColumn.unequalColumnInWhere[i].selectivity = selectivityMap[a.drivingTableColumn.unequalColumnInWhere[i].columnName.Name.Name.L]
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
func (a *threeStarAdvisor) extractColumnInSelect() error {
	selectStmt, ok := a.currentNode.(*ast.SelectStmt)
	if !ok {
		return fmt.Errorf("in three star advisor, type of current node is not ast.SelectStmt")
	}
	if selectStmt.Where != nil {
		// 访问Where子句，解析并存储属于驱动表等值列和非等值列
		selectStmt.Where.Accept(a)
		for _, col := range a.drivingTableColumn.equalColumnInWhere {
			a.possibleColumns.add(col.columnName)
		}
		for _, col := range a.drivingTableColumn.unequalColumnInWhere {
			a.columnLastAdd.add(col.columnName)
			a.possibleColumns.add(col.columnName)
		}
	}
	if selectStmt.OrderBy != nil {
		// 遍历Order By的对象切片，存储其中属于驱动表的裸列
		for _, item := range selectStmt.OrderBy.Items {
			if col, ok := item.Expr.(*ast.ColumnNameExpr); ok {
				if !a.isColumnInDrivingTable(col) {
					continue
				}
				a.drivingTableColumn.columnInOrderBy = append(
					a.drivingTableColumn.columnInOrderBy,
					columnWithSelectivity{columnName: col},
				)
				a.possibleColumns.add(col)
			}
		}
	}
	if selectStmt.Fields != nil {
		// 遍历Select子句，存储其中属于驱动表的裸列
		for _, field := range selectStmt.Fields.Fields {
			if col, ok := field.Expr.(*ast.ColumnNameExpr); ok {
				if !a.isColumnInDrivingTable(col) {
					continue
				}
				a.drivingTableColumn.columnInFieldList = append(
					a.drivingTableColumn.columnInFieldList,
					columnWithSelectivity{columnName: col},
				)
				a.possibleColumns.add(col)
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
			if _, exist := a.possibleColumns[primaryColumnName]; !exist {
				return nil
			}
			var orderByPrimaryKey bool
			// 当没有Order By时，按照主键排序，覆盖索引可以不包含主键
			if len(a.drivingTableColumn.columnInOrderBy) == 0 {
				orderByPrimaryKey = true
			}
			// 当Order By主键的时候，按照主键排序，覆盖索引可以不包含主键
			if len(a.drivingTableColumn.columnInOrderBy) == 1 && a.drivingTableColumn.columnInOrderBy[0].columnName.Name.Name.L == primaryColumnName {
				orderByPrimaryKey = true

			}
			if orderByPrimaryKey {
				a.possibleColumns.delete(&ast.ColumnNameExpr{Name: primaryColumn})
			}
		}
		// 2. 遍历列的类型
		// 把不适合作为索引的列添加到columnShouldNotAdd中
		for _, columnDefine := range a.drivingTableCreateStmt.Cols {
			if _, exist := a.possibleColumns[columnDefine.Name.Name.L]; !exist {
				continue
			}
			if columnDefine.Tp.Tp == mysql.TypeBlob {
				a.columnShouldNotAdd.add(&ast.ColumnNameExpr{Name: columnDefine.Name})
			}
		}
	}
	return nil
}

func (a *threeStarAdvisor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	switch currentNode := in.(type) {
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
				a.drivingTableColumn.equalColumnInWhere = append(
					a.drivingTableColumn.equalColumnInWhere,
					columnWithSelectivity{
						columnName: col,
					},
				)
			}
		case opcode.GE, opcode.GT, opcode.LE, opcode.LT, opcode.NE:
			if _, ok := currentNode.R.(*parser_driver.ValueExpr); !ok {
				return in, false
			}
			if col, ok := currentNode.L.(*ast.ColumnNameExpr); ok {
				if !a.isColumnInDrivingTable(col) {
					return in, false
				}
				a.drivingTableColumn.unequalColumnInWhere = append(
					a.drivingTableColumn.unequalColumnInWhere,
					columnWithSelectivity{
						columnName: col,
					},
				)
			}
		}
	}
	return in, false
}

func (a *threeStarAdvisor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

type columnWithSelectivity struct {
	columnName  *ast.ColumnNameExpr
	selectivity float64
}

type columnsWithSelectivity []columnWithSelectivity

func (c columnsWithSelectivity) Len() int {
	return len(c)
}

func (c columnsWithSelectivity) Less(i, j int) bool {
	return c[i].selectivity > c[j].selectivity
}

func (c columnsWithSelectivity) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
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
	currentNode       *ast.Join
	drivenTableSource map[string] /*table name*/ *ast.TableSource
}

func newJoinAdvisor(ctx *session.Context, log *logrus.Entry, drivenTableSource map[string]*ast.TableSource) CreateIndexAdvisor {
	return joinIndexAdvisor{
		sqlContext:        ctx,
		log:               log,
		drivenTableSource: drivenTableSource,
	}
}

func (a joinIndexAdvisor) setCurrentNode(node interface{}) CreateIndexAdvisor {
	a.currentNode, _ = node.(*ast.Join)
	return a
}

func (a joinIndexAdvisor) GiveAdvice() *OptimizeResult {
	if a.currentNode == nil {
		return nil
	}
	indexColumnMap := make(columnMap)
	drivenTableName := a.getDrivenTableName()
	if drivenTableName == "" {
		return nil
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
		return nil
	}

	tableSource := a.drivenTableSource[drivenTableName]
	tableName, ok := tableSource.Source.(*ast.TableName)
	if !ok {
		a.log.Warn("in join index advisor driven tableSource.Source is not ast.TableName")
		return nil
	}
	createTable, exist, err := a.sqlContext.GetCreateTableStmt(tableName)
	if err != nil {
		a.log.Warnf("join index advisor get create table statement failed,err %v", err)
		return nil
	}
	if !exist {
		a.log.Warnf("join index advisor get create table statement failed,table not exist %s", drivenTableName)
		return nil
	}
	if util.IsIndex(indexColumnMap, createTable.Constraints) {
		return nil
	}
	indexColumn := make([]string, 0, len(indexColumnMap))
	for column := range indexColumnMap {
		indexColumn = append(indexColumn, column)
	}
	return &OptimizeResult{
		TableName:      drivenTableName,
		IndexedColumns: indexColumn,
		Reason:         fmt.Sprintf("索引建议 | SQL：%s 中，字段 %s 为被驱动表 %s 上的关联字段", restore(a.currentNode), strings.Join(indexColumn, "，"), drivenTableName),
	}
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
	log                *logrus.Entry
	sqlContext         *session.Context
	currentNode        *ast.BinaryOperationExpr
	drivingTableSource *ast.TableSource // 驱动表的TableSource
}

func newFunctionIndexAdvisor(ctx *session.Context, log *logrus.Entry, drivingTableSource *ast.TableSource) CreateIndexAdvisor {
	return &functionIndexAdvisor{
		log:                log,
		sqlContext:         ctx,
		drivingTableSource: drivingTableSource,
	}
}

func (a *functionIndexAdvisor) setCurrentNode(node interface{}) CreateIndexAdvisor {
	a.currentNode, _ = node.(*ast.BinaryOperationExpr)
	return a
}

func (a *functionIndexAdvisor) GiveAdvice() *OptimizeResult {
	if a.currentNode == nil {
		return nil
	}
	if _, ok := a.currentNode.L.(*ast.FuncCallExpr); !ok {
		if _, ok := a.currentNode.R.(*ast.FuncCallExpr); !ok {
			return nil
		}
	}
	columnNameVisitor := util.ColumnNameVisitor{}
	a.currentNode.L.Accept(&columnNameVisitor)
	if len(columnNameVisitor.ColumnNameList) == 0 {
		return nil
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
				return nil
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
		return &OptimizeResult{
			TableName:      tableName,
			IndexedColumns: columns,
			Reason:         fmt.Sprintf("索引建议 | SQL：%s 中，使用了函数作为查询条件，在MySQL5.7以上的版本，可以在虚拟列上创建索引", restore(a.currentNode.L)),
		}
	}
	if curVersion != nil && curVersion.GreaterThan(semver.MustParse("8.0.12")) {
		return &OptimizeResult{
			TableName:      tableName,
			IndexedColumns: columns,
			Reason:         fmt.Sprintf("索引建议 | SQL：%s 中，使用了函数作为查询条件，在MySQL8.0.13以上的版本，可以创建函数索引", restore(a.currentNode.L)),
		}
	}
	// 某些版本解析会出错，例如"8.0.35-0<system_name>0.22.04.1"
	return &OptimizeResult{
		TableName:      tableName,
		IndexedColumns: columns,
		Reason:         fmt.Sprintf("索引建议 | SQL：%s 中，使用了函数作为查询条件，在MySQL5.7以上的版本，可以在虚拟列上创建索引，在MySQL8.0.13以上的版本，可以创建函数索引", restore(a.currentNode.L)),
	}
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
	currentNode            *ast.SelectField
	drivingTableCreateStmt *ast.CreateTableStmt // 驱动表的建表语句
	drivingTableSource     *ast.TableSource     // 驱动表的TableSource
}

func newExtremalIndexAdvisor(drivingTableSource *ast.TableSource, drivingTableCreateStmt *ast.CreateTableStmt) CreateIndexAdvisor {
	return &extremalIndexAdvisor{
		drivingTableSource:     drivingTableSource,
		drivingTableCreateStmt: drivingTableCreateStmt,
	}
}

func (a *extremalIndexAdvisor) setCurrentNode(node interface{}) CreateIndexAdvisor {
	a.currentNode, _ = node.(*ast.SelectField)
	return a
}

func (a *extremalIndexAdvisor) GiveAdvice() *OptimizeResult {
	if a.currentNode == nil {
		return nil
	}
	var indexColumn string
	var node *ast.AggregateFuncExpr
	var ok bool
	if node, ok = a.currentNode.Expr.(*ast.AggregateFuncExpr); !ok {
		return nil
	}
	if len(node.Args) == 0 {
		return nil
	}
	var column *ast.ColumnNameExpr
	if strings.ToLower(node.F) == ast.AggFuncMin || strings.ToLower(node.F) == ast.AggFuncMax {
		if column, ok = node.Args[0].(*ast.ColumnNameExpr); ok {
			indexColumn = column.Name.Name.L
		}
	} else {
		return nil
	}
	var tableName string
	if column.Name.Table.L != "" {
		tableName = column.Name.Table.L
	} else {
		tableName = util.GetTableNameFromTableSource(a.drivingTableSource)
	}
	if util.IsIndex(map[string]struct{}{indexColumn: {}}, a.drivingTableCreateStmt.Constraints) {
		return nil
	}
	return &OptimizeResult{
		TableName:      tableName,
		IndexedColumns: []string{indexColumn},
		Reason:         fmt.Sprintf("索引建议 | SQL：%s 中，使用了最值函数，可以利用索引有序的性质快速找到最值", restore(a.currentNode)),
	}
}

/*
prefixIndexAdvisor 前缀索引建议者

	触发条件:
		1. WHERE语句中等值条件包含Like子句
		2. Like子句使用了前缀匹配
*/
type prefixIndexAdvisor struct {
	currentNode        *ast.PatternLikeExpr
	drivingTableSource *ast.TableSource // 驱动表的TableSource
}

func newPrefixIndexAdvisor(drivingTableSource *ast.TableSource) CreateIndexAdvisor {
	return &prefixIndexAdvisor{drivingTableSource: drivingTableSource}
}

func (a *prefixIndexAdvisor) setCurrentNode(node interface{}) CreateIndexAdvisor {
	a.currentNode, _ = node.(*ast.PatternLikeExpr)
	return a
}

func (a *prefixIndexAdvisor) GiveAdvice() *OptimizeResult {
	if a.currentNode == nil {
		return nil
	}
	if !util.CheckWhereFuzzySearch(a.currentNode) {
		return nil
	}
	column, ok := a.currentNode.Expr.(*ast.ColumnNameExpr)
	if !ok {
		return nil
	}
	var tableName string
	if column.Name.Table.L != "" {
		tableName = column.Name.Table.L
	} else {
		tableName = util.GetTableNameFromTableSource(a.drivingTableSource)
	}
	return &OptimizeResult{
		TableName:      tableName,
		IndexedColumns: []string{column.Name.Name.L},
		Reason:         fmt.Sprintf("索引建议 | SQL：%s 中，使用了前缀模式匹配，在数据量大的时候，可以建立翻转函数索引", restore(a.currentNode)),
	}
}
