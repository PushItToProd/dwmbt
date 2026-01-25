package macaddr

import "regexp"

var nonHexChars = regexp.MustCompile(`[^0-9A-Fa-f]`)

func removeNonHex(s string) string {
	return nonHexChars.ReplaceAllString(s, "")
}
