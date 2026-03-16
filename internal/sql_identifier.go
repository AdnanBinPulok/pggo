package internal

import (
	"fmt"
	"regexp"
	"strings"
)

var identifierPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// ValidateIdentifier ensures table and column names are safe SQL identifiers.
func ValidateIdentifier(identifier string) error {
	if !identifierPattern.MatchString(identifier) {
		return fmt.Errorf("invalid sql identifier: %s", identifier)
	}
	return nil
}

// QuoteIdentifier safely quotes one SQL identifier.
func QuoteIdentifier(identifier string) (string, error) {
	if err := ValidateIdentifier(identifier); err != nil {
		return "", err
	}
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`, nil
}
