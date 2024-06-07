package builder

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
		s = t.Name
	}

	s = s + " AS " + t.Alias

	return s, binds, nil
}

// Creating Table struct for use in Builder
func NewTable(name string) *Table {
	return &Table{
		Name:  name,
		Alias: name + "_" + randStr(),
	}
}

// Using Query as subquery in FROM
func NewTableSub(q query) *Table {
	return &Table{
		Alias: randStr() + "_" + randStr(),
		Query: q,
	}
}
