package assert

import (
	"github.com/yinhylin/validator/validation"
	"reflect"
)

type NotBlank struct {
	Assert
}

func NewNotBlank() *NotBlank {
	n := BuildAssert(&NotBlank{}).(*NotBlank)
	n.Message("{{name}} should not be blank")
	return n
}

func (n *NotBlank) Validate(input interface{}) (valid *validation.Validation, output interface{}) {
	it := reflect.TypeOf(input)
	if input == nil || (it.Kind() == reflect.String && input.(string) == "") {
		valid = n.BuildValidation(n.GetMessage(), M{"name": n.GetTagField()})
		return
	}
	output = input
	return
}
