package internal

import "fmt"

// ValidateColumns checks all provided columns exist in allowed set.
func ValidateColumns(allowed map[string]struct{}, columns []string) error {
	for _, c := range columns {
		if _, ok := allowed[c]; !ok {
			return fmt.Errorf("unknown column: %s", c)
		}
	}
	return nil
}
