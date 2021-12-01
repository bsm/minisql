# MiniSQL

[![Go Reference](https://pkg.go.dev/badge/github.com/bsm/minisql.svg)](https://pkg.go.dev/github.com/bsm/minisql)
[![Test](https://github.com/bsm/minisql/actions/workflows/test.yml/badge.svg)](https://github.com/bsm/minisql/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A minimal SQL query builder for [Go](https://golang.org/).

## Usage

```go
func GetUserName(ctx context.Context, db *sql.DB, userID int64) (string, error) {
  query := minisql.Pooled()
  defer minisql.Release(query)

  query.AppendString(`SELECT user.name FROM users WHERE id = `)
  query.AppendValue(userID)

  var name string
  err := query.QueryRowContext(ctx, db).Scan(&name)
  return name, err
}
```

## Documentation

Please see the [API documentation](https://pkg.go.dev/github.com/bsm/minisql) for package and API descriptions and examples.
