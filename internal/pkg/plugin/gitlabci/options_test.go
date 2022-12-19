package gitlabci

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("action struct", func() {
	var (
		pipelineConfig               *ci.PipelineConfig
		imageRepoURL, user, repoName string
		configLocation               downloader.ResourceLocation
		repoInfo                     *git.RepoInfo
		ciType                       server.CIServerType
	)
	BeforeEach(func() {
		imageRepoURL = "exmaple.com"
		user = "test_user"
		repoName = "test_repo"
		configLocation = "123/workflows"
		ciType = server.CIServerType("gitlab")
		pipelineConfig = &ci.PipelineConfig{
			ConfigLocation: configLocation,
			ImageRepo: &step.ImageRepoStepConfig{
				URL:  imageRepoURL,
				User: user,
			},
		}
		repoInfo = &git.RepoInfo{
			Repo:     repoName,
			RepoType: "gitlab",
		}
	})
	Context("buildCIFileConfig method", func() {
		It("should work normal", func() {
			var nilStepConfig *step.SonarQubeStepConfig
			var nilDingTalkConfig *step.DingtalkStepConfig
			var nilBool *bool
			var nilArray []string
			CIFileConfig := pipelineConfig.BuildCIFileConfig(ciType, repoInfo)
			Expect(string(CIFileConfig.Type)).Should(Equal("gitlab"))
			Expect(CIFileConfig.ConfigLocation).Should(Equal(configLocation))
			expectVars := cifile.CIFileVarsMap{
				"SonarqubeSecretKey":    "SONAR_SECRET_TOKEN",
				"AppName":               "test_repo",
				"ImageRepoSecret":       "IMAGE_REPO_SECRET",
				"ImageRepoDockerSecret": "image-repo-auth",
				"RepoType":              "gitlab",
				"imageRepo": map[string]interface{}{
					"url":  "exmaple.com",
					"user": "test_user",
				},
				"dingTalk":            nilDingTalkConfig,
				"DingTalkSecretKey":   "DINGTALK_SECURITY_VALUE",
				"DingTalkSecretToken": "DINGTALK_SECURITY_TOKEN",
				"StepGlobalVars":      "",
				"configLocation":      downloader.ResourceLocation("123/workflows"),
				"sonarqube":           nilStepConfig,
				"GitlabConnectionID":  "gitlabConnection",
				"language": map[string]interface{}{
					"name":      "",
					"version":   "",
					"frameWork": "",
				},
				"test": map[string]interface{}{
					"enable":                nilBool,
					"command":               nilArray,
					"containerName":         "",
					"coverageCommand":       "",
					"CoverageStatusCommand": "",
				},
			}
			Expect(CIFileConfig.Vars).Should(Equal(expectVars))
		})
	})
})
