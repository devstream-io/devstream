package jenkins

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
)

var _ = Describe("Pipeline struct", func() {
	var (
		pipeline                 *Pipeline
		jobName, jenkinsFilePath string
	)
	BeforeEach(func() {
		jobName = "test_pipeline"
		jenkinsFilePath = "test_path"
		pipeline = &Pipeline{
			JobName:         jobName,
			JenkinsfilePath: jenkinsFilePath,
		}
	})

	Context("getJobName method", func() {
		When("jobName has Slash", func() {
			BeforeEach(func() {
				jobName = "testFolderJob"
				pipeline.JobName = fmt.Sprintf("folder/%s", jobName)
			})
			It("should return later item", func() {
				Expect(pipeline.getJobName()).Should(Equal(jobName))
			})
		})
	})

	Context("getJobFolder method", func() {
		When("folder name exist", func() {
			BeforeEach(func() {
				jobName = "testFolderJob"
				pipeline.JobName = fmt.Sprintf("folder/%s", jobName)
			})
			It("should return later item", func() {
				Expect(pipeline.getJobFolder()).Should(Equal("folder"))
			})
		})
		When("folder name not exist", func() {
			It("should return empty string", func() {
				Expect(pipeline.getJobFolder()).Should(Equal(""))
			})
		})
	})

	Context("extractPipelinePlugins method", func() {
		When("pipeline plugins is config", func() {
			BeforeEach(func() {
				pipeline.ImageRepo = &plugins.ImageRepoJenkinsConfig{
					AuthNamespace: "test",
					URL:           "http://test.com",
					User:          "test",
				}
			})
			It("should return plugin config", func() {
				plugins := pipeline.extractPipelinePlugins()
				Expect(len(plugins)).Should(Equal(1))
			})
		})
	})

	Context("setDefaultValue method", func() {
		When("joName is empty and imageRepo namespace is empty", func() {
			BeforeEach(func() {
				pipeline.JobName = ""
				pipeline.ImageRepo = &plugins.ImageRepoJenkinsConfig{
					URL:  "http://test.com",
					User: "test",
				}
			})
			It("should set default value", func() {
				testJobName := "default_job"
				testNamespace := "test_namespace"
				pipeline.setDefaultValue(testJobName, testNamespace)
				Expect(pipeline.JobName).Should(Equal(testJobName))
				Expect(pipeline.ImageRepo).ShouldNot(BeNil())
				Expect(pipeline.ImageRepo.AuthNamespace).Should(Equal(testNamespace))
			})
		})
	})
})
