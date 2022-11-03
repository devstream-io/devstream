package general

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("validate func", func() {
	var (
		options                  configmanager.RawOptions
		cloneURL, configLocation string
	)
	BeforeEach(func() {
		cloneURL = "git@github.com/root/test-exmaple.git"
		configLocation = "workflows"
		options = configmanager.RawOptions{
			"scm": map[string]interface{}{
				"cloneURL": cloneURL,
			},
			"pipeline": map[string]interface{}{
				"configLocation": configLocation,
			},
			"projectRepo": map[string]interface{}{
				"repoType": "github",
				"owner":    "test",
				"repo":     "gg",
			},
		}
	})
	When("scm token is not setted", func() {
		It("should return err", func() {
			_, err := validate(options)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("scm repo is gitlab", func() {
		BeforeEach(func() {
			options["scm"] = map[string]interface{}{
				"cloneURL": "http://exmaple.gitlab.com",
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
