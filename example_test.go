package minisql_test

import (
	"fmt"

	"github.com/bsm/minisql"
)

func Example() {
	query := new(minisql.Query)
	query.AppendString(`SELECT * FROM users WHERE id > `)
	query.AppendValue(33)
	query.AppendString(` LIMIT `)
	query.AppendInt(100)

	fmt.Println(query.SQL())
	fmt.Println(query.Args())

	// Output:
	// SELECT * FROM users WHERE id > $1 LIMIT 100
	// [33]
}
