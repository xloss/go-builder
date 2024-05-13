package builder

import (
	"fmt"
	"strings"
)

type On interface {
	gen(query query) (string, error)
}

type OnAnd struct {
	List []On
}

func (o OnAnd) gen(q query) (string, error) {
	if q == nil {
		return "", fmt.Errorf("query cannot be nil")
	}

	if len(o.List) == 0 {
		return "", nil
	}

	list := make([]string, len(o.List))

	for i, on := range o.List {
		sql, err := on.gen(q)
		if err != nil {
			return "", err
		}

		list[i] = sql
	}

	return "(" + strings.Join(list, " AND ") + ")", nil
}

type OnEq struct {
	Table1  *Table
	Table2  *Table
	Column1 string
	Column2 string
}

func (o OnEq) gen(q query) (string, error) {
	if q == nil {
		return "", fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(o.Table1) {
		return "", fmt.Errorf("table %s does not exist", o.Table1.Name)
	}

	if !q.checkTable(o.Table2) {
		return "", fmt.Errorf("table %s does not exist", o.Table2.Name)
	}

	return o.Table1.Alias + "." + o.Column1 + " = " + o.Table2.Alias + "." + o.Column2, nil
}

type OnLess struct {
	Table1  *Table
	Table2  *Table
	Column1 string
	Column2 string
}

func (o OnLess) gen(q query) (string, error) {
	if q == nil {
		return "", fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(o.Table1) {
		return "", fmt.Errorf("table %s does not exist", o.Table1.Name)
	}

	if !q.checkTable(o.Table2) {
		return "", fmt.Errorf("table %s does not exist", o.Table2.Name)
	}

	return o.Table1.Alias + "." + o.Column1 + " < " + o.Table2.Alias + "." + o.Column2, nil
}

type OnMore struct {
	Table1  *Table
	Table2  *Table
	Column1 string
	Column2 string
}

func (o OnMore) gen(q query) (string, error) {
	if q == nil {
		return "", fmt.Errorf("query cannot be nil")
	}

	if !q.checkTable(o.Table1) {
		return "", fmt.Errorf("table %s does not exist", o.Table1.Name)
	}

	if !q.checkTable(o.Table2) {
		return "", fmt.Errorf("table %s does not exist", o.Table2.Name)
	}

	return o.Table1.Alias + "." + o.Column1 + " > " + o.Table2.Alias + "." + o.Column2, nil
}
