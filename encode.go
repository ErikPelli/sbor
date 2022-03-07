package sbor

import (
	"bytes"
	"github.com/ErikPelli/sbor/internal/encode"
	"reflect"
)

// Marshal returns the MessagePack encoding of v.
//
// External types are supported only through Encoder, to be able to separate
// different session, each one with its own defined external types.
//
// Marshal uses the following type-dependent default encodings:
//
// Boolean values encode as MessagePack boolean.
//
// Nil values encode as MessagePack nil (for example an empty pointer).
//
// Integer values encode as MessagePack int.
//
// Floating point values encode as MessagePack float.
//
// String values encode as MessagePack string.
//
// Array, channel and slice values encode as MessagePack array, except that []byte
// encodes as MessagePack binary, and a nil slice encodes as an empty MessagePack array (length=0).
//
//
// Struct values encode as MessagePack map. Each exported struct field becomes a member
// of the object, using the field name as the object key, unless the field is omitted for
// one of the reasons given below.
//
// The encoding of each struct field can be customized by the format string stored under the "sbor"
// key in the struct field's tag. The format string gives the name of the field, possibly followed by
// a comma-separated list of options. The name may be empty in order to specify options without overriding
// the default field name.
//
// The "omitempty" option specifies that the field should be omitted from the encoding if the field has an
// empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.
//
// The "structarray" option specifies that the current struct must be encoded as an array instead of a map,
// so all the keys will be discarded.
//
// As a special case, if the field tag is "-", the field is always omitted. Note that a field with name "-"
// can still be generated using the tag "-,".
// Examples of struct field tags and their meanings:
//
//   // Field appears in MessagePack as key "myName".
//   Field int `sbor:"myName"`
//
//   // Field appears in MessagePack as key "myName" and
//   // the field is omitted from the object if its value is empty,
//   // as defined above.
//   Field int `sbor:"myName,omitempty"`
//
//   // Field appears in MessagePack as key "Field" (the default), but
//   // the field is skipped if empty.
//   // Note the leading comma.
//   Field int `sbor:",omitempty"`
//
//   // Field appears in MessagePack as key "Field" (the default).
//   // Whole struct is now an array.
//   Field int `sbor:",structarray"`
//
//   // Field is ignored by this package.
//   Field int `sbor:"-"`
//
//   // Field appears in MessagePack as key "-".
//   Field int `sbor:"-,"`
//
// It's possible to use other types as key than string.
// You can use "setcustomkeys" option on an exported field which is a map
// that have a string as its key.
// In that way, with the customkey option, you can use a string as the key and
// change the MessagePack key with the value set with this option.
// If you use this option, this field won't be in the MessagePack, it's automatically skipped.
// This option must not contain duplicated values for different keys, or the behavior maybe bad.
//
// The "customkey" specifies that the current name is only a key for the precedent
// set custom keys collection, and so the MessagePack key is the value correspondent
// to the field name.
//
//   // Now this map is used to get correspondent key for other fields.
//   Keys map[string]int `sbor:",setcustomkeys"`
//
//   // Field appears in MessagePack as key "myName".
//   Field int `sbor:"myName"`
//
//   // "myName" is used as key in the Keys map.
//   // In this case the key will be an integer.
//   // If "myName" doesn't exists in map, en error is returned.
//   Field int `sbor:"myName,customkey"`
//
// If a struct has multiple field with the same name, an error will be returned.
//
// Anonymous struct fields are usually marshaled as sub-maps, with the field
// name as key in their parent, or a custom name given in its field tag.
// To force ignoring of an anonymous struct field, give the field a tag of "-".
//
// Map values encode as MessagePack maps. The map's key type can be any type
// supported by MessagePack (int, float, string, map, array, ...).
// Keys of any type are used directly.
//
// Pointer values encode as the value pointed to.
// A nil pointer encodes as the nil MessagePack value.
//
// Interface values encode as the value contained in the interface.
// A nil interface value encodes as the nil MessagePack value.
//
// Complex and function values cannot be encoded in MessagePack.
// Attempting to encode such a value causes Marshal to return
// an InvalidTypeError.
//
// Marshal handle cyclic data structures representing them only once,
// to avoid an infinite loop.
//
func Marshal(v interface{}) ([]byte, error) {
	state := encode.NewEncoderState()
	result := state.TypeWrapper(reflect.ValueOf(v))
	bufferResult := bytes.NewBuffer(make([]byte, 0, result.Len()))
	_, err := result.WriteTo(bufferResult)

	return bufferResult.Bytes(), err
}
