package file

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("replaceAppNameInPathStr func", func() {
	var (
		placeHolder string
		filePath    string
		appName     string
	)
	BeforeEach(func() {
		placeHolder = "__app__"
		appName = "test"
	})
	When("filePath not contains placeHolder", func() {
		BeforeEach(func() {
			filePath = "/app/dev"
		})
		It("should return same filePath", func() {
			newPath, err := replaceAppNameInPathStr(filePath, placeHolder, appName)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(newPath).Should(Equal(filePath))
		})
	})
	When("filPath contains placeHolder", func() {
		BeforeEach(func() {
			filePath = fmt.Sprintf("app/%s/dev", placeHolder)
		})
		It("should replace placeHolder with app name", func() {
			newPath, err := replaceAppNameInPathStr(filePath, placeHolder, appName)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(newPath).Should(Equal(fmt.Sprintf("app/%s/dev", appName)))
		})
	})
})
