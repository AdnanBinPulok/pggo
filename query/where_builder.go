package query

import (
	"fmt"
	"sort"
	"strings"
)

// BuildWhereClause turns a column-value map into SQL and args.
func BuildWhereClause(where map[string]any, startIndex int, quote func(string) (string, error)) (string, []any, error) {
	if len(where) == 0 {
		return "", nil, nil
	}
	keys := make([]string, 0, len(where))
	for k := range where {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	args := make([]any, 0, len(keys))
	idx := startIndex

	for _, col := range keys {
		qcol, err := quote(col)
		if err != nil {
			return "", nil, err
		}
		val := where[col]
		if cond, ok := val.(Condition); ok {
			sqlPart, condArgs, err := cond.ToSQL(qcol, idx)
			if err != nil {
				return "", nil, err
			}
			parts = append(parts, sqlPart)
			args = append(args, condArgs...)
			idx += len(condArgs)
			continue
		}
		parts = append(parts, fmt.Sprintf("%s = $%d", qcol, idx))
		args = append(args, val)
		idx++
	}

	return "WHERE " + strings.Join(parts, " AND "), args, nil
}
