package mysql

import (
	"context"
	"database/sql"
	_driver "database/sql/driver"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/actiontech/sqle/sqle/driver/mysql/onlineddl"
	"github.com/actiontech/sqle/sqle/driver/mysql/optimizer/index"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/pingcap/parser/ast"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	allRules := make([]*driver.Rule, len(rulepkg.RuleHandlers))
	for i := range rulepkg.RuleHandlers {
		allRules[i] = &rulepkg.RuleHandlers[i].Rule
	}

	driver.RegisterAuditDriver(driver.DriverTypeMySQL, allRules, params.Params{})
	driver.RegisterDriverManger(driver.DriverTypeMySQL, NewDriverManagerFunc)

	if err := LoadPtTemplateFromFile("./scripts/pt-online-schema-change.template"); err != nil {
		panic(err)
	}
}

// MysqlDriverImpl implements driver.Driver interface
type MysqlDriverImpl struct {
	// Ctx is SQL session.
	Ctx *session.Context
	// cnf is task cnf, cnf variables record in rules.
	cnf *Config

	rules []*driver.Rule

	// result keep inspect result for single audited SQL.
	// It refresh on every Audit.
	result *driver.AuditResult
	// HasInvalidSql represent one of the commit sql base-validation failed.
	HasInvalidSql bool

	inst *driver.DSN

	log *logrus.Entry
	// dbConn is a SQL driver for MySQL.
	dbConn *executor.Executor
	// isConnected represent dbConn has Connected.
	isConnected bool
	// isOfflineAudit represent Audit without instance.
	isOfflineAudit bool
}

func NewInspect(log *logrus.Entry, cfg *driver.Config) (*MysqlDriverImpl, error) {
	var inspect = &MysqlDriverImpl{}

	if cfg.DSN != nil {
		conn, err := executor.NewExecutor(log, cfg.DSN, cfg.DSN.DatabaseName)
		if err != nil {
			return nil, errors.Wrap(err, "new executor in inspect")
		}
		inspect.isConnected = true
		inspect.dbConn = conn
		inspect.inst = cfg.DSN

		ctx := session.NewContext(nil, session.WithExecutor(conn))
		ctx.SetCurrentSchema(cfg.DSN.DatabaseName)

		inspect.Ctx = ctx
	} else {
		ctx := session.NewContext(nil)
		inspect.Ctx = ctx
	}

	inspect.log = log
	inspect.rules = cfg.Rules
	inspect.result = driver.NewInspectResults()
	inspect.isOfflineAudit = cfg.DSN == nil

	inspect.cnf = &Config{
		DMLRollbackMaxRows: -1,
		DDLOSCMinSize:      -1,
		DDLGhostMinSize:    -1,
	}
	for _, rule := range cfg.Rules {
		if rule.Name == rulepkg.ConfigDMLRollbackMaxRows {
			max := rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).Int()
			inspect.cnf.DMLRollbackMaxRows = int64(max)
		}
		if rule.Name == rulepkg.ConfigDDLOSCMinSize {
			min := rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).Int()
			inspect.cnf.DDLOSCMinSize = int64(min)
		}
		if rule.Name == rulepkg.ConfigDDLGhostMinSize {
			min := rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName).Int()
			inspect.cnf.DDLGhostMinSize = int64(min)
		}
		if rule.Name == rulepkg.ConfigOptimizeIndexEnabled {
			inspect.cnf.optimizeIndexEnabled = true
			inspect.cnf.calculateCardinalityMaxRow = rule.Params.GetParam(rulepkg.DefaultMultiParamsFirstKeyName).Int()
			inspect.cnf.compositeIndexMaxColumn = rule.Params.GetParam(rulepkg.DefaultMultiParamsSecondKeyName).Int()
		}
		if rule.Name == rulepkg.ConfigDMLExplainPreCheckEnable {
			inspect.cnf.dmlExplainPreCheckEnable = true
		}
		if rule.Name == rulepkg.ConfigSQLIsExecuted {
			inspect.cnf.isExecutedSQL = true
		}
	}

	return inspect, nil
}

