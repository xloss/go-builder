package builder

import (
	"fmt"
	"strings"
)

type InsertQuery struct {
	table    *Table
	values   []insertValue
	conflict []string
	update   []set
	returns  []Column
	binds    map[string]any
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

func (q *InsertQuery) OnConflict(c ...string) *InsertQuery {
	q.conflict = append(q.conflict, c...)

	return q
}

func (q *InsertQuery) UpdateSet(column string, value any) *InsertQuery {
	q.update = append(q.update, set{Column: column, Value: value})

	return q
}

func (q *InsertQuery) UpdateSetNow(column string) *InsertQuery {
	q.update = append(q.update, set{Column: column, Now: true})

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

func (q *InsertQuery) getConflict() string {
	if len(q.conflict) == 0 {
		return ""
	}

	return " ON CONFLICT (" + strings.Join(q.conflict, ", ") + ")"
}

func (q *InsertQuery) getUpdate() string {
	if len(q.update) == 0 || len(q.conflict) == 0 {
		return ""
	}

	s := " DO UPDATE SET "

	for i, st := range q.update {
		s += st.Column + " = "

		if st.Now {
			s += "NOW()"
		} else {
			tag := st.Column + "_" + randStr()

			s += "@" + tag

			q.addBind(tag, st.Value)
		}

		if i != len(q.update)-1 {
			s += ", "
		}
	}

	return s
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

	return "INSERT INTO " + q.table.Name + " AS " + q.table.Alias + values + q.getConflict() + q.getUpdate() + returns, q.binds, nil
}
