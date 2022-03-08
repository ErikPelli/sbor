package types

import (
	"github.com/ErikPelli/sbor/internal/utils"
	"testing"
)

func TestNil_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Input: Nil{}, Expected: []byte{0xC0}},
	}
	utils.TypeWriteToTest(t, data)
}

func TestBoolean_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Input: Boolean(false), Expected: []byte{0xC2}, Name: "false"},
		{Input: Boolean(true), Expected: []byte{0xC3}, Name: "true"},
	}
	utils.TypeWriteToTest(t, data)
}
