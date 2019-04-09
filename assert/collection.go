package assert

import (
	"github.com/yinhylin/validator/validation"
)

type Collection struct {
	Assert
	items []Interface
}

func NewCollection() *Collection {
	n := BuildAssert(&Collection{}).(*Collection)
	n.Message("{{name}} should be map[string]interface{}")
	return n
}

func (c *Collection) Item(items ...Interface) Interface {
	c.items = append(c.items, items...)
	return c
}

func (c *Collection) GetItems() []Interface {
	return c.items
}

func (c *Collection) HasItems() bool {
	return true
}

func (c *Collection) GetItemType() int {
	return TypeSubordinate
}

func (c *Collection) BeforeValidate(parent Interface) {
	// 非 Optional 的约束规则，默认为 Required
	for k, i := range c.items {
		c.items[k] = c.WrapRequired(i)
	}
}

func (c *Collection) Validate(input interface{}) (valid *validation.Validation, output interface{}) {
	if _, ok := input.(map[string]interface{}); !ok {
		if c.GetField() == "" {
			valid = c.BuildValidation(c.GetMessage(), M{"name": "collection"})
		} else {
			valid = c.BuildValidation(c.GetMessage(), M{"name": c.GetTagField()})
		}
		return
	}
	output = input
	return
}
