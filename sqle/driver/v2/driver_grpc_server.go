package driverV2

import (
	"context"
	"fmt"
	"sync"

	protoV2 "github.com/actiontech/sqle/sqle/driver/v2/proto"
	"github.com/actiontech/sqle/sqle/pkg/params"

	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

var (
	ErrNodesCountExceedOne = errors.New("after parse, nodes count exceed one")
	ErrSQLIsNotSupported   = errors.New("SQL is not supported")
	ErrSQLisEmpty          = errors.New("SQL is empty")
)

type DSN struct {
	Host             string
	Port             string
	User             string
	Password         string
	AdditionalParams params.Params

	// DatabaseName is the default database to connect.
	DatabaseName string
}

type Rule struct {
	Name       string
	Desc       string
	Annotation string

	// Category is the category of the rule. Such as "Naming Conventions"...
	// Rules will be displayed on the SQLE rule list page by category.
	Category  string
	Level     RuleLevel
	Params    params.Params
	Knowledge RuleKnowledge
}

type Config struct {
	DSN   *DSN
	Rules []*Rule
}

func NewConfig(dsn *DSN, rules []*Rule) (*Config, error) {
	return &Config{
		DSN:   dsn,
		Rules: rules,
	}, nil
}

type DriverGrpcServer struct {
	Meta          DriverMetas
	DriverFactory func(*Config) (Driver, error)
	Drivers       map[string] /* session id*/ Driver
	sync.Mutex
}

func (d *DriverGrpcServer) getDriverBySession(session *protoV2.Session) (Driver, error) {
	d.Lock()
	defer d.Unlock()
	driver, ok := d.Drivers[session.Id]
	if !ok {
		return nil, fmt.Errorf("session %s not found", session.Id)
	}
	return driver, nil
}

func (d *DriverGrpcServer) Metas(ctx context.Context, req *protoV2.Empty) (*protoV2.MetasResponse, error) {
	rules := make([]*protoV2.Rule, 0, len(d.Meta.Rules))
	for _, r := range d.Meta.Rules {
		rules = append(rules, ConvertRuleFromDriverToProto(r))
	}

	ms := make([]protoV2.OptionalModule, 0, len(d.Meta.EnabledOptionalModule))
	for _, m := range d.Meta.EnabledOptionalModule {
		ms = append(ms, protoV2.OptionalModule(m))
	}

	return &protoV2.MetasResponse{
		PluginName:               d.Meta.PluginName,
		DatabaseDefaultPort:      d.Meta.DatabaseDefaultPort,
		Logo:                     d.Meta.Logo,
		DatabaseAdditionalParams: ConvertParamToProtoParam(d.Meta.DatabaseAdditionalParams),
		Rules:                    rules,
		EnabledOptionalModule:    ms,
	}, nil
}

func (d *DriverGrpcServer) Init(ctx context.Context, req *protoV2.InitRequest) (*protoV2.InitResponse, error) {
	var rules = make([]*Rule, 0, len(req.GetRules()))
	for _, rule := range req.GetRules() {
		rules = append(rules, ConvertRuleFromProtoToDriver(rule))
	}

	var dsn *DSN
	if req.GetDsn() != nil {
		dsn = &DSN{
			Host:             req.GetDsn().GetHost(),
			Port:             req.GetDsn().GetPort(),
			User:             req.GetDsn().GetUser(),
			Password:         req.GetDsn().GetPassword(),
			DatabaseName:     req.GetDsn().GetDatabase(),
			AdditionalParams: ConvertProtoParamToParam(req.GetDsn().GetAdditionalParams()),
		}
	}
	id := RandStr(20)
	driver, err := d.DriverFactory(&Config{
		DSN:   dsn,
		Rules: rules,
	})
	if err != nil {
		return nil, errors.Wrap(err, "init config")
	}
	d.Lock()
	if _, ok := d.Drivers[id]; ok {
		d.Unlock()
		driver.Close(ctx)
		return nil, errors.New("session id is duplicated") // 几乎不会发生.
	}
	d.Drivers[id] = driver
	d.Unlock()

	return &protoV2.InitResponse{
		Session: &protoV2.Session{
			Id: id,
		},
	}, nil
}

func (d *DriverGrpcServer) Close(ctx context.Context, req *protoV2.CloseRequest) (*protoV2.Empty, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.Empty{}, err
	}
	d.Mutex.Lock()
	delete(d.Drivers, req.Session.Id)
	d.Mutex.Unlock()

	driver.Close(ctx)
	return &protoV2.Empty{}, nil
}

func (d *DriverGrpcServer) Parse(ctx context.Context, req *protoV2.ParseRequest) (*protoV2.ParseResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.ParseResponse{}, err
	}
	if req.Sql == nil {
		return &protoV2.ParseResponse{}, ErrSQLisEmpty
	}
	nodes, err := driver.Parse(ctx, req.Sql.Query)
	if err != nil {
		return &protoV2.ParseResponse{}, err
	}

	resp := &protoV2.ParseResponse{}
	for _, node := range nodes {
		resp.Nodes = append(resp.Nodes, &protoV2.Node{
			Text:        node.Text,
			Type:        node.Type,
			Fingerprint: node.Fingerprint,
		})
	}
	return resp, nil
}

