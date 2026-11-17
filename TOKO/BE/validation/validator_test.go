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


// Test struct for dateformat validation
type DateformatTestStruct struct {
	Date string `json:"date" validate:"dateformat"`
}

func TestDateformatValidation(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{"valid date", "2025-01-15", false},
		{"valid leap year date", "2024-02-29", false},
		{"invalid leap year date", "2025-02-29", true},
		{"invalid format dd-mm-yyyy", "15-01-2025", true},
		{"invalid format mm/dd/yyyy", "01/15/2025", true},
		{"invalid month 13", "2025-13-01", true},
		{"invalid day 32", "2025-01-32", true},
		{"empty string fails regex", "", true},
		{"invalid format letters", "abcd-ef-gh", true},
		{"partial date", "2025-01", true},
		{"valid end of month", "2025-12-31", false},
		{"valid start of year", "2025-01-01", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DateformatTestStruct{Date: tt.value}
			err := Validate.Struct(s)
			if tt.expectErr && err == nil {
				t.Errorf("expected error for value %q but got none", tt.value)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error for value %q but got: %v", tt.value, err)
			}
		})
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

// Uji aturan khusus yang dipakai DTO Toko (ditambahkan saat pembersihan dead code, 17 Nov 2026)
type NotBlankTimeTestStruct struct {
	Name  string `json:"name" validate:"notblank"`
	Start string `json:"start" validate:"timeformat"`
}

func TestNotBlankAndTimeformat(t *testing.T) {
	tests := []struct {
		name      string
		input     NotBlankTimeTestStruct
		expectErr bool
	}{
		{"valid HH:MM", NotBlankTimeTestStruct{Name: "Pagi", Start: "08:00"}, false},
		{"valid HH:MM:SS", NotBlankTimeTestStruct{Name: "Pagi", Start: "08:00:30"}, false},
		{"nama hanya spasi", NotBlankTimeTestStruct{Name: "   ", Start: "08:00"}, true},
		{"jam tidak valid", NotBlankTimeTestStruct{Name: "Pagi", Start: "25:00"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.input)
			if tt.expectErr != (err != nil) {
				t.Fatalf("expectErr=%v, err=%v", tt.expectErr, err)
			}
			if err != nil {
				for _, ve := range err.(validator.ValidationErrors) {
					if msg := ve.Translate(ErrTrans); msg == "" || msg == ve.Error() {
						t.Errorf("pesan %s belum diterjemahkan: %q", ve.Tag(), msg)
					}
				}
			}
		})
	}
}
