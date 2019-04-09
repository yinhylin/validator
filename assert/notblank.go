package assert

import (
	"reflect"

	"github.com/yinhylin/validator/validation"
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
	if IsBlank(input) {
		valid = n.BuildValidation(n.GetMessage(), M{"name": n.GetTagField()})
		return
	}
	output = input
	return
}

func IsBlank(input interface{}) bool {
	it := reflect.TypeOf(input)
	return input == nil || (it.Kind() == reflect.String && input.(string) == "")
}
