package tag

import (
	"strings"
)

// Options is the string following a comma in a struct field's "sbor"
// tag, or the empty string. It does not include the leading comma.
type Options string

// ParseTag splits a struct field's json tag into its name and
// comma-separated options.
func ParseTag(tag string) (string, Options) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], Options(tag[idx+1:])
	}
	return tag, ""
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o Options) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
