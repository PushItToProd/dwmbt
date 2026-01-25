package macaddr

import "strings"

type MacAddrFormat struct {
	Uppercase bool
	Separator string
}

func (f MacAddrFormat) String() string {
	var byteFormat string
	if f.Uppercase {
		byteFormat = "%02X"
	} else {
		byteFormat = "%02x"
	}
	format := strings.Repeat(byteFormat+f.Separator, 5) + byteFormat
	return format
}

var defaultMacFormat = MacAddrFormat{
	Uppercase: false,
	Separator: ":",
}
