package builder

import "fmt"

type DeleteQuery struct {
	table *Table
	where Where
	full  bool
	binds map[string]any
}

func NewDelete(table *Table) *DeleteQuery {
	return &DeleteQuery{
		table: table,
		binds: make(map[string]any),
	}
}

func (q *DeleteQuery) checkTable(table *Table) bool {
	return q.table == table
}

func (q *DeleteQuery) addBind(key string, value any) {
	q.binds[key] = value
}

func (q *DeleteQuery) Where(w Where) *DeleteQuery {
	q.where = w

	return q
}

func (q *DeleteQuery) Full() *DeleteQuery {
	q.full = true

	return q
}

func (q *DeleteQuery) getWhere() (string, error) {
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

func (q *DeleteQuery) Get() (string, map[string]any, error) {
	if q.table == nil {
		return "", nil, fmt.Errorf("table not set")
	}

	where, err := q.getWhere()
	if err != nil {
		return "", nil, err
	}

	if !q.full && where == "" {
		return "", nil, fmt.Errorf("use .Full() to delete without WHERE")
	}

	return "DELETE FROM " + q.table.Name + " AS " + q.table.Alias + where, q.binds, nil
}
