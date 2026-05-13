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

// --- Numeric Types ---

// SmallInt creates SMALLINT datatype (2-byte integer).
func (dt DataTypeFactory) SmallInt() *ColumnDef { return &ColumnDef{Type: "SMALLINT"} }

// Decimal creates DECIMAL(precision, scale) datatype.
func (dt DataTypeFactory) Decimal(precision, scale int) *ColumnDef {
	return &ColumnDef{Type: fmt.Sprintf("DECIMAL(%d,%d)", precision, scale)}
}

// Numeric creates NUMERIC(precision, scale) datatype.
func (dt DataTypeFactory) Numeric(precision, scale int) *ColumnDef {
	return &ColumnDef{Type: fmt.Sprintf("NUMERIC(%d,%d)", precision, scale)}
}

// Real creates REAL datatype (4-byte floating point).
func (dt DataTypeFactory) Real() *ColumnDef { return &ColumnDef{Type: "REAL"} }

// DoublePrecision creates DOUBLE PRECISION datatype (8-byte floating point).
func (dt DataTypeFactory) DoublePrecision() *ColumnDef { return &ColumnDef{Type: "DOUBLE PRECISION"} }

// SmallSerial creates SMALLSERIAL datatype (auto-incrementing 2-byte integer).
func (dt DataTypeFactory) SmallSerial() *ColumnDef { return &ColumnDef{Type: "SMALLSERIAL"} }

// BigSerial creates BIGSERIAL datatype (auto-incrementing 8-byte integer).
func (dt DataTypeFactory) BigSerial() *ColumnDef { return &ColumnDef{Type: "BIGSERIAL"} }

// Money creates MONEY datatype.
func (dt DataTypeFactory) Money() *ColumnDef { return &ColumnDef{Type: "MONEY"} }

// --- Character Types ---

// Char creates CHAR(n) fixed-length character datatype.
func (dt DataTypeFactory) Char(length int) *ColumnDef {
	return &ColumnDef{Type: fmt.Sprintf("CHAR(%d)", length)}
}

// --- Binary Types ---

// Bytea creates BYTEA datatype (binary data / byte array).
func (dt DataTypeFactory) Bytea() *ColumnDef { return &ColumnDef{Type: "BYTEA"} }

// --- Date/Time Types ---

// Date creates DATE datatype.
func (dt DataTypeFactory) Date() *ColumnDef { return &ColumnDef{Type: "DATE"} }

// Time creates TIME datatype (without time zone).
func (dt DataTypeFactory) Time() *ColumnDef { return &ColumnDef{Type: "TIME"} }

// TimeTz creates TIMETZ datatype (with time zone).
func (dt DataTypeFactory) TimeTz() *ColumnDef { return &ColumnDef{Type: "TIMETZ"} }

// Interval creates INTERVAL datatype.
func (dt DataTypeFactory) Interval() *ColumnDef { return &ColumnDef{Type: "INTERVAL"} }

// --- UUID ---

// UUID creates UUID datatype.
func (dt DataTypeFactory) UUID() *ColumnDef { return &ColumnDef{Type: "UUID"} }

// --- JSON Types ---

// JSON creates JSON datatype.
func (dt DataTypeFactory) JSON() *ColumnDef { return &ColumnDef{Type: "JSON"} }

// JSONB creates JSONB datatype (binary JSON).
func (dt DataTypeFactory) JSONB() *ColumnDef { return &ColumnDef{Type: "JSONB"} }

// --- XML ---

// XML creates XML datatype.
func (dt DataTypeFactory) XML() *ColumnDef { return &ColumnDef{Type: "XML"} }

// --- Network Address Types ---

// Inet creates INET datatype (IPv4 or IPv6 host address).
func (dt DataTypeFactory) Inet() *ColumnDef { return &ColumnDef{Type: "INET"} }

// Cidr creates CIDR datatype (IPv4 or IPv6 network address).
func (dt DataTypeFactory) Cidr() *ColumnDef { return &ColumnDef{Type: "CIDR"} }

