package builder

import (
	"fmt"
	"testing"
)

func TestSelectQuery_addBind(t *testing.T) {
	q := NewSelect()
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

func TestSelectQuery_checkTable(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)

	if !q.checkTable(table1) {
		t.Errorf("q.checkTable() returned false")
	}

	if q.checkTable(table2) {
		t.Errorf("q.checkTable() returned true")
	}

	q.IsSub()

	if !q.checkTable(table2) {
		t.Errorf("q.checkTable() returned false")
	}
}

func TestSelectQuery_From(t *testing.T) {
	table := NewTable("table")

	q := NewSelect()
	q.From(table)

	if len(q.from) != 1 {
		t.Errorf("q.from should have 1 column")
	}

	if q.from[0] != table {
		t.Errorf("q.from[0] should have table %v", table)
	}
}

func TestSelectQuery_Column(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()

	q.Column(ColumnName{Table: table, Name: "col1"})
	q.Column(ColumnName{Table: table, Name: "col2"})

	if len(q.columns) != 2 {
		t.Errorf("q.columns should have 2 column")
	}
}

func TestSelectQuery_LeftJoin(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)
	q.Column(ColumnName{Table: table2, Name: "col"})
	q.LeftJoin(table2, OnEq{
		Table1:  table1,
		Table2:  table2,
		Column1: "id",
		Column2: "table_id",
	})

	if len(q.joins) != 1 {
		t.Errorf("q.joins should have 1 join")
	}
	if q.joins[0].Table != table2 {
		t.Errorf("q.joins[0].Table should have table %v", table2)
	}
	if !q.joins[0].Left {
		t.Errorf("q.joins[0].Left should have true")
	}
	if q.joins[0].Used {
		t.Errorf("q.joins[0].Used should have false")
	}

	_, err := q.joins[0].Gen(q)
	if err != nil {
		t.Errorf("q.joins[0].Gen should not have error")
	}
	if !q.joins[0].Used {
		t.Errorf("q.joins[0].Used should have true")
	}
}

func TestSelectQuery_Where(t *testing.T) {
	table := NewTable("table")

	where := WhereEq{Table: table, Column: "col", Value: 1}

	q := NewSelect()

	if q.where != nil {
		t.Errorf("q.where should have nil")
	}

	q.Where(where)

	if q.where != where {
		t.Errorf("q.where should have where")
	}
}

func TestSelectQuery_Order(t *testing.T) {
	table := NewTable("table")

	order1 := Order{Table: table, Column: "col1", Desc: true}
	order2 := Order{Table: table, Column: "col2"}

	q := NewSelect()
	q.Order(order1, order2)

	if len(q.order) != 2 {
		t.Errorf("q.order should have 2 values")
	}

	if q.order[0] != order1 {
		t.Errorf("q.order[0] should have order1")
	}
	if q.order[1] != order2 {
		t.Errorf("q.order[1] should have order2")
	}
}

func TestSelectQuery_Limit(t *testing.T) {
	q := NewSelect()
	q.Limit(10)

	var (
		tag   string
		value int
	)

	for k, v := range q.binds {
		tag, value = k, v.(int)
	}

	if value != 10 {
		t.Errorf("value should have 10")
	}

	if q.limit != tag {
		t.Errorf("q.limit should have '%s'", tag)
	}
}

func TestSelectQuery_Offset(t *testing.T) {
	q := NewSelect()
	q.Offset(10)

	var (
		tag   string
		value int
	)

	for k, v := range q.binds {
		tag, value = k, v.(int)
	}

	if value != 10 {
		t.Errorf("value should have 10")
	}

	if q.offset != tag {
		t.Errorf("q.offset should have '%s'", tag)
	}
}

