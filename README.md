# PostreSQL Query Builder for "jackc/pgx"
## Install
Use `go get` to install this package.
```
go get github.com/xloss/go-builder
```

## Usage
### Basic usage
```go
table1 := builder.NewTable("table1")
table2 := builder.NewTable("table2")

q := builder.NewSelect()
q.
    From(table1).
    Column(table1, "id").
    Column(table2, "col").
    LeftJoin(table2, builder.OnEq{Table1: table1, Table2: table2, Column1: "id", Column2: "table_id"}).
    Where(builder.WhereEq{Table: table1, Column: "id", Value: 1}).
    Order(builder.Order{Table: table1, Column: "name", Desc: true}).
    Limit(10).
    Offset(5)

sql, binds, err := q.Get()
fmt.Println(sql)
fmt.Println(binds)

// Output:
// SELECT table1_pxaisxqdvb.id, table2_xftynvknii.col FROM table1 as table1_pxaisxqdvb LEFT JOIN table2 AS table2_xftynvknii ON table1_pxaisxqdvb.id = table2_xftynvknii.table_id WHERE table1_pxaisxqdvb.id = @id_ywfaoazuel ORDER BY table1_pxaisxqdvb.name DESC LIMIT @limit_jlhldzwuei OFFSET @offset_myzuqcgmdn
// map[id_ywfaoazuel:1 limit_jlhldzwuei:10 offset_myzuqcgmdn:5]
```

## Package in Development
### Available
#### SELECT
##### Init
```go
q := builder.NewSelect()
```

##### Columns
* Column (`.Column(table *Table, name string)` `SELECT table_hash.col`)
* Column with Alias (`.ColumnAlias(table *Table, name string, alias string)` `SELECT table_hash.col as alias`)
* COUNT function (`.ColumnCount(table *Table, alias string)` `SELECT COUNT(*) as alias`)
* COALESCE Function (`.ColumnCoalesce(table *Table, name string, alias string, default any)` `SELECT COALESCE(table_hash.col, 10) as alias`)

##### From
```go
q.From(table1)
q.From(table2)
```
```go
q.From(table1, table2)
```

##### Join
* Left Join (`.LeftJoin(table *Table, on On)`)

###### On
* `builder.OnEq{Table1 *Table, Table2 *Table, Column1 string, Column2 string}`
* `builder.OnLess{Table1 *Table, Table2 *Table, Column1 string, Column2 string}`
* `builder.OnMore{Table1 *Table, Table2 *Table, Column1 string, Column2 string}`
* `builder.OnAnd{List: []builder.On{}}`

##### Where
```go
.Where(where builder.Where)
```

* `builder.WhereEq{Table *Table, Column string, Value any}` `table_hash.col = @value_hash`
* `builder.WhereIsNull{Table *Table, Column string}` `table_hash.col IS NULL`
* `builder.WhereIsNotNull{Table *Table, Column string}` `table_hash.col IS NOT NULL`
* `builder.WhereIn{Table *Table, Column string, Values interface{}}` `table_hash.col = ANY(@values_hash)`
* `builder.WhereMore{Table *Table, Column string, Value any}` `table_hash.col > @value_hash`
* `builder.WhereLess{Table *Table, Column string, Value any}`
* `builder.WhereMoreEq{Table *Table, Column string, Value any}` `table_hash.col >= @value_hash`
* `builder.WhereLessEq{Table *Table, Column string, Value any}`
* `builder.WhereMoreColumn{Table1 *Table, Table2 *Table, Column1 string, Column2 string}` `table1_hash.col1 = table2_hash.col2`
* `builder.WhereFullText{Table *Table, Column string, Value string, Language string}` `to_tsvector('language', table_hash.col) @@ plainto_tsquery(@value_hash)`
* `builder.WhereAnd{List: []builder.Where{}}`
* `builder.WhereOr{List: []builder.Where{}}`

##### Order
`.Order(order ...builder.Order)`
```go
builder.Order{Table *Table, Column string, Desc bool}
```

##### Limit
`.Limit(limit int)`

##### OFFSET
`.Offset(offset int)`