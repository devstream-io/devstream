package ci_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("PipelineConfig struct", func() {
	var (
		a                                            *ci.PipelineConfig
		imageRepoURL, user, repoName, configLocation string
		r                                            *git.RepoInfo
		ciType                                       server.CIServerType
	)
	BeforeEach(func() {
		imageRepoURL = "exmaple.com"
		user = "test_user"
		repoName = "test_repo"
		configLocation = "123/workflows"
		ciType = "gitlab"
		a = &ci.PipelineConfig{
			ConfigLocation: configLocation,
			ImageRepo: &step.ImageRepoStepConfig{
				URL:  imageRepoURL,
				User: user,
			},
		}
		r = &git.RepoInfo{
			Repo: repoName,
		}
	})
	Context("BuildCIFileConfig method", func() {
		It("should work normal", func() {
			var nilStepConfig *step.SonarQubeStepConfig
			var nilDingTalkConfig *step.DingtalkStepConfig
			var nilGeneral *step.GeneralStepConfig
			CIFileConfig := a.BuildCIFileConfig(ciType, r)
			Expect(string(CIFileConfig.Type)).Should(Equal("gitlab"))
			Expect(CIFileConfig.ConfigLocation).Should(Equal(configLocation))
			expectVars := cifile.CIFileVarsMap{
				"SonarqubeSecretKey":    "SONAR_SECRET_TOKEN",
				"AppName":               "test_repo",
				"ImageRepoSecret":       "IMAGE_REPO_SECRET",
				"ImageRepoDockerSecret": "image-repo-auth",
				"imageRepo": map[string]interface{}{
					"url":  "exmaple.com",
					"user": "test_user",
				},
				"general":             nilGeneral,
				"dingTalk":            nilDingTalkConfig,
				"DingTalkSecretKey":   "DINGTALK_SECURITY_VALUE",
				"DingTalkSecretToken": "DINGTALK_SECURITY_TOKEN",
				"StepGlobalVars":      "",
				"configLocation":      "123/workflows",
				"sonarqube":           nilStepConfig,
				"GitlabConnectionID":  "gitlabConnection",
			}
			Expect(CIFileConfig.Vars).Should(Equal(expectVars))
		})
	})
})
