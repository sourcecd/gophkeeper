// Package fixederrors with custom errors
package fixederrors

import "errors"

// custom service error types
var (
	ErrAlreadyExists            = errors.New("value already exists")
	ErrNoValue                  = errors.New("value not found")
	ErrUnkType                  = errors.New("unknown data type")
	ErrUserAlreadyExists        = errors.New("user already exist")
	ErrUserNotExists            = errors.New("user not exist")
	ErrInvalidToken             = errors.New("invalid token")
	ErrRecordAlreadyExists      = errors.New("item already exist")
	ErrRecordNotFound           = errors.New("record not found")
	ErrInvalidCreditCard        = errors.New("wrong credit card number")
	ErrInvalidTextFormat        = errors.New("invalid text format")
	ErrWrongLoginPasswordFormat = errors.New("wrong login password format")
)
