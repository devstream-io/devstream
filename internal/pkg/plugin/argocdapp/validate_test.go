package argocdapp

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("setDefault func", func() {
	var option configmanager.RawOptions

	BeforeEach(func() {
		option = configmanager.RawOptions{
			"imageRepo": map[string]any{
				"user": "test_user",
				"url":  "http://test.com",
			},
		}
	})
	It("should return with default value", func() {
		var appOption *app
		var sourceOption *source
		var destinationOption *destination
		opt, err := setDefault(option)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(opt).Should(Equal(configmanager.RawOptions{
			"app":         appOption,
			"destination": destinationOption,
			"source":      sourceOption,
			"imageRepo": map[string]any{
				"url":       "http://test.com/",
				"user":      "test_user",
				"initalTag": "0.0.1",
			},
		}))
	})
})
