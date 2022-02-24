package types

import (
	"encoding/binary"
	"io"
	"math"
	"strconv"
)

// Len returns the length of the MessagePack encoded map.
// It is a negative value if the data inside is invalid.
func (m Map) Len() int {
	length := len(m)
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

	for i := 0; total >= 0 && i < len(m); i++ {
		k := m[i].Key.Len()
		v := m[i].Value.Len()
		total += k + v

		if k < 0 || v < 0 {
			total = -1
		}
	}

	return total
}

// WriteTo writes the encoding of the map value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (m Map) WriteTo(w io.Writer) (int64, error) {
	var header []byte
	length := len(m)

	switch {
	case length < 1<<4:
		header = make([]byte, 1)
		header[0] = FixMap | byte(length)
	case length <= math.MaxUint16:
		header = make([]byte, 3)
		header[0] = Map16
		binary.BigEndian.PutUint16(header[1:], uint16(length))
	case length <= math.MaxUint32:
		header = make([]byte, 5)
		header[0] = Map32
		binary.BigEndian.PutUint32(header[1:], uint32(length))
	default:
		return 0, InvalidTypeError{"Map exceeded max length. Len: " + strconv.Itoa(length)}
	}

	nHeader, err := w.Write(header)
	nTotal := int64(nHeader)

	for i := range m {
		var n int64
		if err == nil {
			n, err = m[i].Key.WriteTo(w)
			nTotal += n
			if err == nil {
				n, err = m[i].Value.WriteTo(w)
				nTotal += n
			}
		}
	}

	return nTotal, err
}
