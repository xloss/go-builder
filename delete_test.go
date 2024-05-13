package builder

import (
	"testing"
)

func TestDeleteQuery_addBind(t *testing.T) {
	q := NewDelete(NewTable("table"))
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

func TestDeleteQuery_checkTable(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewDelete(table1)

	if !q.checkTable(table1) {
		t.Errorf("q.checkTable() returned false")
	}

	if q.checkTable(table2) {
		t.Errorf("q.checkTable() returned true")
	}
}

func TestDeleteQuery_Where(t *testing.T) {
	table := NewTable("table")

	where := WhereEq{Table: table, Column: "col", Value: 1}

	q := NewDelete(table)

	if q.where != nil {
		t.Errorf("q.where should have nil")
	}

	q.Where(where)

	if q.where != where {
		t.Errorf("q.where should have where")
	}
}

func TestDeleteQuery_Full(t *testing.T) {
	table := NewTable("table")
	q := NewDelete(table)

	if q.full {
		t.Errorf("q.full should have false")
	}

	q.Full()

	if !q.full {
		t.Errorf("q.full should have true")
	}
}

func TestDeleteQuery_getWhere(t *testing.T) {
	table := NewTable("table")
	q := NewDelete(table)

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

func TestDeleteQuery_Get(t *testing.T) {
	table := NewTable("table")
	q := NewDelete(table)

	_, _, err := q.Get()
	if err == nil {
		t.Errorf("q.Get() should have returned an error")
	}

	q.Where(WhereEq{Table: table, Column: "col", Value: 5})

	sql, binds, err := q.Get()
	if err != nil {
		t.Errorf("q.Get() should not have returned error. return: %e", err)
	}

	if len(binds) != 1 {
		t.Errorf("binds should have 1 values")
	}

	var tag string

	for k, v := range q.binds {
		if v == 5 {
			tag = k
		}
	}

	if sql != "DELETE FROM "+table.Name+" AS "+table.Alias+" WHERE "+table.Alias+".col = @"+tag {
		t.Errorf("bad returned where. return %s", sql)
	}

	q = NewDelete(table)
	q.Full()

	sql, binds, err = q.Get()
	if err != nil {
		t.Errorf("q.Get() should not have returned error. return: %e", err)
	}

	if len(binds) != 0 {
		t.Errorf("binds should have 0 values")
	}

	if sql != "DELETE FROM "+table.Name+" AS "+table.Alias {
		t.Errorf("bad returned where. return %s", sql)
	}
}
