package builder

type Table struct {
	Name  string
	Alias string
}

func NewTable(t string) *Table {
	return &Table{
		Name:  t,
		Alias: t + "_" + randStr(),
	}
}
