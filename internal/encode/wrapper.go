package encode

import (
	"github.com/ErikPelli/sbor/internal/types"
	"github.com/ErikPelli/sbor/internal/utils"
	"reflect"
)

// EncoderState contains data to correctly encode the current type.
type EncoderState struct {
}

func NewEncoderState() *EncoderState {
	return &EncoderState{}
}

// TypeWrapper convert a primitive type into its messagepack
// correspondent type using reflection.
func (e *EncoderState) TypeWrapper(value reflect.Value) utils.MessagePackTypeEncoder {
	switch value.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return types.Uint(value.Uint())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return types.Int(value.Int())

	case reflect.Float32, reflect.Float64:
		return types.Float(value.Float())

	case reflect.String:
		return types.String(value.String())

	case reflect.Bool:
		return types.Boolean(value.Bool())

	case reflect.Interface:
		return e.TypeWrapper(value.Elem())

	case reflect.Ptr:
		if value.IsNil() {
			return types.Nil{}
		} else {
			return e.TypeWrapper(value.Elem())
		}

	case reflect.Map:
		mapR := make(types.Map, value.Len())
		iter := value.MapRange()
		for i := 0; iter.Next(); i++ {
			mapR[i].Key = e.TypeWrapper(iter.Key())
			mapR[i].Value = e.TypeWrapper(iter.Value())
		}
		return mapR

	case reflect.Slice:
		if value.Type() == reflect.TypeOf([]byte(nil)) {
			// Binary
			return types.Binary(value.Bytes())
		}
		fallthrough // Use reflect.Array code

	case reflect.Array:
		arrayR := make(types.Array, value.Len())
		for i := range arrayR {
			v := value.Index(i)
			arrayR[i] = e.TypeWrapper(v)
		}
		return arrayR

	case reflect.Struct:
		return NewEncodingStruct(types.Struct(value), e)

	case reflect.Chan:
		length := value.Len()
		arrayR := make(types.Array, length)
		if length > 0 {
			var i int

			// Recover from channel panic
			defer func() {
				if errPanic := recover(); errPanic != nil {
					arrayR = arrayR[:i]
				}
			}()

			// Read until channel is closed
			for i = 0; i < length; i++ {
				r, ok := value.Recv()
				if ok {
					arrayR[i] = e.TypeWrapper(r)
				}
			}
		}
		return arrayR

	case reflect.Invalid:
		return types.Nil{}

	default:
		return utils.ErrorMessagePackType("unknown encoder for this type")
	}
}
