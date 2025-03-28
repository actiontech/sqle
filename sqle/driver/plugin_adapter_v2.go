package driver

import (
	"context"
	sqlDriver "database/sql/driver"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	protoV2 "github.com/actiontech/sqle/sqle/driver/v2/proto"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"golang.org/x/text/language"

	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type PluginProcessorV2 struct {
	cfg               func(cmdBase string, cmdArgs []string) *goPlugin.ClientConfig
	cmdBase           string
	cmdArgs           []string
	client            *goPlugin.Client
	meta              *driverV2.DriverMetas
	pluginPidFilePath string
	sync.Mutex
}

func (d *PluginProcessorV2) getDriverClient(l *logrus.Entry) (protoV2.DriverClient, error) {
	var client *goPlugin.Client

	d.Lock()
	if d.client.Exited() {
		l.Infof("plugin process is exited, restart it")
		newClient := goPlugin.NewClient(d.cfg(d.cmdBase, d.cmdArgs))
		_, err := newClient.Client()
		if err != nil {
			d.Unlock()
			return nil, err
		}
		l.Infof("restart plugin success")
		d.client.Kill()
		d.client = newClient
	}

	client = d.client
	d.Unlock()

	cp, err := client.Client()
	if err != nil {
		return nil, err
	}
	rawI, err := cp.Dispense(driverV2.PluginSetName)
	if err != nil {
		return nil, err
	}
	//nolint:forcetypeassert
	s, ok := rawI.(protoV2.DriverClient)
	if !ok {
		return nil, fmt.Errorf("client is not implement protoV2.DriverClient")
	}
	return s, nil
}

func (d *PluginProcessorV2) GetDriverMetas() (*driverV2.DriverMetas, error) {
	c, err := d.getDriverClient(log.NewEntry())
	if err != nil {
		return nil, err
	}

	result, err := c.Metas(context.TODO(), &protoV2.Empty{})
	if err != nil {
		return nil, err
	}

	var isI18n bool
	ms := make([]driverV2.OptionalModule, 0, len(result.EnabledOptionalModule))
	for _, m := range result.EnabledOptionalModule {
		ms = append(ms, driverV2.OptionalModule(m))
		if driverV2.OptionalModule(m) == driverV2.OptionalModuleI18n {
			isI18n = true
		}
	}

	rules := make([]*driverV2.Rule, 0, len(result.Rules))
	ruleVersionIncludedMap := make(map[uint32]struct{})
	var ruleVersinIncluded []uint32
	for _, r := range result.Rules {
		if len(r.I18NRuleInfo) > 0 {
			if _, exist := r.I18NRuleInfo[i18nPkg.DefaultLang.String()]; !exist {
				// 多语言插件必须支持 i18nPkg.DefaultLang 用以默认展示
				return nil, fmt.Errorf("client rule: %s not support language: %s", r.Name, i18nPkg.DefaultLang.String())
			}
		}
		dr, err := driverV2.ConvertI18nRuleFromProtoToDriver(r, result.PluginName, isI18n)
		if err != nil {
			return nil, err
		}
		rules = append(rules, dr)
		if _, exist := ruleVersionIncludedMap[dr.Version]; !exist {
			ruleVersionIncludedMap[dr.Version] = struct{}{}
			ruleVersinIncluded = append(ruleVersinIncluded, dr.Version)
		}
	}

	ps, err := driverV2.ConvertProtoParamToParam(result.DatabaseAdditionalParams)
	if err != nil {
		return nil, fmt.Errorf("plugin Metas rule param err: %w", err)
	}
	meta := &driverV2.DriverMetas{
		PluginName:               result.PluginName,
		DatabaseDefaultPort:      result.DatabaseDefaultPort,
		Logo:                     result.Logo,
		DatabaseAdditionalParams: ps,
		Rules:                    rules,
		RuleVersionIncluded:      ruleVersinIncluded,
		EnabledOptionalModule:    ms,
	}
	d.meta = meta
	return meta, nil
}

func (d *PluginProcessorV2) Open(l *logrus.Entry, cfgV2 *driverV2.Config) (Plugin, error) {
	l = l.WithFields(logrus.Fields{
		"plugin":         d.meta.PluginName,
		"plugin_version": driverV2.ProtocolVersion,
	})
	c, err := d.getDriverClient(l)
	if err != nil {
		return nil, err
	}

	var dsn *protoV2.DSN
	if cfgV2.DSN != nil {
		dsn = &protoV2.DSN{
			Host:             cfgV2.DSN.Host,
			Port:             cfgV2.DSN.Port,
			User:             cfgV2.DSN.User,
			Password:         cfgV2.DSN.Password,
			Database:         cfgV2.DSN.DatabaseName,
			AdditionalParams: driverV2.ConvertParamToProtoParam(cfgV2.DSN.AdditionalParams),
		}
	}

	l.Infof("starting call plugin interface [Init]")
	result, err := c.Init(context.TODO(), &protoV2.InitRequest{
		Dsn:   dsn,
		Rules: driverV2.ConvertI18nRulesFromDriverToProto(cfgV2.Rules),
	})
	if err != nil {
		l.Errorf("fail to call plugin interface [Init], error: %v", err)
		return nil, err
	}
	l.Infof("call plugin interface [Init] success")
	return &PluginImplV2{
		client:  c,
		Session: result.Session,
		l:       l.WithField("session_id", result.Session.Id),
		meta:    d.meta,
	}, nil
}

func (d *PluginProcessorV2) Stop() error {
	d.Lock()
	if d.client != nil {
		d.client.Kill()
	}
	os.Remove(d.pluginPidFilePath)
	d.Unlock()
	return nil
}

type PluginImplV2 struct {
	l       *logrus.Entry
	client  protoV2.DriverClient
	Session *protoV2.Session
	meta    *driverV2.DriverMetas
}

func (s *PluginImplV2) Backup(ctx context.Context, backupStrategy string, sql string, backupMaxRows uint64) (backupSqls []string, executeResult string, err error) {
	api := "Backup"
	s.preLog(api)
	var strategy protoV2.BackupStrategy
	switch backupStrategy {
	case driverV2.BackupStrategyReverseSql:
		strategy = protoV2.BackupStrategy_ReverseSql
	case driverV2.BackupStrategyNone:
		strategy = protoV2.BackupStrategy_None
	case driverV2.BackupStrategyOriginalRow:
		strategy = protoV2.BackupStrategy_OriginalRow
	case driverV2.BackupStrategyManually:
		strategy = protoV2.BackupStrategy_Manually
	default:
		return []string{}, "", fmt.Errorf("unsupported strategy %v", backupStrategy)
	}
	resp, err := s.client.Backup(ctx, &protoV2.BackupReq{
		Session:        s.Session,
		BackupStrategy: strategy,
		Sql:            sql,
		BackupMaxRows:  backupMaxRows,
	})
	s.afterLog(api, err)
	if err != nil {
		return nil, "", err
	}
	return resp.GetBackupSql(), resp.GetExecuteResult(), nil
}

func (p *PluginImplV2) RecommendBackupStrategy(ctx context.Context, sql string) (*RecommendBackupStrategyRes, error) {
	api := "RecommendBackupStrategy"
	p.preLog(api)

	resp, err := p.client.RecommendBackupStrategy(ctx, &protoV2.RecommendBackupStrategyReq{
		Session: p.Session,
		Sql:     sql,
	})
	p.afterLog(api, err)
	if err != nil {
		return nil, err
	}
	return &RecommendBackupStrategyRes{
		BackupStrategy:    resp.BackupStrategy.String(),
		BackupStrategyTip: resp.BackupStrategyTip,
		TablesRefer:       resp.TablesRefer,
		SchemasRefer:      resp.SchemasRefer,
	}, nil
}

func (s *PluginImplV2) preLog(ApiName string) {
	s.l.Infof("starting call plugin interface [%s]", ApiName)
}

func (s *PluginImplV2) afterLog(ApiName string, err error) {
	if err != nil {
		s.l.Errorf("fail to call plugin interface [%s], error: %v", ApiName, err)
	} else {
		s.l.Infof("call plugin interface [%s] success", ApiName)
	}
}

func (s *PluginImplV2) Close(ctx context.Context) {
	api := "Close"
	s.preLog(api)
	_, err := s.client.Close(ctx, &protoV2.CloseRequest{
		Session: s.Session,
	})
	s.afterLog(api, err)
}

func (s *PluginImplV2) KillProcess(ctx context.Context) error {
	api := "Kill Process"
	s.preLog(api)
	rs, err := s.client.KillProcess(ctx, &protoV2.KillProcessRequest{
		Session: s.Session,
	})
	s.afterLog(api, err)
	if err != nil {
		return err
	}
	if rs.ErrMessage != "" {
		return errors.New(rs.ErrMessage)
	}
	return nil
}

// audit

func (s *PluginImplV2) Parse(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
	api := "Parse"
	s.preLog(api)
	resp, err := s.client.Parse(ctx, &protoV2.ParseRequest{
		Session: s.Session,
		Sql: &protoV2.ParsedSQL{
			Query: sqlText,
		}},
	)
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}

	nodes := make([]driverV2.Node, len(resp.Nodes))
	for i, node := range resp.Nodes {
		nodes[i] = driverV2.Node{
			Type:        node.Type,
			Text:        node.Text,
			Fingerprint: node.Fingerprint,
			StartLine:   node.StartLine,
			ExecBatchId: node.BatchId,
		}
	}
	return nodes, nil
}

