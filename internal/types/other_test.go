package types

import (
	"github.com/ErikPelli/sbor/internal/utils"
	"testing"
)

func TestNil_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Nil{}, []byte{0xC0}},
	}
	utils.TypeWriteToTest(t, data)
}

func TestBoolean_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Boolean(false), []byte{0xC2}},
		{Boolean(true), []byte{0xC3}},
	}
	utils.TypeWriteToTest(t, data)
}
