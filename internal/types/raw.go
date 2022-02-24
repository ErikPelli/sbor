package types

import (
	"encoding/binary"
	"io"
	"math"
)

// Len returns the length of the MessagePack encoded string.
// It is a negative value if the data inside is invalid.
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
		length = -1
	}
	return length
}

// WriteTo writes the encoding of the string value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (s String) WriteTo(w io.Writer) (int64, error) {
	var bytes []byte
	length := len(s)

	switch {
	case length <= Max5Bit:
		bytes = make([]byte, 1, length+1)
		bytes[0] = FixStr | byte(length)
	case length <= math.MaxUint8:
		bytes = make([]byte, 2, length+2)
		bytes[0] = Str8
		bytes[1] = byte(length)
	case length <= math.MaxUint16:
		bytes = make([]byte, 3, length+3)
		bytes[0] = Str16
		binary.BigEndian.PutUint16(bytes[1:], uint16(length))
	case length <= math.MaxUint32:
		bytes = make([]byte, 5, length+5)
		bytes[0] = Str32
		binary.BigEndian.PutUint32(bytes[1:], uint32(length))
	}

	// Space in byte slice for the string was allocated using make()
	bytes = append(bytes, s...)
	writtenBytes, err := w.Write(bytes)
	return int64(writtenBytes), err
}

// Len returns the length of the MessagePack encoded string.
// It is a negative value if the data inside is invalid.
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
		length = -1
	}
	return length
}

// WriteTo writes the encoding of the binary value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (b Binary) WriteTo(w io.Writer) (int64, error) {
	var bytes []byte
	length := len(b)

	switch {
	case length <= math.MaxUint8:
		bytes = make([]byte, 2, length+2)
		bytes[0] = Bin8
		bytes[1] = byte(length)
	case length <= math.MaxUint16:
		bytes = make([]byte, 3, length+3)
		bytes[0] = Bin16
		binary.BigEndian.PutUint16(bytes[1:], uint16(length))
	case length <= math.MaxUint32:
		bytes = make([]byte, 5, length+5)
		bytes[0] = Bin32
		binary.BigEndian.PutUint32(bytes[1:], uint32(length))
	}

	bytes = append(bytes, b...)
	writtenBytes, err := w.Write(bytes)
	return int64(writtenBytes), err
}
