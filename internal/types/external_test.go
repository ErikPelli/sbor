package types

import (
	"github.com/ErikPelli/sbor/internal/utils"
	"math/rand"
	"testing"
)

func TestExternal_WriteTo_FixExt(t *testing.T) {
	data := []utils.WriteTestData{
		{Input: External{0x10, []byte{0x09}}, Expected: []byte{0xD4, 0x10, 0x09}, Name: "FixExt1"},
		{Input: External{0x17, []byte{0x88, 0x92}}, Expected: []byte{0xD5, 0x17, 0x88, 0x92}, Name: "FixExt2"},
		{Input: External{0x99, []byte{0x04, 0x05, 0x06, 0x07}},
			Expected: []byte{0xD6, 0x99, 0x04, 0x05, 0x06, 0x07}, Name: "FixExt4"},
		{Input: External{0x10, []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
			Expected: []byte{0xD7, 0x10, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}, Name: "FixExt8"},
		{Input: External{0xAA, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
			Expected: []byte{0xD8, 0xAA, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, Name: "FixExt16"},
	}

	utils.TypeWriteToTest(t, data)
}

func TestExternal_WriteTo_Ext8(t *testing.T) {
	input := External{
		Type: 0x67,
		Data: make([]byte, 100),
	}
	rand.Read(input.Data)
	expected := make([]byte, 3, 103)
	expected[0] = Ext8
	expected[1] = byte(len(input.Data))
	expected[2] = input.Type
	expected = append(expected, input.Data...)

	data := []utils.WriteTestData{
		{Input: input, Expected: expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestExternal_WriteTo_Ext16(t *testing.T) {
	input := External{
		Type: 0x78,
		Data: make([]byte, 1000),
	}
	rand.Read(input.Data)
	expected := make([]byte, 4, 1004)
	expected[0] = Ext16
	expected[1] = 0x03 // Length
	expected[2] = 0xE8 // 1000 (Big Endian)
	expected[3] = input.Type
	expected = append(expected, input.Data...)

	data := []utils.WriteTestData{
		{Input: input, Expected: expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestExternal_WriteTo_Ext32(t *testing.T) {
	input := External{
		Type: 0x38,
		Data: make([]byte, 80000),
	}
	rand.Read(input.Data)
	expected := make([]byte, 6, 80006)
	expected[0] = Ext32
	expected[1] = 0x00 // Length
	expected[2] = 0x01 // 80000 (Big Endian)
	expected[3] = 0x38
	expected[4] = 0x80
	expected[5] = input.Type
	expected = append(expected, input.Data...)

	data := []utils.WriteTestData{
		{Input: input, Expected: expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestExternal_WriteTo_ExtError(t *testing.T) {
	input := External{
		Type: 0x03,
		Data: make([]byte, 4294967298),
	}
	var expected []byte

	data := []utils.WriteTestData{
		{Input: input, Expected: expected},
	}
	utils.TypeWriteToTest(t, data, true)
}
