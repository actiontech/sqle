package mysql

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/pingcap/parser/ast"
	"github.com/sirupsen/logrus"
)

const (
	defaultCompositeIndexMaxColumn = 5
)

type OptimizeResult struct {
	TableName      string
	IndexedColumns []string

	Reason string
}

type indexOptimizer struct {
	sqlContext             *session.Context
	log                    *logrus.Entry
	tablesShouldBeOptimize map[string] /*table name*/ *ast.TableSource
	drivenTableSources     map[string] /*table name*/ *ast.TableSource
	drivingTableSource     *ast.TableSource
	drivingTableCreateStmt *ast.CreateTableStmt
	originNode             ast.Node
	// optimizer options:
	compositeIndexMaxColumn int
}

func NewIndexOptimizer(log *logrus.Entry, ctx *session.Context, node ast.Node, opts ...option) *indexOptimizer {
	log = log.WithField("optimizer", "index")

	optimizer := &indexOptimizer{
		originNode:              node,
		sqlContext:              ctx,
		log:                     log,
		compositeIndexMaxColumn: defaultCompositeIndexMaxColumn,
		tablesShouldBeOptimize:  make(map[string]*ast.TableSource),
		drivenTableSources:      make(map[string]*ast.TableSource),
	}

	for _, opt := range opts {
		opt.apply(optimizer)
	}

	return optimizer
}

type option func(*indexOptimizer)

func (oo option) apply(o *indexOptimizer) {
	oo(o)
}

func WithMaxColumn(column int) option {
	return func(o *indexOptimizer) {
		o.compositeIndexMaxColumn = column
	}
}

func (optimizer *indexOptimizer) Optimize() []*OptimizeResult {
	if !optimizer.canOptimize() {
		return nil
	}

	optimizer.loadTablesShouldBeOptimize()
	if optimizer.hasNoTablesToOptimize() {
		return nil
	}

	return optimizer.generateOptimizeResult()
}

func (opt *indexOptimizer) canOptimize() bool {
	canNotOptimizeWarnf := "can not optimize node: %v, reason: %v"
	if opt.sqlContext == nil {
		return false
	}
	if opt.originNode == nil {
		return false
	}
	selectStmt, ok := opt.originNode.(*ast.SelectStmt)
	if !ok {
		opt.log.Warnf(canNotOptimizeWarnf, opt.originNode, "not select statement")
		return false
	}

	if selectStmt.From == nil {
		opt.log.Warnf(canNotOptimizeWarnf, opt.originNode, "no from clause")
		return false
	}

	extractor := util.TableNameExtractor{TableNames: map[string]*ast.TableName{}}
	selectStmt.Accept(&extractor)
	for name, ast := range extractor.TableNames {
		exist, err := opt.sqlContext.IsTableExistInDatabase(ast)
		if err != nil {
			opt.log.Warnf(canNotOptimizeWarnf, opt.originNode, err)
			return false
		}
		if !exist {
			opt.log.Warnf(canNotOptimizeWarnf, opt.originNode, fmt.Sprintf("table %s not exist", name))
			return false
		}
	}

	return true
}

func (opt *indexOptimizer) loadTablesShouldBeOptimize() {
	executionPlans, err := opt.sqlContext.GetExecutionPlan(opt.originNode.Text())
	if err != nil {
		// explain will executed failure, if SQL is collect from MyBatis, it not executable SQL.
		opt.log.Errorf("get execution plan failed, sql: %v, error: %v", opt.originNode.Text(), err)
		return
	}
	extractor := util.TableSourceExtractor{TableSources: map[string]*ast.TableSource{}}
	opt.originNode.Accept(&extractor)

	for id, record := range executionPlans {
		if record.Type == executor.ExplainRecordAccessTypeAll || record.Type == executor.ExplainRecordAccessTypeIndex {
			recordTableNameLow := strings.ToLower(record.Table)
			tableSource, ok := extractor.TableSources[recordTableNameLow]
			if !ok {
				continue
			}

			if id == 0 {
				opt.drivingTableSource = tableSource
				tableName, _ := tableSource.Source.(*ast.TableName)

				createTableStmt, _, err := opt.sqlContext.GetCreateTableStmt(tableName)
				if err != nil {
					continue
				}
				opt.drivingTableCreateStmt = createTableStmt
			} else {
				opt.drivenTableSources[recordTableNameLow] = tableSource
			}
			opt.tablesShouldBeOptimize[recordTableNameLow] = tableSource
		}
	}
}

func (opt *indexOptimizer) hasNoTablesToOptimize() bool {
	return len(opt.tablesShouldBeOptimize) == 0
}

func (opt *indexOptimizer) generateOptimizeResult() []*OptimizeResult {
	advisor := selectAdvisorVisitor{
		log:                  opt.log,
		extremalIndexAdvisor: newExtremalIndexAdvisor(opt.drivingTableSource),
		threeStarAdvisor: newThreeStarAdvisor(
			opt.sqlContext, opt.log,
			opt.drivingTableSource, opt.drivingTableCreateStmt, opt.originNode,
			opt.compositeIndexMaxColumn,
		),
		whereAdvisorVisitor: whereAdvisorVisitor{
			advices:              make([]*OptimizeResult, 0),
			prefixIndexAdvisor:   newPrefixIndexAdvisor(opt.drivingTableSource),
			functionIndexAdvisor: newFunctionIndexAdvisor(opt.sqlContext, opt.log, opt.drivingTableSource),
		},
		fromAdvisorVisitor: fromAdvisorVisitor{
			joinIndexAdvisor: newJoinAdvisor(opt.sqlContext, opt.drivenTableSources),
			advices:          make([]*OptimizeResult, 0),
		},
		advices: make([]*OptimizeResult, 0),
	}

	opt.originNode.Accept(&advisor)
	return advisor.advices
}
