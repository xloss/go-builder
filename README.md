# PostgreSQL Query Builder for "jackc/pgx"
## Install
Use `go get` to install this package.
```
go get github.com/xloss/go-builder
```

Refer to the documentation here:
http://godoc.org/github.com/xloss/go-builder

## Basic usage
### Select
```go
table1 := builder.NewTable("table1")
table2 := builder.NewTable("table2")

q := builder.NewSelect()
q.From(table1)
q.Column(builder.ColumnName{Table: table1, Name: "id"})
q.Column(builder.ColumnName{Table: table2, Name: "col", Alias: "c1"})
q.LeftJoin(table2, builder.OnEq{Table1: table1, Table2: table2, Column1: "id", Column2: "table_id"})
q.Where(builder.WhereEq{Table: table1, Column: "id", Value: 1})
q.Order(builder.Order{Table: table1, Column: "name", Desc: true})
q.Limit(10)
q.Offset(5)

sql, binds, err := q.Get()
fmt.Println(sql)
fmt.Println(binds)

// Output:
// SELECT table1_pxaisxqdvb.id, table2_xftynvknii.col AS c1 FROM table1 AS table1_pxaisxqdvb LEFT JOIN table2 AS table2_xftynvknii ON table1_pxaisxqdvb.id = table2_xftynvknii.table_id WHERE table1_pxaisxqdvb.id = @id_ywfaoazuel ORDER BY table1_pxaisxqdvb.name DESC LIMIT @limit_jlhldzwuei OFFSET @offset_myzuqcgmdn
// map[id_ywfaoazuel:1 limit_jlhldzwuei:10 offset_myzuqcgmdn:5]
```
Subquery in FROM:
```go
table1 := builder.NewTable("table1")
q1 := builder.NewSelect()
q1.From(table1)
q1.Column(builder.ColumnName{Table: table1, Name: "column1"})

table2 := NewTableSub(q1)

q2 := builder.NewSelect()
q2.From(table2)
q2.Column(builder.ColumnName{Table: table2, Name: "column1"})

sql, binds, err := q2.Get()
fmt.Println(sql)

// Output:
// SELECT wjsfhotgjo_vhncsrjhtg.column1 FROM (SELECT table1_ttwlctmwvj.column1 FROM table1 AS table1_ttwlctmwvj) AS wjsfhotgjo_vhncsrjhtg
```
### Update
```go
table := builder.NewTable("table")

q := builder.NewUpdate(table)
q.Set("col1", "value1")
q.SetNow("col2")
q.Where(builder.WhereEq{Table: table, Column: "col3", Value: 5})

q.Return(builder.ColumnName{Table: table, Name: "col1"})
q.Return(builder.ColumnName{Table: table, Name: "col2", Alias: "a1"})

sql, binds, err := q.Get()
fmt.Println(sql)
fmt.Println(binds)

// Output:
// UPDATE table AS table_kiykrrnhxf SET col1 = @col1_tolhdmocsn, col2 = NOW() WHERE table_kiykrrnhxf.col3 = @col3_tkdyhzjxqb RETURNING table_kiykrrnhxf.col1, table_kiykrrnhxf.col2 AS a1
// map[col1_tolhdmocsn:value1 col3_tkdyhzjxqb:5]
```
### Insert
```go
table := builder.NewTable("table")
q := builder.NewInsert(table)

q.Value("col1", 5)
q.Value("col2", "str")

q.Return(builder.ColumnName{Table: table, Name: "col1"})
q.Return(builder.ColumnName{Table: table, Name: "col2", Alias: "a1"})

sql, binds, err := q.Get()

// Output:
// INSERT INTO table AS table_jhpjqkzvkd (col1, col2) VALUES (@col1_buaonudjkx, @col2_amouztvkkt) RETURNING table_jhpjqkzvkd.col1, table_jhpjqkzvkd.col2 AS a1
// map[col1_buaonudjkx:5 col2_amouztvkkt:str]
```
### Delete
```go
table := builder.NewTable("table")
q := builder.NewDelete(table)

q.Where(builder.WhereEq{Table: table, Column: "col", Value: 5})

sql, binds, err := q.Get()
// Output:
// DELETE FROM table AS table_iuulmrhwnt WHERE table_iuulmrhwnt.col = @col_ujpmtrhymk
// map[col_ujpmtrhymk:5]
```

Without WHERE:
```go
table := builder.NewTable("table")
q := builder.NewDelete(table)

q.Full()

sql, binds, err := q.Get()
// Output:
// DELETE FROM table AS table_iuulmrhwnt
// map[]
```