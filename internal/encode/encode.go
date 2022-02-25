package encode

import (
	"bytes"
	"github.com/ErikPelli/sbor/internal/types"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	value := reflect.ValueOf(v)

	for !value.IsZero() && value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	result := typeWrapper(value)
	bufferResult := bytes.NewBuffer(make([]byte, result.Len()))
	_, err := result.WriteTo(bufferResult)

	return bufferResult.Bytes(), err
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
