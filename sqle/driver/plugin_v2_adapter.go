package driver

import (
	"context"
	sqlDriver "database/sql/driver"
	"fmt"

	v2 "github.com/actiontech/sqle/sqle/driver/v2"
	protoV2 "github.com/actiontech/sqle/sqle/driver/v2/proto"
	"github.com/actiontech/sqle/sqle/pkg/params"

	goPlugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc/status"
)

type PluginBootV2 struct {
	client *goPlugin.Client
}

func (d *PluginBootV2) Register() (*v2.DriverMetas, error) {
	cp, err := d.client.Client()
	if err != nil {
		return nil, err
	}
	rawI, err := cp.Dispense(v2.PluginSetName)
	if err != nil {
		return nil, err //todo
	}
	//nolint:forcetypeassert
	s := rawI.(protoV2.DriverClient)

	result, err := s.Metas(context.TODO(), &protoV2.Empty{})
	if err != nil {
		return nil, err
	}

	rules := make([]*v2.Rule, 0, len(result.Rules))
	for _, r := range result.Rules {
		rules = append(rules, v2.ConvertRuleFromProtoToDriver(r))
	}

	ms := make([]v2.OptionalModule, 0, len(result.EnabledOptionalModule))
	for _, m := range result.EnabledOptionalModule {
		ms = append(ms, v2.OptionalModule(m))
	}

	return &v2.DriverMetas{
		PluginName:               result.PluginName,
		DatabaseDefaultPort:      result.DatabaseDefaultPort,
		DatabaseAdditionalParams: v2.ConvertProtoParamToParam(result.DatabaseAdditionalParams),
		Rules:                    rules,
		EnabledOptionalModule:    ms,
	}, nil
}

func (d *PluginBootV2) Open(cfgV2 *v2.Config) (Plugin, error) {
	cp, err := d.client.Client()
	if err != nil {
		return nil, err
	}
	rawI, err := cp.Dispense(v2.PluginSetName)
	if err != nil {
		return nil, err //todo
	}
	//nolint:forcetypeassert
	s := rawI.(protoV2.DriverClient)
	var dsn *protoV2.DSN
	if cfgV2.DSN != nil {
		dsn = &protoV2.DSN{
			Host:             cfgV2.DSN.Host,
			Port:             cfgV2.DSN.Port,
			User:             cfgV2.DSN.User,
			Password:         cfgV2.DSN.Password,
			Database:         cfgV2.DSN.DatabaseName,
			AdditionalParams: v2.ConvertParamToProtoParam(cfgV2.DSN.AdditionalParams),
		}
	}

	var rules = make([]*protoV2.Rule, 0, len(cfgV2.Rules))
	for _, rule := range cfgV2.Rules {
		rules = append(rules, v2.ConvertRuleFromDriverToProto(rule))
	}

	result, err := s.Init(context.TODO(), &protoV2.InitRequest{
		Dsn:   dsn,
		Rules: rules,
	})
	if err != nil {
		return nil, err //todo
	}
	return &PluginImplV2{
		client:  s,
		Session: result.Session,
	}, nil
}

func (d *PluginBootV2) Stop() error {
	if d.client != nil {
		d.client.Kill()
	}
	return nil
}

type PluginImplV2 struct {
	client  protoV2.DriverClient
	Session *protoV2.Session
}

func (s *PluginImplV2) Close(ctx context.Context) {
	s.client.Close(ctx, &protoV2.CloseResquest{
		Session: s.Session,
	})
}

// audit

func (s *PluginImplV2) Parse(ctx context.Context, sqlText string) ([]v2.Node, error) {
	resp, err := s.client.Parse(ctx, &protoV2.ParseRequest{
		Session: s.Session,
		Sql: &protoV2.ParsedSQL{
			Query: sqlText,
		}})
	if err != nil {
		return nil, err
	}

	nodes := make([]v2.Node, len(resp.Nodes))
	for i, node := range resp.Nodes {
		nodes[i] = v2.Node{
			Type:        node.Type,
			Text:        node.Text,
			Fingerprint: node.Fingerprint,
		}
	}
	return nodes, nil
}

