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

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	"github.com/actiontech/sqle/sqle/driver/mysql/session"
	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MysqlDriverImpl implements driver.Driver interface
type MysqlDriverImpl struct {
	// Ctx is SQL session.
	Ctx *session.Context
	// cnf is task cnf, cnf variables record in rules.
	cnf *Config

	rules []*driverV2.Rule

	// result keep inspect result for single audited SQL.
	// It refresh on every Audit.
	result *driverV2.AuditResults
	// HasInvalidSql represent one of the commit sql base-validation failed.
	HasInvalidSql bool

	inst *driverV2.DSN

	log *logrus.Entry
	// dbConn is a SQL driver for MySQL.
	dbConn *executor.Executor
	// isConnected represent dbConn has Connected.
	isConnected bool
	// isOfflineAudit represent Audit without instance.
	isOfflineAudit bool
}

func NewInspect(log *logrus.Entry, cfg *driverV2.Config) (*MysqlDriverImpl, error) {
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
	inspect.result = driverV2.NewAuditResults()
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
			inspect.cnf.indexSelectivityMinValue = rule.Params.GetParam(rulepkg.DefaultMultiParamsFirstKeyName).Float64()
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

func (i *MysqlDriverImpl) SetExecutor(dbConn *executor.Executor) {
	i.dbConn = dbConn
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

func (i *MysqlDriverImpl) KillProcess(ctx context.Context) error {
	connID := i.dbConn.Db.GetConnectionID()
	if connID == "" {
		return fmt.Errorf("cannot find mysql conn_id, check logs")
	}
	logEntry := log.NewEntry().WithField("mysql_driver", "kill_process")
	killConn, err := executor.NewExecutor(logEntry, i.inst, i.inst.DatabaseName)
	if err != nil {
		return err
	}
	defer killConn.Db.Close()
	killSQL := fmt.Sprintf("KILL %v", connID)
	err = util.KillProcess(ctx, killSQL, killConn, logEntry)
	return err
}

func (i *MysqlDriverImpl) query(ctx context.Context, query string, args ...interface{}) ([]map[string]sql.NullString, error) {
	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}
	return conn.Db.Query(query, args...)
}

func (i *MysqlDriverImpl) Parse(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
	nodes, err := i.ParseSql(sqlText)
	if err != nil {
		return nil, err
	}

	lowerCaseTableNames, err := i.Ctx.GetSystemVariable(session.SysVarLowerCaseTableNames)
	if err != nil {
		return nil, err
	}

	ns := make([]driverV2.Node, len(nodes))
	for idx := range nodes {
		n := driverV2.Node{}
		fingerprint, err := util.Fingerprint(nodes[idx].Text(), lowerCaseTableNames == "0")
		if err != nil {
			return nil, err
		}
		n.Fingerprint = fingerprint
		n.Text = nodes[idx].Text()
		n.StartLine = uint64(nodes[idx].StartLine())
		n.Type = i.assertSQLType(nodes[idx])

		ns[idx] = n
	}
	return ns, nil
}

func (i *MysqlDriverImpl) assertSQLType(stmt ast.Node) string {
	switch stmt.(type) {
	case ast.DMLNode:
		switch stmt.(type) {
		case *ast.SelectStmt, *ast.UnionStmt:
			return driverV2.SQLTypeDQL
		default:
			return driverV2.SQLTypeDML
		}
	default:
		return driverV2.SQLTypeDDL
	}
}

func (i *MysqlDriverImpl) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	for _, sql := range sqls {
		if sql == "" {
			return nil, errors.New("has empty sql")
		}
	}
	results := make([]*driverV2.AuditResults, 0, len(sqls))
	for _, sql := range sqls {
		result, err := i.audit(ctx, sql)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (i *MysqlDriverImpl) audit(ctx context.Context, sql string) (*driverV2.AuditResults, error) {
	i.result = driverV2.NewAuditResults()

	nodes, err := i.ParseSql(sql)
	if err != nil {
		return nil, err
	}

	if i.IsOfflineAudit() || i.IsExecutedSQL() {
		err = i.CheckInvalidOffline(nodes[0])
	} else {
		err = i.CheckInvalid(nodes[0])
	}
	if err != nil && session.IsParseShowCreateTableContentErr(err) {
		i.Logger().Errorf("check invalid failed: %v", err)
		i.result.Add(driverV2.RuleLevelWarn, CheckInvalidError, fmt.Sprintf(CheckInvalidErrorFormat, "解析建表语句失败，部分在线审核规则可能失效，请人工确认"))
	} else if err != nil {
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

	var ghostRule *driverV2.Rule
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
			// todo #1630 临时跳过解析建表语句失败导致的规则
			if session.IsParseShowCreateTableContentErr(err) {
				i.Logger().Errorf("skip rule, rule_desc_name=%v rule_desc=%v err:%v", rule.Name, rule.Desc, err.Error())
				continue
			}
			return nil, err
		}
	}

	if i.cnf.optimizeIndexEnabled {
		params := params.Params{
			{
				Key:   MAX_INDEX_COLUMN,
				Value: fmt.Sprint(i.cnf.compositeIndexMaxColumn),
				Type:  params.ParamTypeInt,
			}, {
				Key:   MIN_COLUMN_SELECTIVITY,
				Value: fmt.Sprint(i.cnf.indexSelectivityMinValue),
				Type:  params.ParamTypeFloat64,
			},
		}
		results := optimize(i.log, i.Ctx, nodes[0], params)
		for _, advice := range results {
			i.result.Add(
				driverV2.RuleLevelNotice,
				rulepkg.ConfigOptimizeIndexEnabled,
				advice.Reason,
			)
		}

	}

	// dry run gh-ost
	useGhost, err := i.onlineddlWithGhost(sql)
	if err != nil {
		return nil, errors.Wrap(err, "check whether use ghost or not")
	}
	if useGhost {
		if _, err := i.executeByGhost(ctx, sql, true); err != nil {
			i.result.Add(driverV2.RuleLevelError, ghostRule.Name, fmt.Sprintf("表空间大小超过%vMB, 将使用gh-ost进行上线, 但是dry-run抛出如下错误: %v", i.cnf.DDLGhostMinSize, err))
		} else {
			i.result.Add(ghostRule.Level, ghostRule.Name, fmt.Sprintf("表空间大小超过%vMB, 将使用gh-ost进行上线", i.cnf.DDLGhostMinSize))
		}
	}

	// print osc
	oscCommandLine, err := i.generateOSCCommandLine(nodes[0])
	if err != nil && session.IsParseShowCreateTableContentErr(err) {
		i.Logger().Errorf("generate osc command failed: %v", err.Error()) // todo #1630 临时跳过创表语句解析错误
	} else if err != nil {
		return nil, err
	}
	if oscCommandLine != "" {
		i.result.Add(driverV2.RuleLevelNotice, rulepkg.ConfigDDLOSCMinSize, fmt.Sprintf("[osc]%s", oscCommandLine))
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

func (i *MysqlDriverImpl) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	if i.IsOfflineAudit() {
		return nil, nil
	}

	conn, err := i.getDbConn()
	if err != nil {
		return nil, err
	}

	num, err := util.GetAffectedRowNum(ctx, sql, conn)
	if err != nil && errors.Is(err, util.ErrUnsupportedSqlType) {
		return &driverV2.EstimatedAffectRows{ErrMessage: err.Error()}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get affected row num failed: %w", err)
	}

	return &driverV2.EstimatedAffectRows{
		Count: num,
	}, nil
}

type Config struct {
	DMLRollbackMaxRows int64
	DDLOSCMinSize      int64
	DDLGhostMinSize    int64

	optimizeIndexEnabled     bool
	dmlExplainPreCheckEnable bool
	compositeIndexMaxColumn  int
	indexSelectivityMinValue float64
	isExecutedSQL            bool
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

func (i *MysqlDriverImpl) GetConn() *executor.Executor {
	return i.dbConn
}

func (i *MysqlDriverImpl) GetDSN() *driverV2.DSN {
	return i.inst
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

type PluginProcessor struct{}

func (p *PluginProcessor) GetDriverMetas() (*driverV2.DriverMetas, error) {
	if err := LoadPtTemplateFromFile("./scripts/pt-online-schema-change.template"); err != nil {
		panic(err)
	}
	allRules := make([]*driverV2.Rule, len(rulepkg.RuleHandlers))
	for i := range rulepkg.RuleHandlers {
		allRules[i] = &rulepkg.RuleHandlers[i].Rule
	}
	return &driverV2.DriverMetas{
		PluginName:               driverV2.DriverTypeMySQL,
		DatabaseDefaultPort:      3306,
		Logo:                     logo,
		Rules:                    allRules,
		DatabaseAdditionalParams: params.Params{},
		EnabledOptionalModule: []driverV2.OptionalModule{
			driverV2.OptionalModuleGenRollbackSQL,
			driverV2.OptionalModuleQuery,
			driverV2.OptionalModuleExplain,
			driverV2.OptionalModuleGetTableMeta,
			driverV2.OptionalModuleExtractTableFromSQL,
			driverV2.OptionalModuleEstimateSQLAffectRows,
			driverV2.OptionalModuleKillProcess,
		},
	}, nil
}

func (p *PluginProcessor) Open(l *logrus.Entry, cfg *driverV2.Config) (driver.Plugin, error) {
	return NewInspect(l, cfg)
}

func (p *PluginProcessor) Stop() error {
	return nil
}

func init() {
	driver.BuiltInPluginProcessors[driverV2.DriverTypeMySQL] = &PluginProcessor{}
}
