package builder

import (
	"fmt"
)

type UpdateQuery struct {
	table *Table
	sets  []set
	where Where
	binds map[string]any
}

func NewUpdate(table *Table) *UpdateQuery {
	return &UpdateQuery{
		table: table,
		binds: make(map[string]any),
	}
}

func (q *UpdateQuery) checkTable(table *Table) bool {
	return q.table == table
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

func (q *UpdateQuery) getSet() (string, error) {
	if len(q.sets) == 0 {
		return "", UpdateNoSets
	}

	s := " SET "

	for i, st := range q.sets {
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

func (q *UpdateQuery) Get() (string, map[string]any, error) {
	if q.table == nil {
		return "", nil, fmt.Errorf("table not set")
	}

	sets, err := q.getSet()
	if err != nil {
		return "", nil, err
	}

	where, err := q.getWhere()
	if err != nil {
		return "", nil, err
	}

	return "UPDATE " + q.table.Name + " AS " + q.table.Alias + sets + where, q.binds, nil
}
