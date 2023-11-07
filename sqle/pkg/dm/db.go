package dm

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	_ "dm"
)

const DriverName = "dm"

type DSN struct {
	Host     string
	Port     string
	User     string
	Password string
}

// NewDB :创建数据库连接
func NewDB(dsn *DSN) (*sql.DB, error) {
	// dm://SYSDBA:SYSDBA@localhost:5236
	dataSourceName := fmt.Sprintf("dm://%s:%s@%s:%s", dsn.User, dsn.Password, dsn.Host, dsn.Port)
	var db *sql.DB
	var err error
	if db, err = sql.Open(DriverName, dataSourceName); err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Close(db *sql.DB) error {
	return db.Close()
}

func QueryCurrentSchema(ctx context.Context, db *sql.DB) (string, error) {
	query := `SELECT SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA') AS schema_name FROM dual`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return "", errors.Wrapf(err, "failed to query %s", query)
	}
	defer rows.Close()

	schema := ""
	if rows.Next() {
		err = rows.Scan(&schema)
		if err != nil {
			return "", errors.Wrapf(err, "failed to scan %s", query)
		}
	}
	return schema, nil
}
