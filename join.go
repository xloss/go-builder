package builder

import "fmt"

type join struct {
	Table *Table
	On    On
	Used  bool
	Left  bool
	Inner bool
}

func (j *join) use(q query) error {
	if j.Used {
		return nil
	}

	j.Used = true

	if j.On == nil {
		return nil
	}

	return j.On.use(q)
}

func (j join) Gen(q query) (string, error) {
	if q == nil {
		return "", fmt.Errorf("query cannot be nil")
	}

	if err := q.checkTable(j.Table); err != nil {
		return "", err
	}

	if err := validateJoinTable(j.Table); err != nil {
		return "", err
	}

	if j.On == nil {
		return "", fmt.Errorf("join on cannot be nil")
	}

	on, err := j.On.gen(q)
	if err != nil {
		return "", err
	}

	table, binds, err := j.Table.gen()
	if err != nil {
		return "", err
	}

	for k, v := range binds {
		q.addBind(k, v)
	}

	s := ""

	if j.Left {
		s += " LEFT"
	} else if j.Inner {
		s += " INNER"
	}

	return s + " JOIN " + table + " ON " + on, nil
}
