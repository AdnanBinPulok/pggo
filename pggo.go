package pggo

import (
	"pggo/connection"
	"pggo/query"
	"pggo/schema"
)

// DatabaseConnection is the shared PostgreSQL pool wrapper used by all tables.
type DatabaseConnection = connection.DatabaseConnection

// Table represents a typed table handle where all CRUD methods return (data, error).
type Table[T any] = schema.Table[T]

// Column defines one table column.
type Column = schema.Column

// ColumnDef defines one SQL datatype and its constraints.
type ColumnDef = schema.ColumnDef

// DataType is a fluent datatype factory used to define table columns.
var DataType = schema.DataType{}

// PageResult carries paginated query results.
type PageResult[T any] = schema.PageResult[T]

// SyncOptions controls schema synchronization behavior.
type SyncOptions = schema.SyncOptions

// NewDatabaseConnection creates and initializes one shared pool.
func NewDatabaseConnection(dbURL string, maxConnections int, reconnect bool) *DatabaseConnection {
	conn := &DatabaseConnection{
		DBURL:          dbURL,
		MaxConnections: maxConnections,
		Reconnect:      reconnect,
	}
	_, err := conn.Connect()
	if err != nil {
		panic(err)
	}
	return conn
}

// In creates an IN condition.
var In = query.In

// Between creates a BETWEEN condition.
var Between = query.Between

// IsNull creates an IS NULL condition.
var IsNull = query.IsNull

// IsNotNull creates an IS NOT NULL condition.
var IsNotNull = query.IsNotNull

// Like creates a LIKE condition.
var Like = query.Like

// Gt creates a > condition.
var Gt = query.Gt

// Lt creates a < condition.
var Lt = query.Lt

// Gte creates a >= condition.
var Gte = query.Gte

// Lte creates a <= condition.
var Lte = query.Lte

// Neq creates a != condition.
var Neq = query.Neq
