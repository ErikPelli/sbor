package types

import (
	"github.com/ErikPelli/sbor/internal/utils"
	"testing"
)

func TestInt_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Int(-1), []byte{0xFF}},
		{Int(-4), []byte{0xFC}},
		{Int(120), []byte{0x78}},
		{Int(-120), []byte{Int8, 0x88}},
		{Int(-28000), []byte{Int16, 0x92, 0xA0}},
		{Int(1 << 24), []byte{Int32, 0x01, 0x00, 0x00, 0x00}},
		{Int(-1 << 46), []byte{Int64, 0xFF, 0xFF, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	utils.TypeWriteToTest(t, data)
}

func TestUint_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Uint(120), []byte{0x78}},
		{Uint(199), []byte{Uint8, 0xC7}},
		{Uint(29000), []byte{Uint16, 0x71, 0x48}},
		{Uint(1 << 24), []byte{Uint32, 0x01, 0x00, 0x00, 0x00}},
		{Uint(1 << 46), []byte{Uint64, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	utils.TypeWriteToTest(t, data)
}

func TestFloat_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Float(4.839), []byte{Float64, 0x40, 0x13, 0x5B, 0x22, 0xD0, 0xE5, 0x60, 0x42}},
		{Float(9.5), []byte{Float32, 0x41, 0x18, 0x00, 0x00}},
		{Float(1.25), []byte{Float32, 0x3F, 0xA0, 0x00, 0x00}},
	}
	utils.TypeWriteToTest(t, data)
}
