package types

import (
	"bytes"
	"fmt"
	"github.com/ErikPelli/sbor/internal/utils"
	"math/rand"
	"reflect"
	"testing"
)

func TestMap_WriteTo_Simple(t *testing.T) {
	data := []utils.WriteTestData{
		{
			Input:    Map([]MessagePackMap{}),
			Expected: []byte{0x80},
			Name:     "empty map",
		},
		{
			Input: func() Map {
				s := String("int")
				v := Uint(1)
				return []MessagePackMap{
					{Key: &s, Value: &v},
				}
			}(),
			Expected: []byte{0x81, 0xA3, 0x69, 0x6E, 0x74, 0x01},
			Name:     "one element",
		},
		{
			Input: func() Map {
				_1 := String("boolean")
				_2 := Boolean(true)
				_3 := String("null")
				_5 := String("string")
				_6 := String("foo bar")
				_7 := String("float")
				_8 := Float(0.5)
				return []MessagePackMap{
					{Key: &_1, Value: &_2},
					{Key: &_3, Value: Nil{}},
					{Key: &_5, Value: &_6},
					{Key: &_7, Value: &_8},
				}
			}(),
			Expected: []byte{0x84, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCA, 0x3F, 0x00, 0x00, 0x00},
			Name:     "multiple types 1",
		},
		{
			Input: func() Map {
				_1 := String("boolean")
				_2 := Boolean(true)
				_3 := String("null")
				_5 := String("string")
				_6 := String("foo bar")
				_7 := String("int")
				_8 := Int(-1)
				_9 := String("float")
				_10 := Float(0.5)
				return []MessagePackMap{
					{Key: &_1, Value: &_2},
					{Key: &_3, Value: Nil{}},
					{Key: &_5, Value: &_6},
					{Key: &_7, Value: &_8},
					{Key: &_9, Value: &_10},
				}
			}(),
			Expected: []byte{0x85, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA3, 0x69, 0x6E, 0x74, 0xFF, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCA, 0x3F, 0x00, 0x00, 0x00},
			Name:     "multiple types 2",
		},
	}
	utils.TypeWriteToTest(t, data)
}

func TestMap_WriteTo_Nested(t *testing.T) {
	data := []utils.WriteTestData{
		{
			Input: func() Map {
				_1 := String("boolean")
				_2 := Boolean(true)
				_3 := String("null")
				_5 := String("string")
				_6 := String("foo bar")
				_7 := String("int")
				_8 := Int(-2)
				_9 := String("float")
				_10 := Float(0.5)

				arr := String("array")
				foo := String("foo")
				bar := String("bar")

				return []MessagePackMap{
					{Key: &_1, Value: &_2},
					{Key: &_3, Value: Nil{}},
					{Key: &_5, Value: &_6},
					{Key: &arr, Value: Array{
						&foo,
						&bar,
					}},
					{Key: &_7, Value: &_8},
					{Key: &_9, Value: &_10},
				}
			}(),
			Expected: []byte{0x86, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA5, 0x61, 0x72, 0x72, 0x61, 0x79, 0x92, 0xA3, 0x66, 0x6F, 0x6F, 0xA3, 0x62, 0x61, 0x72, 0xA3, 0x69, 0x6E, 0x74, 0xFE, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCA, 0x3F, 0x00, 0x00, 0x00},
			Name:     "array inside",
		},
		{
			Input: func() Map {
				_1 := String("int")
				_2 := Int(-2)
				_3 := String("float")
				_4 := Float(0.5)
				_5 := String("boolean")
				_6 := Boolean(true)
				_7 := String("null")

				_9 := String("string")
				_10 := String("foo bar")

				arr := String("array")
				foo := String("foo")
				bar := String("bar")

				obj := String("object")
				minusOne := Int(-1)
				floatValue := Float(0.5)

				return []MessagePackMap{
					{Key: &_1, Value: &_2},
					{Key: &_3, Value: &_4},
					{Key: &_5, Value: &_6},
					{Key: &_7, Value: Nil{}},
					{Key: &_9, Value: &_10},
					{Key: &arr, Value: Array{
						&foo,
						&bar,
					}},
					{Key: &obj, Value: Map{
						{Key: &foo, Value: &minusOne},
						{Key: &bar, Value: &floatValue},
					}},
				}
			}(),
			Expected: []byte{0x87, 0xA3, 0x69, 0x6E, 0x74, 0xFE, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCA, 0x3F, 0x00, 0x00, 0x00, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA5, 0x61, 0x72, 0x72, 0x61, 0x79, 0x92, 0xA3, 0x66, 0x6F, 0x6F, 0xA3, 0x62, 0x61, 0x72, 0xA6, 0x6F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x82, 0xA3, 0x66, 0x6F, 0x6F, 0xFF, 0xA3, 0x62, 0x61, 0x72, 0xCA, 0x3F, 0x00, 0x00, 0x00},
			Name:     "nested map",
		},
	}
	utils.TypeWriteToTest(t, data)
}