// MacAddr creates MACADDR datatype.
func (dt DataTypeFactory) MacAddr() *ColumnDef { return &ColumnDef{Type: "MACADDR"} }

// MacAddr8 creates MACADDR8 datatype (EUI-64 format).
func (dt DataTypeFactory) MacAddr8() *ColumnDef { return &ColumnDef{Type: "MACADDR8"} }

// --- Bit String Types ---

// Bit creates BIT(n) fixed-length bit string datatype.
func (dt DataTypeFactory) Bit(length int) *ColumnDef {
	return &ColumnDef{Type: fmt.Sprintf("BIT(%d)", length)}
}

// VarBit creates BIT VARYING(n) variable-length bit string datatype.
func (dt DataTypeFactory) VarBit(length int) *ColumnDef {
	return &ColumnDef{Type: fmt.Sprintf("BIT VARYING(%d)", length)}
}

// --- Text Search Types ---

// TsVector creates TSVECTOR datatype (text search document).
func (dt DataTypeFactory) TsVector() *ColumnDef { return &ColumnDef{Type: "TSVECTOR"} }

// TsQuery creates TSQUERY datatype (text search query).
func (dt DataTypeFactory) TsQuery() *ColumnDef { return &ColumnDef{Type: "TSQUERY"} }

// --- Geometric Types ---

// Point creates POINT datatype (geometric point on a plane).
func (dt DataTypeFactory) Point() *ColumnDef { return &ColumnDef{Type: "POINT"} }

// Line creates LINE datatype (infinite line on a plane).
func (dt DataTypeFactory) Line() *ColumnDef { return &ColumnDef{Type: "LINE"} }

// Lseg creates LSEG datatype (line segment on a plane).
func (dt DataTypeFactory) Lseg() *ColumnDef { return &ColumnDef{Type: "LSEG"} }

// Box creates BOX datatype (rectangular box on a plane).
func (dt DataTypeFactory) Box() *ColumnDef { return &ColumnDef{Type: "BOX"} }

// Path creates PATH datatype (geometric path on a plane).
func (dt DataTypeFactory) Path() *ColumnDef { return &ColumnDef{Type: "PATH"} }

// Polygon creates POLYGON datatype (closed geometric path on a plane).
func (dt DataTypeFactory) Polygon() *ColumnDef { return &ColumnDef{Type: "POLYGON"} }

// Circle creates CIRCLE datatype (circle on a plane).
func (dt DataTypeFactory) Circle() *ColumnDef { return &ColumnDef{Type: "CIRCLE"} }

// --- Range Types ---

// Int4Range creates INT4RANGE datatype (range of integers).
func (dt DataTypeFactory) Int4Range() *ColumnDef { return &ColumnDef{Type: "INT4RANGE"} }

// Int8Range creates INT8RANGE datatype (range of bigints).
func (dt DataTypeFactory) Int8Range() *ColumnDef { return &ColumnDef{Type: "INT8RANGE"} }

// NumRange creates NUMRANGE datatype (range of numerics).
func (dt DataTypeFactory) NumRange() *ColumnDef { return &ColumnDef{Type: "NUMRANGE"} }

// TsRange creates TSRANGE datatype (range of timestamps without time zone).
func (dt DataTypeFactory) TsRange() *ColumnDef { return &ColumnDef{Type: "TSRANGE"} }

// TsTzRange creates TSTZRANGE datatype (range of timestamps with time zone).
func (dt DataTypeFactory) TsTzRange() *ColumnDef { return &ColumnDef{Type: "TSTZRANGE"} }

// DateRange creates DATERANGE datatype (range of dates).
func (dt DataTypeFactory) DateRange() *ColumnDef { return &ColumnDef{Type: "DATERANGE"} }

// --- Object Identifier Types ---

// OID creates OID datatype (object identifier).
func (dt DataTypeFactory) OID() *ColumnDef { return &ColumnDef{Type: "OID"} }

// --- Log Sequence Number ---

// PgLsn creates PG_LSN datatype (PostgreSQL log sequence number).
func (dt DataTypeFactory) PgLsn() *ColumnDef { return &ColumnDef{Type: "PG_LSN"} }
