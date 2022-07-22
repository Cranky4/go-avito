package hw09structvalidator

import (
	"fmt"
)

type MinValidator struct {
	NumberValidator
	min int
}

func (validator MinValidator) Validate(field string, value interface{}) ValidationError {
	return validator.NumberValidator.Validate(field, value, validator.validateItem)
}

func (validator MinValidator) validateItem(num int) error {
	var err error
	if num < validator.min {
		err = ErrValidaton{
			Message: fmt.Sprintf(
				"cannot be less than %d, actual is %d",
				validator.min,
				num,
			),
		}
	}

	return err
}

func (validator MinValidator) GetName() string {
	return "min"
}

func (validator *MinValidator) SetParam(param interface{}) (err error) {
	num, err := validator.NumberValidator.SetParam(param)
	if err != nil {
		return
	}

	validator.min = num

	return
}

func (validator *MinValidator) SetParamFromString(param string) error {
	return validator.SetParam(param)
}

func NewMinValidator() *MinValidator {
	return &MinValidator{}
}
