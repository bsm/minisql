package minisql_test

import (
	"testing"

	"github.com/bsm/minisql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query", func() {
	It("should pool", func() {
		q1 := minisql.Pooled()
		q2 := minisql.Pooled()
		Expect(q2).ToNot(BeIdenticalTo(q1))

		minisql.Release(q1)
		q3 := minisql.Pooled()
		Expect(q3).To(BeIdenticalTo(q1))
	})

	It("should generate query", func() {
		q := new(minisql.Query)
		q.AppendString(`SELECT user.name FROM users WHERE id > `)
		q.AppendValue(33)
		q.AppendString(` LIMIT `)
		q.AppendInt(100)
		Expect(q.SQL()).To(Equal(`SELECT user.name FROM users WHERE id > $1 LIMIT 100`))
		Expect(q.Args()).To(Equal([]interface{}{33}))
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "minisql")
}
