package hw09structvalidator

import (
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var validatorRegistry = buildValidatorRegistry()

func (v ValidationErrors) Error() string {
	var builder strings.Builder

	for _, validationnError := range v {
		if validationnError.Err.Error() != "" {
			builder.WriteString(validationnError.Field)
			builder.WriteString(": ")
			builder.WriteString(validationnError.Err.Error())
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func Validate(v interface{}) error {
	iv := reflect.ValueOf(v)

	if iv.Kind().String() != "struct" {
		return ErrInvalidInputArgument
	}

	fieldsCount := iv.NumField()
	ivType := iv.Type()

	validationErrors := make(ValidationErrors, 0, fieldsCount)

	for i := 0; i < fieldsCount; i++ {
		fieldReflect := ivType.Field(i)

		fieldTag, tagExists := fieldReflect.Tag.Lookup("validate")

		if !tagExists {
			continue
		}

		manyValidatorsTags := strings.Split(fieldTag, "|")

		for _, validatorTag := range manyValidatorsTags {
			if !strings.ContainsRune(validatorTag, ':') {
				return ErrInvalidValidatorTag{Tag: fieldTag, Field: fieldReflect.Name}
			}

			validatorData := strings.SplitN(validatorTag, ":", 2)
			validator, err := prepareValidator(validatorData)
			if err != nil {
				return err
			}

			validationError := validator.Validate(fieldReflect.Name, iv.Field(i).Interface())

			if validationError.Err != nil {
				validationErrors = append(validationErrors, validationError)
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func buildValidatorRegistry() ValidatorRegistry {
	return NewValidatorRegistry([]Validator{
		NewLengthValidator(),
		NewMinValidator(),
		NewMaxValidator(),
		NewRegexpValidator(),
		NewInStringsValidator(),
	})
}

func prepareValidator(validatorData []string) (Validator, error) {
	validator, validatorExists := validatorRegistry[validatorData[0]]
	if !validatorExists {
		return nil, ErrValidatorNotExists{ValidatorName: validatorData[0]}
	}

	err := validator.SetParamFromString(validatorData[1])
	if err != nil {
		return nil, err
	}

	return validator, nil
}
