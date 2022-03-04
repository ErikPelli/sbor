package types

import (
	"encoding/binary"
	"io"
	"math"
)

// Len returns the length of the MessagePack encoded external type.
// It is 0 if the data inside is invalid.
func (e External) Len() int {
	length := len(e.Data)

	switch length {
	case 1, 2, 4, 8, 16:
		length += 2
	default:
		if length <= math.MaxUint8 {
			length += 3
		} else if length <= math.MaxUint16 {
			length += 4
		} else if length <= math.MaxUint32 {
			length += 6
		} else {
			length = 0
		}
	}

	return length
}

// WriteTo writes the encoding of the external value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (e External) WriteTo(w io.Writer) (int64, error) {
	var header []byte
	length := len(e.Data)

	var fixExtTypeSet bool
	var fixExtTypeCode byte

	setFixCode := func(code byte) {
		if !fixExtTypeSet {
			fixExtTypeSet = true
			fixExtTypeCode = code
		}
	}

	switch {
	// Fixed length External
	case length == 1:
		setFixCode(FixExt1)
		fallthrough
	case length == 2:
		setFixCode(FixExt2)
		fallthrough
	case length == 4:
		setFixCode(FixExt4)
		fallthrough
	case length == 8:
		setFixCode(FixExt8)
		fallthrough
	case length == 16:
		setFixCode(FixExt16)
		header = make([]byte, 2)
		header[0] = fixExtTypeCode
		header[1] = e.Type

	// Variable length external
	case length <= math.MaxUint8:
		header = make([]byte, 3)
		header[0] = Ext8
		header[1] = byte(length)
		header[2] = e.Type
	case length <= math.MaxUint16:
		header = make([]byte, 4)
		header[0] = Ext16
		binary.BigEndian.PutUint16(header[1:], uint16(length))
		header[3] = e.Type
	case length <= math.MaxUint32:
		header = make([]byte, 6)
		header[0] = Ext32
		binary.BigEndian.PutUint32(header[1:], uint32(length))
		header[5] = e.Type
	default:
		return 0, ExceededLengthError{Type: "External", ActualLength: length}
	}

	headerBytes, err := w.Write(header)
	var dataBytes int
	if err == nil {
		dataBytes, err = w.Write(e.Data)
	}

	return int64(headerBytes + dataBytes), err
}
