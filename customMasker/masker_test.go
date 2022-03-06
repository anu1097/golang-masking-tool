package customMasker

import (
	"reflect"
	"testing"
)

func TestMasker_overlay(t *testing.T) {
	type args struct {
		str     string
		overlay string
		start   int
		end     int
	}
	tests := []struct {
		name          string
		m             *Masker
		args          args
		wantOverlayed string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				str:     "",
				overlay: "*",
				start:   0,
				end:     0,
			},
			wantOverlayed: "",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				str:     "abcdefg",
				overlay: "***",
				start:   1,
				end:     5,
			},
			wantOverlayed: "a***fg",
		},
		{
			name: "Start Less Than 0",
			m:    NewMasker(),
			args: args{
				str:     "abcdefg",
				overlay: "***",
				start:   -1,
				end:     5,
			},
			wantOverlayed: "***fg",
		},
		{
			name: "Start Greater Than Length",
			m:    NewMasker(),
			args: args{
				str:     "abcdefg",
				overlay: "***",
				start:   30,
				end:     31,
			},
			wantOverlayed: "abcdefg***",
		},
		{
			name: "End Less Than 0",
			m:    NewMasker(),
			args: args{
				str:     "abcdefg",
				overlay: "***",
				start:   1,
				end:     -5,
			},
			wantOverlayed: "***bcdefg",
		},
		{
			name: "Start Less Than End",
			m:    NewMasker(),
			args: args{
				str:     "abcdefg",
				overlay: "***",
				start:   5,
				end:     1,
			},
			wantOverlayed: "a***fg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Masker{}
			if gotOverlayed := m.overlay(tt.args.str, tt.args.overlay, tt.args.start, tt.args.end); gotOverlayed != tt.wantOverlayed {
				t.Errorf("Masker.overlay() = %v, want %v", gotOverlayed, tt.wantOverlayed)
			}
		})
	}
}

