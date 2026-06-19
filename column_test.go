package builder

import (
	"fmt"
	"strings"
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

	var tag string
	for k, v := range q.binds {
		if v == 10 {
			tag = k
		}
	}

	if tag == "" {
		t.Fatal("tag should not be empty")
	}

	if sql != "COALESCE("+table.Alias+".col1, @"+tag+") AS a1" {
		t.Fatal(sql)
	}

	c.Default = "O'Reilly"
	sql, err = c.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	tag = ""
	for k, v := range q.binds {
		if v == "O'Reilly" {
			tag = k
		}
	}

	if tag == "" {
		t.Fatal("tag should not be empty")
	}

	if sql != "COALESCE("+table.Alias+".col1, @"+tag+") AS a1" {
		t.Fatal(sql)
	}
}

func TestColumnCount_genFilter(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c1 := ColumnCount{
		Alias: "a1",
		Filter: WhereIsNull{
			Table:  table,
			Column: "col1",
		},
	}

	s1, err := c1.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	if s1 != "COUNT(*) FILTER (WHERE "+table.Alias+".col1 IS NULL) AS a1" {
		t.Fatal(s1)
	}

	c2 := ColumnCount{
		Alias: "a2",
		Filter: WhereEq{
			Table:  table,
			Column: "col2",
			Value:  1,
		},
	}

	s2, err := c2.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	var tag string

	for k, v := range q.binds {
		if v == 1 {
			tag = k
		}
	}

	if tag == "" {
		t.Fatal("tag should not be empty")
	}

	if s2 != "COUNT(*) FILTER (WHERE "+table.Alias+".col2 = @"+tag+") AS a2" {
		t.Fatal(s2)
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

func TestColumnValue_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c1 := ColumnValue{Value: 1}
	c2 := ColumnValue{Value: "O'Reilly", Alias: "a1"}
	c3 := ColumnValue{}

	s1, err := c1.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	var tag1 string
	for k, v := range q.binds {
		if v == 1 {
			tag1 = k
		}
	}

	if tag1 == "" {
		t.Fatal("tag1 should not be empty")
	}

	if s1 != "@"+tag1 {
		t.Fatal(s1)
	}

	s2, err := c2.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	var tag2 string
	for k, v := range q.binds {
		if v == "O'Reilly" {
			tag2 = k
		}
	}

	if tag2 == "" {
		t.Fatal("tag2 should not be empty")
	}

	if s2 != "@"+tag2+" AS a1" {
		t.Fatal(s2)
	}

	_, err = c3.gen(q)
	if err == nil {
		t.Error("expected error")
	}
}

func TestColumnValue_genDoesNotInlineString(t *testing.T) {
	table := NewTable("table")

	q := NewSelect().
		From(table).
		Column(ColumnValue{Value: "O'Reilly", Alias: "value"})

	sql, binds, err := q.Get()
	if err != nil {
		t.Fatal(err)
	}

	if len(binds) != 1 {
		t.Fatalf("binds should have 1 value")
	}

	if strings.Contains(sql, "O'Reilly") {
		t.Fatalf("sql should not contain raw string value: %s", sql)
	}

	for _, v := range binds {
		if v != "O'Reilly" {
			t.Fatalf("bind value should have O'Reilly")
		}
	}
}

func TestColumnName_genInvalidName(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnName{Table: table, Name: "bad name"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnName_genInvalidAlias(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnName{Table: table, Name: "col", Alias: "bad alias"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnCount_genInvalidName(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnCount{Table: table, Name: "bad name", Alias: "a1"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnCount_genInvalidAlias(t *testing.T) {
	q := NewSelect()

	c := ColumnCount{Alias: "bad alias"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnCoalesce_genInvalidName(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnCoalesce{Table: table, Name: "bad name", Alias: "a1", Default: 1}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnCoalesce_genInvalidAlias(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnCoalesce{Table: table, Name: "col", Alias: "bad alias", Default: 1}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnJsonbArrayElementsText_genInvalidName(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnJsonbArrayElementsText{Table: table, Name: "bad name", Alias: "a1"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnJsonbArrayElementsText_genInvalidAlias(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnJsonbArrayElementsText{Table: table, Name: "col", Alias: "bad alias"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnValue_genInvalidAlias(t *testing.T) {
	q := NewSelect()

	c := ColumnValue{Value: 1, Alias: "bad alias"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnMin_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c1 := ColumnMin{Table: table, Name: "col1", Alias: "a1"}

	s1, err := c1.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	if s1 != "MIN("+table.Alias+".col1) AS a1" {
		t.Fatal(s1)
	}

	c2 := ColumnMin{
		Table: table,
		Name:  "col2",
		Alias: "a2",
		Filter: WhereIsNull{
			Table:  table,
			Column: "col3",
		},
	}

	s2, err := c2.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	if s2 != "MIN("+table.Alias+".col2) FILTER (WHERE "+table.Alias+".col3 IS NULL) AS a2" {
		t.Fatal(s2)
	}
}

func TestColumnMax_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c1 := ColumnMax{Table: table, Name: "col1", Alias: "a1"}

	s1, err := c1.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	if s1 != "MAX("+table.Alias+".col1) AS a1" {
		t.Fatal(s1)
	}

	c2 := ColumnMax{
		Table: table,
		Name:  "col2",
		Alias: "a2",
		Filter: WhereIsNotNull{
			Table:  table,
			Column: "col3",
		},
	}

	s2, err := c2.gen(q)
	if err != nil {
		t.Fatal(err)
	}

	if s2 != "MAX("+table.Alias+".col2) FILTER (WHERE "+table.Alias+".col3 IS NOT NULL) AS a2" {
		t.Fatal(s2)
	}
}

func TestColumnMin_genInvalidName(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnMin{Table: table, Name: "bad name", Alias: "a1"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnMin_genInvalidAlias(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnMin{Table: table, Name: "col", Alias: "bad alias"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnMax_genInvalidName(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnMax{Table: table, Name: "bad name", Alias: "a1"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
	}
}

func TestColumnMax_genInvalidAlias(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	c := ColumnMax{Table: table, Name: "col", Alias: "bad alias"}

	_, err := c.gen(q)
	if err == nil {
		t.Errorf("c.gen should have returned error")
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
	// SELECT COALESCE(table1_epzedlhnbl.column1, @column1_default_gznyaknpce) AS a1, COALESCE(table1_epzedlhnbl.column2, @column2_default_odlytjcxsz) AS a2 FROM table1 AS table1_epzedlhnbl
	// map[column1_default_gznyaknpce:5 column2_default_odlytjcxsz:text]
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

func ExampleColumnValue() {
	table1 := NewTable("table1")
	query1 := NewSelect()
	query1.Column(
		ColumnValue{Value: 1},
		ColumnValue{Value: "1", Alias: "a1"},
	)
	query1.From(table1)

	fmt.Println(query1.Get())

	// Result:
	// SELECT @value_ufglfjoels, @value_tbotyrxlvp AS a1 FROM table1 AS table1_luroihnrel
	// map[value_tbotyrxlvp:1 value_ufglfjoels:1]
	// <nil>
}
