package sbor

import (
	"io"
)

// An Encoder writes MessagePack values to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the MessagePack encoding of v to the stream.
// See the documentation for Marshal for details about the conversion of Go values to MessagePack.
func (e *Encoder) Encode(v interface{}) error {
	bytes, err := Marshal(v)
	if err == nil {
		_, err = e.w.Write(bytes)
	}

	return err
}
