package encode

import (
	"encoding/binary"
	"time"
)

func convertTimestampToBytes(t time.Time) []byte {
	var result []byte
	seconds := uint64(t.Unix())
	nanoSeconds := t.Nanosecond()

	if (seconds >> 34) == 0 {
		data := (uint64(nanoSeconds) << 34) | seconds

		if (data & 0xffffffff00000000) == 0 {
			// timestamp 32
			result = make([]byte, 4)
			binary.BigEndian.PutUint32(result, uint32(data))
		} else {
			// timestamp 64
			result = make([]byte, 8)
			binary.BigEndian.PutUint64(result, data)
		}

		return result
	}

	// timestamp 96
	result = make([]byte, 12)
	binary.BigEndian.PutUint32(result, uint32(nanoSeconds))
	binary.BigEndian.PutUint64(result[4:], seconds)

	return result
}
