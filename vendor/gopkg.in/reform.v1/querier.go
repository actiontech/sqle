package reform

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Querier performs queries and commands.
type Querier struct {
	ctx     context.Context
	dbtxCtx DBTXContext
	tag     string
	Dialect
	Logger Logger
}

func newQuerier(ctx context.Context, dbtxCtx DBTXContext, tag string, dialect Dialect, logger Logger) *Querier {
	return &Querier{
		ctx:     ctx,
		dbtxCtx: dbtxCtx,
		tag:     tag,
		Dialect: dialect,
		Logger:  logger,
	}
}

func (q *Querier) clone() *Querier {
	return newQuerier(q.ctx, q.dbtxCtx, q.tag, q.Dialect, q.Logger)
}

func (q *Querier) logBefore(query string, args []interface{}) {
	if q.Logger != nil {
		q.Logger.Before(query, args)
	}
}

func (q *Querier) logAfter(query string, args []interface{}, d time.Duration, err error) {
	if q.Logger != nil {
		q.Logger.After(query, args, d, err)
	}
}

func (q *Querier) startQuery(command string) string {
	if q.tag == "" {
		return command
	}
	return command + " /* " + q.tag + " */"
}

// Tag returns Querier's tag. Default tag is empty.
func (q *Querier) Tag() string {
	return q.tag
}

// WithTag returns a copy of Querier with set tag. Returned Querier is tied to the same DB or TX.
// See Tagging section in documentation for details.
func (q *Querier) WithTag(format string, args ...interface{}) *Querier {
	newQ := q.clone()
	if len(args) == 0 {
		newQ.tag = format
	} else {
		newQ.tag = fmt.Sprintf(format, args...)
	}
	return newQ
}

// QualifiedView returns quoted qualified view name.
func (q *Querier) QualifiedView(view View) string {
	v := q.QuoteIdentifier(view.Name())
	if view.Schema() != "" {
		v = q.QuoteIdentifier(view.Schema()) + "." + v
	}
	return v
}

// Context returns Querier's context. Default context is context.Background().
func (q *Querier) Context() context.Context {
	return q.ctx
}

// WithContext returns a copy of Querier with set context. Returned Querier is tied to the same DB or TX.
// See Context section in documentation for details.
func (q *Querier) WithContext(ctx context.Context) *Querier {
	newQ := q.clone()
	newQ.ctx = ctx
	return newQ
}

// QualifiedColumns returns a slice of quoted qualified column names for given view.
func (q *Querier) QualifiedColumns(view View) []string {
	v := q.QualifiedView(view)
	res := view.Columns()
	for i := 0; i < len(res); i++ {
		res[i] = v + "." + q.QuoteIdentifier(res[i])
	}
	return res
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (q *Querier) Exec(query string, args ...interface{}) (sql.Result, error) {
	q.logBefore(query, args)
	start := time.Now()
	res, err := q.dbtxCtx.ExecContext(q.ctx, query, args...)
	q.logAfter(query, args, time.Since(start), err)
	return res, err
}

// ExecContext just calls q.WithContext(ctx).Exec(query, args...), and that form should be used instead.
// This method exists to satisfy various standard interfaces for advanced use-cases.
func (q *Querier) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return q.WithContext(ctx).Exec(query, args...)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (q *Querier) Query(query string, args ...interface{}) (*sql.Rows, error) {
	q.logBefore(query, args)
	start := time.Now()
	rows, err := q.dbtxCtx.QueryContext(q.ctx, query, args...)
	q.logAfter(query, args, time.Since(start), err)
	return rows, err
}

// QueryContext just calls q.WithContext(ctx).Query(query, args...), and that form should be used instead.
// This method exists to satisfy various standard interfaces for advanced use-cases.
func (q *Querier) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return q.WithContext(ctx).Query(query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called.
func (q *Querier) QueryRow(query string, args ...interface{}) *sql.Row {
	q.logBefore(query, args)
	start := time.Now()
	row := q.dbtxCtx.QueryRowContext(q.ctx, query, args...)
	q.logAfter(query, args, time.Since(start), nil)
	return row
}

// QueryRowContext just calls q.WithContext(ctx).QueryRow(query, args...), and that form should be used instead.
// This method exists to satisfy various standard interfaces for advanced use-cases.
func (q *Querier) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return q.WithContext(ctx).QueryRow(query, args...)
}

// check interfaces
var (
	_ DBTX        = (*Querier)(nil)
	_ DBTXContext = (*Querier)(nil)
)
