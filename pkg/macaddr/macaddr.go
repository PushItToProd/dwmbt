package macaddr

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

type MacAddr [6]byte

var ErrInvalidMacAddr = errors.New("invalid MAC address")
var ErrWrongDecodeLength = errors.New("wrong decode length (this should never happen!)")

func New(a string) (MacAddr, error) {
	if len(a) < 12 || len(a) > 17 {
		return MacAddr{}, ErrInvalidMacAddr
	}

	a = removeNonHex(a)
	if len(a) != 12 {
		return MacAddr{}, ErrInvalidMacAddr
	}

	mac, err := hex.DecodeString(a)
	if err != nil {
		return MacAddr{}, err
	}
	// check the length of the decoded string - this should always be 6 since we check the string length above, but
	// we check it here to be safe
	if len(mac) != 6 {
		log.Panicf("hex string %q of length %d incorrectly decoded to %d bytes: %v", a, len(a), len(mac), mac)
	}

	var addr MacAddr
	copy(addr[:], mac)
	return addr, nil
}

func (m MacAddr) Format(f MacAddrFormat) string {
	return fmt.Sprintf(f.String(), m[0], m[1], m[2], m[3], m[4], m[5])
}

func (m MacAddr) String() string {
	return m.Format(defaultMacFormat)
}
