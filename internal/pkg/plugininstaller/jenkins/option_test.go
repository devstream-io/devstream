package jenkins

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
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
			"jenkins": map[string]interface{}{
				"url":  jenkinsURL,
				"user": userName,
			},
			"scm": map[string]interface{}{
				"cloneURL": projectURL,
			},
			"pipeline": map[string]interface{}{
				"jobName":         jobName,
				"jenkinsfilePath": jenkinsFilePath,
			},
		}

	})
	It("should work normal", func() {
		job, err := newJobOptions(rawOptions)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(job.Pipeline.Job).Should(Equal(jobName))
		Expect(job.Pipeline.JenkinsfilePath).Should(Equal(jenkinsFilePath))
		Expect(job.SCM.CloneURL).Should(Equal(projectURL))
		Expect(job.Jenkins.URL).Should(Equal(jenkinsURL))
		Expect(job.Jenkins.User).Should(Equal(userName))
	})
})
