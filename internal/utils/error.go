package utils

import (
	"fmt"
	"io"
	"strconv"
)

// ErrorMessagePackType returns always an error if you try to write it.
type ErrorMessagePackType string

func (e ErrorMessagePackType) Len() int {
	return 0
}

func (e ErrorMessagePackType) WriteTo(w io.Writer) (int64, error) {
	return 0, InvalidTypeError{string(e)}
}

type InvalidTypeError struct {
	Type string
}

func (i InvalidTypeError) Error() string {
	return "Unexpected type: " + i.Type
}

type ExceededLengthError struct {
	Type         string
	ActualLength int
}

func (e ExceededLengthError) Error() string {
	return e.Type + " exceeded max length (len: " + strconv.Itoa(e.ActualLength) + ")"
}

type DuplicatedKeyError struct {
	Key interface{}
}

func (d DuplicatedKeyError) Error() string {
	return fmt.Sprintf("Duplicated key %v", d.Key)
}

type OutOfBoundError struct {
	Key int
}

func (o OutOfBoundError) Error() string {
	return fmt.Sprintf("Index %d out of bound", o.Key)
}
