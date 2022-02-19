package types

import (
	"encoding/binary"
	"io"
	"math"
)

// WriteTo writes a string value to the Writer.
// It returns the number of the written bytes
// and an optional error.
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

// WriteTo writes a binary slice value to the Writer.
// It returns the number of the written bytes
// and an optional error.
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
