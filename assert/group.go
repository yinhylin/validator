package assert

import (
	"github.com/yinhylin/validator/validation"
)

type Group struct {
	Assert
	items []Interface
}

func NewGroup() *Group {
	n := BuildAssert(&Group{}).(*Group)
	n.Message("{{name}} should be []interface{}")
	return n
}

func (c *Group) Item(items ...Interface) Interface {
	c.items = append(c.items, items...)
	return c
}

func (c *Group) GetItems() []Interface {
	return c.items
}

func (c *Group) HasItems() bool {
	return true
}

func (c *Group) GetItemType() int {
	return TypeGroup
}

func (c *Group) BeforeValidate(parent Interface) {
	// 非 Optional 的约束规则，默认为 Required
	for k, i := range c.items {
		c.items[k] = c.WrapRequired(i)
	}
}

func (c *Group) Validate(input interface{}) (valid *validation.Validation, output interface{}) {
	if _, ok := input.([]map[string]interface{}); !ok {
		valid = c.BuildValidation(c.GetMessage(), M{"name": c.GetTagField()})
		return
	}
	output = input

	return
}
