package handler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

func validationMessage(err error) string {
	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) {
		return "invalid request body"
	}

	messages := make([]string, 0, len(fieldErrors))
	for _, fe := range fieldErrors {
		switch fe.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", fe.Field()))
		case "datetime":
			messages = append(messages, fmt.Sprintf("%s must be a valid date in YYYY-MM-DD format", fe.Field()))
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", fe.Field()))
		}
	}

	return strings.Join(messages, "; ")
}
