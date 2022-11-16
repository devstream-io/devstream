package jenkinspipeline

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("newJobOptions func", func() {
	var (
		jenkinsURL, jobName, projectURL, jenkinsFilePath, userName string
		rawOptions                                                 configmanager.RawOptions
	)
	BeforeEach(func() {
		jenkinsURL = "http://test.com"
		userName = "test_user"
		jobName = "test_folder/test_job"
		projectURL = "http://127.0.0.1:300/test/project"
		jenkinsFilePath = "http://raw.content.com/Jenkinsfile"
		rawOptions = configmanager.RawOptions{
			"jobName": jobName,
			"jenkins": map[string]interface{}{
				"url":  jenkinsURL,
				"user": userName,
			},
			"scm": map[string]interface{}{
				"cloneURL": projectURL,
			},
			"pipeline": map[string]interface{}{
				"configLocation": jenkinsFilePath,
			},
		}

	})
	It("should work normal", func() {
		job, err := newJobOptions(rawOptions)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(string(job.JobName)).Should(Equal(jobName))
		Expect(job.Pipeline.ConfigLocation).Should(Equal(jenkinsFilePath))
		Expect(job.SCM.CloneURL).Should(Equal(projectURL))
		Expect(job.Jenkins.URL).Should(Equal(jenkinsURL))
		Expect(job.Jenkins.User).Should(Equal(userName))
	})
})

var _ = Describe("options struct", func() {
	var (
		jobName, jenkinsFilePath, repoOwner, repoName, secretToken, errMsg string
		repoInfo                                                           *git.RepoInfo
		j                                                                  *jenkins.MockClient
		opts                                                               *jobOptions
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
		opts = &jobOptions{
			JobName: jenkinsJobName(jobName),
			CIConfig: ci.CIConfig{
				ProjectRepo: repoInfo,
				Pipeline: ci.PipelineConfig{
					ConfigLocation: jenkinsFilePath,
				},
			},
		}
		j = &jenkins.MockClient{}
	})

	Context("getJobName method", func() {
		When("jobName has slash", func() {
			BeforeEach(func() {
				jobName = "testFolderJob"
				opts.JobName = jenkinsJobName(fmt.Sprintf("folder/%s", jobName))
			})
			It("should return later item", func() {
				Expect(opts.JobName.getJobName()).Should(Equal(jobName))
			})
		})
		When("jobName does'nt have slash", func() {
			BeforeEach(func() {
				opts.JobName = jenkinsJobName("testJob")
			})
			It("should return jobName", func() {
				Expect(opts.JobName.getJobName()).Should(Equal("testJob"))
			})
		})
	})

	Context("getJobFolder method", func() {
		When("folder name exist", func() {
			BeforeEach(func() {
				opts.JobName = jenkinsJobName(fmt.Sprintf("folder/%s", jobName))
			})
			It("should return later item", func() {
				Expect(opts.JobName.getJobFolder()).Should(Equal("folder"))
			})
		})
		When("folder name not exist", func() {
			It("should return empty string", func() {
				Expect(opts.JobName.getJobFolder()).Should(Equal(""))
			})
		})
	})

	Context("extractpipelinePlugins method", func() {
		When("repo type is github", func() {
			BeforeEach(func() {
				opts.CIConfig.Pipeline.ImageRepo = &step.ImageRepoStepConfig{
					URL:  "http://test.com",
					User: "test",
				}
				opts.CIConfig.ProjectRepo.RepoType = "github"
			})
			It("should return plugin config", func() {
				plugins := opts.extractPlugins()
				Expect(len(plugins)).Should(Equal(2))
			})
		})
		When("repo type is gitlab", func() {
			BeforeEach(func() {
				opts.CIConfig.Pipeline.ImageRepo = &step.ImageRepoStepConfig{
					URL:  "http://test.com",
					User: "test",
				}
				opts.CIConfig.ProjectRepo.RepoType = "gitlab"
			})
			It("should return plugin config", func() {
				plugins := opts.extractPlugins()
				Expect(len(plugins)).Should(Equal(2))
			})
		})
	})

	Context("createOrUpdateJob method", func() {
		When("jenkins client return normal", func() {
			BeforeEach(func() {
				opts.CIConfig.ProjectRepo.SSHPrivateKey = "test"
			})
			It("should work noraml", func() {
				err := opts.createOrUpdateJob(j, secretToken)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("jenkins client script error", func() {
			BeforeEach(func() {
				errMsg = "script err"
				j.ExecuteScriptError = fmt.Errorf(errMsg)
			})
			It("should return error", func() {
				err := opts.createOrUpdateJob(j, secretToken)
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
				err := opts.remove(j)
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
				err := opts.remove(j)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("install method", func() {
		When("install plugin failed", func() {
			BeforeEach(func() {
				errMsg = "install plugin failed"
				j.InstallPluginsIfNotExistsError = fmt.Errorf(errMsg)
			})
			It("should return error", func() {
				err := opts.install(j, secretToken)
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
				err := opts.install(j, secretToken)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("all config valid", func() {
			BeforeEach(func() {
				j = &jenkins.MockClient{}
			})
			It("should work normal", func() {
				err := opts.install(j, secretToken)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
