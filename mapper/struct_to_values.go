package mapper

import (
	"fmt"
	"reflect"
)

// StructToMap converts a typed model to column-value map based on db tags.
func StructToMap[T any](value T) (map[string]any, error) {
	meta, err := ExtractModelMeta[T]()
	if err != nil {
		return nil, err
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, fmt.Errorf("model pointer is nil")
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct")
	}
	out := make(map[string]any, len(meta.Fields))
	for _, f := range meta.Fields {
		out[f.Column] = rv.Field(f.Index).Interface()
	}
	return out, nil
}