func (i *MysqlDriverImpl) IsOfflineAudit() bool {
	return i.isOfflineAudit
}

func (i *MysqlDriverImpl) IsExecutedSQL() bool {
	return i.cnf.isExecutedSQL
}

func (i *MysqlDriverImpl) executeByGhost(ctx context.Context, query string, isDryRun bool) (_driver.Result, error) {
	node, err := i.ParseSql(query)
	if err != nil {
		return nil, errors.Wrap(err, "parse SQL")
	}

	stmt, ok := node[0].(*ast.AlterTableStmt)
	if !ok {
		return nil, errors.New("type assertion failed, unable to convert to expected type")
	}
	schema := i.Ctx.GetSchemaName(stmt.Table)

	run := func(dryRun bool) error {
		executor, err := onlineddl.NewExecutor(i.log, i.inst, schema, query)
		if err != nil {
			return err
		}

		err = executor.Execute(ctx, dryRun)
		if err != nil {
			return err
		}
		return nil
	}

	actionStr := "run"
	if isDryRun {
		actionStr = "dry-run"
	}

	i.log.Infof("%s gh-ost", actionStr)
	if err := run(isDryRun); err != nil {
		i.log.Errorf("%s gh-ost error:%v", actionStr, err)
		return nil, errors.Wrap(err, fmt.Sprintf("%s gh-ost", actionStr))
	}
	i.log.Infof("%s OK!", actionStr)
	return _driver.ResultNoRows, nil
}

func (i *MysqlDriverImpl) Exec(ctx context.Context, query string) (_driver.Result, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}

	useGhost, err := i.onlineddlWithGhost(query)
	if err != nil {
		return nil, errors.Wrap(err, "check whether use ghost or not")
	}

	if useGhost {
		if _, err := i.executeByGhost(ctx, query, true); err != nil {
			return nil, err
		}
		return i.executeByGhost(ctx, query, false)
	}

	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.Db.Exec(query)
}

func (i *MysqlDriverImpl) onlineddlWithGhost(query string) (bool, error) {
	if i.cnf.DDLGhostMinSize == -1 {
		return false, nil
	}

	node, err := i.ParseSql(query)
	if err != nil {
		return false, errors.Wrap(err, "parse SQL")
	}

	stmt, ok := node[0].(*ast.AlterTableStmt)
	if !ok {
		return false, nil
	}

	tableSize, err := i.Ctx.GetTableSize(stmt.Table)
	if err != nil {
		return false, errors.Wrap(err, "get table size")
	}

	return int64(tableSize) > i.cnf.DDLGhostMinSize, nil
}

func (i *MysqlDriverImpl) Tx(ctx context.Context, queries ...string) ([]_driver.Result, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.Db.Transact(queries...)
}

func (i *MysqlDriverImpl) query(ctx context.Context, query string, args ...interface{}) ([]map[string]sql.NullString, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.Db.Query(query, args...)
}

func (i *MysqlDriverImpl) Parse(ctx context.Context, sqlText string) ([]driver.Node, error) {
	nodes, err := i.ParseSql(sqlText)
	if err != nil {
		return nil, err
	}

	lowerCaseTableNames, err := i.Ctx.GetSystemVariable(session.SysVarLowerCaseTableNames)
	if err != nil {
		return nil, err
	}

	ns := make([]driver.Node, len(nodes))
	for i := range nodes {
		n := driver.Node{}
		fingerprint, err := util.Fingerprint(nodes[i].Text(), lowerCaseTableNames == "0")
		if err != nil {
			return nil, err
		}
		n.Fingerprint = fingerprint
		n.Text = nodes[i].Text()
		switch nodes[i].(type) {
		case ast.DMLNode:
			n.Type = driver.SQLTypeDML
		default:
			n.Type = driver.SQLTypeDDL
		}

		ns[i] = n
	}
	return ns, nil
}

