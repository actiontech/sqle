package driver

import (
	"context"
	sqlDriver "database/sql/driver"
	"fmt"

	driverV1 "github.com/actiontech/sqle/sqle/driver/v1"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/pkg/params"
	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
)

type PluginProcessorV1 struct {
	cfg     func(cmdBase string, cmdArgs []string) *goPlugin.ClientConfig
	cmdBase string
	cmdArgs []string
	client  *goPlugin.Client // this client will be killed after Register.
	metas   *driverV2.DriverMetas
}

func convertRuleFromV1ToV2(rule *driverV1.Rule) *driverV2.Rule {
	var ps = make(params.Params, 0, len(rule.Params))
	for _, p := range rule.Params {
		ps = append(ps, &params.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  p.Type,
		})
	}
	return &driverV2.Rule{
		Name:   rule.Name,
		Level:  driverV2.RuleLevel(rule.Level),
		Params: ps,
		I18nRuleInfo: map[string]*driverV2.RuleInfo{
			locale.DefaultLang.String(): {
				Desc:       rule.Desc,
				Annotation: rule.Annotation,
				Category:   rule.Category,
				Params:     ps,
				Knowledge:  driverV2.RuleKnowledge{},
			},
		},
	}
}

func (d *PluginProcessorV1) GetDriverMetas() (*driverV2.DriverMetas, error) {
	defer d.client.Kill()
	name, rules, params, enableQuery, enableSQLAnalysis, err := driverV1.RegisterDrivers(d.client, d.cfg, d.cmdBase, d.cmdArgs)
	if err != nil {
		return nil, err
	}

	m := []driverV2.OptionalModule{
		driverV2.OptionalModuleGenRollbackSQL, //在V1版本里面这个接口无法判断是否实现，需要插件返回空数据代表不生成。
	}
	if enableQuery {
		m = append(m, driverV2.OptionalModuleQuery)
	}
	if enableSQLAnalysis {
		m = append(m, []driverV2.OptionalModule{
			driverV2.OptionalModuleExplain,
			driverV2.OptionalModuleGetTableMeta,
			driverV2.OptionalModuleExtractTableFromSQL,
		}...)
	}

	rulesV2 := make([]*driverV2.Rule, 0, len(rules))
	for _, rule := range rules {
		rulesV2 = append(rulesV2, convertRuleFromV1ToV2(rule))
	}
	meta := &driverV2.DriverMetas{
		PluginName:               name,
		DatabaseDefaultPort:      0,
		Rules:                    rulesV2,
		DatabaseAdditionalParams: params,
		EnabledOptionalModule:    m,
	}
	d.metas = meta
	return meta, nil
}

func (d *PluginProcessorV1) Open(l *logrus.Entry, cfgV2 *driverV2.Config) (Plugin, error) {
	l = l.WithFields(logrus.Fields{
		"plugin":         d.metas.PluginName,
		"plugin_version": driverV1.ProtocolVersion,
	})
	var cfg = &driverV1.Config{}
	if cfgV2.DSN != nil {
		cfg.DSN = &driverV1.DSN{
			Host:             cfgV2.DSN.Host,
			Port:             cfgV2.DSN.Port,
			User:             cfgV2.DSN.User,
			Password:         cfgV2.DSN.Password,
			DatabaseName:     cfgV2.DSN.DatabaseName,
			AdditionalParams: cfgV2.DSN.AdditionalParams,
		}
	}
	for _, rule := range cfgV2.Rules {
		cfg.Rules = append(cfg.Rules, &driverV1.Rule{
			Name:       rule.Name,
			Desc:       rule.I18nRuleInfo.GetRuleInfoByLangTag(locale.DefaultLang.String()).Desc,
			Annotation: rule.I18nRuleInfo.GetRuleInfoByLangTag(locale.DefaultLang.String()).Annotation,
			Category:   rule.I18nRuleInfo.GetRuleInfoByLangTag(locale.DefaultLang.String()).Category,
			Level:      driverV1.RuleLevel(rule.Level),
			Params:     rule.Params,
		})
	}
	dm, err := driverV1.NewDriverManger(l, d.metas.PluginName, cfg)
	if err != nil {
		return nil, err
	}
	p := &PluginImplV1{
		dm,
	}
	return p, nil
}

func (d *PluginProcessorV1) Stop() error {
	return nil
}

type PluginImplV1 struct {
	driverV1.DriverManager
}

func (p *PluginImplV1) Close(ctx context.Context) {
	p.DriverManager.Close(ctx)
}

func (p *PluginImplV1) KillProcess(ctx context.Context) error {
	return errors.NewNotImplementedError("KillProcess not support yet")
}

func (p *PluginImplV1) Parse(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return nil, err
	}
	nodes, err := client.Parse(ctx, sqlText)
	if err != nil {
		return nil, err
	}
	nodesV2 := []driverV2.Node{}
	for _, node := range nodes {
		nodesV2 = append(nodesV2, driverV2.Node{
			Text:        node.Text,
			Type:        node.Type,
			Fingerprint: node.Fingerprint,
		})
	}
	return nodesV2, nil
}

