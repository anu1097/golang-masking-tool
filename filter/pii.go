package filter

import (
	"regexp"

	"github.com/anu1097/golang-masking-tool/customMasker"
)

const defaultPhoneRegex = `^((\+\d{1,3}(-| )?\(?\d\)?(-| )?\d{1,5})|(\(?\d{2,6}\)?))(-| )?(\d{3,4})(-| )?(\d{4})(( x| ext)\d{1,5}){0,1}$`
const defaultEmailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

type piiRegexFilter struct {
	RegexList []regexp.Regexp
	mtype     customMasker.Mtype
}

// Get Phone Filter.
func PhoneFilter() *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(defaultPhoneRegex),
		},
	}
}

// Get Custom Phone Filter with custom masking type.
func CustomPhoneFilter(mtype customMasker.Mtype) *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(defaultPhoneRegex),
		},
		mtype: mtype,
	}
}

// Get Email Filter.
func EmailFilter() *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(defaultEmailRegex),
		},
	}
}

// Get Custom Email Filter with custom masking type.
func CustomEmailFilter(mtype customMasker.Mtype) *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(defaultEmailRegex),
		},
		mtype: mtype,
	}
}

// Get Custom Regex Filter.
func CustomRegexFilter(regexPattern string) *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(regexPattern),
		},
	}
}

// Get Custom Regex Filter with custom masking type
func CustomRegexFilterWithMType(regexPattern string, mtype customMasker.Mtype) *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(regexPattern),
		},
		mtype: mtype,
	}
}

func (x *piiRegexFilter) ReplaceString(s string) string {
	for _, p := range x.RegexList {
		s = p.ReplaceAllString(s, customMaskerInstance.String(x.mtype, s, GetFilteredLabel()))
	}
	return s
}

func (x *piiRegexFilter) MaskString(s string) string {
	return s
}

func (x *piiRegexFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return false
}
