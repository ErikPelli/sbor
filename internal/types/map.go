package types

import (
	"encoding/binary"
	"io"
	"math"
)

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
	}

	nHeader, err := w.Write(header)
	nTotal := int64(nHeader)

	for k, v := range m {
		var n int64
		if err == nil {
			n, err = k.WriteTo(w)
			nTotal += n
			if err == nil {
				n, err = v.WriteTo(w)
				nTotal += n
			}
		}
	}

	return nTotal, err
}
