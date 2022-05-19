package validator

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/util/validation"
)

var v *validator.Validate

func init() {
	v = validator.New()

	if err := v.RegisterValidation("dns1123subdomain", dns1123SubDomain); err != nil {
		log.Fatal(err)
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

func dns1123SubDomain(fl validator.FieldLevel) bool {
	return len(validation.IsDNS1123Subdomain(fl.Field().String())) == 0
}
