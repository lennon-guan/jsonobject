package jsonobject

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
)

// Constants representing the type of JSON values.
const (
	Invalid = iota
	Nil
	Number
	String
	Boolean
	Array
	Object
)

// JsonValue represents a JSON value of any type.
type JsonValue struct {
	kind int
	val  any
}

// newValue creates a new JsonValue based on the type of the provided value.
func newValue(v any) *JsonValue {
	switch v.(type) {
	case float64:
		return &JsonValue{kind: Number, val: v}
	case float32:
		return &JsonValue{kind: Number, val: v}
	case int:
		return &JsonValue{kind: Number, val: v}
	case int8:
		return &JsonValue{kind: Number, val: v}
	case int16:
		return &JsonValue{kind: Number, val: v}
	case int32:
		return &JsonValue{kind: Number, val: v}
	case int64:
		return &JsonValue{kind: Number, val: v}
	case uint:
		return &JsonValue{kind: Number, val: v}
	case uint8:
		return &JsonValue{kind: Number, val: v}
	case uint16:
		return &JsonValue{kind: Number, val: v}
	case uint32:
		return &JsonValue{kind: Number, val: v}
	case uint64:
		return &JsonValue{kind: Number, val: v}
	case string:
		return &JsonValue{kind: String, val: v}
	case bool:
		return &JsonValue{kind: Boolean, val: v}
	case []any:
		return &JsonValue{kind: Array, val: v}
	case map[string]any:
		return &JsonValue{kind: Object, val: v}
	default:
		return &JsonValue{kind: Invalid, val: nil}
	}
}

// transformValue transforms the given value into its basic equivalent.
func (v *JsonValue) transformValue(val any) any {
	if val == nil {
		return &JsonValue{kind: Nil, val: nil}
	}
	switch value := val.(type) {
	case *JsonValue:
		return value.val
	case int:
		return value
	case int8:
		return value
	case int16:
		return value
	case int32:
		return value
	case int64:
		return value
	case uint:
		return value
	case uint8:
		return value
	case uint16:
		return value
	case uint32:
		return value
	case uint64:
		return value
	case float32:
		return value
	case float64:
		return value
	case bool:
		return value
	case string:
		return value
	case []string:
		rv := make([]any, len(value))
		for i := range value {
			rv[i] = value[i]
		}
		return rv
	case []int:
		rv := make([]any, len(value))
		for i := range value {
			rv[i] = value[i]
		}
		return rv
	case map[string]string:
		rv := make(map[string]any, len(value))
		for k, v := range value {
			rv[k] = v
		}
		return rv
	case map[string]int:
		rv := make(map[string]any, len(value))
		for k, v := range value {
			rv[k] = v
		}
		return rv
	default:
		r := reflect.ValueOf(val)
		switch r.Kind() {
		case reflect.Map:
			rv := make(map[string]any, r.Len())
			for _, k := range r.MapKeys() {
				rv[k.String()] = v.transformValue(r.MapIndex(k).Interface())
			}
			return rv
		case reflect.Slice, reflect.Array:
			rv := make([]any, r.Len())
			for i := range rv {
				rv[i] = v.transformValue(r.Index(i).Interface())
			}
			return rv
		}
	}
	return val
}

// Field returns the JsonValue associated with the given field name.
func (v *JsonValue) Field(name string) *JsonValue {
	if m, is := v.val.(map[string]any); !is {
		return &JsonValue{kind: Invalid, val: nil}
	} else if vv, found := m[name]; !found {
		return &JsonValue{kind: Invalid, val: nil}
	} else {
		return newValue(vv)
	}
}

// SetField sets the field with the provided name to the specified value.
func (v *JsonValue) SetField(name string, val any) *JsonValue {
	if m, is := v.val.(map[string]any); is {
		m[name] = v.transformValue(val)
	}
	return v
}

// getRealIndex calculates the real index in an array.
func (v *JsonValue) getRealIndex(index int) ([]any, int) {
	if m, is := v.val.([]any); !is {
		return nil, -1
	} else if n := len(m); index >= n {
		return nil, -1
	} else {
		if index >= 0 {
			return m, index
		}
		index += n
		if index >= 0 {
			return m, index
		}
	}
	return nil, -1
}

// Index returns the JsonValue at the specified index in an array.
func (v *JsonValue) Index(index int) *JsonValue {
	if m, index := v.getRealIndex(index); index < 0 {
		return &JsonValue{kind: Invalid, val: nil}
	} else {
		return newValue(m[index])
	}
}

// SetIndex sets the array element at the specified index to the provided value.
func (v *JsonValue) SetIndex(index int, val any) *JsonValue {
	if m, index := v.getRealIndex(index); index >= 0 {
		m[index] = v.transformValue(val)
	}
	return v
}

// Append appends the provided value to the array.
func (v *JsonValue) Append(val any) *JsonValue {
	if v.kind != Array {
		return v
	}
	a, is := v.val.([]any)
	if !is {
		return v
	}
	v.val = append(a, v.transformValue(val))
	return v
}

// Each iterates over array elements and applies the provided function on each element.
func (v *JsonValue) Each(f func(item *JsonValue, index int)) {
	if v.kind != Array {
		return
	}
	for i, item := range v.val.([]any) {
		f(newValue(item), i)
	}
}

// EachPair iterates over key-value pairs in an object and applies the provided function on each pair.
func (v *JsonValue) EachPair(f func(value *JsonValue, key string)) {
	if v.kind != Object {
		return
	}
	for k, item := range v.val.(map[string]any) {
		f(newValue(item), k)
	}
}

