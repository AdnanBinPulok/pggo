package schema

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/AdnanBinPulok/pggo/cache"
	"github.com/AdnanBinPulok/pggo/connection"
	"github.com/AdnanBinPulok/pggo/internal"
	"github.com/AdnanBinPulok/pggo/mapper"
	"github.com/AdnanBinPulok/pggo/query"

	"github.com/jackc/pgx/v5"
)

// PageResult holds paginated typed records and metadata.
//
// Fields:
// - Items: The page items.
// - Page: The current page number (1-based).
// - Limit: The requested limit per page.
// - Total: The total matching rows across all pages.
// - TotalPages: The computed total number of pages.
type PageResult[T any] struct {
	Items      []T
	Page       int
	Limit      int
	Total      int64
	TotalPages int
}

// Table is the strict typed table handle.
//
// All public CRUD methods return (data, error) and use one shared database pool
// through Connection. Cache is optional per table and can be enabled through
// EnableCache.
type Table[T any] struct {
	Name        string
	Connection  *connection.DatabaseConnection
	Columns     []Column
	DebugMode   bool
	CacheKey    string
	cache       *cache.Manager[T]
	cacheTTL    time.Duration
	SyncOptions SyncOptions
}

// EnableCache enables typed in-memory caching for this table.
//
// Parameters:
// - ttl: Cache time-to-live duration.
func (t *Table[T]) EnableCache(ttl time.Duration) {
	t.cacheTTL = ttl
	t.cache = cache.NewManager[T](ttl)
}

// SetSchemaSync configures schema synchronization behavior used by CreateTable.
//
// Parameters:
// - opts.SyncColumns: When true, compare table definition with DB columns.
// - opts.DropMissing: When true, missing-in-definition DB columns are candidates for drop.
// - opts.DryRun: When true, plan actions but do not execute changes.
// - opts.AllowDestructive: Must be true before DropMissing can execute.
func (t *Table[T]) SetSchemaSync(opts SyncOptions) {
	t.SyncOptions = opts
}

