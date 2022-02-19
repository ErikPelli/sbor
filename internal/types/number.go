package encode

import (
	"encoding/binary"
	"github.com/ErikPelli/sbor/internal/types"
	"io"
	"math"
)

const (
	Max7Bit = 127
	Max5Bit = 31
)

func encodeInteger(i int64, w io.Writer) (int, error) {
	switch {
	case i >= -Max5Bit && i < 0:
		// negative fix int
		writtenBytes, err := w.Write([]byte{types.NegativeFixInt | byte(-i)})
		return writtenBytes, err

	case i >= 0 && i <= Max7Bit:
		// positive fix int
		writtenBytes, err := w.Write([]byte{types.FixInt | byte(i)})
		return writtenBytes, err

	case i >= math.MinInt8 && i <= math.MaxInt8:
		// int 8
		writtenBytes, err := w.Write([]byte{types.Uint8, byte(i)})
		return writtenBytes, err

	case i >= math.MinInt16 && i <= math.MaxInt16:
		// int 16
		bigInt := make([]byte, 3)
		bigInt[0] = types.Int16
		binary.BigEndian.PutUint16(bigInt[1:], uint16(i))
		writtenBytes, err := w.Write(bigInt)
		return writtenBytes, err

	case i >= math.MinInt32 && i <= math.MaxInt32:
		// int 32
		bigInt := make([]byte, 5)
		bigInt[0] = types.Int32
		binary.BigEndian.PutUint32(bigInt[1:], uint32(i))
		writtenBytes, err := w.Write(bigInt)
		return writtenBytes, err

	default:
		// int 64
		bigInt := make([]byte, 9)
		bigInt[0] = types.Int64
		binary.BigEndian.PutUint64(bigInt[1:], uint64(i))
		writtenBytes, err := w.Write(bigInt)
		return writtenBytes, err
	}
}

func encodeUnsignedInteger(i uint64, w io.Writer) int {

}
