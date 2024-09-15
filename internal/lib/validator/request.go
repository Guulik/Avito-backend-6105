package validator

import (
	"github.com/go-playground/validator/v10"
)

func Validate(request interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(request)
	if err != nil {
		return err
	}
	return nil
}
