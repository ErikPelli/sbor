package encode

import (
	"github.com/ErikPelli/sbor/internal/types"
	"github.com/ErikPelli/sbor/internal/utils"
	"reflect"
	"time"
)

// ExtUserHandler is a function that handle a custom encode defined by the user,
// and serializes to a MessagePack External.
type ExtUserHandler struct {
	Type    byte
	Encoder func(interface{}) ([]byte, error)
}

// EncoderState contains data to correctly encode the current type.
type EncoderState struct {
	extUserHandlers map[reflect.Type]ExtUserHandler
}

// NewEncoderState returns a new empty EncoderState, to create a new encoding context.
func NewEncoderState() *EncoderState {
	return &EncoderState{
		extUserHandlers: make(map[reflect.Type]ExtUserHandler),
	}
}

// SetExternalTypeHandler associate a specific data type with a custom encoding
// function provided by the user.
// Code is a number between 0 and 127 and indicate the correspondent MessagePack
// External type code.
// The handler function receives the value to encode as an empty interface, so it
// needs to do a type assertion and provide a byte array as result, along with an
// eventual error (error = nil if there were no errors).
func (e *EncoderState) SetExternalTypeHandler(typeInvolved interface{}, handler ExtUserHandler) error {
	// Max value is 127
	if handler.Type > 0x7F {
		return utils.OutOfBoundError{Key: int(handler.Type)}
	}

	if handler.Encoder == nil {
		return utils.InvalidTypeError{Type: "nil as function"}
	}

	e.extUserHandlers[reflect.TypeOf(typeInvolved)] = handler

	return nil
}

// TypeWrapper convert a primitive type into its messagepack
// correspondent type using reflection.
func (e *EncoderState) TypeWrapper(value reflect.Value) utils.MessagePackTypeEncoder {
	if value.IsValid() {
		// Reserved external
		if value.Type() == reflect.TypeOf(time.Time{}) {
			return types.External{
				Type: byte(Timestamp),
				Data: convertTimestampToBytes(value.Interface().(time.Time)),
			}
		}

		// User external
		if len(e.extUserHandlers) > 0 {
			handler, ok := e.extUserHandlers[value.Type()]
			if ok {
				bytes, err := handler.Encoder(value.Interface())
				if err != nil {
					return utils.ErrorMessagePackType(err.Error())
				} else {
					return types.External{Type: handler.Type, Data: bytes}
				}
			}
		}
	}

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
