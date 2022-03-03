package types

import (
	"encoding/binary"
	"io"
	"math"
)

// Len returns the length of the MessagePack encoded string.
// It is 0 if the data inside is invalid.
func (s String) Len() int {
	length := len(s)

	switch {
	case length <= Max5Bit:
		length += 1
	case length <= math.MaxUint8:
		length += 2
	case length <= math.MaxUint16:
		length += 3
	case length <= math.MaxUint32:
		length += 5
	default:
		length = 0
	}

	return length
}

// WriteTo writes the encoding of the string value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (s String) WriteTo(w io.Writer) (int64, error) {
	var header []byte
	length := len(s)

	switch {
	case length <= Max5Bit:
		header = make([]byte, 1)
		header[0] = FixStr | byte(length)
	case length <= math.MaxUint8:
		header = make([]byte, 2)
		header[0] = Str8
		header[1] = byte(length)
	case length <= math.MaxUint16:
		header = make([]byte, 3)
		header[0] = Str16
		binary.BigEndian.PutUint16(header[1:], uint16(length))
	case length <= math.MaxUint32:
		header = make([]byte, 5)
		header[0] = Str32
		binary.BigEndian.PutUint32(header[1:], uint32(length))
	default:
		return 0, ExceededLengthError{Type: "String", ActualLength: length}
	}

	headerBytes, err := w.Write(header)
	var dataBytes int
	if err == nil {
		dataBytes, err = io.WriteString(w, string(s))
	}

	return int64(headerBytes + dataBytes), err
}

// Len returns the length of the MessagePack encoded string.
// It is 0 if the data inside is invalid.
func (b Binary) Len() int {
	length := len(b)

	switch {
	case length <= math.MaxUint8:
		length += 2
	case length <= math.MaxUint16:
		length += 3
	case length <= math.MaxUint32:
		length += 5
	default:
		length = 0
	}

	return length
}

// WriteTo writes the encoding of the binary value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (b Binary) WriteTo(w io.Writer) (int64, error) {
	var header []byte
	length := len(b)

	switch {
	case length <= math.MaxUint8:
		header = make([]byte, 2)
		header[0] = Bin8
		header[1] = byte(length)
	case length <= math.MaxUint16:
		header = make([]byte, 3)
		header[0] = Bin16
		binary.BigEndian.PutUint16(header[1:], uint16(length))
	case length <= math.MaxUint32:
		header = make([]byte, 5)
		header[0] = Bin32
		binary.BigEndian.PutUint32(header[1:], uint32(length))
	default:
		return 0, ExceededLengthError{Type: "Binary", ActualLength: length}
	}

	headerBytes, err := w.Write(header)
	var dataBytes int
	if err == nil {
		dataBytes, err = w.Write(b)
	}

	return int64(headerBytes + dataBytes), err
}
