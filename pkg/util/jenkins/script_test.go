package jenkins

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("scriptError struct", func() {
	var (
		e *scriptError
	)
	BeforeEach(func() {
		e = &scriptError{}
	})
	Context("Error method", func() {
		When("errorMsg has colon", func() {
			BeforeEach(func() {
				e.errorMsg = "test: test with colon"
			})
			It("should return err msg", func() {
				result := e.Error()
				Expect(result).Should(Equal("execute groovy script failed: test with colon()"))
			})
		})
		When("errorMsg not has colon", func() {
			BeforeEach(func() {
				e.errorMsg = "this is test"
			})
			It("should return script execute err msg", func() {
				result := e.Error()
				Expect(result).Should(Equal(fmt.Sprintf("execute groovy script failed: %s()", e.errorMsg)))
			})
		})
	})
})

var _ = Describe("jenkins script method", func() {
	var (
		s   *ghttp.Server
		j   JenkinsAPI
		err error
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
	Context("ConfigCascForRepo method", func() {
		var (
			cascConfig *RepoCascConfig
		)
		When("response text not contain verifier", func() {
			BeforeEach(func() {
				cascConfig = &RepoCascConfig{
					RepoType:     "github",
					CredentialID: "credId",
					JenkinsURL:   "testURL",
					SecretToken:  "secretToken",
				}
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/scriptText"),
						ghttp.RespondWith(http.StatusOK, "ok"),
					),
				)
			})
			It("should return err", func() {
				err = j.ConfigCascForRepo(cascConfig)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("execute groovy script failed: script verifier error(ok)"))
			})
		})
	})
})
