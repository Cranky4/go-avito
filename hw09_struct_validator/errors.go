package hw09structvalidator

import (
	"errors"
	"fmt"
)

var ErrInvalidInputArgument = errors.New("input is not struct")

// ErrInvalidValidatorTag.
type ErrInvalidValidatorTag struct {
	Field, Tag string
}

func (e ErrInvalidValidatorTag) Error() string {
	return fmt.Sprintf("invalid validatior tag %s for field %s", e.Tag, e.Field)
}

// ErrInvalidValidatorTagValue.
type ErrInvalidValidatorTagValue struct {
	ExpectedType string
	CurrentValue interface{}
}

func (e ErrInvalidValidatorTagValue) Error() string {
	return fmt.Sprintf("invalid validatior tag value %v, expected type is %s", e.CurrentValue, e.ExpectedType)
}

// ErrValidatorNotExists.
type ErrValidatorNotExists struct {
	ValidatorName string
}

func (e ErrValidatorNotExists) Error() string {
	return fmt.Sprintf("validator with name '%s' not exisis", e.ValidatorName)
}

// ErrValidaton.
type ErrValidaton struct {
	Message string
}

func (e ErrValidaton) Error() string {
	return e.Message
}
