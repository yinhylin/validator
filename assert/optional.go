package assert

import (
	"github.com/yinhylin/validator/validation"
)

type Optional struct {
	Assert
	items []Interface
}

func NewOptional() *Optional {
	r := BuildAssert(&Optional{}).(*Optional)
	return r
}

func (o *Optional) Item(items ...Interface) Interface {
	o.items = append(o.items, items...)
	return o
}

func (o *Optional) GetItems() []Interface {
	return o.items
}

func (o *Optional) GetItemType() int {
	return TypeSame
}

func (o *Optional) HasItems() bool {
	return true
}

func (o *Optional) Validate(input interface{}) (valid *validation.Validation, output interface{}) {
	if input == ErrNotFound {
		o.Abort()
	}
	output = input
	return
}
