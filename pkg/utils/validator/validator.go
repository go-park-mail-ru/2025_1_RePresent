package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) (string, error) {
	err := validate.Struct(s)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var errorMessages string
			for _, verr := range validationErrors {
				errorMessages += verr.Field() + ": " + verr.Tag() + " \n "
			}
			return errorMessages, err
		}
	}
	return "", nil
}
