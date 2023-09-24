package hw09structvalidator

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
)

func TestValidate(t *testing.T) { //nolint: funlen
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name:        "structn't",
			in:          false,
			expectedErr: ErrValNotStruct,
		},
		{
			name: "supportn't field",
			in: struct {
				NoSupport byte `validate:"max:42"`
			}{
				NoSupport: 'w',
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "NoSupport",
					Err:   ErrValUnsupportedType,
				},
			},
		},
		{
			name: "supportn't slice field",
			in: struct {
				NoSupport []byte `validate:"len:32"`
			}{
				NoSupport: []byte{'w'},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "NoSupport",
					Err:   ErrValUnsupportedType,
				},
			},
		},
		{
			name: "invalid tag slice",
			in: struct {
				Value string `validate:"in:"`
			}{
				Value: "for",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValIncorrectRule,
				},
			},
		},
		{
			name: "invalid tag len",
			in: struct {
				Value string `validate:"len:"`
			}{
				Value: "for",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValIncorrectRule,
				},
			},
		},
		{
			name: "invalid tag regex",
			in: struct {
				Value string `validate:"regexp:l-aksfk(as[fk,oasm$fopasmf"`
			}{
				Value: "for",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValIncorrectTagRegexPattern,
				},
			},
		},
		{
			name: "invalid tag min",
			in: struct {
				Value int `validate:"min:z"`
			}{
				Value: 4,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValTagValueShouldBeInteger,
				},
			},
		},
		{
			name: "invalid tag max",
			in: struct {
				Value int `validate:"max:z"`
			}{
				Value: 4,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Value",
					Err:   ErrValTagValueShouldBeInteger,
				},
			},
		},
		{
			name: "invalid user",
			in: User{
				ID:     strings.Repeat("a", 100),
				Name:   "first name",
				Age:    5,
				Email:  "email at sign cool dot com",
				Role:   "guest",
				Phones: []string{"8800"},
				meta:   []byte{'1', '2', '3'},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValIncorrectStringLength,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValIntLessThanMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValRegexString,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValNoMatchingElementInSlice,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValIncorrectStringLength,
				},
			},
		},
		{
			name: "token",
			in: Token{
				Header:    []byte{1, 2, 3, 4},
				Payload:   []byte{1, 2, 3, 4},
				Signature: []byte{1, 2, 3, 4},
			},
		},
		{
			name: "invalid response",
			in: Response{
				Code: 418,
				Body: "I'm a teapot",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValNoMatchingElementInSlice,
				},
			},
		},
		{
			name: "invalid app",
			in: App{
				Version: "2.0.19",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValIncorrectStringLength,
				},
			},
		},
		{
			name: "valid response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
		},
		{
			name: "valid app",
			in: App{
				Version: "1.0.0",
			},
		},
		{
			name: "valid user",
			in: User{
				ID:     strings.Repeat("a", 36),
				Name:   "Firstname Lastname",
				Age:    20,
				Email:  "email@cool.com",
				Role:   "admin",
				Phones: []string{"88005553535"},
				meta:   []byte{'1', '2', '3'},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
