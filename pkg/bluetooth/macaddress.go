package bluetooth

import (
	"regexp"
	"strings"
)

var nonHexChars = regexp.MustCompile(`[^0-9A-Fa-f]`)

// chunkString splits a string into chunks of length chunkSize. The last chunk may be shorter than chunkSize.
// XXX: This might not be safe for multi-byte characters.
func chunkString(str string, chunkSize int) []string {
	chunks := make([]string, 0, (len(str)+chunkSize-1)/chunkSize)
	for i := 0; i < len(str); i += chunkSize {
		end := i + chunkSize
		if end > len(str) {
			end = len(str)
		}
		chunks = append(chunks, str[i:end])
	}
	return chunks
}

// NormalizeMac checks if a MAC Address looks valid and, if so, returns a
// normalized string with 12 hex digits in pairs of two digits delimited by
// colons. If the address doesn't look valid, it returns an empty string and
// false.
func NormalizeMac(mac string) (string, bool) {
	if len(mac) > 17 {
		return "", false
	}
	mac = nonHexChars.ReplaceAllString(mac, "")
	if len(mac) != 12 {
		return "", false
	}

	mac = strings.ToLower(mac)
	macParts := chunkString(mac, 2)
	mac = strings.Join(macParts, ":")
	return mac, true
}
