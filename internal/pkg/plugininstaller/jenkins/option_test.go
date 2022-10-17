package jenkins

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

var (
	testError = errors.New("test error")
)

var _ = Describe("JobOptions struct", func() {
	var (
		jenkinsURL, secretToken, jobName, projectURL, jenkinsFilePath, userName, password, repoOwner, repoName string
		jobOptions                                                                                             *JobOptions
		basicAuth                                                                                              *jenkins.BasicAuth
		projectRepo                                                                                            *scm.Repo
		ciConfig                                                                                               *ci.CIConfig
		mockClient                                                                                             jenkins.JenkinsAPI
	)
	BeforeEach(func() {
		jenkinsURL = "http://test.com"
		userName = "test_user"
		password = "test_password"
		repoOwner = "owner"
		repoName = "repo"
		jobName = "test_folder/test_job"
		projectURL = "http://127.0.0.1:300/test/project"
		jenkinsFilePath = "http://raw.content.com/Jenkinsfile"
		basicAuth = &jenkins.BasicAuth{
			Password: password,
			Username: userName,
		}
		projectRepo = &scm.Repo{
			Owner:    repoOwner,
			Repo:     repoName,
			Branch:   "test",
			BaseURL:  "http://127.0.0.1:300",
			RepoType: "gitlab",
		}
		ciConfig = &ci.CIConfig{
			Type:           "jenkins",
			ConfigLocation: jenkinsFilePath,
		}
		secretToken = "secret"
		jobOptions = &JobOptions{
			Jenkins: Jenkins{
				URL:           jenkinsURL,
				User:          userName,
				Namespace:     "jenkins",
				EnableRestart: false,
			},
			SCM: scm.SCMInfo{
				CloneURL: projectURL,
				Branch:   "test",
			},
			Pipeline: Pipeline{
				JobName:         jobName,
				JenkinsfilePath: jenkinsFilePath,
				ImageRepo:       &plugins.ImageRepoJenkinsConfig{},
			},
			BasicAuth:   basicAuth,
			ProjectRepo: projectRepo,
			CIConfig:    ciConfig,
			SecretToken: secretToken,
		}
	})
	Context("createOrUpdateJob method", func() {
		When("jenkins client return normal", func() {
			BeforeEach(func() {
				jobOptions.SCM.SSHprivateKey = "test"
				mockClient = &jenkins.MockClient{}
			})
			It("should work noraml", func() {
				err := jobOptions.createOrUpdateJob(mockClient)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("jenkins client return error", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{
					ExecuteScriptError: testError,
				}
			})
			It("should return error", func() {
				err := jobOptions.createOrUpdateJob(mockClient)
				Expect(err).Should(HaveOccurred())
			})
		})
	})
	Context("buildWebhookInfo method", func() {
		When("repo is gitlab", func() {
			It("should work for gitlab", func() {
				webHookInfo := jobOptions.buildWebhookInfo()
				Expect(webHookInfo.Address).Should(Equal(fmt.Sprintf("%s/project/%s", jobOptions.Jenkins.URL, jobOptions.Pipeline.JobName)))
				Expect(webHookInfo.SecretToken).Should(Equal(secretToken))
			})
		})
		When("repo is github", func() {
			BeforeEach(func() {
				jobOptions.ProjectRepo.RepoType = "github"
			})
			It("should work for github", func() {
				webHookInfo := jobOptions.buildWebhookInfo()
				Expect(webHookInfo.Address).Should(Equal(fmt.Sprintf("%s/github-webhook/", jobOptions.Jenkins.URL)))
				Expect(webHookInfo.SecretToken).Should(Equal(secretToken))
			})
		})
	})
	Context("installPlugins method", func() {
		When("jenkins client return error", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{
					InstallPluginsIfNotExistsError: testError,
				}
			})
			It("should return error", func() {
				var pluginConfigs []pluginConfigAPI
				err := installPlugins(mockClient, pluginConfigs, false)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})
	Context("deleteJob method", func() {
		When("jenkins job get error", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{
					GetFolderJobError: fmt.Errorf("job error"),
				}
			})
			It("should return error", func() {
				err := jobOptions.deleteJob(mockClient)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("jenkins job is not exist", func() {
			var errMsg string
			BeforeEach(func() {
				errMsg = "404"
				mockClient = &jenkins.MockClient{
					GetFolderJobError: fmt.Errorf(errMsg),
				}
			})
			It("should return delete error", func() {
				err := jobOptions.deleteJob(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("buildCIConfig method", func() {
		When("jenkinsfilePath is local path", func() {
			BeforeEach(func() {
				jobOptions.Pipeline.JenkinsfilePath = "test/local"
			})
			It("should use localPath", func() {
				ciConfig, err := jobOptions.buildCIConfig()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(ciConfig.ConfigLocation).Should(Equal(jobOptions.Pipeline.JenkinsfilePath))
			})
		})
	})

	Context("extractJenkinsPlugins method", func() {
		When("repo is gitlab", func() {
			BeforeEach(func() {
				jobOptions.ProjectRepo.RepoType = "gitlab"
				jobOptions.Pipeline = Pipeline{
					JobName:         jobName,
					JenkinsfilePath: jenkinsFilePath,
				}
			})
			It("should return pluginConfig", func() {
				configs := jobOptions.extractJenkinsPlugins()
				Expect(len(configs)).Should(Equal(1))
			})
		})
		When("repo is github", func() {
			BeforeEach(func() {
				jobOptions.ProjectRepo.RepoType = "github"
				jobOptions.Pipeline = Pipeline{
					JobName:         jobName,
					JenkinsfilePath: jenkinsFilePath,
				}
			})
			It("should return pluginConfig", func() {
				configs := jobOptions.extractJenkinsPlugins()
				Expect(len(configs)).Should(Equal(1))
			})
		})
	})
})
