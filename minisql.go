package minisql

import (
	"context"
	"database/sql"
	"strconv"
)

// Execer interface, may apply to an sql.DB or an sql.Tx.
type Execer interface {
	// ExecContext executes a query that doesn't return rows. For example: an INSERT and UPDATE.
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}

// Querier interface, may apply to an sql.DB or an sql.Tx.
type Querier interface {
	// QueryContext executes a query that returns rows, typically a SELECT.
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
}

// RowQuerier interface, may apply to an sql.DB or an sql.Tx.
type RowQuerier interface {
	// QueryRowContext executes a query that is expected to return at most one row. QueryRowContext always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// Query builder.
type Query struct {
	Placeholder PlaceholderFormat
	stmt        []byte
	args        []interface{}
}

// UsePlaceholder set a custom placeholder format.
func (q *Query) UsePlaceholder(f PlaceholderFormat) *Query {
	q.Placeholder = f
	return q
}

// Reset resets the query.
func (q *Query) Reset() {
	*q = Query{stmt: q.stmt[:0], args: q.args[:0]}
}

// ExecContext executes the query.
func (q *Query) ExecContext(ctx context.Context, target Execer) (sql.Result, error) {
	return target.ExecContext(ctx, q.SQL(), q.args...)
}

// QueryContext queries sql.Rows.
func (q *Query) QueryContext(ctx context.Context, target Querier) (*sql.Rows, error) {
	return target.QueryContext(ctx, q.SQL(), q.args...)
}

// QueryRowContext queries an sql.Row.
func (q *Query) QueryRowContext(ctx context.Context, target RowQuerier) *sql.Row {
	return target.QueryRowContext(ctx, q.SQL(), q.args...)
}

// SQL exposes the raw SQL.
func (q *Query) SQL() string {
	return string(q.stmt)
}

// Args exposes the collected arguments.
func (q *Query) Args() []interface{} {
	return q.args
}

// AppendString appends a raw SQL string.
func (q *Query) AppendString(str string) {
	q.stmt = append(q.stmt, str...)
}

// AppendByte appends a raw byte.
func (q *Query) AppendByte(c byte) {
	q.stmt = append(q.stmt, c)
}

// AppendInt appends a single raw integer.
func (q *Query) AppendInt(n int64) {
	q.stmt = strconv.AppendInt(q.stmt, n, 10)
}

// AppendValue appends an argument value.
func (q *Query) AppendValue(value interface{}) {
	q.appendPlacholder()
	q.args = append(q.args, value)
}

func (q *Query) appendPlacholder() {
	switch q.Placeholder {
	case Question:
		q.AppendByte('?')
	case Colon:
		q.AppendByte(':')
		q.AppendInt(int64(len(q.args)) + 1)
	default:
		q.AppendByte('$')
		q.AppendInt(int64(len(q.args)) + 1)
	}
}
