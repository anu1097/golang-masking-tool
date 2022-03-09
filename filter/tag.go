package filter

import "github.com/anu1097/golang-masking-tool/customMasker"

type tagFilter struct {
	SecureTags []string
	maskType   customMasker.Mtype
}

var tagKey = "mask"

func SetTagKey(tag string) {
	tagKey = tag
}

func GetTagKey() string {
	return tagKey
}

// Get Tag Filter. Need to pass custom masker type string.
//
// Example:
//   input: secret
//   output: [filtered]
func TagFilter(tags ...customMasker.Mtype) *tagFilter {
	if len(tags) == 0 {
		tags = []customMasker.Mtype{customMasker.MSecret}
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
	return customMaskerInstance.String(x.maskType, s, GetFilteredLabel())
}

func (x *tagFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	for i := range x.SecureTags {
		if x.SecureTags[i] == tag {
			x.maskType = customMasker.Mtype(tag)
			return true
		}
	}
	return false
}
