package driver

import (
	"context"
	sqlDriver "database/sql/driver"

	v1 "github.com/actiontech/sqle/sqle/driver/v1"
	v2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	goPlugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
)

type PluginBootV1 struct {
	cfg    func(path string) *goPlugin.ClientConfig
	path   string
	client *goPlugin.Client // this client will be killed after Register.
	metas  *v2.DriverMetas
}

func convertRuleFromV1ToV2(rule *v1.Rule) *v2.Rule {
	var ps = make(params.Params, 0, len(rule.Params))
	for _, p := range rule.Params {
		ps = append(ps, &params.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  p.Type,
		})
	}
	return &v2.Rule{
		Name:       rule.Name,
		Category:   rule.Category,
		Desc:       rule.Desc,
		Annotation: rule.Annotation,
		Level:      v2.RuleLevel(rule.Level),
		Params:     ps,
	}
}

func (d *PluginBootV1) Register() (*v2.DriverMetas, error) {
	defer d.client.Kill()
	name, rules, params, err := v1.RegisterDrivers(d.client, d.cfg, d.path)
	if err != nil {
		return nil, err
	}

	rulesV2 := make([]*v2.Rule, 0, len(rules))
	for _, rule := range rules {
		rulesV2 = append(rulesV2, convertRuleFromV1ToV2(rule))
	}
	meta := &v2.DriverMetas{
		PluginName:               name,
		DatabaseDefaultPort:      0,
		Rules:                    rulesV2,
		DatabaseAdditionalParams: params,
	}
	d.metas = meta
	return meta, nil
}

func (d *PluginBootV1) Open(l *logrus.Entry, cfgV2 *v2.Config) (Plugin, error) {
	l = l.WithFields(logrus.Fields{
		"plugin":         d.metas.PluginName,
		"plugin_version": v1.ProtocolVersion,
	})
	cfg := &v1.Config{
		DSN: &v1.DSN{
			Host:             cfgV2.DSN.Host,
			Port:             cfgV2.DSN.Port,
			User:             cfgV2.DSN.User,
			Password:         cfgV2.DSN.Password,
			DatabaseName:     cfgV2.DSN.DatabaseName,
			AdditionalParams: cfgV2.DSN.AdditionalParams,
		},
	}
	for _, rule := range cfgV2.Rules {
		cfg.Rules = append(cfg.Rules, &v1.Rule{
			Name:       rule.Name,
			Desc:       rule.Desc,
			Annotation: rule.Annotation,
			Category:   rule.Category,
			Level:      v1.RuleLevel(rule.Level),
			Params:     rule.Params,
		})
	}
	dm, err := v1.NewDriverManger(l, d.metas.PluginName, cfg)
	if err != nil {
		return nil, err
	}
	p := &PluginImplV1{
		dm,
	}
	return p, nil
}

func (d *PluginBootV1) Stop() error {
	return nil
}

type PluginImplV1 struct {
	v1.DriverManager
}

func (p *PluginImplV1) Close(ctx context.Context) {
	p.DriverManager.Close(ctx)
}

func (p *PluginImplV1) Parse(ctx context.Context, sqlText string) ([]v2.Node, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return nil, err
	}
	nodes, err := client.Parse(ctx, sqlText)
	if err != nil {
		return nil, err
	}
	nodesV2 := []v2.Node{}
	for _, node := range nodes {
		nodesV2 = append(nodesV2, v2.Node{
			Text:        node.Text,
			Type:        node.Type,
			Fingerprint: node.Fingerprint,
		})
	}
	return nodesV2, nil
}

func (p *PluginImplV1) Audit(ctx context.Context, sqls []string) ([]*v2.AuditResults, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return nil, err
	}
	resultsV2 := []*v2.AuditResults{}
	for _, sql := range sqls {
		resultV1, err := client.Audit(ctx, sql)
		if err != nil {
			return nil, err
		}
		resultV2 := &v2.AuditResults{}
		for _, result := range resultV1.Results {
			resultV2.Results = append(resultV2.Results, &v2.AuditResult{
				Level:   v2.RuleLevel(result.Level),
				Message: result.Message,
			})
		}
		resultsV2 = append(resultsV2, resultV2)
	}
	return resultsV2, nil
}

func (p *PluginImplV1) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return "", "", err
	}
	return client.GenRollbackSQL(ctx, sql)
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

func (p *PluginImplV1) Tx(ctx context.Context, queries ...string) ([]sqlDriver.Result, error) {
	client, err := p.DriverManager.GetAuditDriver()
	if err != nil {
		return nil, err
	}
	return client.Tx(ctx, queries...)
}

func (p *PluginImplV1) Query(ctx context.Context, sql string, conf *v2.QueryConf) (*v2.QueryResult, error) {
	client, err := p.DriverManager.GetSQLQueryDriver()
	if err != nil {
		return nil, err
	}
	resultV1, err := client.Query(ctx, sql, &v1.QueryConf{TimeOutSecond: conf.TimeOutSecond})
	if err != nil {
		return nil, err
	}
	rowsV2 := []*v2.QueryResultRow{}
	for _, row := range resultV1.Rows {
		rowV2 := &v2.QueryResultRow{}
		for _, v := range row.Values {
			rowV2.Values = append(rowV2.Values, &v2.QueryResultValue{
				Value: v.Value,
			})
		}
		rowsV2 = append(rowsV2, rowV2)
	}

	return &v2.QueryResult{
		Column: resultV1.Column,
		Rows:   rowsV2,
	}, nil
}

func (p *PluginImplV1) Explain(ctx context.Context, conf *v2.ExplainConf) (*v2.ExplainResult, error) {
	client, err := p.DriverManager.GetAnalysisDriver()
	if err != nil {
		return nil, err
	}
	resultV1, err := client.Explain(ctx, &v1.ExplainConf{
		Sql: conf.Sql,
	})
	if err != nil {
		return nil, err
	}

	columnsV2 := []v2.TabularDataHead{}
	for _, column := range resultV1.ClassicResult.Columns {
		columnsV2 = append(columnsV2, v2.TabularDataHead{
			Name: column.Name,
			Desc: column.Desc,
		})
	}

	resultV2 := &v2.ExplainResult{}
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
	resultV1, err := client.GetTableMetaBySQL(ctx, &v1.GetTableMetaBySQLConf{
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

		columnV2 := []v2.TabularDataHead{}
		for _, column := range tm.ColumnsInfo.Columns {
			columnV2 = append(columnV2, v2.TabularDataHead{
				Name: column.Name,
				Desc: column.Desc,
			})
		}
		tmV2.ColumnsInfo.Columns = columnV2
		tmV2.ColumnsInfo.Rows = tm.ColumnsInfo.Rows

		indexesColV2 := []v2.TabularDataHead{}
		for _, column := range tm.IndexesInfo.Columns {
			indexesColV2 = append(indexesColV2, v2.TabularDataHead{
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
