package assert

import (
	"reflect"
	"testing"
)

type TestNotBlankInput struct {
	Foo string `json:"foo"`
}

func TestNotBlank_Validate(t *testing.T) {
	in := &TestNotBlankInput{}

	n := NewNotBlank().Field("Foo")
	n.SetFieldType(reflect.TypeOf(in))
	valid, _ := n.Validate(nil)
	if valid == nil {
		t.Error("notblank validate(nil) want got validation")
	}

	val := "foo"
	valid, out := n.Validate(val)
	if valid != nil || out != val {
		t.Errorf("notblank validate(%s) should be pass", val)
	}
}
