package builder

import "fmt"

type Group interface {
	gen(q query) (string, error)
}

type GroupColumn struct {
	Table  *Table
	Column string
}

func (g GroupColumn) gen(q query) (string, error) {
	if g.Table != nil {
		if !q.checkTable(g.Table) {
			return "", fmt.Errorf("table %s does not exist", g.Table.Name)
		}

		return g.Table.Alias + "." + g.Column, nil
	}

	return g.Column, nil
}
