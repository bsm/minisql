package minisql_test

import (
	"reflect"
	"testing"

	. "github.com/bsm/minisql"
)

func TestQuery(t *testing.T) {
	q := new(Query)
	q.AppendString(`SELECT user.name FROM users WHERE id > `)
	q.AppendValue(33)
	q.AppendString(` LIMIT `)
	q.AppendInt(100)

	if exp, got := `SELECT user.name FROM users WHERE id > $1 LIMIT 100`, q.SQL(); exp != got {
		t.Errorf("expected %q, got %q", exp, got)
	}
	if exp, got := []any{33}, q.Args(); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestQuery_dialects(t *testing.T) {
	tests := []struct {
		dialect Dialect
		exp     string
	}{
		{Postgres, `INSERT INTO "app"."users" ("id", "na""me") VALUES ($1, $2)`},
		{MySQL, "INSERT INTO `app`.`users` (`id`, `na\"me`) VALUES (?, ?)"},
		{SQLite, `INSERT INTO "app"."users" ("id", "na""me") VALUES (?, ?)`},
		{Oracle, `INSERT INTO "app"."users" ("id", "na""me") VALUES (:1, :2)`},
	}
	for _, tt := range tests {
		q := new(Query).UseDialect(tt.dialect)
		q.AppendString(`INSERT INTO `)
		q.AppendIdent("app", "users")
		q.AppendString(` (`)
		q.AppendColumns("id", `na"me`)
		q.AppendString(`) VALUES (`)
		q.AppendValues(1, "x")
		q.AppendByte(')')

		if got := q.SQL(); got != tt.exp {
			t.Errorf("dialect %d: expected %q, got %q", tt.dialect, tt.exp, got)
		}
	}
}

func TestQuery_AppendIdent_escapesMySQL(t *testing.T) {
	q := new(Query).UseDialect(MySQL)
	q.AppendIdent("a`b")
	if exp, got := "`a``b`", q.SQL(); exp != got {
		t.Errorf("expected %q, got %q", exp, got)
	}
}

func TestQuery_AppendValues(t *testing.T) {
	q := new(Query)
	q.AppendValue("x")
	q.AppendString(` IN (`)
	q.AppendValues(1, 2, 3)
	q.AppendByte(')')

	if exp, got := `$1 IN ($2, $3, $4)`, q.SQL(); exp != got {
		t.Errorf("expected %q, got %q", exp, got)
	}
	if exp, got := []any{"x", 1, 2, 3}, q.Args(); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestQuery_AppendParsed(t *testing.T) {
	q := new(Query)
	if err := q.AppendParsed(Integer, "1", "2"); err != nil {
		t.Fatal(err)
	}
	if exp, got := `$1, $2`, q.SQL(); exp != got {
		t.Errorf("expected %q, got %q", exp, got)
	}
	if exp, got := []any{int64(1), int64(2)}, q.Args(); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}

	// a bad value leaves the query as it was
	if err := q.AppendParsed(Integer, "3", "x"); err == nil {
		t.Fatal("expected an error")
	}
	if exp, got := `$1, $2`, q.SQL(); exp != got {
		t.Errorf("expected %q, got %q", exp, got)
	}
	if n := len(q.Args()); n != 2 {
		t.Errorf("expected 2 args, got %d", n)
	}
}
