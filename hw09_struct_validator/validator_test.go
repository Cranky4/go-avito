package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	InvalidValidatorTag struct {
		Code int `validate:"zz:200,404,500"`
	}

	InvalidLenTag struct {
		Code int `validate:"len:ten"`
	}

	InvalidMaxTag struct {
		Code int `validate:"max:six"`
	}

	InvalidMinTag struct {
		Code int `validate:"min:seven"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "unexpected_string",
			expectedErr: ErrInvalidInputArgument,
		},
		{
			in:          InvalidValidatorTag{Code: 415},
			expectedErr: ErrInvalidValidatorTag{},
		},
		{
			in:          InvalidLenTag{Code: 415},
			expectedErr: ErrInvalidValidatorTagValue{},
		},
		{
			in:          InvalidMaxTag{Code: 415},
			expectedErr: ErrInvalidValidatorTagValue{},
		},
		{
			in:          InvalidMinTag{Code: 415},
			expectedErr: ErrInvalidValidatorTagValue{},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			assert.NotNil(t, err)
			_ = tt
		})
	}
	validationTests := []struct {
		in             interface{}
		expectedErrors string
	}{
		{
			in: User{
				ID:     "short-id", // too short
				Name:   "John",
				Age:    32,
				Email:  "john@smith", // invalid
				Role:   "admin",
				Phones: []string{"88005553535"},
				meta:   nil,
			},
			expectedErrors: "ID: expected size is 36, actual is 8\nEmail: invalid format for john@smith\n",
		},
		{
			in: User{
				ID:     "5f56fb38-4ba3-40b3-98d8-3119c7061a86",
				Name:   "John",
				Age:    12, // too young
				Email:  "john@smith.com",
				Role:   "stuff",
				Phones: []string{"88005553535"},
				meta:   nil,
			},
			expectedErrors: "Age: cannot be less than 18, actual is 12\n",
		},
		{
			in: User{
				ID:     "5f56fb38-4ba3-40b3-98d8-3119c7061a86",
				Name:   "John",
				Age:    99, // too old
				Email:  "john@smith.com",
				Role:   "stuff",
				Phones: []string{"880055535"}, // too short
				meta:   nil,
			},
			expectedErrors: "Age: cannot be greater than 50, actual is 99\nPhones[0]: expected size is 11, actual is 9\n",
		},
		{
			in: User{
				ID:     "5f56fb38-4ba3-40b3-98d8-3119c7061a86",
				Name:   "John",
				Age:    33,
				Email:  "john@smith.com",
				Role:   "customer", // invalid
				Phones: []string{"88005553535"},
				meta:   nil,
			},
			expectedErrors: "Role: must be one of [admin stuff], actual is customer\n",
		},
		{
			in: Response{
				Code: 232,
				Body: "What?",
			},
			expectedErrors: "Code: must be one of [200 404 500], actual is 232\n",
		},
		{
			in: App{
				Version: "v1.23.23-alpha",
			},
			expectedErrors: "Version: expected size is 5, actual is 14\n",
		},
	}
	for i, tt := range validationTests {
		t.Run(fmt.Sprintf("validation case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			assert.True(t, errors.As(err, &ValidationErrors{}))
			assert.Equal(t, tt.expectedErrors, err.Error())

			_ = tt
		})
	}
	validTests := []interface{}{
		User{
			ID:     "5f56fb38-4ba3-40b3-98d8-3119c7061a86",
			Name:   "John",
			Age:    32,
			Email:  "john@smith.com",
			Role:   "admin",
			Phones: []string{"88005553535"},
			meta:   nil,
		},
		App{
			Version: "1.0.1",
		},
		Token{
			Header:    []byte("Header"),
			Payload:   []byte("Payload"),
			Signature: []byte("Signature"),
		},
		Response{
			Code: 200,
			Body: "OK",
		},
	}
	for i, tt := range validTests {
		t.Run(fmt.Sprintf("validation case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)

			assert.Nil(t, err)

			_ = tt
		})
	}
}
