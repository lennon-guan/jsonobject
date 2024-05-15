package jsonobject

import (
	"encoding/json"
	"fmt"
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

func (v *JsonValue) Field(name string) *JsonValue {
	if m, is := v.val.(map[string]any); !is {
		return &JsonValue{kind: Invalid, val: nil}
	} else if vv, found := m[name]; !found {
		return &JsonValue{kind: Invalid, val: nil}
	} else {
		return newValue(vv)
	}
}

func (v *JsonValue) SetField(name string, val any) bool {
	if m, is := v.val.(map[string]any); !is {
		return false
	} else {
		m[name] = val
	}
	return true
}

func (v *JsonValue) Index(index int) *JsonValue {
	if index < 0 {
		return &JsonValue{kind: Invalid, val: nil}
	} else if m, is := v.val.([]any); !is {
		return &JsonValue{kind: Invalid, val: nil}
	} else if index >= len(m) {
		return &JsonValue{kind: Invalid, val: nil}
	} else {
		return newValue(m[index])
	}
}

func (v *JsonValue) SetIndex(index int, val any) bool {
	if index < 0 {
		return false
	} else if m, is := v.val.([]any); !is {
		return false
	} else if index >= len(m) {
		return false
	} else {
		m[index] = val
	}
	return true
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

func (v *JsonValue) Marshal() ([]byte, error) {
	return json.Marshal(v.val)
}

func (v *JsonValue) MarshalIndent(prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v.val, prefix, indent)
}

func Parse(jsonBytes []byte) (*JsonValue, error) {
	var v any
	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		return nil, err
	}
	return newValue(v), nil
}
