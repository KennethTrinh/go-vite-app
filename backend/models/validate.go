package models

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

func ValidateStruct[T any](s T) []string {
	var errors []string
	err := validate.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Field()+" does not satisfy "+err.Tag()+" condition")
		}
	}
	return errors
}
