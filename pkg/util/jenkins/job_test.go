package jenkins

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("jenkins job method", func() {
	var (
		s                           *ghttp.Server
		jobName, jobFolder, reqPath string
		j                           JenkinsAPI
		err                         error
	)

	BeforeEach(func() {
		s = ghttp.NewServer()
		s.RouteToHandler("GET", "/api/json", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
		s.RouteToHandler("GET", "/crumbIssuer/api/json/api/json", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
		opts := &JenkinsConfigOption{
			URL:       s.URL(),
			Namespace: "test",
			BasicAuth: &BasicAuth{
				Username: "test_user",
				Password: "test_password",
			},
		}
		j, err = NewClient(opts)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Context("GetFolderJob method", func() {
		When("jobFolder is empty", func() {
			BeforeEach(func() {
				jobName = "test_no_folder"
				jobFolder = ""
				reqPath = fmt.Sprintf("/job/%s/api/json", jobName)
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", reqPath),
						ghttp.RespondWith(http.StatusOK, "ok"),
					),
				)

			})
			It("should get job", func() {
				_, err := j.GetFolderJob(jobName, jobFolder)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})

		When("jobFolder is setted", func() {
			BeforeEach(func() {
				jobName = "test_job"
				jobFolder = "folder"
				reqPath = fmt.Sprintf("/job/%s/job/%s/api/json", jobFolder, jobName)
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", reqPath),
						ghttp.RespondWith(http.StatusOK, "ok"),
					),
				)

			})
			It("should get job", func() {
				_, err := j.GetFolderJob(jobName, jobFolder)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})
})

var _ = Describe("BuildRenderedScript func", func() {
	var j *JobScriptRenderInfo
	It("should work normal", func() {
		j = &JobScriptRenderInfo{
			FolderName: "gg",
			RepoType:   "test",
		}
		_, err := BuildRenderedScript(j)
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("IsNotFoundError func", func() {
	var err error
	When("It is not found err", func() {
		BeforeEach(func() {
			err = errorNotFound
		})
		It("should return true", func() {
			Expect(IsNotFoundError(err)).Should(BeTrue())
		})
	})
	When("It is other err", func() {
		BeforeEach(func() {
			err = fmt.Errorf("test error")
		})
		It("should return false", func() {
			Expect(IsNotFoundError(err)).Should(BeFalse())
		})
	})
})
