package filter

import "github.com/anu1097/golang-mask-utility/masker"

type tagFilter struct {
	SecureTags []string
	maskType   masker.Mtype
}

var tagKey = "mask"

func SetTagKey(tag string) {
	tagKey = tag
}

func GetTagKey() string {
	return tagKey
}

func TagFilter(tags ...masker.Mtype) *tagFilter {
	if len(tags) == 0 {
		tags = []masker.Mtype{masker.MSecret}
	}
	var secureTags []string
	for _, tag := range tags {
		secureTags = append(secureTags, string(tag))
	}
	return &tagFilter{
		SecureTags: secureTags,
	}
}

func (x *tagFilter) ReplaceString(s string) string { return s }

func (x *tagFilter) MaskString(s string) string {
	return maskerInstance.String(x.maskType, s, GetFilteredLabel())
}

func (x *tagFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	for i := range x.SecureTags {
		if x.SecureTags[i] == tag {
			x.maskType = masker.Mtype(tag)
			return true
		}
	}
	return false
}
