package validator

import (
	"errors"
	"retarget/pkg/entity"

	"github.com/go-playground/validator/v10"
	"gopkg.in/inf.v0"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("gt_decimal_01", func(fl validator.FieldLevel) bool {
		decimal, ok := fl.Field().Interface().(entity.Decimal)
		if !ok || decimal.Dec == nil {
			return false
		}
		threshold := inf.NewDec(1, 1) // 0.1
		return decimal.Cmp(threshold) == 1
	})
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
