package types

import (
	"bytes"
	"testing"
)

func TestFloat_WriteTo(t *testing.T) {
	data := []struct {
		input  Float
		output []byte
	}{
		{Float{F: 4.839}, []byte{Float64, 0x40, 0x13, 0x5B, 0x22, 0xD0, 0xE5, 0x60, 0x42}},
		{Float{F: 9.5, SinglePrecision: true}, []byte{Float32, 0x41, 0x18, 0x00, 0x00}},
		{Float{F: 1.25}, []byte{Float64, 0x3F, 0xF4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, test := range data {
		var buffer bytes.Buffer

		_, _ = test.input.WriteTo(&buffer)
		result := buffer.Bytes()

		if !bytes.Equal(result, test.output) {
			t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.output)
		}
	}
}
