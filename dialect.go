package minisql

import "strconv"

// Dialect selects the placeholder format and the identifier quoting of a
// database. The zero value is Postgres.
//
//	Dialect   Placeholders    Identifiers
//	Postgres  $1, $2, $3      "name"
//	MySQL     ?, ?, ?         `name`
//	SQLite    ?, ?, ?         "name"
//	Oracle    :1, :2, :3      "name"
type Dialect uint8

const (
	Postgres Dialect = iota
	MySQL
	SQLite
	Oracle
)

// appendPlaceholder appends the placeholder for argument number n (from 1).
func (d Dialect) appendPlaceholder(b []byte, n int) []byte {
	switch d {
	case MySQL, SQLite:
		return append(b, '?')
	case Oracle:
		return strconv.AppendInt(append(b, ':'), int64(n), 10)
	default:
		return strconv.AppendInt(append(b, '$'), int64(n), 10)
	}
}

// identQuote returns the character that wraps an identifier.
func (d Dialect) identQuote() byte {
	if d == MySQL {
		return '`'
	}
	return '"'
}