func TestMap_WriteTo_Map16(t *testing.T) {
	expected := make([]byte, 3, 1003)
	expected[0] = Map16
	expected[1] = 0x03 // Length
	expected[2] = 0xE8 // 1000 (Big Endian)
	e := bytes.NewBuffer(expected)

	input := make([]MessagePackMap, 1000)
	for i := range input {
		iv := Int(i)
		bv := Boolean(rand.Uint32()%2 == 0)
		elem := MessagePackMap{
			Key:   &iv,
			Value: &bv,
		}
		input[i] = elem
		_, _ = elem.Key.WriteTo(e)
		_, _ = elem.Value.WriteTo(e)
	}

	data := []utils.WriteTestData{
		{Input: Map(input), Expected: e.Bytes()},
	}
	utils.TypeWriteToTest(t, data)
}

func TestMap_WriteTo_Map32(t *testing.T) {
	expected := make([]byte, 5, 80005)
	expected[0] = Map32
	expected[1] = 0x00 // Length
	expected[2] = 0x01 // 80000 (Big Endian)
	expected[3] = 0x38
	expected[4] = 0x80
	e := bytes.NewBuffer(expected)

	input := make([]MessagePackMap, 80000)
	for i := range input {
		iv := Int(i)
		bv := Boolean(rand.Uint32()%2 == 0)
		elem := MessagePackMap{
			Key:   &iv,
			Value: &bv,
		}
		input[i] = elem
		_, _ = elem.Key.WriteTo(e)
		_, _ = elem.Value.WriteTo(e)
	}

	data := []utils.WriteTestData{
		{Input: Map(input), Expected: e.Bytes()},
	}
	utils.TypeWriteToTest(t, data)
}

func BenchmarkMap_Array_Duplication_Check_Small(b *testing.B) {
	keys := [10]Int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	values := [10]Int{4, 3, 2, 1, 9, 8, 7, 6, 5, 4}
	mapValues := []MessagePackMap{
		{Key: &keys[0], Value: &values[0]},
		{Key: &keys[1], Value: &values[1]},
		{Key: &keys[2], Value: &values[2]},
		{Key: &keys[3], Value: &values[3]},
		{Key: &keys[4], Value: &values[4]},
		{Key: &keys[5], Value: &values[5]},
		{Key: &keys[6], Value: &values[6]},
		{Key: &keys[7], Value: &values[7]},
		{Key: &keys[8], Value: &values[8]},
		{Key: &keys[9], Value: &values[9]},
	}

	length := len(mapValues)

	for n := 0; n < b.N; n++ {
		keys := make([]utils.MessagePackType, 0, length)
		for i := 0; i < length; i++ {
			currentKey := mapValues[i].Key
			for j := range keys {
				if reflect.DeepEqual(currentKey, keys[j]) {
					b.Errorf("Key collision")
				}
			}
			keys = append(keys, currentKey)
		}
	}
}

func BenchmarkMap_Map_Duplication_Check_Small(b *testing.B) {
	keys := [10]Int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	values := [10]Int{4, 3, 2, 1, 9, 8, 7, 6, 5, 4}
	mapValues := []MessagePackMap{
		{Key: &keys[0], Value: &values[0]},
		{Key: &keys[1], Value: &values[1]},
		{Key: &keys[2], Value: &values[2]},
		{Key: &keys[3], Value: &values[3]},
		{Key: &keys[4], Value: &values[4]},
		{Key: &keys[5], Value: &values[5]},
		{Key: &keys[6], Value: &values[6]},
		{Key: &keys[7], Value: &values[7]},
		{Key: &keys[8], Value: &values[8]},
		{Key: &keys[9], Value: &values[9]},
	}

	length := len(mapValues)

	for n := 0; n < b.N; n++ {
		keys := make(map[string]struct{}, length)
		for i := 0; i < length; i++ {
			currentKey := mapValues[i].Key
			stringKey := fmt.Sprint(currentKey)
			_, ok := keys[stringKey]
			if ok {
				b.Errorf("Key collision")
			}
			keys[stringKey] = struct{}{}
		}
	}
}

func BenchmarkMap_Map_Duplication_Check_Large(b *testing.B) {
	mapValues := make([]MessagePackMap, 80000)
	for i := range mapValues {
		iv := Int(i)
		bv := Boolean(rand.Uint32()%2 == 0)
		elem := MessagePackMap{
			Key:   Array{&iv},
			Value: &bv,
		}
		mapValues[i] = elem
	}

	length := len(mapValues)

	for n := 0; n < b.N; n++ {
		keys := make(map[string]struct{}, length)
		for i := 0; i < length; i++ {
			currentKey := mapValues[i].Key
			stringKey := fmt.Sprint(currentKey)
			_, ok := keys[stringKey]
			if ok {
				b.Errorf("Key collision")
			}
			keys[stringKey] = struct{}{}
		}
	}
}

func TestArray_Len_MapError1(t *testing.T) {
	input := Map(make([]MessagePackMap, 1))
	b := Boolean(false)
	input[0] = MessagePackMap{
		Key:   &b,
		Value: utils.ErrorMessagePackType("test"),
	}

	if input.Len() != 0 {
		t.Error("Error was expected.")
	}
}
