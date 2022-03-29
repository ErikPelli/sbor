package types

import (
	"github.com/ErikPelli/sbor/internal/utils"
	"io"
)

// Len returns the length of the MessagePack encoded null.
func (n Nil) Len() int {
	return 1
}

// WriteTo writes the encoding of the null value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
func (n Nil) WriteTo(w io.Writer) (int64, error) {
	writtenBytes, err := w.Write([]byte{NilCode})
	return int64(writtenBytes), err
}

func (n Nil) ReadFrom(code byte, r io.Reader) (int64, error) {
	// TODO
	return 0, utils.InvalidArgumentError{}
}

// Len returns the length of the MessagePack encoded boolean.
func (b Boolean) Len() int {
	return 1
}

// WriteTo writes the encoding of the boolean value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
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

func (b *Boolean) ReadFrom(code byte, r io.Reader) (int64, error) {
	// TODO
	return 0, utils.InvalidArgumentError{}
}
