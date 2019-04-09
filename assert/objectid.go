package assert

import (
	"encoding/hex"
	"reflect"

	"github.com/yinhylin/validator/validation"
)

type ObjectId struct {
	Assert
}

func NewObjectId() *ObjectId {
	n := BuildAssert(&ObjectId{}).(*ObjectId)
	n.Message("{{name}} should be of type mongodb ObjectId string")
	return n
}

func (n *ObjectId) Validate(input interface{}) (valid *validation.Validation, output interface{}) {
	// 空值检测应交给专门的空值判断规则进行处理
	if IsBlank(input) {
		output = input
		return
	}

	it := reflect.TypeOf(input)
	if it.Kind() != reflect.String || !IsObjectIdHex(input.(string)) {
		n.BuildValidation(n.GetMessage(), M{"name": n.GetTagField()})
		return
	}

	output = input
	return
}

func IsObjectIdHex(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
