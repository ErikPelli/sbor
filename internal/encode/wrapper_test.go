package encode

import (
	"bytes"
	"github.com/ErikPelli/sbor/internal/types"
	"reflect"
	"sort"
	"testing"
)

func Test_TypeWrapper_Simple(t *testing.T) {
	a := -4
	ch := make(chan int, 10)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5

	data := []types.WriteTestData{
		{TypeWrapper(reflect.ValueOf(uint32(90))), []byte{0x5A}},
		{TypeWrapper(reflect.ValueOf(a)), []byte{0xFC}},
		{TypeWrapper(reflect.ValueOf(float32(9.5))), []byte{0xCA, 0x41, 0x18, 0x00, 0x00}},
		{TypeWrapper(reflect.ValueOf(1.25)), []byte{0xCB, 0x3F, 0xF4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{TypeWrapper(reflect.ValueOf(types.MessagePackType(types.Int(-4)))), []byte{0xFC}},

		{TypeWrapper(reflect.ValueOf(&a)), []byte{0xFC}},
		{TypeWrapper(reflect.ValueOf((*int)(nil))), []byte{0xC0}},
		{TypeWrapper(reflect.ValueOf(nil)), []byte{0xC0}},
		{TypeWrapper(reflect.ValueOf(true)), []byte{0xC3}},

		{TypeWrapper(reflect.ValueOf("hello world")), []byte{0xAB, 0x68, 0x65, 0x6C, 0x6C, 0x6F,
			0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64}},
		{TypeWrapper(reflect.ValueOf([]byte{0x01, 0x02, 0x03, 0x04, 0x05})), []byte{0xC4, 0x05, 0x01, 0x02,
			0x03, 0x04, 0x05}},
		{TypeWrapper(reflect.ValueOf([]string{"foo", "bar"})), []byte{0x92, 0xA3, 0x66, 0x6F, 0x6F,
			0xA3, 0x62, 0x61, 0x72}},
		{TypeWrapper(reflect.ValueOf((*[2]string)([]string{"foo", "bar"}))), []byte{0x92, 0xA3, 0x66, 0x6F,
			0x6F, 0xA3, 0x62, 0x61, 0x72}},
		{TypeWrapper(reflect.ValueOf([]interface{}{123, nil, 5.5})), []byte{0x93, 0x7B, 0xC0, 0xCB, 0x40,
			0x16, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},

		{TypeWrapper(reflect.ValueOf(map[string]int{"int": 1})), []byte{0x81, 0xA3, 0x69, 0x6E, 0x74, 0x01}},
		{TypeWrapper(reflect.ValueOf(ch)), []byte{0x95, 0x01, 0x02, 0x03, 0x04, 0x05}},
		{TypeWrapper(reflect.ValueOf(struct {
			Hello int
		}{1})), []byte{0x81, 0xA5, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x01}},
	}

	types.TypeWriteToTest(t, data)
}

func Test_TypeWrapper_Nested(t *testing.T) {
	data := []types.WriteTestData{
		{TypeWrapper(reflect.ValueOf(
			map[interface{}]interface{}{
				"int":     -2,
				"float":   0.5,
				"boolean": true,
				"null":    nil,
				"string":  "foo bar",
				"array":   []string{"foo", "bar"},
				"object": map[string]interface{}{
					"foo": -1,
					"bar": 0.5,
				},
			})),
			[]byte{0x87, 0xA3, 0x69, 0x6E, 0x74, 0xFE, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCB, 0x3F, 0xE0,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E,
				0x75, 0x6C, 0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62,
				0x61, 0x72, 0xA5, 0x61, 0x72, 0x72, 0x61, 0x79, 0x92, 0xA3, 0x66, 0x6F, 0x6F, 0xA3, 0x62, 0x61, 0x72,
				0xA6, 0x6F, 0x62, 0x6A, 0x65, 0x63, 0x74, 0x82, 0xA3, 0x66, 0x6F, 0x6F, 0xFF, 0xA3, 0x62, 0x61, 0x72,
				0xCB, 0x3F, 0xE0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}

	for _, test := range data {
		var buffer bytes.Buffer

		_, err := test.Input.WriteTo(&buffer)
		if err != nil {
			t.Errorf(err.Error())
		}

		result := buffer.Bytes()

		// Map encoding is not ordered, so we can't compare data byte per byte
		if len(result) != len(test.Expected) || result[0] != test.Expected[0] {
			t.Errorf("Length different than expected")
		}

		// Check bytes by removing correct data
		sort.Slice(result, func(i int, j int) bool { return result[i] < result[j] })
		sort.Slice(test.Expected, func(i int, j int) bool { return test.Expected[i] < test.Expected[j] })

		if !bytes.Equal(result, test.Expected) {
			t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.Expected)
		}
	}
}
