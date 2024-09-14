package dataparser

import (
	"github.com/asaskevich/govalidator"
	fixederrors "github.com/sourcecd/gophkeeper/internal/fixed_errors"
)

type ParsedData interface {
	Parse() ([]byte, error)
}

type parsedText struct {
	b []byte
}
type parsedCard struct {
	b []byte
}
type parsedBin struct {
	b []byte
}
type parsedCredentials struct {
	b []byte
}

func NewTextParser(b []byte) *parsedText {
	return &parsedText{
		b: b,
	}
}
func NewCardParser(b []byte) *parsedCard {
	return &parsedCard{
		b: b,
	}
}
func NewBinParser(b []byte) *parsedBin {
	return &parsedBin{
		b: b,
	}
}
func NewCredentialsParser(b []byte) *parsedCredentials {
	return &parsedCredentials{
		b: b,
	}
}

func (p *parsedText) Parse() ([]byte, error) {
	return p.b, nil
}
func (p *parsedCard) Parse() ([]byte, error) {
	if !govalidator.IsCreditCard(string(p.b)) {
		return nil, fixederrors.ErrInvalidCreditCard
	}
	return p.b, nil
}
func (p *parsedBin) Parse() ([]byte, error) {
	return p.b, nil
}
func (p *parsedCredentials) Parse() ([]byte, error) {
	return p.b, nil
}

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
