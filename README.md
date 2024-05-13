# PostgreSQL Query Builder for "jackc/pgx"
## Install
Use `go get` to install this package.
```
go get github.com/xloss/go-builder
```

## Usage
### Basic usage
#### Select
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
#### Update
```go
table := builder.NewTable("table")

q := builder.NewUpdate(table)
q.Set("col1", "value1")
q.SetNow("col2")
q.Where(builder.WhereEq{Table: table, Column: "col3", Value: 5})

sql, binds, err := q.Get()
fmt.Println(sql)
fmt.Println(binds)

// Output:
// UPDATE table AS table_kiykrrnhxf SET col1 = @col1_tolhdmocsn, col2 = NOW() WHERE table_kiykrrnhxf.col3 = @col3_tkdyhzjxqb
// map[col1_tolhdmocsn:value1 col3_tkdyhzjxqb:5]
```
#### Insert
```go
table := NewTable("table")
q := NewInsert(table)

q.Value("col1", 5)
q.Value("col2", "str")

q.Column(table, "col1")
q.ColumnAlias(table, "col2", "a1")

sql, binds, err := q.Get()

// Output:
// INSERT INTO table AS table_jhpjqkzvkd (col1, col2) VALUES (@col1_buaonudjkx, @col2_amouztvkkt) RETURNING table_jhpjqkzvkd.col1, table_jhpjqkzvkd.col2 as a1
// map[col1_buaonudjkx:5 col2_amouztvkkt:str]
```
#### Delete
```go
table := NewTable("table")
q := NewDelete(table)

q.Where(WhereEq{Table: table, Column: "col", Value: 5})

sql, binds, err := q.Get()
// Output:
// DELETE FROM table AS table_iuulmrhwnt WHERE table_iuulmrhwnt.col = @col_ujpmtrhymk
// map[col_ujpmtrhymk:5]

// Without WHERE
table := NewTable("table")
q := NewDelete(table)

q.Full()

sql, binds, err := q.Get()
// Output:
// DELETE FROM table AS table_iuulmrhwnt
// map[]
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

* `builder.WhereEq{Table *Table, Column string, Value interface{}}` - `table_hash.col = @value_hash`
* `builder.WhereIsNull{Table *Table, Column string}` - `table_hash.col IS NULL`
* `builder.WhereIsNotNull{Table *Table, Column string}` - `table_hash.col IS NOT NULL`
* `builder.WhereIn{Table *Table, Column string, Values interface{}}` - `table_hash.col = ANY(@values_hash)`
* `builder.WhereMore{Table *Table, Column string, Value interface{}}` - `table_hash.col > @value_hash`
* `builder.WhereLess{Table *Table, Column string, Value interface{}}`
* `builder.WhereMoreEq{Table *Table, Column string, Value interface{}}` - `table_hash.col >= @value_hash`
* `builder.WhereLessEq{Table *Table, Column string, Value interface{}}`
* `builder.WhereMoreColumn{Table1 *Table, Table2 *Table, Column1 string, Column2 string}` - `table1_hash.col1 = table2_hash.col2`
* `builder.WhereILike{Table *Table, Column string, Value interface{}}` - `table_hash.col ILIKE @value_hash`
* `builder.WhereFullText{Table *Table, Column string, Value string, Language string}` - `to_tsvector('language', table_hash.col) @@ plainto_tsquery(@value_hash)`
* `builder.WhereAnd{List: []builder.Where{}}`
* `builder.WhereOr{List: []builder.Where{}}`

##### Order
`.Order(order ...builder.Order)`
```go
builder.Order{Table *Table, Column string, Desc bool}
```

##### Limit
`.Limit(limit int)`

##### Offset
`.Offset(offset int)`

#### UPDATE
##### Init
```go
t := builder.NewTable("table")
q := builder.NewUpdate(t)
```

##### Set
* `.Set(column string, value any)` `table_hash.col = @value_hash`
* `.SetNow(column string)` `table_hash.col = NOW()`

##### Where
[See](#where)

#### INSERT
##### Init
```go
t := builder.NewTable("table")
q := builder.NewInsert(t)
```
##### Values
`.Value(name string, value any)`

##### Returning
* `.Column(table *Table, name string)` - `RETURNING table_hash.col`
* `.ColumnAlias(table *Table, name string, alias string)` - `RETURNING table_hash.col as alias`

#### DELETE
##### Init
```go
t := builder.NewTable("table")
q := builder.NewDelete(t)
```
##### Where
[See](#where)

##### Protection against complete cleaning
Without a WHERE block, function `.Get()` will throw an error. You need to call function `.Full()` to generate sql without a WHERE block.