func (s *PluginImplV2) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	api := "Audit"
	s.preLog(api)
	auditSqls := make([]*protoV2.AuditSQL, 0, len(sqls))
	for _, sql := range sqls {
		auditSqls = append(auditSqls, &protoV2.AuditSQL{Query: sql})
	}
	resp, err := s.client.Audit(ctx, &protoV2.AuditRequest{
		Session: s.Session,
		Sqls:    auditSqls,
	})
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}

	rets := make([]*driverV2.AuditResults, 0, len(resp.AuditResults))
	for _, results := range resp.AuditResults {
		dResult, err := driverV2.ConvertI18nAuditResultsFromProtoToDriver(results.Results, s.meta.IsOptionalModuleEnabled(driverV2.OptionalModuleI18n))
		if err != nil {
			return nil, err
		}
		ret := driverV2.NewAuditResults()
		ret.Results = dResult
		rets = append(rets, ret)
	}
	return rets, nil
}

func (s *PluginImplV2) GenRollbackSQL(ctx context.Context, sql string) (string, i18nPkg.I18nStr, error) {
	api := "GenRollbackSQL"
	s.preLog(api)
	resp, err := s.client.GenRollbackSQL(ctx, &protoV2.GenRollbackSQLRequest{
		Session: s.Session,
		Sql: &protoV2.NeedRollbackSQL{
			Query: sql,
		},
	})
	s.afterLog(api, err)
	if err != nil {
		return "", nil, err
	}

	var i18nReason i18nPkg.I18nStr
	if resp.Sql.Message != "" && len(resp.Sql.I18NRollbackSQLInfo) == 0 {
		i18nReason = i18nPkg.ConvertStr2I18nAsDefaultLang(resp.Sql.Message)
	} else if len(resp.Sql.I18NRollbackSQLInfo) > 0 {
		i18nReason = make(i18nPkg.I18nStr, len(resp.Sql.I18NRollbackSQLInfo))
		for langTag, v := range resp.Sql.I18NRollbackSQLInfo {
			tag, err := language.Parse(langTag)
			if err != nil {
				return "", nil, fmt.Errorf("fail to parse I18NRollbackSQLInfo tag [%s], error: %v", langTag, err)
			}
			i18nReason[tag] = v.Message
		}
	}
	return resp.Sql.Query, i18nReason, nil
}

