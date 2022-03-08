# MessagePack serializer in Go

[![Go Report Card](https://goreportcard.com/badge/github.com/ErikPelli/sbor)](https://goreportcard.com/report/github.com/ErikPelli/sbor)
[![CodeQL](https://github.com/ErikPelli/sbor/actions/workflows/codeql.yml/badge.svg)](https://github.com/ErikPelli/sbor/actions/workflows/codeql.yml)
[![Linter](https://github.com/ErikPelli/sbor/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/ErikPelli/sbor/actions/workflows/golangci-lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ErikPelli/sbor.svg)](https://pkg.go.dev/github.com/ErikPelli/sbor)
[![codecov](https://codecov.io/gh/ErikPelli/sbor/branch/master/graph/badge.svg?token=bK8mgSNKwF)](https://codecov.io/gh/ErikPelli/sbor)

[SBOR](https://github.com/ErikPelli/sbor) is a modern and straightforward MessagePack serializer written completely in
Go, without the use of code that use "unsafe" package and thus preserving the cross-compilation characteristics of the
language.

The aim of the project is to make a library that balances the performance with the _ease of use_. Its code must be easy
and understandable, and the tests must be adequate, with unit test code coverage greater than 95%.

## Resources

- [Reference](https://pkg.go.dev/github.com/ErikPelli/sbor)
- [Examples](https://pkg.go.dev/github.com/ErikPelli/sbor#pkg-examples)

## What is MessagePack?

MessagePack is an efficient binary serialization format that lets you exchange data among multiple languages, like JSON,
but it's faster and smaller, and provide support to custom types defined by the user, called extension. In addition, it
supports the transmission of raw bytes, unlike JSON. If you are curious, check
its [specification](https://github.com/msgpack/msgpack/blob/master/spec.md).
MessagePack is not human-readable directly.

## Meaning of the project name

SBOR is an acronym that stands for Serializer to Binary Object Representation, to use a short name for the library and
highlight the fact that the format is binary. Moreover, in the Czech language this word means "choir", and this could
represent the set of different functions that compose the library and that work together to provide the correct output.

Contrary to what the name might say, this project currently has no connection with CBOR format, but potential support in
the future for that too, as some of its features are similar to MessagePack, should not be ruled out.

## Installation

You can install this library using:

```
go get github.com/ErikPelli/sbor
```

At the moment it's still in unstable version.

## Features

- Encoding of primitives, time.Time, arrays, slices, maps, structs and value contained in an interface
- Encoding of struct as an array or as a map (key value)
- Omit only specified fields using sbor:"-"
- Renaming of fields using sbor:"new_field_name"
- Support for every type as the key (it could be an integer, a map, an array, etc.), using custom keys

## TODO

- Cache intermediate results to avoid repetition of certain operations when encoding
- Decode from MessagePack bytes to Go data types

## Quickstart

Please read carefully the documentation contained in reference to learn more about the available functions.

```go
package example

import (
	"fmt"
	"github.com/ErikPelli/sbor"
)

func StructMarshalExample() {
	type Test struct {
		Keys  map[string][]byte `sbor:",setcustomkeys"`
		Hello string            `sbor:",customkey"`
	}
	
	s := Test{
		Keys: map[string][]byte{"Hello" : {0xCC, 0x11, 0xAA, 0x00}},
		Hello: "world",
	}

	result, err := sbor.Marshal(s)
	if err != nil {
		panic(err)
	}
	
	// Byte slice result equivalent to:
	// {
	//      [0xCC, 0x11, 0xAA, 0x00]: "world"
	// }
	fmt.Println(result)
}
```

## License

This project is licensed under the MIT License. See [LICENSE](https://github.com/ErikPelli/sbor/blob/master/LICENSE) for
the full license text.