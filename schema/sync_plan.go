package schema

// SyncPlan summarizes table sync actions.
type SyncPlan struct {
	AddColumns  []string
	DropColumns []string
}
