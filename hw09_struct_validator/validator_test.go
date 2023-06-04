package hw09structvalidator

import (
	"encoding/json"
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
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     strings.Repeat("1q2w3e", 6),
				Name:   "Jay",
				Age:    45,
				Email:  "oneplusone@equalstube.id",
				Role:   "Surfer",
				Phones: []string{"+6844432555"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "12",
				Ports:   []int{5001, 5002, 50003, 4040},
			},
		},
		// ...
		// Place your code here.
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
