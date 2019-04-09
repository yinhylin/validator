package assert

import (
	"fmt"
	"reflect"
	"testing"
)

type TestNumberValidate struct {
	Foo int `json:"foo"`
}

func TestNumber_Validate(t *testing.T) {
	c := NewNumber().Field("Foo")
	s := &TestNumberValidate{}
	c.SetFieldType(reflect.TypeOf(s))

	ins := []interface{}{[]int{1}, "1a", "foo"}
	for _, i := range ins {
		valid, _ := c.Validate(i)
		if valid == nil {
			t.Errorf("non-number value %v(%T) should return validation", i, i)
		}
	}

	ins = []interface{}{1, 1.0}
	for _, i := range ins {
		valid, _ := c.Validate(i)
		if valid != nil {
			t.Errorf("number value %v(%T) return validation, want PASS", i, i)
		}
	}

	_, out := c.Validate("1.0")
	if fmt.Sprintf("%T", out) != "float64" {
		t.Errorf("string value should be converted to type float64")
	}
}
