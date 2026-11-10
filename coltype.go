package minisql

import (
	"fmt"
	"strconv"
	"time"
)

// Parser turns a text value into the Go value the driver binds. ColumnType
// is the built-in one; implement it for anything else, such as a UUID or
// JSON type of your driver.
type Parser interface {
	Parse(s string) (any, error)
}

// ParserFunc lets an ordinary function be used as a Parser.
type ParserFunc func(s string) (any, error)

// Parse calls f(s).
func (f ParserFunc) Parse(s string) (any, error) { return f(s) }

// ColumnType tells how to turn a text value into the Go value the driver
// binds for a column. Use it for input that arrives as text, such as URL
// query parameters or filter values in a JSON body. The zero value is String.
type ColumnType uint8

const (
	String    ColumnType = iota // kept as it is
	Integer                     // int64
	Float                       // float64
	Decimal                     // checked to be e.g. -12.50, bound as a string
	Boolean                     // bool, e.g. "true", "false", "1", "0"
	Date                        // checked to be YYYY-MM-DD, bound as a string
	Timestamp                   // RFC 3339, e.g. 2026-11-10T14:30:00Z, as time.Time
)

func (t ColumnType) String() string {
	switch t {
	case Integer:
		return "integer"
	case Float:
		return "float"
	case Decimal:
		return "decimal"
	case Boolean:
		return "boolean"
	case Date:
		return "date"
	case Timestamp:
		return "timestamp"
	default:
		return "string"
	}
}

// Parse converts s into the Go value for the column type, e.g. "42" into
// int64(42) for Integer.
//
// Some types stay strings after the check, so the database parses them:
//   - Decimal, because float64 loses digits, e.g. on money or numeric(38,10).
//   - Date, because a time.Time carries a time zone, and the driver could
//     shift it to the day before or after.
//
// Timestamp needs no such care: with its offset it is one fixed instant.
func (t ColumnType) Parse(s string) (any, error) {
	switch t {
	case Integer:
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return n, nil
		}
	case Float:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f, nil
		}
	case Decimal:
		if isDecimal(s) {
			return s, nil
		}
	case Boolean:
		if b, err := strconv.ParseBool(s); err == nil {
			return b, nil
		}
	case Date:
		if _, err := time.Parse(time.DateOnly, s); err == nil {
			return s, nil
		}
	case Timestamp:
		// RFC3339Nano also accepts values without fractional seconds.
		if ts, err := time.Parse(time.RFC3339Nano, s); err == nil {
			return ts, nil
		}
	default:
		return s, nil
	}
	return nil, fmt.Errorf("minisql: %q is not a valid %s", s, t)
}

// isDecimal reports whether s is a plain decimal number: an optional sign,
// digits and an optional fraction, e.g. "-12.50", "3", "0.5" or ".5".
// Exponents, NaN and Infinity are not accepted, as databases differ on them.
func isDecimal(s string) bool {
	if s != "" && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}
	digits, dots := 0, 0
	for i := range len(s) {
		switch c := s[i]; {
		case c >= '0' && c <= '9':
			digits++
		case c == '.' && dots == 0:
			dots++
		default:
			return false
		}
	}
	return digits > 0
}
