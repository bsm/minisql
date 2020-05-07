# MiniSQL

[![GoDoc](https://godoc.org/github.com/bsm/minisql?status.svg)](https://godoc.org/github.com/bsm/minisql)
[![Build Status](https://travis-ci.org/bsm/minisql.svg)](https://travis-ci.org/bsm/minisql)
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

Please see the [API documentation](https://godoc.org/github.com/bsm/minisql) for package and API descriptions and examples.
