package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
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
			in:          123,
			expectedErr: ErrNotStruct,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in:          Response{200, ""},
			expectedErr: nil,
		},
		{
			in:          Response{403, "forbidden"},
			expectedErr: ValidationError{"Code", ErrValueNotInSet},
		},
		{
			in:          App{Version: "abc"},
			expectedErr: ValidationError{"Version", ErrValueLength},
		},
		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "d5e66c85-f3cf-4281-916e-86b98d513fcc",
				Name:   "Foo Bar",
				Age:    30,
				Email:  "example@test.com",
				Role:   "baz",
				Phones: nil,
				meta:   nil,
			},
			expectedErr: ValidationError{"Role", ErrValueNotInSet},
		},
		{
			in: User{
				ID:     "d5e66c85-f3cf-4281-916e-86b98d513fcc",
				Name:   "Foo Bar",
				Age:    30,
				Email:  "example@test.com",
				Role:   "admin",
				Phones: []string{"+7999111222", "+999111222"},
				meta:   nil,
			},
			expectedErr: ValidationError{"Phones[1]", ErrValueLength},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			var validationErr ValidationError
			if errors.As(tt.expectedErr, &validationErr) {
				var got ValidationError
				if !errors.As(err, &got) {
					t.Fatalf("expected %v, got %v", validationErr, err)
				}
				if validationErr.Field != got.Field {
					t.Fatalf("expected field %v, got %v", validationErr.Field, got.Field)
				}
				if !errors.Is(validationErr.Err, got.Err) {
					t.Fatalf("expected err %v, got %v", validationErr.Err, got.Err)
				}
				return
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
