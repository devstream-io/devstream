package validator

import (
	"errors"
	"log"

	"gopkg.in/yaml.v3"

	"github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/util/validation"
)

var v *validator.Validate

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
}

func Struct(s interface{}) []error {
	if err := v.Struct(s); err != nil {
		errs := make([]error, 0)
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, errors.New(e.Error()))
		}
		return errs
	}
	return nil
}

func StructAllError(s interface{}) error {
	return v.Struct(s)
}

func dns1123SubDomain(fl validator.FieldLevel) bool {
	return len(validation.IsDNS1123Subdomain(fl.Field().String())) == 0
}

func isYaml(fl validator.FieldLevel) bool {
	return yaml.Unmarshal([]byte(fl.Field().String()), &struct{}{}) == nil
}
