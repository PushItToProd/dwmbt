package utils

import (
	"fmt"
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
