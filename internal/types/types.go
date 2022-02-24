package types

import (
	"io"
	"reflect"
)

// MessagePack types
const (
	FixMap   = 0x80
	FixArray = 0x90
	FixStr   = 0xA0

	NilCode = 0xC0
	False   = 0xC2
	True    = 0xC3

	Bin8  = 0xC4
	Bin16 = 0xC5
	Bin32 = 0xC6

	Ext8  = 0xC7
	Ext16 = 0xC8
	Ext32 = 0xC9

	Float32 = 0xCA
	Float64 = 0xCB

	Uint8  = 0xCC
	Uint16 = 0xCD
	Uint32 = 0xCE
	Uint64 = 0xCF

	Int8  = 0xD0
	Int16 = 0xD1
	Int32 = 0xD2
	Int64 = 0xD3

	FixExt1  = 0xD4
	FixExt2  = 0xD5
	FixExt4  = 0xD6
	FixExt8  = 0xD7
	FixExt16 = 0xD8

	Str8  = 0xD9
	Str16 = 0xDA
	Str32 = 0xDB

	Array16 = 0xDC
	Array32 = 0xDD

	Map16 = 0xDE
	Map32 = 0xDF

	Timestamp = -1
)

// General constants
const (
	NegativeFixIntMin = -32
	Max5Bit           = 0b00011111
)

// Go types
type (
	Boolean bool
	Nil     struct{}
	Int     int64
	Uint    uint64
	Float   struct {
		F               float64
		SinglePrecision bool
	}
	String string
	Binary []byte
	Array  []MessagePackType
	Map    []MessagePackMap
	Struct reflect.Value
)

type MessagePackType interface {
	// Len() int
	// io.ReaderFrom
	io.WriterTo
}

type MessagePackMap struct {
	key   MessagePackType
	value MessagePackType
}
