package decode

import (
	"github.com/ErikPelli/sbor/internal/types"
	"github.com/ErikPelli/sbor/internal/utils"
	"io"
	"reflect"
	"strconv"
)

// ExtUserHandler is a function that handle a custom decode defined by the user,
// and serializes to a MessagePack External.
type ExtUserHandler struct {
	Type    byte
	Decoder func(b []byte) (interface{}, error)
}

// DecoderState contains data to correctly decode the current type.
type DecoderState struct {
	extUserHandlers map[reflect.Type]ExtUserHandler
}

// NewDecoderState returns a new empty DecoderState, to create a new decoding context.
func NewDecoderState() *DecoderState {
	return &DecoderState{
		extUserHandlers: make(map[reflect.Type]ExtUserHandler),
	}
}

// SetExternalTypeHandler associate a specific data type with a custom decoding
// function provided by the user.
// Code is a number between 0 and 127 and indicate the correspondent MessagePack
// External type code.
// The handler function receives the value to decode as a byte slice, and
// needs to provide a compatible type with the target variable as result, along with an
// eventual error (error = nil if there were no errors).
func (d *DecoderState) SetExternalTypeHandler(typeInvolved interface{}, handler ExtUserHandler) error {
	// Max value is 127
	if handler.Type > 0x7F {
		return utils.OutOfBoundError{Key: int(handler.Type)}
	}

	if handler.Decoder == nil {
		return utils.InvalidTypeError{Type: "nil as function"}
	}

	d.extUserHandlers[reflect.TypeOf(typeInvolved)] = handler

	return nil
}

// NextType returns the next MessagePack type.
// It returns io.EOF if the parsing is finished.
func (d *DecoderState) NextType(r io.Reader) (utils.MessagePackType, error) {
	codeSlice := make([]byte, 1)
	_, err := r.Read(codeSlice)
	if err != nil {
		return nil, err
	}

	typeCode := codeSlice[0]
	var currentType *utils.MessagePackType

	switch typeCode {
	case types.Int8, types.Int16, types.Int32, types.Int64:
		*currentType = new(types.Int)
	case types.Uint8, types.Uint16, types.Uint32, types.Uint64:
		*currentType = new(types.Uint)
	// TODO finish types
	// case ....
	default:
		// Check ranges (fixed values)
		switch {
		case typeCode <= 0x7F || typeCode >= 0xE0:
			// Positive and negative fixed integer
			*currentType = new(types.Int)
		case typeCode <= types.FixMap+15:
			// Fixed Map
			*currentType = new(types.Map)
		case typeCode <= types.FixArray+15:
			// Fixed Array
			*currentType = new(types.Array)
		case typeCode <= types.FixStr+31:
			// Fixed String
			*currentType = new(types.String)
		}
	}

	if currentType == nil {
		return nil, utils.InvalidTypeError{Type: "Unable to decode " + strconv.Itoa(int(typeCode))}
	}

	_, err = (*currentType).ReadFrom(typeCode, r)
	return *currentType, err
}
