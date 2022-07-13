package driver

import (
	"context"
	"database/sql"
	_driver "database/sql/driver"
	"sync"
	"time"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"

	"github.com/percona/go-mysql/query"
	"github.com/pkg/errors"
	"vitess.io/vitess/go/vt/sqlparser"
)

var pluginImpls = make(map[string]*pluginImpl)
var pluginImplsMu = &sync.Mutex{}

type pluginImpl struct {
	auditAdaptor    *AuditAdaptor
	queryAdaptor    *QueryAdaptor
	analysisAdaptor *AnalysisAdaptor
	db              *sql.DB
	conn            *sql.Conn
}

func (p *pluginImpl) Close(ctx context.Context) {
	for name, pluginImpl := range pluginImpls {
		if pluginImpl.conn != nil {
			if err := pluginImpl.conn.Close(); err != nil {
				pluginImpl.auditAdaptor.l.Error("failed to close connection in driver adaptor", "err", err, "plugin_name", name)
			}
		}

		if pluginImpl.db != nil {
			if err := pluginImpl.db.Close(); err != nil {
				pluginImpl.auditAdaptor.l.Error("failed to close database in driver adaptor", "err", err, "plugin_name", name)
			}
		}

		pluginImplsMu.Lock()
		delete(pluginImpls, name)
		pluginImplsMu.Unlock()
	}
}

func (p *pluginImpl) Ping(ctx context.Context) error {
	if err := p.conn.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping in driver adaptor")
	}
	return nil
}

func (p *pluginImpl) Exec(ctx context.Context, sql string) (_driver.Result, error) {
	res, err := p.conn.ExecContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "exec sql in driver adaptor")
	}
	return res, nil
}

