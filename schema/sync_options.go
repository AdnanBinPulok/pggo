package schema

// SyncOptions controls how CreateTable synchronizes schema.
// SyncColumns: If true, it will add missing columns and drop extra columns.
// DropMissing: If true, it will drop columns that are in the database but not in the struct.
// DryRun: If true, it will only log the SQL statements without executing them.
// AllowDestructive: If true, it will allow destructive changes like dropping columns or tables.
type SyncOptions struct {
	SyncColumns      bool
	DropMissing      bool
	DryRun           bool
	AllowDestructive bool
}
