package ci_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var _ = Describe("SetDefault func", func() {
	var (
		options             configmanager.RawOptions
		url, configLocation string
	)
	BeforeEach(func() {
		url = "https://github.com/root/test-exmaple.git"
		configLocation = "workflows"
		options = configmanager.RawOptions{
			"scm": map[string]interface{}{
				"url": url,
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
		opts, err := ci.SetDefault(server.CIGithubType)(options)
		Expect(err).ShouldNot(HaveOccurred())
		log.Infof("----> %+v", opts)
		projectRepo, ok := opts["scm"]
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
