package hw09structvalidator

import (
	"fmt"
)

type MaxValidator struct {
	NumberValidator
	max int
}

func (validator MaxValidator) Validate(field string, value interface{}) ValidationError {
	return validator.NumberValidator.Validate(field, value, validator.validateItem)
}

func (validator MaxValidator) validateItem(num int) error {
	var err error
	if num > validator.max {
		err = ErrValidaton{
			Message: fmt.Sprintf(
				"cannot be greater than %d, actual is %d",
				validator.max,
				num,
			),
		}
	}

	return err
}

func (validator MaxValidator) GetName() string {
	return "max"
}

func NewMaxValidator() *MaxValidator {
	return &MaxValidator{}
}

func (validator *MaxValidator) SetParam(param interface{}) (err error) {
	num, err := validator.NumberValidator.SetParam(param)
	if err != nil {
		return
	}

	validator.max = num

	return
}

func (validator *MaxValidator) SetParamFromString(param string) error {
	return validator.SetParam(param)
}
