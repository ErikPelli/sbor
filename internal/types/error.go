package types

import (
	"fmt"
	"strconv"
)

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
