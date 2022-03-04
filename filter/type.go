package filter

import (
	"reflect"

	"github.com/anu1097/golang-mask-utility/masker"
)

type typeFilter struct {
	target   reflect.Type
	maskType masker.Mtype
}

func TypeFilter(t interface{}) *typeFilter {
	return &typeFilter{
		target: reflect.TypeOf(t),
	}
}

func CustomTypeFilter(t interface{}, maskType masker.Mtype) *typeFilter {
	return &typeFilter{
		target:   reflect.TypeOf(t),
		maskType: maskType,
	}
}

func (x *typeFilter) ReplaceString(s string) string { return s }

func (x *typeFilter) MaskString(s string) string {
	return maskerInstance.String(x.maskType, s, GetFilteredLabel())
}

func (x *typeFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return x.target == reflect.TypeOf(value)
}
