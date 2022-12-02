package ci

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("PipelineConfig struct", func() {
	var (
		a                            *PipelineConfig
		imageRepoURL, user, repoName string
		configLocation               downloader.ResourceLocation
		r                            *git.RepoInfo
		ciType                       server.CIServerType
	)
	BeforeEach(func() {
		imageRepoURL = "exmaple.com"
		user = "test_user"
		repoName = "test_repo"
		configLocation = "123/workflows"
		ciType = "gitlab"
		a = &PipelineConfig{
			ConfigLocation: configLocation,
			ImageRepo: &step.ImageRepoStepConfig{
				URL:  imageRepoURL,
				User: user,
			},
		}
		r = &git.RepoInfo{
			Repo:     repoName,
			RepoType: "gitlab",
		}
	})
	Context("BuildCIFileConfig method", func() {
		It("should work normal", func() {
			var nilStepConfig *step.SonarQubeStepConfig
			var nilDingTalkConfig *step.DingtalkStepConfig
			var emptyBool *bool
			var emptyArray []string
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
				"dingTalk":            nilDingTalkConfig,
				"DingTalkSecretKey":   "DINGTALK_SECURITY_VALUE",
				"DingTalkSecretToken": "DINGTALK_SECURITY_TOKEN",
				"StepGlobalVars":      "",
				"configLocation":      downloader.ResourceLocation("123/workflows"),
				"sonarqube":           nilStepConfig,
				"RepoType":            "gitlab",
				"GitlabConnectionID":  "gitlabConnection",
				"test": map[string]interface{}{
					"enable":                emptyBool,
					"command":               emptyArray,
					"containerName":         "",
					"coverageCommand":       "",
					"CoverageStatusCommand": "",
				},
				"language": map[string]interface{}{
					"name":      "",
					"version":   "",
					"frameWork": "",
				},
			}
			Expect(CIFileConfig.Vars).Should(Equal(expectVars))
		})
	})
	It("should return file Vars", func() {
		varMap := a.GenerateCIFileVars(r)
		var emptyDingtalk *step.DingtalkStepConfig
		var emptySonar *step.SonarQubeStepConfig
		var emptyBool *bool
		var emptyArray []string
		Expect(varMap).Should(Equal(cifile.CIFileVarsMap{
			"configLocation":        downloader.ResourceLocation("123/workflows"),
			"DingTalkSecretToken":   "DINGTALK_SECURITY_TOKEN",
			"ImageRepoSecret":       "IMAGE_REPO_SECRET",
			"ImageRepoDockerSecret": "image-repo-auth",
			"AppName":               "test_repo",
			"StepGlobalVars":        "",
			"RepoType":              "gitlab",
			"imageRepo": map[string]interface{}{
				"url":  "exmaple.com",
				"user": "test_user",
			},
			"dingTalk":           emptyDingtalk,
			"sonarqube":          emptySonar,
			"SonarqubeSecretKey": "SONAR_SECRET_TOKEN",
			"GitlabConnectionID": "gitlabConnection",
			"DingTalkSecretKey":  "DINGTALK_SECURITY_VALUE",
			"test": map[string]interface{}{
				"enable":                emptyBool,
				"command":               emptyArray,
				"containerName":         "",
				"coverageCommand":       "",
				"CoverageStatusCommand": "",
			},
			"language": map[string]interface{}{
				"name":      "",
				"version":   "",
				"frameWork": "",
			},
		}))
	})
})