// executor

func (s *PluginImplV2) Ping(ctx context.Context) error {
	api := "Ping"
	s.preLog(api)
	_, err := s.client.Ping(ctx, &protoV2.PingRequest{
		Session: s.Session,
	})
	s.afterLog(api, err)
	return err
}

func (s *PluginImplV2) Exec(ctx context.Context, sql string) (sqlDriver.Result, error) {
	api := "Exec"
	s.preLog(api)
	resp, err := s.client.Exec(ctx, &protoV2.ExecRequest{
		Session: s.Session,
		Sql:     &protoV2.ExecSQL{Query: sql},
	})
	s.afterLog(api, err)
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

func (s *PluginImplV2) ExecBatch(ctx context.Context, sqls ...string) ([]sqlDriver.Result, error) {
	api := "ExecBatch"
	s.preLog(api)
	execSqls := make([]*protoV2.ExecSQL, 0, len(sqls))
	for _, sql := range sqls {
		execSqls = append(execSqls, &protoV2.ExecSQL{Query: sql})
	}
	resp, err := s.client.ExecBatch(ctx, &protoV2.ExecBatchRequest{
		Session: s.Session,
		Sqls:    execSqls,
	})
	s.afterLog(api, err)
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

func (s *PluginImplV2) Tx(ctx context.Context, sqls ...string) ([]sqlDriver.Result, error) {
	api := "Tx"
	s.preLog(api)
	execSqls := make([]*protoV2.ExecSQL, 0, len(sqls))
	for _, sql := range sqls {
		execSqls = append(execSqls, &protoV2.ExecSQL{Query: sql})
	}
	resp, err := s.client.Tx(ctx, &protoV2.TxRequest{
		Session: s.Session,
		Sqls:    execSqls,
	})
	s.afterLog(api, err)
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

func (s *PluginImplV2) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	api := "Query"
	s.preLog(api)
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
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}
	result := &driverV2.QueryResult{
		Column: params.Params{},
		Rows:   []*driverV2.QueryResultRow{},
	}
	for _, p := range res.GetColumn() {
		i18nDesc, err := i18nPkg.ConvertStrMap2I18nStr(p.I18NDesc)
		if err != nil {
			return nil, fmt.Errorf("PluginImplV2 Query fail to convert i18nDesc to I18nStrMap, error: %v", err)
		}
		result.Column = append(result.Column, &params.Param{
			Key:      p.GetKey(),
			Value:    p.GetValue(),
			Desc:     p.GetDesc(),
			I18nDesc: i18nDesc,
			Type:     params.ParamType(p.GetType()),
		})
	}
	for _, row := range res.GetRows() {
		r := &driverV2.QueryResultRow{
			Values: []*driverV2.QueryResultValue{},
		}
		for _, value := range row.GetValues() {
			r.Values = append(r.Values, &driverV2.QueryResultValue{
				Value: value.GetValue(),
			})
		}
		result.Rows = append(result.Rows, r)
	}
	return result, nil
}

func (s *PluginImplV2) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	api := "Explain"
	s.preLog(api)
	req := &protoV2.ExplainRequest{
		Session: s.Session,
		Sql: &protoV2.ExplainSQL{
			Query: conf.Sql,
		},
	}
	res, err := s.client.Explain(ctx, req)
	s.afterLog(api, err)
	if err != nil && status.Code(err) == driverV2.GrpcErrSQLIsNotSupported {
		return nil, driverV2.ErrSQLIsNotSupported
	} else if err != nil {
		return nil, err
	}

	td, err := driverV2.ConvertProtoTabularDataToDriver(res.ClassicResult.Data, s.meta.IsOptionalModuleEnabled(driverV2.OptionalModuleI18n))
	if err != nil {
		return nil, fmt.Errorf("ClassicResult: %w", err)
	}
	return &driverV2.ExplainResult{
		ClassicResult: driverV2.ExplainClassicResult{
			TabularData: td,
		},
	}, nil
}

