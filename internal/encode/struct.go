package encode

import (
	"github.com/ErikPelli/sbor/internal/types"
	"io"
	"reflect"
	"strings"
)

// EncodingStruct is the internal representation of a Go struct
// which is used to do the encoding to MessagePack.
type EncodingStruct struct {
	visited **struct{} // Indicate if this struct has been already written.
	types.Struct
}

func NewEncodingStruct(s types.Struct) EncodingStruct {
	visitedPtr := (*struct{})(nil)
	return EncodingStruct{
		visited: &visitedPtr,
		Struct:  s,
	}
}

// Len returns the length of the MessagePack encoded struct.
// It returns 0 if the struct has been already written.
func (e EncodingStruct) Len() int {
	// TODO Improve Struct Len() implementation
	if e.visited != nil && *e.visited != nil {
		return 0
	}

	vStruct := reflect.Value(e.Struct)
	result := make(types.Map, 0, vStruct.NumField())

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
					if options[0] != "" {
						name = options[0]
					}
				} else if options[1] == "" {
					name = options[0]
				}
			}
		}

		fieldValue := vStruct.Field(i)
		if omitempty && fieldValue.IsZero() {
			continue
		}

		result = append(result, types.MessagePackMap{
			Key:   types.String(name), // TODO: Support more key types
			Value: TypeWrapper(fieldValue),
		})
	}

	return result.Len()
}

// WriteTo writes the encoding of the struct value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
// It doesn't write anything if it has been already written, to
// prevent cyclic reference between multiple structs.
func (e EncodingStruct) WriteTo(w io.Writer) (int64, error) {
	// Skip an already parsed struct (avoid infinite parse in cyclic graph)
	//
	// visited is a double pointer to be able to modify its value with
	// a passed by value EncodingStruct
	if e.visited != nil && *e.visited != nil {
		return 0, nil
	}

	// Use an empty struct as flag to avoid waste of memory
	if e.visited != nil {
		*e.visited = &struct{}{}
	}

	vStruct := reflect.Value(e.Struct)
	result := make(types.Map, 0, vStruct.NumField())
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
					if options[0] != "" {
						name = options[0]
					}
				} else if options[1] == "" {
					name = options[0]
				}
				// TODO: Convert value to another type
			}
		}

		fieldValue := vStruct.Field(i)
		if omitempty && fieldValue.IsZero() {
			continue
		}

		result = append(result, types.MessagePackMap{
			Key:   types.String(name), // TODO: Support more key types
			Value: TypeWrapper(fieldValue),
		})
	}

	return result.WriteTo(w)
}
