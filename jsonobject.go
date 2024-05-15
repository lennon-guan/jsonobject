package jsonobject

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	Invalid = iota
	Null
	Number
	String
	Boolean
	Array
	Object
)

type JsonValue struct {
	kind int
	val  any
}

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

func (v *JsonValue) transformValue(val any) any {
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

func (v *JsonValue) Field(name string) *JsonValue {
	if m, is := v.val.(map[string]any); !is {
		return &JsonValue{kind: Invalid, val: nil}
	} else if vv, found := m[name]; !found {
		return &JsonValue{kind: Invalid, val: nil}
	} else {
		return newValue(vv)
	}
}

func (v *JsonValue) SetField(name string, val any) *JsonValue {
	if m, is := v.val.(map[string]any); is {
		m[name] = v.transformValue(val)
	}
	return v
}

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

func (v *JsonValue) Index(index int) *JsonValue {
	if m, index := v.getRealIndex(index); index < 0 {
		return &JsonValue{kind: Invalid, val: nil}
	} else {
		return newValue(m[index])
	}
}

func (v *JsonValue) SetIndex(index int, val any) *JsonValue {
	if m, index := v.getRealIndex(index); index < 0 {
	} else {
		m[index] = v.transformValue(val)
	}
	return v
}

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

func (v *JsonValue) Each(f func(item *JsonValue, index int)) {
	if v.kind != Array {
		return
	}
	for i, item := range v.val.([]any) {
		f(newValue(item), i)
	}
}

func (v *JsonValue) EachPair(f func(value *JsonValue, key string)) {
	if v.kind != Object {
		return
	}
	for k, item := range v.val.(map[string]any) {
		f(newValue(item), k)
	}
}

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
func (v *JsonValue) String() string {
	return fmt.Sprint(v.val)
}

func (v *JsonValue) Int() (int, error) {
	if v.kind == Number {
		switch vv := v.val.(type) {
		case int:
			return vv, nil
		case float64:
			return int(vv), nil
		case float32:
			return int(vv), nil
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
	return 0, fmt.Errorf("not a number")
}

func (v *JsonValue) DefaultInt(defaultValue int) int {
	if r, e := v.Int(); e == nil {
		return r
	}
	return defaultValue
}

func (v *JsonValue) MustInt() int {
	if r, e := v.Int(); e == nil {
		return r
	} else {
		panic(e)
	}
}

func (v *JsonValue) Bool() (bool, error) {
	if v.kind == Boolean {
		switch vv := v.val.(type) {
		case bool:
			return vv, nil
		}
	}
	return false, fmt.Errorf("not a boolean")
}

func (v *JsonValue) DefaultBool(defaultValue bool) bool {
	if r, e := v.Bool(); e == nil {
		return r
	}
	return defaultValue
}

func (v *JsonValue) MustBool() bool {
	if r, e := v.Bool(); e == nil {
		return r
	} else {
		panic(e)
	}
}

func (v *JsonValue) Marshal() ([]byte, error) {
	return json.Marshal(v.val)
}

func (v *JsonValue) MarshalIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v.val, prefix, indent)
}

func (v *JsonValue) MarshalToString() string {
	if b, err := v.Marshal(); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (v *JsonValue) MarshalIndentToString(prefix, indent string) string {
	if b, err := v.MarshalIndent(prefix, indent); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func Parse(jsonBytes []byte) (*JsonValue, error) {
	var v any
	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		return nil, err
	}
	return newValue(v), nil
}

func NewObject() *JsonValue {
	return newValue(map[string]any{})
}

func NewArray() *JsonValue {
	return newValue([]any{})
}
