package validation

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

type (
	CallbackValidationMessage func(validator.FieldError) string
)

var (
	Validate *validator.Validate
	ErrTrans ut.Translator
)

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	idn := id.New()
	uni := ut.New(idn, idn)
	errTrans, _ := uni.GetTranslator("id")
	ErrTrans = errTrans
	idTranslations.RegisterDefaultTranslations(Validate, ErrTrans)
	Validate.RegisterTranslation("nefield", ErrTrans,
		func(ut ut.Translator) error {
			return ut.Add("nefield", "{0} tidak boleh sama dengan {1}", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("nefield", fe.Field(), fe.Param())
			return t
		},
	)

	// 4. Date Format
	Validate.RegisterValidation("dateformat", func(f1 validator.FieldLevel) bool {
		dateStr := f1.Field().String()

		// Quick regex check for format
		regex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
		if !regex.MatchString(dateStr) {
			return false
		}

		// Parse with time to ensure valid date (not 2025-02-30)
		_, err := time.Parse("2006-01-02", dateStr)
		return err == nil
	})

	AddCustomErrorMessage("dateformat", "{0} hanya boleh menggunakan format YYYY-MM-DD", func(fe validator.FieldError) string {
		return fe.Param()
	})

	// 11. Not Blank - rejects strings that are empty after trimming whitespace
	Validate.RegisterValidation("notblank", func(f1 validator.FieldLevel) bool {
		value, ok := f1.Field().Interface().(string)
		if !ok {
			return true
		}
		return strings.TrimSpace(value) != ""
	})

	AddCustomErrorMessage("notblank", "{0} tidak boleh kosong atau hanya berisi spasi", func(fe validator.FieldError) string {
		return ""
	})

	// 12. Time Format "HH:MM" or "HH:MM:SS" (24-hour) validation
	Validate.RegisterValidation("timeformat", func(f1 validator.FieldLevel) bool {
		timeStr := f1.Field().String()

		if _, err := time.Parse("15:04:05", timeStr); err == nil {
			return true
		}
		_, err := time.Parse("15:04", timeStr)
		return err == nil
	})

	AddCustomErrorMessage("timeformat", "{0} harus dalam format jam yang valid (HH:MM atau HH:MM:SS)", func(fe validator.FieldError) string {
		return ""
	})

}

func AddCustomErrorMessage(tag string, errMessage string, cb CallbackValidationMessage) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}
	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := cb(fe)
		tag := fe.Tag()
		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	Validate.RegisterTranslation(tag, ErrTrans, registerFn, transFn)
}
