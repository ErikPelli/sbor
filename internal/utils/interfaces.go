package utils

import "io"

// MessagePackTypeEncoder contains the methods used to convert the type into bytes
type MessagePackTypeEncoder interface {
	Len() int
	io.WriterTo
}

// MessagePackTypeDecoder contains the methods used to convert the bytes into the type
type MessagePackTypeDecoder interface {
	io.Writer
}

// MessagePackType is a MessagePack-compatible type
type MessagePackType interface {
	MessagePackTypeEncoder
	MessagePackTypeDecoder
}
