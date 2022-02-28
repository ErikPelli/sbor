package sbor

import (
	"bytes"
	"github.com/ErikPelli/sbor/internal/encode"
	"reflect"
)

// Marshal returns the MessagePack encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	value := reflect.ValueOf(v)

	for !value.IsZero() && value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	result := encode.TypeWrapper(value)
	bufferResult := bytes.NewBuffer(make([]byte, 0, result.Len()))
	_, err := result.WriteTo(bufferResult)

	return bufferResult.Bytes(), err
}
