package jenkins

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("pipeline struct", func() {
	var (
		p                                                                  *pipeline
		jobName, jenkinsFilePath, repoOwner, repoName, secretToken, errMsg string
		repoInfo                                                           *git.RepoInfo
		j                                                                  *jenkins.MockClient
	)
	BeforeEach(func() {
		repoOwner = "owner"
		repoName = "repo"
		jobName = "test_pipeline"
		jenkinsFilePath = "test_path"
		repoInfo = &git.RepoInfo{
			Owner:    repoOwner,
			Repo:     repoName,
			Branch:   "test",
			BaseURL:  "http://127.0.0.1:300",
			RepoType: "gitlab",
		}
		secretToken = "test_secret"
		p = &pipeline{
			Job:             jobName,
			JenkinsfilePath: jenkinsFilePath,
		}
		j = &jenkins.MockClient{}
	})

	Context("getJobName method", func() {
		When("jobName has slash", func() {
			BeforeEach(func() {
				jobName = "testFolderJob"
				p.Job = fmt.Sprintf("folder/%s", jobName)
			})
			It("should return later item", func() {
				Expect(p.getJobName()).Should(Equal(jobName))
			})
		})
		When("jobName does'nt have slash", func() {
			BeforeEach(func() {
				p.Job = "testJob"
			})
			It("should return jobName", func() {
				Expect(p.getJobName()).Should(Equal("testJob"))
			})
		})
	})

	Context("getJobFolder method", func() {
		When("folder name exist", func() {
			BeforeEach(func() {
				jobName = "testFolderJob"
				p.Job = fmt.Sprintf("folder/%s", jobName)
			})
			It("should return later item", func() {
				Expect(p.getJobFolder()).Should(Equal("folder"))
			})
		})
		When("folder name not exist", func() {
			It("should return empty string", func() {
				Expect(p.getJobFolder()).Should(Equal(""))
			})
		})
	})

	Context("extractpipelinePlugins method", func() {
		When("repo type is github", func() {
			BeforeEach(func() {
				p.ImageRepo = &step.ImageRepoStepConfig{
					URL:  "http://test.com",
					User: "test",
				}
			})
			It("should return plugin config", func() {
				plugins := p.extractPlugins(&git.RepoInfo{
					RepoType: "github",
				})
				Expect(len(plugins)).Should(Equal(2))
			})
		})
		When("repo type is gitlab", func() {
			BeforeEach(func() {
				p.ImageRepo = &step.ImageRepoStepConfig{
					URL:  "http://test.com",
					User: "test",
				}
			})
			It("should return plugin config", func() {
				plugins := p.extractPlugins(&git.RepoInfo{
					RepoType: "gitlab",
				})
				Expect(len(plugins)).Should(Equal(2))
			})
		})
	})

	Context("createOrUpdateJob method", func() {
		When("jenkins client return normal", func() {
			BeforeEach(func() {
				repoInfo.SSHPrivateKey = "test"
			})
			It("should work noraml", func() {
				err := p.createOrUpdateJob(j, repoInfo, secretToken)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("jenkins client script error", func() {
			BeforeEach(func() {
				errMsg = "script err"
				j.ExecuteScriptError = fmt.Errorf(errMsg)
			})
			It("should return error", func() {
				err := p.createOrUpdateJob(j, repoInfo, secretToken)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
	})

	Context("remove method", func() {
		When("jenkins job get error", func() {
			BeforeEach(func() {
				errMsg = "get job error"
				j.GetFolderJobError = fmt.Errorf(errMsg)
			})
			It("should return error", func() {
				err := p.remove(j, repoInfo)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("jenkins job is not exist", func() {
			BeforeEach(func() {
				errMsg = "404"
				j.GetFolderJobError = fmt.Errorf(errMsg)
			})
			It("should return delete error", func() {
				err := p.remove(j, repoInfo)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("buildCIConfig method", func() {
		var (
			customMap map[string]interface{}
		)
		BeforeEach(func() {
			p.JenkinsfilePath = "test/local"
			p.Job = "test"
			p.ImageRepo = &step.ImageRepoStepConfig{
				URL:  "testurl",
				User: "testuser",
			}
			customMap = map[string]interface{}{
				"test": "gg",
			}
		})
		It("should work normal", func() {
			var emptyDingTalk *step.DingtalkStepConfig
			var emptySonar *step.SonarQubeStepConfig
			var emptyGeneral *step.GeneralStepConfig
			ciConfig := p.buildCIConfig(repoInfo, customMap)
			Expect(ciConfig.ConfigLocation).Should(Equal(p.JenkinsfilePath))
			Expect(string(ciConfig.Type)).Should(Equal("jenkins"))
			expectedMap := cifile.CIFileVarsMap{
				"AppName":         "test",
				"jobName":         "test",
				"jenkinsfilePath": "test/local",
				"imageRepo": map[string]interface{}{
					"url":  "testurl",
					"user": "testuser",
				},
				"dingTalk":            emptyDingTalk,
				"sonarqube":           emptySonar,
				"general":             emptyGeneral,
				"ImageRepoSecret":     "IMAGE_REPO_SECRET",
				"DingTalkSecretKey":   "DINGTALK_SECURITY_VALUE",
				"DingTalkSecretToken": "",
				"StepGlobalVars":      "",
				"SonarqubeSecretKey":  "SONAR_SECRET_TOKEN",
			}
			Expect(ciConfig.Vars).Should(Equal(expectedMap))
		})
	})

	Context("install method", func() {
		When("install plugin failed", func() {
			BeforeEach(func() {
				errMsg = "install plugin failed"
				j.InstallPluginsIfNotExistsError = fmt.Errorf(errMsg)
			})
			It("should return error", func() {
				err := p.install(j, repoInfo, secretToken)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("config plugin failed", func() {
			BeforeEach(func() {
				errMsg = "config plugin failed"
				j = &jenkins.MockClient{}
				j.ConfigCascForRepoError = fmt.Errorf(errMsg)
			})
			It("should return error", func() {
				err := p.install(j, repoInfo, secretToken)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("all config valid", func() {
			BeforeEach(func() {
				j = &jenkins.MockClient{}
			})
			It("should work normal", func() {
				err := p.install(j, repoInfo, secretToken)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
