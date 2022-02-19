package types

import (
	"bytes"
	"testing"
)

func TestNil_WriteTo(t *testing.T) {
	input := Nil{}
	expected := []byte{NilCode}
	var buffer bytes.Buffer

	_, _ = input.WriteTo(&buffer)
	result := buffer.Bytes()

	if !bytes.Equal(result, expected) {
		t.Errorf("Invalid result. Function returned %v. Expected %v.", result, expected)
	}
}

func TestBoolean_WriteTo(t *testing.T) {
	data := []struct {
		input  bool
		output []byte
	}{
		{false, []byte{False}},
		{true, []byte{True}},
	}

	for _, test := range data {
		var buffer bytes.Buffer

		_, _ = Boolean(test.input).WriteTo(&buffer)
		result := buffer.Bytes()

		if !bytes.Equal(result, test.output) {
			t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.output)
		}
	}
}
