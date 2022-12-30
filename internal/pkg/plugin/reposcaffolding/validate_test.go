package reposcaffolding

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("validate func", func() {
	var (
		rawOpts configmanager.RawOptions
	)
	When("reposcaffolding option is not valid", func() {
		BeforeEach(func() {
			rawOpts = configmanager.RawOptions{
				"not_exist": true,
			}
		})
		It("should return err", func() {
			_, err := validate(rawOpts)
			Expect(err).Should(HaveOccurred())
		})
	})
	When("reposcaffolding option is valid", func() {
		BeforeEach(func() {
			rawOpts = configmanager.RawOptions{
				"sourceRepo": map[string]string{
					"owner":   "test_user",
					"name":    "test_repo",
					"scmType": "github",
					"branch":  "main",
				},
				"destinationRepo": map[string]string{
					"owner":   "dst_user",
					"name":    "dst_repo",
					"scmType": "github",
					"branch":  "main",
				},
			}
		})
		It("should return noraml", func() {
			_, err := validate(rawOpts)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
})
