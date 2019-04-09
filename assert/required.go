package assert

import (
	"github.com/yinhylin/validator/validation"
)

type Required struct {
	Assert
	items []Interface
}

func NewRequired() *Required {
	r := BuildAssert(&Required{}).(*Required)
	r.Message("{{name}} is required")
	return r
}

func (r *Required) Item(items ...Interface) Interface {
	r.items = append(r.items, items...)
	return r
}

func (r *Required) GetItems() []Interface {
	return r.items
}

func (r *Required) HasItems() bool {
	return true
}

func (r *Required) GetItemType() int {
	return TypeSame
}

func (r *Required) Validate(input interface{}) (valid *validation.Validation, output interface{}) {
	if input == ErrNotFound {
		valid = r.BuildValidation(r.GetMessage(), M{"name": r.GetTagField()})
		return
	}
	output = input
	return
}
