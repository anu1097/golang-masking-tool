package customMasker

// String mask input string of the mask type
//
// Example:
//
//   masker.String(masker.MName, "ggwhite")
//   masker.String(masker.MID, "A123456789")
//   masker.String(masker.MMobile, "0987987987")
func String(t Mtype, i string) string {
	return instance.String(t, i, i)
}

// Name mask the second letter and the third letter
//
// Example:
//   input: ABCD
//   output: A**D
func Name(i string) string {
	return instance.Name(i)
}

// ID mask last 4 digits of ID number
//
// Example:
//   input: A123456789
//   output: A12345****
func ID(i string) string {
	return instance.ID(i)
}

// Address keep first 6 letters, mask the rest
//
// Example:
//   input: 台北市內湖區內湖路一段737巷1號1樓
//   output: 台北市內湖區******
func Address(i string) string {
	return instance.Address(i)
}

// CreditCard mask 6 digits from the 7'th digit
//
// Example:
//   input1: 1234567890123456 (VISA, JCB, MasterCard)(len = 16)
//   output1: 123456******3456
//   input2: 123456789012345 (American Express)(len = 15)
//   output2: 123456******345
func CreditCard(i string) string {
	return instance.CreditCard(i)
}

// Email keep domain and the first 3 letters
//
// Example:
//   input: ggw.chang@gmail.com
//   output: ggw****@gmail.com
func Email(i string) string {
	return instance.Email(i)
}

// Mobile mask 3 digits from the 4'th digit
//
// Example:
//   input: 0987654321
//   output: 0987***321
func Mobile(i string) string {
	return instance.Mobile(i)
}

// Telephone remove "(", ")", " ", "-" chart, and mask last 4 digits of telephone number, format to "(??)????-????"
//
// Example:
//   input: 0227993078
//   output: (02)2799-****
func Telephone(i string) string {
	return instance.Telephone(i)
}

// Password always return "************"
func Password(i string) string {
	return instance.Password(i)
}

// SetMask sets Masking Character used my custom masker
func SetMask(mask string) {
	instance.mask = mask
}
