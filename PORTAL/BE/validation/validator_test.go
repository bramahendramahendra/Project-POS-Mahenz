package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

// TestValidatorInitialized ensures the validator is properly initialized
func TestValidatorInitialized(t *testing.T) {
	if Validate == nil {
		t.Fatal("Validate should not be nil after init()")
	}
	if ErrTrans == nil {
		t.Fatal("ErrTrans should not be nil after init()")
	}
}

// Test json tag name extraction
type JSONTagTestStruct struct {
	UserName string `json:"user_name" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required"`
	Ignored  string `json:"-" validate:"required"`
}

func TestJSONTagNameFunc(t *testing.T) {
	t.Run("validation errors use json tag names", func(t *testing.T) {
		s := JSONTagTestStruct{UserName: "", Email: "", Ignored: ""}
		err := Validate.Struct(s)
		if err == nil {
			t.Fatal("expected validation errors")
		}

		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			t.Fatal("expected validator.ValidationErrors type")
		}

		// Check that field names from json tags are used
		for _, ve := range validationErrors {
			field := ve.Field()
			if field == "UserName" {
				t.Error("expected json tag name 'user_name' but got struct field 'UserName'")
			}
			if field == "Email" {
				t.Error("expected json tag name 'email' but got struct field 'Email'")
			}
		}
	})
}

// Test error translation
func TestErrorTranslation(t *testing.T) {
	t.Run("required field translation", func(t *testing.T) {
		type TestStruct struct {
			Name string `json:"name" validate:"required"`
		}
		s := TestStruct{Name: ""}
		err := Validate.Struct(s)
		if err != nil {
			validationErrors, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("expected validator.ValidationErrors type")
			}
			translated := validationErrors[0].Translate(ErrTrans)
			if translated == "" {
				t.Error("translation should not be empty")
			}
		}
	})

	t.Run("nefield translation", func(t *testing.T) {
		type TestStruct struct {
			Password        string `json:"password" validate:"required"`
			ConfirmPassword string `json:"confirm_password" validate:"nefield=Password"`
		}
		// Same password should fail nefield
		s := TestStruct{Password: "secret", ConfirmPassword: "secret"}
		err := Validate.Struct(s)
		if err != nil {
			validationErrors, ok := err.(validator.ValidationErrors)
			if !ok {
				t.Fatal("expected validator.ValidationErrors type")
			}
			for _, ve := range validationErrors {
				if ve.Tag() == "nefield" {
					translated := ve.Translate(ErrTrans)
					if translated == "" {
						t.Error("nefield translation should not be empty")
					}
				}
			}
		}
	})
}
