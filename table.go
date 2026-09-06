package builder

import "fmt"

type Table struct {
	Name  string
	Alias string
	Query query
}

func (t Table) gen() (string, map[string]any, error) {
	var (
		s     = ""
		binds = make(map[string]any)
		err   error
	)

	if t.Query != nil {
		s, binds, err = t.Query.Get()
		if err != nil {
			return "", nil, err
		}

		s = "(" + s + ")"
	} else {
		if errValidate := validateQualifiedIdentifier(t.Name, "table name"); errValidate != nil {
			return "", nil, errValidate
		}

		s = t.Name
	}

	if errValidate := validateIdentifier(t.Alias, "table alias"); errValidate != nil {
		return "", nil, errValidate
	}

	s = s + " AS " + t.Alias

	return s, binds, nil
}

// Creating Table struct for use in Builder
func NewTable(name string) *Table {
	return &Table{
		Name:  name,
		Alias: tableAlias(name),
	}
}

// Using Query as subquery in FROM
func NewTableSub(q query) *Table {
	return &Table{
		Alias: randStr() + "_" + randStr(),
		Query: q,
	}
}

func validateDMLTargetTable(table *Table) error {
	if table == nil {
		return fmt.Errorf("table not set")
	}

	if table.Query != nil {
		return fmt.Errorf("subquery cannot be used as DML table")
	}

	if err := validateQualifiedIdentifier(table.Name, "table name"); err != nil {
		return err
	}

	if err := validateIdentifier(table.Alias, "table alias"); err != nil {
		return err
	}

	return nil
}

func validateJoinTable(table *Table) error {
	if table == nil {
		return fmt.Errorf("join table cannot be nil")
	}

	if table.Query == nil {
		if err := validateQualifiedIdentifier(table.Name, "join table name"); err != nil {
			return err
		}
	}

	if err := validateIdentifier(table.Alias, "join table alias"); err != nil {
		return err
	}

	return nil
}
