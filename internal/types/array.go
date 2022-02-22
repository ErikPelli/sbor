package types

import (
	"encoding/binary"
	"io"
	"math"
)

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
	}

	nHeader, err := w.Write(header)
	nTotal := int64(nHeader)

	for i := 0; err == nil && i < len(a); i++ {
		var currentN int64
		currentN, err = a[i].WriteTo(w)
		nTotal += currentN
	}

	return nTotal, err
}
