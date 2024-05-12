package builder

import (
	"strings"
	"testing"
)

func TestNewTable(t *testing.T) {
	table := NewTable("test_table")

	if table.Name != "test_table" {
		t.Errorf("table name is wrong")
	}

	if strings.HasPrefix("test_table_", table.Alias) {
		t.Errorf("table alias is wrong")
	}

	if len(table.Alias) != (len(table.Name) + randStrLen + 1) {
		t.Errorf("table alias length is wrong")
	}
}
