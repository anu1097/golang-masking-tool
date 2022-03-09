package filter

import (
	"strings"

	"github.com/anu1097/golang-masking-tool/customMasker"
)

type fieldFilter struct {
	target   string
	maskType customMasker.Mtype
}

// Get a Custom Field Filter. Pass custom Masker type to define filter mechanism
func CustomFieldFilter(target string, maskType customMasker.Mtype) *fieldFilter {
	return &fieldFilter{
		target:   target,
		maskType: maskType,
	}
}

// Get a Field Filter.
func FieldFilter(target string) *fieldFilter {
	return &fieldFilter{
		target: target,
	}
}

func (x *fieldFilter) MaskString(s string) string {
	return customMaskerInstance.String(x.maskType, s, GetFilteredLabel())
}

func (x *fieldFilter) ReplaceString(s string) string {
	return s
}

func (x *fieldFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return x.target == fieldName
}

type fieldPrefixFilter struct {
	prefix   string
	maskType customMasker.Mtype
}

// Get a Field Prefix Filter.
func FieldPrefixFilter(prefix string) *fieldPrefixFilter {
	return &fieldPrefixFilter{
		prefix: prefix,
	}
}

// Get a Custom Field Prefix Filter. Pass custom Masker type to define filter mechanism
func CustomFieldPrefixFilter(prefix string, maskType customMasker.Mtype) *fieldPrefixFilter {
	return &fieldPrefixFilter{
		prefix:   prefix,
		maskType: maskType,
	}
}

func (x *fieldPrefixFilter) MaskString(s string) string {
	return customMaskerInstance.String(x.maskType, s, GetFilteredLabel())
}

func (x *fieldPrefixFilter) ReplaceString(s string) string {
	return s
}

func (x *fieldPrefixFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return strings.HasPrefix(fieldName, x.prefix)
}
