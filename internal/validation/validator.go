package validation

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return ""
		}
		return name
	})

	Validate.RegisterValidation("rfc3339", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		_, err := time.Parse(time.RFC3339, fl.Field().String())
		return err == nil
	})
}
