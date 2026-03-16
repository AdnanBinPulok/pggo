package query

import (
	"fmt"
	"strings"
)

// BuildSearchClause builds a grouped ILIKE clause for the provided columns.
func BuildSearchClause(searchColumns []string, searchText string, startIndex int, quote func(string) (string, error)) (string, []any, error) {
	if len(searchColumns) == 0 || strings.TrimSpace(searchText) == "" {
		return "", nil, nil
	}
	parts := make([]string, 0, len(searchColumns))
	args := make([]any, 0, len(searchColumns))
	for i, col := range searchColumns {
		qcol, err := quote(col)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, fmt.Sprintf("CAST(%s AS TEXT) ILIKE $%d", qcol, startIndex+i))
		args = append(args, "%"+searchText+"%")
	}
	return "(" + strings.Join(parts, " OR ") + ")", args, nil
}
