package types

import (
	"encoding/binary"
	"io"
	"math"
)

// Len returns the length of the MessagePack encoded integer.
func (i Int) Len() int {
	var length int

	switch {
	case i >= NegativeFixIntMin && i <= math.MaxInt8:
		length = 1
	case i >= math.MinInt8 && i < NegativeFixIntMin:
		length = 2
	case i >= math.MinInt16 && i <= math.MaxInt16:
		length = 3
	case i >= math.MinInt32 && i <= math.MaxInt32:
		length = 5
	default:
		length = 9
	}

	return length
}

// WriteTo writes the encoding of the integer value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (i Int) WriteTo(w io.Writer) (int64, error) {
	bytes := make([]byte, i.Len())

	switch {
	case i >= NegativeFixIntMin && i <= math.MaxInt8:
		// negative and positive fix int
		bytes[0] = byte(i)
	case i >= math.MinInt8 && i < NegativeFixIntMin:
		// int8
		bytes[0] = Int8
		bytes[1] = byte(i)
	case i >= math.MinInt16 && i <= math.MaxInt16:
		// int16
		bytes[0] = Int16
		binary.BigEndian.PutUint16(bytes[1:], uint16(i))
	case i >= math.MinInt32 && i <= math.MaxInt32:
		// int32
		bytes[0] = Int32
		binary.BigEndian.PutUint32(bytes[1:], uint32(i))
	default:
		// int64
		bytes[0] = Int64
		binary.BigEndian.PutUint64(bytes[1:], uint64(i))
	}

	writtenBytes, err := w.Write(bytes)
	return int64(writtenBytes), err
}

// Len returns the length of the MessagePack encoded unsigned integer.
func (u Uint) Len() int {
	var length int

	switch {
	case u <= math.MaxInt8:
		length = 1
	case u <= math.MaxUint8:
		length = 2
	case u <= math.MaxUint16:
		length = 3
	case u <= math.MaxUint32:
		length = 5
	default:
		length = 9
	}

	return length
}

// WriteTo writes the encoding of the unsigned integer value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (u Uint) WriteTo(w io.Writer) (int64, error) {
	bytes := make([]byte, u.Len())

	switch {
	case u <= math.MaxInt8:
		// positive fix int
		bytes[0] = byte(u)
	case u <= math.MaxUint8:
		// uint8
		bytes[0] = Uint8
		bytes[1] = byte(u)
	case u <= math.MaxUint16:
		// uint16
		bytes[0] = Uint16
		binary.BigEndian.PutUint16(bytes[1:], uint16(u))
	case u <= math.MaxUint32:
		// uint32
		bytes[0] = Uint32
		binary.BigEndian.PutUint32(bytes[1:], uint32(u))
	default:
		// uint64
		bytes[0] = Uint64
		binary.BigEndian.PutUint64(bytes[1:], uint64(u))
	}

	writtenBytes, err := w.Write(bytes)
	return int64(writtenBytes), err
}

// Len returns the length of the MessagePack encoded float.
func (f Float) Len() int {
	var length int

	if f.SinglePrecision {
		// Header [1 byte] + 32 bit data [4 byte]
		length = 1 + 4
	} else {
		// Header [1 byte] + 64 bit data [8 byte]
		length = 1 + 8
	}

	return length
}

// WriteTo writes the encoding of the floating point value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (f Float) WriteTo(w io.Writer) (int64, error) {
	bytes := make([]byte, f.Len())

	if f.SinglePrecision {
		bytes[0] = Float32
		binary.BigEndian.PutUint32(bytes[1:], math.Float32bits(float32(f.F)))
	} else {
		bytes[0] = Float64
		binary.BigEndian.PutUint64(bytes[1:], math.Float64bits(f.F))
	}

	writtenBytes, err := w.Write(bytes)
	return int64(writtenBytes), err
}
