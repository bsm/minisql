// Package minisql is a minimal SQL query builder.
//
// A Query collects raw SQL and arguments. Values never go into the SQL text:
// each one is added to the argument list and replaced by a placeholder in
// the format of the chosen Dialect.
//
// The package only handles what differs between databases or is easy to get
// wrong: placeholders, identifier quoting, list separators and turning text
// input into typed values. It never writes SQL keywords or clauses; a helper
// that only saves the caller one AppendString does not belong here.
package minisql

import (
	"context"
	"database/sql"
	"strconv"
)

// Execer interface, may apply to an sql.DB or an sql.Tx.
type Execer interface {
	// ExecContext executes a query that doesn't return rows. For example: an INSERT and UPDATE.
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

// Querier interface, may apply to an sql.DB or an sql.Tx.
type Querier interface {
	// QueryContext executes a query that returns rows, typically a SELECT.
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

// RowQuerier interface, may apply to an sql.DB or an sql.Tx.
type RowQuerier interface {
	// QueryRowContext executes a query that is expected to return at most one row. QueryRowContext always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// Query builder. The zero value is an empty Postgres query.
type Query struct {
	Dialect Dialect
	stmt    []byte
	args    []any
}

// UseDialect sets the dialect and returns the query, so it can be chained:
//
//	q := minisql.Pooled().UseDialect(minisql.MySQL)
func (q *Query) UseDialect(d Dialect) *Query {
	q.Dialect = d
	return q
}

// Reset clears the SQL, the arguments and the dialect, but keeps the
// allocated memory.
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
func (q *Query) Args() []any {
	return q.args
}

// AppendString appends a raw SQL string. It is not escaped, so never pass
// user input here: use AppendValue or AppendIdent instead.
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

// AppendValue appends an argument value as a placeholder, e.g. $1.
func (q *Query) AppendValue(value any) {
	q.appendPlaceholder()
	q.args = append(q.args, value)
}

// AppendValues appends a comma-separated list of argument values, e.g.
// $1, $2, $3, for IN (...) and VALUES (...) lists. It appends nothing when
// there are no values; note that IN () is not valid SQL.
func (q *Query) AppendValues(values ...any) {
	for i, v := range values {
		if i > 0 {
			q.AppendString(", ")
		}
		q.AppendValue(v)
	}
}

// AppendParsed parses each text value with p and appends them as a
// comma-separated list of placeholders, like AppendValues:
//
//	err := q.AppendParsed(minisql.Integer, "1", "2") // $1, $2 with int64 args
//
// If any value fails to parse, it returns the error and appends nothing.
func (q *Query) AppendParsed(p Parser, values ...string) error {
	parsed := make([]any, len(values))
	for i, s := range values {
		v, err := p.Parse(s)
		if err != nil {
			return err
		}
		parsed[i] = v
	}
	q.AppendValues(parsed...)
	return nil
}

// AppendIdent appends a quoted identifier, such as a table or column name.
// Several parts are joined with a dot into a qualified name:
//
//	q.AppendIdent("users")           // "users"
//	q.AppendIdent("public", "users") // "public"."users"
//
// Quote characters inside a part are escaped, so the result is always a
// single identifier. Note that quoted identifiers are case-sensitive in
// Postgres, SQLite and Oracle.
func (q *Query) AppendIdent(parts ...string) {
	for i, part := range parts {
		if i > 0 {
			q.AppendByte('.')
		}
		q.appendQuoted(part)
	}
}

// AppendColumns appends a comma-separated list of quoted identifiers, e.g.
// "id", "name", "email".
func (q *Query) AppendColumns(names ...string) {
	for i, name := range names {
		if i > 0 {
			q.AppendString(", ")
		}
		q.appendQuoted(name)
	}
}

// appendQuoted wraps name in the dialect's quote character and escapes the
// quote character inside it by doubling it, e.g. a"b becomes "a""b".
func (q *Query) appendQuoted(name string) {
	quote := q.Dialect.identQuote()
	q.AppendByte(quote)
	for i := range len(name) {
		if name[i] == quote {
			q.AppendByte(quote)
		}
		q.AppendByte(name[i])
	}
	q.AppendByte(quote)
}

func (q *Query) appendPlaceholder() {
	switch q.Dialect {
	case MySQL, SQLite:
		q.AppendByte('?')
	case Oracle:
		q.AppendByte(':')
		q.AppendInt(int64(len(q.args)) + 1)
	default:
		q.AppendByte('$')
		q.AppendInt(int64(len(q.args)) + 1)
	}
}