func (p *PluginImplV1) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return nil, err
	}
	resultsV2 := []*driverV2.AuditResults{}
	for _, sql := range sqls {
		resultV1, err := client.Audit(ctx, sql)
		if err != nil {
			return nil, err
		}
		resultV2 := &driverV2.AuditResults{}
		for _, result := range resultV1.Results {
			resultV2.Results = append(resultV2.Results, &driverV2.AuditResult{
				Level: driverV2.RuleLevel(result.Level),
				I18nAuditResultInfo: map[string]driverV2.AuditResultInfo{
					locale.DefaultLang.String(): {
						Message: result.Message,
					},
				},
			})
		}
		resultsV2 = append(resultsV2, resultV2)
	}
	return resultsV2, nil
}

func (p *PluginImplV1) GenRollbackSQL(ctx context.Context, sql string) (string, driverV2.I18nStr, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return "", nil, err
	}
	rollbackSql, reason, err := client.GenRollbackSQL(ctx, sql)
	return rollbackSql, driverV2.I18nStr{locale.DefaultLang.String(): reason}, err
}

func (p *PluginImplV1) Ping(ctx context.Context) error {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return err
	}
	return client.Ping(ctx)
}

func (p *PluginImplV1) Exec(ctx context.Context, query string) (sqlDriver.Result, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return nil, err
	}
	return client.Exec(ctx, query)
}

func (s *PluginImplV1) ExecBatch(ctx context.Context, sqls ...string) ([]sqlDriver.Result, error) {
	return nil, fmt.Errorf("unimplement this method")
}

func (p *PluginImplV1) Tx(ctx context.Context, queries ...string) ([]sqlDriver.Result, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return nil, err
	}
	return client.Tx(ctx, queries...)
}

func (p *PluginImplV1) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	client, err := p.DriverManager.GetSQLQueryDriver()
	if err != nil {
		return nil, err
	}
	resultV1, err := client.Query(ctx, sql, &driverV1.QueryConf{TimeOutSecond: conf.TimeOutSecond})
	if err != nil {
		return nil, err
	}
	rowsV2 := []*driverV2.QueryResultRow{}
	for _, row := range resultV1.Rows {
		rowV2 := &driverV2.QueryResultRow{}
		for _, v := range row.Values {
			rowV2.Values = append(rowV2.Values, &driverV2.QueryResultValue{
				Value: v.Value,
			})
		}
		rowsV2 = append(rowsV2, rowV2)
	}

	return &driverV2.QueryResult{
		Column: resultV1.Column,
		Rows:   rowsV2,
	}, nil
}

func (p *PluginImplV1) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	client, err := p.DriverManager.GetAnalysisDriver()
	if err != nil {
		return nil, err
	}
	resultV1, err := client.Explain(ctx, &driverV1.ExplainConf{
		Sql: conf.Sql,
	})
	if err != nil {
		return nil, err
	}

	columnsV2 := []driverV2.TabularDataHead{}
	for _, column := range resultV1.ClassicResult.Columns {
		columnsV2 = append(columnsV2, driverV2.TabularDataHead{
			Name: column.Name,
			Desc: column.Desc,
		})
	}

	resultV2 := &driverV2.ExplainResult{}
	resultV2.ClassicResult.Rows = resultV1.ClassicResult.Rows
	resultV2.ClassicResult.Columns = columnsV2
	return resultV2, nil

}

func (p *PluginImplV1) Schemas(ctx context.Context) ([]string, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return nil, err
	}
	return client.Schemas(ctx)
}

func (p *PluginImplV1) GetTableMetaBySQL(ctx context.Context, conf *GetTableMetaBySQLConf) (*GetTableMetaBySQLResult, error) {
	client, err := p.DriverManager.GetAnalysisDriver()
	if err != nil {
		return nil, err
	}
	resultV1, err := client.GetTableMetaBySQL(ctx, &driverV1.GetTableMetaBySQLConf{
		Sql: conf.Sql,
	})
	if err != nil {
		return nil, err
	}

	resultV2 := &GetTableMetaBySQLResult{}
	for _, tm := range resultV1.TableMetas {
		tmV2 := &TableMeta{}
		tmV2.Table.Name = tm.Name
		tmV2.Table.Schema = tm.Schema
		tmV2.CreateTableSQL = tm.CreateTableSQL
		tmV2.Message = tm.Message

		columnV2 := []driverV2.TabularDataHead{}
		for _, column := range tm.ColumnsInfo.Columns {
			columnV2 = append(columnV2, driverV2.TabularDataHead{
				Name: column.Name,
				Desc: column.Desc,
			})
		}
		tmV2.ColumnsInfo.Columns = columnV2
		tmV2.ColumnsInfo.Rows = tm.ColumnsInfo.Rows

		indexesColV2 := []driverV2.TabularDataHead{}
		for _, column := range tm.IndexesInfo.Columns {
			indexesColV2 = append(indexesColV2, driverV2.TabularDataHead{
				Name: column.Name,
				Desc: column.Desc,
			})

		}
		tmV2.IndexesInfo.Columns = indexesColV2
		tmV2.IndexesInfo.Rows = tm.IndexesInfo.Rows

		resultV2.TableMetas = append(resultV2.TableMetas, tmV2)
	}
	return resultV2, nil
}

func (p *PluginImplV1) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	return nil, NewErrPluginAPINotImplement(driverV2.OptionalModuleEstimateSQLAffectRows)
}