func (p *pluginImpl) Tx(ctx context.Context, sqls ...string) ([]_driver.Result, error) {
	var (
		err error
		tx  *sql.Tx
	)

	tx, err = p.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin tx in driver adaptor")
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				err = errors.Wrap(err, "rollback tx in driver adaptor")
				return
			}
		} else {
			if err = tx.Commit(); err != nil {
				err = errors.Wrap(err, "commit tx in driver adaptor")
				return
			}
		}
	}()

	results := make([]_driver.Result, 0, len(sqls))
	for _, sql := range sqls {
		result, e := tx.ExecContext(ctx, sql)
		if e != nil {
			err = errors.Wrap(e, "exec sql in driver adaptor")
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (p *pluginImpl) Schemas(ctx context.Context) ([]string, error) {
	rows, err := p.conn.QueryContext(ctx, p.auditAdaptor.dt.ShowDatabaseSQL())
	if err != nil {
		return nil, errors.Wrap(err, "query database in driver adaptor")
	}
	defer rows.Close()

	var schemas []string
	for rows.Next() {
		var schema string
		if err := rows.Scan(&schema); err != nil {
			return nil, errors.Wrap(err, "scan database in driver adaptor")
		}
		schemas = append(schemas, schema)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "scan database in driver adaptor")
	}

	return schemas, nil
}

func (p *pluginImpl) Parse(ctx context.Context, sql string) ([]driver.Node, error) {
	sqls, err := sqlparser.SplitStatementToPieces(sql)
	if err != nil {
		return nil, errors.Wrap(err, "split sql")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "split sql %s error", sql)
	}

	nodes := make([]driver.Node, 0, len(sqls))
	for _, sql := range sqls {
		n := driver.Node{
			Text:        sql,
			Type:        classifySQL(sql),
			Fingerprint: query.Fingerprint(sql),
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func classifySQL(sql string) (sqlType string) {
	if utils.HasPrefix(sql, "update", false) ||
		utils.HasPrefix(sql, "insert", false) ||
		utils.HasPrefix(sql, "delete", false) {
		return driver.SQLTypeDML
	}

	return driver.SQLTypeDDL
}

func (p *pluginImpl) Audit(ctx context.Context, sql string) (*driver.AuditResult, error) {
	var err error
	var ast interface{}
	if p.auditAdaptor.ao.sqlParser != nil {
		ast, err = p.auditAdaptor.ao.sqlParser(sql)
		if err != nil {
			return nil, errors.Wrap(err, "parse sql")
		}
	}

	result := driver.NewInspectResults()
	for _, rule := range p.auditAdaptor.cfg.Rules {
		handler, ok := p.auditAdaptor.ruleToRawHandler[rule.Name]
		if ok {
			msg, err := handler(ctx, rule, sql)
			if err != nil {
				return nil, errors.Wrapf(err, "audit SQL %s in driver adaptor", sql)
			}
			result.Add(rule.Level, msg)
		} else {
			handler, ok := p.auditAdaptor.ruleToASTHandler[rule.Name]
			if ok {
				msg, err := handler(ctx, rule, ast)
				if err != nil {
					return nil, errors.Wrapf(err, "audit SQL %s in driver adaptor", sql)
				}
				result.Add(rule.Level, msg)
			}
		}
	}

	return result, nil
}

func (p *pluginImpl) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	return "", "", nil
}

func (p *pluginImpl) QueryPrepare(ctx context.Context, sql string, conf *driver.QueryPrepareConf) (*driver.QueryPrepareResult, error) {
	if p.queryAdaptor.queryPrepare != nil {
		return p.queryAdaptor.queryPrepare(ctx, sql, conf, DbConf{
			Db:   p.db,
			Conn: p.conn,
		})
	}
	return &driver.QueryPrepareResult{
		NewSQL:    sql,
		ErrorType: driver.ErrorTypeNotError,
		Error:     "",
	}, nil
}

func (p *pluginImpl) Query(ctx context.Context, query string, conf *driver.QueryConf) (*driver.QueryResult, error) {
	if p.queryAdaptor.query != nil {
		return p.queryAdaptor.query(ctx, query, conf, DbConf{
			Db:   p.db,
			Conn: p.conn,
		})
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(conf.TimeOutSecond)*time.Second)
	defer cancel()
	rows, err := p.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := &driver.QueryResult{
		Column: params.Params{},
		Rows:   []*driver.QueryResultRow{},
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for _, column := range columns {
		result.Column = append(result.Column, &params.Param{
			Key:   column,
			Value: column,
			Desc:  column,
		})
	}

	for rows.Next() {
		buf := make([]interface{}, len(columns))
		data := make([]sql.NullString, len(columns))
		for i := range buf {
			buf[i] = &data[i]
		}
		if err := rows.Scan(buf...); err != nil {
			return nil, err
		}
		value := &driver.QueryResultRow{
			Values: []*driver.QueryResultValue{},
		}
		for i := 0; i < len(columns); i++ {
			value.Values = append(value.Values, &driver.QueryResultValue{Value: data[i].String})
		}
		result.Rows = append(result.Rows, value)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *pluginImpl) ListTablesInSchema(ctx context.Context, conf *driver.ListTablesInSchemaConf) (*driver.ListTablesInSchemaResult, error) {
	if p.analysisAdaptor.listTablesInSchemaFunc != nil {
		return p.analysisAdaptor.listTablesInSchemaFunc(ctx, conf, DbConf{
			Db:   p.db,
			Conn: p.conn,
		})
	}

	return &driver.ListTablesInSchemaResult{}, nil
}
func (p *pluginImpl) GetTableMetaByTableName(ctx context.Context, conf *driver.GetTableMetaByTableNameConf) (*driver.GetTableMetaByTableNameResult, error) {
	if p.analysisAdaptor.getTableMetaByTableNameFunc != nil {
		return p.analysisAdaptor.getTableMetaByTableNameFunc(ctx, conf, DbConf{
			Db:   p.db,
			Conn: p.conn,
		})
	}

	return &driver.GetTableMetaByTableNameResult{}, nil
}
func (p *pluginImpl) GetTableMetaBySQL(ctx context.Context, conf *driver.GetTableMetaBySQLConf) (*driver.GetTableMetaBySQLResult, error) {
	if p.analysisAdaptor.getTableMetaBySQLFunc != nil {
		return p.analysisAdaptor.getTableMetaBySQLFunc(ctx, conf, DbConf{
			Db:   p.db,
			Conn: p.conn,
		})
	}

	return &driver.GetTableMetaBySQLResult{}, nil
}
func (p *pluginImpl) Explain(ctx context.Context, conf *driver.ExplainConf) (*driver.ExplainResult, error) {
	if p.analysisAdaptor.explainFunc != nil {
		return p.analysisAdaptor.explainFunc(ctx, conf, DbConf{
			Db:   p.db,
			Conn: p.conn,
		})
	}

	return &driver.ExplainResult{}, nil
}
