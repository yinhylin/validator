package assert

import (
	"reflect"
	"strconv"

	"github.com/yinhylin/validator/validation"
)

type Number struct {
	Assert
}

func NewNumber() *Number {
	n := BuildAssert(&Number{}).(*Number)
	n.Message("{{name}} should be of type number")
	return n
}

func (n *Number) Validate(input interface{}) (valid *validation.Validation, output interface{}) {
	// 空值检测应交给专门的空值判断规则进行处理
	if IsBlank(input) {
		output = input
		return
	}

	it := reflect.TypeOf(input)
	switch it.Kind() {
	default:
		valid = n.BuildValidation(n.GetMessage(), M{"name": n.GetTagField()})
		return
	case reflect.String:
		out, err := strconv.ParseFloat(input.(string), 64)
		if err != nil {
			valid = n.BuildValidation(n.GetMessage(), M{"name": n.GetTagField()})
			return
		}
		output = out
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		output = input
	}
	return
}
