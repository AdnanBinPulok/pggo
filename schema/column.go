package schema

import (
	"fmt"
	"strings"
)

// Column defines one table column with name and datatype.
type Column struct {
	Name     string
	DataType ColumnDef
}

// ColumnDef defines SQL type and optional constraints.
type ColumnDef struct {
	Type        string
	Constraints []string
}

// NotNull marks a column as NOT NULL.
func (cd *ColumnDef) NotNull() *ColumnDef {
	cd.Constraints = append(cd.Constraints, "NOT NULL")
	return cd
}

// Unique marks a column as UNIQUE.
func (cd *ColumnDef) Unique() *ColumnDef {
	cd.Constraints = append(cd.Constraints, "UNIQUE")
	return cd
}

// PrimaryKey marks a column as PRIMARY KEY.
func (cd *ColumnDef) PrimaryKey() *ColumnDef {
	cd.Constraints = append(cd.Constraints, "PRIMARY KEY")
	return cd
}

// DefaultValue sets DEFAULT expression.
func (cd *ColumnDef) DefaultValue(value any) *ColumnDef {
	cd.Constraints = append(cd.Constraints, "DEFAULT "+formatDefaultValue(value))
	return cd
}

func formatDefaultValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "NULL"
	case string:
		if isKnownDefaultExpression(v) {
			return v
		}
		return fmt.Sprintf("'%s'", escapeSQLString(v))
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	default:
		return fmt.Sprintf("%v", value)
	}
}

func isKnownDefaultExpression(v string) bool {
	switch v {
	case "CURRENT_TIMESTAMP", "CURRENT_DATE", "CURRENT_TIME", "LOCALTIME", "LOCALTIMESTAMP", "NOW()", "now()":
		return true
	default:
		return false
	}
}

func escapeSQLString(v string) string {
	return strings.ReplaceAll(v, "'", "''")
}

// String returns SQL fragment for this datatype definition.
func (cd ColumnDef) String() string {
	out := cd.Type
	for _, c := range cd.Constraints {
		out += " " + c
	}
	return out
}
