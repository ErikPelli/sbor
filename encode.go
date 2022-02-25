package sbor

import (
	"bytes"
	"github.com/ErikPelli/sbor/internal/encode"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	value := reflect.ValueOf(v)

	for !value.IsZero() && value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	result := encode.TypeWrapper(value)
	bufferResult := bytes.NewBuffer(make([]byte, result.Len()))
	_, err := result.WriteTo(bufferResult)

	return bufferResult.Bytes(), err
}
