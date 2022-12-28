package step

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("GithubStepConfig", func() {
	var (
		c          *GithubStepConfig
		mockClient *jenkins.MockClient
		repoOwner  string
	)
	BeforeEach(func() {
		repoOwner = "test_user"
		c = &GithubStepConfig{
			RepoOwner: repoOwner,
		}
	})

	Context("GetJenkinsPlugins method", func() {
		It("should return github plugin", func() {
			plugins := c.GetJenkinsPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("github-branch-source"))
		})
	})

	Context("ConfigJenkins method", func() {
		When("create password failed", func() {
			var (
				createErrMsg string
			)
			BeforeEach(func() {
				createErrMsg = "create password failed"
				mockClient = &jenkins.MockClient{
					CreatePasswordCredentialError: fmt.Errorf(createErrMsg),
				}
			})
			It("should return error", func() {
				_, err := c.ConfigJenkins(mockClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(createErrMsg))
			})
		})
		When("create password success", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return repoCascConfig", func() {
				cascConfig, err := c.ConfigJenkins(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(cascConfig.RepoType).Should(Equal("github"))
			})
		})
	})
	Context("ConfigSCM method", func() {
		It("should return nil", func() {
			scmClient := &scm.MockScmClient{}
			err := c.ConfigSCM(scmClient)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("newGithubStep func", func() {
	var (
		pluginConfig *StepGlobalOption
	)
	BeforeEach(func() {
		pluginConfig = &StepGlobalOption{
			RepoInfo: &git.RepoInfo{
				Owner: "test",
			},
		}
	})
	It("should return github plugin", func() {
		githubPlugin := newGithubStep(pluginConfig)
		Expect(githubPlugin.RepoOwner).Should(Equal("test"))
	})
})
