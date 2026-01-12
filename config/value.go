package config

import (
	"encoding/json"
	"fmt"
	pkgJson "go-pkg/encoding/json"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/proto"
)

type errValue struct {
	err error
}

func (v errValue) Bool() (bool, error)              { return false, v.err }
func (v errValue) Int() (int64, error)              { return 0, v.err }
func (v errValue) Float() (float64, error)          { return 0.0, v.err }
func (v errValue) Duration() (time.Duration, error) { return 0, v.err }
func (v errValue) String() (string, error)          { return "", v.err }
func (v errValue) Scan(any) error                   { return v.err }
func (v errValue) Load() any                        { return nil }
func (v errValue) Store(any)                        {}
func (v errValue) Slice() ([]Value, error)          { return nil, v.err }
func (v errValue) Map() (map[string]Value, error)   { return nil, v.err }

type Value interface {
	Bool() (bool, error)
	Int() (int64, error)
	Float() (float64, error)
	String() (string, error)
	Duration() (time.Duration, error)
	Slice() ([]Value, error)
	Map() (map[string]Value, error)
	Scan(v any) error
	Load() any
	Store(any)
}

type atomicValue struct {
	atomic.Value
}

var (
	_ Value = (*atomicValue)(nil)
	_ Value = (*errValue)(nil)
)

// typeAssertError returns a type assertion error.
func (v *atomicValue) typeAssertError() error {
	return fmt.Errorf("type assert to %v failed", reflect.TypeOf(v.Load()))
}

// Bool retrieves the boolean value.
func (v *atomicValue) Bool() (bool, error) {
	switch val := v.Load().(type) {
	case bool:
		return val, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return strconv.ParseBool(fmt.Sprint(val))
	case string:
		return strconv.ParseBool(val)
	default:
		return false, v.typeAssertError()
	}
}

// Int retrieves the integer value.
func (v *atomicValue) Int() (int64, error) {
	switch val := v.Load().(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		return 0, v.typeAssertError()
	}
}

// Float retrieves the float value.
func (v *atomicValue) Float() (float64, error) {
	switch val := v.Load().(type) {
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, v.typeAssertError()
	}
}

// String retrieves the string value.
func (v *atomicValue) String() (string, error) {
	switch val := v.Load().(type) {
	case string:
		return val, nil
	case bool:
		return strconv.FormatBool(val), nil
	case int, int8, int16, int32, int64:
		return fmt.Sprint(val), nil
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprint(val), nil
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	default:
		return "", v.typeAssertError()
	}
}

// Duration retrieves the duration value.
func (v *atomicValue) Duration() (time.Duration, error) {
	switch val := v.Load().(type) {
	case time.Duration:
		return val, nil
	case int, int8, int16, int32, int64:
		return time.Duration(reflect.ValueOf(val).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return time.Duration(reflect.ValueOf(val).Uint()), nil
	case float32:
		return time.Duration(val), nil
	case float64:
		return time.Duration(val), nil
	case string:
		return time.ParseDuration(val)
	default:
		return 0, v.typeAssertError()
	}
}

// Slice retrieves the slice value.
func (v *atomicValue) Slice() ([]Value, error) {
	vals, ok := v.Load().([]any)
	if !ok {
		return nil, v.typeAssertError()
	}
	slices := make([]Value, 0, len(vals))
	for _, val := range vals {
		a := new(atomicValue)
		a.Store(val)
		slices = append(slices, a)
	}
	return slices, nil
}

func (v *atomicValue) Map() (map[string]Value, error) {
	vals, ok := v.Load().(map[string]any)
	if !ok {
		return nil, v.typeAssertError()
	}
	m := make(map[string]Value, len(vals))
	for key, val := range vals {
		a := new(atomicValue)
		a.Store(val)
		m[key] = a
	}
	return m, nil
}

func (v *atomicValue) Scan(obj any) error {
	data, err := json.Marshal(v.Load())
	if err != nil {
		return err
	}
	if pb, ok := obj.(proto.Message); ok {
		return pkgJson.UnmarshalOptions.Unmarshal(data, pb)
	}
	return json.Unmarshal(data, obj)
}
