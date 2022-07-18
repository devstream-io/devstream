package validator

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
)

type Tool struct {
	Name       string                 `yaml:"name" validate:"required"`
	InstanceID string                 `yaml:"instanceID" validate:"required,dns1123subdomain"`
	DependsOn  []string               `yaml:"dependsOn"`
	Options    map[string]interface{} `yaml:"options"`
}

func TestStruct(t *testing.T) {
	tests := []struct {
		name         string
		s            interface{}
		wantErrCount int
	}{
		// TODO: Add test cases.
		{"base", struct{}{}, 0},
		{"base Tool instance", Tool{}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Struct(tt.s); len(got) != tt.wantErrCount {
				t.Errorf("Struct() = %v\n, got err count: %d\n, want err count:%d", got, len(got), tt.wantErrCount)
			}
		})
	}
}

type FakerFieldLeveler struct {
	validator.FieldLevel
	field string
}

func (fl *FakerFieldLeveler) Field() reflect.Value {
	return reflect.ValueOf(fl.field)
}

func Test_dns1123SubDomain(t *testing.T) {
	goodValues, badValues := "a", ""
	tests := []struct {
		name string
		fl   validator.FieldLevel
		want bool
	}{
		// TODO: Add test cases.
		{"base", &FakerFieldLeveler{field: goodValues}, true},
		{"base", &FakerFieldLeveler{field: badValues}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dns1123SubDomain(tt.fl); got != tt.want {
				t.Errorf("dns1123SubDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
