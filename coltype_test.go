package minisql_test

import (
	"testing"
	"time"

	. "github.com/bsm/minisql"
)

func TestColumnType_Parse(t *testing.T) {
	ts := time.Date(2026, 11, 10, 14, 30, 0, 500_000_000, time.FixedZone("", 3600))

	tests := []struct {
		typ ColumnType
		in  string
		exp any // nil means an error is expected
	}{
		{String, "a'b", "a'b"},
		{Integer, "-42", int64(-42)},
		{Integer, "4.2", nil},
		{Float, "4.2", 4.2},
		{Float, "abc", nil},
		{Float, "NaN", nil},
		{Float, "-Inf", nil},
		{Decimal, "-12345678901234567890.0000000001", "-12345678901234567890.0000000001"},
		{Decimal, ".5", ".5"},
		{Decimal, "1e5", nil},
		{Decimal, "1.2.3", nil},
		{Decimal, "-", nil},
		{Decimal, ".", nil},
		{Boolean, "true", true},
		{Boolean, "0", false},
		{Boolean, "yes", nil},
		{Date, "2026-11-10", "2026-11-10"},
		{Date, "2026-02-30", nil},
		{Date, "10/11/2026", nil},
		{Timestamp, "2026-11-10T14:30:00.5+01:00", ts},
		{Timestamp, "2026-11-10 14:30:00", nil},
		{ColumnType(99), "x", nil},
	}
	for _, tt := range tests {
		got, err := tt.typ.Parse(tt.in)
		switch {
		case tt.exp == nil && err == nil:
			t.Errorf("%s %q: expected an error, got %#v", tt.typ, tt.in, got)
		case tt.exp == nil:
		case err != nil:
			t.Errorf("%s %q: unexpected error %v", tt.typ, tt.in, err)
		case tt.typ == Timestamp:
			if !got.(time.Time).Equal(tt.exp.(time.Time)) {
				t.Errorf("%s %q: expected %v, got %v", tt.typ, tt.in, tt.exp, got)
			}
		case got != tt.exp:
			t.Errorf("%s %q: expected %#v, got %#v", tt.typ, tt.in, tt.exp, got)
		}
	}
}

func TestParserFunc(t *testing.T) {
	exclaim := ParserFunc(func(s string) (any, error) { return s + "!", nil })
	if got, _ := exclaim.Parse("a"); got != "a!" {
		t.Errorf("expected %q, got %#v", "a!", got)
	}
}
