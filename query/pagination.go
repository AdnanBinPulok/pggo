package query

import "fmt"

// NormalizePage validates page and limit and returns normalized values.
func NormalizePage(page, limit int) (int, int, error) {
	if page < 1 {
		return 0, 0, fmt.Errorf("page must be >= 1")
	}
	if limit < 1 {
		return 0, 0, fmt.Errorf("limit must be >= 1")
	}
	return page, limit, nil
}
