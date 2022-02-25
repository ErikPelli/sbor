package sbor

import (
	"io"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(v interface{}) error {
	bytes, err := Marshal(v)
	if err != nil {
		return err
	}

	_, err = e.w.Write(bytes)
	return err
}
