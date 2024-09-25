// Package dataparser manage data types
package dataparser

import (
	"strings"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
)

// ParsedData data parse interface
type ParsedData interface {
	Parse() ([]byte, error)
}

// manage text data type
type parsedText struct {
	b []byte
}

// manage creditcard data type
type parsedCard struct {
	b []byte
}

// manage binary data type
type parsedBin struct {
	b []byte
}

// manage login/pass data type
type parsedCredentials struct {
	b []byte
}

// NewTextParser create text parser instance
func NewTextParser(b []byte) *parsedText {
	return &parsedText{
		b: b,
	}
}

// NewCardParser create creditcard parser instance
func NewCardParser(b []byte) *parsedCard {
	return &parsedCard{
		b: b,
	}
}

// NewBinParser create binary parser instance
func NewBinParser(b []byte) *parsedBin {
	return &parsedBin{
		b: b,
	}
}

// NewCredentialsParser create login/password parser instance
func NewCredentialsParser(b []byte) *parsedCredentials {
	return &parsedCredentials{
		b: b,
	}
}

// Parse method for text input
func (p *parsedText) Parse() ([]byte, error) {
	if !utf8.Valid(p.b) {
		return nil, fixederrors.ErrInvalidTextFormat
	}
	return p.b, nil
}

// Parse method for creditcard input
func (p *parsedCard) Parse() ([]byte, error) {
	if !govalidator.IsCreditCard(string(p.b)) {
		return nil, fixederrors.ErrInvalidCreditCard
	}
	return p.b, nil
}

// Parse method for binary input
func (p *parsedBin) Parse() ([]byte, error) {
	return p.b, nil
}

// Parse method for login/password input
func (p *parsedCredentials) Parse() ([]byte, error) {
	if !utf8.Valid(p.b) {
		return nil, fixederrors.ErrInvalidTextFormat
	}
	if len(strings.Split(string(p.b), " ")) != 2 {
		return nil, fixederrors.ErrWrongLoginPasswordFormat
	}
	return p.b, nil
}

// Dataparser select parser
func Dataparser(dType string, payload []byte) ParsedData {
	switch dType {
	case "TEXT":
		return NewTextParser(payload)
	case "CREDENTIALS":
		return NewCredentialsParser(payload)
	case "BINARY":
		return NewBinParser(payload)
	case "CARD":
		return NewCardParser(payload)
	}
	return nil
}
