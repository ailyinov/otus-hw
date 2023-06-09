package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
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
		Ports   []int  `validate:"max:5000|min:4000"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code   int    `validate:"in:200,404,500"`
		Body   string `json:"omitempty"`
		Header string `validate:"test:1251"`
	}
)

func TestValidate(t *testing.T) {
	var tests = []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     strings.Repeat("1q2w3e", 6),
				Name:   "Jay",
				Age:    40,
				Email:  "oneplusone@equalstube.id",
				Role:   "admin",
				Phones: []string{"+6844432555"},
				meta:   nil,
			},
			expectedErr: *new(ValidationErrors),
		},
		{
			in: App{
				Version: "11111",
				Ports:   []int{5000, 4000, 4100, 4040},
			},
			expectedErr: *new(ValidationErrors),
		},
		{
			in: Token{
				Header:    []byte{1, 2},
				Payload:   []byte{1, 2, 3},
				Signature: []byte{3, 21, 2},
			},
			expectedErr: *new(ValidationErrors),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(&tt.in)
			require.Equal(t, tt.expectedErr, err)
			_ = tt
		})
	}
}

func TestValidateErr(t *testing.T) {
	var tests = []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     strings.Repeat("1q2w3e1", 6),
				Name:   "Jay",
				Age:    17,
				Email:  "wrong-mail",
				Role:   "wrong-role",
				Phones: []string{"+6844432555dd"},
				meta:   nil,
			},
			expectedErr: ErrNotValidValue,
		},
		{
			in: App{
				Version: "1111111",
				Ports:   []int{5100, 4000, 4100, 4040},
			},
			expectedErr: ErrNotValidValue,
		},
		{
			in: Response{
				Code:   500,
				Body:   "body",
				Header: "header",
			},
			expectedErr: ErrNoValidator,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(&tt.in)
			if errs, ok := err.(ValidationErrors); ok {
				for _, e := range errs {
					require.True(t, errors.Is(e.Err, tt.expectedErr))
				}
			}
			_ = tt
		})
	}
}
