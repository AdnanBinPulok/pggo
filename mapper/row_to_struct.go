package mapper

import "reflect"

// MapToStruct converts one row map to model T.
func MapToStruct[T any](row map[string]any) (T, error) {
	var result T
	meta, err := ExtractModelMeta[T]()
	if err != nil {
		return result, err
	}
	rv := reflect.ValueOf(&result).Elem()
	for col, val := range row {
		f, ok := meta.FieldsByColumn[col]
		if !ok {
			continue
		}
		if err := CastAssign(rv.Field(f.Index), val); err != nil {
			return result, err
		}
	}
	return result, nil
}

// MapsToStructs converts list of row maps to model list.
func MapsToStructs[T any](rows []map[string]any) ([]T, error) {
	out := make([]T, 0, len(rows))
	for _, row := range rows {
		item, err := MapToStruct[T](row)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}
