package encode

import (
	"github.com/ErikPelli/sbor/internal/types"
	"reflect"
)

func TypeWrapper(value reflect.Value) types.MessagePackTypeEncoder {
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
		return TypeWrapper(value.Elem())

	case reflect.Ptr:
		if value.IsNil() {
			return types.Nil{}
		} else {
			return TypeWrapper(value.Elem())
		}

	case reflect.Map:
		if value.IsNil() {
			return types.Nil{}
		}

		mapR := make(types.Map, value.Len())

		iter := value.MapRange()
		for i := 0; iter.Next(); i++ {
			mapR[i].Key = TypeWrapper(iter.Key())
			mapR[i].Value = TypeWrapper(iter.Value())
		}
		return mapR

	case reflect.Slice:
		if value.IsNil() {
			return types.Nil{}
		}

		if value.Type() == reflect.TypeOf([]byte(nil)) {
			// Binary
			return types.Binary(value.Bytes())
		}

		fallthrough // Use reflect.Array code

	case reflect.Array:
		arrayR := make(types.Array, value.Len())
		for i := range arrayR {
			v := value.Index(i)
			arrayR[i] = TypeWrapper(v)
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
		var i int

		defer func() {
			// Recover from channel panic
			if errPanic := recover(); errPanic != nil {
				arrayR = arrayR[:i]
			}
		}()

		// Read until channel is closed
		for i = 0; i < length; i++ {
			r, ok := value.Recv()
			if ok {
				arrayR[i] = TypeWrapper(r)
			}
		}
		return arrayR

	case reflect.Invalid:
		return types.Nil{}
	}

	return types.Nil{}
}
