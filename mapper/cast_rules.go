package mapper

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// CastAssign assigns src to dst with strict but practical conversions.
func CastAssign(dst reflect.Value, src any) error {
	if !dst.CanSet() {
		return fmt.Errorf("destination is not settable")
	}
	if src == nil {
		dst.Set(reflect.Zero(dst.Type()))
		return nil
	}

	sv := reflect.ValueOf(src)
	if sv.Type().AssignableTo(dst.Type()) {
		dst.Set(sv)
		return nil
	}
	if sv.Type().ConvertibleTo(dst.Type()) {
		dst.Set(sv.Convert(dst.Type()))
		return nil
	}

	if dst.Kind() == reflect.String {
		switch v := src.(type) {
		case []byte:
			dst.SetString(string(v))
			return nil
		}
	}

	if dst.Type() == reflect.TypeOf(time.Time{}) {
		if t, ok := src.(time.Time); ok {
			dst.Set(reflect.ValueOf(t))
			return nil
		}
	}

	if dst.Kind() >= reflect.Int && dst.Kind() <= reflect.Int64 {
		i64, err := toInt64(src)
		if err != nil {
			return err
		}
		dst.SetInt(i64)
		return nil
	}

	if dst.Kind() >= reflect.Uint && dst.Kind() <= reflect.Uint64 {
		u64, err := toUint64(src)
		if err != nil {
			return err
		}
		dst.SetUint(u64)
		return nil
	}

	if dst.Kind() == reflect.Float32 || dst.Kind() == reflect.Float64 {
		f64, err := toFloat64(src)
		if err != nil {
			return err
		}
		dst.SetFloat(f64)
		return nil
	}

	return fmt.Errorf("cannot cast %T to %s", src, dst.Type().String())
}

func toInt64(v any) (int64, error) {
	switch n := v.(type) {
	case int:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	case uint:
		return int64(n), nil
	case uint8:
		return int64(n), nil
	case uint16:
		return int64(n), nil
	case uint32:
		return int64(n), nil
	case uint64:
		return int64(n), nil
	case float32:
		return int64(n), nil
	case float64:
		return int64(n), nil
	case string:
		return strconv.ParseInt(n, 10, 64)
	default:
		return 0, fmt.Errorf("cannot cast %T to int64", v)
	}
}

func toUint64(v any) (uint64, error) {
	switch n := v.(type) {
	case int:
		return uint64(n), nil
	case int8:
		return uint64(n), nil
	case int16:
		return uint64(n), nil
	case int32:
		return uint64(n), nil
	case int64:
		return uint64(n), nil
	case uint:
		return uint64(n), nil
	case uint8:
		return uint64(n), nil
	case uint16:
		return uint64(n), nil
	case uint32:
		return uint64(n), nil
	case uint64:
		return n, nil
	case float32:
		return uint64(n), nil
	case float64:
		return uint64(n), nil
	case string:
		return strconv.ParseUint(n, 10, 64)
	default:
		return 0, fmt.Errorf("cannot cast %T to uint64", v)
	}
}

func toFloat64(v any) (float64, error) {
	switch n := v.(type) {
	case int:
		return float64(n), nil
	case int8:
		return float64(n), nil
	case int16:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint8:
		return float64(n), nil
	case uint16:
		return float64(n), nil
	case uint32:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case float32:
		return float64(n), nil
	case float64:
		return n, nil
	case string:
		return strconv.ParseFloat(n, 64)
	default:
		return 0, fmt.Errorf("cannot cast %T to float64", v)
	}
}
