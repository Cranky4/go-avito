package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type InStringsValidator struct {
	expectedStrings []string
}

func (validator InStringsValidator) Validate(field string, value interface{}) ValidationError {
	var err error
	switch v := value.(type) {
	case []string:
		for i, val := range v {
			err = validator.validateItem(val)
			if err != nil {
				field = fmt.Sprintf("%s[%d]", field, i)
				break
			}
		}
	case string:
		err = validator.validateItem(v)
	case fmt.Stringer:
		err = validator.validateItem(v.String()) // ?
	case int:
		err = validator.validateItem(strconv.Itoa(v))
	case []int:
		for i, val := range v {
			err = validator.validateItem(strconv.Itoa(val))
			if err != nil {
				field = fmt.Sprintf("%s[%d]", field, i)
				break
			}
		}
	default:
		rf := reflect.ValueOf(value)
		if rf.Kind().String() == "string" {
			err = validator.validateItem(rf.String())
		} else {
			err = ErrValidaton{Message: "value must by type of string"}
		}
	}

	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func (validator InStringsValidator) validateItem(str string) error {
	for _, expectedString := range validator.expectedStrings {
		if expectedString == str {
			return nil
		}
	}

	return ErrValidaton{Message: fmt.Sprintf("must be one of %s, actual is %s", validator.expectedStrings, str)}
}

func (validator InStringsValidator) GetName() string {
	return "in"
}

func (validator *InStringsValidator) SetParamFromString(param string) error {
	return validator.SetParam(strings.Split(param, ","))
}

func (validator *InStringsValidator) SetParam(param interface{}) error {
	strings, ok := param.([]string)

	if !ok {
		return ErrInvalidValidatorTagValue{ExpectedType: "[]string", CurrentValue: param}
	}

	validator.expectedStrings = strings

	return nil
}

func NewInStringsValidator() *InStringsValidator {
	return &InStringsValidator{}
}
