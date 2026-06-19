package builder

type Group interface {
	gen(q query) (string, error)
}

type GroupColumn struct {
	Table  *Table
	Column string
}

func (g GroupColumn) gen(q query) (string, error) {
	if err := validateIdentifier(g.Column, "group column"); err != nil {
		return "", err
	}

	if g.Table != nil {
		if err := q.checkTable(g.Table); err != nil {
			return "", err
		}

		return g.Table.Alias + "." + g.Column, nil
	}

	return g.Column, nil
}
