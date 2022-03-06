package utils

import (
	"bytes"
	"testing"
)

type WriteTestData struct {
	Input    MessagePackType
	Expected []byte
}

func TypeWriteToTest(t *testing.T, data []WriteTestData, errorExpected ...bool) {
	var isErrorInvalid bool
	if len(errorExpected) == 0 || !errorExpected[0] {
		isErrorInvalid = true
	}

	for _, test := range data {
		var buffer bytes.Buffer

		inputLen := test.Input.Len()
		expectedLen := len(test.Expected)

		_, err := test.Input.WriteTo(&buffer)
		if isErrorInvalid && err != nil {
			t.Error(err.Error())
		} else if !isErrorInvalid && err == nil {
			t.Error("Error was expected.")
		}

		result := buffer.Bytes()

		if !bytes.Equal(result, test.Expected) {
			t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.Expected)
		}

		if inputLen != expectedLen {
			t.Errorf("Invalid result length. Function returned %v. Expected %v.", inputLen, expectedLen)
		}
	}
}
