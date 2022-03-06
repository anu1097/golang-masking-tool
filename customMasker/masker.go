// Package masker Provide mask format of Taiwan usually used(Name, Address, Email, ID ...etc.),
package customMasker

import (
	"math"
	"net/url"
	"strings"
)

type MaskerInterface interface {
	overlay(str string, overlay string, start int, end int) (overlayed string)
	String(t Mtype, i string) string
	Name(i string) string
	ID(i string) string
	Address(i string) string
	CreditCard(i string) string
	Email(i string) string
	Mobile(i string) string
	Telephone(i string) string
	Password(i string) string
	URL(i string) string
	UpdateMaskingCharacter(maskingCharacter MaskingCharacter)
}

// Masker is a instance to marshal masked string
type Masker struct {
	mask string
}

func strLoop(str string, length int) string {
	var mask string
	for i := 1; i <= length; i++ {
		mask += str
	}
	return mask
}

func (m *Masker) overlay(str string, overlay string, start int, end int) (overlayed string) {
	r := []rune(str)
	l := len([]rune(r))

	if l == 0 {
		return ""
	}

	if start < 0 {
		start = 0
	}
	if start > l {
		start = l
	}
	if end < 0 {
		end = 0
	}
	if end > l {
		end = l
	}
	if start > end {
		tmp := start
		start = end
		end = tmp
	}

	overlayed = ""
	overlayed += string(r[:start])
	overlayed += overlay
	overlayed += string(r[end:])
	return overlayed
}

// String mask input string of the mask type
//
// Example:
//
//   masker.String(masker.MName, "ggwhite")
//   masker.String(masker.MID, "A123456789")
//   masker.String(masker.MMobile, "0987987987")
func (m *Masker) String(t Mtype, i string, defaultFilteredString string) string {
	switch t {
	default:
		return defaultFilteredString
	case MPassword:
		return m.Password(i)
	case MName:
		return m.Name(i)
	case MAddress:
		return m.Address(i)
	case MEmail:
		return m.Email(i)
	case MMobile:
		return m.Mobile(i)
	case MID:
		return m.ID(i)
	case MTelephone:
		return m.Telephone(i)
	case MCreditCard:
		return m.CreditCard(i)
	case MURL:
		return m.URL(i)
	}
}

// Name mask the second letter and the third letter
//
// Example:
//   input: ABCD
//   output: A**D
func (m *Masker) Name(i string) string {
	l := len([]rune(i))

	if l == 0 {
		return ""
	}

	// if has space
	if strs := strings.Split(i, " "); len(strs) > 1 {
		tmp := make([]string, len(strs))
		for idx, str := range strs {
			tmp[idx] = m.Name(str)
		}
		return strings.Join(tmp, " ")
	}

	if l == 2 || l == 3 {
		return m.overlay(i, strLoop(m.mask, len("**")), 1, 2)
	}

	if l > 3 {
		return m.overlay(i, strLoop(m.mask, len("**")), 1, 3)
	}

	return strLoop(m.mask, len("**"))
}

// ID mask last 4 digits of ID number
//
// Example:
//   input: A123456789
//   output: A12345****
func (m *Masker) ID(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return m.overlay(i, strLoop(m.mask, len("****")), 6, 10)
}

// Address keep first 6 letters, mask the rest
//
// Example:
//   input: 台北市內湖區內湖路一段737巷1號1樓
//   output: 台北市內湖區******
func (m *Masker) Address(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	if l <= 6 {
		return strLoop(m.mask, len("******"))
	}
	return m.overlay(i, strLoop(m.mask, len("******")), 6, math.MaxInt64)
}

// CreditCard mask 6 digits from the 7'th digit
//
// Example:
//   input1: 1234567890123456 (VISA, JCB, MasterCard)(len = 16)
//   output1: 123456******3456
//   input2: 123456789012345` (American Express)(len = 15)
//   output2: 123456******345`
func (m *Masker) CreditCard(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return m.overlay(i, strLoop(m.mask, len("******")), 6, 12)
}

// Email keep domain and the first 3 letters
//
// Example:
//   input: ggw.chang@gmail.com
//   output: ggw****@gmail.com
func (m *Masker) Email(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	tmp := strings.Split(i, "@")

	switch len(tmp) {
	case 0:
		return ""
	case 1:
		return m.overlay(i, strLoop(m.mask, len("****")), 3, 7)
	}

	addr := tmp[0]
	domain := tmp[1]

	addr = m.overlay(addr, strLoop(m.mask, len("****")), 3, 7)

	return addr + "@" + domain
}

// Mobile mask 3 digits from the 4'th digit
//
// Example:
//   input: 0987654321
//   output: 0987***321
func (m *Masker) Mobile(i string) string {
	if len(i) == 0 {
		return ""
	}
	return m.overlay(i, strLoop(m.mask, len("***")), 4, 7)
}

// Telephone remove "(", ")", " ", "-" chart, and mask last 4 digits of telephone number, format to "(??)????-????"
//
// Example:
//   input: 0227993078
//   output: (02)2799-****
func (m *Masker) Telephone(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	i = strings.Replace(i, " ", "", -1)
	i = strings.Replace(i, "(", "", -1)
	i = strings.Replace(i, ")", "", -1)
	i = strings.Replace(i, "-", "", -1)

	l = len([]rune(i))

	if l != 10 && l != 8 {
		return i
	}

	ans := ""

	if l == 10 {
		ans += "("
		ans += i[:2]
		ans += ")"
		i = i[2:]
	}

	ans += i[:4]
	ans += "-"
	ans += "****"

	return ans
}

// Password always return "************"
func (m *Masker) Password(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return strLoop(m.mask, len("************"))
}

// URL mask the password part of the URL if exists
//
// Example:
//   input: http://admin:mysecretpassword@localhost:1234/uri
//   output:http://admin:xxxxx@localhost:1234/uri
func (m *Masker) URL(i string) string {
	u, err := url.Parse(i)
	if err != nil {
		return i
	}
	return u.Redacted()
}

// Update Masking Character Used by Custom Masker
func (m *Masker) UpdateMaskingCharacter(maskingCharacter MaskingCharacter) {
	m.mask = string(maskingCharacter)
}

// NewMasker create Masker
func NewMasker() *Masker {
	return &Masker{
		mask: string(PStar),
	}
}

var instance *Masker

func init() {
	instance = NewMasker()
}
