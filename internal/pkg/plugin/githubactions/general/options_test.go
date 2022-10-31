package general

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("action struct", func() {
	var (
		a                                            *action
		imageRepoURL, user, repoName, configLocation string
		r                                            *git.RepoInfo
	)
	BeforeEach(func() {
		imageRepoURL = "exmaple.com"
		user = "test_user"
		repoName = "test_repo"
		configLocation = "123/workflows"
		a = &action{
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
	Context("buildCIConfig method", func() {
		It("should work normal", func() {
			var nilStepConfig *step.SonarQubeStepConfig
			var nilDingTalkConfig *step.DingtalkStepConfig
			ciConfig := a.buildCIConfig(r)
			Expect(string(ciConfig.Type)).Should(Equal("github"))
			Expect(ciConfig.ConfigLocation).Should(Equal(configLocation))
			expectVars := cifile.CIFileVarsMap{
				"SonarqubeSecretKey": "SONAR_SECRET_TOKEN",
				"AppName":            "test_repo",
				"imageRepo": map[string]interface{}{
					"url":  "exmaple.com",
					"user": "test_user",
				},
				"dingTalk":            nilDingTalkConfig,
				"ImageRepoSecret":     "IMAGE_REPO_SECRET",
				"DingTalkSecretKey":   "DINGTALK_SECURITY_VALUE",
				"DingTalkSecretToken": "",
				"StepGlobalVars":      "",
				"configLocation":      "123/workflows",
				"sonarqube":           nilStepConfig,
			}
			Expect(ciConfig.Vars).Should(Equal(expectVars))
		})
	})
})
