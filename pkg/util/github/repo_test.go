package github_test

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v42/github"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	gh "github.com/devstream-io/devstream/pkg/util/github"
)

var _ = Describe("Repo", func() {
	var s *ghttp.Server
	var rightClient, wrongClient *gh.Client
	var owner, repo, org = "o", "r", "or"
	// var rep *go_github.Repository
	defaultBranch := "db"
	mergeBranch := "mb"
	mainBranch := "mab"
	filePath := ".placeholder"
	rightOpt := &gh.Option{
		Owner: owner,
		Repo:  repo,
		Org:   org,
	}
	wrongOpt := &gh.Option{
		Owner: owner,
		Repo:  "",
		Org:   org,
	}
	BeforeEach(func() {
		s = ghttp.NewServer()
		rightClient, _ = gh.NewClientWithOption(rightOpt, s.URL())
		Expect(rightClient).NotTo(Equal(nil))
		wrongClient, _ = gh.NewClientWithOption(wrongOpt, s.URL())
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
			r, err := rightClient.GetRepoDescription()
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
			var wantR *github.Repository
			Expect(r).To(Equal(wantR))
		})
		It("GetRepoDescription with no error and status 200", func() {
			u := fmt.Sprintf("/repos/%v/%v", org, repo)
			s.Reset()
			s.RouteToHandler("GET", gh.BaseURLPath+u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, ``)
			})
			r, err := rightClient.GetRepoDescription()
			Expect(err).To(Succeed())
			var wantR *github.Repository = &github.Repository{}
			Expect(r).To(Equal(wantR))
		})
	})

	Context("PushLocalPathToBranch", func() {
		BeforeEach(func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
		})

		It("1. create new branch from main", func() {
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			r, err := rightClient.PushLocalPathToBranch(mergeBranch, mainBranch, filePath)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
			Expect(r).To(Equal(false))
		})
		It("2. create new branch from main", func() {
			// u := fmt.Sprintf("/repos/%v/%v/git/ref/heads/%s", org, repo, filePath)
			u := fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, strings.Trim(os.TempDir(), "/"))
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			s.RouteToHandler("GET", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "")
			})
			r, err := rightClient.PushLocalPathToBranch(mergeBranch, mainBranch, filePath)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
			Expect(r).To(Equal(false))
		})
	})

	Context("InitRepo", func() {
		It("CreateRepo with status 500", func() {
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			err := rightClient.InitRepo(mainBranch)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
		})
		It("CreateFile with status 500", func() {
			u := gh.BaseURLPath + fmt.Sprintf("/orgs/%v/repos", org)
			s.Reset()
			s.SetAllowUnhandledRequests(true)
			s.SetUnhandledRequestStatusCode(http.StatusInternalServerError)
			s.RouteToHandler("POST", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{}`)
			})
			err := rightClient.InitRepo(mainBranch)
			fmt.Println(err)
			Expect(err).NotTo(Succeed())
			Expect(err.Error()).To(ContainSubstring(strconv.Itoa(http.StatusInternalServerError)))
		})
		It("CreateFile with status 200", func() {
			u := gh.BaseURLPath + fmt.Sprintf("/orgs/%v/repos", org)
			u2 := gh.BaseURLPath + fmt.Sprintf("/repos/%s/%s/contents/%s", org, repo, ".placeholder")
			s.Reset()
			s.RouteToHandler("POST", u, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{}`)
			})
			s.RouteToHandler("PUT", u2, func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `{}`)
			})
			err := rightClient.InitRepo(mainBranch)
			Expect(err).To(Succeed())
		})
	})
})
