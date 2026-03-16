package internal

import "errors"

var (
	// ErrNotFound is returned when no row matches the query.
	ErrNotFound = errors.New("row not found")
)

