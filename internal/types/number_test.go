package types

import (
	"github.com/ErikPelli/sbor/internal/utils"
	"testing"
)

func TestInt_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Int(-1), []byte{0xFF}, "negative fixint 1"},
		{Int(-4), []byte{0xFC}, "negative fixint 2"},
		{Int(120), []byte{0x78}, "positive fixint"},
		{Int(-120), []byte{Int8, 0x88}, "int8"},
		{Int(-28000), []byte{Int16, 0x92, 0xA0}, "int16"},
		{Int(1 << 24), []byte{Int32, 0x01, 0x00, 0x00, 0x00}, "int32"},
		{Int(-1 << 46), []byte{Int64, 0xFF, 0xFF, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00}, "int64"},
	}
	utils.TypeWriteToTest(t, data)
}

func TestUint_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Uint(120), []byte{0x78}, "fixed uint"},
		{Uint(199), []byte{Uint8, 0xC7}, "uint8"},
		{Uint(29000), []byte{Uint16, 0x71, 0x48}, "uint16"},
		{Uint(1 << 24), []byte{Uint32, 0x01, 0x00, 0x00, 0x00}, "uint32"},
		{Uint(1 << 46), []byte{Uint64, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00}, "uint64"},
	}
	utils.TypeWriteToTest(t, data)
}

func TestFloat_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Float(4.839), []byte{Float64, 0x40, 0x13, 0x5B, 0x22, 0xD0, 0xE5, 0x60, 0x42}, "float64"},
		{Float(9.5), []byte{Float32, 0x41, 0x18, 0x00, 0x00}, "float32-1"},
		{Float(1.25), []byte{Float32, 0x3F, 0xA0, 0x00, 0x00}, "float32-2"},
	}
	utils.TypeWriteToTest(t, data)
}
