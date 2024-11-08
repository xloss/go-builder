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
	if !q.checkTable(c.Table) {
		return "", fmt.Errorf("table %s is not exist", c.Table)
	}

	if c.Name == "" {
		return "", fmt.Errorf("name is empty")
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
}

func (c ColumnCount) gen(q query) (string, error) {
	if c.Alias == "" {
		return "", fmt.Errorf("alias is empty")
	}

	s := "COUNT("

	if q.checkTable(c.Table) && c.Name != "" {
		if c.Distinct {
			s += "DISTINCT "
		}

		s += c.Table.Alias + "." + c.Name
	} else {
		s += "*"
	}

	s += ") AS " + c.Alias

	return s, nil
}

type ColumnCoalesce struct {
	Table   *Table // required
	Name    string // required
	Alias   string // required
	Default any    // required
}

func (c ColumnCoalesce) gen(q query) (string, error) {
	if !q.checkTable(c.Table) {
		return "", fmt.Errorf("table %s is not exist", c.Table)
	}

	if c.Name == "" {
		return "", fmt.Errorf("name is empty")
	}

	if c.Alias == "" {
		return "", fmt.Errorf("alias is empty")
	}

	if c.Default == nil {
		return "", fmt.Errorf("default is empty")
	}

	d := ""

	switch c.Default.(type) {
	case string:
		d = "'" + c.Default.(string) + "'"
	default:
		d = fmt.Sprintf("%v", c.Default)
	}

	return "COALESCE(" + c.Table.Alias + "." + c.Name + ", " + d + ") AS " + c.Alias, nil
}

type ColumnJsonbArrayElementsText struct {
	Table    *Table // required
	Name     string // required
	Alias    string // required
	Distinct bool
}

func (c ColumnJsonbArrayElementsText) gen(q query) (string, error) {
	if !q.checkTable(c.Table) {
		return "", fmt.Errorf("table %s is not exist", c.Table)
	}

	if c.Name == "" {
		return "", fmt.Errorf("name is empty")
	}

	if c.Alias == "" {
		return "", fmt.Errorf("alias is empty")
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

func (c ColumnValue) gen(_ query) (string, error) {
	if c.Value == nil {
		return "", fmt.Errorf("value is empty")
	}

	s := ""

	switch c.Value.(type) {
	case string:
		s = fmt.Sprintf("'%s'", c.Value)
	default:
		s = fmt.Sprintf("%v", c.Value)
	}

	if c.Alias != "" {
		s += " AS " + c.Alias
	}

	return s, nil
}