func (d *DriverGrpcServer) Audit(ctx context.Context, req *protoV2.AuditRequest) (*protoV2.AuditResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.AuditResponse{}, err
	}
	if len(req.Sqls) == 0 {
		return &protoV2.AuditResponse{}, ErrSQLisEmpty
	}
	sqls := make([]string, 0, len(req.Sqls))
	for _, sql := range req.Sqls {
		sqls = append(sqls, sql.Query)
	}

	auditResults, err := driver.Audit(ctx, sqls)
	if err != nil {
		return &protoV2.AuditResponse{}, err
	}

	resp := &protoV2.AuditResponse{}
	for _, results := range auditResults {
		rets := &protoV2.AuditResults{
			Results: []*protoV2.AuditResult{},
		}
		for _, result := range results.Results {
			rets.Results = append(rets.Results, &protoV2.AuditResult{
				Level:    string(result.Level),
				Message:  result.Message,
				RuleName: result.RuleName,
			})
		}
		resp.AuditResults = append(resp.AuditResults, rets)
	}
	return resp, nil
}

func (d *DriverGrpcServer) GenRollbackSQL(ctx context.Context, req *protoV2.GenRollbackSQLRequest) (*protoV2.GenRollbackSQLResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.GenRollbackSQLResponse{}, err
	}
	if req.Sql == nil {
		return &protoV2.GenRollbackSQLResponse{}, ErrSQLisEmpty
	}
	rollbackSQL, reason, err := driver.GenRollbackSQL(ctx, req.Sql.Query)
	return &protoV2.GenRollbackSQLResponse{
		Sql: &protoV2.RollbackSQL{
			Query:   rollbackSQL,
			Message: reason,
		},
	}, err
}

func (d *DriverGrpcServer) Ping(ctx context.Context, req *protoV2.PingRequest) (*protoV2.Empty, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.Empty{}, err
	}
	return &protoV2.Empty{}, driver.Ping(ctx)
}

func (d *DriverGrpcServer) Exec(ctx context.Context, req *protoV2.ExecRequest) (*protoV2.ExecResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.ExecResponse{}, err
	}
	if req.Sql == nil {
		return &protoV2.ExecResponse{}, ErrSQLisEmpty
	}
	result, err := driver.Exec(ctx, req.Sql.Query)
	if err != nil {
		return &protoV2.ExecResponse{}, err
	}

	execResult := &protoV2.ExecResult{}
	lastInsertId, lastInsertIdErr := result.LastInsertId()
	execResult.LastInsertId = lastInsertId
	if lastInsertIdErr != nil {
		execResult.LastInsertIdError = lastInsertIdErr.Error()
	}
	rowsAffected, rowsAffectedErr := result.RowsAffected()
	execResult.RowsAffected = rowsAffected
	if rowsAffectedErr != nil {
		execResult.RowsAffectedError = rowsAffectedErr.Error()
	}
	return &protoV2.ExecResponse{Result: execResult}, nil
}

func (d *DriverGrpcServer) Tx(ctx context.Context, req *protoV2.TxRequest) (*protoV2.TxResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.TxResponse{}, err
	}
	if len(req.Sqls) == 0 {
		return &protoV2.TxResponse{}, ErrSQLisEmpty
	}
	sqls := make([]string, 0, len(req.Sqls))
	for _, sql := range req.Sqls {
		sqls = append(sqls, sql.Query)
	}

	results, err := driver.Tx(ctx, sqls...)
	if err != nil {
		return &protoV2.TxResponse{}, err
	}

	txResults := make([]*protoV2.ExecResult, 0, len(results))
	for _, result := range results {
		txResult := &protoV2.ExecResult{}

		lastInsertId, lastInsertIdErr := result.LastInsertId()
		txResult.LastInsertId = lastInsertId
		if lastInsertIdErr != nil {
			txResult.LastInsertIdError = lastInsertIdErr.Error()
		}
		rowsAffected, rowsAffectedErr := result.RowsAffected()
		txResult.RowsAffected = rowsAffected
		if rowsAffectedErr != nil {
			txResult.RowsAffectedError = rowsAffectedErr.Error()
		}

		txResults = append(txResults, txResult)
	}
	return &protoV2.TxResponse{Results: txResults}, nil
}

