package encode

import (
	"bytes"
	"github.com/ErikPelli/sbor/internal/types"
	"github.com/ErikPelli/sbor/internal/utils"
	"reflect"
	"testing"
)

func TestEncodingStruct_WriteTo(t *testing.T) {
	exampleStruct := struct {
		Hello      int     `sbor:"-"`
		F          float64 `sbor:"float64"`
		Hyphen     string  `sbor:"-,"`
		Bytes      []byte  `sbor:",omitempty"`
		Apple      uint    `sbor:"unsigned,omitempty"`
		unexported bool
	}{
		Hello:  66,
		F:      9.5,
		Hyphen: "hyphen",
		Apple:  32,
	}

	enc := NewEncodingStruct(types.Struct(reflect.ValueOf(exampleStruct)), NewEncoderState())
	data := []utils.WriteTestData{
		{Input: enc, Expected: []byte{0x83, 0xA7, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0x36, 0x34, 0xCA, 0x41, 0x18, 0x00, 0x00, 0xA1,
			0x2D, 0xA6, 0x68, 0x79, 0x70, 0x68, 0x65, 0x6E, 0xA8, 0x75, 0x6E, 0x73, 0x69, 0x67, 0x6E, 0x65, 0x64, 0x20}, Name: "example struct"},

		// Already visited
		{Input: enc, Expected: []byte{}, Name: "already visited"},
	}

	utils.TypeWriteToTest(t, data)
}

func BenchmarkEncodingStruct_WriteTo(b *testing.B) {
	exampleStruct := struct {
		Hello      int     `sbor:"-"`
		F          float64 `sbor:"float64"`
		Hyphen     string  `sbor:"-,"`
		Bytes      []byte  `sbor:",omitempty"`
		Apple      uint    `sbor:"unsigned,omitempty"`
		unexported bool
	}{
		Hello:  66,
		F:      9.5,
		Hyphen: "hyphen",
		Apple:  32,
	}

	for i := 0; i < b.N; i++ {
		state := NewEncoderState()
		enc := state.TypeWrapper(reflect.ValueOf(exampleStruct))
		var buffer bytes.Buffer
		_, _ = enc.WriteTo(&buffer)
	}
}

func TestEncodingStruct_WriteTo_Nested(t *testing.T) {
	type Integers struct {
		A int8   `sbor:"a"`
		B uint16 `sbor:"b"`
		C int32  `sbor:"c"`
	}

	exampleStruct := struct {
		Hyphen string   `sbor:"h"`
		I      Integers `sbor:"i"`
	}{
		Hyphen: "hyphen",
		I: Integers{
			A: -8,
			B: 32000,
			C: -40000,
		},
	}

	enc := NewEncodingStruct(types.Struct(reflect.ValueOf(exampleStruct)), NewEncoderState())
	data := []utils.WriteTestData{
		{Input: enc, Expected: []byte{0x82, 0xA1, 0x68, 0xA6, 0x68, 0x79, 0x70, 0x68, 0x65, 0x6E, 0xA1, 0x69, 0x83, 0xA1,
			0x61, 0xF8, 0xA1, 0x62, 0xCD, 0x7D, 0x00, 0xA1, 0x63, 0xD2, 0xFF, 0xFF, 0x63, 0xC0}, Name: "nested struct"},

		// Already visited
		{Input: enc, Expected: []byte{}, Name: "already visited nested struct"},
	}

	utils.TypeWriteToTest(t, data)
}

func TestEncodingStruct_WriteTo_Array(t *testing.T) {
	type Integers struct {
		A int8   `sbor:"a,structarray"`
		B uint16 `sbor:"b"`
		C int32  `sbor:"c"`
	}

	exampleStruct := struct {
		Hyphen string   `sbor:"h"`
		I      Integers `sbor:"i"`
	}{
		Hyphen: "hyphen",
		I: Integers{
			A: -8,
			B: 32000,
			C: -40000,
		},
	}

	enc := NewEncodingStruct(types.Struct(reflect.ValueOf(exampleStruct)), NewEncoderState())
	data := []utils.WriteTestData{
		{Input: enc, Expected: []byte{0x82, 0xA1, 0x68, 0xA6, 0x68, 0x79, 0x70, 0x68, 0x65, 0x6E, 0xA1, 0x69, 0x93, 0xF8, 0xCD, 0x7D, 0x00, 0xD2, 0xFF, 0xFF, 0x63, 0xC0}, Name: "nested array struct"},
	}

	utils.TypeWriteToTest(t, data)
}

func TestEncodingStruct_WriteTo_DuplicatedKey(t *testing.T) {
	type Integers struct {
		A int8   `sbor:"a"`
		B uint16 `sbor:"a"`
		C int32  `sbor:"c"`
	}

	exampleStruct := Integers{
		A: -8,
		B: 32000,
		C: -40000,
	}

	enc := NewEncodingStruct(types.Struct(reflect.ValueOf(exampleStruct)), NewEncoderState())
	data := []utils.WriteTestData{
		{Input: enc, Expected: []byte{}, Name: "duplicated key"},
	}

	utils.TypeWriteToTest(t, data, true)
}

func TestEncodingStruct_WriteTo_SetCustomKeys(t *testing.T) {
	type Integers struct {
		A map[string]int8 `sbor:",setcustomkeys"`
		B uint16          `sbor:"b,customkey"`
		C int32           `sbor:"c,customkey"`
	}

	exampleStruct := Integers{
		A: map[string]int8{"b": -8, "c": 38},
		B: 32000,
		C: -40000,
	}

	enc := NewEncodingStruct(types.Struct(reflect.ValueOf(exampleStruct)), NewEncoderState())
	data := []utils.WriteTestData{
		{Input: enc, Expected: []byte{0x82, 0xF8, 0xCD, 0x7D, 0x00, 0x26, 0xD2, 0xFF, 0xFF, 0x63, 0xC0}, Name: "custom key"},
	}

	utils.TypeWriteToTest(t, data)
}

func TestEncodingStruct_WriteTo_SetCustomKeys_Duplicated(t *testing.T) {
	type Integers struct {
		A map[string]interface{} `sbor:",setcustomkeys"`
		B uint16                 `sbor:"b,customkey"`
		C int32                  `sbor:"b,customkey"`
	}

	exampleStruct := Integers{
		A: map[string]interface{}{"b": -8, "c": 38},
		B: 32000,
		C: -40000,
	}

	enc := NewEncodingStruct(types.Struct(reflect.ValueOf(exampleStruct)), NewEncoderState())
	data := []utils.WriteTestData{
		{Input: enc, Expected: []byte{}, Name: "duplicated custom key"},
	}

	utils.TypeWriteToTest(t, data, true)
}

func TestEncodingStruct_WriteTo_SetCustomKeys_Invalid(t *testing.T) {
	type Integers struct {
		A int8   `sbor:",setcustomkeys"`
		B uint16 `sbor:"b"`
		C int32  `sbor:"c"`
	}

	exampleStruct := Integers{
		A: -8,
		B: 32000,
		C: -40000,
	}

	enc := NewEncodingStruct(types.Struct(reflect.ValueOf(exampleStruct)), NewEncoderState())
	data := []utils.WriteTestData{
		{Input: enc, Expected: []byte{}, Name: "invalid setcustomkeys"},
	}

	utils.TypeWriteToTest(t, data, true)
}
