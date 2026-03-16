package schema

import "fmt"

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
	cd.Constraints = append(cd.Constraints, fmt.Sprintf("DEFAULT '%v'", value))
	return cd
}

// String returns SQL fragment for this datatype definition.
func (cd ColumnDef) String() string {
	out := cd.Type
	for _, c := range cd.Constraints {
		out += " " + c
	}
	return out
}
