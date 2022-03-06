package encode

import (
	"bytes"
	"testing"
	"time"
)

func TestExternal_WriteTo_FixExt(t *testing.T) {
	data := []struct {
		t        time.Time
		expected []byte
		name     string
	}{
		{time.Unix(1646580000, 0), []byte{0x62, 0x24, 0xD1, 0x20}, "Timestamp 32"},
		{time.Unix(1646580000, 12345), []byte{0x00, 0x00, 0xC0, 0xE4, 0x62, 0x24, 0xD1, 0x20}, "Timestamp 64"},
		{time.Unix(-100, 98765), []byte{0x00, 0x01, 0x81, 0xCD, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x9C}, "Timestamp 96"},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			result := convertTimestampToBytes(test.t)

			expectedLen := len(test.expected)
			resultLen := len(result)

			if resultLen != expectedLen {
				t.Errorf("Invalid result length. Function returned %v. Expected %v.", resultLen, expectedLen)
			}

			if !bytes.Equal(result, test.expected) {
				t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.expected)
			}
		})
	}
}
