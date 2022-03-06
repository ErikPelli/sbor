package types

import (
	"bytes"
	"github.com/ErikPelli/sbor/internal/utils"
	"math/rand"
	"testing"
)

func TestArray_WriteTo(t *testing.T) {
	data := []utils.WriteTestData{
		{Array([]utils.MessagePackType{
			String("foo"),
			String("bar")}),
			[]byte{0x92, 0xA3, 0x66, 0x6F, 0x6F, 0xA3, 0x62, 0x61, 0x72},
			"only strings",
		},
		{Array([]utils.MessagePackType{
			Uint(123),
			Nil{},
			Float(5.5)}),
			[]byte{0x93, 0x7B, 0xC0, 0xCA, 0x40, 0xB0, 0x00, 0x00},
			"mixed types",
		},
		{Array([]utils.MessagePackType{
			Array([]utils.MessagePackType{Nil{}}),
			Array([]utils.MessagePackType{Uint(1), Uint(2), Uint(3), Uint(4), Uint(5)}),
			Boolean(false)}),
			[]byte{0x93, 0x91, 0xC0, 0x95, 0x01, 0x02, 0x03, 0x04, 0x05, 0xC2},
			"nested arrays",
		},
	}
	utils.TypeWriteToTest(t, data)
}

func TestArray_WriteTo_Arr16(t *testing.T) {
	expected := make([]byte, 3, 1003)
	expected[0] = Array16
	expected[1] = 0x03 // Length
	expected[2] = 0xE8 // 1000 (Big Endian)
	e := bytes.NewBuffer(expected)

	input := make([]utils.MessagePackType, 1000)
	for i := range input {
		input[i] = Boolean(rand.Uint32()%2 == 0)
		_, _ = input[i].WriteTo(e)
	}

	data := []utils.WriteTestData{
		{Array(input), e.Bytes(), ""},
	}
	utils.TypeWriteToTest(t, data)
}

func TestArray_Len_WriteTo_Arr32(t *testing.T) {
	expected := make([]byte, 5, 80005)
	expected[0] = Array32
	expected[1] = 0x00 // Length
	expected[2] = 0x01 // 80000 (Big Endian)
	expected[3] = 0x38
	expected[4] = 0x80
	e := bytes.NewBuffer(expected)

	input := make([]utils.MessagePackType, 80000)
	for i := range input {
		input[i] = Boolean(rand.Uint32()%2 == 0)
		_, _ = input[i].WriteTo(e)
	}

	data := []utils.WriteTestData{
		{Array(input), e.Bytes(), ""},
	}
	utils.TypeWriteToTest(t, data)
}

func TestArray_Len_ArrError1(t *testing.T) {
	input := Array(make([]utils.MessagePackType, 1))
	input[0] = utils.ErrorMessagePackType("test")

	if input.Len() != 0 {
		t.Error("Error was expected.")
	}
}