func TestMasker_String(t *testing.T) {
	type args struct {
		t Mtype
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Error Mask Type",
			m:    NewMasker(),
			args: args{
				t: Mtype(""),
				i: "abcdefg",
			},
			want: "abcdefg",
		},
		{
			name: "Password",
			m:    NewMasker(),
			args: args{
				t: MPassword,
				i: "ggwhite",
			},
			want: "************",
		},
		{
			name: "Name",
			m:    NewMasker(),
			args: args{
				t: MName,
				i: "ggwhite",
			},
			want: "g**hite",
		},
		{
			name: "Address",
			m:    NewMasker(),
			args: args{
				t: MAddress,
				i: "1 AB Road Paradise",
			},
			want: "1 AB R******",
		},
		{
			name: "Email",
			m:    NewMasker(),
			args: args{
				t: MEmail,
				i: "dummy@dummy.com",
			},
			want: "dum****@dummy.com",
		},
		{
			name: "Mobile",
			m:    NewMasker(),
			args: args{
				t: MMobile,
				i: "0978978978",
			},
			want: "0978***978",
		},
		{
			name: "ID",
			m:    NewMasker(),
			args: args{
				t: MID,
				i: "A123456789",
			},
			want: "A12345****",
		},
		{
			name: "Telephone",
			m:    NewMasker(),
			args: args{
				t: MTelephone,
				i: "0288889999",
			},
			want: "(02)8888-****",
		},
		{
			name: "CreditCard",
			m:    NewMasker(),
			args: args{
				t: MCreditCard,
				i: "1234567890123456",
			},
			want: "123456******3456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMasker()
			if got := m.String(tt.args.t, tt.args.i, tt.args.i); got != tt.want {
				t.Errorf("Masker.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Name(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Length 1",
			m:    NewMasker(),
			args: args{
				i: "A",
			},
			want: "**",
		},
		{
			name: "Length 2",
			m:    NewMasker(),
			args: args{
				i: "AB",
			},
			want: "A**",
		},
		{
			name: "Length 3",
			m:    NewMasker(),
			args: args{
				i: "ABC",
			},
			want: "A**C",
		},
		{
			name: "Length 4",
			m:    NewMasker(),
			args: args{
				i: "ABCD",
			},
			want: "A**D",
		},
		{
			name: "Length 5",
			m:    NewMasker(),
			args: args{
				i: "ABCDE",
			},
			want: "A**DE",
		},
		{
			name: "Length 6",
			m:    NewMasker(),
			args: args{
				i: "ABCDEF",
			},
			want: "A**DEF",
		},
		{
			name: "English Full Name",
			m:    NewMasker(),
			args: args{
				i: "Jorge Marry",
			},
			want: "J**ge M**ry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Name(tt.args.i); got != tt.want {
				t.Errorf("Masker.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_ID(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				i: "A123456789",
			},
			want: "A12345****",
		},
		{
			name: "Length Less Than 6",
			m:    NewMasker(),
			args: args{
				i: "A12",
			},
			want: "A12****",
		},
		{
			name: "Length Less Than 6",
			m:    NewMasker(),
			args: args{
				i: "A",
			},
			want: "A****",
		},
		{
			name: "Length Between 6 and 10",
			m:    NewMasker(),
			args: args{
				i: "A123456",
			},
			want: "A12345****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.ID(tt.args.i); got != tt.want {
				t.Errorf("Masker.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Address(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Long Address",
			m:    NewMasker(),
			args: args{
				i: "1 AB Road Paradise",
			},
			want: "1 AB R******",
		},
		{
			name: "Length Less Than 6",
			m:    NewMasker(),
			args: args{
				i: "ABC",
			},
			want: "******",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Address(tt.args.i); got != tt.want {
				t.Errorf("Masker.Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_CreditCard(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "VISA JCB MasterCard",
			m:    NewMasker(),
			args: args{
				i: "1234567890123456",
			},
			want: "123456******3456",
		},
		{
			name: "American Express",
			m:    NewMasker(),
			args: args{
				i: "123456789012345",
			},
			want: "123456******345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.CreditCard(tt.args.i); got != tt.want {
				t.Errorf("Masker.CreditCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Email(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				i: "dummy@dummy.com",
			},
			want: "dum****@dummy.com",
		},
		{
			name: "Address Less Than 3",
			m:    NewMasker(),
			args: args{
				i: "qq@gmail.com",
			},
			want: "qq****@gmail.com",
		},
		{
			name: "Invalid email format",
			m:    NewMasker(),
			args: args{
				i: "abcd",
			},
			want: "abc****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Email(tt.args.i); got != tt.want {
				t.Errorf("Masker.Email() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Mobile(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				i: "0978978978",
			},
			want: "0978***978",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				i: "0912345678",
			},
			want: "0912***678",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Mobile(tt.args.i); got != tt.want {
				t.Errorf("Masker.Mobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Telephone(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "With Special Chart",
			m:    NewMasker(),
			args: args{
				i: "(02-)27   99-3--078",
			},
			want: "(02)2799-****",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				i: "0227993078",
			},
			want: "(02)2799-****",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				i: "0788079966",
			},
			want: "(07)8807-****",
		},
		{
			name: "Length Not Eq 8 or 10",
			m:    NewMasker(),
			args: args{
				i: "2349966",
			},
			want: "2349966",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Telephone(tt.args.i); got != tt.want {
				t.Errorf("Masker.Telephone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMasker_Password(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		m    *Masker
		args args
		want string
	}{
		{
			name: "Empty Input",
			m:    NewMasker(),
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				i: "1234567",
			},
			want: "************",
		},
		{
			name: "Happy Pass",
			m:    NewMasker(),
			args: args{
				i: "abcd!@#$%321",
			},
			want: "************",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Password(tt.args.i); got != tt.want {
				t.Errorf("Masker.Password() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Masker
	}{
		{
			name: "New Instance",
			want: &Masker{mask: "*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMasker(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		t Mtype
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Error Mask Type",
			args: args{
				t: Mtype(""),
				i: "abcdefg",
			},
			want: "abcdefg",
		},
		{
			name: "Password",
			args: args{
				t: MPassword,
				i: "ggwhite",
			},
			want: "************",
		},
		{
			name: "Name",
			args: args{
				t: MName,
				i: "ggwhite",
			},
			want: "g**hite",
		},
		{
			name: "Address",
			args: args{
				t: MAddress,
				i: "1 AB Paradise Road",
			},
			want: "1 AB P******",
		},
		{
			name: "Email",
			args: args{
				t: MEmail,
				i: "dummy@dummy.com",
			},
			want: "dum****@dummy.com",
		},
		{
			name: "Mobile",
			args: args{
				t: MMobile,
				i: "0978978978",
			},
			want: "0978***978",
		},
		{
			name: "ID",
			args: args{
				t: MID,
				i: "A123456789",
			},
			want: "A12345****",
		},
		{
			name: "Telephone",
			args: args{
				t: MTelephone,
				i: "0288889999",
			},
			want: "(02)8888-****",
		},
		{
			name: "CreditCard",
			args: args{
				t: MCreditCard,
				i: "1234567890123456",
			},
			want: "123456******3456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.t, tt.args.i); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestName(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Length 1",
			args: args{
				i: "A",
			},
			want: "**",
		},
		{
			name: "Length 2",
			args: args{
				i: "AB",
			},
			want: "A**",
		},
		{
			name: "Length 3",
			args: args{
				i: "ABC",
			},
			want: "A**C",
		},
		{
			name: "Length 4",
			args: args{
				i: "ABCD",
			},
			want: "A**D",
		},
		{
			name: "Length 5",
			args: args{
				i: "ABCDE",
			},
			want: "A**DE",
		},
		{
			name: "Length 6",
			args: args{
				i: "ABCDEF",
			},
			want: "A**DEF",
		},
		{
			name: "English Full Name",
			args: args{
				i: "Jorge Marry",
			},
			want: "J**ge M**ry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Name(tt.args.i); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Happy Pass",
			args: args{
				i: "A123456789",
			},
			want: "A12345****",
		},
		{
			name: "Length Less Than 6",
			args: args{
				i: "A12",
			},
			want: "A12****",
		},
		{
			name: "Length Less Than 6",
			args: args{
				i: "A",
			},
			want: "A****",
		},
		{
			name: "Length Between 6 and 10",
			args: args{
				i: "A123456",
			},
			want: "A12345****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ID(tt.args.i); got != tt.want {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddress(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Long Address",
			args: args{
				i: "1 AB Paradise Road",
			},
			want: "1 AB P******",
		},
		{
			name: "Length Less Than 6",
			args: args{
				i: "ABC",
			},
			want: "******",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Address(tt.args.i); got != tt.want {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreditCard(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "VISA JCB MasterCard",
			args: args{
				i: "1234567890123456",
			},
			want: "123456******3456",
		},
		{
			name: "American Express",
			args: args{
				i: "123456789012345",
			},
			want: "123456******345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreditCard(tt.args.i); got != tt.want {
				t.Errorf("CreditCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Happy Pass",
			args: args{
				i: "ggw.chang@gmail.com",
			},
			want: "ggw****ng@gmail.com",
		},
		{
			name: "Address Less Than 3",
			args: args{
				i: "qq@gmail.com",
			},
			want: "qq****@gmail.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Email(tt.args.i); got != tt.want {
				t.Errorf("Email() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMobile(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Happy Pass",
			args: args{
				i: "0978978978",
			},
			want: "0978***978",
		},
		{
			name: "Happy Pass",
			args: args{
				i: "0912345678",
			},
			want: "0912***678",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mobile(tt.args.i); got != tt.want {
				t.Errorf("Mobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTelephone(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "With Special Chart",
			args: args{
				i: "(02-)27   99-3--078",
			},
			want: "(02)2799-****",
		},
		{
			name: "Happy Pass",
			args: args{
				i: "0227993078",
			},
			want: "(02)2799-****",
		},
		{
			name: "Happy Pass",
			args: args{
				i: "0788079966",
			},
			want: "(07)8807-****",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Telephone(tt.args.i); got != tt.want {
				t.Errorf("Telephone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPassword(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty Input",
			args: args{
				i: "",
			},
			want: "",
		},
		{
			name: "Happy Pass",
			args: args{
				i: "1234567",
			},
			want: "************",
		},
		{
			name: "Happy Pass",
			args: args{
				i: "abcd!@#$%321",
			},
			want: "************",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Password(tt.args.i); got != tt.want {
				t.Errorf("Password() = %v, want %v", got, tt.want)
			}
		})
	}
}
