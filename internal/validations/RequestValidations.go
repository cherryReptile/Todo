package validations

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func CreatingValidate(i interface{}) error {
	validate = validator.New()
	err := validate.Struct(i)

	if err != nil {
		err = errorIterate(err)
	}

	return err
}

func UpdatingValidate(i interface{}) error {
	validate = validator.New()
	err := validate.Struct(i)
	err = errorIterate(err)

	return err
}

func errorIterate(err error) error {
	for _, errFields := range err.(validator.ValidationErrors) {
		err = errors.New(fmt.Sprintf("%v field is %v", errFields.StructField(), errFields.Tag()))
	}
	return err
}
