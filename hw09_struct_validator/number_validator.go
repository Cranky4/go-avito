package hw09structvalidator

import (
	"fmt"
	"strconv"
)

type NumberValidator struct{}

func (validator NumberValidator) Validate(
	field string,
	value interface{},
	callback func(num int) error,
) ValidationError {
	var err error
	switch v := value.(type) {
	case []int:
		for i, num := range v {
			err = callback(num)
			if err != nil {
				field = fmt.Sprintf("%s[%d]", field, i)
				break
			}
		}
	case int:
		err = callback(v)
	default:
		return ValidationError{
			Field: field,
			Err:   ErrValidaton{Message: "value must by type of int or []int"},
		}
	}

	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func (validator *NumberValidator) SetParam(param interface{}) (int, error) {
	var num int

	switch p := param.(type) {
	case string:
		n, err := strconv.Atoi(p)
		if err != nil {
			return num, ErrInvalidValidatorTagValue{ExpectedType: "int or numeric string", CurrentValue: param}
		}
		num = n
	case int:
		num = p
	default:
		return num, ErrInvalidValidatorTagValue{ExpectedType: "int or numeric string", CurrentValue: param}
	}

	return num, nil
}
