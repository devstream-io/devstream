package github

import (
	"fmt"
	"net/http"

	// "github.com/devstream-io/devstream/pkg/util/github"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("WorkflowHelper", func() {
	var s *ghttp.Server
	var rightClient, wrongClient *Client
	owner, repo, f, org := "o", "r", ".github/workflows/test", "or"
	u := fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, generateGitHubWorkflowFileByName(f))
	rightOpt := &Option{
		Owner: owner,
		Repo:  repo,
		Org:   org,
	}
	wrongOpt := &Option{
		Owner: owner,
		Repo:  "",
		Org:   org,
	}
	BeforeEach(func() {
		s = ghttp.NewServer()
		rightClient, _ = NewClientWithOption(rightOpt, s.URL())
		Expect(rightClient).NotTo(Equal(nil))
		wrongClient, _ = NewClientWithOption(wrongOpt, s.URL())
		Expect(wrongClient).NotTo(Equal(nil))
	})

	AfterEach(func() {
		//shut down the server between tests
		s.Close()
	})

	Describe("fetching getFileSHA", func() {

		It("error reason is not 404", func() {
			s.Reset()
			s.RouteToHandler("GET", BaseURLPath+u, func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprint(w, "")
			})

			res, err := rightClient.getFileSHA(f)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err.Error()).NotTo(ContainSubstring("404"))
			Expect(res).To(Equal(""))
		})

		It("error reason is 404", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusNotFound)
			res, err := rightClient.getFileSHA(f)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
			Expect(res).To(Equal(""))
		})

		It("no error occurred", func() {
			s.Reset()
			s.RouteToHandler("GET", BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{
					"sha": "21212"
				  }`)
			})
			res, err := rightClient.getFileSHA(f)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err).To(Succeed())
			Expect(res).To(Equal("21212"))
		})

		It("finally got unexpected error", func() {
			s.Reset()
			s.RouteToHandler("GET", BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusAlreadyReported)
				fmt.Fprint(w, `{}`)
			})
			res, err := rightClient.getFileSHA(f)
			Expect(s.ReceivedRequests()).Should(HaveLen(1))
			Expect(err.Error()).To(ContainSubstring("unexpected error"))
			Expect(res).To(Equal(""))
		})
	})
})
