package builder

import (
	"fmt"
	"strings"
)

type Where interface {
	Gen(q query) (string, map[string]any, error)
}

type WhereEq struct {
	Table  *Table
	Column string
	Value  interface{}
}

func (w WhereEq) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	tag := w.Column + "_" + randStr()

	return w.Table.Alias + "." + w.Column + " = @" + tag, map[string]any{tag: w.Value}, nil
}

type WhereIsNull struct {
	Table  *Table
	Column string
}

func (w WhereIsNull) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	return w.Table.Alias + "." + w.Column + " IS NULL", map[string]any{}, nil
}

type WhereIsNotNull struct {
	Table  *Table
	Column string
}

func (w WhereIsNotNull) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	return w.Table.Alias + "." + w.Column + " IS NOT NULL", nil, nil
}

type WhereIn struct {
	Table  *Table
	Column string
	Values interface{}
}

func (w WhereIn) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	tag := w.Column + "_" + randStr()

	return w.Table.Alias + "." + w.Column + " = ANY(@" + tag + ")", map[string]any{tag: w.Values}, nil
}

type WhereMore struct {
	Table  *Table
	Column string
	Value  interface{}
}

func (w WhereMore) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	tag := w.Column + "_" + randStr()

	return w.Table.Alias + "." + w.Column + " > @" + tag, map[string]any{tag: w.Value}, nil
}

type WhereLess struct {
	Table  *Table
	Column string
	Value  interface{}
}

func (w WhereLess) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	tag := w.Column + "_" + randStr()

	q.addBind(tag, w.Value)

	return w.Table.Alias + "." + w.Column + " < @" + tag, map[string]any{tag: w.Value}, nil
}

type WhereMoreEq struct {
	Table  *Table
	Column string
	Value  interface{}
}

func (w WhereMoreEq) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	tag := w.Column + "_" + randStr()

	q.addBind(tag, w.Value)

	return w.Table.Alias + "." + w.Column + " >= @" + tag, map[string]any{tag: w.Value}, nil
}

type WhereLessEq struct {
	Table  *Table
	Column string
	Value  interface{}
}

func (w WhereLessEq) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	tag := w.Column + "_" + randStr()

	q.addBind(tag, w.Value)

	return w.Table.Alias + "." + w.Column + " <= @" + tag, map[string]any{tag: w.Value}, nil
}

type WhereMoreColumn struct {
	Table1, Table2   *Table
	Column1, Column2 string
}

func (w WhereMoreColumn) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table1) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table1.Name)
	}

	if !q.checkTable(w.Table2) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table2.Name)
	}

	return w.Table1.Alias + "." + w.Column1 + " > " + w.Table2.Alias + "." + w.Column2, map[string]any{}, nil
}

type WhereILike struct {
	Table  *Table
	Column string
	Value  interface{}
}

func (w WhereILike) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	tag := w.Column + "_" + randStr()

	return w.Table.Alias + "." + w.Column + " ILIKE @" + tag, map[string]any{tag: w.Value}, nil
}

type WhereFullText struct {
	Table    *Table
	Column   string
	Language string
	Value    string
}

func (w WhereFullText) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(w.Table) {
		return "", nil, fmt.Errorf("table %s does not exist", w.Table.Name)
	}

	tag := w.Column + "_" + randStr()

	return "to_tsvector('" + w.Language + "', " + w.Table.Alias + "." + w.Column + ") @@ plainto_tsquery(@" + tag + ")", map[string]any{tag: w.Value}, nil
}

type WhereAnd struct {
	List []Where
}

func (w WhereAnd) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if len(w.List) == 0 {
		return "", nil, nil
	}

	var (
		list  = make([]string, len(w.List))
		binds = make(map[string]any)
	)

	for i, where := range w.List {
		sql, bind, err := where.Gen(q)
		if err != nil {
			return "", nil, err
		}

		list[i] = sql

		for k, v := range bind {
			binds[k] = v
		}
	}

	return "(" + strings.Join(list, " AND ") + ")", binds, nil
}

type WhereOr struct {
	List []Where
}

func (w WhereOr) Gen(q query) (string, map[string]any, error) {
	if q == nil {
		return "", nil, fmt.Errorf("query cannot be nil")
	}

	if len(w.List) == 0 {
		return "", nil, nil
	}

	var (
		list  = make([]string, len(w.List))
		binds = make(map[string]any)
	)

	for i, where := range w.List {
		sql, bind, err := where.Gen(q)
		if err != nil {
			return "", nil, err
		}

		list[i] = sql

		for k, v := range bind {
			binds[k] = v
		}
	}

	return "(" + strings.Join(list, " OR ") + ")", binds, nil
}
