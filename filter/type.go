package filter

import (
	"reflect"

	"github.com/anu1097/golang-masking-tool/customMasker"
)

type typeFilter struct {
	target   reflect.Type
	maskType customMasker.Mtype
}

// Get Type Filter.
func TypeFilter(t interface{}) *typeFilter {
	return &typeFilter{
		target: reflect.TypeOf(t),
	}
}

// Get Custom Type Filter with custom masking type.
func CustomTypeFilter(t interface{}, maskType customMasker.Mtype) *typeFilter {
	return &typeFilter{
		target:   reflect.TypeOf(t),
		maskType: maskType,
	}
}

func (x *typeFilter) ReplaceString(s string) string { return s }

func (x *typeFilter) MaskString(s string) string {
	return customMaskerInstance.String(x.maskType, s, GetFilteredLabel())
}

func (x *typeFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return x.target == reflect.TypeOf(value)
}
