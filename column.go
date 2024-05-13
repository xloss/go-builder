package builder

import (
	"fmt"
)

type Column interface {
	gen(q query) (string, error)
}

type ColumnName struct {
	Table *Table
	Name  string
	Alias string
}

func (c ColumnName) gen(q query) (string, error) {
	if !q.checkTable(c.Table) {
		return "", fmt.Errorf("table %s is not exist", c.Table)
	}

	if c.Name == "" {
		return "", fmt.Errorf("name is empty")
	}

	s := c.Table.Alias + "." + c.Name

	if c.Alias != "" {
		s += " AS " + c.Alias
	}

	return s, nil
}

type ColumnCount struct {
	Table *Table
	Name  string
	Alias string
}

func (c ColumnCount) gen(q query) (string, error) {
	if c.Alias == "" {
		return "", fmt.Errorf("alias is empty")
	}

	s := "COUNT("

	if q.checkTable(c.Table) && c.Name != "" {
		s += c.Table.Alias + "." + c.Name
	} else {
		s += "*"
	}

	s += ") AS " + c.Alias

	return s, nil
}

type ColumnCoalesce struct {
	Table   *Table
	Name    string
	Alias   string
	Default any
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