// CreateTable creates the table if needed and optionally synchronizes columns.
//
// Dangerous behavior:
//   - Dropping columns is possible only when all are true: SyncColumns, DropMissing,
//     and AllowDestructive. This explicit gate prevents accidental data loss.
//
// Returns:
// - error: Non-nil when SQL execution or validation fails.
func (t *Table[T]) CreateTable() error {
	if err := internal.ValidateIdentifier(t.Name); err != nil {
		return err
	}
	pool, err := t.Connection.Pool()
	if err != nil {
		return err
	}
	parts := make([]string, 0, len(t.Columns))
	for _, c := range t.Columns {
		qcol, err := internal.QuoteIdentifier(c.Name)
		if err != nil {
			return err
		}
		parts = append(parts, fmt.Sprintf("%s %s", qcol, c.DataType.String()))
	}
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return err
	}
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tbl, strings.Join(parts, ", "))
	_, err = pool.Exec(context.Background(), stmt)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	if t.SyncOptions.SyncColumns {
		plan, err := t.PlanSchemaSync()
		if err != nil {
			return err
		}
		if !t.SyncOptions.DryRun {
			if err := t.applySchemaSync(plan); err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteTable drops the table.
//
// Dangerous behavior:
// - This permanently removes table data.
//
// Returns:
// - error: Non-nil on SQL errors.
func (t *Table[T]) DeleteTable() error {
	pool, err := t.Connection.Pool()
	if err != nil {
		return err
	}
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return err
	}
	_, err = pool.Exec(context.Background(), "DROP TABLE IF EXISTS "+tbl)
	if err != nil {
		return fmt.Errorf("drop table: %w", err)
	}
	if t.cache != nil {
		t.cache.Clear()
	}
	return nil
}

// DropTable is a backward-compatible alias for DeleteTable.
func (t *Table[T]) DropTable() error {
	return t.DeleteTable()
}

// PlanSchemaSync compares DB columns and defined columns and returns a sync plan.
func (t *Table[T]) PlanSchemaSync() (SyncPlan, error) {
	dbCols, err := t.GetColumnsFromDB()
	if err != nil {
		return SyncPlan{}, err
	}
	defCols := t.definedColumnsSet()
	plan := SyncPlan{}
	for _, c := range t.Columns {
		if !contains(dbCols, c.Name) {
			plan.AddColumns = append(plan.AddColumns, c.Name)
		}
	}
	for _, c := range dbCols {
		if _, ok := defCols[c]; !ok {
			plan.DropColumns = append(plan.DropColumns, c)
		}
	}
	return plan, nil
}

// GetColumnsFromDB returns current database columns for this table.
func (t *Table[T]) GetColumnsFromDB() ([]string, error) {
	pool, err := t.Connection.Pool()
	if err != nil {
		return nil, err
	}
	rows, err := pool.Query(context.Background(),
		"SELECT column_name FROM information_schema.columns WHERE table_name = $1", t.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// InsertOne inserts one typed record and returns the inserted row.
//
// Parameters:
// - data: Typed model instance to insert.
//
// Returns:
// - T: Inserted row mapped back to model type.
// - error: Non-nil when validation, SQL, or mapping fails.
func (t *Table[T]) InsertOne(data T) (T, error) {
	var zero T
	valuesMap, err := mapper.StructToMap(data)
	if err != nil {
		return zero, err
	}
	filtered := t.filterAllowed(valuesMap)
	filtered = t.filterInsertable(filtered)
	if len(filtered) == 0 {
		return zero, fmt.Errorf("no insertable columns provided")
	}
	keys := sortedKeys(filtered)
	columns := make([]string, 0, len(keys))
	ph := make([]string, 0, len(keys))
	args := make([]any, 0, len(keys))
	for i, k := range keys {
		q, err := internal.QuoteIdentifier(k)
		if err != nil {
			return zero, err
		}
		columns = append(columns, q)
		ph = append(ph, fmt.Sprintf("$%d", i+1))
		args = append(args, filtered[k])
	}
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return zero, err
	}
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING *", tbl, strings.Join(columns, ", "), strings.Join(ph, ", "))
	rows, err := t.Connection.Pool()
	if err != nil {
		return zero, err
	}
	pgRows, err := rows.Query(context.Background(), stmt, args...)
	if err != nil {
		return zero, err
	}
	defer pgRows.Close()
	resultMaps, err := scanRowsToMaps(pgRows)
	if err != nil {
		return zero, err
	}
	if len(resultMaps) == 0 {
		return zero, internal.ErrNotFound
	}
	result, err := mapper.MapToStruct[T](resultMaps[0])
	if err != nil {
		return zero, err
	}
	t.cacheRecord(result)
	return result, nil
}

// InsertMany inserts many typed records and returns inserted rows.
func (t *Table[T]) InsertMany(data []T) ([]T, error) {
	out := make([]T, 0, len(data))
	for _, item := range data {
		inserted, err := t.InsertOne(item)
		if err != nil {
			return nil, err
		}
		out = append(out, inserted)
	}
	return out, nil
}

// Insert is a compatibility alias for InsertOne.
func (t *Table[T]) Insert(data T) (T, error) {
	return t.InsertOne(data)
}

// ToColumns converts one typed model to a raw column-value map.
//
// Notes:
// - Keys are DB column names resolved from `db` tags or snake_case field names.
// - Output is filtered to table-defined columns for safety.
func (t *Table[T]) ToColumns(data T) (map[string]any, error) {
	valuesMap, err := mapper.StructToMap(data)
	if err != nil {
		return nil, err
	}
	return t.filterAllowed(valuesMap), nil
}

// ToColumnsList converts multiple typed models to raw column-value maps.
func (t *Table[T]) ToColumnsList(data []T) ([]map[string]any, error) {
	out := make([]map[string]any, 0, len(data))
	for _, item := range data {
		row, err := t.ToColumns(item)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, nil
}

// ToValues is a backward-compatible alias for ToColumns.
func (t *Table[T]) ToValues(data T) (map[string]any, error) {
	return t.ToColumns(data)
}

// ToValuesList is a backward-compatible alias for ToColumnsList.
func (t *Table[T]) ToValuesList(data []T) ([]map[string]any, error) {
	return t.ToColumnsList(data)
}

// RawData is an alias of ToColumns for a more response-oriented naming style.
func (t *Table[T]) RawData(data T) (map[string]any, error) {
	return t.ToColumns(data)
}

// RawDataList is an alias of ToColumnsList.
func (t *Table[T]) RawDataList(data []T) ([]map[string]any, error) {
	return t.ToColumnsList(data)
}

// FetchOne returns one typed row matching where conditions.
//
// Accepted inputs:
// - map[string]any where each key is column name.
// - (string key, any value) shorthand using args.
//
// Returns:
// - T: Matched row.
// - error: ErrNotFound when no row exists, or SQL/mapping error.
func (t *Table[T]) FetchOne(where any, args ...any) (T, error) {
	var zero T
	whereMap, err := parseWhere(where, args...)
	if err != nil {
		return zero, err
	}
	if t.cache != nil && len(whereMap) == 1 {
		if raw, ok := whereMap[t.CacheKey]; ok {
			key := fmt.Sprint(raw)
			if cached, found := t.cache.Get(key); found {
				return cached, nil
			}
		}
	}
	items, err := t.FetchMany(whereMap)
	if err != nil {
		return zero, err
	}
	if len(items) == 0 {
		return zero, internal.ErrNotFound
	}
	item := items[0]
	t.cacheRecord(item)
	return item, nil
}

// FetchMany returns typed rows matching where conditions.
func (t *Table[T]) FetchMany(where map[string]any) ([]T, error) {
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return nil, err
	}
	whereSQL, whereArgs, err := query.BuildWhereClause(where, 1, internal.QuoteIdentifier)
	if err != nil {
		return nil, err
	}
	stmt := fmt.Sprintf("SELECT * FROM %s %s", tbl, whereSQL)
	pool, err := t.Connection.Pool()
	if err != nil {
		return nil, err
	}
	rows, err := pool.Query(context.Background(), stmt, whereArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resultMaps, err := scanRowsToMaps(rows)
	if err != nil {
		return nil, err
	}
	out, err := mapper.MapsToStructs[T](resultMaps)
	if err != nil {
		return nil, err
	}
	for _, item := range out {
		t.cacheRecord(item)
	}
	return out, nil
}

// Count returns number of rows matching where conditions.
func (t *Table[T]) Count(where map[string]any) (int64, error) {
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return 0, err
	}
	whereSQL, args, err := query.BuildWhereClause(where, 1, internal.QuoteIdentifier)
	if err != nil {
		return 0, err
	}
	stmt := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tbl, whereSQL)
	pool, err := t.Connection.Pool()
	if err != nil {
		return 0, err
	}
	var count int64
	if err := pool.QueryRow(context.Background(), stmt, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// FetchPages returns paginated typed rows using where/order rules.
func (t *Table[T]) FetchPages(page, limit int, orderBy, order string, where map[string]any) (PageResult[T], error) {
	var out PageResult[T]
	page, limit, err := query.NormalizePage(page, limit)
	if err != nil {
		return out, err
	}
	total, err := t.Count(where)
	if err != nil {
		return out, err
	}
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return out, err
	}
	whereSQL, whereArgs, err := query.BuildWhereClause(where, 1, internal.QuoteIdentifier)
	if err != nil {
		return out, err
	}
	orderSQL, err := query.BuildOrderBy(orderBy, order, internal.QuoteIdentifier)
	if err != nil {
		return out, err
	}
	offset := (page - 1) * limit
	args := append(whereArgs, limit, offset)
	stmt := fmt.Sprintf("SELECT * FROM %s %s %s LIMIT $%d OFFSET $%d", tbl, whereSQL, orderSQL, len(whereArgs)+1, len(whereArgs)+2)
	pool, err := t.Connection.Pool()
	if err != nil {
		return out, err
	}
	rows, err := pool.Query(context.Background(), stmt, args...)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	maps, err := scanRowsToMaps(rows)
	if err != nil {
		return out, err
	}
	items, err := mapper.MapsToStructs[T](maps)
	if err != nil {
		return out, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	out = PageResult[T]{Items: items, Page: page, Limit: limit, Total: total, TotalPages: totalPages}
	return out, nil
}

// PageSearch returns paginated typed rows with optional text search.
//
// Parameters:
// - page, limit, orderBy, order: Same semantics as FetchPages.
// - where: Base AND conditions.
// - searchColumns: Columns used for text search.
// - searchText: Text matched with ILIKE pattern.
func (t *Table[T]) PageSearch(page, limit int, orderBy, order string, where map[string]any, searchColumns []string, searchText string) (PageResult[T], error) {
	var out PageResult[T]
	page, limit, err := query.NormalizePage(page, limit)
	if err != nil {
		return out, err
	}
	total, err := t.SearchCount(where, searchColumns, searchText)
	if err != nil {
		return out, err
	}
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return out, err
	}
	whereSQL, whereArgs, err := query.BuildWhereClause(where, 1, internal.QuoteIdentifier)
	if err != nil {
		return out, err
	}
	searchSQL, searchArgs, err := query.BuildSearchClause(searchColumns, searchText, len(whereArgs)+1, internal.QuoteIdentifier)
	if err != nil {
		return out, err
	}
	fullWhere := whereSQL
	if searchSQL != "" {
		if fullWhere == "" {
			fullWhere = "WHERE " + searchSQL
		} else {
			fullWhere = fullWhere + " AND " + searchSQL
		}
	}
	orderSQL, err := query.BuildOrderBy(orderBy, order, internal.QuoteIdentifier)
	if err != nil {
		return out, err
	}
	offset := (page - 1) * limit
	args := append(whereArgs, searchArgs...)
	args = append(args, limit, offset)
	stmt := fmt.Sprintf("SELECT * FROM %s %s %s LIMIT $%d OFFSET $%d", tbl, fullWhere, orderSQL, len(args)-1, len(args))
	pool, err := t.Connection.Pool()
	if err != nil {
		return out, err
	}
	rows, err := pool.Query(context.Background(), stmt, args...)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	maps, err := scanRowsToMaps(rows)
	if err != nil {
		return out, err
	}
	items, err := mapper.MapsToStructs[T](maps)
	if err != nil {
		return out, err
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	out = PageResult[T]{Items: items, Page: page, Limit: limit, Total: total, TotalPages: totalPages}
	return out, nil
}

// SearchCount returns number of rows matching where + search rules without pagination.
func (t *Table[T]) SearchCount(where map[string]any, searchColumns []string, searchText string) (int64, error) {
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return 0, err
	}
	whereSQL, whereArgs, err := query.BuildWhereClause(where, 1, internal.QuoteIdentifier)
	if err != nil {
		return 0, err
	}
	searchSQL, searchArgs, err := query.BuildSearchClause(searchColumns, searchText, len(whereArgs)+1, internal.QuoteIdentifier)
	if err != nil {
		return 0, err
	}
	fullWhere := whereSQL
	if searchSQL != "" {
		if fullWhere == "" {
			fullWhere = "WHERE " + searchSQL
		} else {
			fullWhere = fullWhere + " AND " + searchSQL
		}
	}
	args := append(whereArgs, searchArgs...)
	stmt := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tbl, fullWhere)
	pool, err := t.Connection.Pool()
	if err != nil {
		return 0, err
	}
	var count int64
	if err := pool.QueryRow(context.Background(), stmt, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Delete removes rows matching where and returns deleted rows.
func (t *Table[T]) Delete(where map[string]any) ([]T, error) {
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return nil, err
	}
	whereSQL, args, err := query.BuildWhereClause(where, 1, internal.QuoteIdentifier)
	if err != nil {
		return nil, err
	}
	stmt := fmt.Sprintf("DELETE FROM %s %s RETURNING *", tbl, whereSQL)
	pool, err := t.Connection.Pool()
	if err != nil {
		return nil, err
	}
	rows, err := pool.Query(context.Background(), stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	maps, err := scanRowsToMaps(rows)
	if err != nil {
		return nil, err
	}
	out, err := mapper.MapsToStructs[T](maps)
	if err != nil {
		return nil, err
	}
	if t.cache != nil {
		t.cache.Clear()
	}
	return out, nil
}

// Update updates rows matching where and returns updated rows.
func (t *Table[T]) Update(set map[string]any, where map[string]any) ([]T, error) {
	if len(set) == 0 {
		return nil, fmt.Errorf("set cannot be empty")
	}
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return nil, err
	}
	set = t.filterAllowed(set)
	if len(set) == 0 {
		return nil, fmt.Errorf("set contains no known columns")
	}
	keys := sortedKeys(set)
	parts := make([]string, 0, len(keys))
	args := make([]any, 0, len(keys))
	for i, k := range keys {
		q, err := internal.QuoteIdentifier(k)
		if err != nil {
			return nil, err
		}
		parts = append(parts, fmt.Sprintf("%s = $%d", q, i+1))
		args = append(args, set[k])
	}
	whereSQL, whereArgs, err := query.BuildWhereClause(where, len(args)+1, internal.QuoteIdentifier)
	if err != nil {
		return nil, err
	}
	args = append(args, whereArgs...)
	stmt := fmt.Sprintf("UPDATE %s SET %s %s RETURNING *", tbl, strings.Join(parts, ", "), whereSQL)
	pool, err := t.Connection.Pool()
	if err != nil {
		return nil, err
	}
	rows, err := pool.Query(context.Background(), stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	maps, err := scanRowsToMaps(rows)
	if err != nil {
		return nil, err
	}
	out, err := mapper.MapsToStructs[T](maps)
	if err != nil {
		return nil, err
	}
	if t.cache != nil {
		t.cache.Clear()
	}
	return out, nil
}

func (t *Table[T]) applySchemaSync(plan SyncPlan) error {
	pool, err := t.Connection.Pool()
	if err != nil {
		return err
	}
	tbl, err := internal.QuoteIdentifier(t.Name)
	if err != nil {
		return err
	}
	for _, c := range plan.AddColumns {
		def, ok := t.columnDef(c)
		if !ok {
			continue
		}
		qcol, err := internal.QuoteIdentifier(c)
		if err != nil {
			return err
		}
		stmt := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s %s", tbl, qcol, def.String())
		if _, err := pool.Exec(context.Background(), stmt); err != nil {
			return err
		}
	}
	if t.SyncOptions.DropMissing && t.SyncOptions.AllowDestructive {
		for _, c := range plan.DropColumns {
			qcol, err := internal.QuoteIdentifier(c)
			if err != nil {
				return err
			}
			stmt := fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s", tbl, qcol)
			if _, err := pool.Exec(context.Background(), stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Table[T]) definedColumnsSet() map[string]struct{} {
	out := make(map[string]struct{}, len(t.Columns))
	for _, c := range t.Columns {
		out[c.Name] = struct{}{}
	}
	return out
}

func (t *Table[T]) columnDef(name string) (ColumnDef, bool) {
	for _, c := range t.Columns {
		if c.Name == name {
			return c.DataType, true
		}
	}
	return ColumnDef{}, false
}

func (t *Table[T]) filterAllowed(values map[string]any) map[string]any {
	allowed := t.definedColumnsSet()
	out := map[string]any{}
	for k, v := range values {
		if _, ok := allowed[k]; ok {
			out[k] = v
		}
	}
	return out
}

func (t *Table[T]) filterInsertable(values map[string]any) map[string]any {
	out := map[string]any{}
	for _, c := range t.Columns {
		v, ok := values[c.Name]
		if !ok {
			continue
		}
		if shouldSkipAutoColumn(c, v) {
			continue
		}
		out[c.Name] = v
	}
	return out
}

func shouldSkipAutoColumn(column Column, value any) bool {
	columnType := strings.ToUpper(column.DataType.Type)
	isAutoType := strings.Contains(columnType, "SERIAL")
	isPrimary := false
	for _, c := range column.DataType.Constraints {
		if strings.EqualFold(strings.TrimSpace(c), "PRIMARY KEY") {
			isPrimary = true
			break
		}
	}
	if !(isAutoType || isPrimary) {
		return false
	}
	return isZeroValue(value)
}

func isZeroValue(value any) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func:
		return rv.IsNil()
	default:
		return rv.IsZero()
	}
}

func (t *Table[T]) cacheRecord(item T) {
	if t.cache == nil || t.CacheKey == "" {
		return
	}
	m, err := mapper.StructToMap(item)
	if err != nil {
		return
	}
	v, ok := m[t.CacheKey]
	if !ok {
		return
	}
	t.cache.Set(fmt.Sprint(v), item)
}

func parseWhere(where any, args ...any) (map[string]any, error) {
	if where == nil {
		return map[string]any{}, nil
	}
	if m, ok := where.(map[string]any); ok {
		return m, nil
	}
	if k, ok := where.(string); ok && len(args) == 1 {
		return map[string]any{k: args[0]}, nil
	}
	return nil, fmt.Errorf("invalid where format")
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func scanRowsToMaps(rows pgx.Rows) ([]map[string]any, error) {
	result := []map[string]any{}
	fields := rows.FieldDescriptions()
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		row := map[string]any{}
		for i, f := range fields {
			row[string(f.Name)] = vals[i]
		}
		result = append(result, row)
	}
	return result, rows.Err()
}
