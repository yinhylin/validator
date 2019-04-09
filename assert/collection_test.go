package assert

import (
	"testing"
)

func TestCollection_BeforeValidate(t *testing.T) {
	c := NewCollection().Item(
		NewNotBlank().Field("foo"),
	)
	c.BeforeValidate(nil)
	if _, ok := c.GetItems()[0].(*Required); !ok {
		t.Errorf("non-optional assert should be requied assert default")
	}
}

func TestCollection_Validate(t *testing.T) {
	c := NewCollection().Item(
		NewNotBlank().Field("foo"),
	)
	for _, i := range []interface{}{"123", nil, 123, []int{1, 2}} {
		valid, _ := c.Validate(i)
		if valid == nil {
			t.Errorf("non-collection value %v should return validation", i)
		}
	}

	v := make(map[string]interface{})
	v["foo"] = "bar"
	valid, out := c.Validate(v)
	if valid != nil {
		t.Errorf("%v should not return validation", v)
	}
	if v["foo"] != out.(map[string]interface{})["foo"] {
		t.Errorf("validate should return input if PASS")
	}
}
