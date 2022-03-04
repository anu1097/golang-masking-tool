package filter

import (
	"strings"

	"github.com/anu1097/golang-mask-utility/masker"
)

type fieldFilter struct {
	target   string
	maskType masker.Mtype
}

func CustomFieldFilter(target string, maskType masker.Mtype) *fieldFilter {
	return &fieldFilter{
		target:   target,
		maskType: maskType,
	}
}

func FieldFilter(target string) *fieldFilter {
	return &fieldFilter{
		target: target,
	}
}

func (x *fieldFilter) MaskString(s string) string {
	return maskerInstance.String(x.maskType, s, GetFilteredLabel())
}

func (x *fieldFilter) ReplaceString(s string) string {
	return s
}

func (x *fieldFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return x.target == fieldName
}

type fieldPrefixFilter struct {
	prefix   string
	maskType masker.Mtype
}

func FieldPrefixFilter(prefix string) *fieldPrefixFilter {
	return &fieldPrefixFilter{
		prefix: prefix,
	}
}

func CustomFieldPrefixFilter(prefix string, maskType masker.Mtype) *fieldPrefixFilter {
	return &fieldPrefixFilter{
		prefix:   prefix,
		maskType: maskType,
	}
}

func (x *fieldPrefixFilter) MaskString(s string) string {
	return maskerInstance.String(x.maskType, s, GetFilteredLabel())
}

func (x *fieldPrefixFilter) ReplaceString(s string) string {
	return s
}

func (x *fieldPrefixFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return strings.HasPrefix(fieldName, x.prefix)
}