func (d *DriverGrpcServer) Query(ctx context.Context, req *protoV2.QueryRequest) (*protoV2.QueryResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.QueryResponse{}, err
	}
	if req.Sql == nil {
		return &protoV2.QueryResponse{}, ErrSQLisEmpty
	}
	conf := &QueryConf{
		TimeOutSecond: req.GetConf().GetTimeoutSecond(),
	}
	res, err := driver.Query(ctx, req.Sql.Query, conf)
	if err != nil {
		return &protoV2.QueryResponse{}, err
	}
	resp := &protoV2.QueryResponse{
		Column: []*protoV2.Param{},
		Rows:   []*protoV2.QueryResultRow{},
	}
	for _, param := range res.Column {
		resp.Column = append(resp.Column, &protoV2.Param{
			Key:   param.Key,
			Value: param.Value,
			Desc:  param.Desc,
			Type:  string(param.Type),
		})
	}
	for _, row := range res.Rows {
		rows := &protoV2.QueryResultRow{
			Values: []*protoV2.QueryResultValue{},
		}
		for _, value := range row.Values {
			rows.Values = append(rows.Values, &protoV2.QueryResultValue{
				Value: value.Value,
			})
		}
		resp.Rows = append(resp.Rows, rows)
	}
	return resp, nil
}

func (d *DriverGrpcServer) Explain(ctx context.Context, req *protoV2.ExplainRequest) (*protoV2.ExplainResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.ExplainResponse{}, err
	}
	if req.Sql == nil {
		return &protoV2.ExplainResponse{}, ErrSQLisEmpty
	}
	conf := &ExplainConf{
		Sql: req.Sql.Query,
	}
	res, err := driver.Explain(ctx, conf)
	if err != nil && err == ErrSQLIsNotSupported {
		return &protoV2.ExplainResponse{}, status.Error(GrpcErrSQLIsNotSupported, err.Error())
	} else if err != nil {
		return &protoV2.ExplainResponse{}, err
	}

	resp := &protoV2.ExplainResponse{
		ClassicResult: &protoV2.ExplainClassicResult{
			Data: ConvertTabularDataToProto(res.ClassicResult.TabularData),
		},
	}
	return resp, nil
}

func (d *DriverGrpcServer) GetDatabases(ctx context.Context, req *protoV2.GetDatabasesRequest) (*protoV2.GetDatabasesResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.GetDatabasesResponse{}, err
	}
	databases, err := driver.GetDatabases(ctx)
	protoDatabases := make([]*protoV2.Database, 0, len(databases))
	for _, d := range databases {
		protoDatabases = append(protoDatabases, &protoV2.Database{Name: d})
	}
	return &protoV2.GetDatabasesResponse{Databases: protoDatabases}, err
}

func (d *DriverGrpcServer) GetTableMeta(ctx context.Context, req *protoV2.GetTableMetaRequest) (*protoV2.GetTableMetaResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.GetTableMetaResponse{}, err
	}
	if req.Table == nil && req.Table.Name == "" {
		return &protoV2.GetTableMetaResponse{}, errors.New("table name is empty")
	}

	tm, err := driver.GetTableMeta(ctx, &Table{Name: req.Table.Name, Schema: req.Table.Schema})
	if err != nil {
		return &protoV2.GetTableMetaResponse{}, err
	}

	return &protoV2.GetTableMetaResponse{
		TableMeta: ConvertTableMetaToProto(tm),
	}, nil
}

func (d *DriverGrpcServer) ExtractTableFromSQL(ctx context.Context, req *protoV2.ExtractTableFromSQLRequest) (*protoV2.ExtractTableFromSQLResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.ExtractTableFromSQLResponse{}, err
	}
	if req.Sql == nil {
		return &protoV2.ExtractTableFromSQLResponse{}, ErrSQLisEmpty
	}
	tables, err := driver.ExtractTableFromSQL(ctx, req.Sql.Query)
	if err != nil {
		return &protoV2.ExtractTableFromSQLResponse{}, err
	}
	protoTables := make([]*protoV2.Table, 0, len(tables))
	for _, t := range tables {
		protoTables = append(protoTables, &protoV2.Table{
			Name:   t.Name,
			Schema: t.Schema,
		})
	}
	return &protoV2.ExtractTableFromSQLResponse{
		Tables: protoTables,
	}, nil
}

func (d *DriverGrpcServer) EstimateSQLAffectRows(ctx context.Context, req *protoV2.EstimateSQLAffectRowsRequest) (*protoV2.EstimateSQLAffectRowsResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.EstimateSQLAffectRowsResponse{}, err
	}
	if req.Sql == nil {
		return &protoV2.EstimateSQLAffectRowsResponse{}, ErrSQLisEmpty
	}
	ar, err := driver.EstimateSQLAffectRows(ctx, req.Sql.Query)
	if err != nil {
		return &protoV2.EstimateSQLAffectRowsResponse{}, err
	}
	return &protoV2.EstimateSQLAffectRowsResponse{
		Count:      ar.Count,
		ErrMessage: ar.ErrMessage,
	}, nil
}

func (d *DriverGrpcServer) KillProcess(ctx context.Context, req *protoV2.KillProcessRequest) (*protoV2.KillProcessResponse, error) {
	driver, err := d.getDriverBySession(req.Session)
	if err != nil {
		return &protoV2.KillProcessResponse{}, err
	}
	info, err := driver.KillProcess(ctx)
	if err != nil {
		return &protoV2.KillProcessResponse{}, err
	}
	return &protoV2.KillProcessResponse{
		ErrMessage: info.ErrMessage,
	}, nil
}
