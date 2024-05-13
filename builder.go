package builder

type query interface {
	checkTable(table *Table) bool
	addBind(key string, value any)
	Get() (string, map[string]any, error)
}

type column struct {
	Table     *Table
	Name      string
	Alias     string
	Aggregate bool
}

type set struct {
	Column string
	Value  interface{}
	Now    bool
}

type insertValue struct {
	Column string
	Value  interface{}
}

type Order struct {
	Table  *Table
	Column string
	Desc   bool
}
