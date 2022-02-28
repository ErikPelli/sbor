package sbor

import (
	"bytes"
	"testing"
)

func TestMarshal_And_Encoder(t *testing.T) {
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

	expected := []byte{0x83, 0xA7, 0x66, 0x6C, 0x6F, 0x61, 0x74, 0x36, 0x34, 0xCB, 0x40, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0xA1, 0x2D, 0xA6, 0x68, 0x79, 0x70, 0x68, 0x65, 0x6E, 0xA8, 0x75, 0x6E, 0x73, 0x69, 0x67, 0x6E, 0x65, 0x64, 0x20}

	r, err := Marshal(exampleStruct)
	if err != nil {
		t.Errorf("Marshal Error: %v", err)
	}

	var b bytes.Buffer
	e := NewEncoder(&b)
	err = e.Encode(exampleStruct)
	if err != nil {
		t.Errorf("Encoder Error: %v", err)
	}

	if !bytes.Equal(r, expected) {
		t.Errorf("Marshal output different than expected. Returned %v. Expected %v.", r, expected)
	}

	if !bytes.Equal(b.Bytes(), expected) {
		t.Errorf("Encoder output different than expected. Returned %v. Expected %v.", b.Bytes(), expected)
	}
}
