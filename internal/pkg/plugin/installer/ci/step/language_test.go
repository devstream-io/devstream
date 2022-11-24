package step

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Language struct", func() {
	var (
		l *language
	)
	When("lanuage is go", func() {
		BeforeEach(func() {
			l = &language{
				Name: "go",
			}
		})
		It("should return go default config", func() {
			defaultConfig := l.getGeneralDefaultOption()
			Expect(defaultConfig).ShouldNot(BeNil())
			Expect(defaultConfig.testOption.Command).Should(Equal([]string{"go test ./..."}))
		})
	})
	When("lanuage is java", func() {
		BeforeEach(func() {
			l = &language{
				Name: "java",
			}
		})
		It("should return java default config", func() {
			defaultConfig := l.getGeneralDefaultOption()
			Expect(defaultConfig).ShouldNot(BeNil())
			Expect(defaultConfig.testOption.Command).Should(Equal([]string{"mvn -B test"}))
		})
	})
	When("language is not exist", func() {
		BeforeEach(func() {
			l = &language{
				Name: "not_exist",
			}
		})
		It("should return go default config", func() {
			defaultConfig := l.getGeneralDefaultOption()
			Expect(defaultConfig).Should(BeNil())
		})
	})
})
