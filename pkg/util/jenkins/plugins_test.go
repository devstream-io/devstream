package jenkins

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("jenins plugin methods", func() {
	var (
		p   []*JenkinsPlugin
		s   *ghttp.Server
		j   JenkinsAPI
		err error
	)
	BeforeEach(func() {
		p = []*JenkinsPlugin{
			{
				Name:    "test",
				Version: "test_version",
			},
		}
		s = ghttp.NewServer()
		s.RouteToHandler("GET", "/api/json", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
		s.RouteToHandler("GET", "/crumbIssuer/api/json/api/json", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
		s.RouteToHandler("GET", "/pluginManager/api/json", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
		opts := &JenkinsConfigOption{
			URL:       s.URL(),
			Namespace: "test",
			BasicAuth: &BasicAuth{
				Username: "test_user",
				Password: "test_password",
			},
			EnableRestart: false,
		}
		j, err = NewClient(opts)
		Expect(err).ShouldNot(HaveOccurred())
	})
	Context("InstallPluginsIfNotExists method", func() {
		When("jenkins restart enable is false", func() {
			BeforeEach(func() {
				s.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/scriptText"),
						ghttp.RespondWith(http.StatusOK, "ok"),
					),
				)
			})
			It("should return error", func() {
				err := j.InstallPluginsIfNotExists(p)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("installed new plugins need to restart jenkins"))
			})
		})
	})
})
