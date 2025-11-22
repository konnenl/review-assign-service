package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()
	_ = v.RegisterValidation("is_active", validateIsActive)
	return &CustomValidator{Validator: v}

}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func validateIsActive(field validator.FieldLevel) bool {
	isActive := field.Field()
	if isActive.Kind() == reflect.Ptr {
		return !isActive.IsNil()
	}
	return true
}

func GetValidationErrors(err error) error {
	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	var msgs []string
	for _, fieldErr := range valErrors {
		msg := fmt.Sprintf("%s: %s",
			strings.ToLower(fieldErr.Field()),
			getErrorText(fieldErr),
		)
		msgs = append(msgs, msg)
	}

	return fmt.Errorf("validation failed: %s", strings.Join(msgs, ", "))
}

func getErrorText(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return "min lenght is " + err.Param()
	default:
		return "invalid"
	}
}
