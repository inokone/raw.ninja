package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationMessage is a function to convert Gin-Gonic validation errors to `StatusMessage`.
func ValidationMessage(err error) StatusMessage {
	return StatusMessage{
		Code:    400,
		Message: strings.Join(ValidationMessages(err), " "),
	}
}

// ValidationMessages is a function to convert Gin-Gonic validation error to human readable.
func ValidationMessages(err error) []string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		out := make([]string, len(ve))
		for i, fe := range ve {
			out[i] = asMsg(fe)
		}
		return out
	} else if je, ok := err.(*json.UnmarshalTypeError); ok {
		return []string{fmt.Sprintf("The field %s must be a %s", je.Field, je.Type.String())}
	}
	return nil
}

func asMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s field is required.", fe.Field())
	case "len":
		return fmt.Sprintf("%s should be at least %s characters.", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("%s should be an email.", fe.Field())
	}
	return "Unknown error"
}
