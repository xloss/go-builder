package builder

type column struct {
	Table     *Table
	Name      string
	Alias     string
	Aggregate bool
}

type query interface {
	checkTable(table *Table) bool
	addBind(key string, value any)
	Get() (string, map[string]any, error)
}

type Order struct {
	Table  *Table
	Column string
	Desc   bool
}