func (s *PluginImplV2) Audit(ctx context.Context, sqls []string) ([]*v2.AuditResults, error) {
	auditSqls := make([]*protoV2.AuditSQL, 0, len(sqls))
	for _, sql := range sqls {
		auditSqls = append(auditSqls, &protoV2.AuditSQL{Query: sql})
	}
	resp, err := s.client.Audit(ctx, &protoV2.AuditRequest{
		Session: s.Session,
		Sqls:    auditSqls,
	})
	if err != nil {
		return nil, err
	}

	rets := []*v2.AuditResults{}
	for _, results := range resp.AuditResults {
		ret := &v2.AuditResults{}
		for _, result := range results.Results {
			ret.Results = append(ret.Results, &v2.AuditResult{
				Level:   v2.RuleLevel(result.Level),
				Message: result.Message,
			})
		}
		rets = append(rets, ret)
	}
	return rets, nil
}

func (s *PluginImplV2) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	resp, err := s.client.GenRollbackSQL(ctx, &protoV2.GenRollbackSQLRequest{
		Session: s.Session,
		Sql: &protoV2.NeedRollbackSQL{
			Query: sql,
		},
	})
	if err != nil {
		return "", "", err
	}

	return resp.Sql.Query, resp.Sql.Message, nil
}

// executor

func (s *PluginImplV2) Ping(ctx context.Context) error {
	_, err := s.client.Ping(ctx, &protoV2.PingRequest{
		Session: s.Session,
	})
	return err
}

func (s *PluginImplV2) Exec(ctx context.Context, sql string) (sqlDriver.Result, error) {
	resp, err := s.client.Exec(ctx, &protoV2.ExecRequest{
		Session: s.Session,
		Sql:     &protoV2.ExecSQL{Query: sql},
	})
	if err != nil {
		return nil, err
	}
	return &dbDriverResult{
		lastInsertId:    resp.Result.LastInsertId,
		lastInsertIdErr: resp.Result.LastInsertIdError,
		rowsAffected:    resp.Result.RowsAffected,
		rowsAffectedErr: resp.Result.RowsAffectedError,
	}, nil
}

func (s *PluginImplV2) Tx(ctx context.Context, sqls ...string) ([]sqlDriver.Result, error) {
	execSqls := make([]*protoV2.ExecSQL, 0, len(sqls))
	for _, sql := range sqls {
		execSqls = append(execSqls, &protoV2.ExecSQL{Query: sql})
	}
	resp, err := s.client.Tx(ctx, &protoV2.TxRequest{
		Session: s.Session,
		Sqls:    execSqls,
	})
	if err != nil {
		return nil, err
	}

	ret := make([]sqlDriver.Result, len(resp.Results))
	for i, result := range resp.Results {
		ret[i] = &dbDriverResult{
			lastInsertId:    result.LastInsertId,
			lastInsertIdErr: result.LastInsertIdError,
			rowsAffected:    result.RowsAffected,
			rowsAffectedErr: result.RowsAffectedError,
		}
	}
	return ret, nil
}

func (s *PluginImplV2) Query(ctx context.Context, sql string, conf *v2.QueryConf) (*v2.QueryResult, error) {
	req := &protoV2.QueryRequest{
		Session: s.Session,
		Sql: &protoV2.QuerySQL{
			Query: sql,
		},
		Conf: &protoV2.QueryConf{
			TimeoutSecond: conf.TimeOutSecond,
		},
	}
	res, err := s.client.Query(ctx, req)
	if err != nil {
		return nil, err
	}
	result := &v2.QueryResult{
		Column: params.Params{},
		Rows:   []*v2.QueryResultRow{},
	}
	for _, p := range res.GetColumn() {
		result.Column = append(result.Column, &params.Param{
			Key:   p.GetKey(),
			Value: p.GetValue(),
			Desc:  p.GetDesc(),
			Type:  params.ParamType(p.GetType()),
		})
	}
	for _, row := range res.GetRows() {
		r := &v2.QueryResultRow{
			Values: []*v2.QueryResultValue{},
		}
		for _, value := range row.GetValues() {
			r.Values = append(r.Values, &v2.QueryResultValue{
				Value: value.GetValue(),
			})
		}
		result.Rows = append(result.Rows, r)
	}
	return result, nil
}

