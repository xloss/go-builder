package builder

import (
	"fmt"
	"testing"
)

func TestJoin_Gen(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)
	q.From(table2)

	j := join{
		Table: table1,
		On: OnEq{
			Table1:  table1,
			Table2:  table2,
			Column1: "id",
			Column2: "table_id",
		},
	}

	sql, err := j.Gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != fmt.Sprintf(" JOIN %[1]s AS %[2]s ON %[2]s.id = %[3]s.table_id", table1.Name, table1.Alias, table2.Alias) {
		t.Fatal(sql)
	}

	j.Left = true

	sql, err = j.Gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != fmt.Sprintf(" LEFT JOIN %[1]s AS %[2]s ON %[2]s.id = %[3]s.table_id", table1.Name, table1.Alias, table2.Alias) {
		t.Fatal(sql)
	}

	j.Left = false
	j.Inner = true

	sql, err = j.Gen(q)
	if err != nil {
		t.Fatal(err)
	}

	if sql != fmt.Sprintf(" INNER JOIN %[1]s AS %[2]s ON %[2]s.id = %[3]s.table_id", table1.Name, table1.Alias, table2.Alias) {
		t.Fatal(sql)
	}
}

func TestJoin_GenInvalidTableName(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("bad table")

	q := NewSelect()
	q.From(table1)
	q.LeftJoin(table2, OnEq{
		Table1:  table1,
		Column1: "id",
		Table2:  table2,
		Column2: "table_id",
	})
	q.Column(ColumnName{Table: table2, Name: "id"})

	_, err := q.joins[0].Gen(q)
	if err == nil {
		t.Errorf("j.Gen should have returned error")
	}
}

func TestJoin_GenInvalidTableAlias(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")
	table2.Alias = "bad alias"

	q := NewSelect()
	q.From(table1)
	q.LeftJoin(table2, OnEq{
		Table1:  table1,
		Column1: "id",
		Table2:  table2,
		Column2: "table_id",
	})
	q.Column(ColumnName{Table: table2, Name: "id"})

	_, err := q.joins[0].Gen(q)
	if err == nil {
		t.Errorf("j.Gen should have returned error")
	}
}

func TestJoin_GenSubqueryTable(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTableSub(NewSelect())

	q := NewSelect()
	q.From(table1)
	q.LeftJoin(table2, OnEq{
		Table1:  table1,
		Column1: "id",
		Table2:  table2,
		Column2: "table_id",
	})
	q.Column(ColumnName{Table: table2, Name: "id"})

	_, err := q.joins[0].Gen(q)
	if err == nil {
		t.Errorf("j.Gen should have returned error")
	}
}