func (i *MysqlDriverImpl) Audit(ctx context.Context, sql string) (*driver.AuditResult, error) {
	i.result = driver.NewInspectResults()

	if sql == "" {
		return nil, errors.New("sql is empty")
	}

	nodes, err := i.ParseSql(sql)
	if err != nil {
		return nil, err
	}

	if i.IsOfflineAudit() || i.IsExecutedSQL() {
		err = i.CheckInvalidOffline(nodes[0])
	} else {
		err = i.CheckInvalid(nodes[0])
	}
	if err != nil {
		return nil, err
	}

	if !i.result.HasResult() && i.cnf.dmlExplainPreCheckEnable {
		if err = i.CheckExplain(nodes[0]); err != nil {
			return nil, err
		}
	}

	if i.result.HasResult() {
		i.HasInvalidSql = true
		i.Logger().Warnf("SQL %s invalid, %s", nodes[0].Text(), i.result.Message())
	}

	var ghostRule *driver.Rule
	for _, rule := range i.rules {
		if rule.Name == rulepkg.ConfigDDLGhostMinSize {
			ghostRule = rule
		}

		handler, ok := rulepkg.RuleHandlerMap[rule.Name]
		if !ok || handler.Func == nil {
			continue
		}
		if i.IsOfflineAudit() && !handler.IsAllowOfflineRule(nodes[0]) {
			continue
		}
		if i.cnf.isExecutedSQL {
			if handler.OnlyAuditNotExecutedSQL {
				continue
			}
			if handler.IsDisableExecutedSQLRule(nodes[0]) {
				continue
			}
		}

		input := &rulepkg.RuleHandlerInput{
			Ctx:  i.Ctx,
			Rule: *rule,
			Res:  i.result,
			Node: nodes[0],
		}

		if err := handler.Func(input); err != nil {
			return nil, err
		}
	}

	if i.cnf.optimizeIndexEnabled && index.CanOptimize(i.log, i.Ctx, nodes[0]) {
		optimizer := index.NewOptimizer(
			i.log, i.Ctx,
			index.WithCalculateCardinalityMaxRow(i.cnf.calculateCardinalityMaxRow),
			index.WithCompositeIndexMaxColumn(i.cnf.compositeIndexMaxColumn),
		)

		advices, err := optimizer.Optimize(ctx, nodes[0].(*ast.SelectStmt))
		if err != nil {
			// ignore error, source: https://github.com/actiontech/sqle/issues/416
			i.log.Errorf("optimize sqle failed: %v", err)
		}

		var buf strings.Builder
		for _, advice := range advices {
			buf.WriteString(fmt.Sprintf("建议为表 %s 列 %s 添加索引", advice.TableName, strings.Join(advice.IndexedColumns, ",")))
			if advice.Reason != "" {
				buf.WriteString(fmt.Sprintf(", 原因(%s)", advice.Reason))
			}
		}
		i.result.Add(driver.RuleLevelNotice, buf.String())
	}

	// dry run gh-ost
	useGhost, err := i.onlineddlWithGhost(sql)
	if err != nil {
		return nil, errors.Wrap(err, "check whether use ghost or not")
	}
	if useGhost {
		if _, err := i.executeByGhost(ctx, sql, true); err != nil {
			i.result.Add(driver.RuleLevelError, fmt.Sprintf("表空间大小超过%vMB, 将使用gh-ost进行上线, 但是dry-run抛出如下错误: %v", i.cnf.DDLGhostMinSize, err))
		} else {
			i.result.Add(ghostRule.Level, fmt.Sprintf("表空间大小超过%vMB, 将使用gh-ost进行上线", i.cnf.DDLGhostMinSize))
		}
	}

	// print osc
	oscCommandLine, err := i.generateOSCCommandLine(nodes[0])
	if err != nil {
		return nil, err
	}
	if oscCommandLine != "" {
		i.result.Add(driver.RuleLevelNotice, fmt.Sprintf("[osc]%s", oscCommandLine))
	}

	if !i.IsExecutedSQL() {
		i.Ctx.UpdateContext(nodes[0])
	}

	return i.result, nil
}

func (i *MysqlDriverImpl) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	if i.IsOfflineAudit() {
		return "", "", nil
	}
	if i.HasInvalidSql {
		return "", "", nil
	}

	nodes, err := i.ParseSql(sql)
	if err != nil {
		return "", "", err
	}

	rollback, reason, err := i.GenerateRollbackSql(nodes[0])
	if err != nil {
		return "", "", err
	}

	i.Ctx.UpdateContext(nodes[0])

	return rollback, reason, nil
}

