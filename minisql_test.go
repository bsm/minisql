package minisql_test

import (
	"reflect"
	"testing"

	. "github.com/bsm/minisql"
)

func TestPooled(t *testing.T) {
	q1 := Pooled()
	q2 := Pooled()
	if q1 == q2 {
		t.Fatal("expected different")
	}

	Release(q1)
	if q3 := Pooled(); q1 != q3 {
		t.Fatal("expected same")
	}
}

func TestQuery(t *testing.T) {
	q := new(Query)
	q.AppendString(`SELECT user.name FROM users WHERE id > `)
	q.AppendValue(33)
	q.AppendString(` LIMIT `)
	q.AppendInt(100)

	if exp, got := `SELECT user.name FROM users WHERE id > $1 LIMIT 100`, q.SQL(); exp != got {
		t.Errorf("expected %q, got %q", exp, got)
	}
	if exp, got := []interface{}{33}, q.Args(); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
