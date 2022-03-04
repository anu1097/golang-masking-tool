package masker

type Mtype string

// Maske Types of format string
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
