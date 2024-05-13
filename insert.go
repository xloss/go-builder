package builder

import "fmt"

type InsertQuery struct {
	table   *Table
	values  []insertValue
	columns []column
	binds   map[string]any
}

func NewInsert(table *Table) *InsertQuery {
	return &InsertQuery{
		table: table,
		binds: make(map[string]any),
	}
}

func (q *InsertQuery) checkTable(table *Table) bool {
	return q.table == table
}

func (q *InsertQuery) addBind(key string, value any) {
	q.binds[key] = value
}

func (q *InsertQuery) Value(column string, v any) *InsertQuery {
	q.values = append(q.values, insertValue{Column: column, Value: v})

	return q
}

func (q *InsertQuery) Column(table *Table, name string) *InsertQuery {
	q.columns = append(q.columns, column{Table: table, Name: name})

	return q
}

func (q *InsertQuery) ColumnAlias(table *Table, name, alias string) *InsertQuery {
	q.columns = append(q.columns, column{Table: table, Name: name, Alias: alias})

	return q
}

func (q *InsertQuery) getValues() (string, error) {
	if len(q.values) == 0 {
		return "", fmt.Errorf("no values")
	}

	var c, t string

	for i, v := range q.values {
		tag := v.Column + "_" + randStr()

		c += v.Column
		t += "@" + tag

		q.addBind(tag, v.Value)

		if i < len(q.values)-1 {
			c += ", "
			t += ", "
		}
	}

	return " (" + c + ") VALUES (" + t + ")", nil
}

func (q *InsertQuery) getReturns() string {
	if len(q.columns) == 0 {
		return ""
	}

	var s string

	for i, v := range q.columns {
		s += v.Table.Alias + "." + v.Name

		if v.Alias != "" {
			s += " as " + v.Alias
		}

		if i < len(q.columns)-1 {
			s += ", "
		}
	}

	return " RETURNING " + s
}

func (q *InsertQuery) Get() (string, map[string]any, error) {
	if q.table == nil {
		return "", nil, fmt.Errorf("table not set")
	}

	values, err := q.getValues()
	if err != nil {
		return "", nil, err
	}

	returns := q.getReturns()

	return "INSERT INTO " + q.table.Name + " AS " + q.table.Alias + values + returns, q.binds, nil
}