// Size returns the size of the array or object represented by the JsonValue.
func (v *JsonValue) Size() int {
	switch v.kind {
	case Array:
		return len(v.val.([]any))
	case Object:
		return len(v.val.(map[string]any))
	default:
		return 0
	}
}

// IsValid checks if the JsonValue is valid.
func (v *JsonValue) IsValid() bool {
	return v.kind != Invalid
}

// IsNil checks if the JsonValue is nil.
func (v *JsonValue) IsNil() bool {
	return v.kind == Nil
}

// String returns the string representation of the JsonValue.
func (v *JsonValue) String() string {
	return fmt.Sprint(v.val)
}

// Float returns the value as a float64, if possible.
func (v *JsonValue) Float() (float64, error) {
	if v.kind == Number {
		switch vv := v.val.(type) {
		case float64:
			return vv, nil
		case float32:
			return float64(vv), nil
		case int:
			return float64(vv), nil
		case int8:
			return float64(vv), nil
		case int16:
			return float64(vv), nil
		case int32:
			return float64(vv), nil
		case int64:
			return float64(vv), nil
		case uint:
			return float64(vv), nil
		case uint8:
			return float64(vv), nil
		case uint16:
			return float64(vv), nil
		case uint32:
			return float64(vv), nil
		case uint64:
			return float64(vv), nil
		}
	}
	return 0, fmt.Errorf("not a number")
}

// DefaultFloat returns the float value or a default if it can't be converted.
func (v *JsonValue) DefaultFloat(defaultValue float64) float64 {
	if r, e := v.Float(); e == nil {
		return r
	}
	return defaultValue
}

// MustFloat returns the float value or panics if it can't be converted.
func (v *JsonValue) MustFloat() float64 {
	if r, e := v.Float(); e == nil {
		return r
	} else {
		panic(e)
	}
}

// Int returns the value as an int, if possible.
func (v *JsonValue) Int() (int, error) {
	if v.kind == Number {
		switch vv := v.val.(type) {
		case int:
			return vv, nil
		case float64:
			if vv == math.Floor(vv) {
				return int(vv), nil
			}
		case float32:
			if v64 := float64(vv); v64 == math.Floor(v64) {
				return int(vv), nil
			}
		case int8:
			return int(vv), nil
		case int16:
			return int(vv), nil
		case int32:
			return int(vv), nil
		case int64:
			return int(vv), nil
		case uint:
			return int(vv), nil
		case uint8:
			return int(vv), nil
		case uint16:
			return int(vv), nil
		case uint32:
			return int(vv), nil
		case uint64:
			return int(vv), nil
		}
	}
	return 0, fmt.Errorf("not an integer")
}

// DefaultInt returns the int value or a default if it can't be converted.
func (v *JsonValue) DefaultInt(defaultValue int) int {
	if r, e := v.Int(); e == nil {
		return r
	}
	return defaultValue
}

// MustInt returns the int value or panics if it can't be converted.
func (v *JsonValue) MustInt() int {
	if r, e := v.Int(); e == nil {
		return r
	} else {
		panic(e)
	}
}

// Bool returns the value as a bool, if possible.
func (v *JsonValue) Bool() (bool, error) {
	if v.kind == Boolean {
		switch vv := v.val.(type) {
		case bool:
			return vv, nil
		}
	}
	return false, fmt.Errorf("not a boolean")
}

// DefaultBool returns the boolean value or a default if it can't be converted.
func (v *JsonValue) DefaultBool(defaultValue bool) bool {
	if r, e := v.Bool(); e == nil {
		return r
	}
	return defaultValue
}

// MustBool returns the boolean value or panics if it can't be converted.
func (v *JsonValue) MustBool() bool {
	if r, e := v.Bool(); e == nil {
		return r
	} else {
		panic(e)
	}
}

// Str returns the value as a string, if possible.
func (v *JsonValue) Str() (string, error) {
	if v.kind == String {
		switch vv := v.val.(type) {
		case string:
			return vv, nil
		}
	}
	return "", fmt.Errorf("not a string")
}

// DefaultStr returns the string value or a default if it can't be converted.
func (v *JsonValue) DefaultStr(defaultValue string) string {
	if r, e := v.Str(); e == nil {
		return r
	}
	return defaultValue
}

// MustStr returns the string value or panics if it can't be converted.
func (v *JsonValue) MustStr() string {
	if r, e := v.Str(); e == nil {
		return r
	} else {
		panic(e)
	}
}

// Marshal marshals the JsonValue into a JSON byte slice.
func (v *JsonValue) Marshal() ([]byte, error) {
	return json.Marshal(v.val)
}

// MarshalIndent marshals the JsonValue into a JSON byte slice with indentation.
func (v *JsonValue) MarshalIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v.val, prefix, indent)
}

// MarshalToString marshals the JsonValue into a JSON string.
func (v *JsonValue) MarshalToString() string {
	if b, err := v.Marshal(); err != nil {
		return ""
	} else {
		return string(b)
	}
}

// MarshalIndentToString marshals the JsonValue into a JSON string with indentation.
func (v *JsonValue) MarshalIndentToString(prefix, indent string) string {
	if b, err := v.MarshalIndent(prefix, indent); err != nil {
		return ""
	} else {
		return string(b)
	}
}

// Parse parses a JSON byte slice into a JsonValue.
func Parse(jsonBytes []byte) (*JsonValue, error) {
	var v any
	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		return nil, err
	}
	return newValue(v), nil
}

// NewObject creates a new JsonValue representing an empty JSON object.
func NewObject() *JsonValue {
	return newValue(map[string]any{})
}

// NewArray creates a new JsonValue representing an empty JSON array.
func NewArray() *JsonValue {
	return newValue([]any{})
}
