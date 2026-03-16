package mapper

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// FieldMeta stores one model field mapping entry.
type FieldMeta struct {
	Index  int
	Column string
	Type   reflect.Type
}

// ModelMeta stores all field mappings for one model type.
type ModelMeta struct {
	FieldsByColumn map[string]FieldMeta
	Fields         []FieldMeta
}

var metaCache sync.Map

// ExtractModelMeta returns cached model metadata for T.
func ExtractModelMeta[T any]() (ModelMeta, error) {
	var zero T
	rt := reflect.TypeOf(zero)
	if rt == nil {
		return ModelMeta{}, fmt.Errorf("model type is nil")
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return ModelMeta{}, fmt.Errorf("model must be struct")
	}
	if cached, ok := metaCache.Load(rt); ok {
		return cached.(ModelMeta), nil
	}

	meta := ModelMeta{FieldsByColumn: map[string]FieldMeta{}}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.PkgPath != "" {
			continue
		}
		tag := strings.TrimSpace(f.Tag.Get("db"))
		if tag == "-" {
			continue
		}
		col := tag
		if idx := strings.Index(tag, ","); idx > -1 {
			col = tag[:idx]
		}
		if col == "" {
			col = ToSnakeCase(f.Name)
		}
		m := FieldMeta{Index: i, Column: col, Type: f.Type}
		meta.Fields = append(meta.Fields, m)
		meta.FieldsByColumn[col] = m
	}
	metaCache.Store(rt, meta)
	return meta, nil
}
