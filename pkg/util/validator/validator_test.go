package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type mockTool struct {
	Name       string `validate:"required"`
	InstanceID string `validate:"required,dns1123subdomain"`
	URL        string `validate:"url"`
	TestOne    string `validate:"required_without=TestTwo"`
	TestTwo    string `validate:"required_without=TestOne"`
}

type mockFieldLeveler struct {
	validator.FieldLevel
	field string
}

func (fl *mockFieldLeveler) Field() reflect.Value {
	return reflect.ValueOf(fl.field)
}

var _ = Describe("Check func", func() {
	var (
		mockTest *mockTool
	)
	When("all field is empty", func() {
		BeforeEach(func() {
			mockTest = &mockTool{}
		})
		It("should return all field error", func() {
			errs := CheckStructError(mockTest)
			Expect(len(errs)).Should(Equal(5))
			Expect(errs.Combine().Error()).Should(Equal("config options are not valid:\n  field mockTool.Name is required\n  field mockTool.InstanceID is required\n  field mockTool.URL is a not valid url\n  field mockTool.TestOne validation failed on the 'required_without' tag\n  field mockTool.TestTwo validation failed on the 'required_without' tag"))
		})
	})
	When("all field is valid", func() {
		BeforeEach(func() {
			mockTest = &mockTool{
				Name:       "test",
				InstanceID: "test",
				URL:        "http://www.com",
				TestOne:    "without TestTwo",
			}
		})
		It("should return empty", func() {
			errs := CheckStructError(mockTest)
			Expect(errs).Should(BeEmpty())
		})
	})
})

var _ = Describe("dns1123SubDomain func", func() {
	var (
		testVal *mockFieldLeveler
	)
	When("value is valid", func() {
		BeforeEach(func() {
			testVal = &mockFieldLeveler{field: "a"}
		})
		It("should return true", func() {
			Expect(dns1123SubDomain(testVal)).Should(BeTrue())
		})
	})
	When("valid is invalid", func() {
		BeforeEach(func() {
			testVal = &mockFieldLeveler{field: ""}
		})
		It("should return false", func() {
			Expect(dns1123SubDomain(testVal)).Should(BeFalse())
		})
	})
})

var _ = Describe("isYaml func", func() {
	var (
		testVal *mockFieldLeveler
	)
	When("value is valid", func() {
		BeforeEach(func() {
			testVal = &mockFieldLeveler{field: `
name: Martin D'vloper`}
		})
		It("should return true", func() {
			Expect(isYaml(testVal)).Should(BeTrue())
		})
	})
	When("valid is invalid", func() {
		BeforeEach(func() {
			testVal = &mockFieldLeveler{field: `
---
# An employee record
job: Developer
skill:Elite
foods:
- Apple
- Orange
languages:
perl: Elite
python: Elite
education: ||
4 GCSEs
`}
		})
		It("should return false", func() {
			Expect(isYaml(testVal)).Should(BeFalse())
		})
	})
})
