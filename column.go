package builder

import (
	"fmt"
)

type Column interface {
	gen(q query) (string, error)
}

type ColumnName struct {
	Table    *Table // required
	Name     string // required
	Alias    string
	Distinct bool
}

func (c ColumnName) gen(q query) (string, error) {
	if err := q.checkTable(c.Table); err != nil {
		return "", err
	}

	if err := validateIdentifier(c.Name, "column name"); err != nil {
		return "", err
	}

	if err := validateIdentifierIfNotEmpty(c.Alias, "column alias"); err != nil {
		return "", err
	}

	s := ""
	if c.Distinct {
		s += "DISTINCT "
	}

	s += c.Table.Alias + "." + c.Name

	if c.Alias != "" {
		s += " AS " + c.Alias
	}

	return s, nil
}

type ColumnCount struct {
	Table    *Table // required
	Name     string
	Alias    string // required
	Distinct bool
	Filter   Where
}

func (c ColumnCount) gen(q query) (string, error) {
	if err := validateIdentifier(c.Alias, "column alias"); err != nil {
		return "", err
	}

	s := "COUNT("

	if c.Name != "" {
		if err := q.checkTable(c.Table); err != nil {
			return "", err
		}

		if err := validateIdentifier(c.Name, "column name"); err != nil {
			return "", err
		}

		if c.Distinct {
			s += "DISTINCT "
		}

		s += c.Table.Alias + "." + c.Name
	} else {
		s += "*"
	}

	s += ")"

	if c.Filter != nil {
		filter, binds, err := c.Filter.gen(q)
		if err != nil {
			return "", err
		}

		if filter != "" {
			for k, v := range binds {
				q.addBind(k, v)
			}

			s += " FILTER (WHERE " + filter + ")"
		}
	}

	s += " AS " + c.Alias

	return s, nil
}

type ColumnCoalesce struct {
	Table   *Table // required
	Name    string // required
	Alias   string // required
	Default any    // required
}

func (c ColumnCoalesce) gen(q query) (string, error) {
	if err := q.checkTable(c.Table); err != nil {
		return "", err
	}

	if err := validateIdentifier(c.Name, "column name"); err != nil {
		return "", err
	}

	if err := validateIdentifier(c.Alias, "column alias"); err != nil {
		return "", err
	}

	if c.Default == nil {
		return "", fmt.Errorf("default is empty")
	}

	tag := c.Name + "_default_" + randStr()
	q.addBind(tag, c.Default)

	return "COALESCE(" + c.Table.Alias + "." + c.Name + ", @" + tag + ") AS " + c.Alias, nil
}

type ColumnJsonbArrayElementsText struct {
	Table    *Table // required
	Name     string // required
	Alias    string // required
	Distinct bool
}

func (c ColumnJsonbArrayElementsText) gen(q query) (string, error) {
	if err := q.checkTable(c.Table); err != nil {
		return "", err
	}

	if err := validateIdentifier(c.Name, "column name"); err != nil {
		return "", err
	}

	if err := validateIdentifier(c.Alias, "column alias"); err != nil {
		return "", err
	}

	s := ""
	if c.Distinct {
		s += "DISTINCT "
	}

	return s + "JSONB_ARRAY_ELEMENTS_TEXT(" + c.Table.Alias + "." + c.Name + ") AS " + c.Alias, nil
}

type ColumnValue struct {
	Value any // required
	Alias string
}

func (c ColumnValue) gen(q query) (string, error) {
	if q == nil {
		return "", fmt.Errorf("query cannot be nil")
	}

	if c.Value == nil {
		return "", fmt.Errorf("value is empty")
	}

	if err := validateIdentifierIfNotEmpty(c.Alias, "column alias"); err != nil {
		return "", err
	}

	tag := "value_" + randStr()
	q.addBind(tag, c.Value)

	s := "@" + tag

	if c.Alias != "" {
		s += " AS " + c.Alias
	}

	return s, nil
}

type columnRaw struct {
	Value string // required
}

func (c columnRaw) gen(_ query) (string, error) {
	if c.Value == "" {
		return "", fmt.Errorf("value is empty")
	}

	return c.Value, nil
}
