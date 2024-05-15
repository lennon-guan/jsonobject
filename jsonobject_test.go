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
	if v.Field("age").DefaultInt(-1) != 37 {
		t.Fatalf("age is not 37")
	}
	if v.Field("fav.movie").String() != "Deer Hunter" {
		t.Fatal("fav.movie is not Deer Hunter")
	}
	if fv := v.Field("friends").Index(1).Field("first").String(); fv != "Roger" {
		t.Fatalf("friends[1].first is not Roger but %s", fv)
	}
	if fv := v.Field("friends").Index(-1).Field("first").String(); fv != "Jane" {
		t.Fatalf("friends[-1].first is not Jane but %s", fv)
	}
}

func TestSetValue(t *testing.T) {
	v, err := Parse(jsonBytes)
	if err != nil {
		t.Fatalf("parse error: %+v", err)
	}
	v.Field("friends").Index(1).SetField("last", "ZGG")
	v.Field("friends").Each(func(item *JsonValue, index int) {
		item.SetField("age", 99)
	})
	r, _ := v.MarshalIndent("", "  ")
	v2, err := Parse(r)
	if err != nil {
		t.Fatalf("parse error: %+v", err)
	}
	if v2.Field("friends").Index(1).Field("last").String() != "ZGG" {
		t.Fatal("friends[1].last is not ZGG")
	}
	if v2.Field("friends").Size() != 3 {
		t.Fatal("friends size is not 3")
	}
	for i := 0; i < v2.Field("friends").Size(); i++ {
		if v2.Field("friends").Index(i).Field("age").DefaultInt(-1) != 99 {
			t.Fatal("friends[i].age is not 99")
		}
	}
}

func TestNewObject(t *testing.T) {
	r := NewObject()
	r.SetField("name", "ZGG")
	r.SetField("level", "SSR")
	r.SetField("tags", []string{"数据分析", "Web服务", "运维工具"})
	r.SetField("comments",
		NewArray().
			Append(NewObject().SetField("id", 1).SetField("content", "非常好用")).
			Append(NewObject().SetField("id", 2).SetField("content", "相见恨晚")),
	)
	t.Log(r.MarshalIndentToString("", "  "))
}
