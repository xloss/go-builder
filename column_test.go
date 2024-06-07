package builder

import (
	"fmt"
	"testing"
)

func TestColumnName_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c1 := ColumnName{Table: table, Name: "col1"}
	c2 := ColumnName{Table: table, Name: "col2", Alias: "a1"}
	c3 := ColumnName{Table: table, Name: "col3", Distinct: true}

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

	s3, err := c3.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if s3 != "DISTINCT "+table.Alias+".col3" {
		t.Fatal(s3)
	}
}

func TestColumnCount_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c0 := ColumnCount{}
	c1 := ColumnCount{Alias: "col1"}
	c2 := ColumnCount{Table: table, Name: "col2", Alias: "a1"}
	c3 := ColumnCount{Table: table, Name: "col3", Alias: "a2", Distinct: true}

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

	s3, err := c3.gen(q)
	if err != nil {
		t.Fatal(err)
	}
	if s3 != "COUNT(DISTINCT "+table.Alias+".col3) AS a2" {
		t.Fatal(s3)
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

func TestColumnJsonbArrayElementsText_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect().From(table)
	q.From(table)

	c := ColumnJsonbArrayElementsText{Name: "col", Alias: "a"}

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

	c.Name = "col"
	c.Alias = ""

	_, err = c.gen(q)
	if err == nil {
		t.Error("expected error")
	}

	c.Alias = "a1"

	sql, err := c.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	if sql != "JSONB_ARRAY_ELEMENTS_TEXT("+table.Alias+"."+c.Name+") AS "+c.Alias {
		t.Fatal(sql)
	}

	c.Distinct = true

	sql, err = c.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	if sql != "DISTINCT JSONB_ARRAY_ELEMENTS_TEXT("+table.Alias+"."+c.Name+") AS "+c.Alias {
		t.Fatal(sql)
	}
}

func ExampleColumnName() {
	table1 := NewTable("table1")
	query1 := NewSelect()
	query1.Column(
		ColumnName{Table: table1, Name: "column1"},
		ColumnName{Table: table1, Name: "column2", Alias: "a2"},
		ColumnName{Table: table1, Name: "column3", Distinct: true},
	)
	query1.From(table1)

	fmt.Println(query1.Get())

	// Result:
	// SELECT table1_punanojozl.column1, table1_punanojozl.column2 AS a2, DISTINCT table1_punanojozl.column3 FROM table1 AS table1_punanojozl
	// map[]
	// <nil>
}

func ExampleColumnCount() {
	table1 := NewTable("table1")
	query1 := NewSelect()
	query1.Column(
		ColumnCount{Table: table1, Alias: "a1"},
		ColumnCount{Table: table1, Name: "column2", Alias: "a2"},
		ColumnCount{Table: table1, Name: "column3", Alias: "a3", Distinct: true},
	)
	query1.From(table1)

	fmt.Println(query1.Get())

	// Result:
	// SELECT COUNT(*) AS a1, COUNT(table1_yyapxlsrva.column2) AS a2, COUNT(DISTINCT table1_yyapxlsrva.column3) AS a3 FROM table1 AS table1_yyapxlsrva
	// map[]
	// <nil>
}

func ExampleColumnCoalesce() {
	table1 := NewTable("table1")
	query1 := NewSelect()
	query1.Column(
		ColumnCoalesce{Table: table1, Name: "column1", Alias: "a1", Default: 5},
		ColumnCoalesce{Table: table1, Name: "column2", Alias: "a2", Default: "text"},
	)
	query1.From(table1)

	fmt.Println(query1.Get())

	// Result:
	// SELECT COALESCE(table1_yzethlflca.column1, 5) AS a1, COALESCE(table1_yzethlflca.column2, 'text') AS a2 FROM table1 AS table1_yzethlflca
	// map[]
	// <nil>
}

func ExampleColumnJsonbArrayElementsText() {
	table1 := NewTable("table1")
	query1 := NewSelect()
	query1.Column(
		ColumnJsonbArrayElementsText{Table: table1, Name: "column1", Alias: "a1"},
		ColumnJsonbArrayElementsText{Table: table1, Name: "column2", Alias: "a2", Distinct: true},
	)
	query1.From(table1)

	fmt.Println(query1.Get())

	// Result:
	// SELECT JSONB_ARRAY_ELEMENTS_TEXT(table1_yifmeamxwi.column1) AS a1, DISTINCT JSONB_ARRAY_ELEMENTS_TEXT(table1_yifmeamxwi.column2) AS a2 FROM table1 AS table1_yifmeamxwi
	// map[]
	// <nil>
}
