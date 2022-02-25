package encode

import (
	"github.com/ErikPelli/sbor/internal/types"
	"io"
	"reflect"
	"strings"
)

type EncodingStruct struct {
	visited **struct{} // Indicate if this struct has been already written.
	types.Struct
}

func (e EncodingStruct) Len() int {
	// TODO Improve Struct Len() implementation
	vStruct := reflect.Value(e.Struct)
	// Remove pointers to struct
	for !vStruct.IsZero() && vStruct.Kind() == reflect.Ptr {
		vStruct = vStruct.Elem()
	}

	result := make(types.Map, vStruct.NumField())

	for i := 0; i < vStruct.NumField(); i++ {
		field := vStruct.Type().Field(i)
		tag := field.Tag.Get("sbor")

		if !field.IsExported() {
			continue
		}

		name := field.Name
		omitempty := false

		if tag != "" {
			options := strings.Split(tag, ",")

			if len(options) == 1 {
				if options[0] == "-" {
					continue
				} else {
					name = options[0]
				}
			}

			if len(options) >= 2 {
				if options[1] == "omitempty" {
					omitempty = true
				} else if options[1] == "" {
					name = options[0]
				}
			}
		}

		fieldValue := vStruct.Field(i)
		if omitempty && fieldValue.IsZero() {
			continue
		}

		result[i].Key = types.String(name)
		result[i].Value = TypeWrapper(fieldValue)
	}

	return result.Len()
}

func (e EncodingStruct) WriteTo(w io.Writer) (int64, error) {
	// Skip an already parsed struct (avoid infinite parse in cyclic graph)
	//
	// visited is a double pointer to be able to modify its value with
	// a passed by value EncodingStruct
	if e.visited != nil && *e.visited != nil {
		return 0, nil
	}

	// Use an empty struct to avoid waste of memory
	if e.visited != nil {
		*e.visited = &struct{}{}
	}

	vStruct := reflect.Value(e.Struct)
	// Remove pointers to struct
	for !vStruct.IsZero() && vStruct.Kind() == reflect.Ptr {
		vStruct = vStruct.Elem()
	}

	result := make(types.Map, vStruct.NumField())
	// TODO: Encode as array

	for i := 0; i < vStruct.NumField(); i++ {
		field := vStruct.Type().Field(i)
		tag := field.Tag.Get("sbor")

		if !field.IsExported() {
			continue
		}

		name := field.Name
		omitempty := false

		if tag != "" {
			options := strings.Split(tag, ",")

			if len(options) == 1 {
				if options[0] == "-" {
					continue
				} else {
					name = options[0]
				}
			}

			if len(options) >= 2 {
				if options[1] == "omitempty" {
					omitempty = true
				} else if options[1] == "" {
					name = options[0]
				} else {
					// TODO: Convert value to another type
				}
			}
		}

		fieldValue := vStruct.Field(i)
		if omitempty && fieldValue.IsZero() {
			continue
		}

		result[i].Key = types.String(name) // TODO: Support more key types
		result[i].Value = TypeWrapper(fieldValue)
	}

	return result.WriteTo(w)
}