func (p *PluginImplV2) ExplainJSONFormat(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainJSONResult, error) {
	return nil, errors.New("ExplainJSONFormat not support yet")
}

// metadata

func (s *PluginImplV2) Schemas(ctx context.Context) ([]string, error) {
	api := "GetDatabases"
	s.preLog(api)
	resp, err := s.client.GetDatabases(ctx, &protoV2.GetDatabasesRequest{
		Session: s.Session,
	})
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}
	databases := make([]string, 0, len(resp.Databases))
	for _, d := range resp.Databases {
		databases = append(databases, d.Name)
	}
	return databases, nil
}

func (s *PluginImplV2) getTableMeta(ctx context.Context, table *driverV2.Table) (*driverV2.TableMeta, error) {
	api := "GetTableMeta"
	s.preLog(api)
	result, err := s.client.GetTableMeta(ctx, &protoV2.GetTableMetaRequest{
		Session: s.Session,
		Table: &protoV2.Table{
			Name:   table.Name,
			Schema: table.Schema,
		},
	})
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}
	return driverV2.ConvertProtoTableMetaToDriver(result.TableMeta, s.meta.IsOptionalModuleEnabled(driverV2.OptionalModuleI18n))
}

func (s *PluginImplV2) extractTableFromSQL(ctx context.Context, sql string) ([]*driverV2.Table, error) {
	api := "ExtractTableFromSQL"
	s.preLog(api)
	result, err := s.client.ExtractTableFromSQL(ctx, &protoV2.ExtractTableFromSQLRequest{
		Session: s.Session,
		Sql:     &protoV2.ExtractedSQL{Query: sql},
	})
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}
	tables := make([]*driverV2.Table, 0, len(result.Tables))
	for _, table := range result.Tables {
		tables = append(tables, &driverV2.Table{
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

func (s *PluginImplV2) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	api := "EstimateSQLAffectRows"
	s.preLog(api)
	ar, err := s.client.EstimateSQLAffectRows(ctx, &protoV2.EstimateSQLAffectRowsRequest{
		Session: s.Session,
		Sql: &protoV2.AffectRowsSQL{
			Query: sql,
		},
	})
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}
	return &driverV2.EstimatedAffectRows{
		Count:      ar.Count,
		ErrMessage: ar.ErrMessage,
	}, nil
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

func (s *PluginImplV2) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabaseSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {
	api := "GetDatabaseObjectDDL"
	s.preLog(api)
	dbInfoReq := make([]*protoV2.DatabaseSchemaInfo, len(objInfos))
	for i, dbSchema := range objInfos {
		dbObjs := make([]*protoV2.DatabaseObject, len(dbSchema.DatabaseObjects))
		for j, dbObj := range dbSchema.DatabaseObjects {
			dbObjs[j] = &protoV2.DatabaseObject{
				ObjectName: dbObj.ObjectName,
				ObjectType: dbObj.ObjectType,
			}
		}
		dbInfoReq[i] = &protoV2.DatabaseSchemaInfo{
			SchemaName:     dbSchema.SchemaName,
			DatabaseObject: dbObjs,
		}
	}
	resp, err := s.client.GetDatabaseObjectDDL(ctx, &protoV2.DatabaseObjectInfoRequest{
		Session:            s.Session,
		DatabaseSchemaInfo: dbInfoReq,
	})
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}
	ret := make([]*driverV2.DatabaseSchemaObjectResult, len(resp.DatabaseSchemaObject))
	for i, info := range resp.DatabaseSchemaObject {
		ObjDDL := make([]*driverV2.DatabaseObjectDDL, len(info.DatabaseObjectDDL))
		for j, obj := range info.DatabaseObjectDDL {
			ObjDDL[j] = &driverV2.DatabaseObjectDDL{
				DatabaseObject: &driverV2.DatabaseObject{
					ObjectName: obj.DatabaseObject.ObjectName,
					ObjectType: obj.DatabaseObject.ObjectType,
				},
				ObjectDDL: obj.ObjectDDL,
			}
		}
		ret[i] = &driverV2.DatabaseSchemaObjectResult{
			SchemaName:         info.SchemaName,
			SchemaDDL:          info.SchemaDDL,
			DatabaseObjectDDLs: ObjDDL,
		}
	}
	return ret, nil
}

func (s *PluginImplV2) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	api := "GetDatabaseDiffModifySQL"
	s.preLog(api)
	resp, err := s.client.GetDatabaseDiffModifySQL(ctx, &protoV2.DatabaseDiffModifyRequest{
		Session: s.Session,
		CalibratedDSN: &protoV2.DSN{
			Host:             calibratedDSN.Host,
			Port:             calibratedDSN.Port,
			User:             calibratedDSN.User,
			Password:         calibratedDSN.Password,
			AdditionalParams: driverV2.ConvertParamToProtoParam(calibratedDSN.AdditionalParams),
			Database:         calibratedDSN.DatabaseName,
		},
		ObjInfos: driverV2.ConvertDatabasSchemaInfoToProto(objInfos),
	})
	s.afterLog(api, err)
	if err != nil {
		return nil, err
	}

	dbDiffSQLs := make([]*driverV2.DatabaseDiffModifySQLResult, len(resp.SchemaDiffModify))
	for i, schemaDiff := range resp.SchemaDiffModify {
		dbDiffSQLs[i] = &driverV2.DatabaseDiffModifySQLResult{
			SchemaName: schemaDiff.SchemaName,
			ModifySQLs: schemaDiff.ModifySQLs,
		}
	}
	return dbDiffSQLs, nil
}
