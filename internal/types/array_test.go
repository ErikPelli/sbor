package types

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestArray_WriteTo(t *testing.T) {
	data := []WriteTestData{
		{Array([]MessagePackType{
			String("foo"),
			String("bar")}),
			[]byte{0x92, 0xA3, 0x66, 0x6F, 0x6F, 0xA3, 0x62, 0x61, 0x72},
		},
		{Array([]MessagePackType{
			Uint(123),
			Nil{},
			Float{F: 5.5}}),
			[]byte{0x93, 0x7B, 0xC0, 0xCB, 0x40, 0x16, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{Array([]MessagePackType{
			Array([]MessagePackType{Nil{}}),
			Array([]MessagePackType{Uint(1), Uint(2), Uint(3), Uint(4), Uint(5)}),
			Boolean(false)}),
			[]byte{0x93, 0x91, 0xC0, 0x95, 0x01, 0x02, 0x03, 0x04, 0x05, 0xC2},
		},
	}
	TypeWriteToTest(t, data)
}

func TestArray_WriteTo_Arr16(t *testing.T) {
	expected := make([]byte, 3, 1003)
	expected[0] = Array16
	expected[1] = 0x03 // Length
	expected[2] = 0xE8 // 1000 (Big Endian)
	e := bytes.NewBuffer(expected)

	input := make([]MessagePackType, 1000)
	for i := range input {
		input[i] = Boolean(rand.Uint32()%2 == 0)
		_, _ = input[i].WriteTo(e)
	}

	data := []WriteTestData{
		{Array(input), e.Bytes()},
	}
	TypeWriteToTest(t, data)
}

func TestArray_Len_WriteTo_Arr32(t *testing.T) {
	expected := make([]byte, 5, 80005)
	expected[0] = Array32
	expected[1] = 0x00 // Length
	expected[2] = 0x01 // 80000 (Big Endian)
	expected[3] = 0x38
	expected[4] = 0x80
	e := bytes.NewBuffer(expected)

	input := make([]MessagePackType, 80000)
	for i := range input {
		input[i] = Boolean(rand.Uint32()%2 == 0)
		_, _ = input[i].WriteTo(e)
	}

	data := []WriteTestData{
		{Array(input), e.Bytes()},
	}
	TypeWriteToTest(t, data)
}
