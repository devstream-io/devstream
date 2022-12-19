package general

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("validate func", func() {
	var (
		options        configmanager.RawOptions
		configLocation string
	)
	BeforeEach(func() {
		configLocation = "workflows"
		options = configmanager.RawOptions{
			"pipeline": map[string]interface{}{
				"configLocation": configLocation,
			},
			"scm": map[string]interface{}{
				"scmType":  "github",
				"owner":    "test",
				"name":     "gg",
				"branch":   "main",
				"needAuth": true,
			},
		}
	})

	When("scm repo is gitlab", func() {
		BeforeEach(func() {
			options["scm"] = map[string]interface{}{
				"url": "http://exmaple.gitlab.com",
			}
		})
		It("should return error", func() {
			_, err := validate(options)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("all valid", func() {
		BeforeEach(func() {
			os.Setenv("GITHUB_TOKEN", "test")
		})
		It("should not raise error", func() {
			_, err := validate(options)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			os.Unsetenv("GITHUB_TOKEN")
		})
	})
})
