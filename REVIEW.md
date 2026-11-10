# REVIEW.md

Review against `AGENTS.md`. Focus on what tests and linters miss.

## Flag

- A value written into the SQL text instead of bound as an argument.
- Identifier quoting that a name can break out of, in any `Dialect`.
- Wrong placeholder numbers, e.g. after mixing `AppendValue` and
  `AppendParsed`, or a failed `AppendParsed` that changes the query.
- Parsing that loses data, e.g. a decimal through float64.
- Helpers that write SQL keywords or clauses (see the package doc).
- Unflagged breaking changes to exported API.

## Skip

- Formatting and anything `golangci-lint` reports.
- Naming, unless a name is misleading.
- Requests for tests of trivial paths.

## Format

One comment per issue, with a severity: XS nit, S minor, M should fix,
L bug, XL serious bug. Say what breaks and for which input.
