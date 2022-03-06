package utils

import (
	"testing"
)

func TestSkipTag(t *testing.T) {
	name, options := ParseTag("-")
	if name != "-" || string(options) != "" {
		t.Errorf("Invalid value. Name: %s, Options: %s", name, options)
	}

	if options.Contains("") {
		t.Error("Invalid tag contains, Empty string must be false.")
	}
}

func TestHyphenTag(t *testing.T) {
	name, options := ParseTag("-,b,c,d")
	if name != "-" || string(options) != "b,c,d" {
		t.Errorf("Invalid value. Name: %s, Options: %s", name, options)
	}

	if !(options.Contains("b") && options.Contains("c") && options.Contains("d")) {
		t.Error("Invalid tag contains")
	}

}

func TestTrailingComma(t *testing.T) {
	name, options := ParseTag("trailing,")
	if name != "trailing" || string(options) != "" {
		t.Errorf("Invalid value. Name: %s, Options: %s", name, options)
	}

	if options.Contains("") {
		t.Error("Invalid tag contains, Empty string must be false.")
	}
}

func TestOptionNotFound(t *testing.T) {
	_, options := ParseTag("trailing,pizza,mafia,mandolino")
	if options.Contains("coconut") {
		t.Error("Option should not be found.")
	}
}
