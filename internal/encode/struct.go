package encode

import (
	"github.com/ErikPelli/sbor/internal/types"
	"github.com/ErikPelli/sbor/internal/utils"
	"io"
	"reflect"
)

// EncodingStruct is the internal representation of a Go struct
// which is used to do the encoding to MessagePack.
type EncodingStruct struct {
	visited **struct{} //  Indicate if this struct already been written
	state   *EncoderState
	types.Struct
}

func NewEncodingStruct(s types.Struct, e *EncoderState) EncodingStruct {
	visitedPtr := (*struct{})(nil)
	return EncodingStruct{
		visited: &visitedPtr,
		state:   e,
		Struct:  s,
	}
}

// Len returns the length of the MessagePack encoded struct.
// It returns 0 if the struct has already been written.
func (e EncodingStruct) Len() int {
	if e.visited != nil && *e.visited != nil {
		return 0
	}

	valueStruct := reflect.Value(e.Struct)
	if result, encodeAsArray, err := e.structParse(valueStruct); err == nil {
		if encodeAsArray {
			array := make(types.Array, len(result))
			for i := range result {
				array[i] = result[i].Value
			}
			return array.Len()
		}
		return result.Len()
	}
	return 0
}

// WriteTo writes the encoding of the struct value to io.Writer.
// It implements io.WriterTo interface.
// It returns the number of written bytes and an optional error.
// It doesn't write anything if it has already been written, to
// prevent cyclic reference between multiple structs.
func (e EncodingStruct) WriteTo(w io.Writer) (int64, error) {
	// Skip an already parsed struct (avoid infinite parse in cyclic graph)
	//
	// visited is a double pointer to be able to modify its value with
	// a passed by value EncodingStruct
	if e.visited != nil && *e.visited != nil {
		return 0, nil
	}

	valueStruct := reflect.Value(e.Struct)
	if result, encodeAsArray, err := e.structParse(valueStruct); err == nil {
		// Flag this struct as visited
		// Use an empty struct as flag to avoid waste of memory
		if e.visited != nil {
			*e.visited = &struct{}{}
		}

		if encodeAsArray {
			array := make(types.Array, len(result))
			for i := range result {
				array[i] = result[i].Value
			}
			return array.WriteTo(w)
		}
		return result.WriteTo(w)
	} else {
		return 0, err
	}
}

func (e EncodingStruct) structParse(valueStruct reflect.Value) (result types.Map, encodeAsArray bool, err error) {
	numFields := valueStruct.NumField()
	result = make(types.Map, 0, numFields)
	valueStructType := valueStruct.Type()

	var customKeysMap map[string]interface{}
	usedKeysMap := make(map[string]struct{}, numFields)

	for i := 0; i < numFields; i++ {
		field := valueStructType.Field(i)
		fieldValue := valueStruct.Field(i)

		if !field.IsExported() {
			continue
		}

		// Tag parsing
		tagValue := field.Tag.Get("sbor")
		tagName, tagOptions := utils.ParseTag(tagValue)

		var name utils.MessagePackTypeEncoder
		name = types.String(field.Name)

		if tagName == "-" && len(tagValue) == 1 {
			// Skip "-"
			continue
		}

		if tagName != "" {
			// Set name of field using specified name
			name = types.String(tagName)
		}

		if tagOptions.Contains("omitempty") {
			// Skip zero value with omitempty option
			if fieldValue.IsZero() {
				continue
			}
		}

		if tagOptions.Contains("structarray") {
			encodeAsArray = true
		}

		if tagOptions.Contains("setcustomkeys") {
			var ok bool
			mapInterface := fieldValue.Interface()
			customKeysMap, ok = mapInterface.(map[string]interface{})
			if !ok {
				err = utils.InvalidTypeError{Type: "invalid custom keys type"}
				return
			}
			continue
		}

		if tagOptions.Contains("customkey") {
			// Change MessagePack field name using current name as map key
			oldName := string(name.(types.String))
			newName, ok := customKeysMap[oldName]
			if ok {
				name = e.state.TypeWrapper(reflect.ValueOf(newName))
				delete(customKeysMap, oldName)
			} else {
				err = utils.InvalidTypeError{Type: "invalid key " + oldName + " using customkey option"}
				return
			}
		} else {
			// Check duplicated key in standard tag
			checkName := string(name.(types.String))
			_, already := usedKeysMap[checkName]
			if already {
				err = utils.DuplicatedKeyError{Key: name}
				return
			}
			usedKeysMap[checkName] = struct{}{}
		}

		result = append(result, types.MessagePackMap{
			Key:   name,
			Value: e.state.TypeWrapper(fieldValue),
		})
	}
	return
}
