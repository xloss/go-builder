package builder

import (
	"fmt"
	"strings"
)

type InsertQuery struct {
	table             *Table
	values            []insertValue
	conflict          []string
	conflictDoNothing bool
	update            []set
	returns           []Column
	binds             map[string]any
}

func NewInsert(table *Table) *InsertQuery {
	return &InsertQuery{
		table: table,
		binds: make(map[string]any),
	}
}

func (q *InsertQuery) checkTable(table *Table) error {
	if table == nil {
		return fmt.Errorf("table cannot be nil")
	}

	if q.table != table {
		return fmt.Errorf("table %s does not exist", table.Name)
	}

	return nil
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

func (q *InsertQuery) OnConflictDoNothing(c ...string) *InsertQuery {
	q.conflict = append(q.conflict, c...)
	q.conflictDoNothing = true

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
		if err := validateIdentifier(v.Column, "insert column"); err != nil {
			return "", err
		}

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

func (q *InsertQuery) getConflict() (string, error) {
	if len(q.conflict) == 0 {
		return "", nil
	}

	for _, column := range q.conflict {
		if err := validateIdentifier(column, "conflict column"); err != nil {
			return "", err
		}
	}

	s := " ON CONFLICT (" + strings.Join(q.conflict, ", ") + ")"

	if q.conflictDoNothing {
		s += " DO NOTHING"
	}

	return s, nil
}

func (q *InsertQuery) getUpdate() (string, error) {
	if len(q.update) == 0 || len(q.conflict) == 0 || q.conflictDoNothing {
		return "", nil
	}

	s := " DO UPDATE SET "

	for i, st := range q.update {
		if err := validateIdentifier(st.Column, "update column"); err != nil {
			return "", err
		}

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

	return s, nil
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
	if err := validateDMLTargetTable(q.table); err != nil {
		return "", nil, err
	}

	q.binds = make(map[string]any)

	values, err := q.getValues()
	if err != nil {
		return "", nil, err
	}

	conflict, err := q.getConflict()
	if err != nil {
		return "", nil, err
	}

	update, err := q.getUpdate()
	if err != nil {
		return "", nil, err
	}

	returns, err := q.getReturns()
	if err != nil {
		return "", nil, err
	}

	return "INSERT INTO " + q.table.Name + " AS " + q.table.Alias + values + conflict + update + returns, q.binds, nil
}
