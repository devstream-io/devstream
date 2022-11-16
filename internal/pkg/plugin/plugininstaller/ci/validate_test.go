package ci_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci"
)

var _ = Describe("SetDefault func", func() {
	var (
		options                  configmanager.RawOptions
		cloneURL, configLocation string
	)
	BeforeEach(func() {
		cloneURL = "http://github.com/root/test-exmaple.git"
		configLocation = "workflows"
		options = configmanager.RawOptions{
			"scm": map[string]interface{}{
				"cloneURL": cloneURL,
			},
			"pipeline": map[string]interface{}{
				"configLocation": configLocation,
			},
		}
	})
	BeforeEach(func() {
		os.Setenv("GITHUB_TOKEN", "test")
	})
	It("should work normal", func() {
		opts, err := ci.SetSCMDefault(options)
		Expect(err).ShouldNot(HaveOccurred())
		projectRepo, ok := opts["projectRepo"]
		Expect(ok).Should(BeTrue())
		Expect(projectRepo).ShouldNot(BeNil())
		CIFileConfig, ok := opts["ci"]
		Expect(ok).Should(BeTrue())
		Expect(CIFileConfig).ShouldNot(BeNil())
	})
	AfterEach(func() {
		os.Unsetenv("GITHUB_TOKEN")
	})
})
