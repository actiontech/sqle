package driver

import (
	"context"
	"database/sql"
	_driver "database/sql/driver"
	"fmt"
	"time"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	hclog "github.com/hashicorp/go-hclog"
	
	"github.com/percona/go-mysql/query"
	"github.com/pkg/errors"
	"vitess.io/vitess/go/vt/sqlparser"
)

type DriverImpl struct {
	Log    hclog.Logger
	Config *driverV2.Config
	Ah     *AuditHandler

	Dt   Dialector
	DB   *sql.DB
	Conn *sql.Conn
}

func NewDriverImpl(l hclog.Logger, dt Dialector, ah *AuditHandler, cfg *driverV2.Config) (driverV2.Driver, error) {
	di := &DriverImpl{
		Log:    l,
		Config: cfg,
		Ah:     ah,
		Dt:     dt,
	}
	if cfg.DSN == nil {
		return di, nil
	}
	db, conn, err := dt.Open(cfg.DSN)
	if err != nil {
		return nil, err
	}
	di.DB = db     // will be closed by DriverImpl.Close
	di.Conn = conn // will be closed by DriverImpl.Close
	return di, nil
}

func (p *DriverImpl) GetConn() (*sql.Conn, error) {
	if p.Conn == nil {
		return nil, fmt.Errorf("database conn not initialized")
	}
	return p.Conn, nil
}

// check pluginImpl is implement driver.plugin
var _ driverV2.Driver = &DriverImpl{}

func (p *DriverImpl) Close(ctx context.Context) {
	if p.Conn != nil {
		p.Conn.Close()
	}
	if p.DB != nil {
		p.DB.Close()
	}
}

func (p *DriverImpl) Ping(ctx context.Context) error {
	conn, err := p.GetConn()
	if err != nil {
		return err
	}
	if err := conn.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping in driver adaptor")
	}
	return nil
}

func (p *DriverImpl) Exec(ctx context.Context, sql string) (_driver.Result, error) {
	conn, err := p.GetConn()
	if err != nil {
		return nil, err
	}
	res, err := conn.ExecContext(ctx, sql)
	if err != nil {
		return nil, errors.Wrap(err, "exec sql in driver adaptor")
	}
	return res, nil
}

func (p *DriverImpl) Tx(ctx context.Context, sqls ...string) ([]_driver.Result, error) {
	var (
		err error
		tx  *sql.Tx
	)
	conn, err := p.GetConn()
	if err != nil {
		return nil, err
	}
	tx, err = conn.BeginTx(ctx, nil)
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

func (p *DriverImpl) Query(ctx context.Context, query string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	conn, err := p.GetConn()
	if err != nil {
		return nil, err
	}

	var cancel func()
	if conf != nil && conf.TimeOutSecond > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(conf.TimeOutSecond)*time.Second)
		defer cancel()
	}
	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := &driverV2.QueryResult{
		Column: params.Params{},
		Rows:   []*driverV2.QueryResultRow{},
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
		value := &driverV2.QueryResultRow{
			Values: []*driverV2.QueryResultValue{},
		}
		for i := 0; i < len(columns); i++ {
			value.Values = append(value.Values, &driverV2.QueryResultValue{Value: data[i].String})
		}
		result.Rows = append(result.Rows, value)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *DriverImpl) Parse(ctx context.Context, sql string) ([]driverV2.Node, error) {
	sqls, err := sqlparser.SplitStatementToPieces(sql)
	if err != nil {
		return nil, errors.Wrap(err, "split sql")
	}
	if err != nil {
		return nil, errors.Wrapf(err, "split sql %s error", sql)
	}

	nodes := make([]driverV2.Node, 0, len(sqls))
	for _, sql := range sqls {
		n := driverV2.Node{
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
		return driverV2.SQLTypeDML
	}
	return driverV2.SQLTypeDDL
}

func (p *DriverImpl) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	results := make([]*driverV2.AuditResults, 0, len(sqls))
	for i, sql := range sqls {
		ruleResults := driverV2.NewAuditResults()
		for j, rule := range p.Config.Rules {
			result, err := p.Ah.Audit(ctx, rule, sql, sqls[i+1:])
			if err != nil {
				return nil, err
			}
			ruleResults.Results[j] = result
		}
		results = append(results, ruleResults)
	}
	return results, nil
}

func (p *DriverImpl) GenRollbackSQL(ctx context.Context, sql string) (string, string, error) {
	return "", "", nil
}

func (p *DriverImpl) GetDatabases(ctx context.Context) ([]string, error) {
	conn, err := p.GetConn()
	if err != nil {
		return nil, err
	}
	rows, err := conn.QueryContext(ctx, p.Dt.ShowDatabaseSQL())
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

func (p *DriverImpl) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	return &driverV2.ExplainResult{}, nil
}

func (p *DriverImpl) GetTableMeta(ctx context.Context, table *driverV2.Table) (*driverV2.TableMeta, error) {
	return &driverV2.TableMeta{}, nil

}
func (p *DriverImpl) ExtractTableFromSQL(ctx context.Context, sql string) ([]*driverV2.Table, error) {
	return []*driverV2.Table{}, nil
}

func (p *DriverImpl) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	return &driverV2.EstimatedAffectRows{}, nil
}

func (p *DriverImpl) KillProcess(ctx context.Context) (*driverV2.KillProcessInfo, error) {
	return &driverV2.KillProcessInfo{}, nil
}
