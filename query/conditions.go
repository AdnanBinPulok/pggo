package query

import (
	"fmt"
	"strings"
)

// Condition renders SQL with positional arguments starting at startIndex.
type Condition interface {
	ToSQL(column string, startIndex int) (string, []any, error)
}

type simpleCondition struct {
	op     string
	values []any
}

func (c simpleCondition) ToSQL(column string, startIndex int) (string, []any, error) {
	switch c.op {
	case "IS NULL", "IS NOT NULL":
		return fmt.Sprintf("%s %s", column, c.op), nil, nil
	case "IN":
		if len(c.values) == 0 {
			return "", nil, fmt.Errorf("IN requires at least one value")
		}
		ph := make([]string, 0, len(c.values))
		for i := range c.values {
			ph = append(ph, fmt.Sprintf("$%d", startIndex+i))
		}
		return fmt.Sprintf("%s IN (%s)", column, strings.Join(ph, ", ")), c.values, nil
	case "BETWEEN":
		if len(c.values) != 2 {
			return "", nil, fmt.Errorf("BETWEEN requires exactly two values")
		}
		return fmt.Sprintf("%s BETWEEN $%d AND $%d", column, startIndex, startIndex+1), c.values, nil
	default:
		if len(c.values) != 1 {
			return "", nil, fmt.Errorf("%s requires exactly one value", c.op)
		}
		return fmt.Sprintf("%s %s $%d", column, c.op, startIndex), c.values, nil
	}
}

// In builds a WHERE IN condition.
func In(values ...any) Condition { return simpleCondition{op: "IN", values: values} }

// Between builds a BETWEEN condition.
func Between(min, max any) Condition { return simpleCondition{op: "BETWEEN", values: []any{min, max}} }

// IsNull builds an IS NULL condition.
func IsNull() Condition { return simpleCondition{op: "IS NULL"} }

// IsNotNull builds an IS NOT NULL condition.
func IsNotNull() Condition { return simpleCondition{op: "IS NOT NULL"} }

// Like builds a LIKE condition.
func Like(value any) Condition { return simpleCondition{op: "LIKE", values: []any{value}} }

// Gt builds a greater-than condition.
func Gt(value any) Condition { return simpleCondition{op: ">", values: []any{value}} }

// Lt builds a less-than condition.
func Lt(value any) Condition { return simpleCondition{op: "<", values: []any{value}} }

// Gte builds a greater-than-or-equal condition.
func Gte(value any) Condition { return simpleCondition{op: ">=", values: []any{value}} }

// Lte builds a less-than-or-equal condition.
func Lte(value any) Condition { return simpleCondition{op: "<=", values: []any{value}} }

// Neq builds a not-equal condition.
func Neq(value any) Condition { return simpleCondition{op: "!=", values: []any{value}} }
