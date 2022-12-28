package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/config"
)

var _ = Describe("Language struct", func() {
	var (
		l config.LanguageOption
	)
	Context("GetGeneralDefaultOpt func", func() {
		When("lanuage is go", func() {
			BeforeEach(func() {
				l = config.LanguageOption{
					Name: "go",
				}
			})
			It("should return go default config", func() {
				defaultConfig := l.GetGeneralDefaultOpt()
				Expect(defaultConfig).ShouldNot(BeNil())
				Expect(defaultConfig.Test.Command).Should(Equal([]string{"go test ./..."}))
			})
		})
		When("lanuage is java", func() {
			BeforeEach(func() {
				l = config.LanguageOption{
					Name: "java",
				}
			})
			It("should return java default config", func() {
				defaultConfig := l.GetGeneralDefaultOpt()
				Expect(defaultConfig).ShouldNot(BeNil())
				Expect(defaultConfig.Test.Command).Should(Equal([]string{"mvn -B test"}))
			})
		})
		When("language is not exist", func() {
			BeforeEach(func() {
				l = config.LanguageOption{
					Name: "not_exist",
				}
			})
			It("should return go default config", func() {
				defaultConfig := l.GetGeneralDefaultOpt()
				Expect(defaultConfig).Should(BeNil())
			})
		})
	})
	Context("IsConfigured func", func() {
		When("is empty", func() {
			BeforeEach(func() {
				l = config.LanguageOption{}
			})
			It("should return false", func() {
				Expect(l.IsConfigured()).Should(BeFalse())
			})
		})
		When("not empty", func() {
			BeforeEach(func() {
				l = config.LanguageOption{Name: "test"}
			})
			It("should return true", func() {
				Expect(l.IsConfigured()).Should(BeTrue())
			})
		})
	})
})
