package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type DSN struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (d *DSN) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.User, d.Password, d.Host, d.Port, d.Database)
}

type DB struct {
	db *sql.DB
}

func NewDB(dsn *DSN) (*DB, error) {
	// 创建一个数据库连接池
	db, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return nil, err
	}

	// 设置连接池的最大连接数和空闲连接数
	db.SetMaxOpenConns(100) // 设置最大连接数
	db.SetMaxIdleConns(10)  // 设置空闲连接数

	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

func (o *DB) Close() error {
	return o.db.Close()
}

func (o *DB) QueryTopSQLs(ctx context.Context, topN int, orderBy string) ([]*DynPerformancePgColumns, error) {
	query := fmt.Sprintf(DynPerformanceViewPgSQLTpl, orderBy, topN)
	rows, err := o.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query %s", query)
	}
	defer rows.Close()

	var ret []*DynPerformancePgColumns
	for rows.Next() {
		res := DynPerformancePgColumns{}
		err = rows.Scan(&res.SQLFullText, &res.Executions, &res.ElapsedTime, &res.CPUTime,
			&res.DiskReads, &res.BufferGets, &res.UserIOWaitTime)
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