func TestSelectQuery_Group(t *testing.T) {
	table := NewTable("table")

	group1 := GroupColumn{Table: table, Column: "col1"}
	group2 := GroupColumn{Column: "col2"}

	q := NewSelect()
	q.Group(group1, group2)

	if len(q.group) != 2 {
		t.Errorf("q.group should have 2 values")
	}

	if q.group[0] != group1 {
		t.Errorf("q.group[0] should have group1")
	}
	if q.group[1] != group2 {
		t.Errorf("q.group[1] should have group2")
	}
}

func TestSelectQuery_IsSub(t *testing.T) {
	q := NewSelect()
	q.IsSub()

	if !q.isSub {
		t.Errorf("q.isSub should have true")
	}
}

func TestSelectQuery_getSelect(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1, table2)

	if _, err := q.getSelect(); err == nil {
		t.Errorf("q.getSelect should have error")
	}

	q.Column(ColumnName{Table: table1, Name: "col1"})
	q.Column(ColumnCoalesce{Table: table2, Name: "col2", Alias: "a1", Default: 10})

	s, err := q.getSelect()
	if err != nil {
		t.Errorf("q.getSelect should not have returned error. return: %e", err)
	}
	if s != fmt.Sprintf("SELECT %[1]s.col1, COALESCE(%[2]s.col2, 10) AS a1", table1.Alias, table2.Alias) {
		t.Errorf("bad returned select. return %s", s)
	}
}

func TestSelectQuery_getFrom(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()

	if _, err := q.getFrom(); err == nil {
		t.Errorf("q.getFrom should have error")
	}

	q.From(table1, table2)

	s, err := q.getFrom()
	if err != nil {
		t.Errorf("q.getFrom should not have returned error. return: %e", err)
	}
	if s != fmt.Sprintf(" FROM table1 AS %s, table2 AS %s", table1.Alias, table2.Alias) {
		t.Errorf("bad returned from. return '%s'", s)
	}
}

func TestSelectQuery_getWhere(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)
	q.LeftJoin(table2, OnEq{Table1: table2, Table2: table1, Column1: "id", Column2: "table_id"})

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

	q.Where(WhereEq{Table: table2, Column: "col", Value: 5})

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
	if where != " WHERE "+table2.Alias+".col = @"+tag {
		t.Errorf("bad returned where. return %s", where)
	}
}

func TestSelectQuery_getOrder(t *testing.T) {
	table := NewTable("table")

	q := NewSelect()
	q.From(table)

	order, err := q.getOrder()
	if err != nil {
		t.Errorf("q.getOrder should not have returned error. return: %e", err)
	}
	if order != "" {
		t.Errorf("order should have empty string")
	}

	q.Order(Order{Table: table, Column: "col1"})

	if len(q.order) != 1 {
		t.Errorf("q.order should have 1 values")
	}

	order, err = q.getOrder()
	if err != nil {
		t.Errorf("q.getOrder should not have returned error. return: %e", err)
	}
	if order != " ORDER BY "+table.Alias+".col1" {
		t.Errorf("bad returned select. return %s", order)
	}

	q.Order(Order{Column: "col2", Desc: true})

	if len(q.order) != 2 {
		t.Errorf("q.order should have 2 values")
	}

	order, err = q.getOrder()
	if err != nil {
		t.Errorf("q.getOrder should not have returned error. return: %e", err)
	}
	if order != " ORDER BY "+table.Alias+".col1, col2 DESC" {
		t.Errorf("bad returned order. return %s", order)
	}
}

