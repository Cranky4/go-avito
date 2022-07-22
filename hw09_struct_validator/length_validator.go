package hw09structvalidator

import (
	"fmt"
	"strconv"
)

type LengthValidator struct {
	expectedLength int
}

func (validator LengthValidator) Validate(field string, value interface{}) ValidationError {
	var err error
	switch v := value.(type) {
	case []string:
		for i, str := range v {
			err = validator.validateItem(str)
			if err != nil {
				field = fmt.Sprintf("%s[%d]", field, i)
				break
			}
		}
	case string:
		err = validator.validateItem(v)
	default:
		return ValidationError{
			Field: field,
			Err:   ErrValidaton{Message: "value must by type of string or []string"},
		}
	}

	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func (validator LengthValidator) validateItem(str string) error {
	var err error

	if len(str) != validator.expectedLength {
		err = ErrValidaton{
			Message: fmt.Sprintf(
				"expected size is %d, actual is %d",
				validator.expectedLength,
				len(str),
			),
		}
	}

	return err
}

func (validator LengthValidator) GetName() string {
	return "len"
}

func (validator *LengthValidator) SetParamFromString(param string) error {
	return validator.SetParam(param)
}

func (validator *LengthValidator) SetParam(param interface{}) error {
	var length int

	switch p := param.(type) {
	case string:
		l, err := strconv.Atoi(p)
		if err != nil {
			return ErrInvalidValidatorTagValue{ExpectedType: "int or numeric string", CurrentValue: param}
		}
		length = l
	case int:
		length = p
	default:
		return ErrInvalidValidatorTagValue{ExpectedType: "int or numeric string", CurrentValue: param}
	}

	validator.expectedLength = length

	return nil
}

func NewLengthValidator() *LengthValidator {
	return &LengthValidator{}
}
