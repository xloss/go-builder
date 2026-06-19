package builder

import (
	"fmt"
	"strings"
)

func validateIdentifier(name, field string) error {
	if name == "" {
		return fmt.Errorf("%s is empty", field)
	}

	if !isIdentifier(name) {
		return fmt.Errorf("%s %s is invalid", field, name)
	}

	return nil
}

func validateQualifiedIdentifier(name, field string) error {
	if name == "" {
		return fmt.Errorf("%s is empty", field)
	}

	list := strings.Split(name, ".")
	for _, item := range list {
		if !isIdentifier(item) {
			return fmt.Errorf("%s %s is invalid", field, name)
		}
	}

	return nil
}

func isIdentifier(name string) bool {
	for i, r := range name {
		if r > 127 {
			return false
		}

		if i == 0 {
			if isIdentifierLetter(r) || r == '_' {
				continue
			}

			return false
		}

		if isIdentifierLetter(r) || isIdentifierDigit(r) || r == '_' {
			continue
		}

		return false
	}

	return name != ""
}

func isIdentifierLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func isIdentifierDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func tableAlias(name string) string {
	return strings.ReplaceAll(name, ".", "_") + "_" + randStr()
}

func validateIdentifierIfNotEmpty(name, field string) error {
	if name == "" {
		return nil
	}

	return validateIdentifier(name, field)
}
