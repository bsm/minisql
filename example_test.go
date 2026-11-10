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

func Example_insert() {
	query := new(minisql.Query).UseDialect(minisql.MySQL)
	query.AppendString(`INSERT INTO `)
	query.AppendIdent("app", "users")
	query.AppendString(` (`)
	query.AppendColumns("name", "settings")
	query.AppendString(`) VALUES (`)
	query.AppendValues("Alice", `{"theme":"dark"}`)
	query.AppendByte(')')

	fmt.Println(query.SQL())
	fmt.Println(query.Args())

	// Output:
	// INSERT INTO `app`.`users` (`name`, `settings`) VALUES (?, ?)
	// [Alice {"theme":"dark"}]
}

func Example_filter() {
	// Filter values often arrive as text, e.g. ?id=7&id=9 in a URL.
	query := new(minisql.Query)
	query.AppendString(`SELECT * FROM users WHERE id IN (`)
	if err := query.AppendParsed(minisql.Integer, "7", "9"); err != nil {
		panic(err)
	}
	query.AppendByte(')')

	fmt.Println(query.SQL())
	fmt.Printf("%#v\n", query.Args())

	// Output:
	// SELECT * FROM users WHERE id IN ($1, $2)
	// []interface {}{7, 9}
}
