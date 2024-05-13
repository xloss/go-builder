package builder

import (
	"fmt"
	"testing"
)

func TestInsertQuery_addBind(t *testing.T) {
	q := NewInsert(NewTable("table"))
	q.addBind("key1", 1)
	q.addBind("key2", 2.2)
	q.addBind("key3", "v")

	if len(q.binds) != 3 {
		t.Errorf("q.binds should have 3 values")
	}
	if q.binds["key1"] != 1 {
		t.Errorf("q.binds[\"key1\"] should have 1")
	}
	if q.binds["key2"] != 2.2 {
		t.Errorf("q.binds[\"key2\"] should have 2")
	}
	if q.binds["key3"] != "v" {
		t.Errorf("q.binds[\"key3\"] should have v")
	}
}

func TestInsertQuery_checkTable(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewInsert(table1)

	if !q.checkTable(table1) {
		t.Errorf("q.checkTable() returned false")
	}

	if q.checkTable(table2) {
		t.Errorf("q.checkTable() returned true")
	}
}

func TestInsertQuery_Value(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.Value("col1", 5)
	q.Value("col2", 2.2)

	if len(q.values) != 2 {
		t.Errorf("q.values should have 2 values")
	}
	if q.values[0].Column != "col1" {
		t.Errorf("q.values[0].Column should have col1")
	}
	if q.values[0].Value != 5 {
		t.Errorf("q.values[0].Value should have 5")
	}
	if q.values[1].Column != "col2" {
		t.Errorf("q.values[1].Column should have col2")
	}
	if q.values[1].Value != 2.2 {
		t.Errorf("q.values[1].Value should have 2.2")
	}
}

func TestInsertQuery_Column(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.Column(table, "col1")
	q.ColumnAlias(table, "col2", "a1")

	if len(q.columns) != 2 {
		t.Errorf("q.values should have 2 values")
	}
	if q.columns[0].Table != table {
		t.Errorf("q.columns[0].Table should have table")
	}
	if q.columns[0].Name != "col1" {
		t.Errorf("q.columns[0].Name should have col1")
	}
	if q.columns[1].Table != table {
		t.Errorf("q.columns[1].Table should have table")
	}
	if q.columns[1].Name != "col2" {
		t.Errorf("q.columns[1].Name should have col2")
	}
	if q.columns[1].Alias != "a1" {
		t.Errorf("q.columns[1].Alias should have a1")
	}
}

func TestInsertQuery_ColumnAlias(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.ColumnAlias(table, "name", "desc")

	if len(q.columns) != 1 {
		t.Errorf("q.columns should have 1 column")
	}

	if q.columns[0].Name != "name" {
		t.Errorf("q.columns[0].Name should have name string")
	}
	if q.columns[0].Table != table {
		t.Errorf("q.columns[0].Table should have table %v", table)
	}
	if q.columns[0].Alias != "desc" {
		t.Errorf("q.columns[0].Alias should have desc string")
	}
	if q.columns[0].Aggregate != false {
		t.Errorf("q.columns[0].Aggregate should have false")
	}
}

func TestInsertQuery_getValues(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.Value("col1", 5)
	q.Value("col2", 2.2)

	values, err := q.getValues()
	if err != nil {
		t.Errorf("q.getValues() returned %v", err)
	}

	var tag1, tag2 string

	for k, v := range q.binds {
		if v == 5 {
			tag1 = k
		} else if v == 2.2 {
			tag2 = k
		}
	}

	if values != " (col1, col2) VALUES (@"+tag1+", @"+tag2+")" {
		t.Errorf("q.getValues() returned %v", values)
	}
}

func TestInsertQuery_getReturns(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.Column(table, "col1")
	q.ColumnAlias(table, "col2", "a1")

	returns := q.getReturns()

	if returns != " RETURNING "+table.Alias+".col1, "+table.Alias+".col2 as a1" {
		t.Errorf("q.getReturns() returned '%v'", returns)
	}
}

func TestInsertQuery_Get(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.Value("col1", 5)
	q.Value("col2", "str")

	q.Column(table, "col1")
	q.ColumnAlias(table, "col2", "a1")

	sql, binds, err := q.Get()
	if err != nil {
		t.Errorf("q.Get() returned %v", err)
	}

	var tag1, tag2 string

	for k, v := range binds {
		if v == 5 {
			tag1 = k
		} else if v == "str" {
			tag2 = k
		}
	}

	if len(binds) != 2 {
		t.Errorf("q.Get() should have 2 values")
	}

	if sql != fmt.Sprintf("INSERT INTO %[1]s AS %[2]s (col1, col2) VALUES (@%[3]s, @%[4]s) RETURNING %[2]s.col1, %[2]s.col2 as a1", table.Name, table.Alias, tag1, tag2) {
		t.Errorf("q.Get() returned '%v'", sql)
	}
}
