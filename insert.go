package builder

import "fmt"

type InsertQuery struct {
	table   *Table
	values  []insertValue
	returns []Column
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

func (q *InsertQuery) Return(c ...Column) *InsertQuery {
	q.returns = append(q.returns, c...)

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

func (q *InsertQuery) getReturns() (string, error) {
	if len(q.returns) == 0 {
		return "", nil
	}

	var s string

	for i, v := range q.returns {
		c, err := v.gen(q)
		if err != nil {
			return "", err
		}

		s += c

		if i != len(q.returns)-1 {
			s += ", "
		}
	}

	return " RETURNING " + s, nil
}

func (q *InsertQuery) Get() (string, map[string]any, error) {
	if q.table == nil {
		return "", nil, fmt.Errorf("table not set")
	}

	values, err := q.getValues()
	if err != nil {
		return "", nil, err
	}

	returns, err := q.getReturns()
	if err != nil {
		return "", nil, err
	}

	return "INSERT INTO " + q.table.Name + " AS " + q.table.Alias + values + returns, q.binds, nil
}