func (i *MysqlDriverImpl) Close(ctx context.Context) {
	i.closeDbConn()
}

func (i *MysqlDriverImpl) Ping(ctx context.Context) error {
	if i.IsOfflineAudit() {
		return nil
	}

	conn, err := i.getDbConn()
	if err != nil {
		return err
	}
	return conn.Db.Ping()
}

func (i *MysqlDriverImpl) Schemas(ctx context.Context) ([]string, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.ShowDatabases(true)
}

type Config struct {
	DMLRollbackMaxRows int64
	DDLOSCMinSize      int64
	DDLGhostMinSize    int64

	optimizeIndexEnabled       bool
	dmlExplainPreCheckEnable   bool
	calculateCardinalityMaxRow int
	compositeIndexMaxColumn    int
	isExecutedSQL              bool
}

func (i *MysqlDriverImpl) Context() *session.Context {
	return i.Ctx
}

func (i *MysqlDriverImpl) ParseSql(sql string) ([]ast.Node, error) {
	stmts, err := util.ParseSql(sql)
	if err != nil {
		i.Logger().Errorf("parse sql failed, error: %v, sql: %s", err, sql)
		return nil, err
	}
	nodes := make([]ast.Node, 0, len(stmts))
	for _, stmt := range stmts {
		// node can only be ast.Node
		//nolint:forcetypeassert
		node := stmt.(ast.Node)
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (i *MysqlDriverImpl) Logger() *logrus.Entry {
	return i.log
}

// getDbConn get db conn and just connect once.
func (i *MysqlDriverImpl) getDbConn() (*executor.Executor, error) {
	if i.isConnected {
		return i.dbConn, nil
	}
	conn, err := executor.NewExecutor(i.log, i.inst, i.Ctx.CurrentSchema())
	if err == nil {
		i.isConnected = true
		i.dbConn = conn
	}
	return conn, err
}

// closeDbConn close db conn and just close once.
func (i *MysqlDriverImpl) closeDbConn() {
	if i.isConnected {
		i.dbConn.Db.Close()
		i.isConnected = false
	}
}

// getTableName get table name from TableName ast.
func (i *MysqlDriverImpl) getTableName(stmt *ast.TableName) string {
	schema := i.Ctx.GetSchemaName(stmt)
	if schema == "" {
		return stmt.Name.String()
	}
	return fmt.Sprintf("%s.%s", schema, stmt.Name)
}

// getTableNameWithQuote get table name with quote.
func (i *MysqlDriverImpl) getTableNameWithQuote(stmt *ast.TableName) string {
	name := strings.Replace(i.getTableName(stmt), ".", "`.`", -1)
	return fmt.Sprintf("`%s`", name)
}

// getPrimaryKey get table's primary key.
func (i *MysqlDriverImpl) getPrimaryKey(stmt *ast.CreateTableStmt) (map[string]struct{}, bool, error) {
	pkColumnsName, hasPk := util.GetPrimaryKey(stmt)
	if !hasPk {
		return pkColumnsName, hasPk, nil
	}
	return pkColumnsName, hasPk, nil
}

type DriverManager struct {
	inspect *MysqlDriverImpl
}

func (d *DriverManager) GetAuditDriver() (driver.Driver, error) {
	return d.inspect, nil
}

func (d *DriverManager) GetSQLQueryDriver() (driver.SQLQueryDriver, error) {
	return d.getSQLQueryDriver()
}

func (d *DriverManager) GetAnalysisDriver() (driver.AnalysisDriver, error) {
	return d.getAnalysisDriver()
}

func (d *DriverManager) Close(ctx context.Context) {
	d.inspect.Close(ctx)
}

func NewDriverManagerFunc(log *logrus.Entry, dbType string, config *driver.Config) (driver.DriverManager, error) {
	inspect, err := NewInspect(log, config)
	if err != nil {
		return nil, err
	}
	return &DriverManager{
		inspect: inspect,
	}, nil
}
