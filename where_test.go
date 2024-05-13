package builder

import (
	"testing"
)

func TestWhereEq_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereEq{
		Table:  table,
		Column: "col",
		Value:  "value",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag, value string
	)

	for k, v := range binds {
		tag, value = k, v.(string)
	}

	if value != where.Value {
		t.Errorf("value is wrong")
	}

	if sql != table.Alias+".col = @"+tag {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereIsNull_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereIsNull{
		Table:  table,
		Column: "col",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 0 {
		t.Errorf("bind len should be 0, but got %v", len(binds))
	}

	if sql != table.Alias+".col IS NULL" {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereIsNotNull_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereIsNotNull{
		Table:  table,
		Column: "col",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 0 {
		t.Errorf("bind len should be 0, but got %v", len(binds))
	}

	if sql != table.Alias+".col IS NOT NULL" {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereIn_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereIn{
		Table:  table,
		Column: "col",
		Values: []int{
			1, 2, 3,
		},
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag   string
		value []int
	)

	for k, v := range binds {
		tag, value = k, v.([]int)
	}

	if len(value) != 3 {
		t.Errorf("value is wrong")
	}

	if value[0] != 1 || value[1] != 2 || value[2] != 3 {
		t.Errorf("bind is wrong")
	}

	if sql != table.Alias+".col = ANY(@"+tag+")" {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereMore_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereMore{
		Table:  table,
		Column: "col",
		Value:  "value",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag, value string
	)

	for k, v := range binds {
		tag, value = k, v.(string)
	}

	if value != where.Value {
		t.Errorf("value is wrong")
	}

	if sql != table.Alias+".col > @"+tag {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereLess_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereLess{
		Table:  table,
		Column: "col",
		Value:  "value",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag, value string
	)

	for k, v := range binds {
		tag, value = k, v.(string)
	}

	if value != where.Value {
		t.Errorf("value is wrong")
	}

	if sql != table.Alias+".col < @"+tag {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereMoreEq_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereMoreEq{
		Table:  table,
		Column: "col",
		Value:  "value",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag, value string
	)

	for k, v := range binds {
		tag, value = k, v.(string)
	}

	if value != where.Value {
		t.Errorf("value is wrong")
	}

	if sql != table.Alias+".col >= @"+tag {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereLessEq_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereLessEq{
		Table:  table,
		Column: "col",
		Value:  "value",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag, value string
	)

	for k, v := range binds {
		tag, value = k, v.(string)
	}

	if value != where.Value {
		t.Errorf("value is wrong")
	}

	if sql != table.Alias+".col <= @"+tag {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereMoreColumn_gen(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q := NewSelect()
	q.From(table1)
	q.From(table2)

	where := WhereMoreColumn{
		Table1:  table1,
		Column1: "col1",
		Table2:  table2,
		Column2: "col2",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 0 {
		t.Errorf("bind len should be 0, but got %v", len(binds))
	}

	if sql != table1.Alias+".col1 > "+table2.Alias+".col2" {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereILike_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereILike{
		Table:  table,
		Column: "col",
		Value:  "value",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag, value string
	)

	for k, v := range binds {
		tag, value = k, v.(string)
	}

	if value != where.Value {
		t.Errorf("value is wrong")
	}

	if sql != table.Alias+".col ILIKE @"+tag {
		t.Errorf("sql is wrong, sql is '%s'", sql)
	}
}

func TestWhereFullText_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereFullText{
		Table:    table,
		Language: "simple",
		Column:   "col",
		Value:    "value",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag, value string
	)

	for k, v := range binds {
		tag, value = k, v.(string)
	}

	if value != where.Value {
		t.Errorf("value is wrong")
	}

	if sql != "to_tsvector('simple', "+table.Alias+".col) @@ plainto_tsquery(@"+tag+")" {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereAnd_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereAnd{
		List: []Where{
			WhereEq{
				Table:  table,
				Column: "col1",
				Value:  "value1",
			},
			WhereEq{
				Table:  table,
				Column: "col2",
				Value:  "value2",
			},
		},
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 2 {
		t.Errorf("bind len should be 2, but got %v", len(binds))
	}

	var (
		tag1, tag2 string
	)

	for k, v := range binds {
		if v == "value1" {
			tag1 = k
		} else if v == "value2" {
			tag2 = k
		}
	}

	if sql != "("+table.Alias+".col1 = @"+tag1+" AND "+table.Alias+".col2 = @"+tag2+")" {
		t.Errorf("wrong sql: %s", sql)
	}
}

func TestWhereOr_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereOr{
		List: []Where{
			WhereEq{
				Table:  table,
				Column: "col1",
				Value:  "value1",
			},
			WhereEq{
				Table:  table,
				Column: "col1",
				Value:  "value2",
			},
		},
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}

	if len(binds) != 2 {
		t.Errorf("bind len should be 2, but got %v", len(binds))
	}

	var (
		tag1, tag2 string
	)

	for k, v := range binds {
		if v == "value1" {
			tag1 = k
		} else if v == "value2" {
			tag2 = k
		}
	}

	if sql != "("+table.Alias+".col1 = @"+tag1+" OR "+table.Alias+".col1 = @"+tag2+")" {
		t.Errorf("wrong sql: %s", sql)
	}
}

func TestWhereJsonbTextExist_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereJsonbTextExist{
		Table:  table,
		Column: "col",
		Value:  "value",
	}

	sql, binds, err := where.gen(q)
	if err != nil {
		t.Error(err)
	}
	if len(binds) != 1 {
		t.Errorf("bind len should be 1, but got %v", len(binds))
	}

	var (
		tag, value string
	)

	for k, v := range binds {
		tag, value = k, v.(string)
	}

	if value != where.Value {
		t.Errorf("value is wrong")
	}

	if sql != table.Alias+".col ? @"+tag {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}
