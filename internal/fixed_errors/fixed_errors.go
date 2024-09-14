package fixederrors

import "errors"

var (
	ErrAlreadyExists       = errors.New("value already exists")
	ErrNoValue             = errors.New("value not found")
	ErrUnkType             = errors.New("unknown data type")
	ErrUserAlreadyExists   = errors.New("user already exist")
	ErrUserNotExists       = errors.New("user not exist")
	ErrInvalidToken        = errors.New("invalid token")
	ErrRecordAlreadyExists = errors.New("item already exist")
	ErrRecordNotFound      = errors.New("record not found")
	ErrInvalidCreditCard   = errors.New("wrong credit card number")
	ErrInvalidTextFormat   = errors.New("invalid text format")
)
