package validator

import (
	"testing"

	"github.com/yinhylin/validator/assert"
)

type TestPerson struct {
	UserId   string `json:"user_id"`
	Mobile   string `json:"mobile"`
	Password int
	Spouse   struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	} `json:"spouse"`
	Children []TestChild `json:"children"`
}

type TestChild struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestValidate(t *testing.T) {
	in := &TestPerson{}

	val := make(map[string]interface{})
	val["mobile"] = "123"
	val["user_id"] = "212332"
	val["Password"] = 123
	val["spouse"] = map[string]interface{}{"name": "foo", "age": 2}
	val["children"] = []map[string]interface{}{{"name": "child1", "age": "8"}, {"name": "child2", "age": 3}}

	coll := assert.NewCollection().Item(
		assert.NewNotBlank().Field("UserId"),
		assert.NewNotBlank().Field("Password"),
		assert.NewOptional().Field("Mobile").Item(
			assert.NewNotBlank(),
		),
		assert.NewCollection().Field("Spouse").Item(
			assert.NewNotBlank().Field("Name"),
		),
		assert.NewGroup().Field("Children").Item(
			assert.NewNotBlank().Field("Name"),
			assert.NewRequired().Field("Age").Item(
				assert.NewNotBlank(),
				assert.NewNumber(),
			),
		),
	)

	v := New(coll)
	valid := v.Validate(val, in)
	if len(valid) > 0 {
		t.Errorf(valid[0].Message)
	}
}
