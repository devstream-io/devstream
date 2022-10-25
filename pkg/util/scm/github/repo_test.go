package github_test

import (
	"fmt"
	"net/http"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

var _ = Describe("Repo", func() {
	var (
		s                        *ghttp.Server
		rightClient, wrongClient *github.Client
		owner, repoName, org     = "o", "r", "or"
	)
	// var rep *go_github.Repository
	defaultBranch := "db"
	mainBranch := "mab"
	rightOpt := &git.RepoInfo{
		Owner:  owner,
		Repo:   repoName,
		Org:    org,
		Branch: mainBranch,
	}
	wrongOpt := &git.RepoInfo{
		Owner: owner,
		Repo:  "",
		Org:   org,
	}
	BeforeEach(func() {
		s = ghttp.NewServer()
		rightClient, _ = github.NewClientWithOption(rightOpt, s.URL())
		Expect(rightClient).NotTo(Equal(nil))
		wrongClient, _ = github.NewClientWithOption(wrongOpt, s.URL())
		Expect(wrongClient).NotTo(Equal(nil))
	})

	AfterEach(func() {
		s.Close()
	})

	Context("CreateRepo", func() {
		BeforeEach(func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
		})
		It("create with status 500", func() {
			err := wrongClient.CreateRepo(org, defaultBranch)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
		})
		It("create with status 200", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusOK)
			err := wrongClient.CreateRepo(org, defaultBranch)
			Expect(err).To(Succeed())
		})
	})

	Context("DeleteRepo", func() {
		BeforeEach(func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
		})
		It("DeleteRepo with status 500", func() {
			err := rightClient.DeleteRepo()
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
		})
		It("DeleteRepo with status 404", func() {
			s.SetUnhandledRequestStatusCode(http.StatusNotFound)
			err := wrongClient.DeleteRepo()
			Expect(err).To(Succeed())
		})
		It("DeleteRepo with status 200", func() {
			s.SetUnhandledRequestStatusCode(http.StatusOK)
			err := wrongClient.DeleteRepo()
			Expect(err).To(Succeed())
		})
	})

	Context("GetRepoDescription", func() {
		BeforeEach(func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
		})

		It("GetRepoDescription with status 500", func() {
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			r, err := rightClient.DescribeRepo()
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
			var wantR *git.RepoInfo
			Expect(r).To(Equal(wantR))
		})
		It("GetRepoDescription with no error and status 200", func() {
			u := fmt.Sprintf("/repos/%v/%v", org, repoName)
			s.Reset()
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, ``)
			})
			r, err := rightClient.DescribeRepo()
			Expect(err).To(Succeed())
			wantR := &git.RepoInfo{
				RepoType: "github",
			}
			Expect(r).To(Equal(wantR))
		})
	})

	Context("InitRepo", func() {
		It("CreateRepo with status 500", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			err := rightClient.InitRepo()
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
		})
		It("CreateFile with status 500", func() {
			u := github.BaseURLPath + fmt.Sprintf("/orgs/%v/repos", org)
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			s.RouteToHandler("POST", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{}`)
			})
			err := rightClient.InitRepo()
			fmt.Println(err)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
		})
		It("CreateFile with status 200", func() {
			u := github.BaseURLPath + fmt.Sprintf("/orgs/%v/repos", org)
			u2 := github.BaseURLPath + fmt.Sprintf("/repos/%s/%s/contents/%s", org, repoName, ".placeholder")
			s.Reset()
			s.RouteToHandler("POST", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{}`)
			})
			s.RouteToHandler("PUT", u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{}`)
			})
			err := rightClient.InitRepo()
			Expect(err).To(Succeed())
		})
	})

	Context("ProtectBranch", func() {
		BeforeEach(func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
		})
		It("ProtectBranch with status 500", func() {
			u := fmt.Sprintf("/repos/%v/%v", org, repoName)
			s.Reset()
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			s.SetAllowUnhandledRequests(true)
			s.RouteToHandler("GET", github.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, ``)
			})
			err := rightClient.ProtectBranch(mainBranch)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
		})
		It("ProtectBranch with status 200", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusOK)
			err := rightClient.ProtectBranch(mainBranch)
			Expect(err).To(Succeed())
		})
	})
})
