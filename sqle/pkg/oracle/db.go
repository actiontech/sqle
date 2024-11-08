package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	_ "github.com/sijms/go-ora/v2"
)

type DSN struct {
	Host        string
	Port        string
	User        string
	Password    string
	ServiceName string
}

func (d *DSN) String() string {
	return fmt.Sprintf("%s:%s/%s", d.Host, d.Port, d.ServiceName)
}

type DB struct {
	db *sql.DB
}

func NewDB(dsn *DSN) (*DB, error) {
	if dsn.ServiceName == "" {
		dsn.ServiceName = "xe"
	}

	sqlDB, err := sql.Open("oracle", fmt.Sprintf("oracle://%s:%s@%s:%s/%s", dsn.User, url.QueryEscape(dsn.Password), dsn.Host, dsn.Port, dsn.ServiceName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to %s", dsn.String())
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to ping %s", dsn.String())
	}

	return &DB{db: sqlDB}, nil
}

func (o *DB) Close() error {
	return o.db.Close()
}

func (o *DB) QueryTopSQLs(ctx context.Context, topN int, orderBy string) ([]*DynPerformanceSQLArea, error) {
	query := fmt.Sprintf(DynPerformanceViewSQLAreaTpl, orderBy, topN)
	rows, err := o.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query %s", query)
	}
	defer rows.Close()

	var ret []*DynPerformanceSQLArea
	for rows.Next() {
		res := DynPerformanceSQLArea{}
		err = rows.Scan(&res.SQLFullText, &res.Executions, &res.ElapsedTime, &res.UserIOWaitTime, &res.CPUTime, &res.DiskReads, &res.BufferGets)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to scan %s", query)
		}
		ret = append(ret, &res)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to iterate %s", query)
	}

	return ret, nil
}
