package builder

import (
	"fmt"
)

type UpdateQuery struct {
	table   *Table
	sets    []set
	where   Where
	binds   map[string]any
	returns []Column
}

func NewUpdate(table *Table) *UpdateQuery {
	return &UpdateQuery{
		table: table,
		binds: make(map[string]any),
	}
}

func (q *UpdateQuery) checkTable(table *Table) error {
	if table == nil {
		return fmt.Errorf("table cannot be nil")
	}

	if q.table != table {
		return fmt.Errorf("table %s does not exist", table.Name)
	}

	return nil
}

func (q *UpdateQuery) addBind(key string, value any) {
	q.binds[key] = value
}

func (q *UpdateQuery) Set(column string, value any) *UpdateQuery {
	q.sets = append(q.sets, set{
		Value:  value,
		Column: column,
	})

	return q
}

func (q *UpdateQuery) SetNow(column string) *UpdateQuery {
	q.sets = append(q.sets, set{
		Column: column,
		Now:    true,
	})

	return q
}

func (q *UpdateQuery) Where(w Where) *UpdateQuery {
	q.where = w

	return q
}

func (q *UpdateQuery) Return(c ...Column) *UpdateQuery {
	q.returns = append(q.returns, c...)

	return q
}

func (q *UpdateQuery) getSet() (string, error) {
	if len(q.sets) == 0 {
		return "", UpdateNoSets
	}

	s := " SET "

	for i, st := range q.sets {
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

		if i != len(q.sets)-1 {
			s += ", "
		}
	}

	return s, nil
}

func (q *UpdateQuery) getWhere() (string, error) {
	if q.where == nil {
		return "", nil
	}

	where, binds, err := q.where.gen(q)
	if err != nil {
		return "", err
	}

	if where == "" {
		return "", nil
	}

	for k, v := range binds {
		q.addBind(k, v)
	}

	return " WHERE " + where, nil
}

func (q *UpdateQuery) getReturns() (string, error) {
	if len(q.returns) == 0 {
		return "", nil
	}

	var s string

	for i, v := range q.returns {
		if v == nil {
			return "", fmt.Errorf("return column cannot be nil")
		}

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

func (q *UpdateQuery) Get() (string, map[string]any, error) {
	if err := validateDMLTargetTable(q.table); err != nil {
		return "", nil, err
	}

	q.binds = make(map[string]any)

	sets, err := q.getSet()
	if err != nil {
		return "", nil, err
	}

	where, err := q.getWhere()
	if err != nil {
		return "", nil, err
	}

	returns, err := q.getReturns()
	if err != nil {
		return "", nil, err
	}

	return "UPDATE " + q.table.Name + " AS " + q.table.Alias + sets + where + returns, q.binds, nil
}
