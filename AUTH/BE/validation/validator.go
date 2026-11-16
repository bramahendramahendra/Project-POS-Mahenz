package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
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
}
