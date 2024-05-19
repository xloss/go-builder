package builder

import "testing"

func TestGroupColumn_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	group := GroupColumn{
		Table:  table,
		Column: "column1",
	}

	sql, err := group.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != table.Alias+"."+group.Column {
		t.Fatal(sql)
	}

	group = GroupColumn{
		Column: "column2",
	}
	sql, err = group.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != group.Column {
		t.Fatal(sql)
	}
}
