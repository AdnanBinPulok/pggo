# PgGo Strict API

PgGo is a strict typed PostgreSQL toolkit for Go with:

- Generic tables: Table[T]
- Shared connection pool across all tables
- Consistent return style: (data, error)
- Safe SQL construction with PostgreSQL placeholders ($1, $2, ...)
- Typed CRUD, pagination, text search, and schema sync controls

This document is a full usage guide with practical examples.

## 1. Install

Add PgGo module from your workspace replace or your git source.

Example with local replace:

```go
module your-app

go 1.25

require pggo v0.0.0

replace pggo => ../pggo
```

Then run:

```bash
go mod tidy
```

## 2. Core Concepts

### 2.1 Shared Connection Pool

Create one connection once, then reuse for all tables.

```go
conn := pggo.NewDatabaseConnection(dbURL, 20, true)
defer conn.Close()
```

### 2.2 Typed Table

Each table is typed with your model.

```go
usersTable := pggo.Table[User]{
	Name:       "users",
	Connection: conn,
	Columns: []pggo.Column{
		{Name: "id", DataType: *pggo.DataType.Serial().PrimaryKey()},
		{Name: "name", DataType: *pggo.DataType.Text().NotNull()},
		{Name: "email", DataType: *pggo.DataType.Text().Unique().NotNull()},
		{Name: "age", DataType: *pggo.DataType.Integer()},
		{Name: "created_at", DataType: *pggo.DataType.TimestampTz()},
	},
}
```

### 2.3 Model Mapping

Struct tags map fields to DB columns.

