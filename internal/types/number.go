package types

import (
	"encoding/binary"
	"io"
	"math"
)

// WriteTo writes an integer value to the Writer.
// It returns the number of the written bytes
// and an optional error.
func (i Int) WriteTo(w io.Writer) (int64, error) {
	var bytes []byte

	switch {
	case i >= NegativeFixIntMin && i <= math.MaxInt8:
		// negative and positive fix int
		bytes = []byte{
			byte(i),
		}
	case i >= math.MinInt8 && i < NegativeFixIntMin:
		// int8
		bytes = []byte{
			Int8,
			byte(i),
		}
	case i >= math.MinInt16 && i <= math.MaxInt16:
		// int16
		bytes = make([]byte, 3)
		bytes[0] = Int16
		binary.BigEndian.PutUint16(bytes[1:], uint16(i))
	case i >= math.MinInt32 && i <= math.MaxInt32:
		// int32
		bytes = make([]byte, 5)
		bytes[0] = Int32
		binary.BigEndian.PutUint32(bytes[1:], uint32(i))
	default:
		// int64
		bytes = make([]byte, 9)
		bytes[0] = Int64
		binary.BigEndian.PutUint64(bytes[1:], uint64(i))
	}

	writtenBytes, err := w.Write(bytes)
	return int64(writtenBytes), err
}

// WriteTo writes an unsigned integer value to the Writer.
// It returns the number of the written bytes
// and an optional error.
func (i Uint) WriteTo(w io.Writer) (int64, error) {
	var bytes []byte

	switch {
	case i <= math.MaxInt8:
		// positive fix int
		bytes = []byte{
			byte(i),
		}
	case i <= math.MaxUint8:
		// uint8
		bytes = []byte{
			Uint8,
			byte(i),
		}
	case i <= math.MaxUint16:
		// uint16
		bytes = make([]byte, 3)
		bytes[0] = Uint16
		binary.BigEndian.PutUint16(bytes[1:], uint16(i))
	case i <= math.MaxUint32:
		// uint32
		bytes = make([]byte, 5)
		bytes[0] = Uint32
		binary.BigEndian.PutUint32(bytes[1:], uint32(i))
	default:
		// uint64
		bytes = make([]byte, 9)
		bytes[0] = Uint64
		binary.BigEndian.PutUint64(bytes[1:], uint64(i))
	}

	writtenBytes, err := w.Write(bytes)
	return int64(writtenBytes), err
}

// WriteTo writes a floating point value to the Writer.
// It returns the number of the written bytes
// and an optional error.
func (f Float) WriteTo(w io.Writer) (int64, error) {
	var bytes []byte

	if f.SinglePrecision {
		bytes = make([]byte, 5)
		bytes[0] = Float32
		binary.BigEndian.PutUint32(bytes[1:], math.Float32bits(float32(f.F)))
	} else {
		bytes = make([]byte, 9)
		bytes[0] = Float64
		binary.BigEndian.PutUint64(bytes[1:], math.Float64bits(f.F))
	}

	writtenBytes, err := w.Write(bytes)
	return int64(writtenBytes), err
}
