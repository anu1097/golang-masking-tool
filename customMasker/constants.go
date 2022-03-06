package customMasker

type Mtype string

// Mask Types of format string
const (
	MPassword   Mtype = "password"
	MName       Mtype = "name"
	MAddress    Mtype = "addr"
	MEmail      Mtype = "email"
	MMobile     Mtype = "mobile"
	MTelephone  Mtype = "tel"
	MID         Mtype = "id"
	MCreditCard Mtype = "credit"
	MURL        Mtype = "url"
	MSecret     Mtype = "secret"
)

type MaskingCharacter string

// MaskingCharacter Types of placeholder string
const (
	PStar       MaskingCharacter = "*"
	PHyphen     MaskingCharacter = "-"
	PUnderscore MaskingCharacter = "_"
	PCross      MaskingCharacter = "x"
)
