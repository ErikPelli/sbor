package types

import (
	"encoding/binary"
	"io"
	"math"
)

// Len returns the length of the MessagePack encoded array.
// It is a negative value if the data inside is invalid.
func (a Array) Len() int {
	length := len(a)
	var total int

	switch {
	case length < 1<<4:
		total = 1
	case length <= math.MaxUint16:
		total = 3
	case length <= math.MaxUint32:
		total = 5
	default:
		total = -1
	}

	for i := 0; total > 0 && i < len(a); i++ {
		currentValue := a[i].Len()
		total += currentValue

		if currentValue < 0 {
			total = -1
		}
	}

	return total
}

// WriteTo writes the encoding of the array value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (a Array) WriteTo(w io.Writer) (int64, error) {
	var header []byte
	length := len(a)

	switch {
	case length < 1<<4:
		header = make([]byte, 1)
		header[0] = FixArray | byte(length)
	case length <= math.MaxUint16:
		header = make([]byte, 3)
		header[0] = Array16
		binary.BigEndian.PutUint16(header[1:], uint16(length))
	case length <= math.MaxUint32:
		header = make([]byte, 5)
		header[0] = Array32
		binary.BigEndian.PutUint32(header[1:], uint32(length))
	default:
		return 0, ExceededLengthError{Type: "Array", ActualLength: length}
	}

	nHeader, err := w.Write(header)
	nTotal := int64(nHeader)

	// Write each element to w
	for i := 0; err == nil && i < length; i++ {
		var currentN int64
		currentN, err = a[i].WriteTo(w)
		nTotal += currentN
	}

	return nTotal, err
}