```go
type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"created_at"`
}
```

If db tag is omitted, PgGo maps field names to snake_case.
Example: UserID -> user_id.

## 3. Table Lifecycle APIs

### 3.1 CreateTable

Creates table if it does not exist.

```go
err := usersTable.CreateTable()
if err != nil {
	return err
}
```

### 3.2 DeleteTable (Dangerous)

Drops table and data.

```go
err := usersTable.DeleteTable()
if err != nil {
	return err
}
```

### 3.3 DropTable Alias

DropTable is a compatibility alias for DeleteTable.

```go
_ = usersTable.DropTable()
```

## 4. Schema Sync Controls

Configure how CreateTable syncs schema:

```go
usersTable.SetSchemaSync(pggo.SyncOptions{
	SyncColumns:      true,
	DropMissing:      false,
	AllowDestructive: false,
	DryRun:           false,
})
```

Fields:

- SyncColumns: compare model columns with DB columns.
- DropMissing: allow dropping DB columns not defined in current table config.
- AllowDestructive: required gate for destructive drop behavior.
- DryRun: plan sync but skip execution.

Dangerous drops only execute when all of these are true:

- SyncColumns = true
- DropMissing = true
- AllowDestructive = true

## 5. Caching

Enable typed in-memory cache:

```go
usersTable.CacheKey = "id"
usersTable.EnableCache(5 * time.Second)
```

Notes:

- Cache key is taken from the mapped model field value for CacheKey column.
- Update and Delete clear cache for consistency.

## 6. CRUD APIs (All return data, error)

### 6.1 InsertOne

```go
inserted, err := usersTable.InsertOne(User{
	Name:  "Alice",
	Email: "alice@example.com",
	Age:   25,
})
if err != nil {
	return err
}
fmt.Println(inserted.ID)
```

InsertOne automatically skips zero-valued SERIAL/PRIMARY KEY id fields so Postgres can auto-generate IDs.

### 6.2 InsertMany

```go
rows, err := usersTable.InsertMany([]User{
	{Name: "Bob", Email: "bob@example.com", Age: 30},
	{Name: "Carol", Email: "carol@example.com", Age: 22},
})
if err != nil {
	return err
}
fmt.Println(len(rows))
```

### 6.3 Insert (Alias)

```go
inserted, err := usersTable.Insert(user)
```

### 6.4 FetchOne

Supports two forms:

1. Map condition

```go
user, err := usersTable.FetchOne(map[string]any{"id": 1})
```

2. Key/value shorthand

```go
user, err := usersTable.FetchOne("id", 1)
```

Returns internal.ErrNotFound when no row matches.

### 6.5 FetchMany

```go
users, err := usersTable.FetchMany(map[string]any{
	"age": pggo.Gte(18),
})
```

### 6.6 Update

```go
updatedRows, err := usersTable.Update(
	map[string]any{"age": 26},
	map[string]any{"id": 1},
)
```

### 6.7 Delete

```go
deletedRows, err := usersTable.Delete(map[string]any{"id": 1})
```

## 7. Count and Pagination

### 7.1 Count

```go
count, err := usersTable.Count(map[string]any{
	"age": pggo.Gt(20),
})
```

### 7.2 FetchPages

```go
pageData, err := usersTable.FetchPages(
	1,
	20,
	"id",
	"asc",
	map[string]any{"age": pggo.Gte(18)},
)
if err != nil {
	return err
}
fmt.Println(pageData.Items, pageData.Total, pageData.TotalPages)
```

PageResult[T] fields:

- Items
- Page
- Limit
- Total
- TotalPages

## 8. Search APIs

### 8.1 PageSearch

```go
searchPage, err := usersTable.PageSearch(
	1,
	10,
	"id",
	"desc",
	map[string]any{"age": pggo.Gte(18)},
	[]string{"name", "email"},
	"ali",
)
```

### 8.2 SearchCount

```go
searchCount, err := usersTable.SearchCount(
	map[string]any{"age": pggo.Gte(18)},
	[]string{"name", "email"},
	"example.com",
)
```

Search uses ILIKE against CAST(column AS TEXT).

## 9. Conditions

Available helpers:

- pggo.In(values...)
- pggo.Between(min, max)
- pggo.IsNull()
- pggo.IsNotNull()
- pggo.Like(value)
- pggo.Gt(value)
- pggo.Gte(value)
- pggo.Lt(value)
- pggo.Lte(value)
- pggo.Neq(value)

Example:

```go
users, err := usersTable.FetchMany(map[string]any{
	"age":   pggo.Between(18, 30),
	"email": pggo.Like("%example.com"),
})
```

## 10. SQL Safety Model

PgGo protects queries by:

- Using placeholders for values ($1, $2, ...)
- Validating identifiers (table/column names)
- Quoting identifiers safely

Important:

- PostgreSQL does not allow binding identifiers as query params.
- Therefore identifiers are validated and quoted, while values are parameterized.

## 11. Error Handling Guidelines

Recommended patterns:

```go
user, err := usersTable.FetchOne("id", 123)
if err != nil {
	if errors.Is(err, internal.ErrNotFound) {
		// handle not found
	}
	return err
}
_ = user
```

For now, internal.ErrNotFound is exposed from internal package path in module internals.
If you need a public sentinel error, add one at top-level pggo package.

## 12. Full End-to-End Example

See full runnable integration example in:

- ../test.go (workspace root)
- examples/strict_demo.go (module example function)

## 13. API Reference Summary

Table[T] methods:

- EnableCache(ttl)
- SetSchemaSync(opts)
- CreateTable() (error)
- DeleteTable() (error)
- DropTable() (error)
- PlanSchemaSync() (SyncPlan, error)
- GetColumnsFromDB() ([]string, error)
- InsertOne(data T) (T, error)
- InsertMany(data []T) ([]T, error)
- Insert(data T) (T, error)
- FetchOne(where any, args ...any) (T, error)
- FetchMany(where map[string]any) ([]T, error)
- Count(where map[string]any) (int64, error)
- FetchPages(page, limit int, orderBy, order string, where map[string]any) (PageResult[T], error)
- PageSearch(page, limit int, orderBy, order string, where map[string]any, searchColumns []string, searchText string) (PageResult[T], error)
- SearchCount(where map[string]any, searchColumns []string, searchText string) (int64, error)
- Update(set map[string]any, where map[string]any) ([]T, error)
- Delete(where map[string]any) ([]T, error)

All methods follow the project rule: return data, error.

