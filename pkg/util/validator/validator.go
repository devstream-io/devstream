package validator

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/util/validation"
)

var v *validator.Validate

type StructFieldErrors []error

func (errs StructFieldErrors) Combine() error {
	if len(errs) == 0 {
		return nil
	}
	var totalErr = "config options are not valid:\n"
	for _, fieldErr := range errs {
		totalErr = fmt.Sprintf("%s  %s\n", totalErr, fieldErr)
	}
	return fmt.Errorf(strings.TrimSpace(totalErr))
}

func init() {
	v = validator.New()

	validations := []struct {
		tag string
		fn  validator.Func
	}{
		{"dns1123subdomain", dns1123SubDomain},
		{"yaml", isYaml},
	}

	for _, vt := range validations {
		if err := v.RegisterValidation(vt.tag, vt.fn); err != nil {
			log.Fatal(err)
		}
	}
	v.RegisterTagNameFunc(getMapstructureOrYamlTagName)
}

// CheckStructError will check s, and return StructValidationError if this struct is not valid
func CheckStructError(s interface{}) StructFieldErrors {
	fieldErrs := make(StructFieldErrors, 0)
	if err := v.Struct(s); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			var fieldCustomErr error
			switch fieldErr.Tag() {
			case "required":
				fieldCustomErr = fmt.Errorf("field %s is required", fieldErr.Namespace())
			case "url":
				fieldCustomErr = fmt.Errorf("field %s is a not valid url", fieldErr.Namespace())
			case "oneof":
				fieldCustomErr = fmt.Errorf("field %s must be one of [%s]", fieldErr.Namespace(), fieldErr.Param())
			default:
				fieldCustomErr = fmt.Errorf("field %s validation failed on the '%s' tag", fieldErr.Namespace(), fieldErr.Tag())
			}
			fieldErrs = append(fieldErrs, fieldCustomErr)
		}
	}
	return fieldErrs
}

func dns1123SubDomain(fl validator.FieldLevel) bool {
	return len(validation.IsDNS1123Subdomain(fl.Field().String())) == 0
}

func isYaml(fl validator.FieldLevel) bool {
	return yaml.Unmarshal([]byte(fl.Field().String()), &struct{}{}) == nil
}

func getMapstructureOrYamlTagName(fld reflect.StructField) string {
	// 1. get tag name from mapstructure or yaml
	tagName := fld.Tag.Get("mapstructure")
	if tagName == "" {
		tagName = fld.Tag.Get("yaml")
	}
	// 2. else get yaml tag name
	name := strings.SplitN(tagName, ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}
