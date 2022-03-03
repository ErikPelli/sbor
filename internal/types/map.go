package types

import (
	"encoding/binary"
	"io"
	"math"
	"reflect"
)

// Len returns the length of the MessagePack encoded map.
// It is 0 if the data inside is invalid.
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
		total = 0
	}

	for i := 0; total > 0 && i < len(m); i++ {
		k := m[i].Key.Len()
		v := m[i].Value.Len()
		total += k + v

		if k < 0 || v < 0 {
			total = 0
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
		return 0, ExceededLengthError{Type: "Map", ActualLength: length}
	}

	nHeader, err := w.Write(header)
	nTotal := int64(nHeader)

	keys := make([]MessagePackType, 0, length)

	// Write each element to w (key and value)
	for i := 0; err == nil && i < length; i++ {
		currentKey := m[i].Key
		for j := range keys {
			if reflect.DeepEqual(currentKey, keys[j]) {
				return 0, DuplicatedKeyError{Key: currentKey}
			}
		}
		keys = append(keys, currentKey)

		var nKey int64
		var nValue int64

		nKey, err = currentKey.WriteTo(w)
		if err == nil {
			nValue, err = m[i].Value.WriteTo(w)
		}

		nTotal += nKey + nValue
	}

	return nTotal, err
}
