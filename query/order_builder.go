package query

import (
	"fmt"
	"strings"
)

// BuildOrderBy validates order direction and builds ORDER BY clause.
func BuildOrderBy(orderBy, order string, quote func(string) (string, error)) (string, error) {
	if orderBy == "" {
		return "", nil
	}
	qcol, err := quote(orderBy)
	if err != nil {
		return "", err
	}
	dir := strings.ToUpper(strings.TrimSpace(order))
	if dir == "" {
		dir = "ASC"
	}
	if dir != "ASC" && dir != "DESC" {
		return "", fmt.Errorf("invalid order: %s", order)
	}
	return "ORDER BY " + qcol + " " + dir, nil
}
