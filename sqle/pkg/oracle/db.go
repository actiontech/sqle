package oracle

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"text/template"

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
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s:%s/%s", dsn.User, url.QueryEscape(dsn.Password), dsn.Host, dsn.Port, dsn.ServiceName)
	if dsn.User == "sys" {
		dataSourceName = fmt.Sprintf("%s%s", dataSourceName, "?dba privilege=sysdba")
	}

	sqlDB, err := sql.Open("oracle", dataSourceName)
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

func (o *DB) QueryTopSQLs(ctx context.Context, collectIntervalMinute string, topN int, notInUsers []string, orderBy string) ([]*DynPerformanceSQLArea, error) {
	// if notInUsers is empty, notInUsersStr will be empty
	// if notInUsers is not empty, notInUsersStr will be formatted as "AND u.username NOT IN ('user1', 'user2')"
	var notInUsersStr string
	if len(notInUsers) > 0 {
		var notInUsersFormatted []string
		var notInUserSqlTpl = `AND u.username NOT IN (%v)`
		for _, user := range notInUsers {
			notInUsersFormatted = append(notInUsersFormatted, fmt.Sprintf("'%s'", user))
		}
		notInUsersStr = strings.Join(notInUsersFormatted, ",")
		notInUsersStr = fmt.Sprintf(notInUserSqlTpl, notInUsersStr)
	}
	metrics := []string{DynPerformanceViewSQLAreaColumnElapsedTime, DynPerformanceViewSQLAreaColumnCPUTime, DynPerformanceViewSQLAreaColumnDiskReads, DynPerformanceViewSQLAreaColumnBufferGets}
	if orderBy != "" {
		metrics = []string{orderBy}
	}
	if topN == 0 {
		topN = 10
	}
	var ret []*DynPerformanceSQLArea
	for _, metric := range metrics {
		query := fmt.Sprintf(DynPerformanceViewSQLAreaTpl, collectIntervalMinute, notInUsersStr, metric, topN)
		err := func() error {
			rows, err := o.db.QueryContext(ctx, query)
			if err != nil {
				return errors.Wrapf(err, "failed to query %s", query)
			}
			defer rows.Close()

			for rows.Next() {
				res := DynPerformanceSQLArea{}
				if err := rows.Scan(
					&res.SQLFullText,
					&res.Executions,
					&res.ElapsedTime,
					&res.UserIOWaitTime,
					&res.CPUTime,
					&res.DiskReads,
					&res.BufferGets,
					&res.UserName,
				); err != nil {
					return errors.Wrapf(err, "failed to scan %s", query)
				}
				ret = append(ret, &res)
			}
			if err := rows.Err(); err != nil {
				return errors.Wrapf(err, "failed to iterate %s", query)
			}
			return nil
		}()
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (o *DB) QuerySlowSQLs(ctx context.Context, collectIntervalMinute string, topN int, slowSQLThresholdMicroseconds int64, notInUsers []string) ([]*DynPerformanceSQLArea, error) {
	var notInUsersStr string
	if len(notInUsers) > 0 {
		var notInUsersFormatted []string
		var notInUserSqlTpl = `AND u.username NOT IN (%v)`
		for _, user := range notInUsers {
			notInUsersFormatted = append(notInUsersFormatted, fmt.Sprintf("'%s'", user))
		}
		notInUsersStr = strings.Join(notInUsersFormatted, ",")
		notInUsersStr = fmt.Sprintf(notInUserSqlTpl, notInUsersStr)
	}
	if topN == 0 {
		topN = 100
	}

	query := fmt.Sprintf(DynPerformanceViewSQLAreaSlowLogTpl, collectIntervalMinute, slowSQLThresholdMicroseconds, notInUsersStr, topN)
	rows, err := o.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query %s", query)
	}
	defer rows.Close()

	var ret []*DynPerformanceSQLArea
	for rows.Next() {
		res := DynPerformanceSQLArea{}
		if err := rows.Scan(
			&res.SQLFullText,
			&res.Executions,
			&res.ElapsedTime,
			&res.UserIOWaitTime,
			&res.CPUTime,
			&res.DiskReads,
			&res.BufferGets,
			&res.UserName,
		); err != nil {
			return nil, errors.Wrapf(err, "failed to scan %s", query)
		}
		ret = append(ret, &res)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to iterate %s", query)
	}
	return ret, nil
}

// QueryProcessList queries active sessions from V$SESSION joined with V$SQL.
// sqlMinSecond filters sessions whose LAST_CALL_ET >= sqlMinSecond (0 means no filter).
func (o *DB) QueryProcessList(ctx context.Context, sqlMinSecond int) ([]*ProcessListSession, error) {
	tpl, err := template.New("processlist").Parse(DynPerformanceViewSessionTpl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse process list sql template")
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, struct{ MinSecond int }{MinSecond: sqlMinSecond}); err != nil {
		return nil, errors.Wrap(err, "failed to execute process list sql template")
	}

	query := buf.String()
	rows, err := o.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query process list")
	}
	defer rows.Close()

	var ret []*ProcessListSession
	for rows.Next() {
		var (
			sid         sql.NullInt64
			username    sql.NullString
			schemaName  sql.NullString
			sqlFullText sql.NullString
			lastCallET  sql.NullInt64
		)
		if err := rows.Scan(&sid, &username, &schemaName, &sqlFullText, &lastCallET); err != nil {
			return nil, errors.Wrapf(err, "failed to scan process list row")
		}
		ret = append(ret, &ProcessListSession{
			SID:         sid.Int64,
			Username:    username.String,
			SchemaName:  schemaName.String,
			SQLFullText: sqlFullText.String,
			LastCallET:  lastCallET.Int64,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to iterate process list rows")
	}
	return ret, nil
}
