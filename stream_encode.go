package sbor

import (
	"bytes"
	"github.com/ErikPelli/sbor/internal/encode"
	"github.com/ErikPelli/sbor/internal/utils"
	"io"
	"reflect"
)

// MessagePackCustom is used to encode an external type
// for which the user have to implement his own code for the encoding.
type MessagePackCustom interface {
	MarshalMsgpack() ([]byte, error)
	UnmarshalMsgpack([]byte) error
}

// CustomEncoder specifies a function to encode the type
// when the MessagePackCustom interface is not implemented by the type.
// This is useful for primitive types unhandled by this library, such
// as complex64 and complex128.
// Encoder receives as input the associated type and must return its byte slice
// representation, or an error if the input is invalid.
// To get the type you should do a type assertion without checking for the success,
// because this function will be used only with the associated type.
// Ex: value := i.(complex64)
type CustomEncoder struct {
	Encoder func(i interface{}) ([]byte, error)
}

// An Encoder writes MessagePack values to an output stream.
type Encoder struct {
	w     io.Writer
	state *encode.EncoderState
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:     w,
		state: encode.NewEncoderState(),
	}
}

// SetExternalType associate a type with an external type code, for this Encoder.
// This link between them is used when encoding a custom type in MessagePack, to give the user the freedom
// to do his own encoding for the current external type.
//
// ID is the correspondent MessagePack External type, and must be a number between 0 and 127.
// Value is an instance of the type, it can be zero value or can contain a value, the important thing is
// that it belongs to the type we have to encode.
// If value implements MessagePackCustom interface, the corresponding method MarshalMsgpack will be used
// for encoding, else you have to provide a CustomEncoder function.
func (e *Encoder) SetExternalType(id int8, value interface{}, c ...CustomEncoder) error {
	_, customInterface := value.(MessagePackCustom)
	if customInterface {
		return e.state.SetExternalTypeHandler(value, encode.ExtUserHandler{
			Type: byte(id),
			Encoder: func(i interface{}) ([]byte, error) {
				encoder := i.(MessagePackCustom)
				return encoder.MarshalMsgpack()
			},
		})
	}

	if len(c) != 1 {
		return utils.InvalidArgumentError{Desc: "CustomEncoder expected"}
	}

	return e.state.SetExternalTypeHandler(value, encode.ExtUserHandler{
		Type:    byte(id),
		Encoder: c[0].Encoder,
	})
}

// Encode writes the MessagePack encoding of v to the stream.
// See the documentation for Marshal for details about the conversion of Go values to MessagePack.
func (e *Encoder) Encode(v interface{}) error {
	result := e.state.TypeWrapper(reflect.ValueOf(v))
	bufferResult := bytes.NewBuffer(make([]byte, 0, result.Len()))
	_, err := result.WriteTo(bufferResult)

	if err == nil {
		_, err = e.w.Write(bufferResult.Bytes())
	}

	return err
}
