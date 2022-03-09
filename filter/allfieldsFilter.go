package filter

import "github.com/anu1097/golang-masking-tool/customMasker"

type allFieldsFilter struct {
	mtype customMasker.Mtype
}

// Get All Fields Filter.
func AllFieldFilter() *allFieldsFilter {
	return &allFieldsFilter{}
}

// Get All Fields Custom Filter with custom masking type
func CustomAllFieldFilter(mtype customMasker.Mtype) *allFieldsFilter {
	return &allFieldsFilter{
		mtype: mtype,
	}
}

func (x *allFieldsFilter) ReplaceString(s string) string {
	return GetFilteredLabel()
}

func (x *allFieldsFilter) MaskString(s string) string {
	return customMaskerInstance.String(x.mtype, s, GetFilteredLabel())
}

func (x *allFieldsFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return fieldName != ""
}
