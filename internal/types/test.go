package types

import (
	"bytes"
	"testing"
)

type writeTestData struct {
	input    MessagePackType
	expected []byte
}

func testTypeWriteTo(t *testing.T, data []writeTestData) {
	for _, test := range data {
		var buffer bytes.Buffer

		_, err := test.input.WriteTo(&buffer)
		if err != nil {
			t.Errorf(err.Error())
		}

		result := buffer.Bytes()

		if !bytes.Equal(result, test.expected) {
			t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.expected)
		}

		inputLen := test.input.Len()
		expectedLen := len(test.expected)
		if inputLen != expectedLen {
			t.Errorf("Invalid result length. Function returned %v. Expected %v.", inputLen, expectedLen)
		}
	}
}
