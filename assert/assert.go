package assert

import (
	"errors"
	"fmt"
	"github.com/yinhylin/validator/validation"
	"reflect"
	"sort"
	"strings"
)

var (
	TagName     = "json"
	ErrNotFound = errors.New("field not found")
	ErrAbort    = errors.New("abort field")
	ErrAbortAll = errors.New("abort all")
)

const (
	TypeSubordinate = iota // 待处理数据是当前数据的子数据
	TypeSame               // 待处理数据与当前数据相同，如：同时对一个数据验证非空和整数型
	TypeGroup              // 待处理数据是当前数据的子数据，且分为结构相同的多组
)

type M map[string]interface{}

type Interface interface {
	Field(field string) Interface
	Message(msg string) Interface
	Group(group ...string) Interface
	Payload(payload interface{}) Interface
	Item(items ...Interface) Interface

	GetField() string
	GetMessage() string
	GetGroups() []string
	GetPayload() interface{}
	GetItems() []Interface

	HasItems() bool
	GetFieldType() reflect.Type
	GetTagField() string

	RealAssert(p Interface)
	SetFieldType(p reflect.Type)
	AddGroup(group ...string)
	GetItemType() int

	ExtendField(field string)
	ExtendGroup(group []string)

	BeforeValidate(parent Interface)
	Validate(input interface{}) (valid *validation.Validation, output interface{})
}

type Assert struct {
	realAssert Interface
	field      string
	fieldType  reflect.Type
	message    string
	groups     []string
	payload    interface{}
}

func BuildAssert(p Interface) Interface {
	p.RealAssert(p)
	return p
}

func (a *Assert) Field(field string) Interface {
	a.field = field
	return a.realAssert
}

func (a *Assert) Message(msg string) Interface {
	a.message = msg
	return a.realAssert
}

func (a *Assert) Group(group ...string) Interface {
	a.groups = group
	return a.realAssert
}

func (a *Assert) Payload(payload interface{}) Interface {
	a.payload = payload
	return a.realAssert
}

func (a *Assert) Item(items ...Interface) Interface {
	panic(fmt.Sprintf(
		"the assert %s does not implement method: Item",
		reflect.TypeOf(a.realAssert).Elem().Name()),
	)
}

func (a *Assert) GetItemType() int {
	panic(fmt.Sprintf(
		"the assert %s does not implement method: GetItemType",
		reflect.TypeOf(a.realAssert).Elem().Name()),
	)
}

func (a *Assert) RealAssert(p Interface) {
	a.realAssert = p
}

func (a *Assert) GetField() string {
	return a.field
}

func (a *Assert) GetTagField() string {
	if a.field == "" {
		panic(fmt.Errorf("assert field not setted"))
	}

	fmt.Errorf("%v____________\n", a.GetFieldType().Kind().String())
	sf, ok := a.GetFieldType().FieldByName(a.field)
	if !ok {
		panic(fmt.Errorf("struct field - %s not found", a.field))
	}
	tag := strings.TrimSpace(sf.Tag.Get(TagName))
	if tag == "-" {
		panic(fmt.Errorf("struct field - %s is omitted in %s tag", a.field, TagName))
	}
	if n := strings.Index(tag, ","); n > -1 {
		tag = tag[0:n]
	}
	if tag == "" {
		tag = a.field
	}
	return tag
}

func (a *Assert) GetFieldType() reflect.Type {
	if a.fieldType.Kind() == reflect.Ptr {
		return a.fieldType.Elem()
	}
	return a.fieldType
}

func (a *Assert) SetFieldType(p reflect.Type) {
	a.fieldType = p
}

func (a *Assert) GetMessage() string {
	return a.message
}

func (a *Assert) GetGroups() []string {
	return a.groups
}

func (a *Assert) AddGroup(group ...string) {
	a.groups = append(a.groups, group...)

	// 分组名去重
	var (
		groups []string
		last   string
	)
	sort.Strings(a.groups)
	for _, g := range a.groups {
		if last != g {
			groups = append(groups, g)
		}
		last = g
	}
}

func (a *Assert) ExtendField(field string) {
	if a.field == "" {
		a.field = field
	}
}

func (a *Assert) ExtendGroup(group []string) {
	a.AddGroup(group...)
}

func (a *Assert) GetPayload() interface{} {
	return a.payload
}

func (a *Assert) BeforeValidate(parent Interface) {
}

func (a *Assert) Validate(input interface{}) (valid *validation.Validation, output interface{}) {
	panic(errors.New("implement the assert"))
}

func (a *Assert) HasItems() bool {
	return false
}

func (a *Assert) GetItems() []Interface {
	panic(fmt.Errorf(
		"the assert %s does not implement method: GetItems",
		reflect.TypeOf(a.realAssert).Elem().Name()),
	)
}

// WrapRequired 为规则包裹 Required 规则，例如 Collection 中的规则默认包裹 Required
func (a *Assert) WrapRequired(p Interface) Interface {
	if _, ok := p.(*Optional); ok {
		return p
	}

	return NewRequired().Field(p.GetField()).Group(p.GetGroups()...).Item(p)
}

func (a *Assert) BuildValidation(msg string, args M) *validation.Validation {
	for k, v := range args {
		msg = strings.Replace(msg, fmt.Sprintf("{{%s}}", k), fmt.Sprintf("%v", v), -1)
	}
	return &validation.Validation{Message: msg, Payload: a.payload}
}

func (a *Assert) Abort() {
	panic(ErrAbort)
}
