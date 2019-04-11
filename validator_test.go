package validator

import (
	"encoding/json"
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
	text := `{"mobile":"123","user_id":"212332","Password":123,"spouse":{"name":"foo","age": 2},"children":[{"name":"child1","age":"8"},{"name":"foo","age":2}]}`
	json.Unmarshal([]byte(text), &val)

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
		assert.NewNotBlank().Field("Other").Group("other"),
	)

	v := New(coll)
	valid := v.Validate(val, in)
	if len(valid) > 0 {
		t.Errorf(valid[0].Message)
	}
}
