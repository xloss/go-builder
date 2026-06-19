# go-builder

`go-builder` is a small PostgreSQL query builder for Go.

It is designed for projects that use [`pgx`](https://github.com/jackc/pgx) and `pgx.NamedArgs`. The generated SQL uses named placeholders such as `@id`, `@limit`, and `@offset`.

The package is intentionally small. It is not an ORM, not a schema mapper, and not a general-purpose SQL abstraction layer for every database engine.

## Features

* `SELECT`
* `INSERT`
* `UPDATE`
* `DELETE`
* `LEFT JOIN`
* `WHERE`
* `AND` / `OR` where groups
* `EXISTS`
* `GROUP BY`
* `ORDER BY`
* `LIMIT`
* `OFFSET`
* `RETURNING`
* `ON CONFLICT DO UPDATE`
* `ON CONFLICT DO NOTHING`
* JSONB existence checks
* PostgreSQL full-text search helpers
* Subqueries in `FROM`
* Named bind arguments for query values

## Installation

```bash
go get github.com/xloss/go-builder
```

## Basic usage

```go
package main

import (
	"context"

	builder "github.com/xloss/go-builder"
	"github.com/jackc/pgx/v5"
)

func example(ctx context.Context, db interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}) error {
	users := builder.NewTable("users")

	sql, args, err := builder.NewSelect().
		From(users).
		Column(builder.ColumnName{Table: users, Name: "id"}).
		Column(builder.ColumnName{Table: users, Name: "email"}).
		Where(builder.WhereEq{Table: users, Column: "id", Value: 10}).
		Limit(1).
		Get()
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, sql, pgx.NamedArgs(args))
	if err != nil {
		return err
	}
	defer rows.Close()

	return rows.Err()
}
```

Generated SQL will look similar to this:

```sql
SELECT users_xxxxxxxxxx.id, users_xxxxxxxxxx.email
FROM users AS users_xxxxxxxxxx
WHERE users_xxxxxxxxxx.id = @id_xxxxxxxxxx
LIMIT @limit_xxxxxxxxxx
```

The exact aliases and bind names are generated automatically.

## Security model

`go-builder` separates query values from SQL text by using named bind parameters.

Safe:

```go
q.Where(builder.WhereEq{
	Table:  users,
	Column: "id",
	Value:  userID,
})
```

The value is added to the bind map and is not inserted directly into the SQL string.

However, SQL identifiers are different from values. Table names, column names, aliases, and similar SQL identifiers cannot be passed as normal query arguments in PostgreSQL.

`go-builder` validates identifiers using a conservative identifier format. Still, identifiers should normally come from trusted application code, not directly from user input.

Do not do this:

```go
table := builder.NewTable(userInput)
```

Prefer constants or controlled mappings:

```go
var allowedSortColumns = map[string]string{
	"name":       "name",
	"created_at": "created_at",
}

column, ok := allowedSortColumns[userSort]
if !ok {
	return errors.New("invalid sort column")
}

q.Order(builder.Order{Table: users, Column: column})
```

## SELECT

```go
users := builder.NewTable("users")

sql, args, err := builder.NewSelect().
	From(users).
	Column(builder.ColumnName{Table: users, Name: "id"}).
	Column(builder.ColumnName{Table: users, Name: "name"}).
	Where(builder.WhereEq{Table: users, Column: "active", Value: true}).
	Order(builder.Order{Table: users, Column: "created_at", Desc: true}).
	Limit(20).
	Offset(40).
	Get()
```

## JOIN

```go
users := builder.NewTable("users")
posts := builder.NewTable("posts")

sql, args, err := builder.NewSelect().
	From(users).
	LeftJoin(posts, builder.OnEq{
		Table1:  users,
		Column1: "id",
		Table2:  posts,
		Column2: "user_id",
	}).
	Column(builder.ColumnName{Table: users, Name: "id"}).
	Column(builder.ColumnName{Table: posts, Name: "title"}).
	Get()
```

Joins are included only when the joined table is used by the query. For example, a joined table may be used in selected columns, `WHERE`, `GROUP BY`, or `ORDER BY`.

## WHERE

```go
users := builder.NewTable("users")

sql, args, err := builder.NewSelect().
	From(users).
	Column(builder.ColumnName{Table: users, Name: "id"}).
	Where(builder.WhereAnd{List: []builder.Where{
		builder.WhereEq{Table: users, Column: "active", Value: true},
		builder.WhereILike{Table: users, Column: "email", Value: "%@example.com"},
	}}).
	Get()
```

## INSERT

```go
users := builder.NewTable("users")

sql, args, err := builder.NewInsert(users).
	Value("email", "user@example.com").
	Value("name", "John").
	Return(builder.ColumnName{Table: users, Name: "id"}).
	Get()
```

## INSERT with ON CONFLICT

```go
users := builder.NewTable("users")

sql, args, err := builder.NewInsert(users).
	Value("email", "user@example.com").
	Value("name", "John").
	OnConflict("email").
	UpdateSet("name", "John").
	Return(builder.ColumnName{Table: users, Name: "id"}).
	Get()
```

## ON CONFLICT DO NOTHING

```go
users := builder.NewTable("users")

sql, args, err := builder.NewInsert(users).
	Value("email", "user@example.com").
	OnConflictDoNothing("email").
	Get()
```

## UPDATE

```go
users := builder.NewTable("users")

sql, args, err := builder.NewUpdate(users).
	Set("name", "John").
	SetNow("updated_at").
	Where(builder.WhereEq{Table: users, Column: "id", Value: 10}).
	Return(builder.ColumnName{Table: users, Name: "id"}).
	Get()
```

## DELETE

```go
users := builder.NewTable("users")

sql, args, err := builder.NewDelete(users).
	Where(builder.WhereEq{Table: users, Column: "id", Value: 10}).
	Get()
```

By default, `DELETE` without `WHERE` is rejected.

To intentionally generate a full-table delete, call `Full()`:

```go
sql, args, err := builder.NewDelete(users).
	Full().
	Get()
```

## Subqueries

```go
users := builder.NewTable("users")

sub := builder.NewSelect().
	From(users).
	Column(builder.ColumnName{Table: users, Name: "id"}).
	Where(builder.WhereEq{Table: users, Column: "active", Value: true})

activeUsers := builder.NewTableSub(sub)

sql, args, err := builder.NewSelect().
	From(activeUsers).
	Column(builder.ColumnName{Table: activeUsers, Name: "id"}).
	Get()
```

## EXISTS

```go
users := builder.NewTable("users")
posts := builder.NewTable("posts")

sub := builder.NewSelect().
	From(posts).
	Where(builder.WhereEqColumn{
		Table1:  users,
		Column1: "id",
		Table2:  posts,
		Column2: "user_id",
	})

sql, args, err := builder.NewSelect().
	From(users).
	Column(builder.ColumnName{Table: users, Name: "id"}).
	Where(builder.WhereExists{Query: sub}).
	Get()
```

`WhereExists` generates `SELECT 1` inside the `EXISTS` query.

## JSONB

Check that a JSONB value contains a key or string value:

```go
items := builder.NewTable("items")

sql, args, err := builder.NewSelect().
	From(items).
	Column(builder.ColumnName{Table: items, Name: "id"}).
	Where(builder.WhereJsonbTextExist{
		Table:  items,
		Column: "tags",
		Value:  "featured",
	}).
	Get()
```

Check that a JSONB value contains any value from a list:

```go
sql, args, err := builder.NewSelect().
	From(items).
	Column(builder.ColumnName{Table: items, Name: "id"}).
	Where(builder.WhereJsonbTextInExist{
		Table:  items,
		Column: "tags",
		Values: []string{"featured", "popular"},
	}).
	Get()
```

## Full-text search

```go
documents := builder.NewTable("documents")

sql, args, err := builder.NewSelect().
	From(documents).
	Column(builder.ColumnName{Table: documents, Name: "id"}).
	Where(builder.WhereFullText{
		Table:    documents,
		Column:   "body",
		Language: "english",
		Value:    "search text",
	}).
	Get()
```

The `Language` field is a PostgreSQL text search configuration name, such as `simple`, `english`, or a schema-qualified configuration like `public.english`.

## Query objects

Query builders are mutable.

Create a new query object for each query you want to build. Reusing the same query object and modifying it in multiple places can make code harder to reason about.

Calling `Get()` more than once is supported, but generated aliases and bind names are implementation details and should not be relied on.

## Limitations

`go-builder` intentionally does not try to cover every SQL feature.

Some advanced SQL expressions may require new builder types before they can be represented cleanly. The package currently prefers explicit typed helpers over raw SQL fragments.

Current limitations include:

* PostgreSQL-focused SQL generation
* `pgx.NamedArgs`-style placeholders
* no ORM-style model mapping
* no automatic schema introspection
* no SQL dialect abstraction
* no public raw SQL expression API

## Stability

The package is intentionally small and evolves based on real usage.

The API is kept simple, but changes are possible when they improve correctness, safety, or long-term maintainability.

## License

MIT
