package hw09structvalidator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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

func TestValidateUser(t *testing.T) {
	t.Run("testValidUser", func(t *testing.T) {
		input := User{
			ID:     "012345678901234567890123456789012345",
			Name:   "name01234567890123456789",
			Age:    33,
			Email:  "test@mail.ru",
			Role:   "stuff",
			Phones: []string{"12334567891"},
		}
		err := Validate(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		require.Nil(t, err)
	})
}

func TestValidateUserErrors(t *testing.T) {
	t.Run("testValidUserErrors", func(t *testing.T) {
		in := User{
			ID:     "0",
			Name:   "0",
			Age:    0,
			Email:  "0",
			Role:   "0",
			Phones: []string{"0"},
		}
		errs := ValidationErrors{
			ValidationError{
				Field: "ID",
				Err:   ErrStringLengthMismatch(36),
			},
			ValidationError{
				Field: "Age",
				Err:   ErrValueOutOfRange("greater", 18),
			},
			ValidationError{
				Field: "Email",
				Err:   ErrStringDoesNotMatchPattern("^\\w+@\\w+\\.\\w+$"),
			},
			ValidationError{
				Field: "Role",
				Err:   ErrValueOutOfValues("admin,stuff"),
			},
			ValidationError{
				Field: "Phones",
				Err:   ErrStringLengthMismatch(11),
			},
		}
		err := Validate(in)
		if err == nil {
			t.Errorf("validation failed, no errors: %v", err)
		}
		require.Equal(t, errs, err)
	})
}

func TestValidateApp(t *testing.T) {
	t.Run("testValidApp", func(t *testing.T) {
		input := App{
			Version: "0.1.5",
		}
		err := Validate(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		require.Nil(t, err)
	})
}

func TestValidateAppErrors(t *testing.T) {
	t.Run("testValidAppErrors", func(t *testing.T) {
		input := App{
			Version: "0",
		}
		errs := ValidationErrors{
			ValidationError{
				Field: "Version",
				Err:   ErrStringLengthMismatch(5),
			},
		}
		err := Validate(input)
		if err == nil {
			t.Errorf("validation failed, no errors: %v", err)
		}
		require.Equal(t, errs, err)
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("testValidToken", func(t *testing.T) {
		input := Token{
			Header:    []byte("dsfdsfds"),
			Payload:   []byte(""),
			Signature: nil,
		}
		err := Validate(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		require.Nil(t, err)
	})
}

func TestValidateResponse(t *testing.T) {
	t.Run("testValidResponse", func(t *testing.T) {
		input := Response{
			Code: 200,
			Body: "",
		}
		err := Validate(input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		require.Nil(t, err)
	})
}

func TestValidateResponseErrors(t *testing.T) {
	t.Run("testValidResponseErrors", func(t *testing.T) {
		input := Response{
			Code: 100,
			Body: "dddddd",
		}
		errs := ValidationErrors{
			ValidationError{
				Field: "Code",
				Err:   ErrValueOutOfValues("200,404,500"),
			},
		}
		err := Validate(input)
		if err == nil {
			t.Errorf("validation failed, no errors: %v", err)
		}
		require.Equal(t, errs, err)
	})
}
