package validater

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validater struct {
	vn *validator.Validate
	ut ut.Translator
}

func NewValidater() (*Validater, error) {

	enlocal := en.New()
	uni := ut.New(enlocal, enlocal)
	engtrans, _ := uni.GetTranslator("en")
	validater := validator.New()
	if err := en_translations.RegisterDefaultTranslations(validater, engtrans); err != nil {
		return nil, err
	}
	validater.RegisterValidation("phone", phoneValidation)
	validater.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name != "-" {
			return name
		}

		return ""
	})

	validater.RegisterTranslation("phone", engtrans,
		func(ut ut.Translator) error {
			return ut.Add("phone", "{0} must be valid phone number", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("phone", fe.Field())
			return t
		},
	)
	
	return &Validater{
		vn: validater,
		ut: engtrans,
	}, nil
}

func (v *Validater) Validation(input interface{}) validator.ValidationErrorsTranslations {

	err := v.vn.Struct(input)
	if err == nil {
		return nil
	}

	translated := err.(validator.ValidationErrors).Translate(v.ut)
	return translated
}
