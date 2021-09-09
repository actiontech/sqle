package reform

import (
	"context"
	"database/sql"
	"time"
)

// DBInterface is a subset of *sql.DB used by reform.
// Can be used together with NewDBFromInterface for easier integration with existing code or for passing test doubles.
//
// It may grow and shrink over time to include only needed *sql.DB methods,
// and is excluded from SemVer compatibility guarantees.
type DBInterface interface {
	DBTXContext
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

	// Deprecated: do not use, it will be removed in v1.6.
	DBTX
	// Deprecated: do not use, it will be removed in v1.6.
	Begin() (*sql.Tx, error)
}

// check interface
var _ DBInterface = (*sql.DB)(nil)

// DB represents a connection to SQL database.
type DB struct {
	*Querier
	db DBInterface
}

// NewDB creates new DB object for given SQL database connection.
// Logger can be nil.
func NewDB(db *sql.DB, dialect Dialect, logger Logger) *DB {
	return NewDBFromInterface(db, dialect, logger)
}

// NewDBFromInterface creates new DB object for given DBInterface.
// Can be used for easier integration with existing code or for passing test doubles.
// Logger can be nil.
func NewDBFromInterface(db DBInterface, dialect Dialect, logger Logger) *DB {
	return &DB{
		Querier: newQuerier(context.Background(), db, "", dialect, logger),
		db:      db,
	}
}

// DBInterface returns DBInterface associated with a given DB object.
func (db *DB) DBInterface() DBInterface {
	return db.db
}

// Begin starts transaction with Querier's context and default options.
func (db *DB) Begin() (*TX, error) {
	return db.BeginTx(db.Querier.ctx, nil)
}

// BeginTx starts transaction with given context and options (can be nil).
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TX, error) {
	db.logBefore("BEGIN", nil)
	start := time.Now()
	tx, err := db.db.BeginTx(ctx, opts)
	db.logAfter("BEGIN", nil, time.Since(start), err)
	if err != nil {
		return nil, err
	}
	return newTX(ctx, tx, db.Dialect, db.Logger), nil
}

// InTransaction wraps function execution in transaction with Querier's context and default options,
// rolling back it in case of error or panic, committing otherwise.
func (db *DB) InTransaction(f func(t *TX) error) error {
	return db.InTransactionContext(db.Querier.ctx, nil, f)
}

// InTransactionContext wraps function execution in transaction with given context and options (can be nil),
// rolling back it in case of error or panic, committing otherwise.
func (db *DB) InTransactionContext(ctx context.Context, opts *sql.TxOptions, f func(t *TX) error) error {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	var committed bool
	defer func() {
		if !committed {
			// always return f() or Commit() error, not possible Rollback() error
			_ = tx.Rollback()
		}
	}()

	err = f(tx)
	if err == nil {
		err = tx.Commit()
	}
	if err == nil {
		committed = true
	}
	return err
}

// check interfaces
var (
	_ DBTX        = (*DB)(nil)
	_ DBTXContext = (*DB)(nil)
)
