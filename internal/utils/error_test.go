package utils

import (
	"bytes"
	"testing"
)

func TestErrorMessagePackType(t *testing.T) {
	errT := ErrorMessagePackType("error test")
	if errT.Len() != 0 {
		t.Error("Invalid length.")
	}

	n, err := errT.WriteTo(&bytes.Buffer{})
	if n != 0 || err == nil {
		t.Errorf("Invalid result. Len: %d, Error: %v", n, err)
	}
}

func TestInvalidTypeError(t *testing.T) {
	errT := InvalidTypeError{Type: "string"}
	if errT.Error() == "" {
		t.Errorf("Empty error. Error: %v", errT)
	}
}

func TestExceededLengthError(t *testing.T) {
	errT := ExceededLengthError{Type: "string", ActualLength: 2 << 45}
	if errT.Error() == "" {
		t.Errorf("Empty error. Error: %v", errT)
	}
}

func TestDuplicatedKeyError(t *testing.T) {
	errT := DuplicatedKeyError{Key: 98}
	if errT.Error() == "" {
		t.Errorf("Empty error. Error: %v", errT)
	}
}

func TestOutOfBoundError(t *testing.T) {
	errT := OutOfBoundError{Key: 9}
	if errT.Error() == "" {
		t.Errorf("Empty error. Error: %v", errT)
	}
}
