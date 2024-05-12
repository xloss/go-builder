package builder

type join struct {
	Table *Table
	On    On
	Used  bool
	Left  bool
}

func (j join) Gen(query query) (string, error) {
	on, err := j.On.Gen(query)
	if err != nil {
		return "", err
	}

	s := ""

	if j.Left {
		s += " LEFT"
	}

	return s + " JOIN " + j.Table.Name + " AS " + j.Table.Alias + " ON " + on, nil
}
