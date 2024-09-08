package fixederrors

import "errors"

var (
	ErrAlreadyExists = errors.New("value already exists")
	ErrNoValue       = errors.New("value not found")
	ErrUnkType       = errors.New("unknown data type")
)
