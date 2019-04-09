package assert

import (
	"reflect"
	"testing"
)

func TestGroup_BeforeValidate(t *testing.T) {
	c := NewGroup().Field("bar").Item(
		NewNotBlank().Field("foo"),
	)
	c.BeforeValidate(nil)
	if _, ok := c.GetItems()[0].(*Required); !ok {
		t.Errorf("non-optional assert should be requied assert default")
	}
}

type TestGroupValidate struct {
	Foo struct {
		Bar string `json:"bar"`
	} `json:"foo"`
}

func TestGroup_Validate(t *testing.T) {
	c := NewGroup().Field("Foo").Item(
		NewNotBlank().Field("Bar"),
	)
	s := &TestGroupValidate{}
	c.SetFieldType(reflect.TypeOf(s))

	ins := []interface{}{"foo", 123, nil}
	for _, i := range ins {
		valid, _ := c.Validate(i)
		if valid == nil {
			t.Errorf("non-collection value %v should return validation", i)
		}
	}

	in := []map[string]interface{}{{"bar": "bar"}}
	valid, out := c.Validate(in)
	if valid != nil {
		t.Errorf("%v should not return validation", in)
	}
	if out == nil {
		t.Errorf("validate should return input if PASS")
	}
}
