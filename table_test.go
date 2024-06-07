package builder

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewTable(t *testing.T) {
	table := NewTable("test_table")

	if table.Name != "test_table" {
		t.Errorf("table name is wrong")
	}

	if strings.HasPrefix("test_table_", table.Alias) {
		t.Errorf("table alias is wrong")
	}

	if len(table.Alias) != (len(table.Name) + randStrLen + 1) {
		t.Errorf("table alias length is wrong")
	}
}

func TestNewTableSub(t *testing.T) {
	q := NewSelect()
	q.From(NewTable("table1"))

	table2 := NewTableSub(q)

	if table2.Query != q {
		t.Errorf("table query is wrong")
	}
}

func TestTable_gen(t *testing.T) {
	table1 := NewTable("table1")

	sql, binds, err := table1.gen()
	if err != nil {
		t.Error(err)
	}
	if len(binds) != 0 {
		t.Errorf("table bind length is wrong")
	}
	if sql != "table1 AS "+table1.Alias {
		t.Errorf("table sql is wrong, sql is '%s'", sql)
	}

	q := NewSelect()
	q.From(table1)
	q.Column(ColumnName{Table: table1, Name: "column1"})

	table2 := NewTableSub(q)

	sql, binds, err = table2.gen()
	if err != nil {
		t.Error(err)
	}
	if len(binds) != 0 {
		t.Errorf("table bind length is wrong")
	}
	if sql != fmt.Sprintf("(SELECT %[2]s.column1 FROM %[1]s AS %[2]s) AS %[3]s", table1.Name, table1.Alias, table2.Alias) {
		t.Errorf("table sql is wrong, sql is '%s'", sql)
	}
}

func TestName(t *testing.T) {
	table1 := NewTable("table1")
	q1 := NewSelect()
	q1.From(table1)
	q1.Column(ColumnName{Table: table1, Name: "column1"})

	table2 := NewTableSub(q1)

	q2 := NewSelect()
	q2.From(table2)
	q2.Column(ColumnName{Table: table2, Name: "column1"})

	_, _, err := q2.Get()
	if err != nil {
		t.Error(err)
	}
}

func ExampleNewTableSub() {
	table1 := NewTable("table1")
	query1 := NewSelect()
	query1.Column(ColumnName{Table: table1, Name: "column1"})
	query1.From(table1)

	table2 := NewTableSub(query1)
	query2 := NewSelect()
	query2.Column(ColumnName{Table: table2, Name: "column2"})
	query2.From(table2)

	fmt.Println(query2.Get())

	// Result:
	// SELECT yccakzcfbx_dxaolgmqrw.column2 FROM (SELECT table1_fmxwghcgnt.column1 FROM table1 AS table1_fmxwghcgnt) AS yccakzcfbx_dxaolgmqrw
	// map[]
	// <nil>
}
