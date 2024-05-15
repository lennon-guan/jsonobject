package jsonobject

import "testing"

var jsonBytes = []byte(`
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}
`)

func TestGetValue(t *testing.T) {
	v, err := Parse(jsonBytes)
	if err != nil {
		t.Fatalf("parse error: %+v", err)
	}
	if age, err := v.Field("age").Int(); err != nil || age != 37 {
		t.Fatalf("age is not 37")
	}
	if v.Field("fav.movie").String() != "Deer Hunter" {
		t.Fatal("fav.movie is not Deer Hunter")
	}
	if fv := v.Field("friends").Index(1).Field("first").String(); fv != "Roger" {
		t.Fatalf("friends[1].first is not Roger but %s", fv)
	}
}

func TestSetValue(t *testing.T) {
	v, err := Parse(jsonBytes)
	if err != nil {
		t.Fatalf("parse error: %+v", err)
	}
	v.Field("friends").Index(1).SetField("last", "ZGG")
	r, _ := v.MarshalIndent("", "  ")
	if v2, err := Parse(r); err != nil {
		t.Fatalf("parse error: %+v", err)
	} else {
		if v2.Field("friends").Index(1).Field("last").String() != "ZGG" {
			t.Fatal("friends[1].last is not ZGG")
		}
	}
}
