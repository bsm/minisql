package minisql_test

import (
	"testing"

	. "github.com/bsm/minisql"
)

func TestPooled(t *testing.T) {
	q1 := Pooled().UseDialect(MySQL)
	q1.AppendValue(1)

	Release(q1)
	q2 := Pooled()
	if q2.Dialect != Postgres || q2.SQL() != "" || len(q2.Args()) != 0 {
		t.Fatalf("expected a clean query, got %+v", q2)
	}
}
