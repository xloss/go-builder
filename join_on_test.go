package builder

import (
	"fmt"
	"testing"
)

func TestOnEq_gen(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)
	q.From(table2)

	on := OnEq{
		Table1:  table1,
		Table2:  table2,
		Column1: "id",
		Column2: "table_id",
	}

	sql, err := on.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != fmt.Sprintf("%s.%s = %s.%s", table1.Alias, on.Column1, table2.Alias, on.Column2) {
		t.Fatal(sql)
	}
}

func TestOnMore_gen(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)
	q.From(table2)

	on := OnMore{
		Table1:  table1,
		Table2:  table2,
		Column1: "id",
		Column2: "table_id",
	}

	sql, err := on.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != fmt.Sprintf("%s.%s > %s.%s", table1.Alias, on.Column1, table2.Alias, on.Column2) {
		t.Fatal(sql)
	}
}

func TestOnLess_gen(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)
	q.From(table2)

	on := OnLess{
		Table1:  table1,
		Table2:  table2,
		Column1: "id",
		Column2: "table_id",
	}

	sql, err := on.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != fmt.Sprintf("%s.%s < %s.%s", table1.Alias, on.Column1, table2.Alias, on.Column2) {
		t.Fatal(sql)
	}
}

func TestOnAnd_gen(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)
	q.From(table2)

	on := OnAnd{
		List: []On{
			OnEq{
				Table1:  table1,
				Table2:  table2,
				Column1: "id",
				Column2: "table_id",
			},
			OnLess{
				Table1:  table1,
				Table2:  table2,
				Column1: "len",
				Column2: "len",
			},
		},
	}

	sql, err := on.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != fmt.Sprintf("(%[1]s.id = %[2]s.table_id AND %[1]s.len < %[2]s.len)", table1.Alias, table2.Alias) {
		t.Fatal(sql)
	}

}
