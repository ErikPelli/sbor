package types

import (
	"testing"
)

func TestNil_WriteTo(t *testing.T) {
	data := []writeTestData{
		{Nil{}, []byte{0xC0}},
	}
	testTypeWriteTo(t, data)
}

func TestBoolean_WriteTo(t *testing.T) {
	data := []writeTestData{
		{Boolean(false), []byte{0xC2}},
		{Boolean(true), []byte{0xC3}},
	}
	testTypeWriteTo(t, data)
}
