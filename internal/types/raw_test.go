package types

import (
	"github.com/ErikPelli/sbor/internal/utils"
	"math/rand"
	"strings"
	"testing"
)

func TestString_WriteTo_FixStr(t *testing.T) {
	input := String("hello world")
	expected := []byte{0xAB, 0x68, 0x65, 0x6C, 0x6C, 0x6F,
		0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64}

	data := []utils.WriteTestData{
		{input, expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestString_WriteTo_Str8(t *testing.T) {
	input := String(strings.Repeat("#", 130))
	expected := make([]byte, 132)
	expected[0] = Str8
	expected[1] = byte(len(input))
	for i := 2; i < len(expected); i++ {
		expected[i] = 0x23
	}

	data := []utils.WriteTestData{
		{input, expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestString_WriteTo_Str16(t *testing.T) {
	input := String(strings.Repeat("-", 500))
	expected := make([]byte, 503)
	expected[0] = Str16
	expected[1] = 0x01 // Length
	expected[2] = 0xF4 // 500 (Big Endian)
	for i := 3; i < len(expected); i++ {
		expected[i] = 0x2D
	}

	data := []utils.WriteTestData{
		{input, expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestString_WriteTo_Str32(t *testing.T) {
	input := String(strings.Repeat("9", 70000))
	expected := make([]byte, 70005)
	expected[0] = Str32
	expected[1] = 0x00 // Length
	expected[2] = 0x01 // 17000 (Big Endian)
	expected[3] = 0x11
	expected[4] = 0x70
	for i := 5; i < len(expected); i++ {
		expected[i] = 0x39
	}

	data := []utils.WriteTestData{
		{input, expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestBinary_WriteTo_Bin8(t *testing.T) {
	input := Binary(make([]byte, 100))
	rand.Read(input)
	expected := make([]byte, 2, 102)
	expected[0] = Bin8
	expected[1] = byte(len(input))
	expected = append(expected, input...)

	data := []utils.WriteTestData{
		{input, expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestBinary_WriteTo_Bin16(t *testing.T) {
	input := Binary(make([]byte, 1000))
	rand.Read(input)
	expected := make([]byte, 3, 1003)
	expected[0] = Bin16
	expected[1] = 0x03 // Length
	expected[2] = 0xE8 // 1000 (Big Endian)
	expected = append(expected, input...)

	data := []utils.WriteTestData{
		{input, expected},
	}
	utils.TypeWriteToTest(t, data)
}

func TestBinary_WriteTo_Bin32(t *testing.T) {
	input := Binary(make([]byte, 80000))
	rand.Read(input)
	expected := make([]byte, 5, 80005)
	expected[0] = Bin32
	expected[1] = 0x00 // Length
	expected[2] = 0x01 // 80000 (Big Endian)
	expected[3] = 0x38
	expected[4] = 0x80
	expected = append(expected, input...)

	data := []utils.WriteTestData{
		{input, expected},
	}
	utils.TypeWriteToTest(t, data)
}
