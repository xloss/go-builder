package builder

import "testing"

func TestColumnName_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c1 := ColumnName{Table: table, Name: "col1"}
	c2 := ColumnName{Table: table, Name: "col2", Alias: "a1"}

	s1, err := c1.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if s1 != table.Alias+".col1" {
		t.Fatal(s1)
	}

	s2, err := c2.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if s2 != table.Alias+".col2 AS a1" {
		t.Fatal(s2)
	}
}

func TestColumnCount_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c0 := ColumnCount{}
	c1 := ColumnCount{Alias: "col1"}
	c2 := ColumnCount{Table: table, Name: "col2", Alias: "a1"}

	_, err := c0.gen(q)
	if err == nil {
		t.Error("expected error")
	}

	s1, err := c1.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if s1 != "COUNT(*) AS col1" {
		t.Fatal(s1)
	}

	s2, err := c2.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if s2 != "COUNT("+table.Alias+".col2) AS a1" {
		t.Fatal(s2)
	}
}

func TestColumnCoalesce_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnCoalesce{Name: "col1", Alias: "a1", Default: "10"}

	_, err := c.gen(q)
	if err == nil {
		t.Error("expected error")
	}

	c.Table = table
	c.Name = ""

	_, err = c.gen(q)
	if err == nil {
		t.Error("expected error")
	}

	c.Name = "col1"
	c.Alias = ""

	_, err = c.gen(q)
	if err == nil {
		t.Error("expected error")
	}

	c.Alias = "a1"
	c.Default = nil

	_, err = c.gen(q)
	if err == nil {
		t.Error("expected error")
	}

	c.Default = 10
	sql, err := c.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != "COALESCE("+table.Alias+".col1, 10) AS a1" {
		t.Fatal(sql)
	}

	c.Default = "str"
	sql, err = c.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql != "COALESCE("+table.Alias+".col1, 'str') AS a1" {
		t.Fatal(sql)
	}
}
