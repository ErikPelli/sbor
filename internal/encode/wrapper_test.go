package encode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/ErikPelli/sbor/internal/types"
	"github.com/ErikPelli/sbor/internal/utils"
	"math"
	"reflect"
	"sort"
	"testing"
	"time"
)

func Test_TypeWrapper_Simple(t *testing.T) {
	a := -4
	ch := make(chan int, 10)
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5

	state := NewEncoderState()

	data := []utils.WriteTestData{
		{state.TypeWrapper(reflect.ValueOf(uint32(90))), []byte{0x5A}, "type conversion"},
		{state.TypeWrapper(reflect.ValueOf(a)), []byte{0xFC}, "integer in var"},
		{state.TypeWrapper(reflect.ValueOf(9.5)), []byte{0xCA, 0x41, 0x18, 0x00, 0x00}, "float32"},
		{state.TypeWrapper(reflect.ValueOf(1.37)), []byte{0xCB, 0x3F, 0xF5, 0xEB, 0x85, 0x1E, 0xB8, 0x51, 0xEC}, "float64"},
		{state.TypeWrapper(reflect.ValueOf(utils.MessagePackType(types.Int(-4)))), []byte{0xFC}, "interface"},

		{state.TypeWrapper(reflect.ValueOf(&a)), []byte{0xFC}, "pointer to int"},
		{state.TypeWrapper(reflect.ValueOf((*int)(nil))), []byte{0xC0}, "empty pointer"},
		{state.TypeWrapper(reflect.ValueOf(nil)), []byte{0xC0}, "nil"},
		{state.TypeWrapper(reflect.ValueOf(true)), []byte{0xC3}, "boolean"},

		{state.TypeWrapper(reflect.ValueOf("hello world")), []byte{0xAB, 0x68, 0x65, 0x6C, 0x6C, 0x6F,
			0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64}, "string"},
		{state.TypeWrapper(reflect.ValueOf([]byte{0x01, 0x02, 0x03, 0x04, 0x05})), []byte{0xC4, 0x05, 0x01, 0x02,
			0x03, 0x04, 0x05}, "byte slice"},
		{state.TypeWrapper(reflect.ValueOf([]string{"foo", "bar"})), []byte{0x92, 0xA3, 0x66, 0x6F, 0x6F,
			0xA3, 0x62, 0x61, 0x72}, "string slice"},
		{state.TypeWrapper(reflect.ValueOf((*[2]string)([]string{"foo", "bar"}))), []byte{0x92, 0xA3, 0x66, 0x6F,
			0x6F, 0xA3, 0x62, 0x61, 0x72}, "pointer to array"},
		{state.TypeWrapper(reflect.ValueOf([]interface{}{123, nil, 5.76})), []byte{0x93, 0x7B, 0xC0, 0xCB, 0x40,
			0x17, 0x0A, 0x3D, 0x70, 0xA3, 0xD7, 0x0A}, "interface slice"},

		{state.TypeWrapper(reflect.ValueOf(map[string]int{"int": 1})), []byte{0x81, 0xA3, 0x69, 0x6E, 0x74, 0x01}, "map"},
		{state.TypeWrapper(reflect.ValueOf(ch)), []byte{0x95, 0x01, 0x02, 0x03, 0x04, 0x05}, "channel"},
		{state.TypeWrapper(reflect.ValueOf(struct {
			Hello int
		}{1})), []byte{0x81, 0xA5, 0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x01}, "struct"},
	}

	utils.TypeWriteToTest(t, data)
}

