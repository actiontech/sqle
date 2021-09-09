package reform

import (
	"context"
	"database/sql"
	"time"
)

// TXInterface is a subset of *sql.Tx used by reform.
// Can be used together with NewTXFromInterface for easier integration with existing code or for passing test doubles.
//
// It may grow and shrink over time to include only needed *sql.Tx methods,
// and is excluded from SemVer compatibility guarantees.
type TXInterface interface {
	DBTXContext
	Commit() error
	Rollback() error

	// Deprecated: do not use, it will be removed in v1.6.
	DBTX
}

// check interface
var _ TXInterface = (*sql.Tx)(nil)

// TX represents a SQL database transaction.
type TX struct {
	*Querier
	tx TXInterface
}

// NewTX creates new TX object for given SQL database transaction.
// Logger can be nil.
func NewTX(tx *sql.Tx, dialect Dialect, logger Logger) *TX {
	return NewTXFromInterface(tx, dialect, logger)
}

// NewTXFromInterface creates new TX object for given TXInterface.
// Can be used for easier integration with existing code or for passing test doubles.
// Logger can be nil.
func NewTXFromInterface(tx TXInterface, dialect Dialect, logger Logger) *TX {
	return newTX(context.Background(), tx, dialect, logger)
}

func newTX(ctx context.Context, tx TXInterface, dialect Dialect, logger Logger) *TX {
	return &TX{
		Querier: newQuerier(ctx, tx, "", dialect, logger),
		tx:      tx,
	}
}

// Commit commits the transaction.
func (tx *TX) Commit() error {
	tx.logBefore("COMMIT", nil)
	start := time.Now()
	err := tx.tx.Commit()
	tx.logAfter("COMMIT", nil, time.Since(start), err)
	return err
}

// Rollback aborts the transaction.
func (tx *TX) Rollback() error {
	tx.logBefore("ROLLBACK", nil)
	start := time.Now()
	err := tx.tx.Rollback()
	tx.logAfter("ROLLBACK", nil, time.Since(start), err)
	return err
}

// check interfaces
var (
	_ DBTX        = (*TX)(nil)
	_ DBTXContext = (*TX)(nil)
)
