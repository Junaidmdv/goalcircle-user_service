package validater

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func  phoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// E.164 format: + followed by 1–15 digits
	regex := `^\+[1-9]\d{1,14}$`

	match, _ := regexp.MatchString(regex, phone)
	return match
}
