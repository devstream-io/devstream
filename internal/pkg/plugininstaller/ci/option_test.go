package ci

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/server"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("Options struct", func() {
	var (
		rawOptions        configmanager.RawOptions
		defaultOpts, opts *Options
	)
	BeforeEach(func() {
		opts = &Options{}
		defaultCIConfig := &CIConfig{
			Type:           "github",
			ConfigLocation: "http://www.test.com",
		}
		defaultRepo := &git.RepoInfo{
			Owner:    "test",
			Repo:     "test_repo",
			Branch:   "test_branch",
			RepoType: "gitlab",
		}
		defaultOpts = &Options{
			CIConfig:    defaultCIConfig,
			ProjectRepo: defaultRepo,
		}
	})

	Context("NewOptions method", func() {
		When("options is valid", func() {
			BeforeEach(func() {
				rawOptions = configmanager.RawOptions{
					"ci": map[string]interface{}{
						"type":    "gitlab",
						"content": "test",
					},
				}
			})
			It("should not raise error", func() {
				_, err := NewOptions(rawOptions)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("fillDefaultValue method", func() {
		When("ci config and repo are all empty", func() {
			It("should set default value", func() {
				opts.fillDefaultValue(defaultOpts)
				Expect(opts.CIConfig).ShouldNot(BeNil())
				Expect(opts.ProjectRepo).ShouldNot(BeNil())
				Expect(opts.CIConfig.ConfigLocation).Should(Equal("http://www.test.com"))
				Expect(opts.ProjectRepo.Repo).Should(Equal("test_repo"))
			})
		})
		When("ci config and repo has some value", func() {
			BeforeEach(func() {
				opts.CIConfig = &CIConfig{
					ConfigLocation: "http://exist.com",
				}
				opts.ProjectRepo = &git.RepoInfo{
					Branch: "new_branch",
				}
			})
			It("should update empty value", func() {
				opts.fillDefaultValue(defaultOpts)
				Expect(opts.CIConfig).ShouldNot(BeNil())
				Expect(opts.ProjectRepo).ShouldNot(BeNil())
				Expect(opts.CIConfig.ConfigLocation).Should(Equal("http://exist.com"))
				Expect(opts.CIConfig.Type).Should(Equal(server.CIServerType("github")))
				Expect(opts.ProjectRepo.Branch).Should(Equal("new_branch"))
				Expect(opts.ProjectRepo.Repo).Should(Equal("test_repo"))
			})
		})
	})
})
