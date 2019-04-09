package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/yinhylin/validator/assert"
	"github.com/yinhylin/validator/validation"
)

const (
	GroupDefault = "default"
)

type Validator struct {
	collection      *assert.Collection
	validations     []*validation.Validation
	groups          []string
	continueOnError bool // 验证错误时是否需要继续验证其它字段
}

func New(collection assert.Interface) *Validator {
	c, ok := collection.(*assert.Collection)
	if !ok {
		panic(fmt.Errorf("validator assert should start with Collection"))
	}

	return &Validator{
		collection:      c,
		continueOnError: false,
	}
}

// ContinueOnError 设置验证错误是是否需要继续进行验证
func (v *Validator) ContinueOnError(val bool) {
	v.continueOnError = val
}

func (v *Validator) Validate(input interface{}, output interface{}, groups ...string) (valid []*validation.Validation) {
	if output == nil {
		panic(errors.New("validate: output is nil"))
	}
	ot := reflect.TypeOf(output)
	if ot.Kind() != reflect.Ptr {
		panic(errors.New("validate: output is non-pointer"))
	}

	v.reset()
	if len(groups) == 0 {
		groups = append(groups, GroupDefault)
	}
	v.groups = groups
	v.collection.SetFieldType(ot)

	var out interface{}
	defer func() {
		if r := recover(); r != nil {
			if r == assert.ErrAbortAll {
				valid = v.validations
				return
			}
			panic(r)
		}

		if len(v.validations) == 0 {
			data, err := json.Marshal(out)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(data, output)
			if err != nil {
				panic(err)
			}
		}
	}()
	v.BeforeValidate(v.collection, nil)
	out = v.validate(input, v.collection, nil)
	valid = v.validations
	return
}

func (v *Validator) validate(input interface{}, c assert.Interface, p assert.Interface) (output interface{}) {
	defer func() {
		if r := recover(); r != nil {
			if r == assert.ErrAbort {
				if !v.continueOnError && len(v.validations) > 0 {
					panic(assert.ErrAbortAll)
				}

				if p == nil || c.GetField() != p.GetField() {
					return
				}
			}
			panic(r)
		}
	}()

	valid, out := c.Validate(input)
	if valid != nil {
		v.validations = append(v.validations, valid)
		panic(assert.ErrAbort)
	}

	if c.HasItems() {
		switch c.GetItemType() {
		default:
			panic(fmt.Errorf("unknow ItemType - %v", c.GetItemType()))
		case assert.TypeSubordinate:
			out = v.processSubordinate(out, c, p)
		case assert.TypeSame:
			out = v.processSame(out, c, p)
		case assert.TypeGroup:
			out = v.processGroup(out, c, p)
		}
	}
	output = out
	return
}

func (v *Validator) BeforeValidate(child assert.Interface, parent assert.Interface) {
	if parent != nil && len(child.GetGroups()) == 0 {
		child.AddGroup(GroupDefault)
	}
	if parent != nil {
		// 子规则继承父级字段和分组
		child.ExtendField(parent.GetField())
		child.ExtendGroup(parent.GetGroups())
	}
	child.BeforeValidate(parent)

	// 将不符合请求的group的规则进行过滤
	if parent != nil {
		if !v.InGroup(v.groups, child.GetGroups()) {
			panic(assert.ErrAbort)
		}
	}
}

func (v *Validator) reset() {
	v.validations = []*validation.Validation{}
}

func (v *Validator) processSame(input interface{}, c assert.Interface, p assert.Interface) interface{} {
	for _, i := range c.GetItems() {
		v.BeforeValidate(i, c)
		i.SetFieldType(c.GetFieldType())
		input = v.validate(input, i, c)
	}
	return input
}

func (v *Validator) processSubordinate(input interface{}, c assert.Interface, p assert.Interface) interface{} {
	in := input.(map[string]interface{})
	for _, i := range c.GetItems() {
		v.BeforeValidate(i, c)
		if p != nil {
			sf, _ := c.GetFieldType().FieldByName(c.GetField())
			i.SetFieldType(sf.Type)
		} else {
			i.SetFieldType(c.GetFieldType())
		}
		field := i.GetTagField()

		val, ok := in[field]
		if ok {
			in[field] = v.validate(val, i, c)
		} else {
			in[field] = v.validate(assert.ErrNotFound, i, c)
		}
	}
	return in
}

func (v *Validator) processGroup(input interface{}, c assert.Interface, p assert.Interface) interface{} {
	in := make([]map[string]interface{}, len(input.([]map[string]interface{})))

	for k, item := range input.([]map[string]interface{}) {
		in[k] = make(map[string]interface{})
		sf, _ := c.GetFieldType().FieldByName(c.GetField())
		for _, i := range c.GetItems() {
			v.BeforeValidate(i, c)
			i.SetFieldType(sf.Type.Elem())
			field := i.GetTagField()

			val, ok := item[field]
			if ok {
				in[k][field] = v.validate(val, i, c)
			} else {
				in[k][field] = v.validate(assert.ErrNotFound, i, c)
			}
		}
	}
	return in
}

// InGroup 判断两个group中是否存在相同的值
func (v *Validator) InGroup(requestGroups []string, groups []string) bool {
	for _, rg := range requestGroups {
		for _, g := range groups {
			if rg == g {
				return true
			}
		}
	}
	return false
}