func Test_TypeWrapper_Nested(t *testing.T) {
	state := NewEncoderState()
	data := []utils.WriteTestData{
		{state.TypeWrapper(reflect.ValueOf(
			map[interface{}]interface{}{
				"int":     -2,
				"float":   1.99,
				"boolean": true,
				"null":    nil,
				"string":  "foo bar",
				"array":   []string{"foo", "bar"},
				"object": map[string]interface{}{
					"foo": -1,
					"bar": 0.5,
				},
			})),
			[]byte{0x87, 0xA3, 0x69, 0x6E, 0x74, 0xFE, 0xA5, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0xCB, 0x3F, 0xFF, 0xD7,
				0x0A, 0x3D, 0x70, 0xA3, 0xD7, 0xA7, 0x62, 0x6F, 0x6F, 0x6C, 0x65, 0x61, 0x6E, 0xC3, 0xA4, 0x6E, 0x75, 0x6C,
				0x6C, 0xC0, 0xA6, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0xA7, 0x66, 0x6F, 0x6F, 0x20, 0x62, 0x61, 0x72, 0xA5,
				0x61, 0x72, 0x72, 0x61, 0x79, 0x92, 0xA3, 0x66, 0x6F, 0x6F, 0xA3, 0x62, 0x61, 0x72, 0xA6, 0x6F, 0x62, 0x6A,
				0x65, 0x63, 0x74, 0x82, 0xA3, 0x66, 0x6F, 0x6F, 0xFF, 0xA3, 0x62, 0x61, 0x72, 0xCA, 0x3F, 0x00, 0x00, 0x00},
			"",
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

func Test_TypeWrapper_Error(t *testing.T) {
	state := NewEncoderState()

	data := []utils.WriteTestData{
		{state.TypeWrapper(reflect.ValueOf(complex64(5.0))), []byte{}, "complex64"},
	}

	utils.TypeWriteToTest(t, data, true)
}

func Test_TypeWrapper_UserHandler(t *testing.T) {
	state := NewEncoderState()
	err := state.SetExternalTypeHandler(complex64(0), ExtUserHandler{Type: 0x10, Encoder: func(i interface{}) ([]byte, error) {
		v, _ := i.(complex64)
		result := make([]byte, 8)

		binary.BigEndian.PutUint32(result, math.Float32bits(real(v)))
		binary.BigEndian.PutUint32(result[4:], math.Float32bits(imag(v)))

		return result, nil
	}})

	if err != nil {
		t.Errorf("Unable to set type handler.")
	}

	data := []utils.WriteTestData{
		{state.TypeWrapper(reflect.ValueOf(complex(float32(9.5), float32(9.5)))),
			[]byte{0xD7, 0x10, 0x41, 0x18, 0x00, 0x00, 0x41, 0x18, 0x00, 0x00},
			"complex64"},
	}

	utils.TypeWriteToTest(t, data)
}

func Test_TypeWrapper_UserHandler_Error1(t *testing.T) {
	state := NewEncoderState()
	err := state.SetExternalTypeHandler(0, ExtUserHandler{Type: 0x9F, Encoder: func(i interface{}) ([]byte, error) {
		return nil, nil
	}})

	if err == nil {
		t.Errorf("Error was expected.")
	}
}

func Test_TypeWrapper_UserHandler_Error2(t *testing.T) {
	state := NewEncoderState()
	err := state.SetExternalTypeHandler(0, ExtUserHandler{Type: 0x10, Encoder: nil})

	if err == nil {
		t.Errorf("Error was expected.")
	}
}

func Test_TypeWrapper_UserHandler_Error3(t *testing.T) {
	state := NewEncoderState()
	err := state.SetExternalTypeHandler(complex64(0), ExtUserHandler{Type: 0x10, Encoder: func(i interface{}) ([]byte, error) {
		return nil, errors.New("test error")
	}})

	if err != nil {
		t.Errorf("Unable to set type handler.")
	}

	data := []utils.WriteTestData{
		{state.TypeWrapper(reflect.ValueOf(complex(float32(9.5), float32(9.5)))),
			[]byte{},
			"complex64"},
	}

	utils.TypeWriteToTest(t, data, true)
}

func Test_TypeWrapper_External(t *testing.T) {
	state := NewEncoderState()

	data := []utils.WriteTestData{
		{state.TypeWrapper(reflect.ValueOf(time.Unix(0, 0))), []byte{0xD6, 0xFF, 0x00, 0x00, 0x00, 0x00}, "time"},
	}

	utils.TypeWriteToTest(t, data)
}
