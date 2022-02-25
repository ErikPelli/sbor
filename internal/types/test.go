package types

import (
	"bytes"
	"testing"
)

type WriteTestData struct {
	Input    MessagePackType
	Expected []byte
}

func TypeWriteToTest(t *testing.T, data []WriteTestData) {
	for _, test := range data {
		var buffer bytes.Buffer

		_, err := test.Input.WriteTo(&buffer)
		if err != nil {
			t.Errorf(err.Error())
		}

		result := buffer.Bytes()

		if !bytes.Equal(result, test.Expected) {
			t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.Expected)
		}

		inputLen := test.Input.Len()
		expectedLen := len(test.Expected)
		if inputLen != expectedLen {
			t.Errorf("Invalid result length. Function returned %v. Expected %v.", inputLen, expectedLen)
		}
	}
}