func (s *PluginImplV2) Explain(ctx context.Context, conf *v2.ExplainConf) (*v2.ExplainResult, error) {
	req := &protoV2.ExplainRequest{
		Session: s.Session,
		Sql: &protoV2.ExplainSQL{
			Query: conf.Sql,
		},
	}
	res, err := s.client.Explain(ctx, req)
	if err != nil && status.Code(err) == v2.GrpcErrSQLIsNotSupported {
		return nil, v2.ErrSQLIsNotSupported
	} else if err != nil {
		return nil, err
	}

	return &v2.ExplainResult{
		ClassicResult: v2.ExplainClassicResult{
			TabularData: v2.ConvertProtoTabularDataToDriver(res.ClassicResult.Data),
		},
	}, nil
}

// metadata

func (s *PluginImplV2) Schemas(ctx context.Context) ([]string, error) {
	resp, err := s.client.GetDatabases(ctx, &protoV2.GetDatabasesRequest{
		Session: s.Session,
	})
	if err != nil {
		return nil, err
	}
	databases := make([]string, 0, len(resp.Databases))
	for _, d := range resp.Databases {
		databases = append(databases, d.Name)
	}
	return databases, nil
}

func (s *PluginImplV2) getTableMeta(ctx context.Context, table *v2.Table) (*v2.TableMeta, error) {
	result, err := s.client.GetTableMeta(ctx, &protoV2.GetTableMetaRequest{
		Session: s.Session,
		Table: &protoV2.Table{
			Name:   table.Name,
			Schema: table.Schema,
		},
	})
	if err != nil {
		return nil, err
	}
	return v2.ConvertProtoTableMetaToDriver(result.TableMeta), nil
}

func (s *PluginImplV2) extractTableFromSQL(ctx context.Context, sql string) ([]*v2.Table, error) {
	result, err := s.client.ExtractTableFromSQL(ctx, &protoV2.ExtractTableFromSQLRequest{
		Session: s.Session,
		Sql:     &protoV2.ExtractedSQL{Query: sql},
	})
	if err != nil {
		return nil, err
	}
	tables := make([]*v2.Table, 0, len(result.Tables))
	for _, table := range result.Tables {
		tables = append(tables, &v2.Table{
			Name:   table.Name,
			Schema: table.Schema,
		})
	}
	return tables, nil
}

func (s *PluginImplV2) GetTableMetaBySQL(ctx context.Context, conf *GetTableMetaBySQLConf) (*GetTableMetaBySQLResult, error) {
	tables, err := s.extractTableFromSQL(ctx, conf.Sql)
	if err != nil {
		return nil, err
	}
	tableMetas := make([]*TableMeta, 0, len(tables))
	for _, table := range tables {
		tableMeta, err := s.getTableMeta(ctx, table)
		if err != nil {
			return nil, err
		}
		tableMetas = append(tableMetas, &TableMeta{
			Table:     *table,
			TableMeta: *tableMeta,
		})
	}
	return &GetTableMetaBySQLResult{TableMetas: tableMetas}, nil
}

type dbDriverResult struct {
	lastInsertId    int64
	lastInsertIdErr string
	rowsAffected    int64
	rowsAffectedErr string
}

func (s *dbDriverResult) LastInsertId() (int64, error) {
	if s.lastInsertIdErr != "" {
		return s.lastInsertId, fmt.Errorf(s.lastInsertIdErr)
	}
	return s.lastInsertId, nil
}

func (s *dbDriverResult) RowsAffected() (int64, error) {
	if s.rowsAffectedErr != "" {
		return s.rowsAffected, fmt.Errorf(s.rowsAffectedErr)
	}
	return s.rowsAffected, nil
}
