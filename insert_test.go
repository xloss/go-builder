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

func TestInsertQuery_Return(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.Return(ColumnName{Table: table, Name: "col1"})
	q.Return(ColumnName{Table: table, Name: "col2", Alias: "a1"})

	if len(q.returns) != 2 {
		t.Errorf("q.values should have 2 values")
	}
}

func TestInsertQuery_OnConflict(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.OnConflict("col1", "col2")

	if len(q.conflict) != 2 {
		t.Errorf("q.conflict should have 2 values")
	}
}

func TestInsertQuery_UpdateSet(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.UpdateSet("col1", 1)
	q.UpdateSet("col2", "a")

	if len(q.update) != 2 {
		t.Errorf("q.update should have 2 values")
	}
	if q.update[0].Column != "col1" {
		t.Errorf("q.update[0].Column should have col1")
	}
	if q.update[0].Value != 1 {
		t.Errorf("q.update[0].Value should have 1")
	}
	if q.update[1].Column != "col2" {
		t.Errorf("q.update[1].Column should have col2")
	}
	if q.update[1].Value != "a" {
		t.Errorf("q.update[1].Value should have a")
	}
}

func TestInsertQuery_UpdateSetNow(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.UpdateSetNow("col1")

	if len(q.update) != 1 {
		t.Errorf("q.update should have 2 values")
	}
	if q.update[0].Column != "col1" {
		t.Errorf("q.update[0].Column should have col1")
	}
	if q.update[0].Value != nil {
		t.Errorf("q.update[0].Value should have 1")
	}
	if !q.update[0].Now {
		t.Errorf("q.update[0].Now should have true")
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

func TestInsertQuery_getConflict(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	sql := q.getConflict()
	if sql != "" {
		t.Errorf("q.getConflict() returned %v", sql)
	}

	q.OnConflict("col1", "col2")

	sql = q.getConflict()
	if sql != " ON CONFLICT (col1, col2)" {
		t.Errorf("q.getConflict() returned %v", sql)
	}
}

func TestInsertQuery_getUpdate(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	sql := q.getUpdate()
	if sql != "" {
		t.Errorf("q.getUpdate() returned %v", sql)
	}

	q.OnConflict("col1")
	sql = q.getUpdate()
	if sql != "" {
		t.Errorf("q.getUpdate() returned %v", sql)
	}

	q.UpdateSet("col1", "value1")
	q.UpdateSetNow("col2")

	sql = q.getUpdate()

	var tag string

	for k, v := range q.binds {
		if v == "value1" {
			tag = k
		}
	}

	if sql != " DO UPDATE SET col1 = @"+tag+", col2 = NOW()" {
		t.Errorf("q.getSet() returned '%v'", sql)
	}
}

func TestInsertQuery_getReturns(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	returns, err := q.getReturns()
	if err != nil {
		t.Errorf("q.getReturns() returned %v", err)
	}
	if returns != "" {
		t.Errorf("q.getReturns() returned %v", returns)
	}

	q.Return(ColumnName{Table: table, Name: "col1"})
	q.Return(ColumnName{Table: table, Name: "col2", Alias: "a1"})

	returns, err = q.getReturns()
	if err != nil {
		t.Errorf("q.getReturns() returned %v", err)
	}
	if returns != " RETURNING "+table.Alias+".col1, "+table.Alias+".col2 AS a1" {
		t.Errorf("q.getReturns() returned '%v'", returns)
	}
}

func TestInsertQuery_Get(t *testing.T) {
	table := NewTable("table")
	q := NewInsert(table)

	q.Value("col1", 5)
	q.Value("col2", "str")

	q.OnConflict("col1")
	q.UpdateSet("col2", "val")

	q.Return(ColumnName{Table: table, Name: "col1"})
	q.Return(ColumnName{Table: table, Name: "col2", Alias: "a1"})

	sql, binds, err := q.Get()
	if err != nil {
		t.Errorf("q.Get() returned %v", err)
	}

	var tag1, tag2, tag3 string

	for k, v := range binds {
		if v == 5 {
			tag1 = k
		} else if v == "str" {
			tag2 = k
		} else if v == "val" {
			tag3 = k
		}
	}

	if len(binds) != 3 {
		t.Errorf("q.Get() should have 2 values")
	}

	if sql != fmt.Sprintf("INSERT INTO %[1]s AS %[2]s (col1, col2) VALUES (@%[3]s, @%[4]s) ON CONFLICT (col1) DO UPDATE SET col2 = @%[5]s RETURNING %[2]s.col1, %[2]s.col2 AS a1", table.Name, table.Alias, tag1, tag2, tag3) {
		t.Errorf("q.Get() returned '%v'", sql)
	}
}
