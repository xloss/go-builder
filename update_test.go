package builder

import (
	"errors"
	"fmt"
	"testing"
)

func TestUpdateQuery_addBind(t *testing.T) {
	q := NewUpdate(NewTable("table"))
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

func TestUpdateQuery_checkTable(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewUpdate(table1)

	if !q.checkTable(table1) {
		t.Errorf("q.checkTable() returned false")
	}

	if q.checkTable(table2) {
		t.Errorf("q.checkTable() returned true")
	}
}

func TestUpdateQuery_Set(t *testing.T) {
	table := NewTable("table")

	q := NewUpdate(table)
	q.Set(table, "col1", "value1")

	if len(q.sets) != 1 {
		t.Errorf("q.sets should have 1 set")
	}
	if q.sets[0].Table != table {
		t.Errorf("Table should have table")
	}
	if q.sets[0].Column != "col1" {
		t.Errorf("Column should have col1")
	}
	if q.sets[0].Value.(string) != "value1" {
		t.Errorf("Value should have value1")
	}
	if q.sets[0].Now {
		t.Errorf("now should have false")
	}
}

func TestUpdateQuery_SetNow(t *testing.T) {
	table := NewTable("table")

	q := NewUpdate(table)
	q.SetNow(table, "col1")

	if len(q.sets) != 1 {
		t.Errorf("q.sets should have 1 set")
	}
	if q.sets[0].Table != table {
		t.Errorf("Table should have table")
	}
	if q.sets[0].Column != "col1" {
		t.Errorf("Column should have col1")
	}
	if !q.sets[0].Now {
		t.Errorf("now should have true")
	}
}

func TestUpdateQuery_Where(t *testing.T) {
	table := NewTable("table")

	where := WhereEq{Table: table, Column: "col", Value: 1}

	q := NewUpdate(table)

	if q.where != nil {
		t.Errorf("q.where should have nil")
	}

	q.Where(where)

	if q.where != where {
		t.Errorf("q.where should have where")
	}
}

func TestUpdateQuery_getSet(t *testing.T) {
	table := NewTable("table")

	q := NewUpdate(table)

	sql, err := q.getSet()
	if !errors.Is(UpdateNoSets, err) {
		t.Errorf("q.getSet() returned %v", err)
	}

	q.Set(table, "col1", "value1")
	q.SetNow(table, "col2")

	sql, err = q.getSet()
	if err != nil {
		t.Errorf("q.getSet() returned %v", err)
	}

	var tag string

	for k, v := range q.binds {
		if v == "value1" {
			tag = k
		}
	}

	if sql != fmt.Sprintf(" SET %[1]s.col1 = @%[2]s, %[1]s.col2 = NOW()", table.Alias, tag) {
		t.Errorf("q.getSet() returned '%v'", sql)
	}
}

func TestUpdateQuery_getWhere(t *testing.T) {
	table := NewTable("table")

	q := NewUpdate(table)

	where, err := q.getWhere()
	if err != nil {
		t.Errorf("q.getWhere should not have returned error. return: %e", err)
	}
	if where != "" {
		t.Errorf("where not have empty string")
	}
	if len(q.binds) != 0 {
		t.Errorf("q.binds should have 0 values")
	}

	q.Where(WhereAnd{})

	where, err = q.getWhere()
	if err != nil {
		t.Errorf("q.getWhere should not have returned error. return: %e", err)
	}
	if where != "" {
		t.Errorf("where not have empty string")
	}
	if len(q.binds) != 0 {
		t.Errorf("q.binds should have 0 values")
	}

	q.Where(WhereEq{Table: table, Column: "col", Value: 5})

	where, err = q.getWhere()
	if err != nil {
		t.Errorf("q.getWhere should not have returned error. return: %e", err)
	}
	if where == "" {
		t.Errorf("where have empty string")
	}
	if len(q.binds) != 1 {
		t.Errorf("q.binds should have 1 values")
	}

	var (
		tag   string
		value int
	)

	for k, v := range q.binds {
		tag, value = k, v.(int)
	}

	if value != 5 {
		t.Errorf("value should have 5")
	}

	if err != nil {
		t.Errorf("q.getWhere should not have returned error. return: %e", err)
	}
	if where != " WHERE "+table.Alias+".col = @"+tag {
		t.Errorf("bad returned where. return %s", where)
	}
}

func TestUpdateQuery_Get(t *testing.T) {
	table := NewTable("table")

	q := NewUpdate(table)
	q.Set(table, "col1", "value1")
	q.SetNow(table, "col2")
	q.Where(WhereEq{Table: table, Column: "col3", Value: 5})

	sql, binds, err := q.Get()
	if err != nil {
		t.Errorf("q.Get() returned %v", err)
	}
	if len(binds) != 2 {
		t.Errorf("binds should have 1 values")
	}

	var val, where string

	for k, v := range q.binds {
		if v == "value1" {
			val = k
		} else if v == 5 {
			where = k
		}
	}

	st := fmt.Sprintf("UPDATE %[1]s AS %[2]s SET %[2]s.col1 = @%[3]s, %[2]s.col2 = NOW() WHERE %[2]s.col3 = @%[4]s", table.Name, table.Alias, val, where)
	if sql != st {
		t.Errorf("bad returned sql. return:\n'%s'\n'%s'", sql, st)
	}
}
