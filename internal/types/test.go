package types

import (
	"bytes"
	"sort"
	"testing"
)

type writeTestData struct {
	input    MessagePackType
	expected []byte
}

func testTypeWriteTo(t *testing.T, data []writeTestData) {
	for _, test := range data {
		var buffer bytes.Buffer

		_, err := test.input.WriteTo(&buffer)
		if err != nil {
			t.Errorf(err.Error())
		}

		result := buffer.Bytes()

		if !bytes.Equal(result, test.expected) {
			t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.expected)
		}
	}
}

func testMapTypeWriteTo(t *testing.T, data []writeTestData) {
	for _, test := range data {
		var buffer bytes.Buffer

		_, err := test.input.WriteTo(&buffer)
		if err != nil {
			t.Errorf(err.Error())
		}

		result := buffer.Bytes()

		// Map encoding is not ordered, so we can't compare data byte per byte
		if len(result) != len(test.expected) || result[0] != test.expected[0] {
			t.Errorf("Length different than expected")
		}

		// Check bytes by removing correct data
		sort.Slice(result, func(i int, j int) bool { return result[i] < result[j] })
		sort.Slice(test.expected, func(i int, j int) bool { return test.expected[i] < test.expected[j] })

		if !bytes.Equal(result, test.expected) {
			t.Errorf("Invalid result. Function returned %v. Expected %v.", result, test.expected)
		}
	}
}
