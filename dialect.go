package minisql

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

// identQuote returns the character that wraps an identifier.
func (d Dialect) identQuote() byte {
	if d == MySQL {
		return '`'
	}
	return '"'
}
