package filter

import "github.com/anu1097/golang-mask-utility/masker"

type allFieldsFilter struct {
	mtype masker.Mtype
}

func AllFieldFilter() *allFieldsFilter {
	return &allFieldsFilter{}
}

func CustomAllFieldFilter(mtype masker.Mtype) *allFieldsFilter {
	return &allFieldsFilter{
		mtype: mtype,
	}
}

func (x *allFieldsFilter) ReplaceString(s string) string {
	return s
}

func (x *allFieldsFilter) MaskString(s string) string {
	return s
}

func (x *allFieldsFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return fieldName != ""
}
