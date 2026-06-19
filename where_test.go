package builder

import (
	"strings"
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

func TestWhereNotEq_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereNotEq{
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

	if sql != table.Alias+".col <> @"+tag {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereEqColumn_gen(t *testing.T) {
	table1 := NewTable("table")
	table2 := NewTable("table")
	q := NewSelect()
	q.From(table1, table2)

	where := WhereEqColumn{
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

	if sql != table1.Alias+".col1 = "+table2.Alias+".col2" {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereNotEqColumn_gen(t *testing.T) {
	table1 := NewTable("table")
	table2 := NewTable("table")
	q := NewSelect()
	q.From(table1, table2)

	where := WhereNotEqColumn{
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

	if sql != table1.Alias+".col1 <> "+table2.Alias+".col2" {
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

func TestWhereFullText_genInvalidLanguage(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	list := []WhereFullText{
		{Table: table, Column: "col", Language: "bad language", Value: "value"},
		{Table: table, Column: "col", Language: "english'); DROP TABLE users; --", Value: "value"},
		{Table: table, Column: "col", Language: ".english", Value: "value"},
		{Table: table, Column: "col", Language: "public.", Value: "value"},
	}

	for _, where := range list {
		_, _, err := where.gen(q)
		if err == nil {
			t.Errorf("where.gen should have returned error")
		}
	}
}

func TestWhereFullText_genQualifiedLanguage(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereFullText{
		Table:    table,
		Column:   "col",
		Language: "public.english",
		Value:    "value",
	}

	sql, _, err := where.gen(q)
	if err != nil {
		t.Errorf("where.gen should not have returned error. return: %e", err)
	}

	if !strings.Contains(sql, "to_tsvector('public.english'") {
		t.Errorf("sql should contain qualified language, sql is %s", sql)
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

func TestWhereAnd_genNilWhere(t *testing.T) {
	q := NewSelect()

	where := WhereAnd{List: []Where{nil}}

	_, _, err := where.gen(q)
	if err == nil {
		t.Errorf("where.gen should have returned error")
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

func TestWhereOr_genNilWhere(t *testing.T) {
	q := NewSelect()

	where := WhereOr{List: []Where{nil}}

	_, _, err := where.gen(q)
	if err == nil {
		t.Errorf("where.gen should have returned error")
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

func TestWhereJsonbTextInExist_gen(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	where := WhereJsonbTextInExist{
		Table:  table,
		Column: "col",
		Values: []string{"value1", "value2"},
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
		value []string
	)

	for k, v := range binds {
		tag, value = k, v.([]string)
	}

	if len(value) != 2 {
		t.Errorf("value is wrong")
	}

	if value[0] != "value1" {
		t.Errorf("value1 is wrong")
	}

	if value[1] != "value2" {
		t.Errorf("value2 is wrong")
	}

	if sql != table.Alias+".col ?| @"+tag {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereExists_gen(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")
	q1 := NewSelect()
	q1.From(table1)
	q2 := NewSelect()
	q2.From(table2)
	q2.IsSub()
	q2.Where(WhereAnd{List: []Where{
		WhereEq{Table: table2, Column: "col1", Value: "value1"},
		WhereEqColumn{Table1: table1, Column1: "col3", Table2: table2, Column2: "col2"},
	}})

	where := WhereExists{
		Query: q2,
	}

	sql, binds, err := where.gen(q1)
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

	if value != "value1" {
		t.Errorf("value is wrong")
	}

	if sql != "EXISTS(SELECT 1 FROM table2 AS "+table2.Alias+" WHERE ("+table2.Alias+".col1 = @"+tag+" AND "+table1.Alias+".col3 = "+table2.Alias+".col2))" {
		t.Errorf("sql is wrong, sql is %s", sql)
	}
}

func TestWhereExists_genDoesNotAccumulateColumns(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q1 := NewSelect()
	q1.From(table1)

	q2 := NewSelect()
	q2.From(table2)
	q2.Where(WhereAnd{List: []Where{
		WhereEq{Table: table2, Column: "col1", Value: "value1"},
		WhereEqColumn{Table1: table1, Column1: "col3", Table2: table2, Column2: "col2"},
	}})

	where := WhereExists{
		Query: q2,
	}

	sql1, binds1, err := where.gen(q1)
	if err != nil {
		t.Error(err)
	}

	if len(q2.columns) != 0 {
		t.Errorf("q2.columns should have 0 values")
	}

	sql2, binds2, err := where.gen(q1)
	if err != nil {
		t.Error(err)
	}

	if len(q2.columns) != 0 {
		t.Errorf("q2.columns should have 0 values")
	}

	if strings.Count(sql1, "SELECT 1") != 1 {
		t.Errorf("sql1 should contain only one SELECT 1, sql is %s", sql1)
	}

	if strings.Count(sql2, "SELECT 1") != 1 {
		t.Errorf("sql2 should contain only one SELECT 1, sql is %s", sql2)
	}

	if len(binds1) != 1 {
		t.Errorf("binds1 should have 1 value")
	}

	if len(binds2) != 1 {
		t.Errorf("binds2 should have 1 value")
	}
}

func TestWhereExists_genRestoresColumns(t *testing.T) {
	table1 := NewTable("table1")
	table2 := NewTable("table2")

	q1 := NewSelect()
	q1.From(table1)

	q2 := NewSelect()
	q2.From(table2)
	q2.Column(ColumnName{Table: table2, Name: "col1"})

	where := WhereExists{
		Query: q2,
	}

	sql, _, err := where.gen(q1)
	if err != nil {
		t.Error(err)
	}

	if len(q2.columns) != 1 {
		t.Errorf("q2.columns should have 1 value")
	}

	if strings.Contains(sql, table2.Alias+".col1") {
		t.Errorf("exists sql should not use original select columns, sql is %s", sql)
	}

	if strings.Count(sql, "SELECT 1") != 1 {
		t.Errorf("sql should contain SELECT 1, sql is %s", sql)
	}
}

func TestWhere_genInvalidColumn(t *testing.T) {
	table := NewTable("table")
	q := NewSelect()
	q.From(table)

	list := []Where{
		WhereEq{Table: table, Column: "bad column", Value: 1},
		WhereNotEq{Table: table, Column: "bad column", Value: 1},
		WhereEqColumn{Table1: table, Column1: "bad column", Table2: table, Column2: "col2"},
		WhereEqColumn{Table1: table, Column1: "col1", Table2: table, Column2: "bad column"},
		WhereNotEqColumn{Table1: table, Column1: "bad column", Table2: table, Column2: "col2"},
		WhereNotEqColumn{Table1: table, Column1: "col1", Table2: table, Column2: "bad column"},
		WhereIsNull{Table: table, Column: "bad column"},
		WhereIsNotNull{Table: table, Column: "bad column"},
		WhereIn{Table: table, Column: "bad column", Values: []int{1, 2}},
		WhereMore{Table: table, Column: "bad column", Value: 1},
		WhereLess{Table: table, Column: "bad column", Value: 1},
		WhereMoreEq{Table: table, Column: "bad column", Value: 1},
		WhereLessEq{Table: table, Column: "bad column", Value: 1},
		WhereMoreColumn{Table1: table, Column1: "bad column", Table2: table, Column2: "col2"},
		WhereMoreColumn{Table1: table, Column1: "col1", Table2: table, Column2: "bad column"},
		WhereILike{Table: table, Column: "bad column", Value: "value"},
		WhereFullText{Table: table, Column: "bad column", Language: "english", Value: "value"},
		WhereJsonbTextExist{Table: table, Column: "bad column", Value: "value"},
		WhereJsonbTextInExist{Table: table, Column: "bad column", Values: []string{"value"}},
	}

	for _, where := range list {
		_, _, err := where.gen(q)
		if err == nil {
			t.Errorf("where.gen should have returned error")
		}
	}
}

func TestWhere_genDoesNotAddBindsToQuery(t *testing.T) {
	table := NewTable("table")

	list := []Where{
		WhereEq{Table: table, Column: "col", Value: 1},
		WhereNotEq{Table: table, Column: "col", Value: 1},
		WhereIn{Table: table, Column: "col", Values: []int{1, 2}},
		WhereMore{Table: table, Column: "col", Value: 1},
		WhereLess{Table: table, Column: "col", Value: 1},
		WhereMoreEq{Table: table, Column: "col", Value: 1},
		WhereLessEq{Table: table, Column: "col", Value: 1},
		WhereILike{Table: table, Column: "col", Value: "value"},
		WhereFullText{Table: table, Column: "col", Language: "english", Value: "value"},
		WhereJsonbTextExist{Table: table, Column: "col", Value: "value"},
		WhereJsonbTextInExist{Table: table, Column: "col", Values: []string{"value"}},
	}

	for _, where := range list {
		q := NewSelect()
		q.From(table)

		_, binds, err := where.gen(q)
		if err != nil {
			t.Errorf("where.gen should not have returned error. return: %e", err)
		}

		if len(binds) != 1 {
			t.Errorf("binds should have 1 value")
		}

		if len(q.binds) != 0 {
			t.Errorf("q.binds should have 0 values")
		}
	}
}
