package bluetooth

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNormalizeMac(t *testing.T) {
	type testcase struct {
		inputMac  string
		outputMac string
		isValid   bool
	}
	tests := []testcase{
		{"aa:bb:cc:dd:ee:ff", "aa:bb:cc:dd:ee:ff", true},
		{"aa:bb:cc:dd:ee:fz", "", false},
		{"aa:bb:ccdd:ee:ff", "aa:bb:cc:dd:ee:ff", true},
		{"AABBCCDDEEFF", "aa:bb:cc:dd:ee:ff", true},
		{"foobar", "", false},
		{"http://example.xyz", "", false},
	}

	getTestname := func(tc testcase) string {
		if tc.isValid {
			return fmt.Sprintf("%q => %q", tc.inputMac, tc.outputMac)
		} else {
			return fmt.Sprintf("%q => false", tc.inputMac)
		}
	}

	for _, tt := range tests {
		t.Run(getTestname(tt), func(t *testing.T) {
			mac, ok := NormalizeMac(tt.inputMac)
			if ok != tt.isValid {
				t.Errorf("validation failed: got %v, wanted %v", ok, tt.isValid)
			}
			if mac != tt.outputMac {
				t.Errorf("normalization failed: got %q, wanted %q", mac, tt.outputMac)
			}
		})
	}
}

func TestChunkString(t *testing.T) {
	type testcase struct {
		input     string
		chunkSize int
		expected  []string
	}
	tests := []testcase{
		{"abcdef", 2, []string{"ab", "cd", "ef"}},
		{"abcdef", 3, []string{"abc", "def"}},
		{"abcdef", 4, []string{"abcd", "ef"}},
		{"abcdef", 6, []string{"abcdef"}},
		{"abcdef", 1, []string{"a", "b", "c", "d", "e", "f"}},
	}

	getTestname := func(tc testcase) string {
		return fmt.Sprintf("%q, %d => %v", tc.input, tc.chunkSize, tc.expected)
	}

	for _, tt := range tests {
		t.Run(getTestname(tt), func(t *testing.T) {
			chunks := chunkString(tt.input, tt.chunkSize)
			if len(chunks) != len(tt.expected) {
				t.Errorf("chunking failed: got %v (len=%d), wanted %v (len=%d)", chunks, len(chunks), tt.expected, len(tt.expected))
			}
			if !reflect.DeepEqual(chunks, tt.expected) {
				t.Errorf("chunking failed: got %v, wanted %v", chunks, tt.expected)
			}
		})
	}
}
