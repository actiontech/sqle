package driver

import (
	"context"
	"database/sql"
	"fmt"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pkg/errors"

	// DRIVER LIST:
	// 	https://github.com/golang/go/wiki/SQLDrivers
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/sijms/go-ora/v2"
)

// Dialector is a interface for database dialect. It used for sql.Open()
type Dialector interface {
	// Dialect return the driver name and dsn detail. The return value is used for sql.Open().
	Dialect(dsn *driverV2.DSN) (driverName string, dsnDetail string)

	// ShowDatabaseSQL return the sql to show all databases.
	ShowDatabaseSQL() string

	// String return the dialect name with more formal name. It is different from driver name.
	// For example, "PostgreSQL" is more formal name than "pgx".
	String() string
}

type PostgresDialector struct {
}

func (d *PostgresDialector) Dialect(dsn *driverV2.DSN) (string, string) {
	if dsn.DatabaseName == "" {
		dsn.DatabaseName = "postgres"
	}

	return "pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.DatabaseName)
}

func (d *PostgresDialector) String() string {
	return "PostgreSQL"
}

func (d *PostgresDialector) ShowDatabaseSQL() string {
	return "select datname from pg_database"
}

type OracleDialector struct {
}

func (d *OracleDialector) Dialect(dsn *driverV2.DSN) (string, string) {
	if dsn.DatabaseName == "" {
		dsn.DatabaseName = "xe"
	}
	return "oracle", fmt.Sprintf("oracle://%s:%s@%s:%s/%s",
		dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.DatabaseName)
}

func (d *OracleDialector) String() string {
	return "Oracle"
}

func (d *OracleDialector) ShowDatabaseSQL() string {
	return "select global_name from global_name"
}

type MssqlDialector struct {
}

func (d *MssqlDialector) Dialect(dsn *driverV2.DSN) (string, string) {
	// connect by:
	// 1. host and port (we used)
	// 2. host and instance
	return "sqlserver", fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.DatabaseName)
}

func (d *MssqlDialector) String() string {
	return "SQL Server"
}

func (d *MssqlDialector) ShowDatabaseSQL() string {
	return "select name from sys.databases"
}

func getDbConn(driverName, dsnDetail string) (db *sql.DB, conn *sql.Conn, err error) {
	db, err = sql.Open(driverName, dsnDetail)
	if err != nil {
		return nil, nil, errors.Wrap(err, "open database failed when new driver")
	}
	conn, err = db.Conn(context.TODO())
	if err != nil {
		db.Close()
		return nil, nil, errors.Wrap(err, "get database connection failed when new driver")
	}
	if err := conn.PingContext(context.TODO()); err != nil {
		db.Close()
		return nil, nil, errors.Wrap(err, "ping database connection failed when new driver")
	}
	return
}
