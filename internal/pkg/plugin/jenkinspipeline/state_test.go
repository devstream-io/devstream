package jenkinspipeline

import (
	"fmt"

	"github.com/bndr/gojenkins"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

var _ = Describe("getJobState func", func() {
	var (
		c jenkins.JenkinsAPI
	)
	When("GetFolderJob return error", func() {
		BeforeEach(func() {
			c = &jenkins.MockClient{
				GetFolderJobError: fmt.Errorf("test error"),
			}
		})
		It("should return err", func() {
			_, err := getJobState(c, "test_job", "test_folder")
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("GetFolderJob return job", func() {
		var (
			jobClass, jobURL string
		)
		BeforeEach(func() {
			jobClass = "testClass"
			jobURL = "testURL"
			c = &jenkins.MockClient{
				GetFolderJobValue: &gojenkins.Job{
					Raw: &gojenkins.JobResponse{
						Class: jobClass,
						URL:   jobURL,
					},
				},
			}
		})
		It("should return job res", func() {
			job, err := getJobState(c, "test_job", "test_folder")
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(job)).ShouldNot(BeZero())
			class, ok := job["Class"]
			Expect(ok).Should(BeTrue())
			Expect(class).Should(Equal(jobClass))
			url, ok := job["URL"]
			Expect(ok).Should(BeTrue())
			Expect(url).Should(Equal(jobURL))
		})
	})
})
