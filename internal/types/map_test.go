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
		{Map([]MessagePackMap{}),
			[]byte{0x80},
		},
		{Map([]MessagePackMap{
			{String("int"), Uint(1)},
		}),
			[]byte{0x81, 0xA3, 0x69, 0x6E, 0x74, 0x01},
		},
		{Map([]MessagePackMap{
			{String("boolean"), Boolean(true)},
			{String("null"), Nil{}},
			{String("string"), String("foo bar")},
			{String("float"), Float(0.5)},
		}),
			[]byte{0x84, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCA, 0x3F, 0x00, 0x00, 0x00},
		},
		{Map([]MessagePackMap{
			{String("boolean"), Boolean(true)},
			{String("null"), Nil{}},
			{String("string"), String("foo bar")},
			{String("int"), Int(-1)},
			{String("float"), Float(0.5)},
		}),
			[]byte{0x85, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA3, 0x69, 0x6E, 0x74, 0xFF, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCA, 0x3F, 0x00, 0x00, 0x00},
		},
	}
	utils.TypeWriteToTest(t, data)
}

func TestMap_WriteTo_Nested(t *testing.T) {
	data := []utils.WriteTestData{
		{Map([]MessagePackMap{
			{String("boolean"), Boolean(true)},
			{String("null"), Nil{}},
			{String("string"), String("foo bar")},
			{String("array"), Array([]utils.MessagePackType{
				String("foo"),
				String("bar"),
			})},
			{String("int"), Int(-2)},
			{String("float"), Float(0.5)},
		}),
			[]byte{0x86, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA5, 0x61, 0x72, 0x72, 0x61, 0x79, 0x92, 0xA3, 0x66, 0x6F, 0x6F, 0xA3, 0x62, 0x61, 0x72, 0xA3, 0x69, 0x6E, 0x74, 0xFE, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCA, 0x3F, 0x00, 0x00, 0x00},
		},
		{Map([]MessagePackMap{
			{String("int"), Int(-2)},
			{String("float"), Float(0.5)},
			{String("boolean"), Boolean(true)},
			{String("null"), Nil{}},
			{String("string"), String("foo bar")},
			{String("array"), Array([]utils.MessagePackType{
				String("foo"),
				String("bar"),
			})},
			{String("object"), Map([]MessagePackMap{
				{String("foo"), Int(-1)},
				{String("bar"), Float(0.5)},
			})},
		}),
			[]byte{0x87, 0xA3, 0x69, 0x6E, 0x74, 0xFE, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCA, 0x3F, 0x00, 0x00, 0x00, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA5, 0x61, 0x72, 0x72, 0x61, 0x79, 0x92, 0xA3, 0x66, 0x6F, 0x6F, 0xA3, 0x62, 0x61, 0x72, 0xA6, 0x6F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x82, 0xA3, 0x66, 0x6F, 0x6F, 0xFF, 0xA3, 0x62, 0x61, 0x72, 0xCA, 0x3F, 0x00, 0x00, 0x00},
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
		elem := MessagePackMap{
			Key:   Int(i),
			Value: Boolean(rand.Uint32()%2 == 0),
		}
		input[i] = elem
		_, _ = elem.Key.WriteTo(e)
		_, _ = elem.Value.WriteTo(e)
	}

	data := []utils.WriteTestData{
		{Map(input), e.Bytes()},
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
		elem := MessagePackMap{
			Key:   Int(i),
			Value: Boolean(rand.Uint32()%2 == 0),
		}
		input[i] = elem
		_, _ = elem.Key.WriteTo(e)
		_, _ = elem.Value.WriteTo(e)
	}

	data := []utils.WriteTestData{
		{Map(input), e.Bytes()},
	}
	utils.TypeWriteToTest(t, data)
}

func BenchmarkMap_Array_Duplication_Check_Small(b *testing.B) {
	mapValues := []MessagePackMap{
		{Key: Int(1), Value: Int(4)},
		{Key: Int(2), Value: Int(3)},
		{Key: Int(3), Value: Int(2)},
		{Key: Int(4), Value: Int(1)},
		{Key: Int(5), Value: Int(9)},
		{Key: Int(6), Value: Int(8)},
		{Key: Int(7), Value: Int(7)},
		{Key: Int(8), Value: Int(6)},
		{Key: Int(9), Value: Int(5)},
		{Key: Int(10), Value: Int(4)},
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
	mapValues := []MessagePackMap{
		{Key: Int(1), Value: Int(4)},
		{Key: Int(2), Value: Int(3)},
		{Key: Int(3), Value: Int(2)},
		{Key: Int(4), Value: Int(1)},
		{Key: Int(5), Value: Int(9)},
		{Key: Int(6), Value: Int(8)},
		{Key: Int(7), Value: Int(7)},
		{Key: Int(8), Value: Int(6)},
		{Key: Int(9), Value: Int(5)},
		{Key: Int(10), Value: Int(4)},
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
		elem := MessagePackMap{
			Key:   Array{Int(i)},
			Value: Boolean(rand.Uint32()%2 == 0),
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