func TestSelectQuery_getJoin(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")
	table3 := NewTable("table3")

	q := NewSelect()
	q.From(table1)
	q.LeftJoin(table2, OnEq{Table1: table1, Table2: table2, Column1: "id", Column2: "table_id"})
	q.LeftJoin(table3, OnEq{Table1: table3, Table2: table2, Column1: "id", Column2: "table_id"})

	j, err := q.getJoin()
	if err != nil {
		t.Errorf("q.getOrder should not have returned error. return: %e", err)
	}
	if j != "" {
		t.Errorf("j should have empty string")
	}

	q.Column(ColumnName{Table: table2, Name: "col1"})
	_, _ = q.getSelect()

	j, err = q.getJoin()
	if err != nil {
		t.Errorf("q.getOrder should not have returned error. return: %e", err)
	}
	if j != " LEFT JOIN "+table2.Name+" AS "+table2.Alias+" ON "+table1.Alias+".id = "+table2.Alias+".table_id" {
		t.Errorf("bad returned join. return %s", j)
	}

	q.Column(ColumnName{Table: table3, Name: "col2"})
	_, _ = q.getSelect()

	j, err = q.getJoin()
	if err != nil {
		t.Errorf("q.getOrder should not have returned error. return: %e", err)
	}
	if j != " LEFT JOIN "+table2.Name+" AS "+table2.Alias+" ON "+table1.Alias+".id = "+table2.Alias+".table_id LEFT JOIN "+table3.Name+" AS "+table3.Alias+" ON "+table3.Alias+".id = "+table2.Alias+".table_id" {
		t.Errorf("bad returned join. return %s", j)
	}
}

func TestSelectQuery_getGroup(t *testing.T) {
	table := NewTable("table")

	q := NewSelect()
	q.From(table)

	group, err := q.getGroup()
	if err != nil {
		t.Errorf("q.getGroup should not have returned error. return: %e", err)
	}
	if group != "" {
		t.Errorf("group should have empty string")
	}

	q.Group(GroupColumn{Table: table, Column: "col1"})

	if len(q.group) != 1 {
		t.Errorf("q.group should have 1 values")
	}

	group, err = q.getGroup()
	if err != nil {
		t.Errorf("q.getOrder should not have returned error. return: %e", err)
	}
	if group != " GROUP BY "+table.Alias+".col1" {
		t.Errorf("bad returned select. return %s", group)
	}

	q.Group(GroupColumn{Column: "col2"})

	if len(q.group) != 2 {
		t.Errorf("q.group should have 2 values")
	}

	group, err = q.getGroup()
	if err != nil {
		t.Errorf("q.getGroup should not have returned error. return: %e", err)
	}
	if group != " GROUP BY "+table.Alias+".col1, col2" {
		t.Errorf("bad returned order. return %s", group)
	}
}

func TestSelectQuery_Get(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.
		From(table1).
		Column(ColumnName{Table: table1, Name: "id"}).
		Column(ColumnName{Table: table2, Name: "col"}).
		LeftJoin(table2, OnEq{Table1: table1, Table2: table2, Column1: "id", Column2: "table_id"}).
		Where(WhereEq{Table: table1, Column: "id", Value: 1}).
		Order(Order{Table: table1, Column: "name", Desc: true}).
		Limit(10).
		Offset(5)

	sql, binds, err := q.Get()
	if err != nil {
		t.Errorf("q.Get should not have returned error. return: %e", err)
	}
	if len(binds) != 3 {
		t.Errorf("binds should have 1 values")
	}

	var w, l, o string

	for k, v := range binds {
		if v == 1 {
			w = k
		} else if v == 10 {
			l = k
		} else if v == 5 {
			o = k
		}
	}

	st := fmt.Sprintf("SELECT %[2]s.id, %[4]s.col FROM %[1]s AS %[2]s LEFT JOIN %[3]s AS %[4]s ON %[2]s.id = %[4]s.table_id WHERE %[2]s.id = @%[5]s ORDER BY %[2]s.name DESC LIMIT @%[6]s OFFSET @%[7]s", table1.Name, table1.Alias, table2.Name, table2.Alias, w, l, o)
	if sql != st {
		t.Errorf("bad returned sql. return:\n'%s'\n'%s'", sql, st)
	}
}

func ExampleNewSelect() {
	table1 := NewTable("table1")
	query1 := NewSelect()
	query1.Column(ColumnName{Table: table1, Name: "column1"})
	query1.From(table1)

	fmt.Println(query1.Get())

	// Result:
	// SELECT table1_pdddspkqfu.column1 FROM table1 AS table1_pdddspkqfu
	// map[]
	// <nil>
}
