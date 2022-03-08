package types

import (
	"github.com/ErikPelli/sbor/internal/utils"
	"testing"
)

func TestInt_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Input: Int(-1), Expected: []byte{0xFF}, Name: "negative fixint 1"},
		{Input: Int(-4), Expected: []byte{0xFC}, Name: "negative fixint 2"},
		{Input: Int(120), Expected: []byte{0x78}, Name: "positive fixint"},
		{Input: Int(-120), Expected: []byte{Int8, 0x88}, Name: "int8"},
		{Input: Int(-28000), Expected: []byte{Int16, 0x92, 0xA0}, Name: "int16"},
		{Input: Int(1 << 24), Expected: []byte{Int32, 0x01, 0x00, 0x00, 0x00}, Name: "int32"},
		{Input: Int(-1 << 46), Expected: []byte{Int64, 0xFF, 0xFF, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00}, Name: "int64"},
	}
	utils.TypeWriteToTest(t, data)
}

func TestUint_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Input: Uint(120), Expected: []byte{0x78}, Name: "fixed uint"},
		{Input: Uint(199), Expected: []byte{Uint8, 0xC7}, Name: "uint8"},
		{Input: Uint(29000), Expected: []byte{Uint16, 0x71, 0x48}, Name: "uint16"},
		{Input: Uint(1 << 24), Expected: []byte{Uint32, 0x01, 0x00, 0x00, 0x00}, Name: "uint32"},
		{Input: Uint(1 << 46), Expected: []byte{Uint64, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00}, Name: "uint64"},
	}
	utils.TypeWriteToTest(t, data)
}

func TestFloat_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Input: Float(4.839), Expected: []byte{Float64, 0x40, 0x13, 0x5B, 0x22, 0xD0, 0xE5, 0x60, 0x42}, Name: "float64"},
		{Input: Float(9.5), Expected: []byte{Float32, 0x41, 0x18, 0x00, 0x00}, Name: "float32-1"},
		{Input: Float(1.25), Expected: []byte{Float32, 0x3F, 0xA0, 0x00, 0x00}, Name: "float32-2"},
	}
	utils.TypeWriteToTest(t, data)
}
