package types

import (
	"testing"
)

func TestNil_WriteTo(t *testing.T) {
	data := []WriteTestData{
		{Nil{}, []byte{0xC0}},
	}
	TypeWriteToTest(t, data)
}

func TestBoolean_WriteTo(t *testing.T) {
	data := []WriteTestData{
		{Boolean(false), []byte{0xC2}},
		{Boolean(true), []byte{0xC3}},
	}
	TypeWriteToTest(t, data)
}
