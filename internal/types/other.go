package types

import "io"

// WriteTo writes a nil value to the Writer.
// It returns the number of the written bytes
// and an optional error.
func (n Nil) WriteTo(w io.Writer) (int64, error) {
	writtenBytes, err := w.Write([]byte{NilCode})
	return int64(writtenBytes), err
}

// WriteTo writes a boolean value to the Writer.
// It returns the number of the written bytes
// and an optional error.
func (b Boolean) WriteTo(w io.Writer) (int64, error) {
	var value byte

	if b {
		value = True
	} else {
		value = False
	}

	writtenBytes, err := w.Write([]byte{value})
	return int64(writtenBytes), err
}
