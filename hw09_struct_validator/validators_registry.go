package hw09structvalidator

type Validator interface {
	Validate(field string, value interface{}) ValidationError
	GetName() string
	SetParam(interface{}) error
	SetParamFromString(string) error
}

type ValidatorRegistry map[string]Validator

func NewValidatorRegistry(validators []Validator) ValidatorRegistry {
	registry := make(ValidatorRegistry)
	for _, validator := range validators {
		registry[validator.GetName()] = validator
	}

	return registry
}
