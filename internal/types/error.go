package types

import "strconv"

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
	return e.Type + " exceeded max length. Len: " + strconv.Itoa(e.ActualLength)
}
