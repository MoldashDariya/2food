package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

func ToMap(err error) map[string]string {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	errs := make(map[string]string)
	for _, fieldError := range validationErrors {
		errs[fieldError.Field()] = msgForTag(fieldError.Tag(), fieldError.Param())
	}

	return errs
}

func msgForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("must be at least %s characters", param)
	case "max":
		return fmt.Sprintf("must be at most %s characters", param)
	case "e164":
		return "invalid phone number format"
	case "len":
		return fmt.Sprintf("must be exactly %s characters long", param)
	default:
		return fmt.Sprintf("invalid value (failed on '%s')", tag)
	}
}
