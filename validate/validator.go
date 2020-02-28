package validate

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate
var translator ut.Translator

func init() {
	validate = validator.New()
}

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func Struct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		errs := err.(validator.ValidationErrors)
		validationErrors := errs.Translate(translator)

		var errors []string
		for _, v := range validationErrors {
			errors = append(errors, v)
		}

		return formattedError{
			messages: errors,
			cause:    err,
		}
	}

	return nil
}

type formattedError struct {
	messages []string
	cause    error
}

func (f formattedError) Error() string {
	return strings.Join(f.messages, ", ")
}

func (f formattedError) Unwrap() error {
	return f.cause
}
