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
	// TODO
	return 0
}

func (e EncodingStruct) WriteTo(w io.Writer) (int64, error) {
	// Skip an already parsed struct (avoid infinite parse in cyclic graph)
	//
	// Visited is a double pointer to be able to modify its value with
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

		result[i].Key = types.String(name)
		result[i].Value = typeWrapper(fieldValue)
	}

	return result.WriteTo(w)
}

func typeWrapper(value reflect.Value) types.MessagePackTypeEncoder {
	switch value.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return types.Uint(value.Uint())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return types.Int(value.Int())

	case reflect.Float32:
		return types.Float{
			F:               value.Float(),
			SinglePrecision: true,
		}

	case reflect.Float64:
		return types.Float{
			F:               value.Float(),
			SinglePrecision: false,
		}

	case reflect.String:
		return types.String(value.String())

	case reflect.Bool:
		return types.Boolean(value.Bool())

	case reflect.Interface:
		return typeWrapper(value.Elem())

	case reflect.Ptr:
		if value.IsNil() {
			return types.Nil{}
		} else {
			return typeWrapper(value.Elem())
		}

	case reflect.Map:
		if value.IsNil() {
			return types.Nil{}
		}

		mapR := make(types.Map, value.Len())

		iter := value.MapRange()
		i := 0
		for iter.Next() {
			mapR[i].Key = typeWrapper(iter.Key())
			mapR[i].Value = typeWrapper(iter.Value())
			i++
		}

		return mapR

	case reflect.Slice:
		if value.IsNil() {
			return types.Nil{}
		}
		fallthrough // Use reflect.Array code

	case reflect.Array:
		arrayR := make(types.Array, value.Len())
		for i := range arrayR {
			v := value.Index(i)
			arrayR[i] = typeWrapper(v)
		}
		return arrayR

	case reflect.Struct:
		visitedPtr := (*struct{})(nil)
		return EncodingStruct{
			visited: &visitedPtr,
			Struct:  types.Struct(value),
		}

	case reflect.Chan:
		length := value.Len()
		arrayR := make(types.Array, length)

		// Read until channel is closed
		for i := 0; i < value.Len(); i++ {
			r, ok := value.Recv()
			if ok {
				arrayR = append(arrayR, typeWrapper(r))
			}
		}

		return arrayR

	case reflect.Invalid:
		return types.Nil{}
	}

	return types.Nil{}
}
