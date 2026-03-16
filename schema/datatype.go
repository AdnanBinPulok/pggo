package schema

import "fmt"

// DataTypeFactory creates common SQL datatype definitions.
type DataTypeFactory struct{}

// DataType is the package-level fluent datatype factory.
type DataType = DataTypeFactory

// Serial creates SERIAL datatype.
func (dt DataTypeFactory) Serial() *ColumnDef { return &ColumnDef{Type: "SERIAL"} }

// Integer creates INTEGER datatype.
func (dt DataTypeFactory) Integer() *ColumnDef { return &ColumnDef{Type: "INTEGER"} }

// BigInt creates BIGINT datatype.
func (dt DataTypeFactory) BigInt() *ColumnDef { return &ColumnDef{Type: "BIGINT"} }

// Text creates TEXT datatype.
func (dt DataTypeFactory) Text() *ColumnDef { return &ColumnDef{Type: "TEXT"} }

// Varchar creates VARCHAR(n) datatype.
func (dt DataTypeFactory) Varchar(length int) *ColumnDef {
	return &ColumnDef{Type: fmt.Sprintf("VARCHAR(%d)", length)}
}

// Bool creates BOOLEAN datatype.
func (dt DataTypeFactory) Bool() *ColumnDef { return &ColumnDef{Type: "BOOLEAN"} }

// Timestamp creates TIMESTAMP datatype.
func (dt DataTypeFactory) Timestamp() *ColumnDef { return &ColumnDef{Type: "TIMESTAMP"} }

// TimestampTz creates TIMESTAMPTZ datatype.
func (dt DataTypeFactory) TimestampTz() *ColumnDef { return &ColumnDef{Type: "TIMESTAMPTZ"} }
