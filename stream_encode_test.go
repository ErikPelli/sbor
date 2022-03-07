package sbor

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
)

type TestExternalCustom struct {
	value string
}

func (t TestExternalCustom) MarshalMsgpack() ([]byte, error) {
	return []byte(t.value), nil
}

func (t *TestExternalCustom) UnmarshalMsgpack(b []byte) error {
	t.value = string(b)
	return nil
}

func TestEncoder_SetExternalType(t *testing.T) {
	exampleStruct := struct {
		Hello *TestExternalCustom `sbor:"hello"`
	}{
		Hello: &TestExternalCustom{"test"},
	}

	expected := []byte{0x81, 0xA5, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0xD6, 0x10, 0x74, 0x65, 0x73, 0x74}

	var b bytes.Buffer
	e := NewEncoder(&b)

	if err := e.SetExternalType(0x10, &TestExternalCustom{}); err != nil {
		t.Errorf("Set External Error: %v", err)
	}

	if err := e.Encode(exampleStruct); err != nil {
		t.Errorf("Encoder Error: %v", err)
	}

	if !bytes.Equal(b.Bytes(), expected) {
		t.Errorf("Encoder output different than expected. Returned %v. Expected %v.", b.Bytes(), expected)
	}
}

func TestEncoder_SetExternalType_CustomEncoder(t *testing.T) {
	exampleStruct := struct {
		Hello complex64 `sbor:"hello"`
	}{
		Hello: complex(float32(9.5), float32(9.5)),
	}

	expected := []byte{0x81, 0xA5, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0xD7, 0x10, 0x41, 0x18, 0x00, 0x00, 0x41, 0x18, 0x00, 0x00}

	var b bytes.Buffer
	e := NewEncoder(&b)

	if err := e.SetExternalType(0x10, complex64(0), CustomEncoder{
		Encoder: func(i interface{}) ([]byte, error) {
			v := i.(complex64)
			result := make([]byte, 8)

			binary.BigEndian.PutUint32(result, math.Float32bits(real(v)))
			binary.BigEndian.PutUint32(result[4:], math.Float32bits(imag(v)))

			return result, nil
		},
	}); err != nil {
		t.Errorf("Set External Error: %v", err)
	}

	if err := e.Encode(exampleStruct); err != nil {
		t.Errorf("Encoder Error: %v", err)
	}

	if !bytes.Equal(b.Bytes(), expected) {
		t.Errorf("Encoder output different than expected. Returned %v. Expected %v.", b.Bytes(), expected)
	}
}

func TestEncoder_SetExternalType_Error(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)

	if err := e.SetExternalType(0x10, int32(0)); err == nil {
		t.Errorf("Set External Error: %v", err)
	}
}
