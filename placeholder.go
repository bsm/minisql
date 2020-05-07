package minisql

// PlaceholderFormat represents a specific format.
type PlaceholderFormat uint8

const (
	// Dollar placeholder e.g. $1, $2, $3, commonly used by PostgreSQL.
	Dollar PlaceholderFormat = iota
	// Question placeholder (e.g. ?, ?, ?), commonly used by MySQL.
	Question
	// Colon placeholder (e.g. :1, :2, :3), commonly used by Oracle.
	Colon
)
