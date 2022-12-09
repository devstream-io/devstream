package jenkinspipeline

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/github"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

var _ = Describe("setDefault func", func() {
	var (
		jenkinsUser, jenkinsPassword, jenkinsURL, jenkinsFilePath, projectURL string
		options                                                               map[string]interface{}
	)
	BeforeEach(func() {
		jenkinsUser = "test"
		jenkinsPassword = "testPassword"
		jenkinsURL = "http://test.jenkins.com/"
		projectURL = "https://test.gitlab.com/test/test_project"
		jenkinsFilePath = "http://raw.content.com/Jenkinsfile"
		err := os.Setenv("JENKINS_PASSWORD", jenkinsPassword)
		Expect(err).NotTo(HaveOccurred())
		options = map[string]interface{}{
			"jenkins": map[string]interface{}{
				"url":  jenkinsURL,
				"user": jenkinsUser,
			},
			"scm": map[string]interface{}{
				"url": projectURL,
			},
			"pipeline": map[string]interface{}{
				"configLocation": jenkinsFilePath,
			},
		}
	})
	When("all input is valid", func() {
		BeforeEach(func() {
			options["scm"] = map[string]interface{}{
				"url":    "git@44.33.22.11:30022:root/spring-demo.git",
				"apiURL": "http://www.app.com",
			}
			options["projectRepo"] = map[string]interface{}{
				"repo": "testRepo",
			}
		})
		It("should set default value", func() {
			newOptions, err := setJenkinsDefault(options)
			Expect(err).Error().ShouldNot(HaveOccurred())
			opts, err := newJobOptions(newOptions)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(opts.JobName)).Should(Equal("spring-demo"))
		})
	})
	AfterEach(func() {
		os.Unsetenv("JENKINS_PASSWORD")
	})
})

var _ = Describe("validate func", func() {
	var (
		githubToken, gitlabToken string
	)
	BeforeEach(func() {
		githubToken = os.Getenv(github.TokenEnvKey)
		gitlabToken = os.Getenv(gitlab.TokenEnvKey)
		err := os.Unsetenv(github.TokenEnvKey)
		Expect(err).Error().ShouldNot(HaveOccurred())
		err = os.Unsetenv(gitlab.TokenEnvKey)
		Expect(err).Error().ShouldNot(HaveOccurred())
	})
	AfterEach(func() {
		if githubToken != "" {
			os.Setenv(github.TokenEnvKey, githubToken)
		}
		if gitlabToken != "" {
			os.Setenv(gitlab.TokenEnvKey, gitlabToken)
		}
	})
	var (
		jenkinsUser, jenkinsURL, jenkinsFilePath, projectURL, repoType string
		options, projectRepo, pipeline                                 map[string]interface{}
	)
	BeforeEach(func() {
		jenkinsUser = "test"
		jenkinsURL = "http://test.jenkins.com/"
		projectURL = "https://test.gitlab.com/test/test_project"
		jenkinsFilePath = "http://raw.content.com/Jenkinsfile"
		pipeline = map[string]interface{}{
			"configLocation": jenkinsFilePath,
		}
		projectRepo = map[string]interface{}{
			"owner": "test_owner",
			"org":   "test_org",
			"repo":  "test_repo",
		}
		options = map[string]interface{}{
			"jobName": "test",
			"jenkins": map[string]interface{}{
				"url":  jenkinsURL,
				"user": jenkinsUser,
			},
			"scm": map[string]interface{}{
				"url": projectURL,
			},
			"pipeline":    pipeline,
			"projectRepo": projectRepo,
		}
	})
	When("Input field miss", func() {
		BeforeEach(func() {
			options = map[string]interface{}{
				"jenkins": map[string]interface{}{
					"url":  jenkinsURL,
					"user": jenkinsUser,
				},
				"scm": map[string]interface{}{
					"url": projectURL,
				},
			}
		})
		It("should return error", func() {
			_, err := validateJenkins(options)
			Expect(err).Error().Should(HaveOccurred())
		})
	})

	When("jobName is not valid", func() {
		BeforeEach(func() {
			options["jobName"] = "folder/not_exist/jobName"
			options["pipeline"] = pipeline
			repoType = "github"
			projectRepo = map[string]any{
				"scmType": repoType,
				"name":    "test",
				"owner":   "test_user",
				"branch":  "main",
			}
			options["scm"] = projectRepo
			os.Setenv(github.TokenEnvKey, "test_env")
		})
		It("should return error", func() {
			_, err := validateJenkins(options)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(Equal(fmt.Sprintf("jenkins jobName illegal: %s", options["jobName"])))
		})
	})
	When("all params is right", func() {
		BeforeEach(func() {
			options["pipeline"] = pipeline
			repoType = "github"
			projectRepo = map[string]any{
				"scmType": repoType,
				"name":    "test",
				"owner":   "test_user",
				"branch":  "main",
			}
			options["scm"] = projectRepo
			os.Setenv(github.TokenEnvKey, "test_env")
		})
		It("should return nil error", func() {
			_, err := validateJenkins(options)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
})
