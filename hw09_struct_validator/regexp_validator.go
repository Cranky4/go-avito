package hw09structvalidator

import (
	"fmt"
	"regexp"
)

type RegexpValidator struct {
	regexp string
}

func (v RegexpValidator) Validate(field string, value interface{}) ValidationError {
	str, ok := value.(string)

	if !ok {
		return ValidationError{
			Field: field,
			Err:   ErrValidaton{Message: "value must by type of string"},
		}
	}

	var err error

	expression, err := regexp.Compile(v.regexp)
	if err != nil {
		return ValidationError{
			Field: field,
			Err:   err,
		}
	}

	if !expression.Match([]byte(str)) {
		message := fmt.Sprintf("invalid format for %s", str)
		err = ErrValidaton{Message: message}
	}

	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func (v RegexpValidator) GetName() string {
	return "regexp"
}

func (v *RegexpValidator) SetParamFromString(param string) error {
	return v.SetParam(param)
}

func (v *RegexpValidator) SetParam(param interface{}) error {
	str, ok := param.(string)

	if !ok {
		return ErrInvalidValidatorTagValue{ExpectedType: "string", CurrentValue: param}
	}

	v.regexp = str

	return nil
}

func NewRegexpValidator() *RegexpValidator {
	return &RegexpValidator{}
}
