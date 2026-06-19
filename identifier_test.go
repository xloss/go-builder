package builder

import "testing"

func TestIsIdentifier(t *testing.T) {
	valid := []string{
		"table",
		"table1",
		"table_1",
		"_table",
		"T",
		"T1",
	}

	for _, name := range valid {
		if !isIdentifier(name) {
			t.Errorf("%s should be valid identifier", name)
		}
	}

	invalid := []string{
		"",
		"1table",
		"table-name",
		"table name",
		"table.name",
		"table;",
		"table$",
		"таблица",
	}

	for _, name := range invalid {
		if isIdentifier(name) {
			t.Errorf("%s should be invalid identifier", name)
		}
	}
}

func TestValidateIdentifier(t *testing.T) {
	if err := validateIdentifier("table1", "field"); err != nil {
		t.Errorf("validateIdentifier should not have returned error. return: %e", err)
	}

	if err := validateIdentifier("", "field"); err == nil {
		t.Errorf("validateIdentifier should have returned error")
	}

	if err := validateIdentifier("table-name", "field"); err == nil {
		t.Errorf("validateIdentifier should have returned error")
	}
}

func TestValidateQualifiedIdentifier(t *testing.T) {
	valid := []string{
		"table",
		"public.table",
		"db_schema.table_name",
		"a.b.c",
	}

	for _, name := range valid {
		if err := validateQualifiedIdentifier(name, "field"); err != nil {
			t.Errorf("validateQualifiedIdentifier should not have returned error. return: %e", err)
		}
	}

	invalid := []string{
		"",
		".table",
		"table.",
		"public..table",
		"public.table-name",
		"public.table name",
		"public.table;",
	}

	for _, name := range invalid {
		if err := validateQualifiedIdentifier(name, "field"); err == nil {
			t.Errorf("validateQualifiedIdentifier should have returned error")
		}
	}
}

func TestTableAlias(t *testing.T) {
	alias := tableAlias("public.users")

	if len(alias) != len("public_users")+randStrLen+1 {
		t.Errorf("alias length is wrong")
	}

	if alias[:len("public_users_")] != "public_users_" {
		t.Errorf("alias prefix is wrong")
	}
}
