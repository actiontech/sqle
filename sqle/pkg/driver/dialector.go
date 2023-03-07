package driver

import (
	"context"
	"database/sql"
	"fmt"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pkg/errors"

	// DRIVER LIST:
	// 	https://github.com/golang/go/wiki/SQLDrivers
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/sijms/go-ora/v2"
)

// Dialector is a interface for database dialect. It used for sql.Open()
type Dialector interface {
	// String return the dialect name with more formal name, it will be used to define plugin name.
	String() string

	// DatabaseAdditionalParam return the database additional param, ex: oracle required service name to connect db;
	// it will be used by Dialector.Open.
	DatabaseAdditionalParam() params.Params

	Open(dsn *driverV2.DSN) (*sql.DB, *sql.Conn, error)

	// ShowDatabaseSQL return the sql to show all databases.
	ShowDatabaseSQL() string
}

type BaseDialector struct {
}

func (d *BaseDialector) String() string {
	return ""
}

func (d *BaseDialector) ShowDatabaseSQL() string {
	return ""
}

func (d *BaseDialector) DatabaseAdditionalParam() params.Params {
	return params.Params{}
}

func (d *BaseDialector) GetConn(driverName, dataSourceName string) (*sql.DB, *sql.Conn, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, nil, err
	}
	conn, err := db.Conn(context.TODO())
	if err != nil {
		db.Close()
		return nil, nil, errors.Wrap(err, "get database connection failed when new driver")
	}
	if err := conn.PingContext(context.TODO()); err != nil {
		conn.Close()
		db.Close()
		return nil, nil, errors.Wrap(err, "ping database connection failed when new driver")
	}
	return db, conn, nil
}

type PostgresDialector struct {
	BaseDialector
}

var _ Dialector = &PostgresDialector{}

func (d *PostgresDialector) String() string {
	return "PostgreSQL"
}

func (d *PostgresDialector) Open(dsn *driverV2.DSN) (*sql.DB, *sql.Conn, error) {
	if dsn.DatabaseName == "" {
		dsn.DatabaseName = "postgres"
	}
	db, conn, err := d.BaseDialector.GetConn("pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.DatabaseName))
	if err != nil {
		return nil, nil, err
	}
	return db, conn, nil
}

func (d *PostgresDialector) ShowDatabaseSQL() string {
	return "select datname from pg_database"
}

const serverNameKey = "service_name"

type OracleDialector struct {
	BaseDialector
}

var _ Dialector = &OracleDialector{}

func (d *OracleDialector) String() string {
	return "Oracle"
}

func (d *OracleDialector) DatabaseAdditionalParam() params.Params {
	return params.Params{
		&params.Param{
			Key:   serverNameKey,
			Value: "XE",
			Desc:  "service name",
			Type:  params.ParamTypeString,
		},
	}
}

func (d *OracleDialector) Open(dsn *driverV2.DSN) (*sql.DB, *sql.Conn, error) {
	serviceName := dsn.AdditionalParams.GetParam(serverNameKey).String()
	if serviceName == "" {
		serviceName = "XE"
	}
	db, conn, err := d.BaseDialector.GetConn("oracle", fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		dsn.User, dsn.Password, dsn.Host, dsn.Port, serviceName))
	if err != nil {
		return nil, nil, err
	}
	if dsn.DatabaseName != "" {
		_, err = conn.ExecContext(context.TODO(), "ALTER SESSION SET CURRENT_SCHEMA = ?", dsn.DatabaseName)
		if err != nil {
			conn.Close()
			db.Close()
			return nil, nil, errors.Wrap(err, fmt.Sprintf("switch to schema %s failed", dsn.DatabaseName))
		}
	}
	return db, conn, nil
}

func (d *OracleDialector) ShowDatabaseSQL() string {
	return "SELECT username FROM all_users"
}

type MssqlDialector struct {
	BaseDialector
}

var _ Dialector = &MssqlDialector{}

func (d *MssqlDialector) String() string {
	return "SQL Server"
}

func (d *MssqlDialector) Open(dsn *driverV2.DSN) (*sql.DB, *sql.Conn, error) {
	return d.BaseDialector.GetConn("sqlserver", fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.DatabaseName))
}

func (d *MssqlDialector) ShowDatabaseSQL() string {
	return "select name from sys.databases"
}
