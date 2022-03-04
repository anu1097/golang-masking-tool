package filter

import (
	"regexp"

	"github.com/anu1097/golang-mask-utility/masker"
)

const defaultPhoneRegex = `^((\+\d{1,3}(-| )?\(?\d\)?(-| )?\d{1,5})|(\(?\d{2,6}\)?))(-| )?(\d{3,4})(-| )?(\d{4})(( x| ext)\d{1,5}){0,1}$`
const defaultEmailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

type piiRegexFilter struct {
	RegexList []regexp.Regexp
	mtype     masker.Mtype
}

func PhoneFilter() *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(defaultPhoneRegex),
		},
	}
}

func CustomPhoneFilter(mtype masker.Mtype) *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(defaultPhoneRegex),
		},
		mtype: mtype,
	}
}

func EmailFilter() *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(defaultEmailRegex),
		},
	}
}

func CustomEmailFilter(mtype masker.Mtype) *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(defaultEmailRegex),
		},
		mtype: mtype,
	}
}

func CustomRegexFilter(regexPattern string) *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(regexPattern),
		},
	}
}

func CustomRegexFilterWithMType(regexPattern string, mtype masker.Mtype) *piiRegexFilter {
	return &piiRegexFilter{
		RegexList: []regexp.Regexp{
			*regexp.MustCompile(regexPattern),
		},
		mtype: mtype,
	}
}

func (x *piiRegexFilter) ReplaceString(s string) string {
	for _, p := range x.RegexList {
		s = p.ReplaceAllString(s, maskerInstance.String(x.mtype, s, GetFilteredLabel()))
	}
	return s
}

func (x *piiRegexFilter) MaskString(s string) string {
	return s
}

func (x *piiRegexFilter) ShouldMask(fieldName string, value interface{}, tag string) bool {
	return false
}
