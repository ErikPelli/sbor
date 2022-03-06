package utils

import (
	"bytes"
	"testing"
)

type WriteTestData struct {
	Input    MessagePackType
	Expected []byte
	Name     string
}

func TypeWriteToTest(t *testing.T, data []WriteTestData, errorExpected ...bool) {
	isErrorInvalid := len(errorExpected) == 0 || !errorExpected[0]

	for _, test := range data {
		t.Run(test.Name, func(t *testing.T) {
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
		})
	}
}
