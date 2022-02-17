package oracle

import (
	"context"
	"database/sql"
	"fmt"

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

	sqlDB, err := sql.Open("oracle", fmt.Sprintf("oracle://%s:%s@%s:%s/%s", dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.ServiceName))
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

func (o *DB) QueryTopSQLs(ctx context.Context, topN int) ([]string, error) {
	deDupSQLs := make(map[string]struct{})
	queryFunc := func(query string) error {
		rows, err := o.db.QueryContext(ctx, query)
		if err != nil {
			return errors.Wrapf(err, "failed to query %s", query)
		}
		defer rows.Close()

		for rows.Next() {
			res := DynPerformanceSQLArea{}
			err = rows.Scan(&res.SQLFullText, &res.Avg)
			if err != nil {
				return errors.Wrapf(err, "failed to scan %s", query)
			}
			deDupSQLs[res.SQLFullText] = struct{}{}
		}

		if err := rows.Err(); err != nil {
			return errors.Wrapf(err, "failed to iterate %s", query)
		}

		return nil
	}

	if err := queryFunc(fmt.Sprintf(DynPerformanceViewSQLAreaTpl, DynPerformanceViewSQLAreaColumnElapsedTime, topN)); err != nil {
		return nil, err
	}
	if err := queryFunc(fmt.Sprintf(DynPerformanceViewSQLAreaTpl, DynPerformanceViewSQLAreaColumnCPUTime, topN)); err != nil {
		return nil, err
	}
	if err := queryFunc(fmt.Sprintf(DynPerformanceViewSQLAreaTpl, DynPerformanceViewSQLAreaColumnBufferGets, topN)); err != nil {
		return nil, err
	}
	if err := queryFunc(fmt.Sprintf(DynPerformanceViewSQLAreaTpl, DynPerformanceViewSQLAreaColumnDiskReads, topN)); err != nil {
		return nil, err
	}
	if err := queryFunc(fmt.Sprintf(DynPerformanceViewSQLAreaTpl, DynPerformanceViewSQLAreaColumnUserIOWaitTime, topN)); err != nil {
		return nil, err
	}

	sqls := make([]string, 0, len(deDupSQLs))
	for sql := range deDupSQLs {
		sqls = append(sqls, sql)
	}
	return sqls, nil
}
