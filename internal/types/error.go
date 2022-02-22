package types

type InvalidTypeError struct {
	Type string
}

func (i InvalidTypeError) Error() string {
	return "Unexpected type: " + i.Type
}
